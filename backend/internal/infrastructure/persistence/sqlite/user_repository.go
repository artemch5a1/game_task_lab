package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"example/web-service-gin/internal/application/abstraction/repository"
	"example/web-service-gin/internal/domain/model"
	specifictype "example/web-service-gin/internal/domain/specific_type"

	"github.com/google/uuid"
)

var _ repository.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	if user == nil {
		return nil, errors.New("user cannot be nil")
	}
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO users (id, username, password, user_role) VALUES (?, ?, ?, ?)`,
		user.ID.String(),
		user.Username,
		user.Password,
		string(user.UserRole),
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, repository.ErrUserAlreadyExists
		}
		return nil, fmt.Errorf("insert user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	if id == uuid.Nil {
		return nil, errors.New("user ID cannot be empty")
	}
	return r.findOne(ctx, `SELECT id, username, password, user_role FROM users WHERE id = ?`, id.String())
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	return r.findOne(ctx, `SELECT id, username, password, user_role FROM users WHERE username = ?`, username)
}

func (r *UserRepository) FindAll(ctx context.Context, limit, offset int) ([]*model.User, error) {
	query := `SELECT id, username, password, user_role FROM users ORDER BY username`
	args := []any{}
	if limit > 0 {
		query += ` LIMIT ? OFFSET ?`
		args = append(args, limit, offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("select users: %w", err)
	}
	defer rows.Close()

	var res []*model.User
	for rows.Next() {
		var idStr, username, password, roleStr string
		if err := rows.Scan(&idStr, &username, &password, &roleStr); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		uid, err := uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("parse user id from db: %w", err)
		}
		res = append(res, &model.User{
			ID:       uid,
			Username: username,
			Password: password,
			UserRole: specifictype.UserRole(roleStr),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}
	return res, nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) (*model.User, error) {
	if user == nil {
		return nil, errors.New("user cannot be nil")
	}
	if user.ID == uuid.Nil {
		return nil, errors.New("user ID cannot be empty")
	}

	result, err := r.db.ExecContext(
		ctx,
		`UPDATE users SET username = ?, password = ?, user_role = ? WHERE id = ?`,
		user.Username,
		user.Password,
		string(user.UserRole),
		user.ID.String(),
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, repository.ErrUserAlreadyExists
		}
		return nil, fmt.Errorf("update user: %w", err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return nil, repository.ErrUserNotFound
	}
	return user, nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("user ID cannot be empty")
	}

	result, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = ?`, id.String())
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return repository.ErrUserNotFound
	}
	return nil
}

func (r *UserRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	if id == uuid.Nil {
		return false, errors.New("user ID cannot be empty")
	}
	var one int
	err := r.db.QueryRowContext(ctx, `SELECT 1 FROM users WHERE id = ?`, id.String()).Scan(&one)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("exists user: %w", err)
	}
	return true, nil
}

func (r *UserRepository) findOne(ctx context.Context, query string, arg any) (*model.User, error) {
	var idStr, username, password, roleStr string
	err := r.db.QueryRowContext(ctx, query, arg).Scan(&idStr, &username, &password, &roleStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrUserNotFound
		}
		return nil, fmt.Errorf("select user: %w", err)
	}
	uid, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("parse user id from db: %w", err)
	}
	return &model.User{
		ID:       uid,
		Username: username,
		Password: password,
		UserRole: specifictype.UserRole(roleStr),
	}, nil
}

