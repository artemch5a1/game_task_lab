package auth

import (
	"context"

	specifictype "example/web-service-gin/internal/domain/specific_type"

	"github.com/google/uuid"
)

// TokenProvider is an application-level abstraction for issuing auth tokens.
// Concrete implementations (JWT, etc.) must live in infrastructure.
type TokenProvider interface {
	Issue(ctx context.Context, userID uuid.UUID, role specifictype.UserRole) (string, error)
}

