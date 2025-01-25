//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	minio "github.com/project-box/db/minio"
	database "github.com/project-box/db/postgres"
	rabbitMQ "github.com/project-box/db/rabbitmq"
	"github.com/project-box/handlers"
	"github.com/project-box/repositories"
	"github.com/project-box/services"
)

func InitializeApp() (*gin.Engine, func(), error) {
	wire.Build(
		AppSet,
		HandlerSet,
		ServiceSet,
		RepositorySet)

	return gin.New(), func() {}, nil
}

var AppSet = wire.NewSet(
	NewApp,
	database.NewPostgresDatabase,
	minio.NewMinIOConnection,
	rabbitMQ.NewRabbitMQConnection,
)

var HandlerSet = wire.NewSet(
	handlers.NewProjectHandler,
	handlers.NewCalendarHandler,
	handlers.NewCourseHandler,
	handlers.NewResourceHandler,
	handlers.NewStaffHandler,
	handlers.NewConfigHandler,
	handlers.NewProjectConfigHandler,
	handlers.NewProjectResourceConfigHandler,
	handlers.NewProgramHandler,
	handlers.NewStudentHandler,
	handlers.NewUploadHandler,
)

var ServiceSet = wire.NewSet(
	services.NewProjectService,
	services.NewCalendarService,
	services.NewCourseService,
	services.NewResourceService,
	services.NewStaffService,
	services.NewConfigService,
	services.NewProjectConfigService,
	services.NewProjectResourceConfigService,
	services.NewProgramService,
	services.NewStudentService,
	services.NewUploadService,
)

var RepositorySet = wire.NewSet(
	repositories.NewProjectRepository,
	repositories.NewProjectStaffRepository,
	repositories.NewProjectNumberCounterRepository,
	repositories.NewStaffRepository,
	repositories.NewFileExtensionRepository,
	repositories.NewProgramRepository,
	repositories.NewCourseRepository,
	repositories.NewSectionRepository,
	repositories.NewResourceRepository,
	repositories.NewResourceTypeRepository,
	repositories.NewCalendarRepository,
	repositories.NewConfigRepository,
	repositories.NewProjectConfigRepository,
	repositories.NewProjectResourceConfigRepository,
	repositories.NewStudentRepository,
	repositories.NewUploadRepository,
)

var RedisSet = wire.NewSet()
