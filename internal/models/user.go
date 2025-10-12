package models

import "github.com/google/uuid"

type User struct {
	UUID           uuid.UUID `json:"uuid,omitempty"`
	Username       string    `json:"username"`
	Role           string    `json:"role"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashed_password,omitempty"`
}
