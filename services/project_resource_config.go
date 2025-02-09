package services

import (
	"context"
	"fmt"

	"github.com/project-box/dtos"
	"github.com/project-box/models"
	"github.com/project-box/repositories"
)

type ProjectResourceConfigService interface {
	GetProjectResourceConfigsByProgramId(ctx context.Context, programID int) ([]dtos.ProjectResourceConfig, error)
	UpsertResourceProjectConfig(config *models.ProjectResourceConfig) error
	UpsertResourceProjectConfigV2(ctx context.Context, programId int, req dtos.CreateProjectResourceConfigRequest) error
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
			Program:   models.Program{ID: config.Program.ID, ProgramNameTH: config.Program.ProgramNameTH},
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

func (s *projectResourceConfigServiceImpl) UpsertResourceProjectConfigV2(ctx context.Context, programId int, req dtos.CreateProjectResourceConfigRequest) error {
	program, err := s.programService.GetProgramByID(ctx, programId)
	if err != nil {
		return err
	}

	if req.Icon != nil {
		filePath := fmt.Sprintf("%s/%s", program.ProgramNameTH, req.Icon.Filename)
		file, err := req.Icon.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		contentType := req.Icon.Header.Get("Content-Type")
		fileSize := req.Icon.Size

		fmt.Println(filePath)
		_, err = s.uploadService.UploadObject(ctx, "icons", filePath, file, fileSize, contentType)
		if err != nil {
			return err
		}
		req.ProjectResourceConfig.IconName = &req.Icon.Filename
	}

	return s.projectResourceConfigRepo.UpsertResourceProjectConfigV2(req.ProjectResourceConfig)
}
