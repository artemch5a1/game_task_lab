package model

import (
	specifictype "example/web-service-gin/internal/domain/specific_type"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	Password string
	Username string
	UserRole specifictype.UserRole
}
