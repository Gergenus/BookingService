package dto

import "github.com/google/uuid"

type UserDTO struct {
	UUID     uuid.UUID `json:"uuid,omitempty"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
