package user

import (
	"neatly/internal/model/user"
	"neatly/pkg/logging"
)

type mapper struct {
	logger logging.Logger
}

func New(logger logging.Logger) *mapper {
	return &mapper{logger: logger}
}

func (m *mapper) MapUserRegisterDTO(dto user.RegisterUserDTO) (user.User, error) {
	phash, err := user.GeneratePasswordHash(dto.Password)
	if err != nil {
		m.logger.Info(err)
		return user.User{}, err
	}

	m.logger.Info("User password hashed")
	return user.User{
		ID:           0,
		Name:         dto.Name,
		Username:     dto.Username,
		Email:        dto.Email,
		Password:     "",
		PasswordHash: phash,
	}, nil
}

func (m *mapper) MapUserLogInUserDTO(dto user.LoginUserDTO) user.User {
	return user.User{
		ID:           0,
		Name:         "",
		Username:     dto.Username,
		Email:        "",
		Password:     dto.Password,
		PasswordHash: "",
	}
}
