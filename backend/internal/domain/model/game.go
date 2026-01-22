package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Game struct {
	ID          uuid.UUID
	Title       string
	Description string
	ReleaseDate time.Time
	GenreID     uuid.UUID
}

func (g *Game) UpdateGameWithValidate(title string,
	description string,
	releaseDate time.Time,
	genreId uuid.UUID) error {

	if releaseDate.After(time.Now().AddDate(0, 0, 30)) {
		return errors.New("release date cannot be more than 30 days in the future")
	}

	g.Title = title
	g.Description = description
	g.ReleaseDate = releaseDate
	g.GenreID = genreId

	return nil
}

func NewGameWithValidate(
	title string,
	description string,
	releaseDate time.Time,
	genreId uuid.UUID,
) (*Game, error) {

	if releaseDate.After(time.Now().AddDate(0, 0, 30)) {
		return nil, errors.New("release date cannot be more than 30 days in the future")
	}

	return &Game{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
		ReleaseDate: releaseDate,
		GenreID:     genreId,
	}, nil
}
