package account

import (
	"neatly/internal/model"
	"neatly/internal/repository"
	"neatly/pkg/jwt"
	"neatly/pkg/logging"
)

type Service struct {
	repository *repository.AccountRepositoryImpl
	logger     logging.Logger
}

func NewService(repository *repository.AccountRepositoryImpl, logger logging.Logger) *Service {
	return &Service{repository: repository, logger: logger}
}

func (s *Service) CreateAccount(a *model.Account) error {
	err := s.repository.CreateAccount(a)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GenerateJWT(a *model.Account) (string, error) {
	err := s.repository.AuthorizeAccount(a)
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
