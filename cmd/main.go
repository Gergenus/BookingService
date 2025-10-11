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
	bookRepo := repository.NewPostgresBookingRepository(db)

	equipService := service.NewEquipmentService(log, &postRepo, miniRepo)
	bookService := service.NewBookingService(&bookRepo, log)

	equipHandler := handler.NewEquipmentHandler(&equipService)
	bookHandler := handler.NewBookingHandler(&bookService)

	e := echo.New()
	eq := e.Group("/api/v1/equipment")
	{
		eq.POST("/create", equipHandler.CreateEquipment)
		eq.GET("", equipHandler.EquipmentByName)
		eq.PUT("/update", nil)
		eq.DELETE("/:id", equipHandler.DeleteEquipment)
		eq.GET("/:id", equipHandler.EquipmentById)
	}
	auth := e.Group("/api/v1/auth")
	{
		auth.POST("/register", nil)
		auth.POST("/login", nil)
		auth.POST("/refresh", nil)
		auth.POST("/logout", nil)
	}
	booking := e.Group("/api/v1/booking")
	{
		booking.POST("/", bookHandler.Createbooking)
		booking.DELETE("/:id", bookHandler.DeleteBooking)
		booking.GET("/:id", bookHandler.Bookings)
	}
	e.Start(":" + cfg.HTTPPort)
}
