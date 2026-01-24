package mapper

import (
	"errors"

	"example/web-service-gin/internal/application/dto"
	"example/web-service-gin/internal/constants"
	"example/web-service-gin/internal/domain/model"

	"github.com/google/uuid"
)

type UserMapper struct{}

func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

func (m *UserMapper) ToUserDto(user *model.User) *dto.UserDto {
	if user == nil {
		return nil
	}
	return &dto.UserDto{
		ID:       user.ID,
		Username: user.Username,
		UserRole: user.UserRole,
	}
}

func (m *UserMapper) ToUserDtoSlice(users []*model.User) []*dto.UserDto {
	if users == nil {
		return []*dto.UserDto{}
	}
	res := make([]*dto.UserDto, len(users))
	for i, u := range users {
		res[i] = m.ToUserDto(u)
	}
	return res
}

func (m *UserMapper) FromCreateUserDto(in *dto.CreateUserDto) (*model.User, error) {
	if in == nil {
		return nil, errors.New(constants.ErrInvalidData)
	}
	return &model.User{
		ID:       uuid.New(),
		Username: in.Username,
		Password: in.Password,
		UserRole: in.UserRole,
	}, nil
}

func (m *UserMapper) FromUpdateUserDto(user *model.User, in *dto.UpdateUserDto) error {
	if user == nil || in == nil {
		return errors.New(constants.ErrInvalidData)
	}
	if user.ID != in.ID {
		return errors.New(constants.ErrIDMismatch)
	}
	user.Username = in.Username
	user.Password = in.Password
	user.UserRole = in.UserRole
	return nil
}

