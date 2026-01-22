package aggregate

import (
	"errors"
	"example/web-service-gin/internal/domain/model"
	"time"

	"github.com/google/uuid"
)

type GameDetailsAggregate struct {
	Game    *model.Game
	Genre   *model.Genre
	Ratings []*model.UserRating

	averageRating float64
	ratingCount   int
}

func NewGameDetailsAggregate(game *model.Game, genre *model.Genre, ratings []*model.UserRating) *GameDetailsAggregate {
	agg := &GameDetailsAggregate{
		Game:    game,
		Genre:   genre,
		Ratings: ratings,
	}
	agg.calculateDerivedFields()
	return agg
}

func (g *GameDetailsAggregate) AddRating(userID uuid.UUID, rating int) error {

	if rating < 1 || rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}

	for _, i := range g.Ratings {
		if i.UserID == userID {
			return errors.New("Пользователь уже оценил игру")
		}
	}

	newRating := &model.UserRating{
		ID:        uuid.New(),
		UserID:    userID,
		GameID:    g.Game.ID,
		Rating:    rating,
		CreatedAt: time.Now().UTC(),
	}

	g.Ratings = append(g.Ratings, newRating)

	g.calculateDerivedFields()

	return nil
}

func (g *GameDetailsAggregate) UpdateRating(userID uuid.UUID, newRating int) error {
	for _, r := range g.Ratings {
		if r.UserID == userID {
			r.Rating = newRating
			r.CreatedAt = time.Now()
			g.calculateDerivedFields()
			return nil
		}
	}
	return errors.New("user hasn't rated this game yet")
}

func (g *GameDetailsAggregate) CalculateAverageRating() float64 {
	if len(g.Ratings) == 0 {
		return 0.0
	}

	var sum int
	for _, r := range g.Ratings {
		sum += r.Rating
	}
	return float64(sum) / float64(len(g.Ratings))
}

func (g *GameDetailsAggregate) GetRatingCount() int {
	return g.ratingCount
}

func (g *GameDetailsAggregate) calculateDerivedFields() {
	g.averageRating = g.CalculateAverageRating()
	g.ratingCount = len(g.Ratings)
}
