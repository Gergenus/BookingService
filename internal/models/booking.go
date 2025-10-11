package models

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	Id          int       `json:"id,omitempty"`
	EquipmentId int       `json:"equipment_id"`
	UserId      uuid.UUID `json:"user_id"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
}
