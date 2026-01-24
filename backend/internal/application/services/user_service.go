package services

import (
	"context"
	"errors"
	"math"
	"strings"

	"example/web-service-gin/internal/application/abstraction/repository"
	"example/web-service-gin/internal/application/dto"
	"example/web-service-gin/internal/application/mapper"
	"example/web-service-gin/internal/constants"
	specifictype "example/web-service-gin/internal/domain/specific_type"

	"github.com/google/uuid"
)

type UserService struct {
	repo       repository.UserRepository
	userMapper *mapper.UserMapper
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo:       repo,
		userMapper: mapper.NewUserMapper(),
	}
}

func (s *UserService) CreateUser(ctx context.Context, in dto.CreateUserDto) (*dto.UserDto, error) {
	if err := s.validateUserData(in.Username, in.Password, in.UserRole); err != nil {
		return nil, err
	}

	u, err := s.userMapper.FromCreateUserDto(&in)
	if err != nil {
		return nil, err
	}

	created, err := s.repo.Create(ctx, u)
	if err != nil {
		return nil, err
	}

	return s.userMapper.ToUserDto(created), nil
}

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*dto.UserDto, error) {
	if id == uuid.Nil {
		return nil, errors.New(constants.ErrUserIDRequired)
	}
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.userMapper.ToUserDto(u), nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]*dto.UserDto, error) {
	users, err := s.repo.FindAll(ctx, math.MaxInt, 0)
	if err != nil {
		return nil, err
	}
	return s.userMapper.ToUserDtoSlice(users), nil
}

func (s *UserService) UpdateUser(ctx context.Context, in dto.UpdateUserDto) (*dto.UserDto, error) {
	if in.ID == uuid.Nil {
		return nil, errors.New(constants.ErrUserIDRequired)
	}
	if err := s.validateUserData(in.Username, in.Password, in.UserRole); err != nil {
		return nil, err
	}

	existing, err := s.repo.FindByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	if err := s.userMapper.FromUpdateUserDto(existing, &in); err != nil {
		return nil, err
	}

	updated, err := s.repo.Update(ctx, existing)
	if err != nil {
		return nil, err
	}
	return s.userMapper.ToUserDto(updated), nil
}

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New(constants.ErrUserIDRequired)
	}
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return repository.ErrUserNotFound
	}
	return s.repo.Delete(ctx, id)
}

func (s *UserService) Authenticate(ctx context.Context, username, password string) (bool, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return false, errors.New(constants.ErrUserUsernameEmpty)
	}
	if strings.TrimSpace(password) == "" {
		return false, errors.New(constants.ErrUserPasswordEmpty)
	}

	u, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return false, nil
		}
		return false, err
	}

	return u.Password == password, nil
}

func (s *UserService) validateUserData(username, password string, role specifictype.UserRole) error {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	if username == "" {
		return errors.New(constants.ErrUserUsernameEmpty)
	}
	if password == "" {
		return errors.New(constants.ErrUserPasswordEmpty)
	}
	if len(username) > 200 || len(password) > 200 {
		return errors.New(constants.ErrValidationTitleLength)
	}
	if role != specifictype.RoleUser && role != specifictype.RoleAdmin {
		return errors.New(constants.ErrUserRoleInvalid)
	}
	return nil
}

