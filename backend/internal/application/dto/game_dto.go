// internal/application/dto/game_dto.go
package dto

import (
	"time"

	"github.com/google/uuid"
)

type GameDto struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ReleaseDate time.Time `json:"releaseDate"`
	GenreID     uuid.UUID `json:"genreId"`
}

type GameDtoWithStats struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	ReleaseDate   time.Time `json:"releaseDate"`
	GenreID       uuid.UUID `json:"genreId"`
	AverageRating float64   `json:"averageRating"`
	RatingCount   int       `json:"ratingCount"`
}

type CreateGameDto struct {
	Title       string    `json:"title" validate:"required,min=1,max=200"`
	Description string    `json:"description" validate:"max=2000"`
	ReleaseDate time.Time `json:"releaseDate" validate:"required"`
	GenreID     uuid.UUID `json:"genreId" validate:"required"`
}

type UpdateGameDto struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	Title       string    `json:"title" validate:"required,min=1,max=200"`
	Description string    `json:"description" validate:"max=2000"`
	ReleaseDate time.Time `json:"releaseDate" validate:"required"`
	GenreID     uuid.UUID `json:"genreId" validate:"required"`
}

// Response DTO для API
type GameResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Items      []*GameDto `json:"items"`
	Total      int        `json:"total"`
	Page       int        `json:"page"`
	PageSize   int        `json:"pageSize"`
	TotalPages int        `json:"totalPages"`
}
