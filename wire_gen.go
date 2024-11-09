// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/project-box/configs"
	"github.com/project-box/db/postgres"
	db2 "github.com/project-box/db/rabbitmq"
	"github.com/project-box/handlers"
	"github.com/project-box/repositories"
	"github.com/project-box/services"
)

// Injectors from wire.go:

func InitializeApp() (*gin.Engine, func(), error) {
	gormDB := db.NewPostgresDatabase()
	channel := db2.NewRabbitMQConnection()
	projectRepository := repositories.NewProjectRepository(gormDB)
	employeeRepository := repositories.NewEmployeeRepository(gormDB)
	majorRepository := repositories.NewMajorRepository(gormDB)
	sectionRepository := repositories.NewSectionRepository(gormDB)
	projectService := services.NewProjectService(channel, projectRepository, employeeRepository, majorRepository, sectionRepository)
	projectHandler := handlers.NewProjectHandler(projectService)
	client, err := configs.InitializeMinioClient()
	if err != nil {
		return nil, nil, err
	}
	resourceRepository := repositories.NewResourceRepository(gormDB)
	resourceService := services.NewResourceService(resourceRepository)
	resourceHandler := handlers.NewResourceHandler(client, resourceService)
	engine, err := NewApp(gormDB, channel, projectHandler, resourceHandler, client)
	if err != nil {
		return nil, nil, err
	}
	return engine, func() {
	}, nil
}

// wire.go:

var AppSet = wire.NewSet(
	NewApp, db.NewPostgresDatabase, db2.NewRabbitMQConnection,
)

var HandlerSet = wire.NewSet(handlers.NewProjectHandler, handlers.NewResourceHandler)

var ServiceSet = wire.NewSet(services.NewProjectService, services.NewResourceService)

var RepositorySet = wire.NewSet(repositories.NewProjectRepository, repositories.NewEmployeeRepository, repositories.NewMajorRepository, repositories.NewSectionRepository, repositories.NewResourceRepository)

var MinioSet = wire.NewSet(configs.InitializeMinioClient)

var RedisSet = wire.NewSet()
