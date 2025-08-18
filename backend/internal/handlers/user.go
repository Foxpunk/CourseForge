package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Foxpunk/courseforge/internal/interfaces"
	"github.com/Foxpunk/courseforge/internal/models"
)

const (
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// UserHandler предоставляет CRUD для пользователей (админ)
type UserHandler struct {
	userManager interfaces.UserManager
}

// NewUserHandler создаёт новый UserHandler
func NewUserHandler(um interfaces.UserManager) *UserHandler {
	return &UserHandler{userManager: um}
}

// CreateUser - создание пользователя (для админа)
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req interfaces.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userManager.CreateUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, interfaces.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	})
}

// ListUsers - список пользователей с пагинацией и фильтрами
func (h *UserHandler) ListUsers(c *gin.Context) {
	q := c.Request.URL.Query()
	var req interfaces.ListUsersRequest

	// Парсим limit
	limit := DefaultPageSize
	if v := q.Get("limit"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			limit = i
		}
	}
	if limit > MaxPageSize {
		limit = MaxPageSize
	}
	// Парсим offset
	offset := 0
	if v := q.Get("offset"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			offset = i
		}
	}

	req.Limit = limit
	req.Offset = offset

	// Фильтр role
	if v := q.Get("role"); v != "" {
		r := models.UserRole(v)
		req.Role = &r
	}
	// Фильтр active
	if v := q.Get("active"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			req.Active = &b
		}
	}

	users, total, err := h.userManager.ListUsers(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Формируем ответ
	resp := make([]interfaces.UserResponse, len(users))
	for i, u := range users {
		resp[i] = interfaces.UserResponse{
			ID:        u.ID,
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Role:      u.Role,
			IsActive:  u.IsActive,
			CreatedAt: u.CreatedAt.Format(time.RFC3339),
		}
	}

	c.JSON(http.StatusOK, interfaces.UserListResponse{Users: resp, Total: int64(total)})
}

// GetUser - информация о конкретном пользователе (по ID)
func (h *UserHandler) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, err := h.userManager.GetUser(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, interfaces.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	})
}

// UpdateUser - редактирование данных пользователя
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req interfaces.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userManager.UpdateUser(c.Request.Context(), uint(id), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, interfaces.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	})
}

// DeleteUser - удаление пользователя
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	err = h.userManager.DeleteUser(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
