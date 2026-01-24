package dto

import (
	"example/web-service-gin/internal/domain/specific_type"

	"github.com/google/uuid"
)

type UserDto struct {
	ID       uuid.UUID             `json:"id"`
	Username string                `json:"username"`
	UserRole specifictype.UserRole `json:"userRole"`
}

type CreateUserDto struct {
	Username string                `json:"username" validate:"required,min=1,max=200"`
	Password string                `json:"password" validate:"required,min=1,max=200"`
	UserRole specifictype.UserRole `json:"userRole" validate:"required"`
}

type UpdateUserDto struct {
	ID       uuid.UUID             `json:"id" validate:"required"`
	Username string                `json:"username" validate:"required,min=1,max=200"`
	Password string                `json:"password" validate:"required,min=1,max=200"`
	UserRole specifictype.UserRole `json:"userRole" validate:"required"`
}

type LoginDto struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

