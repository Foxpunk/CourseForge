package drivers

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Foxpunk/courseforge/internal/interfaces"
	"github.com/Foxpunk/courseforge/internal/models"
	"gorm.io/gorm"
)

type teacherSubjectRepository struct {
	db *gorm.DB
}

func NewTeacherSubjectRepository(db *gorm.DB) interfaces.TeacherSubjectRepository {
	return &teacherSubjectRepository{db: db}
}

// Create создаёт новое назначение преподавателя на дисциплину
func (r *teacherSubjectRepository) Create(ctx context.Context, assignment *models.TeacherSubject) error {
	if assignment == nil {
		return errors.New("assignment cannot be nil")
	}
	if assignment.UserID == 0 || assignment.SubjectID == 0 {
		return errors.New("user ID and subject ID are required")
	}

	log.Printf("Repository Create: UserID=%d, SubjectID=%d", assignment.UserID, assignment.SubjectID)

	result := r.db.WithContext(ctx).Create(assignment)
	if result.Error != nil {
		log.Printf("GORM Create error: %v", result.Error)
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return fmt.Errorf("assignment already exists for teacher %d, subject %d",
				assignment.UserID, assignment.SubjectID)
		}
		return fmt.Errorf("failed to create teacher assignment: %w", result.Error)
	}

	log.Printf("Repository Create: success")
	return nil
}

// GetByID возвращает назначение по его ID (не работает с composite primary key)
func (r *teacherSubjectRepository) GetByID(ctx context.Context, id uint) (*models.TeacherSubject, error) {
	return nil, errors.New("GetByID not supported for composite primary key table")
}

// Delete выполняет удаление назначения (не работает с composite primary key через ID)
func (r *teacherSubjectRepository) Delete(ctx context.Context, id uint) error {
	return errors.New("Delete by ID not supported for composite primary key table")
}

// DeleteByTeacherAndSubject удаляет назначение по teacherID и subjectID
func (r *teacherSubjectRepository) DeleteByTeacherAndSubject(ctx context.Context, teacherID, subjectID uint) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND subject_id = ?", teacherID, subjectID).
		Delete(&models.TeacherSubject{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete assignment: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("assignment not found for teacher %d, subject %d", teacherID, subjectID)
	}
	return nil
}

// GetByTeacher возвращает назначения по преподавателю (игнорируем year)
func (r *teacherSubjectRepository) GetByTeacher(ctx context.Context, teacherID uint, year string) ([]models.TeacherSubject, error) {
	if teacherID == 0 {
		return nil, errors.New("teacher ID is required")
	}

	var list []models.TeacherSubject
	result := r.db.WithContext(ctx).
		Preload("Subject").
		Where("user_id = ?", teacherID).
		Find(&list)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get assignments by teacher: %w", result.Error)
	}
	return list, nil
}

// GetBySubject возвращает назначения по дисциплине (игнорируем year)
func (r *teacherSubjectRepository) GetBySubject(ctx context.Context, subjectID uint, year string) ([]models.TeacherSubject, error) {
	if subjectID == 0 {
		return nil, errors.New("subject ID is required")
	}

	var list []models.TeacherSubject
	result := r.db.WithContext(ctx).
		Preload("Teacher").
		Where("subject_id = ?", subjectID).
		Find(&list)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get assignments by subject: %w", result.Error)
	}
	return list, nil
}

// GetLeadTeacher не поддерживается простой таблицей
func (r *teacherSubjectRepository) GetLeadTeacher(ctx context.Context, subjectID uint, year string) (*models.TeacherSubject, error) {
	return nil, errors.New("lead teacher not supported with simple teacher_subjects table")
}

// SetLead не поддерживается простой таблицей
func (r *teacherSubjectRepository) SetLead(ctx context.Context, assignmentID uint, isLead bool) error {
	return errors.New("lead teacher not supported with simple teacher_subjects table")
}
