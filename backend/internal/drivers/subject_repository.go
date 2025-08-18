package drivers

import (
	"context"
	"errors"
	"fmt"

	"github.com/Foxpunk/courseforge/internal/interfaces"
	"github.com/Foxpunk/courseforge/internal/models"
	"gorm.io/gorm"
)

type subjectRepository struct {
	db *gorm.DB
}

// NewSubjectRepository создаёт новый репозиторий дисциплин
func NewSubjectRepository(db *gorm.DB) interfaces.SubjectRepository {
	return &subjectRepository{db: db}
}

// Create создаёт новую дисциплину
func (r *subjectRepository) Create(ctx context.Context, subject *models.Subject) error {
	if subject == nil {
		return errors.New("subject cannot be nil")
	}
	if subject.Name == "" || subject.Code == "" {
		return errors.New("name and code are required")
	}
	if subject.Semester < 1 {
		return errors.New("semester must be >= 1")
	}

	result := r.db.WithContext(ctx).Create(subject)
	if result.Error != nil {
		return fmt.Errorf("failed to create subject: %w", result.Error)
	}
	return nil
}

// GetByID возвращает дисциплину по ID, подгружая преподавателей и курсовые работы
func (r *subjectRepository) GetByID(ctx context.Context, id uint) (*models.Subject, error) {
	if id == 0 {
		return nil, errors.New("invalid subject ID")
	}

	var subj models.Subject
	result := r.db.WithContext(ctx).
		Preload("Teachers").
		Preload("Courseworks").
		First(&subj, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("subject with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get subject by ID: %w", result.Error)
	}
	return &subj, nil
}

// Update сохраняет изменения в дисциплине
func (r *subjectRepository) Update(ctx context.Context, subject *models.Subject) error {
	if subject == nil {
		return errors.New("subject cannot be nil")
	}
	if subject.ID == 0 {
		return errors.New("subject ID cannot be zero")
	}

	result := r.db.WithContext(ctx).Save(subject)
	if result.Error != nil {
		return fmt.Errorf("failed to update subject: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("subject with ID %d not found", subject.ID)
	}
	return nil
}

// Delete выполняет мягкое удаление дисциплины
func (r *subjectRepository) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid subject ID")
	}

	result := r.db.WithContext(ctx).Delete(&models.Subject{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete subject: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("subject with ID %d not found", id)
	}
	return nil
}

// List возвращает все дисциплины (без фильтрации)
func (r *subjectRepository) List(ctx context.Context) ([]models.Subject, error) {
	var subjects []models.Subject
	result := r.db.WithContext(ctx).
		Preload("Teachers").
		Preload("Courseworks").
		Find(&subjects)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to list subjects: %w", result.Error)
	}
	return subjects, nil
}

// GetByDepartment возвращает дисциплины заданной кафедры
func (r *subjectRepository) GetByDepartment(ctx context.Context, departmentID uint) ([]models.Subject, error) {
	if departmentID == 0 {
		return nil, errors.New("invalid department ID")
	}

	var subjects []models.Subject
	result := r.db.WithContext(ctx).
		Where("department_id = ?", departmentID).
		Preload("Teachers").
		Preload("Courseworks").
		Find(&subjects)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get subjects by department: %w", result.Error)
	}
	return subjects, nil
}

// GetBySemester возвращает дисциплины по семестру
func (r *subjectRepository) GetBySemester(ctx context.Context, semester int) ([]models.Subject, error) {
	if semester < 1 {
		return nil, errors.New("invalid semester")
	}

	var subjects []models.Subject
	result := r.db.WithContext(ctx).
		Where("semester = ?", semester).
		Preload("Teachers").
		Preload("Courseworks").
		Find(&subjects)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get subjects by semester: %w", result.Error)
	}
	return subjects, nil
}

// SetActive устанавливает флаг активности дисциплины
func (r *subjectRepository) SetActive(ctx context.Context, subjectID uint, active bool) error {
	if subjectID == 0 {
		return errors.New("invalid subject ID")
	}

	result := r.db.WithContext(ctx).
		Model(&models.Subject{}).
		Where("id = ?", subjectID).
		Update("is_active", active)

	if result.Error != nil {
		return fmt.Errorf("failed to set active status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("subject with ID %d not found", subjectID)
	}
	return nil
}
