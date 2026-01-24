package data

import (
	"sync"

	"example/web-service-gin/internal/domain/model"

	"github.com/google/uuid"
)

// Data - единый in-memory провайдер данных для репозиториев.
// Репозитории из пакета `inmemory` обязаны использовать один и тот же экземпляр Data.
type Data struct {
	Mu sync.RWMutex

	Games       map[uuid.UUID]*model.Game
	Genres      map[uuid.UUID]*model.Genre
	Users       map[uuid.UUID]*model.User
	UserRatings map[uuid.UUID]*model.UserRating
}

func New() *Data {
	return &Data{
		Games:       make(map[uuid.UUID]*model.Game),
		Genres:      make(map[uuid.UUID]*model.Genre),
		Users:       make(map[uuid.UUID]*model.User),
		UserRatings: make(map[uuid.UUID]*model.UserRating),
	}
}

