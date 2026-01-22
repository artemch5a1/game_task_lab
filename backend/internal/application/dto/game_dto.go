package dto

import (
	"time"

	"github.com/google/uuid"
)

type GameDto struct {
	ID          uuid.UUID
	Title       string
	Description string
	ReleaseDate time.Time
	GenreID     uuid.UUID
}

type GameDtoWithStats struct {
	ID            uuid.UUID
	Title         string
	Description   string
	ReleaseDate   time.Time
	GenreID       uuid.UUID
	AverageRating float64
	RatingCount   int
}

type CreateGameDto struct {
	Title       string
	Description string
	ReleaseDate time.Time
	GenreID     uuid.UUID
}

type UpdateGameDto struct {
	ID          uuid.UUID
	Title       string
	Description string
	ReleaseDate time.Time
	GenreID     uuid.UUID
}
