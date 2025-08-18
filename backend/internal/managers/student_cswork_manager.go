package managers

import (
	"context"
	"errors"
	"time"

	"github.com/Foxpunk/courseforge/internal/interfaces"
	"github.com/Foxpunk/courseforge/internal/models"
)

// StudentCourseworkManagerImpl реализует interfaces.StudentCourseworkManager
type StudentCourseworkManagerImpl struct {
	scRepo interfaces.StudentCourseworkRepository
	cwRepo interfaces.CourseworkRepository
}

// NewStudentCourseworkManager создаёт новый StudentCourseworkManager
func NewStudentCourseworkManager(
	scRepo interfaces.StudentCourseworkRepository,
	cwRepo interfaces.CourseworkRepository,
) interfaces.StudentCourseworkManager {
	return &StudentCourseworkManagerImpl{
		scRepo: scRepo,
		cwRepo: cwRepo,
	}
}

// AssignStudentToCoursework назначает студента на курсовую работу
func (m *StudentCourseworkManagerImpl) AssignStudentToCoursework(ctx context.Context, studentID, courseworkID uint) (*models.StudentCoursework, error) {
	// проверка: студент ещё не назначен
	if _, err := m.scRepo.GetByStudent(ctx, studentID); err == nil {
		return nil, errors.New("student already has an assigned coursework")
	}
	// проверка: свободные места
	cw, count, err := m.cwRepo.GetWithStudentCount(ctx, courseworkID)
	if err != nil {
		return nil, err
	}
	if count >= cw.MaxStudents {
		return nil, errors.New("no slots available for this coursework")
	}
	now := time.Now()
	assign := &models.StudentCoursework{
		StudentID:    studentID,
		CourseworkID: courseworkID,
		Status:       models.StatusAssigned,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := m.scRepo.Create(ctx, assign); err != nil {
		return nil, err
	}
	return assign, nil
}

// GetStudentCoursework возвращает текущее назначение студента
func (m *StudentCourseworkManagerImpl) GetStudentCoursework(ctx context.Context, studentID uint) (*models.StudentCoursework, error) {
	return m.scRepo.GetByStudent(ctx, studentID)
}

// UpdateCourseworkStatus обновляет статус выполнения
func (m *StudentCourseworkManagerImpl) UpdateCourseworkStatus(ctx context.Context, assignmentID uint, status models.CourseworkStatus) error {
	return m.scRepo.UpdateStatus(ctx, assignmentID, status)
}

// SubmitCoursework отмечает отправку курсовой работы
func (m *StudentCourseworkManagerImpl) SubmitCoursework(ctx context.Context, assignmentID uint) error {
	if err := m.scRepo.UpdateStatus(ctx, assignmentID, models.StatusSubmitted); err != nil {
		return err
	}
	return m.scRepo.SetSubmitted(ctx, assignmentID, time.Now())
}

// GradeCoursework выставляет оценку и фидбэк (для преподавателя)
func (m *StudentCourseworkManagerImpl) GradeCoursework(ctx context.Context, assignmentID uint, grade int, feedback string) error {
	if err := m.scRepo.UpdateStatus(ctx, assignmentID, models.StatusReviewed); err != nil {
		return err
	}
	return m.scRepo.SetGrade(ctx, assignmentID, grade, feedback)
}

// CompleteCoursework отмечает выполнение курсовой работы
func (m *StudentCourseworkManagerImpl) CompleteCoursework(ctx context.Context, assignmentID uint) error {
	if err := m.scRepo.UpdateStatus(ctx, assignmentID, models.StatusCompleted); err != nil {
		return err
	}
	return m.scRepo.SetCompleted(ctx, assignmentID, time.Now())
}

// GetTeacherCourseworks возвращает все задания для работ преподавателя
func (m *StudentCourseworkManagerImpl) GetTeacherCourseworks(ctx context.Context, teacherID uint) ([]models.StudentCoursework, error) {
	return m.scRepo.GetByTeacher(ctx, teacherID)
}

// GetCourseworkProgress собирает отчёт по прогрессу курсовой
func (m *StudentCourseworkManagerImpl) GetCourseworkProgress(ctx context.Context, courseworkID uint) (*interfaces.CourseworkProgressReport, error) {
	records, err := m.scRepo.GetByCoursework(ctx, courseworkID)
	if err != nil {
		return nil, err
	}
	report := &interfaces.CourseworkProgressReport{
		CourseworkID:     courseworkID,
		TotalStudents:    len(records),
		AssignedStudents: 0,
		InProgressCount:  0,
		SubmittedCount:   0,
		CompletedCount:   0,
		FailedCount:      0,
		Students:         make([]interfaces.StudentCourseworkSummary, 0, len(records)),
	}
	for _, r := range records {
		switch r.Status {
		case models.StatusAssigned:
			report.AssignedStudents++
		case models.StatusInProgress:
			report.InProgressCount++
		case models.StatusSubmitted:
			report.SubmittedCount++
		case models.StatusCompleted:
			report.CompletedCount++
		case models.StatusFailed:
			report.FailedCount++
		}
		report.Students = append(report.Students, interfaces.StudentCourseworkSummary{
			StudentID:   r.StudentID,
			StudentName: r.Student.GetFullName(),
			Status:      r.Status,
			AssignedAt:  r.CreatedAt,
			SubmittedAt: r.SubmittedAt,
			CompletedAt: r.CompletedAt,
			Grade:       r.Grade,
		})
	}
	return report, nil
}

// UnassignStudentFromCoursework отменяет назначение студента
func (m *StudentCourseworkManagerImpl) UnassignStudentFromCoursework(ctx context.Context, studentID uint) error {
	assignment, err := m.scRepo.GetByStudent(ctx, studentID)
	if err != nil {
		return err
	}
	return m.scRepo.Delete(ctx, assignment.ID)
}
