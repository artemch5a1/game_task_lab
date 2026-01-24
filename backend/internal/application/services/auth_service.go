package services

import (
	"context"
	"errors"
	"strings"

	appauth "example/web-service-gin/internal/application/abstraction/auth"
	"example/web-service-gin/internal/application/abstraction/repository"
	"example/web-service-gin/internal/constants"
	"example/web-service-gin/internal/domain/model"
	specifictype "example/web-service-gin/internal/domain/specific_type"

	"github.com/google/uuid"
)

type AuthService struct {
	users repository.UserRepository
	tokens appauth.TokenProvider
}

func NewAuthService(users repository.UserRepository, tokenProvider appauth.TokenProvider) *AuthService {
	return &AuthService{
		users: users,
		tokens: tokenProvider,
	}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	if username == "" || password == "" {
		return "", errors.New(constants.ErrUnauthorized)
	}

	u, err := s.users.FindByUsername(ctx, username)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return "", errors.New(constants.ErrUnauthorized)
		}
		return "", err
	}

	// TODO: заменить на безопасный хэш (bcrypt/argon2) на следующем шаге.
	if u.Password != password {
		return "", errors.New(constants.ErrUnauthorized)
	}

	return s.tokens.Issue(ctx, u.ID, u.UserRole)
}

func (s *AuthService) Register(ctx context.Context, username, password string) (string, error) {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	if username == "" || password == "" {
		return "", errors.New(constants.ErrInvalidData)
	}
	if len(username) > 200 || len(password) > 200 {
		return "", errors.New(constants.ErrValidationTitleLength)
	}

	u := &model.User{
		ID:       uuid.New(),
		Username: username,
		Password: password, // TODO: hash (bcrypt/argon2)
		UserRole: specifictype.RoleUser,
	}

	created, err := s.users.Create(ctx, u)
	if err != nil {
		return "", err
	}

	return s.tokens.Issue(ctx, created.ID, created.UserRole)
}

