package drivers

import (
	"context"
	"errors"
	"fmt"

	"github.com/Foxpunk/courseforge/internal/interfaces"
	"github.com/Foxpunk/courseforge/internal/models"
	"gorm.io/gorm"
)

type departmentRepository struct {
	db *gorm.DB
}

// NewDepartmentRepository создаёт новый репозиторий кафедр
func NewDepartmentRepository(db *gorm.DB) interfaces.DepartmentRepository {
	return &departmentRepository{db: db}
}

// Create добавляет новую кафедру
func (r *departmentRepository) Create(ctx context.Context, department *models.Department) error {
	if department == nil {
		return errors.New("department cannot be nil")
	}
	if department.DepartmentCode == "" || department.DepartmentName == "" {
		return errors.New("department code and name are required")
	}

	result := r.db.WithContext(ctx).Create(department)
	if result.Error != nil {
		return fmt.Errorf("failed to create department: %w", result.Error)
	}
	return nil
}

// GetByID возвращает кафедру по ID
func (r *departmentRepository) GetByID(ctx context.Context, id uint) (*models.Department, error) {
	if id == 0 {
		return nil, errors.New("invalid department ID")
	}

	var department models.Department
	result := r.db.WithContext(ctx).Preload("TeacherProfiles").Preload("Subjects").Preload("StudentGroups").First(&department, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("department with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get department by ID: %w", result.Error)
	}
	return &department, nil
}

// GetByCode возвращает кафедру по коду
func (r *departmentRepository) GetByCode(ctx context.Context, code string) (*models.Department, error) {
	if code == "" {
		return nil, errors.New("department code is required")
	}

	var department models.Department
	result := r.db.WithContext(ctx).Preload("TeacherProfiles").Preload("Subjects").Preload("StudentGroups").
		Where("department_code = ?", code).
		First(&department)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("department with code %q not found", code)
		}
		return nil, fmt.Errorf("failed to get department by code: %w", result.Error)
	}
	return &department, nil
}

// Update обновляет данные кафедры
func (r *departmentRepository) Update(ctx context.Context, department *models.Department) error {
	if department == nil {
		return errors.New("department cannot be nil")
	}
	if department.ID == 0 {
		return errors.New("department ID is required")
	}

	result := r.db.WithContext(ctx).Save(department)
	if result.Error != nil {
		return fmt.Errorf("failed to update department: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("department with ID %d not found", department.ID)
	}
	return nil
}

// Delete удаляет кафедру по ID (мягкое удаление)
func (r *departmentRepository) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid department ID")
	}

	result := r.db.WithContext(ctx).Delete(&models.Department{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete department: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("department with ID %d not found", id)
	}
	return nil
}

// List возвращает все кафедры
func (r *departmentRepository) List(ctx context.Context) ([]models.Department, error) {
	var departments []models.Department
	result := r.db.WithContext(ctx).Preload("TeacherProfiles").Preload("Subjects").Preload("StudentGroups").Find(&departments)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to list departments: %w", result.Error)
	}
	return departments, nil
}
