package account

import (
	"neatly/internal/model/account"
	"neatly/internal/repository"
	"neatly/pkg/jwt"
	"neatly/pkg/logging"
)

type Service struct {
	repository repository.Account
}

func NewService(repository repository.Account) *Service {
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
	err := s.repository.AuthorizeAccount(a)
	logging.GetLogger().Info(a.ID, a.Name, a.Email)

	err = a.CheckPassword(a.Password)
	if err != nil {
		return "", err
	}

	token, err := jwt.GenerateAccessToken(a.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) GetOne(userID int) (account.Account, error) {
	a, err := s.repository.GetOne(userID)
	logging.GetLogger().Info(a.ID, a.Name, a.Email)

	return a, err
}
