// handlers/auth.go
package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Foxpunk/courseforge/internal/interfaces"
	"github.com/Foxpunk/courseforge/internal/models"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authManager interfaces.AuthManager
	userManager interfaces.UserManager
	jwtSecret   string
}

func NewAuthHandler(authManager interfaces.AuthManager, userManager interfaces.UserManager, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		authManager: authManager,
		userManager: userManager,
		jwtSecret:   jwtSecret,
	}
}

// Login - вход в систему
func (h *AuthHandler) Login(c *gin.Context) {
	log.Println("=== LOGIN HANDLER CALLED ===")
	var req interfaces.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Login bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Login attempt for email: %s", req.Email)
	ctx := context.Background()
	token, user, err := h.authManager.Login(ctx, req.Email, req.Password)
	if err != nil {
		log.Printf("Login failed for %s: %v", req.Email, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные учетные данные"})
		return
	}

	// Проверяем активность пользователя
	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "Аккаунт деактивирован"})
		return
	}
	log.Printf("Login successful for user: %s", user.Email)
	c.JSON(http.StatusOK, interfaces.LoginResponse{
		Token:     token,
		User:      *user,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	})
}

// Register - регистрация пользователя
func (h *AuthHandler) Register(c *gin.Context) {
	log.Println("=== REGISTER HANDLER CALLED ===")

	var req interfaces.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Register bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Registration attempt for email: %s, role: %s", req.Email, req.Role)

	ctx := context.Background()
	user, err := h.authManager.Register(ctx, req)
	if err != nil {
		log.Printf("Registration failed for %s: %v", req.Email, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Registration successful for user: %s", user.Email)

	// Генерируем токен для нового пользователя
	token, err := h.authManager.GenerateToken(user)
	if err != nil {
		log.Printf("Token generation failed for %s: %v", user.Email, err)
		// Возвращаем успех регистрации без токена
		c.JSON(http.StatusCreated, gin.H{
			"message": "Пользователь успешно зарегистрирован. Войдите в систему.",
			"user": interfaces.UserResponse{
				ID:        user.ID,
				Email:     user.Email,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Role:      user.Role,
				IsActive:  user.IsActive,
				CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			},
		})
		return
	}

	log.Printf("Token generated successfully for user: %s", user.Email)

	// Возвращаем LoginResponse для автоматического входа
	c.JSON(http.StatusCreated, interfaces.LoginResponse{
		Token:     token,
		User:      *user,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(), // Установите правильное время истечения
	})
}

// RefreshToken - обновление JWT токена
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req interfaces.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	newToken, err := h.authManager.RefreshToken(ctx, req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Недействительный токен"})
		return
	}

	c.JSON(http.StatusOK, interfaces.RefreshTokenResponse{
		Token:     newToken,
		ExpiresAt: 0, // Установите правильное время истечения токена
	})
}

// ChangePassword - изменение пароля
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req interfaces.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем пользователя из контекста (установленного middleware)
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	userID := user.(*models.User).ID
	ctx := context.Background()
	err := h.authManager.ChangePassword(ctx, userID, req.OldPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пароль успешно изменен"})
}

// ResetPassword - сброс пароля
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req interfaces.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	err := h.authManager.ResetPassword(ctx, req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Инструкции по сбросу пароля отправлены на email"})
}

// GetProfile - получение профиля текущего пользователя
func (h *AuthHandler) GetProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	userID := user.(*models.User).ID
	ctx := context.Background()
	userProfile, err := h.userManager.GetUser(ctx, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	c.JSON(http.StatusOK, userProfile)
}

// Logout - выход из системы (в простой реализации просто подтверждение)
func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Выход выполнен успешно"})
}

// ValidateToken - валидация токена (для внутреннего использования)
func (h *AuthHandler) ValidateToken(tokenString string) (*models.User, error) {
	ctx := context.Background()
	return h.authManager.ValidateToken(ctx, tokenString)
}
