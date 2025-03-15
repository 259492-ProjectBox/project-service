// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	db3 "github.com/project-box/db/minio"
	db2 "github.com/project-box/db/postgres"
	"github.com/project-box/db/rabbitmq"
	"github.com/project-box/handlers"
	"github.com/project-box/repositories"
	"github.com/project-box/services"
)

// Injectors from wire.go:

func InitializeApp() (*gin.Engine, func(), error) {
	channel := db.NewRabbitMQConnection()
	gormDB := db2.NewPostgresDatabase()
	fileExtensionRepository := repositories.NewFileExtensionRepository(gormDB)
	projectStaffRepository := repositories.NewProjectStaffRepository(gormDB)
	projectNumberCounterRepository := repositories.NewProjectNumberCounterRepository(gormDB)
	resourceTypeRepository := repositories.NewResourceTypeRepository(gormDB)
	client, err := db3.NewMinIOConnection()
	if err != nil {
		return nil, nil, err
	}
	uploadRepository := repositories.NewUploadRepository(client)
	resourceRepository := repositories.NewResourceRepository(gormDB, resourceTypeRepository, fileExtensionRepository, uploadRepository)
	projectRepository := repositories.NewProjectRepository(gormDB, fileExtensionRepository, projectStaffRepository, projectNumberCounterRepository, resourceRepository, resourceTypeRepository, uploadRepository)
	staffRepository := repositories.NewStaffRepository(gormDB)
	programRepository := repositories.NewProgramRepository(gormDB)
	projectService := services.NewProjectService(channel, projectRepository, staffRepository, programRepository)
	projectHandler := handlers.NewProjectHandler(projectService)
	resourceService := services.NewResourceService(resourceRepository)
	resourceHandler := handlers.NewResourceHandler(client, resourceService, projectService)
	staffService := services.NewStaffService(staffRepository)
	staffHandler := handlers.NewStaffHandler(staffService)
	configRepository := repositories.NewConfigRepository(gormDB)
	configService := services.NewConfigService(configRepository)
	configHandler := handlers.NewConfigHandler(configService)
	keywordRepository := repositories.NewKeywordRepository(gormDB)
	keywordService := services.NewKeywordService(keywordRepository)
	keywordHandler := handlers.NewKeywordHandler(keywordService)
	projectConfigRepository := repositories.NewProjectConfigRepository(gormDB)
	projectConfigService := services.NewProjectConfigService(projectConfigRepository)
	projectConfigHandler := handlers.NewProjectConfigHandler(projectConfigService)
	projectResourceConfigRepository := repositories.NewProjectResourceConfigRepository(gormDB)
	programService := services.NewProgramService(programRepository)
	projectRoleRepository := repositories.NewProjectRoleRepository(gormDB)
	projectRoleService := services.NewProjectRoleService(projectRoleRepository)
	studentRepository := repositories.NewStudentRepository(gormDB, configRepository)
	studentService := services.NewStudentService(configService, studentRepository, gormDB)
	uploadService := services.NewUploadService(client, programRepository, projectRepository, staffService, projectRoleService, projectService, configService, studentService)
	projectResourceConfigService := services.NewProjectResourceConfigService(projectResourceConfigRepository, programService, uploadService)
	projectResourceConfigHandler := handlers.NewProjectResourceConfigHandler(projectResourceConfigService)
	projectRoleHandler := handlers.NewProjectRoleHandler(projectRoleService)
	programHandler := handlers.NewProgramHandler(programService)
	studentHandler := handlers.NewStudentHandler(studentService)
	uploadHandler := handlers.NewUploadHandler(uploadService)
	engine, err := NewApp(projectHandler, resourceHandler, staffHandler, configHandler, keywordHandler, projectConfigHandler, projectResourceConfigHandler, projectRoleHandler, programHandler, studentHandler, uploadHandler)
	if err != nil {
		return nil, nil, err
	}
	return engine, func() {
	}, nil
}

// wire.go:

var AppSet = wire.NewSet(
	NewApp, db2.NewPostgresDatabase, db3.NewMinIOConnection, db.NewRabbitMQConnection,
)

var HandlerSet = wire.NewSet(handlers.NewProjectHandler, handlers.NewResourceHandler, handlers.NewStaffHandler, handlers.NewConfigHandler, handlers.NewProjectConfigHandler, handlers.NewProjectResourceConfigHandler, handlers.NewProjectRoleHandler, handlers.NewProgramHandler, handlers.NewStudentHandler, handlers.NewUploadHandler, handlers.NewKeywordHandler)

var ServiceSet = wire.NewSet(services.NewProjectService, services.NewResourceService, services.NewStaffService, services.NewConfigService, services.NewProjectConfigService, services.NewProjectResourceConfigService, services.NewProjectRoleService, services.NewProgramService, services.NewStudentService, services.NewUploadService, services.NewKeywordService)

var RepositorySet = wire.NewSet(repositories.NewProjectRepository, repositories.NewProjectStaffRepository, repositories.NewProjectNumberCounterRepository, repositories.NewStaffRepository, repositories.NewFileExtensionRepository, repositories.NewProgramRepository, repositories.NewResourceRepository, repositories.NewResourceTypeRepository, repositories.NewConfigRepository, repositories.NewProjectConfigRepository, repositories.NewProjectResourceConfigRepository, repositories.NewProjectRoleRepository, repositories.NewStudentRepository, repositories.NewUploadRepository, repositories.NewKeywordRepository)

var RedisSet = wire.NewSet()
