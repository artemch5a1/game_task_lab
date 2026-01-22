package services

import (
	"context"
	"errors"
	"example/web-service-gin/internal/application/abstraction/repository"
	"example/web-service-gin/internal/domain/model"
	"time"

	"github.com/google/uuid"
)

type GameService struct {
	repo repository.GameRepository
}

func NewGameService(repo repository.GameRepository) *GameService {
	return &GameService{
		repo: repo,
	}
}

func (s *GameService) CreateGame(ctx context.Context,
	title, description string,
	releaseDate time.Time,
	genreID uuid.UUID) (*model.Game, error) {

	if err := s.validateGameData(title, description, releaseDate); err != nil {
		return nil, err
	}

	game, err := model.NewGameWithValidate(title, description, releaseDate, genreID)

	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, game); err != nil {
		return nil, err
	}

	return game, nil
}

func (s *GameService) validateGameData(title, description string, releaseDate time.Time) error {
	if title == "" || len(title) > 200 {
		return errors.New("title must be between 1 and 200 characters")
	}

	if len(description) > 2000 {
		return errors.New("description too long")
	}

	if releaseDate.After(time.Now().AddDate(1, 0, 0)) {
		return errors.New("release date too far in the future")
	}

	return nil
}
