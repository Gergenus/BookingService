package service

import (
	"context"
	"fmt"
	"log/slog"
	"mime/multipart"

	"github.com/Gergenus/bookingService/internal/models"
	"github.com/Gergenus/bookingService/internal/repository"
	"github.com/minio/minio-go/v7"
)

type EquipmentService struct {
	log  *slog.Logger
	repo repository.LabRepositroy
	mini repository.ImageRepositoryInterface
}

type EquipmentServiceInterface interface {
	CreateEquipment(ctx context.Context, equipment models.Equipment, image *multipart.FileHeader) (int, error)
	Equipment(ctx context.Context, equipment_id int) (*models.Equipment, error)
	EquipmentByName(ctx context.Context, equipmentName string) ([]models.Equipment, error)
	DeleteEquipment(ctx context.Context, equipment_id int) error
	UpdateEquipment(ctx context.Context, equipment models.Equipment) error
	SignURL(ctx context.Context, imagePath string) (*minio.Object, error)
}

func NewEquipmentService(log *slog.Logger, repo repository.LabRepositroy, mini repository.ImageRepositoryInterface) EquipmentService {
	return EquipmentService{log: log, repo: repo, mini: mini}
}

func (e *EquipmentService) SignURL(ctx context.Context, imagePath string) (*minio.Object, error) {
	const op = "equipment_service.SignURL"
	obj, err := e.mini.SignURL(ctx, imagePath)
	if err != nil {
		e.log.Error("signing image to s3 storage error", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	return obj, nil
}

func (e *EquipmentService) CreateEquipment(ctx context.Context, equipment models.Equipment, image *multipart.FileHeader) (int, error) {
	const op = "equipment_service.CreateEquipment"
	e.log.Info("creating equipment", slog.String("equipment_name", equipment.EquipmentName))
	e.log.Info("adding image to s3 storage", slog.String("image", image.Filename))
	url, err := e.mini.AddImage(ctx, image)
	if err != nil {
		e.log.Error("adding image to s3 storage error", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s, %w", op, err)
	}
	/*
		Логика замены equipment.ImageURL url из минио
	*/
	equipment.ImageURL = url

	id, err := e.repo.CreateEquipment(ctx, equipment)
	if err != nil {
		e.log.Error("creating equipment error", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (e *EquipmentService) Equipment(ctx context.Context, equipment_id int) (*models.Equipment, error) {
	const op = "equipment_service.Equipment"
	e.log.Info("getting equipment", slog.Int("equipment_id", equipment_id))
	eq, err := e.repo.Equipment(ctx, equipment_id)
	if err != nil {
		e.log.Error("getting equipment error", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return eq, nil
}

func (e *EquipmentService) DeleteEquipment(ctx context.Context, equipment_id int) error {
	const op = "equipment_service.DeleteEquipment"
	e.log.Info("deleting equipment", slog.Int("equipment_id", equipment_id))

	eq, err := e.repo.Equipment(ctx, equipment_id)
	if err != nil {
		e.log.Error("getting equipment error", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	err = e.mini.DeleteImage(ctx, eq.ImageURL)
	if err != nil {
		e.log.Error("deleting image in miniO error", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}
	err = e.repo.DeleteEquipment(ctx, equipment_id)
	if err != nil {
		e.log.Error("deleting equipment error", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (e *EquipmentService) UpdateEquipment(ctx context.Context, equipment models.Equipment) error {
	const op = "equipment_service.UpdateEquipment"
	e.log.Info("updating equipment", slog.Int("equipment_id", equipment.EquipmentId))
	err := e.repo.UpdateEquipment(ctx, equipment)
	if err != nil {
		e.log.Error("updating equipment error", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (e *EquipmentService) EquipmentByName(ctx context.Context, equipmentName string) ([]models.Equipment, error) {
	const op = "equipment_service.EquipmentByName"
	log := e.log.With(slog.String("op", op))
	log.Info("getting equipment by name", slog.String("name", equipmentName))
	eqs, err := e.repo.EquipmentByName(ctx, equipmentName)
	if err != nil {
		log.Error("getting eq equipment by name error", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return eqs, nil
}
