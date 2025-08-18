package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/Foxpunk/courseforge/internal/interfaces"
	"github.com/Foxpunk/courseforge/internal/models"
	"github.com/gin-gonic/gin"
)

// Middleware хранит зависимости из слоя application
type Middleware struct {
	authManager interfaces.AuthManager
}

// NewMiddleware принимает интерфейс AuthManager из application
func NewMiddleware(am interfaces.AuthManager) *Middleware {
	return &Middleware{authManager: am}
}

// AuthMiddleware проверяет JWT и кладёт доменную сущность User в контекст
func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		hdr := c.GetHeader("Authorization")
		if hdr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		token := strings.TrimPrefix(hdr, "Bearer ")
		if token == hdr {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		user, err := m.authManager.ValidateToken(context.Background(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		c.Set("user", user) // user — это *domain.User
		c.Next()
	}
}

// ниже — guards по ролям
func (m *Middleware) AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		u := m.getUser(c)
		if u == nil || !u.IsAdmin() {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			return
		}
		c.Next()
	}
}

func (m *Middleware) TeacherRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		u := m.getUser(c)
		if u == nil || !u.IsTeacher() {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Teacher access required"})
			return
		}
		c.Next()
	}
}

func (m *Middleware) StudentRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		u := m.getUser(c)
		if u == nil || !u.IsStudent() {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Student access required"})
			return
		}
		c.Next()
	}
}

func (m *Middleware) TeacherOrAdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		u := m.getUser(c)
		if u == nil || (!u.IsTeacher() && !u.IsAdmin()) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Teacher or admin access required"})
			return
		}
		c.Next()
	}
}

/*
// CORS если нужен фронт

	func (m *Middleware) CORS() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type")

			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(http.StatusNoContent)
				return
		}
			c.Next()
		}
	}
*/
func (m *Middleware) CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if origin == "http://localhost:3000" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// getUser достаёт *domain.User из gin.Context
func (m *Middleware) getUser(c *gin.Context) *models.User {
	raw, ok := c.Get("user")
	if !ok {
		return nil
	}
	u, ok := raw.(*models.User)
	if !ok {
		return nil
	}
	return u
}
