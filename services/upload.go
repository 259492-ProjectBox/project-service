package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/project-box/models"
	"github.com/xuri/excelize/v2"
)

type UploadService interface {
	UploadObject(ctx context.Context, bucketName string, objectName string, file io.Reader, fileSize int64, contentType string) (string, error)
	RemoveObject(ctx context.Context, bucketName string, objectName string, opt minio.RemoveObjectOptions) error
	ProcessStudentEnrollmentFile(ctx context.Context, programId int, file *multipart.FileHeader) error
}

type uploadServiceImpl struct {
	client         *minio.Client
	courseService  CourseService
	studentService StudentService
	configService  ConfigService
}

func NewUploadService(client *minio.Client, courseService CourseService, configService ConfigService, studentService StudentService) UploadService {
	return &uploadServiceImpl{
		client:         client,
		courseService:  courseService,
		configService:  configService,
		studentService: studentService,
	}
}

func matchColumn(col, target string) bool {
	return strings.EqualFold(col, target)
}

func (s *uploadServiceImpl) ProcessStudentEnrollmentFile(ctx context.Context, programId int, file *multipart.FileHeader) error {
	academicYear, semester, err := s.getAcademicYearAndSemester(ctx, programId)
	if err != nil {
		return err
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".xlsx" && ext != ".xls" {
		return errors.New("invalid file type: only Excel files are allowed")
	}

	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	excelFile, err := excelize.OpenReader(src)
	if err != nil {
		return fmt.Errorf("failed to read Excel file: %w", err)
	}

	sheetName := excelFile.GetSheetName(0)
	if sheetName == "" {
		return errors.New("no sheets found in the Excel file")
	}

	rows, err := excelFile.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("failed to get rows: %w", err)
	}

	if len(rows) < 4 {
		return errors.New("excel file is empty or does not have enough rows")
	}

	courseNo, err := s.getCourseNo(rows)
	if err != nil {
		return err
	}

	course, err := s.courseService.FindCourseByCourseNo(ctx, courseNo)
	if err != nil {
		return err
	}

	studentInfoColumns, err := s.getStudentInfoColumns(rows[3])
	if err != nil {
		return err
	}

	students, err := s.parseStudents(rows[4:], studentInfoColumns, course.ID, semester, academicYear, programId)
	if err != nil {
		return err
	}

	if err := s.studentService.CreateStudents(ctx, students); err != nil {
		return fmt.Errorf("failed to save student data: %w", err)
	}

	return nil
}

func (s *uploadServiceImpl) getAcademicYearAndSemester(ctx context.Context, programId int) (int, int, error) {
	academicYear, err := s.configService.FindConfigByNameAndProgramId(ctx, "academic year", programId)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to fetch academic year: %w", err)
	}

	academicYearInt, err := strconv.Atoi(academicYear.Value)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to convert academic year: %w", err)
	}

	semester, err := s.configService.FindConfigByNameAndProgramId(ctx, "semester", programId)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to fetch semester: %w", err)
	}

	semesterInt, err := strconv.Atoi(semester.Value)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to convert semester: %w", err)
	}

	return academicYearInt, semesterInt, nil
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
		"secLecColumn":     -1,
		"secLabColumn":     -1,
		"studentColumn":    -1,
		"nameColumn":       -1,
		"cmuAccountColumn": 8,
	}

	for j, col := range headerRow {
		switch {
		case matchColumn(col, "SECLEC"):
			columns["secLecColumn"] = j
		case matchColumn(col, "SECLAB"):
			columns["secLabColumn"] = j
		case matchColumn(col, "รหัสนักศึกษา"):
			columns["studentColumn"] = j
		case matchColumn(col, "ชื่อ - นามสกุล"):
			columns["nameColumn"] = j
		}
	}

	if columns["secLecColumn"] == -1 || columns["secLabColumn"] == -1 || columns["studentColumn"] == -1 || columns["nameColumn"] == -1 {
		return nil, errors.New("missing one or more required columns in the header row")
	}

	return columns, nil
}

func (s *uploadServiceImpl) parseStudents(rows [][]string, columns map[string]int, courseID, semesterInt, academicYearInt, programId int) ([]models.Student, error) {
	var students []models.Student
	for _, row := range rows {
		student := models.Student{
			SecLab:       row[columns["secLabColumn"]],
			StudentID:    row[columns["studentColumn"]],
			FirstName:    row[columns["nameColumn"]],
			LastName:     row[columns["nameColumn"]+1],
			Email:        row[columns["cmuAccountColumn"]],
			CourseID:     courseID,
			Semester:     semesterInt,
			AcademicYear: academicYearInt,
			ProgramID:    programId,
		}
		students = append(students, student)
	}
	return students, nil
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
