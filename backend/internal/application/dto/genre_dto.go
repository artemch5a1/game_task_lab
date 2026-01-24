package dto

import "github.com/google/uuid"

type GenreDto struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}

type CreateGenreDto struct {
	Title string `json:"title" validate:"required,min=1,max=200"`
}

type UpdateGenreDto struct {
	ID    uuid.UUID `json:"id" validate:"required"`
	Title string    `json:"title" validate:"required,min=1,max=200"`
}

