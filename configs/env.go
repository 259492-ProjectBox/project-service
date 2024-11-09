package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func InitialEnv(path string) {
	if err := godotenv.Load(path); err != nil {
		log.Println("No .env file found")
	}
}

func GetPort() string {
	port, ok := os.LookupEnv("PORT")

	if !ok {
		return "8000"
	}

	return port
}
func InitializeMinioClient() (*minio.Client, error) {
	// Get MinIO configuration from environment variables
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")

	// Initialize MinIO client
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false, // Set to true if using HTTPS
	})
	if err != nil {
		// Return the error so it can be handled by the caller
		return nil, fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

	fmt.Println("Successfully connected to MinIO")
	return minioClient, nil
}
