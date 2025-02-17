package services

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
	"github.com/project-box/dtos"
	"github.com/project-box/models"
	"github.com/project-box/repositories"
)

type ProjectResourceConfigService interface {
	GetProjectResourceConfigsByProgramId(ctx context.Context, programID int) ([]dtos.ProjectResourceConfig, error)
	UpsertResourceProjectConfig(config *models.ProjectResourceConfig) error
	UpsertResourceProjectConfigV2(ctx context.Context, req dtos.CreateProjectResourceConfigRequest) error
}

type projectResourceConfigServiceImpl struct {
	uploadService             UploadService
	programService            ProgramService
	projectResourceConfigRepo repositories.ProjectResourceConfigRepository
}

func NewProjectResourceConfigService(projectResourceConfigRepository repositories.ProjectResourceConfigRepository, programService ProgramService, uploadService UploadService) ProjectResourceConfigService {
	return &projectResourceConfigServiceImpl{
		programService:            programService,
		uploadService:             uploadService,
		projectResourceConfigRepo: projectResourceConfigRepository,
	}
}

func (s *projectResourceConfigServiceImpl) GetProjectResourceConfigsByProgramId(ctx context.Context, programId int) ([]dtos.ProjectResourceConfig, error) {
	projectResourcesConfigs, err := s.projectResourceConfigRepo.GetProjectResourceConfigsByProgramId(programId)
	if err != nil {
		return nil, err
	}

	var projectResourceConfigResponses []dtos.ProjectResourceConfig
	for _, config := range projectResourcesConfigs {
		projectResourceConfigResponse := dtos.ProjectResourceConfig{
			ID:             config.ID,
			Title:          config.Title,
			IconName:       config.IconName,
			IsActive:       config.IsActive,
			ResourceTypeID: config.ResourceTypeID,
			ResourceType: models.ResourceType{
				ID:       config.ResourceType.ID,
				TypeName: config.ResourceType.TypeName,
			},
			ProgramID: config.ProgramID,
			Program:   models.Program{ID: config.Program.ID, ProgramNameTH: config.Program.ProgramNameTH, ProgramNameEN: config.Program.ProgramNameEN, Abbreviation: config.Program.Abbreviation},
		}
		if config.IconName != nil {
			objectName := fmt.Sprintf("%s/%s", config.Program.ProgramNameTH, *config.IconName)
			url, err := s.uploadService.GetObjectURL(ctx, "icons", objectName)
			if err != nil {
				return nil, err
			}
			projectResourceConfigResponse.IconURL = url.String()
		}
		projectResourceConfigResponses = append(projectResourceConfigResponses, projectResourceConfigResponse)
	}

	return projectResourceConfigResponses, nil
}

func (s *projectResourceConfigServiceImpl) UpsertResourceProjectConfig(config *models.ProjectResourceConfig) error {
	return s.projectResourceConfigRepo.UpsertResourceProjectConfig(config)
}

func (s *projectResourceConfigServiceImpl) UpsertResourceProjectConfigV2(ctx context.Context, req dtos.CreateProjectResourceConfigRequest) error {
	program, err := s.programService.GetProgramByID(ctx, *req.ProjectResourceConfig.ProgramID)
	if err != nil {
		return err
	}

	if req.Icon != nil {
		if err := s.handleExistingIcon(ctx, program.ProgramNameTH, req.ProjectResourceConfig); err != nil {
			return err
		}

		if err := s.uploadNewIcon(ctx, program.ProgramNameTH, req.Icon, req.ProjectResourceConfig); err != nil {
			return err
		}
	}

	return s.projectResourceConfigRepo.UpsertResourceProjectConfigV2(req.ProjectResourceConfig)
}

func (s *projectResourceConfigServiceImpl) handleExistingIcon(ctx context.Context, programNameTH string, config *models.ProjectResourceConfig) error {
	if config.ID != 0 {
		existingConfig, err := s.projectResourceConfigRepo.Get(ctx, config.ID)
		if err != nil {
			return err
		}
		if existingConfig.IconName != nil {
			err := s.uploadService.RemoveObject(ctx, "icons", fmt.Sprintf("%s/%s", programNameTH, *existingConfig.IconName), minio.RemoveObjectOptions{})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *projectResourceConfigServiceImpl) uploadNewIcon(ctx context.Context, programNameTH string, icon *multipart.FileHeader, config *models.ProjectResourceConfig) error {
	filePath := fmt.Sprintf("%s/%s", programNameTH, icon.Filename)
	file, err := icon.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	contentType := icon.Header.Get("Content-Type")
	fileSize := icon.Size

	_, err = s.uploadService.UploadObject(ctx, "icons", filePath, file, fileSize, contentType)
	if err != nil {
		return err
	}
	config.IconName = &icon.Filename
	return nil
}
