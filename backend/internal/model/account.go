package model

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID           int    `json:"-" db:"id" gorm:"column:id" dbq:"id"`
	Name         string `json:"name" binding:"required" gorm:"column:name" dbq:"name"`
	Username     string `json:"username" binding:"required" gorm:"column:username" dbq:"username"`
	Email        string `json:"email" binding:"required" gorm:"column:email" dbq:"email"`
	Password     string `json:"-"`
	PasswordHash string `json:"password" binding:"required" db:"password_hash" gorm:"column:password_hash" dbq:"password_hash"`
}

// Recommended by dbq
func (a *Account) ScanFast() []interface{} {
	return []interface{}{&a.ID, &a.Name, &a.Username, &a.Email, &a.PasswordHash}
}

// Required by gorm
func (Account) TableName() string {
	return "users"
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
