package mapper

import (
	"neatly/internal/model"
	"neatly/internal/model/dto"
	"neatly/pkg/logging"
)

type AccountMapper struct {
	logger logging.Logger
}

func NewAccountMapper(logger logging.Logger) *AccountMapper {
	return &AccountMapper{logger: logger}
}

func (m *AccountMapper) MapRegisterAccountDTO(dto dto.RegisterAccountDTO) (model.Account, error) {
	phash, err := model.GeneratePasswordHash(dto.Password)
	if err != nil {
		m.logger.Info(err)
		return model.Account{}, err
	}

	m.logger.Info("Account password hashed")
	return model.Account{
		ID:           0,
		Name:         dto.Name,
		Username:     dto.Username,
		Email:        dto.Email,
		Password:     dto.Password,
		PasswordHash: phash,
	}, nil
}

func (m *AccountMapper) MapLogInAccountDTO(dto dto.LoginAccountDTO) model.Account {
	return model.Account{
		ID:           0,
		Name:         "",
		Username:     dto.Username,
		Email:        "",
		Password:     dto.Password,
		PasswordHash: "",
	}
}

func (m *AccountMapper) MapAccountWithTokenDTO(token string, a model.Account) dto.WithTokenDTO {
	return dto.WithTokenDTO{
		Token:    token,
		Name:     a.Name,
		Username: a.Username,
		Email:    a.Email,
	}
}

func (m *AccountMapper) MapAccountDTO(a model.Account) dto.GetAccountDTO {
	return dto.GetAccountDTO{
		Name:     a.Name,
		Username: a.Username,
		Email:    a.Email,
	}
}
