package auth

import (
	"neatly/internal/model/auth"
	"neatly/pkg/logging"
)

type mapper struct {
	logger logging.Logger
}

func New(logger logging.Logger) *mapper {
	return &mapper{logger: logger}
}

func (m *mapper) MapRegisterAccountDTO(dto auth.RegisterAccountDTO) (auth.Account, error) {
	phash, err := auth.GeneratePasswordHash(dto.Password)
	if err != nil {
		m.logger.Info(err)
		return auth.Account{}, err
	}

	m.logger.Info("Account password hashed")
	return auth.Account{
		ID:           0,
		Name:         dto.Name,
		Username:     dto.Username,
		Email:        dto.Email,
		Password:     "",
		PasswordHash: phash,
	}, nil
}

func (m *mapper) MapLogInAccountDTO(dto auth.LoginAccountDTO) auth.Account {
	return auth.Account{
		ID:           0,
		Name:         "",
		Username:     dto.Username,
		Email:        "",
		Password:     dto.Password,
		PasswordHash: "",
	}
}

func (m *mapper) MapJwtDTO(token string) auth.JwtDTO {
	return auth.JwtDTO{
		Token: token,
	}
}
