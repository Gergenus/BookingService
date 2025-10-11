package main

import (
	"github.com/Gergenus/bookingService/internal/config"
	"github.com/Gergenus/bookingService/internal/handler"
	"github.com/Gergenus/bookingService/internal/repository"
	"github.com/Gergenus/bookingService/internal/service"
	"github.com/Gergenus/bookingService/pkg/db"
	"github.com/Gergenus/bookingService/pkg/logger"
	"github.com/Gergenus/bookingService/pkg/s3"
	"github.com/labstack/echo/v4"
)

func main() {
	cfg := config.InitConfig()
	db := db.InitDB(cfg.PostgresURL)
	miniClient := s3.InitS3Storage(cfg.MinioEndpoint, cfg.MinioaccessKeyId, cfg.MinioSecretAccessKey, cfg.MinioBucket)
	log := logger.SetUp(cfg.LogLevel)
	miniRepo := repository.NewMinioImageRepository(miniClient, cfg.MinioBucket, cfg.MinioEndpoint)
	postRepo := repository.NewPostgresLabRepository(db)
	equipService := service.NewEquipmentService(log, &postRepo, miniRepo)
	equipHandler := handler.NewEquipmentHandler(equipService)

	e := echo.New()
	e.POST("/addEquipment", equipHandler.CreateEquipment)

	e.Start(":" + cfg.HTTPPort)
}
