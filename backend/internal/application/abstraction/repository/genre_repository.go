package repository

import (
	"context"
	"errors"

	"example/web-service-gin/internal/domain/model"

	"github.com/google/uuid"
)

var (
	ErrGenreNotFound      = errors.New("genre not found")
	ErrGenreAlreadyExists = errors.New("genre already exists")
)

type GenreRepository interface {
	Create(ctx context.Context, genre *model.Genre) (*model.Genre, error)

	FindByID(ctx context.Context, id uuid.UUID) (*model.Genre, error)

	FindAll(ctx context.Context, limit, offset int) ([]*model.Genre, error)

	Update(ctx context.Context, genre *model.Genre) (*model.Genre, error)

	Delete(ctx context.Context, id uuid.UUID) error

	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}

