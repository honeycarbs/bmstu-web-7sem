package model

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID           int    `json:"-" db:"id"`
	Name         string `json:"name" binding:"required"`
	Username     string `json:"username" binding:"required"`
	Email        string `json:"email" binding:"required"`
	Password     string `json:"-"`
	PasswordHash string `json:"password" binding:"required" db:"password_hash"`
}

func (a *Account) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(password))
	if err != nil {
		return errors.New("password does not match")
	}
	return nil
}

func GeneratePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
