package auth

import (
	"neatly/internal/model/auth"
	"neatly/internal/repository"
	"neatly/pkg/jwt"
)

type Service struct {
	repository repository.Authorisation
}

func NewService(repository repository.Authorisation) *Service {
	return &Service{repository: repository}
}

func (s *Service) CreateAccount(u *auth.Account) error {
	err := s.repository.CreateAccount(u)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GenerateJWT(u *auth.Account) (string, error) {
	err := s.repository.GetAccount(u)
	if err != nil {
		return "", err
	}
	token, err := jwt.GenerateAccessToken(*u)
	if err != nil {
		return "", err
	}

	return token, nil
}
