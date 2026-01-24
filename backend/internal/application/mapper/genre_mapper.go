package mapper

import (
	"errors"

	"example/web-service-gin/internal/application/dto"
	"example/web-service-gin/internal/constants"
	"example/web-service-gin/internal/domain/model"
)

type GenreMapper struct{}

func NewGenreMapper() *GenreMapper {
	return &GenreMapper{}
}

func (m *GenreMapper) ToGenreDto(genre *model.Genre) *dto.GenreDto {
	if genre == nil {
		return nil
	}
	return &dto.GenreDto{
		ID:    genre.ID,
		Title: genre.Title,
	}
}

func (m *GenreMapper) ToGenreDtoSlice(genres []*model.Genre) []*dto.GenreDto {
	if genres == nil {
		return nil
	}

	result := make([]*dto.GenreDto, len(genres))
	for i, g := range genres {
		result[i] = m.ToGenreDto(g)
	}
	return result
}

func (m *GenreMapper) FromCreateGenreDto(in *dto.CreateGenreDto) (*model.Genre, error) {
	if in == nil {
		return nil, errors.New(constants.ErrInvalidData)
	}
	return model.NewGenreWithValidate(in.Title)
}

func (m *GenreMapper) FromUpdateGenreDto(genre *model.Genre, in *dto.UpdateGenreDto) error {
	if genre == nil || in == nil {
		return errors.New(constants.ErrInvalidData)
	}
	if genre.ID != in.ID {
		return errors.New(constants.ErrIDMismatch)
	}
	return genre.UpdateTitleWithValidate(in.Title)
}

