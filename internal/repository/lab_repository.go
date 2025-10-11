package repository

import (
	"context"
	"fmt"

	"github.com/Gergenus/bookingService/internal/models"
	"github.com/Gergenus/bookingService/pkg/db"
)

type PostgresLabRepository struct {
	db db.PostgresDB
}

type LabRepositroy interface {
	CreateEquipment(ctx context.Context, equipment models.Equipment) (int, error)
	Equipment(ctx context.Context, equipment_id int) (*models.Equipment, error)
	DeleteEquipment(ctx context.Context, equipment_id int) error
	UpdateEquipment(ctx context.Context, equipment models.Equipment) error
}

func NewPostgresLabRepository(db db.PostgresDB) PostgresLabRepository {
	return PostgresLabRepository{db: db}
}

// TODO обработку sql ошибок

func (p *PostgresLabRepository) CreateEquipment(ctx context.Context, equipment models.Equipment) (int, error) {
	const op = "lab_repository.CreateEquipment"
	var id int
	err := p.db.DB.QueryRow(ctx, "INSERT INTO equipment (equipment_name, manufacturer, description, image_url) VALUES($1, $2, $3, $4) RETURNING id", equipment.EquipmentName,
		equipment.Manufacturer, equipment.Description, equipment.ImageURL).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (p *PostgresLabRepository) Equipment(ctx context.Context, equipment_id int) (*models.Equipment, error) {
	const op = "lab_repository.Equipment"
	var equipment models.Equipment
	err := p.db.DB.QueryRow(ctx, "SELECT * FROM equipment WHERE id = $1", equipment_id).Scan(&equipment)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &equipment, nil
}

func (p *PostgresLabRepository) DeleteEquipment(ctx context.Context, equipment_id int) error {
	const op = "lab_repository.DeleteEquipment"
	_, err := p.db.DB.Exec(ctx, "DELETE FROM equipment WHERE equipment_id = $1", equipment_id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (p *PostgresLabRepository) UpdateEquipment(ctx context.Context, equipment models.Equipment) error {
	panic("implement me")
}
