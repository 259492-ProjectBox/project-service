package db

import (
	"fmt"
	"log"
	"os"

	"github.com/project-box/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func GetPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_NAME"),
		SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
	}
}

func migrateModel(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&models.Role{},
		&models.Project{},
		&models.Employee{},
		&models.Student{},
		&models.Comment{},
		&models.Course{},
		&models.Section{},
		&models.Major{},
		&models.ImportantDate{},
		&models.ResourceType{},
		&models.Resource{},
	); err != nil {
		return fmt.Errorf("migration error: %w", err)
	}
	return nil
}

func NewPostgresDatabase() *gorm.DB {
	configs := GetPostgresConfig()
	if configs == nil {
		return nil
	}
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		configs.Host, configs.User, configs.Password, configs.DBName, configs.Port, configs.SSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Println(err)
		return nil
	}

	if err = migrateModel(db); err != nil {
		log.Println("migrateModel", err)
		return nil
	}

	return db
}
