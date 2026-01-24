package auth

import (
	"context"

	specifictype "example/web-service-gin/internal/domain/specific_type"

	"github.com/google/uuid"
)

// TokenVerifier validates incoming auth tokens and extracts identity.
// Concrete implementations (JWT, etc.) must live in infrastructure.
type TokenVerifier interface {
	Verify(ctx context.Context, tokenString string) (userID uuid.UUID, role specifictype.UserRole, err error)
}

