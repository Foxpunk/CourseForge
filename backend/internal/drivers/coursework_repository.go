package drivers

import (
	"context"
	"errors"
	"fmt"

	"github.com/Foxpunk/courseforge/internal/interfaces"
	"github.com/Foxpunk/courseforge/internal/models"
	"gorm.io/gorm"
)

type courseworkRepository struct {
	db *gorm.DB
}

// NewCourseworkRepository создаёт новый репозиторий курсовых работ
func NewCourseworkRepository(db *gorm.DB) interfaces.CourseworkRepository {
	return &courseworkRepository{db: db}
}

// Create создаёт новую курсовую работу (тему)
func (r *courseworkRepository) Create(ctx context.Context, cw *models.Coursework) error {
	if cw == nil {
		return errors.New("coursework cannot be nil")
	}
	if cw.Title == "" || cw.Description == "" {
		return errors.New("title and description are required")
	}
	if cw.SubjectID == 0 || cw.TeacherID == 0 {
		return errors.New("subject ID and teacher ID are required")
	}
	if cw.MaxStudents < 1 {
		return errors.New("max_students must be at least 1")
	}
	if cw.DifficultyLevel != models.Easy && cw.DifficultyLevel != models.Medium && cw.DifficultyLevel != models.Hard {
		return fmt.Errorf("invalid difficulty level: %s", cw.DifficultyLevel)
	}

	result := r.db.WithContext(ctx).Create(cw)
	if result.Error != nil {
		return fmt.Errorf("failed to create coursework: %w", result.Error)
	}
	return nil
}

// GetByID возвращает курсовую по ID, подгружая предмет и преподавателя
func (r *courseworkRepository) GetByID(ctx context.Context, id uint) (*models.Coursework, error) {
	if id == 0 {
		return nil, errors.New("invalid coursework ID")
	}

	var cw models.Coursework
	result := r.db.WithContext(ctx).
		Preload("Subject").
		Preload("Teacher").
		First(&cw, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("coursework with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get coursework by ID: %w", result.Error)
	}
	return &cw, nil
}

// Update сохраняет изменения в курсовой работе
func (r *courseworkRepository) Update(ctx context.Context, cw *models.Coursework) error {
	if cw == nil {
		return errors.New("coursework cannot be nil")
	}
	if cw.ID == 0 {
		return errors.New("coursework ID cannot be zero")
	}

	result := r.db.WithContext(ctx).Save(cw)
	if result.Error != nil {
		return fmt.Errorf("failed to update coursework: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("coursework with ID %d not found", cw.ID)
	}
	return nil
}

// Delete выполняет мягкое удаление курсовой работы
func (r *courseworkRepository) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid coursework ID")
	}

	result := r.db.WithContext(ctx).Delete(&models.Coursework{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete coursework: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("coursework with ID %d not found", id)
	}
	return nil
}

// List возвращает список курсовых работ с пагинацией
func (r *courseworkRepository) List(ctx context.Context, limit, offset int) ([]models.Coursework, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	var list []models.Coursework
	result := r.db.WithContext(ctx).
		Preload("Subject").
		Preload("Teacher").
		Limit(limit).
		Offset(offset).
		Find(&list)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to list courseworks: %w", result.Error)
	}
	return list, nil
}

// GetBySubject возвращает курсовые работы по предмету
func (r *courseworkRepository) GetBySubject(ctx context.Context, subjectID uint) ([]models.Coursework, error) {
	if subjectID == 0 {
		return nil, errors.New("invalid subject ID")
	}

	var list []models.Coursework
	result := r.db.WithContext(ctx).
		Where("subject_id = ?", subjectID).
		Preload("Teacher").
		Find(&list)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get courseworks by subject: %w", result.Error)
	}
	return list, nil
}

// GetByTeacher возвращает курсовые работы по преподавателю
func (r *courseworkRepository) GetByTeacher(ctx context.Context, teacherID uint) ([]models.Coursework, error) {
	if teacherID == 0 {
		return nil, errors.New("invalid teacher ID")
	}

	var list []models.Coursework
	result := r.db.WithContext(ctx).
		Where("teacher_id = ?", teacherID).
		Preload("Subject").
		Find(&list)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get courseworks by teacher: %w", result.Error)
	}
	return list, nil
}

// GetAvailable возвращает доступные курсовые работы по предмету
func (r *courseworkRepository) GetAvailable(ctx context.Context, subjectID uint) ([]models.Coursework, error) {
	if subjectID == 0 {
		return nil, errors.New("invalid subject ID")
	}

	var list []models.Coursework
	result := r.db.WithContext(ctx).
		Where("subject_id = ? AND is_available = ?", subjectID, true).
		Preload("Subject").
		Preload("Teacher").
		Find(&list)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get available courseworks: %w", result.Error)
	}
	return list, nil
}

// SetAvailable устанавливает флаг доступности курсовой работы
func (r *courseworkRepository) SetAvailable(ctx context.Context, courseworkID uint, available bool) error {
	if courseworkID == 0 {
		return errors.New("invalid coursework ID")
	}

	result := r.db.WithContext(ctx).
		Model(&models.Coursework{}).
		Where("id = ?", courseworkID).
		Update("is_available", available)

	if result.Error != nil {
		return fmt.Errorf("failed to set availability: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("coursework with ID %d not found", courseworkID)
	}
	return nil
}

// GetWithStudentCount возвращает курсовую работу и число записанных студентов
func (r *courseworkRepository) GetWithStudentCount(ctx context.Context, courseworkID uint) (*models.Coursework, int, error) {
	if courseworkID == 0 {
		return nil, 0, errors.New("invalid coursework ID")
	}

	cw, err := r.GetByID(ctx, courseworkID)
	if err != nil {
		return nil, 0, err
	}

	var count int64
	// Предполагаем, что есть модель StudentCoursework с полем CourseworkID
	result := r.db.WithContext(ctx).
		Model(&models.StudentCoursework{}).
		Where("coursework_id = ?", courseworkID).
		Count(&count)
	if result.Error != nil {
		return cw, 0, fmt.Errorf("failed to count students for coursework %d: %w", courseworkID, result.Error)
	}

	return cw, int(count), nil
}
