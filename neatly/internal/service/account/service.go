package account

import (
	"neatly/internal/model/account"
	"neatly/internal/repository"
	"neatly/pkg/jwt"
	"neatly/pkg/logging"
)

type Service struct {
	repository repository.Authorisation
}

func NewService(repository repository.Authorisation) *Service {
	return &Service{repository: repository}
}

func (s *Service) CreateAccount(a *account.Account) error {
	err := a.Validate()
	if err != nil {
		return err
	}
	err = s.repository.CreateAccount(a)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GenerateJWT(a *account.Account) (string, error) {
	err := s.repository.GetAccount(a)
	logging.GetLogger().Info(a.ID, a.Name, a.Email)

	err = a.CheckPassword(a.Password)
	if err != nil {
		return "", err
	}

	token, err := jwt.GenerateAccessToken(*a)
	if err != nil {
		return "", err
	}

	return token, nil
}
