package model

import (
	"time"

	"github.com/google/uuid"
)

type UserRating struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	GameID    uuid.UUID
	Rating    int
	CreatedAt time.Time
}
