package managers

import (
	"context"
	"errors"

	"github.com/Foxpunk/courseforge/internal/interfaces"
	"github.com/Foxpunk/courseforge/internal/models"
)

// CourseworkManagerImpl реализует interfaces.CourseworkManager
type CourseworkManagerImpl struct {
	cwRepo interfaces.CourseworkRepository
	scRepo interfaces.StudentCourseworkRepository
}

// NewCourseworkManager создаёт новый CourseworkManager
func NewCourseworkManager(
	cwRepo interfaces.CourseworkRepository,
	scRepo interfaces.StudentCourseworkRepository,
) interfaces.CourseworkManager {
	return &CourseworkManagerImpl{
		cwRepo: cwRepo,
		scRepo: scRepo,
	}
}

// CreateCoursework создаёт новую курсовую работу
func (m *CourseworkManagerImpl) CreateCoursework(ctx context.Context, req interfaces.CreateCourseworkRequest) (*models.Coursework, error) {
	// базовая валидация
	if req.SubjectID == 0 || req.TeacherID == 0 {
		return nil, errors.New("subject and teacher IDs are required")
	}
	cw := &models.Coursework{
		Title:           req.Title,
		Description:     req.Description,
		Requirements:    req.Requirements,
		SubjectID:       req.SubjectID,
		TeacherID:       req.TeacherID,
		MaxStudents:     req.MaxStudents,
		DifficultyLevel: req.DifficultyLevel,
		IsAvailable:     true,
	}
	if err := m.cwRepo.Create(ctx, cw); err != nil {
		return nil, err
	}
	return cw, nil
}

// GetCoursework возвращает курсовую работу по ID
func (m *CourseworkManagerImpl) GetCoursework(ctx context.Context, cwID uint) (*models.Coursework, error) {
	return m.cwRepo.GetByID(ctx, cwID)
}

// UpdateCoursework обновляет курсовую работу
func (m *CourseworkManagerImpl) UpdateCoursework(ctx context.Context, cwID uint, req interfaces.UpdateCourseworkRequest) (*models.Coursework, error) {
	cw, err := m.cwRepo.GetByID(ctx, cwID)
	if err != nil {
		return nil, err
	}
	if req.Title != nil {
		cw.Title = *req.Title
	}
	if req.Description != nil {
		cw.Description = *req.Description
	}
	if req.Requirements != nil {
		cw.Requirements = *req.Requirements
	}
	if req.MaxStudents != nil {
		cw.MaxStudents = *req.MaxStudents
	}
	if req.DifficultyLevel != nil {
		cw.DifficultyLevel = *req.DifficultyLevel
	}
	if req.IsAvailable != nil {
		cw.IsAvailable = *req.IsAvailable
	}
	if err := m.cwRepo.Update(ctx, cw); err != nil {
		return nil, err
	}
	return cw, nil
}

// DeleteCoursework удаляет курсовую работу
func (m *CourseworkManagerImpl) DeleteCoursework(ctx context.Context, cwID uint) error {
	return m.cwRepo.Delete(ctx, cwID)
}

// ListCourseworks возвращает список и общее число
func (m *CourseworkManagerImpl) ListCourseworks(ctx context.Context, req interfaces.ListCourseworksRequest) ([]models.Coursework, int, error) {
	list, err := m.cwRepo.List(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, 0, err
	}
	// для MVP считаем total как длину листа
	total := len(list)
	return list, total, nil
}

// GetCourseworksBySubject возвращает по предмету
func (m *CourseworkManagerImpl) GetCourseworksBySubject(ctx context.Context, subjID uint) ([]models.Coursework, error) {
	return m.cwRepo.GetBySubject(ctx, subjID)
}

// GetCourseworksByTeacher возвращает по преподавателю
func (m *CourseworkManagerImpl) GetCourseworksByTeacher(ctx context.Context, teacherID uint) ([]models.Coursework, error) {
	return m.cwRepo.GetByTeacher(ctx, teacherID)
}

// GetAvailableCourseworks возвращает все доступные
func (m *CourseworkManagerImpl) GetAvailableCourseworks(ctx context.Context, studentID uint) ([]models.Coursework, error) {
	list, err := m.cwRepo.List(ctx, 0, 0)
	if err != nil {
		return nil, err
	}
	var available []models.Coursework
	for _, cw := range list {
		if cw.IsAvailable {
			available = append(available, cw)
		}
	}
	return available, nil
}

// SetCourseworkAvailability задаёт доступность
func (m *CourseworkManagerImpl) SetCourseworkAvailability(ctx context.Context, cwID uint, available bool) error {
	return m.cwRepo.SetAvailable(ctx, cwID, available)
}

// CanAssignStudentToCoursework проверяет, можно ли назначить студента
func (m *CourseworkManagerImpl) CanAssignStudentToCoursework(ctx context.Context, studentID, cwID uint) error {
	// проверка дубликата
	if _, err := m.scRepo.GetByStudent(ctx, studentID); err == nil {
		return errors.New("student already assigned")
	}
	// проверка мест
	cw, count, err := m.cwRepo.GetWithStudentCount(ctx, cwID)
	if err != nil {
		return err
	}
	if count >= cw.MaxStudents {
		return errors.New("no available slots")
	}
	return nil
}
