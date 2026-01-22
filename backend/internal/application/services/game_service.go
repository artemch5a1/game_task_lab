package services

import (
	"context"
	"errors"
	"example/web-service-gin/internal/application/abstraction/repository"
	"example/web-service-gin/internal/application/dto"
	"example/web-service-gin/internal/application/mapper"
	"math"
	"time"

	"github.com/google/uuid"
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

func (s *GameService) GetGameByID(ctx context.Context,
	gameID uuid.UUID) (*dto.GameDto, error) {

	if err := s.validateGameID(gameID); err != nil {
		return nil, err
	}

	game, repoError := s.repo.FindByID(ctx, gameID)

	if repoError != nil {
		return nil, repoError
	}

	return s.gameMapper.ToGameDto(game), nil
}

func (s *GameService) GetAllGames(ctx context.Context) ([]*dto.GameDto, error) {

	games, repoError := s.repo.FindAll(ctx, math.MaxInt, 0)

	if repoError != nil {
		return nil, repoError
	}

	return s.gameMapper.ToGameDtoSlice(games), nil
}

func (s *GameService) UpdateGame(ctx context.Context,
	updateGameDto dto.UpdateGameDto) (*dto.GameDto, error) {

	if err := s.validateUpdateData(updateGameDto); err != nil {
		return nil, err
	}

	existingGame, repoError := s.repo.FindByID(ctx, updateGameDto.ID)

	if repoError != nil {
		return nil, repoError
	}

	mapperError := s.gameMapper.FromUpdateGameDto(existingGame, &updateGameDto)

	if mapperError != nil {
		return nil, mapperError
	}

	updatedGame, repoError := s.repo.Update(ctx, existingGame)

	if repoError != nil {
		return nil, repoError
	}

	return s.gameMapper.ToGameDto(updatedGame), nil
}

func (s *GameService) DeleteGame(ctx context.Context,
	gameID uuid.UUID) error {

	if err := s.validateGameID(gameID); err != nil {
		return err
	}

	exists, repoError := s.repo.Exists(ctx, gameID)

	if repoError != nil {
		return repoError
	}

	if !exists {
		return errors.New("game not found")
	}

	repoError = s.repo.Delete(ctx, gameID)

	if repoError != nil {
		return repoError
	}

	return nil
}

func (s *GameService) validateUpdateData(gameDto dto.UpdateGameDto) error {
	if gameDto.ID == uuid.Nil {
		return errors.New("game ID is required")
	}

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

func (s *GameService) validateGameID(gameID uuid.UUID) error {
	if gameID == uuid.Nil {
		return errors.New("game ID is required")
	}

	return nil
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
