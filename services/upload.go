package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"strconv"

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
	studentService StudentService
	configService  ConfigService
}

func NewUploadService(client *minio.Client, configService ConfigService, studentService StudentService) UploadService {
	return &uploadServiceImpl{
		client:         client,
		configService:  configService,
		studentService: studentService,
	}
}

func matchColumn(col, target string) bool {
	return strings.EqualFold(col, target)
}

func (s *uploadServiceImpl) ProcessStudentEnrollmentFile(ctx context.Context, programId int, file *multipart.FileHeader) error {
	academicYear, err := s.configService.FindConfigByNameAndProgramId(ctx, "academic year", programId)
	if err != nil {
		return fmt.Errorf("failed to fetch academic year: %w", err)
	}

	academicYearInt, err := strconv.Atoi(academicYear.Value)
	if err != nil {
		return fmt.Errorf("failed to convert academic year: %w", err)
	}

	semester, err := s.configService.FindConfigByNameAndProgramId(ctx, "semester", programId)
	if err != nil {
		return fmt.Errorf("failed to fetch academic year: %w", err)
	}

	semesterInt, err := strconv.Atoi(semester.Value)
	if err != nil {
		return fmt.Errorf("failed to convert academic year: %w", err)
	}

	var secLecColumn, secLabColumn, studentColumn, nameColumn, cmuAccountColumn int
	secLecColumn, secLabColumn, studentColumn, nameColumn, cmuAccountColumn = -1, -1, -1, -1, 8

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

	headerRow := rows[3]
	for j, col := range headerRow {
		switch {
		case matchColumn(col, "SECLEC"):
			secLecColumn = j
		case matchColumn(col, "SECLAB"):
			secLabColumn = j
		case matchColumn(col, "รหัสนักศึกษา"):
			studentColumn = j
		case matchColumn(col, "ชื่อ - นามสกุล"):
			nameColumn = j
		}
	}

	if secLecColumn == -1 || secLabColumn == -1 || studentColumn == -1 || nameColumn == -1 {
		return errors.New("missing one or more required columns in the header row")
	}

	var students []models.Student
	for _, row := range rows[4:] {

		student := models.Student{
			SecLab:       row[secLabColumn],
			StudentID:    row[studentColumn],
			FirstName:    row[nameColumn],
			LastName:     row[nameColumn+1],
			Email:        row[cmuAccountColumn],
			Semester:     semesterInt,
			AcademicYear: academicYearInt,
			ProgramID:    programId,
		}

		students = append(students, student)
	}

	if err := s.studentService.CreateStudents(ctx, students); err != nil {
		return fmt.Errorf("failed to save student data: %w", err)
	}

	return nil
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
