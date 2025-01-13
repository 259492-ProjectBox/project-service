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
	handlers.NewResourceHandler,
	handlers.NewEmployeeHandler,
	handlers.NewConfigHandler,
	handlers.NewProjectConfigHandler,
	handlers.NewProgramHandler,
)

var ServiceSet = wire.NewSet(
	services.NewProjectService,
	services.NewCalendarService,
	services.NewResourceService,
	services.NewEmployeeService,
	services.NewConfigService,
	services.NewProjectConfigService,
	services.NewProgramService,
)

var RepositorySet = wire.NewSet(
	repositories.NewProjectRepository,
	repositories.NewProjectNumberCounterRepository,
	repositories.NewEmployeeRepository,
	repositories.NewProgramRepository,
	repositories.NewCourseRepository,
	repositories.NewSectionRepository,
	repositories.NewResourceRepository,
	repositories.NewCalendarRepository,
	repositories.NewConfigRepository,
	repositories.NewProjectConfigRepository,
)

var RedisSet = wire.NewSet()
