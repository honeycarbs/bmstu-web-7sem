package auth

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"neatly/pkg/e"
)

type Account struct {
	ID           int    `json:"-" db:"id"`
	Name         string `json:"name" binding:"required"`
	Username     string `json:"username" binding:"required"`
	Email        string `json:"email" binding:"required"`
	Password     string `json:"-"`
	PasswordHash string `json:"password" binding:"required" db:"password_hash"`
}

func (u *Account) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return &e.PasswordDoesNotMatchErr{}
	}
	return nil
}

func GeneratePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password due to error %w", err)
	}
	return string(hash), nil
}
