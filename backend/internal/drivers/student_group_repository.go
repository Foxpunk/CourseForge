package drivers

import (
	"context"
	"errors"
	"fmt"

	"github.com/Foxpunk/courseforge/internal/interfaces"
	"github.com/Foxpunk/courseforge/internal/models"
	"gorm.io/gorm"
)

// studentGroupRepository — реализация interfaces.StudentGroupRepository на GORM
type studentGroupRepository struct {
	db *gorm.DB
}

// NewStudentGroupRepository создаёт новый репозиторий групп студентов
func NewStudentGroupRepository(db *gorm.DB) interfaces.StudentGroupRepository {
	return &studentGroupRepository{db: db}
}

// Create создаёт новую группу студентов
func (r *studentGroupRepository) Create(ctx context.Context, group *models.StudentGroup) error {
	if group == nil {
		return errors.New("group cannot be nil")
	}
	if group.GroupCode == "" {
		return errors.New("group code cannot be empty")
	}

	result := r.db.WithContext(ctx).Create(group)
	if result.Error != nil {
		return fmt.Errorf("failed to create student group: %w", result.Error)
	}
	return nil
}

// GetByID возвращает группу по ID, вместе с данными кафедры
func (r *studentGroupRepository) GetByID(ctx context.Context, id uint) (*models.StudentGroup, error) {
	if id == 0 {
		return nil, errors.New("invalid group ID")
	}

	var group models.StudentGroup
	result := r.db.WithContext(ctx).
		Preload("Department").
		First(&group, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("student group with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get student group by ID: %w", result.Error)
	}
	return &group, nil
}

// GetByCode возвращает группу по её коду
func (r *studentGroupRepository) GetByCode(ctx context.Context, code string) (*models.StudentGroup, error) {
	if code == "" {
		return nil, errors.New("group code cannot be empty")
	}

	var group models.StudentGroup
	result := r.db.WithContext(ctx).
		Preload("Department").
		Where("group_code = ?", code).
		First(&group)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("student group with code %q not found", code)
		}
		return nil, fmt.Errorf("failed to get student group by code: %w", result.Error)
	}
	return &group, nil
}

// Update сохраняет изменения в группе
func (r *studentGroupRepository) Update(ctx context.Context, group *models.StudentGroup) error {
	if group == nil {
		return errors.New("group cannot be nil")
	}
	if group.ID == 0 {
		return errors.New("group ID cannot be zero")
	}

	result := r.db.WithContext(ctx).Save(group)
	if result.Error != nil {
		return fmt.Errorf("failed to update student group: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("student group with ID %d not found", group.ID)
	}
	return nil
}

// Delete выполняет мягкое удаление группы по ID
func (r *studentGroupRepository) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid group ID")
	}

	result := r.db.WithContext(ctx).Delete(&models.StudentGroup{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete student group: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("student group with ID %d not found", id)
	}
	return nil
}

// List возвращает все группы (без пагинации)
func (r *studentGroupRepository) List(ctx context.Context) ([]models.StudentGroup, error) {
	var groups []models.StudentGroup
	result := r.db.WithContext(ctx).
		Preload("Department").
		Find(&groups)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to list student groups: %w", result.Error)
	}
	return groups, nil
}

// GetByDepartment возвращает все группы указанной кафедры
func (r *studentGroupRepository) GetByDepartment(ctx context.Context, departmentID uint) ([]models.StudentGroup, error) {
	if departmentID == 0 {
		return nil, errors.New("invalid department ID")
	}

	var groups []models.StudentGroup
	result := r.db.WithContext(ctx).
		Preload("Department").
		Where("department_id = ?", departmentID).
		Find(&groups)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get groups by department: %w", result.Error)
	}
	return groups, nil
}
