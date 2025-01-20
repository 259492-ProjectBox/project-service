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
	resourceTypeRepository := repositories.NewResourceTypeRepository(gormDB)
	client, err := db3.NewMinIOConnection()
	if err != nil {
		return nil, nil, err
	}
	uploadRepository := repositories.NewUploadRepository(client)
	resourceRepository := repositories.NewResourceRepository(gormDB, resourceTypeRepository, fileExtensionRepository, uploadRepository)
	projectRepository := repositories.NewProjectRepository(gormDB, fileExtensionRepository, projectStaffRepository, resourceRepository, resourceTypeRepository, uploadRepository)
	staffRepository := repositories.NewStaffRepository(gormDB)
	programRepository := repositories.NewProgramRepository(gormDB)
	courseRepository := repositories.NewCourseRepository(gormDB)
	projectNumberCounterRepository := repositories.NewProjectNumberCounterRepository(gormDB)
	projectService := services.NewProjectService(channel, projectRepository, staffRepository, programRepository, courseRepository, projectNumberCounterRepository)
	projectHandler := handlers.NewProjectHandler(projectService)
	calendarRepository := repositories.NewCalendarRepository(gormDB)
	calendarService := services.NewCalendarService(calendarRepository, programRepository)
	calendarHandler := handlers.NewCalendarHandler(calendarService)
	resourceService := services.NewResourceService(resourceRepository)
	resourceHandler := handlers.NewResourceHandler(client, resourceService, projectService)
	staffService := services.NewStaffService(staffRepository)
	staffHandler := handlers.NewStaffHandler(staffService)
	courseService := services.NewCourseService(courseRepository)
	courseHandler := handlers.NewCourseHandler(courseService)
	configRepository := repositories.NewConfigRepository(gormDB)
	configService := services.NewConfigService(configRepository)
	configHandler := handlers.NewConfigHandler(configService)
	projectConfigRepository := repositories.NewProjectConfigRepository(gormDB)
	projectConfigService := services.NewProjectConfigService(projectConfigRepository)
	projectConfigHandler := handlers.NewProjectConfigHandler(projectConfigService)
	projectResourceConfigRepository := repositories.NewProjectResourceConfigRepository(gormDB)
	projectResourceConfigService := services.NewProjectResourceConfigService(projectResourceConfigRepository)
	projectResourceConfigHandler := handlers.NewProjectResourceConfigHandler(projectResourceConfigService)
	programService := services.NewProgramService(programRepository)
	programHandler := handlers.NewProgramHandler(programService)
	engine, err := NewApp(projectHandler, calendarHandler, resourceHandler, staffHandler, courseHandler, configHandler, projectConfigHandler, projectResourceConfigHandler, programHandler)
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

var HandlerSet = wire.NewSet(handlers.NewProjectHandler, handlers.NewCalendarHandler, handlers.NewCourseHandler, handlers.NewResourceHandler, handlers.NewStaffHandler, handlers.NewConfigHandler, handlers.NewProjectConfigHandler, handlers.NewProjectResourceConfigHandler, handlers.NewProgramHandler, handlers.NewStudentHandler)

var ServiceSet = wire.NewSet(services.NewProjectService, services.NewCalendarService, services.NewCourseService, services.NewResourceService, services.NewStaffService, services.NewConfigService, services.NewProjectConfigService, services.NewProjectResourceConfigService, services.NewProgramService)

var RepositorySet = wire.NewSet(repositories.NewProjectRepository, repositories.NewProjectStaffRepository, repositories.NewProjectNumberCounterRepository, repositories.NewStaffRepository, repositories.NewFileExtensionRepository, repositories.NewProgramRepository, repositories.NewCourseRepository, repositories.NewSectionRepository, repositories.NewResourceRepository, repositories.NewResourceTypeRepository, repositories.NewCalendarRepository, repositories.NewConfigRepository, repositories.NewProjectConfigRepository, repositories.NewProjectResourceConfigRepository, repositories.NewUploadRepository)

var RedisSet = wire.NewSet()
