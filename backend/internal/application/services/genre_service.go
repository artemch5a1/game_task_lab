package services

import (
	"context"
	"errors"
	"math"
	"strings"

	"example/web-service-gin/internal/application/abstraction/repository"
	"example/web-service-gin/internal/application/dto"
	"example/web-service-gin/internal/application/mapper"
	"example/web-service-gin/internal/constants"

	"github.com/google/uuid"
)

type GenreService struct {
	repo        repository.GenreRepository
	genreMapper *mapper.GenreMapper
}

func NewGenreService(repo repository.GenreRepository) *GenreService {
	return &GenreService{
		repo:        repo,
		genreMapper: mapper.NewGenreMapper(),
	}
}

func (s *GenreService) CreateGenre(ctx context.Context, in dto.CreateGenreDto) (*dto.GenreDto, error) {
	if err := s.validateGenreTitle(in.Title); err != nil {
		return nil, err
	}

	genre, err := s.genreMapper.FromCreateGenreDto(&in)
	if err != nil {
		return nil, err
	}

	created, err := s.repo.Create(ctx, genre)
	if err != nil {
		return nil, err
	}

	return s.genreMapper.ToGenreDto(created), nil
}

func (s *GenreService) GetGenreByID(ctx context.Context, id uuid.UUID) (*dto.GenreDto, error) {
	if id == uuid.Nil {
		return nil, errors.New(constants.ErrGenreIDRequired)
	}

	genre, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.genreMapper.ToGenreDto(genre), nil
}

func (s *GenreService) GetAllGenres(ctx context.Context) ([]*dto.GenreDto, error) {
	genres, err := s.repo.FindAll(ctx, math.MaxInt, 0)
	if err != nil {
		return nil, err
	}

	return s.genreMapper.ToGenreDtoSlice(genres), nil
}

func (s *GenreService) UpdateGenre(ctx context.Context, in dto.UpdateGenreDto) (*dto.GenreDto, error) {
	if in.ID == uuid.Nil {
		return nil, errors.New(constants.ErrGenreIDRequired)
	}
	if err := s.validateGenreTitle(in.Title); err != nil {
		return nil, err
	}

	existing, err := s.repo.FindByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	if err := s.genreMapper.FromUpdateGenreDto(existing, &in); err != nil {
		return nil, err
	}

	updated, err := s.repo.Update(ctx, existing)
	if err != nil {
		return nil, err
	}

	return s.genreMapper.ToGenreDto(updated), nil
}

func (s *GenreService) DeleteGenre(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New(constants.ErrGenreIDRequired)
	}

	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return repository.ErrGenreNotFound
	}

	return s.repo.Delete(ctx, id)
}

func (s *GenreService) validateGenreTitle(title string) error {
	title = strings.TrimSpace(title)
	if title == "" {
		return errors.New(constants.ErrGenreTitleEmpty)
	}
	if len(title) > 200 {
		return errors.New(constants.ErrValidationTitleLength)
	}
	return nil
}

