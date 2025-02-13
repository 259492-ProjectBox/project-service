package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"strings"
	"time"

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
}

type uploadServiceImpl struct {
	client             *minio.Client
	courseService      CourseService
	studentService     StudentService
	staffService       StaffService
	projectRoleService ProjectRoleService
	configService      ConfigService
	projectService     ProjectService
	programRepo        repositories.ProgramRepository
}

func NewUploadService(client *minio.Client, programRepo repositories.ProgramRepository, staffService StaffService, projectRoleService ProjectRoleService, projectService ProjectService, courseService CourseService, configService ConfigService, studentService StudentService) UploadService {
	return &uploadServiceImpl{
		client:             client,
		programRepo:        programRepo,
		staffService:       staffService,
		projectRoleService: projectRoleService,
		projectService:     projectService,
		courseService:      courseService,
		configService:      configService,
		studentService:     studentService,
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
	defer src.Close()

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

	if err := s.studentService.UpsertStudents(ctx, students); err != nil {
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

	createProjectInfoColumns, err := s.getCreateProjectInfoColumns(rows[1])
	if err != nil {
		return err
	}

	projectRequests, err := s.parseProjects(ctx, rows[2:], createProjectInfoColumns, programId)
	if err != nil {
		return err
	}

	err = s.projectService.CreateProjects(ctx, projectRequests)
	if err != nil {
		return err
	}

	return nil
}

func (s *uploadServiceImpl) getCourseNo(rows [][]string) (string, error) {
	courseRow := rows[1]
	for j, col := range courseRow {
		if matchColumn(col, "COURSE NO :") {
			return rows[1][j+2], nil
		}
	}
	return "", errors.New("course number not found")
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
		"sectionColumn":      -1,
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
		case matchColumn(col, "Section"):
			columns["sectionColumn"] = j
		case matchColumn(col, "Student(s)"):
			columns["studentIdColumn"] = j
		case matchColumn(col, "ชื่อ - นามสกุล"):
			columns["studentNameColumn"] = j
		case matchColumn(col, "Committee(s)"):
			columns["staffColumn"] = j
		case matchColumn(col, "Staff Role"):
			columns["staffRoleColumn"] = j
		}
	}

	if columns["titleTHColumn"] == -1 || columns["titleENColumn"] == -1 || columns["abstractTextColumn"] == -1 || columns["sectionColumn"] == -1 || columns["studentIdColumn"] == -1 || columns["studentNameColumn"] == -1 || columns["staffColumn"] == -1 || columns["staffRoleColumn"] == -1 {
		return nil, errors.New("missing one or more required columns in the header row")
	}

	return columns, nil
}

func (s *uploadServiceImpl) parseStudents(ctx context.Context, rows [][]string, columns map[string]int, programId int) ([]models.Student, error) {
	academicYear, semester, err := s.configService.GetAcademicYearAndSemester(ctx, programId)
	if err != nil {
		return nil, err
	}

	courseNo, err := s.getCourseNo(rows)
	if err != nil {
		return nil, err
	}

	course, err := s.courseService.FindCourseByCourseNo(ctx, courseNo)
	if err != nil {
		return nil, err
	}

	var students []models.Student
	for _, row := range rows {
		student := models.Student{
			StudentID:    row[columns["studentIdColumn"]],
			SecLab:       row[columns["secLabColumn"]],
			FirstName:    row[columns["studentNameColumn"]],
			LastName:     row[columns["studentNameColumn"]+1],
			Email:        row[columns["cmuAccountColumn"]],
			Semester:     academicYear,
			AcademicYear: semester,
			CourseID:     course.ID,
			ProgramID:    programId,
		}
		students = append(students, student)
	}

	return students, nil
}

func (s *uploadServiceImpl) parseProjects(ctx context.Context, rows [][]string, columns map[string]int, programId int) ([]models.ProjectRequest, error) {
	academicYear, semester, err := s.configService.GetAcademicYearAndSemester(ctx, programId)
	if err != nil {
		return nil, err
	}

	var projectRequests []models.ProjectRequest
	for rowIdx, row := range rows {
		if !s.isValidProjectRow(row, columns) {
			continue
		}

		courseNo, err := s.getCourseNo(rows)
		if err != nil {
			return nil, err
		}

		course, err := s.courseService.FindCourseByCourseNo(ctx, courseNo)
		if err != nil {
			return nil, err
		}

		members, err := s.getProjectMembers(rows, rowIdx, columns, semester, academicYear, course.ID, programId)
		if err != nil {
			return nil, err
		}

		s.studentService.UpsertStudents(ctx, members)

		projectStaffs, err := s.getProjectStaffs(ctx, rows, rowIdx, columns)
		if err != nil {
			return nil, err
		}

		project := models.ProjectRequest{
			TitleTH:       &row[columns["titleTHColumn"]],
			TitleEN:       &row[columns["titleENColumn"]],
			AbstractText:  &row[columns["abstractTextColumn"]],
			AcademicYear:  academicYear,
			Semester:      semester,
			SectionID:     &row[columns["sectionColumn"]],
			CourseID:      course.ID,
			ProgramID:     programId,
			ProjectStaffs: projectStaffs,
			Members:       members,
		}
		projectRequests = append(projectRequests, project)
	}

	return projectRequests, nil
}

func (s *uploadServiceImpl) isValidProjectRow(row []string, columns map[string]int) bool {
	return len(row) > columns["titleTHColumn"] &&
		len(row) > columns["titleENColumn"] &&
		len(row) > columns["abstractTextColumn"] &&
		len(row) > columns["sectionColumn"] &&
		len(row) > columns["courseNoColumn"] &&
		!(row[columns["titleTHColumn"]] == "" &&
			row[columns["titleENColumn"]] == "" &&
			row[columns["abstractTextColumn"]] == "" &&
			row[columns["sectionColumn"]] == "" &&
			row[columns["courseNoColumn"]] == "")
}

func (s *uploadServiceImpl) getProjectMembers(rows [][]string, rowIdx int, columns map[string]int, semester, academicYear, courseId, programId int) ([]models.Student, error) {
	var members []models.Student
	memberIdx := 0
	for rowIdx+memberIdx < len(rows) &&
		columns["studentIdColumn"] < len(rows[rowIdx+memberIdx]) &&
		len(rows[rowIdx+memberIdx]) > columns["studentIdColumn"] &&
		rows[rowIdx+memberIdx][columns["studentIdColumn"]] != "" {

		student := models.Student{
			StudentID:    rows[rowIdx+memberIdx][columns["studentIdColumn"]],
			SecLab:       rows[rowIdx+memberIdx][columns["secLabColumn"]],
			FirstName:    rows[rowIdx+memberIdx][columns["studentNameColumn"]],
			LastName:     rows[rowIdx+memberIdx][columns["studentNameColumn"]+1],
			Semester:     semester,
			AcademicYear: academicYear,
			CourseID:     courseId,
			ProgramID:    programId,
		}
		members = append(members, student)
		memberIdx++
	}
	return members, nil
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
