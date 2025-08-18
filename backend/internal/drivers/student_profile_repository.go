package drivers

import (
	"context"
	"errors"
	"fmt"

	"github.com/Foxpunk/courseforge/internal/interfaces"
	"github.com/Foxpunk/courseforge/internal/models"
	"gorm.io/gorm"
)

// studentProfileRepository — реализация interfaces.StudentProfileRepository на GORM
// Работает с таблицей student_profiles

type studentProfileRepository struct {
	db *gorm.DB
}

// NewStudentProfileRepository создаёт новый репозиторий профилей студентов
func NewStudentProfileRepository(db *gorm.DB) interfaces.StudentProfileRepository {
	return &studentProfileRepository{db: db}
}

// Create создаёт новый профиль студента
func (r *studentProfileRepository) Create(ctx context.Context, profile *models.StudentProfile) error {
	if profile == nil {
		return errors.New("profile cannot be nil")
	}
	if profile.UserID == 0 {
		return errors.New("user ID is required")
	}
	if profile.GroupID == 0 {
		return errors.New("group ID is required")
	}

	result := r.db.WithContext(ctx).Create(profile)
	if result.Error != nil {
		return fmt.Errorf("failed to create student profile: %w", result.Error)
	}
	return nil
}

// GetByUserID возвращает профиль студента по ID пользователя
func (r *studentProfileRepository) GetByUserID(ctx context.Context, userID uint) (*models.StudentProfile, error) {
	if userID == 0 {
		return nil, errors.New("user ID is required")
	}

	var profile models.StudentProfile
	result := r.db.WithContext(ctx).
		Preload("User").
		Preload("StudentGroup").
		Where("user_id = ?", userID).
		First(&profile)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("student profile for user %d not found", userID)
		}
		return nil, fmt.Errorf("failed to get student profile by user ID: %w", result.Error)
	}
	return &profile, nil
}

// GetByID возвращает профиль студента по ID
func (r *studentProfileRepository) GetByID(ctx context.Context, id uint) (*models.StudentProfile, error) {
	if id == 0 {
		return nil, errors.New("invalid profile ID")
	}

	var profile models.StudentProfile
	result := r.db.WithContext(ctx).
		Preload("User").
		Preload("StudentGroup").
		First(&profile, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("student profile with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get student profile by ID: %w", result.Error)
	}
	return &profile, nil
}

// Update сохраняет изменения в профиле студента
func (r *studentProfileRepository) Update(ctx context.Context, profile *models.StudentProfile) error {
	if profile == nil {
		return errors.New("profile cannot be nil")
	}
	if profile.ID == 0 {
		return errors.New("profile ID is required")
	}

	result := r.db.WithContext(ctx).Save(profile)
	if result.Error != nil {
		return fmt.Errorf("failed to update student profile: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("student profile with ID %d not found", profile.ID)
	}
	return nil
}

// Delete выполняет мягкое удаление профиля по ID
func (r *studentProfileRepository) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid profile ID")
	}

	result := r.db.WithContext(ctx).Delete(&models.StudentProfile{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete student profile: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("student profile with ID %d not found", id)
	}
	return nil
}

// GetByGroup возвращает все профили студентов указанной группы
func (r *studentProfileRepository) GetByGroup(ctx context.Context, groupID uint) ([]models.StudentProfile, error) {
	if groupID == 0 {
		return nil, errors.New("group ID is required")
	}

	var profiles []models.StudentProfile
	result := r.db.WithContext(ctx).
		Preload("User").
		Preload("StudentGroup").
		Where("group_id = ?", groupID).
		Find(&profiles)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get student profiles by group: %w", result.Error)
	}
	return profiles, nil
}
