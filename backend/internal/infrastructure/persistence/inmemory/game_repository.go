package inmemory

import (
	"context"
	"errors"

	"example/web-service-gin/internal/application/abstraction/repository"
	"example/web-service-gin/internal/infrastructure/persistence/data"
	"example/web-service-gin/internal/domain/model"

	"github.com/google/uuid"
)

// Проверка что реализуем интерфейс
var _ repository.GameRepository = (*GameRepository)(nil)

// GameRepository in-memory реализация
type GameRepository struct {
	data *data.Data
}

// NewGameRepository создает новый in-memory репозиторий
func NewGameRepository(store *data.Data) *GameRepository {
	if store == nil {
		store = data.New()
	}
	return &GameRepository{
		data: store,
	}
}

// Create создает новую игру
func (r *GameRepository) Create(ctx context.Context, game *model.Game) (*model.Game, error) {
	if game == nil {
		return nil, errors.New("game cannot be nil")
	}

	// Генерируем ID если не указан
	if game.ID == uuid.Nil {
		game.ID = uuid.New()
	}

	r.data.Mu.Lock()
	defer r.data.Mu.Unlock()

	// Проверяем уникальность ID
	if _, exists := r.data.Games[game.ID]; exists {
		return nil, repository.ErrAlreadyExists
	}

	// Сохраняем
	r.data.Games[game.ID] = game

	return game, nil
}

// FindByID ищет игру по ID
func (r *GameRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Game, error) {
	if id == uuid.Nil {
		return nil, errors.New("game ID cannot be empty")
	}

	r.data.Mu.RLock()
	defer r.data.Mu.RUnlock()

	game, exists := r.data.Games[id]
	if !exists {
		return nil, repository.ErrNotFound
	}

	return game, nil
}

// FindAll возвращает все игры с пагинацией
func (r *GameRepository) FindAll(ctx context.Context, limit, offset int) ([]*model.Game, error) {
	r.data.Mu.RLock()
	defer r.data.Mu.RUnlock()

	// Получаем все игры
	allGames := make([]*model.Game, 0, len(r.data.Games))
	for _, game := range r.data.Games {
		allGames = append(allGames, game)
	}

	// Сортируем по дате создания (если есть)
	// Для простоты сортируем по ID

	// Применяем пагинацию
	start := offset
	if start > len(allGames) {
		start = len(allGames)
	}

	end := start + limit
	if end > len(allGames) {
		end = len(allGames)
	}

	if limit <= 0 {
		end = len(allGames) // если limit 0 или отрицательный - возвращаем все
	}

	return allGames[start:end], nil
}

// Update обновляет игру
func (r *GameRepository) Update(ctx context.Context, game *model.Game) (*model.Game, error) {
	if game == nil {
		return nil, errors.New("game cannot be nil")
	}

	if game.ID == uuid.Nil {
		return nil, errors.New("game ID cannot be empty")
	}

	r.data.Mu.Lock()
	defer r.data.Mu.Unlock()

	// Проверяем существование
	_, exists := r.data.Games[game.ID]
	if !exists {
		return nil, repository.ErrNotFound
	}

	// Обновляем
	r.data.Games[game.ID] = game

	return game, nil
}

// Delete удаляет игру
func (r *GameRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("game ID cannot be empty")
	}

	r.data.Mu.Lock()
	defer r.data.Mu.Unlock()

	// Проверяем существование
	_, exists := r.data.Games[id]
	if !exists {
		return repository.ErrNotFound
	}

	// Удаляем
	delete(r.data.Games, id)

	return nil
}

// Exists проверяет существование игры
func (r *GameRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	if id == uuid.Nil {
		return false, errors.New("game ID cannot be empty")
	}

	r.data.Mu.RLock()
	defer r.data.Mu.RUnlock()

	_, exists := r.data.Games[id]
	return exists, nil
}

// ========== Дополнительные методы ==========

// Count возвращает количество игр
func (r *GameRepository) Count(ctx context.Context) (int, error) {
	r.data.Mu.RLock()
	defer r.data.Mu.RUnlock()

	return len(r.data.Games), nil
}

// Clear очищает все данные (для тестов)
func (r *GameRepository) Clear() {
	r.data.Mu.Lock()
	defer r.data.Mu.Unlock()

	r.data.Games = make(map[uuid.UUID]*model.Game)
}

// Seed добавляет тестовые данные
func (r *GameRepository) Seed(games []*model.Game) error {
	r.data.Mu.Lock()
	defer r.data.Mu.Unlock()

	for _, game := range games {
		if game.ID == uuid.Nil {
			game.ID = uuid.New()
		}
		r.data.Games[game.ID] = game
	}

	return nil
}
