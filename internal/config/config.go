package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresURL          string
	LogLevel             string
	HTTPPort             string
	JWTSecret            string
	MinioEndpoint        string
	MinioaccessKeyId     string
	MinioSecretAccessKey string
	MinioBucket          string
	// 	AccessTTL   time.Duration
	// 	RefreshTTl  time.Duration
}

func InitConfig() Config {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	// AccessTTL, err := time.ParseDuration(os.Getenv("ACCESS_TTL"))
	// if err != nil {
	// 	panic(err)
	// }
	// RefreshTTl, err := time.ParseDuration(os.Getenv("REFRESH_TTL"))
	// if err != nil {
	// 	panic(err)
	// }
	return Config{
		PostgresURL:          os.Getenv("POSTGRES_URL"),
		LogLevel:             os.Getenv("LOG_LEVEL"),
		HTTPPort:             os.Getenv("HTTP_PORT"),
		JWTSecret:            os.Getenv("JWT_SECRET"),
		MinioEndpoint:        os.Getenv("MINIO_ENDPOINT"),
		MinioaccessKeyId:     os.Getenv("MINIO_ACCESS_KEY_ID"),
		MinioSecretAccessKey: os.Getenv("MINIO_SECRET_ACCCESS_KEY"),
		MinioBucket:          os.Getenv("MINIO_BUCKET"),
		// AccessTTL:   AccessTTL,
		// RefreshTTl:  RefreshTTl,
	}
}
