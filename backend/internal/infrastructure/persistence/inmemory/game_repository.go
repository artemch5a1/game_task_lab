package inmemory

import (
	"context"
	"errors"
	"sync"

	"example/web-service-gin/internal/application/abstraction/repository"
	"example/web-service-gin/internal/domain/model"

	"github.com/google/uuid"
)

// Проверка что реализуем интерфейс
var _ repository.GameRepository = (*GameRepository)(nil)

// GameRepository in-memory реализация
type GameRepository struct {
	games map[uuid.UUID]*model.Game
	mu    sync.RWMutex
}

// NewGameRepository создает новый in-memory репозиторий
func NewGameRepository() *GameRepository {
	return &GameRepository{
		games: make(map[uuid.UUID]*model.Game),
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

	r.mu.Lock()
	defer r.mu.Unlock()

	// Проверяем уникальность ID
	if _, exists := r.games[game.ID]; exists {
		return nil, repository.ErrAlreadyExists
	}

	// Сохраняем
	r.games[game.ID] = game

	return game, nil
}

// FindByID ищет игру по ID
func (r *GameRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Game, error) {
	if id == uuid.Nil {
		return nil, errors.New("game ID cannot be empty")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	game, exists := r.games[id]
	if !exists {
		return nil, repository.ErrNotFound
	}

	return game, nil
}

// FindAll возвращает все игры с пагинацией
func (r *GameRepository) FindAll(ctx context.Context, limit, offset int) ([]*model.Game, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Получаем все игры
	allGames := make([]*model.Game, 0, len(r.games))
	for _, game := range r.games {
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

	r.mu.Lock()
	defer r.mu.Unlock()

	// Проверяем существование
	_, exists := r.games[game.ID]
	if !exists {
		return nil, repository.ErrNotFound
	}

	// Обновляем
	r.games[game.ID] = game

	return game, nil
}

// Delete удаляет игру
func (r *GameRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("game ID cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Проверяем существование
	_, exists := r.games[id]
	if !exists {
		return repository.ErrNotFound
	}

	// Удаляем
	delete(r.games, id)

	return nil
}

// Exists проверяет существование игры
func (r *GameRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	if id == uuid.Nil {
		return false, errors.New("game ID cannot be empty")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.games[id]
	return exists, nil
}

// ========== Дополнительные методы ==========

// Count возвращает количество игр
func (r *GameRepository) Count(ctx context.Context) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.games), nil
}

// Clear очищает все данные (для тестов)
func (r *GameRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.games = make(map[uuid.UUID]*model.Game)
}

// Seed добавляет тестовые данные
func (r *GameRepository) Seed(games []*model.Game) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, game := range games {
		if game.ID == uuid.Nil {
			game.ID = uuid.New()
		}
		r.games[game.ID] = game
	}

	return nil
}
