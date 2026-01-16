package models

import (
	"errors"
)

var (
	ErrUserNotFound       = errors.New("Пользователь не найден")
	ErrEmailAlreadyExists = errors.New("Электронная почта уже существует")
)
