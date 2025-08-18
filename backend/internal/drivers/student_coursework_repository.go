package drivers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Foxpunk/courseforge/internal/interfaces"
	"github.com/Foxpunk/courseforge/internal/models"
	"gorm.io/gorm"
)

type studentCourseworkRepository struct {
	db *gorm.DB
}

// NewStudentCourseworkRepository создаёт новый репозиторий назначений студентов
func NewStudentCourseworkRepository(db *gorm.DB) interfaces.StudentCourseworkRepository {
	return &studentCourseworkRepository{db: db}
}

// Create создаёт новое назначение студента на курсовую работу
func (r *studentCourseworkRepository) Create(ctx context.Context, assignment *models.StudentCoursework) error {
	if assignment == nil {
		return errors.New("assignment cannot be nil")
	}
	if assignment.StudentID == 0 || assignment.CourseworkID == 0 {
		return errors.New("student ID and coursework ID are required")
	}

	result := r.db.WithContext(ctx).Create(assignment)
	if result.Error != nil {
		return fmt.Errorf("failed to create student coursework: %w", result.Error)
	}
	return nil
}

// GetByID возвращает назначение по его ID
func (r *studentCourseworkRepository) GetByID(ctx context.Context, id uint) (*models.StudentCoursework, error) {
	if id == 0 {
		return nil, errors.New("invalid ID")
	}

	var sc models.StudentCoursework
	result := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Coursework").
		First(&sc, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("student coursework with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get by ID: %w", result.Error)
	}
	return &sc, nil
}

// GetByStudent возвращает единственное назначение по студенту
func (r *studentCourseworkRepository) GetByStudent(ctx context.Context, studentID uint) (*models.StudentCoursework, error) {
	if studentID == 0 {
		return nil, errors.New("invalid student ID")
	}

	var sc models.StudentCoursework
	result := r.db.WithContext(ctx).
		Preload("Coursework").
		Where("student_id = ?", studentID).
		First(&sc)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no coursework found for student %d", studentID)
		}
		return nil, fmt.Errorf("failed to get by student: %w", result.Error)
	}
	return &sc, nil
}

// GetByCoursework возвращает список назначений по курсовой работе
func (r *studentCourseworkRepository) GetByCoursework(ctx context.Context, courseworkID uint) ([]models.StudentCoursework, error) {
	if courseworkID == 0 {
		return nil, errors.New("invalid coursework ID")
	}

	var list []models.StudentCoursework
	result := r.db.WithContext(ctx).
		Preload("Student").
		Where("coursework_id = ?", courseworkID).
		Find(&list)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get by coursework: %w", result.Error)
	}
	return list, nil
}

// Update сохраняет изменения в назначении
func (r *studentCourseworkRepository) Update(ctx context.Context, assignment *models.StudentCoursework) error {
	if assignment == nil || assignment.ID == 0 {
		return errors.New("invalid assignment")
	}

	result := r.db.WithContext(ctx).Save(assignment)
	if result.Error != nil {
		return fmt.Errorf("failed to update student coursework: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("student coursework with ID %d not found", assignment.ID)
	}
	return nil
}

// Delete выполняет мягкое удаление назначение по ID
func (r *studentCourseworkRepository) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid ID")
	}

	result := r.db.WithContext(ctx).Delete(&models.StudentCoursework{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete student coursework: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("student coursework with ID %d not found", id)
	}
	return nil
}

// UpdateStatus обновляет статус выполнения курсовой работы
func (r *studentCourseworkRepository) UpdateStatus(ctx context.Context, id uint, status models.CourseworkStatus) error {
	if id == 0 {
		return errors.New("invalid ID")
	}
	// Валидация статуса
	switch status {
	case models.StatusAssigned, models.StatusInProgress, models.StatusSubmitted, models.StatusReviewed, models.StatusCompleted, models.StatusFailed:
	default:
		return fmt.Errorf("invalid status: %s", status)
	}

	result := r.db.WithContext(ctx).
		Model(&models.StudentCoursework{}).
		Where("id = ?", id).
		Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("failed to update status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("student coursework with ID %d not found", id)
	}
	return nil
}

// SetGrade устанавливает оценку и обратную связь
func (r *studentCourseworkRepository) SetGrade(ctx context.Context, id uint, grade int, feedback string) error {
	if id == 0 {
		return errors.New("invalid ID")
	}
	if grade < 2 || grade > 5 {
		return fmt.Errorf("invalid grade: %d", grade)
	}

	result := r.db.WithContext(ctx).
		Model(&models.StudentCoursework{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"grade":    grade,
			"feedback": feedback,
		})
	if result.Error != nil {
		return fmt.Errorf("failed to set grade: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("student coursework with ID %d not found", id)
	}
	return nil
}

// SetSubmitted отмечает назначение как отправленное
func (r *studentCourseworkRepository) SetSubmitted(ctx context.Context, id uint, submittedAt time.Time) error {
	if id == 0 {
		return errors.New("invalid ID")
	}

	result := r.db.WithContext(ctx).
		Model(&models.StudentCoursework{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       models.StatusSubmitted,
			"submitted_at": submittedAt,
		})
	if result.Error != nil {
		return fmt.Errorf("failed to set submitted: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("student coursework with ID %d not found", id)
	}
	return nil
}

// SetCompleted отмечает назначение как завершённое
func (r *studentCourseworkRepository) SetCompleted(ctx context.Context, id uint, completedAt time.Time) error {
	if id == 0 {
		return errors.New("invalid ID")
	}

	result := r.db.WithContext(ctx).
		Model(&models.StudentCoursework{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       models.StatusCompleted,
			"completed_at": completedAt,
		})
	if result.Error != nil {
		return fmt.Errorf("failed to set completed: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("student coursework with ID %d not found", id)
	}
	return nil
}

// GetByTeacher возвращает все назначения для курсовых, которыми руководит преподаватель
func (r *studentCourseworkRepository) GetByTeacher(ctx context.Context, teacherID uint) ([]models.StudentCoursework, error) {
	if teacherID == 0 {
		return nil, errors.New("invalid teacher ID")
	}

	var list []models.StudentCoursework
	result := r.db.WithContext(ctx).
		Joins("JOIN courseworks ON courseworks.id = student_courseworks.coursework_id").
		Where("courseworks.teacher_id = ?", teacherID).
		Preload("Student").
		Preload("Coursework").
		Find(&list)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get by teacher: %w", result.Error)
	}
	return list, nil
}
