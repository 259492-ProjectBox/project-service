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
	client, err := db3.NewMinIOConnection()
	if err != nil {
		return nil, nil, err
	}
	projectRepository := repositories.NewProjectRepository(gormDB, client)
	employeeRepository := repositories.NewEmployeeRepository(gormDB)
	majorRepository := repositories.NewMajorRepository(gormDB)
	courseRepository := repositories.NewCourseRepository(gormDB)
	projectNumberCounterRepository := repositories.NewProjectNumberCounterRepository(gormDB)
	projectService := services.NewProjectService(channel, projectRepository, employeeRepository, majorRepository, courseRepository, projectNumberCounterRepository)
	projectHandler := handlers.NewProjectHandler(projectService)
	engine, err := NewApp(projectHandler)
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

var HandlerSet = wire.NewSet(handlers.NewProjectHandler)

var ServiceSet = wire.NewSet(services.NewProjectService)

var RepositorySet = wire.NewSet(repositories.NewProjectRepository, repositories.NewProjectNumberCounterRepository, repositories.NewEmployeeRepository, repositories.NewMajorRepository, repositories.NewCourseRepository, repositories.NewSectionRepository, repositories.NewResourceRepository)

var RedisSet = wire.NewSet()
