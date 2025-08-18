package drivers

import (
	"context"
	"errors"
	"fmt"

	"github.com/Foxpunk/courseforge/internal/interfaces"
	"github.com/Foxpunk/courseforge/internal/models"
	"gorm.io/gorm"
)

type teacherProfileRepository struct {
	db *gorm.DB
}

// NewTeacherProfileRepository создает новый репозиторий профилей преподавателей
func NewTeacherProfileRepository(db *gorm.DB) interfaces.TeacherProfileRepository {
	return &teacherProfileRepository{db: db}
}

// Create создает новый профиль преподавателя
func (r *teacherProfileRepository) Create(ctx context.Context, profile *models.TeacherProfile) error {
	if profile == nil {
		return errors.New("teacher profile is nil")
	}
	if profile.UserID == 0 || profile.DepartmentID == 0 {
		return errors.New("user ID and department ID are required")
	}

	result := r.db.WithContext(ctx).Create(profile)
	if result.Error != nil {
		return fmt.Errorf("failed to create teacher profile: %w", result.Error)
	}
	return nil
}

// GetByUserID возвращает профиль по ID пользователя
func (r *teacherProfileRepository) GetByUserID(ctx context.Context, userID uint) (*models.TeacherProfile, error) {
	if userID == 0 {
		return nil, errors.New("user ID is required")
	}

	var profile models.TeacherProfile
	result := r.db.WithContext(ctx).
		Preload("User").
		Preload("Department").
		Where("user_id = ?", userID).
		First(&profile)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("teacher profile for user %d not found", userID)
		}
		return nil, fmt.Errorf("failed to get teacher profile by user ID: %w", result.Error)
	}
	return &profile, nil
}

// GetByID возвращает профиль по его ID
func (r *teacherProfileRepository) GetByID(ctx context.Context, id uint) (*models.TeacherProfile, error) {
	if id == 0 {
		return nil, errors.New("profile ID is required")
	}

	var profile models.TeacherProfile
	result := r.db.WithContext(ctx).
		Preload("User").
		Preload("Department").
		First(&profile, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("teacher profile with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get teacher profile by ID: %w", result.Error)
	}
	return &profile, nil
}

// Update обновляет существующий профиль
func (r *teacherProfileRepository) Update(ctx context.Context, profile *models.TeacherProfile) error {
	if profile == nil || profile.ID == 0 {
		return errors.New("invalid profile data")
	}

	result := r.db.WithContext(ctx).Save(profile)
	if result.Error != nil {
		return fmt.Errorf("failed to update teacher profile: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("teacher profile with ID %d not found", profile.ID)
	}
	return nil
}

// Delete удаляет профиль преподавателя по ID (мягкое удаление)
func (r *teacherProfileRepository) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid ID")
	}

	result := r.db.WithContext(ctx).Delete(&models.TeacherProfile{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete teacher profile: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("teacher profile with ID %d not found", id)
	}
	return nil
}

// GetByDepartment возвращает все профили из одной кафедры
func (r *teacherProfileRepository) GetByDepartment(ctx context.Context, departmentID uint) ([]models.TeacherProfile, error) {
	if departmentID == 0 {
		return nil, errors.New("department ID is required")
	}

	var profiles []models.TeacherProfile
	result := r.db.WithContext(ctx).
		Preload("User").
		Preload("Department").
		Where("department_id = ?", departmentID).
		Find(&profiles)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get teacher profiles by department: %w", result.Error)
	}
	return profiles, nil
}
