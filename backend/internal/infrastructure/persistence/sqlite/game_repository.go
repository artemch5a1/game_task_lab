package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"example/web-service-gin/internal/application/abstraction/repository"
	"example/web-service-gin/internal/constants"
	"example/web-service-gin/internal/domain/model"

	"github.com/google/uuid"
)

// compile-time check
var _ repository.GameRepository = (*GameRepository)(nil)

type GameRepository struct {
	db *sql.DB
}

func NewGameRepository(db *sql.DB) *GameRepository {
	return &GameRepository{db: db}
}

func (r *GameRepository) Create(ctx context.Context, game *model.Game) (*model.Game, error) {
	if game == nil {
		return nil, errors.New("game cannot be nil")
	}
	if game.ID == uuid.Nil {
		game.ID = uuid.New()
	}

	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO games (id, title, description, release_date, genre_id) VALUES (?, ?, ?, ?, ?)`,
		game.ID.String(),
		game.Title,
		game.Description,
		game.ReleaseDate.UTC().Format(time.RFC3339Nano),
		game.GenreID.String(),
	)
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "UNIQUE constraint failed") {
			return nil, repository.ErrAlreadyExists
		}
		if strings.Contains(msg, "FOREIGN KEY constraint failed") {
			return nil, errors.New(constants.ErrGenreNotFound)
		}
		return nil, fmt.Errorf("insert game: %w", err)
	}

	return game, nil
}

func (r *GameRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Game, error) {
	if id == uuid.Nil {
		return nil, errors.New("game ID cannot be empty")
	}

	var (
		idStr, title, description, releaseDateStr, genreIDStr string
	)

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, title, description, release_date, genre_id FROM games WHERE id = ?`,
		id.String(),
	).Scan(&idStr, &title, &description, &releaseDateStr, &genreIDStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("select game: %w", err)
	}

	gameID, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("parse game id from db: %w", err)
	}
	genreID, err := uuid.Parse(genreIDStr)
	if err != nil {
		return nil, fmt.Errorf("parse genre_id from db: %w", err)
	}
	releaseDate, err := time.Parse(time.RFC3339Nano, releaseDateStr)
	if err != nil {
		return nil, fmt.Errorf("parse release_date from db: %w", err)
	}

	return &model.Game{
		ID:          gameID,
		Title:       title,
		Description: description,
		ReleaseDate: releaseDate,
		GenreID:     genreID,
	}, nil
}

func (r *GameRepository) FindAll(ctx context.Context, limit, offset int) ([]*model.Game, error) {
	query := `SELECT id, title, description, release_date, genre_id FROM games ORDER BY release_date DESC`
	args := []any{}

	if limit > 0 {
		query += ` LIMIT ? OFFSET ?`
		args = append(args, limit, offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("select games: %w", err)
	}
	defer rows.Close()

	var res []*model.Game
	for rows.Next() {
		var idStr, title, description, releaseDateStr, genreIDStr string
		if err := rows.Scan(&idStr, &title, &description, &releaseDateStr, &genreIDStr); err != nil {
			return nil, fmt.Errorf("scan game: %w", err)
		}
		gameID, err := uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("parse game id from db: %w", err)
		}
		genreID, err := uuid.Parse(genreIDStr)
		if err != nil {
			return nil, fmt.Errorf("parse genre_id from db: %w", err)
		}
		releaseDate, err := time.Parse(time.RFC3339Nano, releaseDateStr)
		if err != nil {
			return nil, fmt.Errorf("parse release_date from db: %w", err)
		}
		res = append(res, &model.Game{
			ID:          gameID,
			Title:       title,
			Description: description,
			ReleaseDate: releaseDate,
			GenreID:     genreID,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate games: %w", err)
	}

	return res, nil
}

func (r *GameRepository) Update(ctx context.Context, game *model.Game) (*model.Game, error) {
	if game == nil {
		return nil, errors.New("game cannot be nil")
	}
	if game.ID == uuid.Nil {
		return nil, errors.New("game ID cannot be empty")
	}

	result, err := r.db.ExecContext(
		ctx,
		`UPDATE games SET title = ?, description = ?, release_date = ?, genre_id = ? WHERE id = ?`,
		game.Title,
		game.Description,
		game.ReleaseDate.UTC().Format(time.RFC3339Nano),
		game.GenreID.String(),
		game.ID.String(),
	)
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "FOREIGN KEY constraint failed") {
			return nil, errors.New(constants.ErrGenreNotFound)
		}
		return nil, fmt.Errorf("update game: %w", err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return nil, repository.ErrNotFound
	}

	return game, nil
}

func (r *GameRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("game ID cannot be empty")
	}

	result, err := r.db.ExecContext(ctx, `DELETE FROM games WHERE id = ?`, id.String())
	if err != nil {
		return fmt.Errorf("delete game: %w", err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (r *GameRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	if id == uuid.Nil {
		return false, errors.New("game ID cannot be empty")
	}

	var one int
	err := r.db.QueryRowContext(ctx, `SELECT 1 FROM games WHERE id = ?`, id.String()).Scan(&one)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("exists game: %w", err)
	}
	return true, nil
}

