package services

import (
	"context"
	"errors"
	"strings"

	appauth "example/web-service-gin/internal/application/abstraction/auth"
	"example/web-service-gin/internal/application/abstraction/repository"
	"example/web-service-gin/internal/constants"
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

