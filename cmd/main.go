package main

import (
	"github.com/Gergenus/bookingService/internal/config"
	"github.com/Gergenus/bookingService/internal/handler"
	"github.com/Gergenus/bookingService/internal/middleware"
	"github.com/Gergenus/bookingService/internal/repository"
	"github.com/Gergenus/bookingService/internal/service"
	"github.com/Gergenus/bookingService/pkg/db"
	"github.com/Gergenus/bookingService/pkg/jwtpkg"
	"github.com/Gergenus/bookingService/pkg/logger"
	"github.com/Gergenus/bookingService/pkg/redispkg"
	"github.com/Gergenus/bookingService/pkg/s3"
	"github.com/labstack/echo/v4"
)

func main() {
	cfg := config.InitConfig()
	db := db.InitDB(cfg.PostgresURL)
	redisDB := redispkg.InitRedisDB(cfg.RedisAddress, cfg.RedisPassword, cfg.RedisDB)
	miniClient := s3.InitS3Storage(cfg.MinioEndpoint, cfg.MinioaccessKeyId, cfg.MinioSecretAccessKey, cfg.MinioBucket)
	log := logger.SetUp(cfg.LogLevel)
	JWT := jwtpkg.NewUserJWTpkg(cfg.JWTSecret, cfg.AccessTTL)
	middle := middleware.NewJWTMiddleware(JWT)

	miniRepo := repository.NewMinioImageRepository(miniClient, cfg.MinioBucket, cfg.MinioEndpoint)
	postRepo := repository.NewPostgresLabRepository(db)
	bookRepo := repository.NewPostgresBookingRepository(db)
	userRepo := repository.NewUserRepository(db, redisDB)

	equipService := service.NewEquipmentService(log, &postRepo, miniRepo)
	bookService := service.NewBookingService(&bookRepo, log)
	userService := service.NewUserService(userRepo, log, JWT, cfg.RefreshTTL)

	equipHandler := handler.NewEquipmentHandler(&equipService)
	bookHandler := handler.NewBookingHandler(&bookService)
	userHandler := handler.NewUserHandler(userService, cfg.AdminSecret)

	e := echo.New()
	eq := e.Group("/api/v1/equipment", middle.Auth)
	{
		eq.POST("/create", equipHandler.CreateEquipment, middle.AdminAuth)
		eq.GET("", equipHandler.EquipmentByName)
		eq.PUT("/update", nil, middle.AdminAuth)
		eq.DELETE("/:id", equipHandler.DeleteEquipment, middle.AdminAuth)
		eq.GET("/:id", equipHandler.EquipmentById)
	}
	auth := e.Group("/api/v1/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
		auth.POST("/refresh", userHandler.Refresh)
		auth.POST("/logout", nil)
	}
	booking := e.Group("/api/v1/booking", middle.Auth, middle.ScientistAuth)
	{
		booking.POST("/", bookHandler.Createbooking)
		booking.DELETE("/:id", bookHandler.DeleteBooking)
		booking.GET("/:id", bookHandler.Bookings)
	}
	e.Start(":" + cfg.HTTPPort)
}
