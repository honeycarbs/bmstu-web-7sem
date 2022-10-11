package account

import (
	"neatly/internal/model/account"
	"neatly/pkg/logging"
)

type mapper struct {
	logger logging.Logger
}

func New(logger logging.Logger) *mapper {
	return &mapper{logger: logger}
}

func (m *mapper) MapRegisterAccountDTO(dto account.RegisterAccountDTO) (account.Account, error) {
	phash, err := account.GeneratePasswordHash(dto.Password)
	if err != nil {
		m.logger.Info(err)
		return account.Account{}, err
	}

	m.logger.Info("Account password hashed")
	return account.Account{
		ID:           0,
		Name:         dto.Name,
		Username:     dto.Username,
		Email:        dto.Email,
		Password:     dto.Password,
		PasswordHash: phash,
	}, nil
}

func (m *mapper) MapLogInAccountDTO(dto account.LoginAccountDTO) account.Account {
	return account.Account{
		ID:           0,
		Name:         "",
		Username:     dto.Username,
		Email:        "",
		Password:     dto.Password,
		PasswordHash: "",
	}
}

func (m *mapper) MapAccountWithTokenDTO(token string, a account.Account) account.WithTokenDTO {
	return account.WithTokenDTO{
		Token:    token,
		Name:     a.Name,
		Username: a.Username,
		Email:    a.Email,
	}
}
