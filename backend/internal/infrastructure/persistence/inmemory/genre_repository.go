package inmemory

import (
	"context"
	"errors"

	"example/web-service-gin/internal/application/abstraction/repository"
	"example/web-service-gin/internal/domain/model"
	"example/web-service-gin/internal/infrastructure/persistence/data"

	"github.com/google/uuid"
)

// Проверка что реализуем интерфейс
var _ repository.GenreRepository = (*GenreRepository)(nil)

// GenreRepository in-memory реализация
type GenreRepository struct {
	data *data.Data
}

// NewGenreRepository создает новый in-memory репозиторий
func NewGenreRepository(store *data.Data) *GenreRepository {
	if store == nil {
		store = data.New()
	}
	return &GenreRepository{data: store}
}

func (r *GenreRepository) Create(ctx context.Context, genre *model.Genre) (*model.Genre, error) {
	if genre == nil {
		return nil, errors.New("genre cannot be nil")
	}
	if genre.ID == uuid.Nil {
		genre.ID = uuid.New()
	}

	r.data.Mu.Lock()
	defer r.data.Mu.Unlock()

	if _, exists := r.data.Genres[genre.ID]; exists {
		return nil, repository.ErrGenreAlreadyExists
	}

	r.data.Genres[genre.ID] = genre
	return genre, nil
}

func (r *GenreRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Genre, error) {
	if id == uuid.Nil {
		return nil, errors.New("genre ID cannot be empty")
	}

	r.data.Mu.RLock()
	defer r.data.Mu.RUnlock()

	genre, exists := r.data.Genres[id]
	if !exists {
		return nil, repository.ErrGenreNotFound
	}

	return genre, nil
}

func (r *GenreRepository) FindAll(ctx context.Context, limit, offset int) ([]*model.Genre, error) {
	r.data.Mu.RLock()
	defer r.data.Mu.RUnlock()

	all := make([]*model.Genre, 0, len(r.data.Genres))
	for _, g := range r.data.Genres {
		all = append(all, g)
	}

	start := offset
	if start > len(all) {
		start = len(all)
	}

	end := start + limit
	if end > len(all) {
		end = len(all)
	}
	if limit <= 0 {
		end = len(all)
	}

	return all[start:end], nil
}

func (r *GenreRepository) Update(ctx context.Context, genre *model.Genre) (*model.Genre, error) {
	if genre == nil {
		return nil, errors.New("genre cannot be nil")
	}
	if genre.ID == uuid.Nil {
		return nil, errors.New("genre ID cannot be empty")
	}

	r.data.Mu.Lock()
	defer r.data.Mu.Unlock()

	if _, exists := r.data.Genres[genre.ID]; !exists {
		return nil, repository.ErrGenreNotFound
	}

	r.data.Genres[genre.ID] = genre
	return genre, nil
}

func (r *GenreRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("genre ID cannot be empty")
	}

	r.data.Mu.Lock()
	defer r.data.Mu.Unlock()

	if _, exists := r.data.Genres[id]; !exists {
		return repository.ErrGenreNotFound
	}

	delete(r.data.Genres, id)
	return nil
}

func (r *GenreRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	if id == uuid.Nil {
		return false, errors.New("genre ID cannot be empty")
	}

	r.data.Mu.RLock()
	defer r.data.Mu.RUnlock()

	_, exists := r.data.Genres[id]
	return exists, nil
}

