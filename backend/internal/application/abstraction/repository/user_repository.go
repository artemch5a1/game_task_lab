package repository

import (
	"context"
	"errors"

	"example/web-service-gin/internal/domain/model"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)

	FindByID(ctx context.Context, id uuid.UUID) (*model.User, error)

	FindByUsername(ctx context.Context, username string) (*model.User, error)

	FindAll(ctx context.Context, limit, offset int) ([]*model.User, error)

	Update(ctx context.Context, user *model.User) (*model.User, error)

	Delete(ctx context.Context, id uuid.UUID) error

	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}

