//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/project-box/configs"
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
		RepositorySet,
		MinioSet)

	return gin.New(), func() {}, nil
}

var AppSet = wire.NewSet(
	NewApp,
	database.NewPostgresDatabase,
	rabbitMQ.NewRabbitMQConnection,
)

var HandlerSet = wire.NewSet(
	handlers.NewProjectHandler,
	handlers.NewResourceHandler,
)

var ServiceSet = wire.NewSet(
	services.NewProjectService,
	services.NewResourceService,
)

var RepositorySet = wire.NewSet(
	repositories.NewProjectRepository,
	repositories.NewEmployeeRepository,
	repositories.NewMajorRepository,
	repositories.NewSectionRepository,
	repositories.NewResourceRepository,
)
var MinioSet = wire.NewSet(
	configs.InitializeMinioClient,
)
var RedisSet = wire.NewSet()
