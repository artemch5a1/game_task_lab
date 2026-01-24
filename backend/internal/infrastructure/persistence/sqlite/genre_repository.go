package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"example/web-service-gin/internal/application/abstraction/repository"
	"example/web-service-gin/internal/domain/model"

	"github.com/google/uuid"
)

// compile-time check
var _ repository.GenreRepository = (*GenreRepository)(nil)

type GenreRepository struct {
	db *sql.DB
}

func NewGenreRepository(db *sql.DB) *GenreRepository {
	return &GenreRepository{db: db}
}

func (r *GenreRepository) Create(ctx context.Context, genre *model.Genre) (*model.Genre, error) {
	if genre == nil {
		return nil, errors.New("genre cannot be nil")
	}
	if genre.ID == uuid.Nil {
		genre.ID = uuid.New()
	}

	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO genres (id, title) VALUES (?, ?)`,
		genre.ID.String(),
		genre.Title,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, repository.ErrGenreAlreadyExists
		}
		return nil, fmt.Errorf("insert genre: %w", err)
	}

	return genre, nil
}

func (r *GenreRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Genre, error) {
	if id == uuid.Nil {
		return nil, errors.New("genre ID cannot be empty")
	}

	var idStr, title string
	err := r.db.QueryRowContext(ctx, `SELECT id, title FROM genres WHERE id = ?`, id.String()).
		Scan(&idStr, &title)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrGenreNotFound
		}
		return nil, fmt.Errorf("select genre: %w", err)
	}

	uid, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("parse genre id from db: %w", err)
	}

	return &model.Genre{ID: uid, Title: title}, nil
}

func (r *GenreRepository) FindAll(ctx context.Context, limit, offset int) ([]*model.Genre, error) {
	query := `SELECT id, title FROM genres ORDER BY title`
	args := []any{}

	if limit > 0 {
		query += ` LIMIT ? OFFSET ?`
		args = append(args, limit, offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("select genres: %w", err)
	}
	defer rows.Close()

	var res []*model.Genre
	for rows.Next() {
		var idStr, title string
		if err := rows.Scan(&idStr, &title); err != nil {
			return nil, fmt.Errorf("scan genre: %w", err)
		}
		uid, err := uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("parse genre id from db: %w", err)
		}
		res = append(res, &model.Genre{ID: uid, Title: title})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate genres: %w", err)
	}

	return res, nil
}

func (r *GenreRepository) Update(ctx context.Context, genre *model.Genre) (*model.Genre, error) {
	if genre == nil {
		return nil, errors.New("genre cannot be nil")
	}
	if genre.ID == uuid.Nil {
		return nil, errors.New("genre ID cannot be empty")
	}

	result, err := r.db.ExecContext(ctx, `UPDATE genres SET title = ? WHERE id = ?`, genre.Title, genre.ID.String())
	if err != nil {
		return nil, fmt.Errorf("update genre: %w", err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return nil, repository.ErrGenreNotFound
	}

	return genre, nil
}

func (r *GenreRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("genre ID cannot be empty")
	}

	result, err := r.db.ExecContext(ctx, `DELETE FROM genres WHERE id = ?`, id.String())
	if err != nil {
		return fmt.Errorf("delete genre: %w", err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return repository.ErrGenreNotFound
	}

	return nil
}

func (r *GenreRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	if id == uuid.Nil {
		return false, errors.New("genre ID cannot be empty")
	}

	var one int
	err := r.db.QueryRowContext(ctx, `SELECT 1 FROM genres WHERE id = ?`, id.String()).Scan(&one)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("exists genre: %w", err)
	}
	return true, nil
}

