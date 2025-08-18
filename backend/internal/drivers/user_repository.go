package drivers

import (
	"context"
	"errors"
	"fmt"

	"github.com/Foxpunk/courseforge/internal/interfaces"
	"github.com/Foxpunk/courseforge/internal/models"
	"gorm.io/gorm"
)

// userRepository - реализация интерфейса UserRepository с использованием GORM
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository создает новый экземпляр репозитория пользователей
func NewUserRepository(db *gorm.DB) interfaces.UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create создает нового пользователя
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return fmt.Errorf("failed to create user: %w", result.Error)
	}

	return nil
}

// GetByID получает пользователя по ID
func (r *userRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	if id == 0 {
		return nil, errors.New("invalid user ID")
	}

	var user models.User
	result := r.db.WithContext(ctx).
		Preload("TeacherSubjects").
		First(&user, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", result.Error)
	}

	return &user, nil
}

// GetByEmail получает пользователя по email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	var user models.User
	result := r.db.WithContext(ctx).
		Preload("TeacherSubjects").
		Where("email = ?", email).
		First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, fmt.Errorf("failed to get user by email: %w", result.Error)
	}

	return &user, nil
}

// Update обновляет пользователя
func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	if user.ID == 0 {
		return errors.New("user ID cannot be zero")
	}

	result := r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", user.ID)
	}

	return nil
}

// Delete удаляет пользователя (мягкое удаление)
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid user ID")
	}

	result := r.db.WithContext(ctx).Delete(&models.User{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", id)
	}

	return nil
}

// List получает список пользователей с пагинацией
func (r *userRepository) List(ctx context.Context, limit, offset int) ([]models.User, error) {
	if limit <= 0 {
		limit = 10 // значение по умолчанию
	}
	if offset < 0 {
		offset = 0
	}

	var users []models.User
	result := r.db.WithContext(ctx).
		Preload("TeacherSubjects").
		Limit(limit).
		Offset(offset).
		Find(&users)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get users list: %w", result.Error)
	}

	return users, nil
}

// GetByRole получает пользователей по роли
func (r *userRepository) GetByRole(ctx context.Context, role models.UserRole) ([]models.User, error) {
	if role == "" {
		return nil, errors.New("role cannot be empty")
	}

	var users []models.User
	result := r.db.WithContext(ctx).
		Preload("TeacherSubjects").
		Where("role = ?", role).
		Find(&users)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get users by role: %w", result.Error)
	}

	return users, nil
}

// UpdateRole обновляет роль пользователя
func (r *userRepository) UpdateRole(ctx context.Context, userID uint, role models.UserRole) error {
	if userID == 0 {
		return errors.New("invalid user ID")
	}

	if role == "" {
		return errors.New("role cannot be empty")
	}

	// Проверяем, что роль валидна
	validRoles := []models.UserRole{models.RoleAdmin, models.RoleTeacher, models.RoleStudent}
	isValidRole := false
	for _, validRole := range validRoles {
		if role == validRole {
			isValidRole = true
			break
		}
	}

	if !isValidRole {
		return fmt.Errorf("invalid role: %s", role)
	}

	result := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("role", role)

	if result.Error != nil {
		return fmt.Errorf("failed to update user role: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", userID)
	}

	return nil
}

// SetActive устанавливает статус активности пользователя
func (r *userRepository) SetActive(ctx context.Context, userID uint, active bool) error {
	if userID == 0 {
		return errors.New("invalid user ID")
	}

	result := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("is_active", active)

	if result.Error != nil {
		return fmt.Errorf("failed to set user active status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", userID)
	}

	return nil
}

// Дополнительные методы для удобства работы

// GetActiveUsers получает только активных пользователей
func (r *userRepository) GetActiveUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	result := r.db.WithContext(ctx).
		Preload("TeacherSubjects").
		Where("is_active = ?", true).
		Find(&users)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get active users: %w", result.Error)
	}

	return users, nil
}

// GetUserCount получает общее количество пользователей
func (r *userRepository) GetUserCount(ctx context.Context) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&models.User{}).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to get user count: %w", result.Error)
	}
	return count, nil
}

// GetUserCountByRole получает количество пользователей по роли
func (r *userRepository) GetUserCountByRole(ctx context.Context, role models.UserRole) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("role = ?", role).
		Count(&count)

	if result.Error != nil {
		return 0, fmt.Errorf("failed to get user count by role: %w", result.Error)
	}
	return count, nil
}
