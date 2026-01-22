package model

import (
	"time"

	"github.com/google/uuid"
)

type Game struct {
	ID          uuid.UUID
	Title       string
	Description string
	releaseDate time.Time
	GenreID     uuid.UUID
}
