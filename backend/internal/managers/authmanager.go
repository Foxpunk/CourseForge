package managers

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/Foxpunk/courseforge/internal/config"
	"github.com/Foxpunk/courseforge/internal/interfaces"
	"github.com/Foxpunk/courseforge/internal/models"
)

type AuthManager struct {
	userRepo interfaces.UserRepository
	jwtCfg   config.JWTConfig
}

func NewAuthManager(userRepo interfaces.UserRepository, jwtCfg config.JWTConfig) interfaces.AuthManager {
	return &AuthManager{userRepo: userRepo, jwtCfg: jwtCfg}
}

func (a *AuthManager) Register(ctx context.Context, req interfaces.RegisterRequest) (*models.User, error) {
	if _, err := a.userRepo.GetByEmail(ctx, req.Email); err == nil {
		return nil, errors.New("user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &models.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         req.Role,
		IsActive:     true,
	}
	if err := a.userRepo.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (a *AuthManager) Login(ctx context.Context, email, password string) (string, *models.User, error) {
	u, err := a.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}
	token, err := a.GenerateToken(u)
	return token, u, err
}

func (a *AuthManager) ValidateToken(ctx context.Context, tokenStr string) (*models.User, error) {
	claims := &jwt.RegisteredClaims{}
	tok, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(a.jwtCfg.SecretKey), nil
	})
	if err != nil || !tok.Valid {
		return nil, errors.New("invalid token")
	}
	id, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return nil, errors.New("invalid token subject")
	}
	return a.userRepo.GetByID(ctx, uint(id))
}

func (a *AuthManager) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
	u, err := a.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(oldPassword)); err != nil {
		return errors.New("wrong password")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return a.userRepo.Update(ctx, u)
}

func (a *AuthManager) ResetPassword(ctx context.Context, email string) error {
	u, err := a.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	tmp := fmt.Sprintf("tmp%d", time.Now().Unix())
	hash, err := bcrypt.GenerateFromPassword([]byte(tmp), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	if err := a.userRepo.Update(ctx, u); err != nil {
		return err
	}
	// TODO: отправить tmp по email
	return nil
}

func (a *AuthManager) RefreshToken(ctx context.Context, tokenStr string) (string, error) {
	u, err := a.ValidateToken(ctx, tokenStr)
	if err != nil {
		return "", err
	}
	return a.GenerateToken(u)
}

func (a *AuthManager) GenerateToken(u *models.User) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   fmt.Sprint(u.ID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.jwtCfg.AccessTokenDuration)),
		Issuer:    a.jwtCfg.Issuer,
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(a.jwtCfg.SecretKey))
}
