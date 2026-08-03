package models

import (
	"errors"
	"time"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidInput = errors.New("invalid input data")
)

type Author struct {
	Id         uint64    `json:"id" db:"id"`
	Username   string    `json:"username" db:"username"`
	Email      string    `json:"email" db:"email"`
	Password   string    `json:"-" db:"password"`
	Created_at time.Time `json:"created_at" db:"created_at"`
}

// Создание автора
type CreateAuthorDTO struct {
	Username string `json:"username" db:"username`
	Email    string `json:"email" db:"email"`
	Password string `json:password" db"password"`
}

// Обновление имени автора
type UpdateUsernameDTO struct {
	Id       uint64 `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
}
