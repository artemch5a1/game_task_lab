package repository

import (
	"context"
	"errors"
	"example/web-service-gin/internal/domain/model"

	"github.com/google/uuid"
)

var (
	ErrNotFound      = errors.New("game not found")
	ErrAlreadyExists = errors.New("game already exists")
)

type GameRepository interface {
	Create(ctx context.Context, game *model.Game) (*model.Game, error)

	FindByID(ctx context.Context, id uuid.UUID) (*model.Game, error)

	FindAll(ctx context.Context, limit, offset int) ([]*model.Game, error)

	Update(ctx context.Context, game *model.Game) (*model.Game, error)

	Delete(ctx context.Context, id uuid.UUID) error

	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}
