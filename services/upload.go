package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/minio/minio-go/v7"
	"github.com/project-box/models"
	"github.com/project-box/repositories"
	"github.com/xuri/excelize/v2"
)

type UploadService interface {
	UploadObject(ctx context.Context, bucketName string, objectName string, file io.Reader, fileSize int64, contentType string) (string, error)
	GetObjectURL(ctx context.Context, bucketName string, objectName string) (*url.URL, error)
	RemoveObject(ctx context.Context, bucketName string, objectName string, opt minio.RemoveObjectOptions) error
	ProcessStudentEnrollmentFile(ctx context.Context, programId int, file *multipart.FileHeader) error
	ProcessCreateProjectFile(ctx context.Context, programId int, file *multipart.FileHeader) error
	ProcessCreateStaffFile(ctx context.Context, programId int, file *multipart.FileHeader) error
}

type uploadServiceImpl struct {
	client             *minio.Client
	studentService     StudentService
	staffService       StaffService
	projectRoleService ProjectRoleService
	configService      ConfigService
	projectService     ProjectService
	projectRepo        repositories.ProjectRepository
	programRepo        repositories.ProgramRepository
	keywordRepo        repositories.KeywordRepository
}

func NewUploadService(client *minio.Client, keywordRepo repositories.KeywordRepository, programRepo repositories.ProgramRepository, projectRepo repositories.ProjectRepository, staffService StaffService, projectRoleService ProjectRoleService, projectService ProjectService, configService ConfigService, studentService StudentService) UploadService {
	return &uploadServiceImpl{
		client:             client,
		programRepo:        programRepo,
		projectRepo:        projectRepo,
		staffService:       staffService,
		projectRoleService: projectRoleService,
		projectService:     projectService,
		configService:      configService,
		studentService:     studentService,
		keywordRepo:        keywordRepo,
	}
}

func matchColumn(col, target string) bool {
	return strings.EqualFold(col, target)
}

func (s *uploadServiceImpl) validateFileType(file *multipart.FileHeader) error {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".xlsx" && ext != ".xls" {
		return errors.New("invalid file type: only Excel files are allowed")
	}
	return nil
}

func (s *uploadServiceImpl) readExcelFile(file *multipart.FileHeader) ([][]string, error) {
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {

		}
	}(src)

	excelFile, err := excelize.OpenReader(src)
	if err != nil {
		return nil, fmt.Errorf("failed to read Excel file: %w", err)
	}

	sheetName := excelFile.GetSheetName(0)
	if sheetName == "" {
		return nil, errors.New("no sheets found in the Excel file")
	}

	rows, err := excelFile.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to get rows: %w", err)
	}

	return rows, nil
}

func (s *uploadServiceImpl) ProcessStudentEnrollmentFile(ctx context.Context, programId int, file *multipart.FileHeader) error {
	if err := s.validateFileType(file); err != nil {
		return err
	}

	rows, err := s.readExcelFile(file)
	if err != nil {
		return err
	}

	if len(rows) < 4 {
		return errors.New("excel file is empty or does not have enough rows")
	}

	studentInfoColumns, err := s.getStudentInfoColumns(rows[3])
	if err != nil {
		return err
	}

	students, err := s.parseStudents(ctx, rows[4:], studentInfoColumns, programId)
	if err != nil {
		return err
	}
	if _, err := s.studentService.UpsertStudents(ctx, students, programId); err != nil {
		return fmt.Errorf("failed to save student data: %w", err)
	}

	return nil
}

func (s *uploadServiceImpl) ProcessCreateProjectFile(ctx context.Context, programId int, file *multipart.FileHeader) error {
	if err := s.validateFileType(file); err != nil {
		return err
	}

	rows, err := s.readExcelFile(file)
	if err != nil {
		return err
	}

	if len(rows) < 3 {
		return errors.New("excel file is empty or does not have enough rows")
	}

	createProjectInfoColumns, err := s.getCreateProjectInfoColumns(rows[2])
	if err != nil {
		return err
	}
	projectRequests, err := s.parseProjects(ctx, rows[3:], createProjectInfoColumns, programId)
	if err != nil {
		return err
	}
	err = s.projectService.CreateProjects(ctx, projectRequests)
	if err != nil {
		return err
	}

	return nil
}

func (s *uploadServiceImpl) ProcessCreateStaffFile(ctx context.Context, programId int, file *multipart.FileHeader) error {
	if err := s.validateFileType(file); err != nil {
		return err
	}

	rows, err := s.readExcelFile(file)
	if err != nil {
		return err
	}

	if len(rows) < 3 {
		return errors.New("excel file is empty or does not have enough rows")
	}

	staffInfoColumns, err := s.getStaffInfoColumns(rows[1])
	if err != nil {
		return err
	}

	staffs, err := s.parseStaffs(rows[2:], staffInfoColumns, programId)
	if err != nil {
		return err
	}

	if err := s.staffService.CreateStaffs(ctx, staffs); err != nil {
		return fmt.Errorf("failed to save staff data: %w", err)
	}

	return nil
}

func (s *uploadServiceImpl) getStaffInfoColumns(headerRow []string) (map[string]int, error) {
	columns := map[string]int{
		"staffPrefixTHColumn": -1,
		"staffPrefixENColumn": -1,
		"staffNameTHColumn":   -1,
		"staffNameENColumn":   -1,
		"staffEmailColumn":    -1,
		"staffIsActiveColumn": -1,
	}

	for j, col := range headerRow {
		switch {
		case matchColumn(col, "Prefix (TH)"):
			columns["staffPrefixTHColumn"] = j
		case matchColumn(col, "Prefix (EN)"):
			columns["staffPrefixENColumn"] = j
		case matchColumn(col, "ชื่อ-นามสกุล (TH)"):
			columns["staffNameTHColumn"] = j
		case matchColumn(col, "ชื่อ-นามสกุล (EN)"):
			columns["staffNameENColumn"] = j
		case matchColumn(col, "Email (required)"):
			columns["staffEmailColumn"] = j
		case matchColumn(col, "InActive"):
			columns["staffIsActiveColumn"] = j
		}
	}

	if columns["staffPrefixTHColumn"] == -1 || columns["staffPrefixENColumn"] == -1 || columns["staffNameTHColumn"] == -1 || columns["staffNameENColumn"] == -1 || columns["staffEmailColumn"] == -1 || columns["staffIsActiveColumn"] == -1 {
		return nil, errors.New("missing one or more required columns in the header row")
	}

	return columns, nil
}

func (s *uploadServiceImpl) getStudentInfoColumns(headerRow []string) (map[string]int, error) {
	columns := map[string]int{
		"secLecColumn":      -1,
		"secLabColumn":      -1,
		"studentIdColumn":   -1,
		"studentNameColumn": -1,
		"cmuAccountColumn":  8,
	}

	for j, col := range headerRow {
		switch {
		case matchColumn(col, "SECLEC"):
			columns["secLecColumn"] = j
		case matchColumn(col, "SECLAB"):
			columns["secLabColumn"] = j
		case matchColumn(col, "รหัสนักศึกษา"):
			columns["studentIdColumn"] = j
		case matchColumn(col, "ชื่อ - นามสกุล"):
			columns["studentNameColumn"] = j
		}
	}

	if columns["secLecColumn"] == -1 || columns["secLabColumn"] == -1 || columns["studentIdColumn"] == -1 || columns["studentNameColumn"] == -1 {
		return nil, errors.New("missing one or more required columns in the header row")
	}

	return columns, nil
}

func (s *uploadServiceImpl) getCreateProjectInfoColumns(headerRow []string) (map[string]int, error) {
	columns := map[string]int{
		"titleTHColumn":      -1,
		"titleENColumn":      -1,
		"abstractTextColumn": -1,
		"secLabColumn":       -1,
		"studentIdColumn":    -1,
		"studentNameColumn":  -1,
		"staffColumn":        -1,
		"staffRoleColumn":    -1,
	}

	for j, col := range headerRow {
		switch {
		case matchColumn(col, "Title (TH)"):
			columns["titleTHColumn"] = j
		case matchColumn(col, "Title (EN)"):
			columns["titleENColumn"] = j
		case matchColumn(col, "Abstract"):
			columns["abstractTextColumn"] = j
		case matchColumn(col, "SECLAB"):
			columns["secLabColumn"] = j
		case matchColumn(col, "Student(s)"):
			columns["studentIdColumn"] = j
		case matchColumn(col, "ชื่อ - นามสกุล"):
			columns["studentNameColumn"] = j
		case matchColumn(col, "Committee(s)"):
			columns["staffColumn"] = j
		case matchColumn(col, "Staff Role"):
			columns["staffRoleColumn"] = j
		case matchColumn(col, "Keywords"):
			columns["keywordsColumn"] = j
		case matchColumn(col, "Is Public"):
			columns["isPublicColumn"] = j
		}

	}

	if columns["titleTHColumn"] == -1 || columns["titleENColumn"] == -1 || columns["abstractTextColumn"] == -1 || columns["secLabColumn"] == -1 || columns["studentIdColumn"] == -1 || columns["studentNameColumn"] == -1 || columns["staffColumn"] == -1 || columns["staffRoleColumn"] == -1 || columns["keywordsColumn"] == -1 || columns["isPublicColumn"] == -1 {
		return nil, errors.New("missing one or more required columns in the header row")
	}

	return columns, nil
}

func (s *uploadServiceImpl) parseStudents(ctx context.Context, rows [][]string, columns map[string]int, programId int) ([]models.Student, error) {
	academicYear, semester, err := s.configService.GetCurrentAcademicYearAndSemester(ctx, programId)
	if err != nil {
		return nil, err
	}

	var students []models.Student
	for _, row := range rows {
		if len(row) <= columns["studentIdColumn"] || len(row) <= columns["secLabColumn"] || len(row) <= columns["studentNameColumn"]+1 {
			return nil, fmt.Errorf("row does not have enough columns: %v", row)
		}

		var email *string
		if len(row) > columns["cmuAccountColumn"] && row[columns["cmuAccountColumn"]] != "" {
			email = &row[columns["cmuAccountColumn"]]
		}
		student := models.Student{
			StudentID:    row[columns["studentIdColumn"]],
			SecLab:       row[columns["secLabColumn"]],
			FirstName:    row[columns["studentNameColumn"]],
			LastName:     row[columns["studentNameColumn"]+1],
			Email:        email,
			Semester:     semester,
			AcademicYear: academicYear,
			ProgramID:    programId,
		}
		students = append(students, student)
	}

	return students, nil
}

func (s *uploadServiceImpl) parseStaffs(rows [][]string, columns map[string]int, programId int) ([]models.Staff, error) {
	var staffs []models.Staff
	for _, row := range rows {
		if len(row) <= columns["staffPrefixTHColumn"] || len(row) <= columns["staffPrefixENColumn"] || len(row) <= columns["staffNameTHColumn"] || len(row) <= columns["staffNameTHColumn"]+1 || len(row) <= columns["staffNameENColumn"] || len(row) <= columns["staffNameENColumn"]+1 || len(row) <= columns["staffEmailColumn"] || len(row) <= columns["staffIsActiveColumn"] {
			continue
		}

		if row[columns["staffEmailColumn"]] == "" {
			return nil, fmt.Errorf("row does not have staff email: %v", row)
		}

		isActive, err := strconv.ParseBool(row[columns["staffIsActiveColumn"]])
		if err != nil {
			return nil, fmt.Errorf("invalid value for Is Resigned: %v", row[columns["staffIsActiveColumn"]])
		}

		staff := models.Staff{
			PrefixTH:    row[columns["staffPrefixTHColumn"]],
			PrefixEN:    row[columns["staffPrefixENColumn"]],
			FirstNameTH: row[columns["staffNameTHColumn"]],
			LastNameTH:  row[columns["staffNameTHColumn"]+1],
			FirstNameEN: row[columns["staffNameENColumn"]],
			LastNameEN:  row[columns["staffNameENColumn"]+1],
			Email:       row[columns["staffEmailColumn"]],
			IsActive:    isActive,
			ProgramID:   programId,
		}
		staffs = append(staffs, staff)
	}

	return staffs, nil
}

func (s *uploadServiceImpl) parseProjects(ctx context.Context, rows [][]string, columns map[string]int, programId int) ([]models.ProjectRequest, error) {
	academicYear, semester, err := s.configService.GetCurrentAcademicYearAndSemester(ctx, programId)
	if err != nil {
		return nil, err
	}

	var projectRequests []models.ProjectRequest
	for rowIdx, row := range rows {
		if !s.isValidProjectRow(ctx, row, columns) {
			continue
		}

		isProjectDuplicate, err := s.projectRepo.CheckDuplicateProjectByTitleAndSemester(ctx, row[columns["titleTHColumn"]], row[columns["titleENColumn"]], academicYear, semester)
		if err != nil {
			return nil, err
		}
		if isProjectDuplicate {
			return nil, fmt.Errorf("project with title TH: %s and title EN: %s already exists", row[columns["titleTHColumn"]], row[columns["titleENColumn"]])
		}

		members, isSecLabSameOverAllMember, err := s.getProjectMembers(rows, rowIdx, columns, semester, academicYear, programId)
		if err != nil {
			return nil, err
		}

		members, err = s.studentService.UpsertStudents(ctx, members, programId)
		if err != nil {
			return nil, err
		}

		projectStaffs, err := s.getProjectStaffs(ctx, rows, rowIdx, columns)
		if err != nil {
			return nil, err
		}

		var keywordArray []models.Keyword
		keywords := strings.Split(row[columns["keywordsColumn"]], ",")
		for _, keyword := range keywords {
			keyword = strings.TrimSpace(keyword)
			if keyword == "" {
				continue
			}

			keywordModel, err := s.keywordRepo.FindByKeywordAndProgramId(ctx, keyword, programId)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
			if errors.Is(err, gorm.ErrRecordNotFound) {
				keywordModel = &models.Keyword{
					Keyword:   keyword,
					ProgramID: programId,
				}
				if err := s.keywordRepo.Create(ctx, keywordModel); err != nil {
					return nil, err
				}
			}
			keywordArray = append(keywordArray, *keywordModel)
		}

		titleTH := &row[columns["titleTHColumn"]]
		titleEN := &row[columns["titleENColumn"]]
		abstractText := &row[columns["abstractTextColumn"]]
		sectionValue := row[columns["secLabColumn"]]
		isPublicValue := row[columns["isPublicColumn"]]

		var sectionID *string
		if isSecLabSameOverAllMember {
			sectionID = &sectionValue
		}
		var isPublic bool
		if isPublicValue == "TRUE" {
			isPublic = true
		} else {
			isPublic = false
		}

		project := models.ProjectRequest{
			TitleTH:       titleTH,
			TitleEN:       titleEN,
			AbstractText:  abstractText,
			AcademicYear:  academicYear,
			Semester:      semester,
			SectionID:     sectionID,
			ProgramID:     programId,
			IsPublic:      isPublic,
			ProjectStaffs: projectStaffs,
			Members:       members,
			Keywords:      keywordArray,
		}
		projectRequests = append(projectRequests, project)
	}

	return projectRequests, nil
}

func (s *uploadServiceImpl) isValidProjectRow(ctx context.Context, row []string, columns map[string]int) bool {
	return len(row) > columns["titleTHColumn"] &&
		len(row) > columns["titleENColumn"] &&
		len(row) > columns["abstractTextColumn"] &&
		len(row) > columns["courseNoColumn"] &&
		!(row[columns["titleTHColumn"]] == "" &&
			row[columns["titleENColumn"]] == "" &&
			row[columns["abstractTextColumn"]] == "" &&
			row[columns["courseNoColumn"]] == "")
}

func (s *uploadServiceImpl) getProjectMembers(rows [][]string, rowIdx int, columns map[string]int, semester, academicYear, programId int) ([]models.Student, bool, error) {
	var members []models.Student
	memberIdx := 0
	currentSeclab := ""
	isSecLabSameOverAllMember := true
	for rowIdx+memberIdx < len(rows) &&
		columns["studentIdColumn"] < len(rows[rowIdx+memberIdx]) &&
		len(rows[rowIdx+memberIdx]) > columns["studentIdColumn"] &&
		rows[rowIdx+memberIdx][columns["studentIdColumn"]] != "" {

		if currentSeclab != "" && currentSeclab != rows[rowIdx+memberIdx][columns["secLabColumn"]] {
			isSecLabSameOverAllMember = false
		}
		currentSeclab = rows[rowIdx+memberIdx][columns["secLabColumn"]]

		student := models.Student{
			StudentID:    rows[rowIdx+memberIdx][columns["studentIdColumn"]],
			SecLab:       rows[rowIdx+memberIdx][columns["secLabColumn"]],
			FirstName:    rows[rowIdx+memberIdx][columns["studentNameColumn"]],
			LastName:     rows[rowIdx+memberIdx][columns["studentNameColumn"]+1],
			Semester:     semester,
			AcademicYear: academicYear,
			ProgramID:    programId,
		}
		members = append(members, student)
		memberIdx++
	}
	return members, isSecLabSameOverAllMember, nil
}

func (s *uploadServiceImpl) getProjectStaffs(ctx context.Context, rows [][]string, rowIdx int, columns map[string]int) ([]models.ProjectStaff, error) {
	var projectStaffs []models.ProjectStaff
	staffIdx := 0
	for rowIdx+staffIdx < len(rows) &&
		columns["staffColumn"] < len(rows[rowIdx+staffIdx]) &&
		len(rows[rowIdx+staffIdx]) > columns["staffColumn"] &&
		rows[rowIdx+staffIdx][columns["staffColumn"]] != "" {

		staff, err := s.staffService.GetStaffByName(ctx, rows[rowIdx+staffIdx][columns["staffColumn"]])
		if err != nil {
			return nil, err
		}

		if len(rows[rowIdx+staffIdx]) <= columns["staffRoleColumn"] {
			return nil, errors.New("staff role column index out of range")
		}
		staffRoleTH := rows[rowIdx+staffIdx][columns["staffRoleColumn"]]
		staffRole, err := s.projectRoleService.GetProjectRoleByRoleName(ctx, staffRoleTH)
		if err != nil {
			return nil, err
		}

		projectStaff := models.ProjectStaff{
			StaffID:       staff.ID,
			ProjectRoleID: staffRole.ID,
		}
		projectStaffs = append(projectStaffs, projectStaff)
		staffIdx++
	}
	return projectStaffs, nil
}

func (s *uploadServiceImpl) UploadObject(ctx context.Context, bucketName string, objectName string, file io.Reader, fileSize int64, contentType string) (string, error) {
	_, err := s.client.PutObject(ctx, bucketName, objectName, file, fileSize, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", fmt.Errorf("failed to upload object: %w", err)
	}

	url, err := s.client.PresignedGetObject(ctx, bucketName, objectName, time.Hour*24*7, nil)
	if err != nil {
		_ = s.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url.String(), nil
}

func (s *uploadServiceImpl) RemoveObject(ctx context.Context, bucketName string, objectName string, opt minio.RemoveObjectOptions) error {
	return s.client.RemoveObject(ctx, bucketName, objectName, opt)
}

func (s *uploadServiceImpl) GetObjectURL(ctx context.Context, bucketName string, objectName string) (*url.URL, error) {
	return s.client.PresignedGetObject(ctx, bucketName, objectName, time.Hour*2, nil)
}
