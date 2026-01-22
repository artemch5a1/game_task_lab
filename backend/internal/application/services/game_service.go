package services

import (
	"context"
	"errors"
	"example/web-service-gin/internal/application/abstraction/repository"
	"example/web-service-gin/internal/application/dto"
	"example/web-service-gin/internal/application/mapper"
	"time"
)

type GameService struct {
	repo       repository.GameRepository
	gameMapper *mapper.GameMapper
}

func NewGameService(repo repository.GameRepository) *GameService {
	return &GameService{
		repo:       repo,
		gameMapper: mapper.NewGameMapper(),
	}
}

func (s *GameService) CreateGame(ctx context.Context,
	gameCreateDto dto.CreateGameDto) (*dto.GameDto, error) {

	if err := s.validateGameData(
		gameCreateDto); err != nil {
		return nil, err
	}

	game, mapperError := s.gameMapper.FromCreateGameDto(&gameCreateDto)

	if mapperError != nil {
		return nil, mapperError
	}

	createdGame, repoError := s.repo.Create(ctx, game)

	if repoError != nil {
		return nil, repoError
	}

	return s.gameMapper.ToGameDto(createdGame), nil
}

func (s *GameService) validateGameData(gameDto dto.CreateGameDto) error {
	if gameDto.Title == "" || len(gameDto.Title) > 200 {
		return errors.New("title must be between 1 and 200 characters")
	}

	if len(gameDto.Description) > 2000 {
		return errors.New("description too long")
	}

	if gameDto.ReleaseDate.After(time.Now().AddDate(1, 0, 0)) {
		return errors.New("release date too far in the future")
	}

	return nil
}
