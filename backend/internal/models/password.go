package models

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword хеширует пароль с помощью bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword проверяет пароль с хешем
func CheckPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
