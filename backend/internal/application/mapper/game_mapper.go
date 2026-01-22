package mapper

import (
	"errors"
	"example/web-service-gin/internal/application/dto"
	"example/web-service-gin/internal/domain/aggregate"
	"example/web-service-gin/internal/domain/model"
)

type GameMapper struct{}

func NewGameMapper() *GameMapper {
	return &GameMapper{}
}

func (m *GameMapper) ToGameDto(game *model.Game) *dto.GameDto {
	if game == nil {
		return nil
	}

	return &dto.GameDto{
		ID:          game.ID,
		Title:       game.Title,
		Description: game.Description,
		ReleaseDate: game.ReleaseDate,
		GenreID:     game.GenreID,
	}
}

func (m *GameMapper) ToGameDtoWithStats(agg *aggregate.GameDetailsAggregate) *dto.GameDtoWithStats {
	if agg == nil || agg.Game == nil {
		return nil
	}

	return &dto.GameDtoWithStats{
		ID:            agg.Game.ID,
		Title:         agg.Game.Title,
		Description:   agg.Game.Description,
		ReleaseDate:   agg.Game.ReleaseDate,
		GenreID:       agg.Game.GenreID,
		AverageRating: agg.CalculateAverageRating(),
		RatingCount:   len(agg.Ratings),
	}
}

func (m *GameMapper) ToGameDtoSlice(games []*model.Game) []*dto.GameDto {
	if games == nil {
		return nil
	}

	result := make([]*dto.GameDto, len(games))
	for i, game := range games {
		result[i] = m.ToGameDto(game)
	}
	return result
}

func (m *GameMapper) FromCreateGameDto(dto *dto.CreateGameDto) (*model.Game, error) {
	if dto == nil {
		return nil, errors.New("create game dto cannot be nil")
	}

	return model.NewGameWithValidate(
		dto.Title,
		dto.Description,
		dto.ReleaseDate,
		dto.GenreID,
	)
}

func (m *GameMapper) FromUpdateGameDto(game *model.Game, dto *dto.UpdateGameDto) error {
	if game == nil {
		return errors.New("game cannot be nil")
	}

	if dto == nil {
		return errors.New("update game dto cannot be nil")
	}

	if game.ID != dto.ID {
		return errors.New("game ID mismatch")
	}

	return game.UpdateGameWithValidate(
		dto.Title,
		dto.Description,
		dto.ReleaseDate,
		dto.GenreID,
	)
}
