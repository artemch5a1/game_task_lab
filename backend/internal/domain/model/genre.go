package model

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

type Genre struct {
	ID    uuid.UUID
	Title string
}

func NewGenreWithValidate(title string) (*Genre, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, errors.New("genre title is required")
	}
	if len(title) > 200 {
		return nil, errors.New("genre title is too long")
	}

	return &Genre{
		ID:    uuid.New(),
		Title: title,
	}, nil
}

func (g *Genre) UpdateTitleWithValidate(title string) error {
	title = strings.TrimSpace(title)
	if title == "" {
		return errors.New("genre title is required")
	}
	if len(title) > 200 {
		return errors.New("genre title is too long")
	}

	g.Title = title
	return nil
}
