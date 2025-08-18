package managers

import (
	"context"
	"errors"
	"log"

	"github.com/Foxpunk/courseforge/internal/interfaces"
	"github.com/Foxpunk/courseforge/internal/models"
)

// SubjectManagerImpl реализует interfaces.SubjectManager
type SubjectManagerImpl struct {
	subjRepo   interfaces.SubjectRepository
	assignRepo interfaces.TeacherSubjectRepository
	profRepo   interfaces.TeacherProfileRepository
}

// NewSubjectManager создаёт новый SubjectManager
func NewSubjectManager(
	subjRepo interfaces.SubjectRepository,
	assignRepo interfaces.TeacherSubjectRepository,
	profRepo interfaces.TeacherProfileRepository,
) interfaces.SubjectManager {
	return &SubjectManagerImpl{
		subjRepo:   subjRepo,
		assignRepo: assignRepo,
		profRepo:   profRepo,
	}
}

// CreateSubject создаёт новую дисциплину
func (m *SubjectManagerImpl) CreateSubject(ctx context.Context, req interfaces.CreateSubjectRequest) (*models.Subject, error) {
	subj := &models.Subject{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Semester:    req.Semester,
		IsActive:    true,
	}
	if err := m.subjRepo.Create(ctx, subj); err != nil {
		return nil, err
	}
	return subj, nil
}

// GetSubject возвращает дисциплину по ID
func (m *SubjectManagerImpl) GetSubject(ctx context.Context, subjectID uint) (*models.Subject, error) {
	return m.subjRepo.GetByID(ctx, subjectID)
}

// UpdateSubject обновляет дисциплину
func (m *SubjectManagerImpl) UpdateSubject(ctx context.Context, subjectID uint, req interfaces.UpdateSubjectRequest) (*models.Subject, error) {
	subj, err := m.subjRepo.GetByID(ctx, subjectID)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		subj.Name = *req.Name
	}
	if req.Description != nil {
		subj.Description = *req.Description
	}
	if req.Semester != nil {
		subj.Semester = *req.Semester
	}
	if req.IsActive != nil {
		subj.IsActive = *req.IsActive
	}
	if err := m.subjRepo.Update(ctx, subj); err != nil {
		return nil, err
	}
	return subj, nil
}

// DeleteSubject удаляет дисциплину
func (m *SubjectManagerImpl) DeleteSubject(ctx context.Context, subjectID uint) error {
	return m.subjRepo.Delete(ctx, subjectID)
}

// ListSubjects возвращает все дисциплины
func (m *SubjectManagerImpl) ListSubjects(ctx context.Context) ([]models.Subject, error) {
	return m.subjRepo.List(ctx)
}

// GetSubjectsByDepartment возвращает дисциплины по кафедре
func (m *SubjectManagerImpl) GetSubjectsByDepartment(ctx context.Context, departmentID uint) ([]models.Subject, error) {
	return m.subjRepo.GetByDepartment(ctx, departmentID)
}

// GetSubjectsBySemester возвращает дисциплины по семестру
func (m *SubjectManagerImpl) GetSubjectsBySemester(ctx context.Context, semester int) ([]models.Subject, error) {
	return m.subjRepo.GetBySemester(ctx, semester)
}

// AssignTeacherToSubject привязывает преподавателя к дисциплине
func (m *SubjectManagerImpl) AssignTeacherToSubject(
	ctx context.Context,
	teacherID, subjectID uint,
	academicYear string,
	isLead bool,
) error {
	if teacherID == 0 || subjectID == 0 {
		return errors.New("invalid parameters for assignment")
	}

	log.Printf("Creating TeacherSubject: UserID=%d, SubjectID=%d", teacherID, subjectID)

	assign := &models.TeacherSubject{
		UserID:    teacherID, // Изменено на UserID
		SubjectID: subjectID,
		// academicYear и isLead игнорируем - их нет в таблице
	}
	log.Printf("TeacherSubject object: %+v", assign)

	err := m.assignRepo.Create(ctx, assign)
	if err != nil {
		log.Printf("Error creating TeacherSubject: %v", err)
		return err
	}

	log.Printf("TeacherSubject created successfully")
	return nil
}

// GetTeacherSubjects возвращает дисциплины преподавателя за год
func (m *SubjectManagerImpl) GetTeacherSubjects(ctx context.Context, teacherID uint, academicYear string) ([]models.Subject, error) {
	assigns, err := m.assignRepo.GetByTeacher(ctx, teacherID, academicYear)
	if err != nil {
		return nil, err
	}
	var subjects []models.Subject
	for _, a := range assigns {
		subjects = append(subjects, a.Subject)
	}
	return subjects, nil
}

// GetSubjectTeachers возвращает профили преподавателей дисциплины за год
func (m *SubjectManagerImpl) GetSubjectTeachers(ctx context.Context, subjectID uint, academicYear string) ([]models.TeacherProfile, error) {
	assigns, err := m.assignRepo.GetBySubject(ctx, subjectID, academicYear)
	if err != nil {
		return nil, err
	}
	var profiles []models.TeacherProfile
	for _, a := range assigns {
		if p, err := m.profRepo.GetByUserID(ctx, a.UserID); err == nil { // Изменено на UserID
			profiles = append(profiles, *p)
		}
	}
	return profiles, nil
}

// SetLeadTeacher меняет ведущего преподавателя (не работает с простой таблицей)
func (m *SubjectManagerImpl) SetLeadTeacher(ctx context.Context, teacherID, subjectID uint, academicYear string) error {
	return errors.New("lead teacher functionality not supported with simple teacher_subjects table")
}

// RemoveTeacherFromSubject удаляет преподавателя с дисциплины
func (m *SubjectManagerImpl) RemoveTeacherFromSubject(ctx context.Context, teacherID, subjectID uint, academicYear string) error {
	return m.assignRepo.DeleteByTeacherAndSubject(ctx, teacherID, subjectID)
}
