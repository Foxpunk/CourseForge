package managers

import (
	"context"
	"errors"

	"github.com/Foxpunk/courseforge/internal/interfaces"
	"github.com/Foxpunk/courseforge/internal/models"
)

// DepartmentManagerImpl реализует interfaces.DepartmentManager
type DepartmentManagerImpl struct {
	deptRepo  interfaces.DepartmentRepository
	teachRepo interfaces.TeacherProfileRepository
}

// NewDepartmentManager создаёт DepartmentManager
func NewDepartmentManager(
	deptRepo interfaces.DepartmentRepository,
	teachRepo interfaces.TeacherProfileRepository,
) interfaces.DepartmentManager {
	return &DepartmentManagerImpl{deptRepo: deptRepo, teachRepo: teachRepo}
}

// CreateDepartment создаёт новую кафедру
func (m *DepartmentManagerImpl) CreateDepartment(
	ctx context.Context,
	req interfaces.CreateDepartmentRequest,
) (*models.Department, error) {
	// проверяем уникальность кода
	if _, err := m.deptRepo.GetByCode(ctx, req.DepartmentCode); err == nil {
		return nil, errors.New("department code already exists")
	}
	dept := &models.Department{
		DepartmentCode: req.DepartmentCode,
		DepartmentName: req.DepartmentName,
		Description:    req.Description,
	}
	if err := m.deptRepo.Create(ctx, dept); err != nil {
		return nil, err
	}
	return dept, nil
}

// GetDepartment возвращает кафедру по ID
func (m *DepartmentManagerImpl) GetDepartment(ctx context.Context, departmentID uint) (*models.Department, error) {
	return m.deptRepo.GetByID(ctx, departmentID)
}

// UpdateDepartment обновляет кафедру
func (m *DepartmentManagerImpl) UpdateDepartment(
	ctx context.Context,
	departmentID uint,
	req interfaces.UpdateDepartmentRequest,
) (*models.Department, error) {
	dept, err := m.deptRepo.GetByID(ctx, departmentID)
	if err != nil {
		return nil, err
	}
	if req.DepartmentCode != nil {
		dept.DepartmentCode = *req.DepartmentCode
	}
	if req.DepartmentName != nil {
		dept.DepartmentName = *req.DepartmentName
	}
	if req.Description != nil {
		dept.Description = *req.Description
	}
	if err := m.deptRepo.Update(ctx, dept); err != nil {
		return nil, err
	}
	return dept, nil
}

// DeleteDepartment удаляет кафедру
func (m *DepartmentManagerImpl) DeleteDepartment(ctx context.Context, departmentID uint) error {
	return m.deptRepo.Delete(ctx, departmentID)
}

// ListDepartments возвращает все кафедры
func (m *DepartmentManagerImpl) ListDepartments(ctx context.Context) ([]models.Department, error) {
	return m.deptRepo.List(ctx)
}

// AssignTeacherToDepartment привязывает преподавателя к кафедре
func (m *DepartmentManagerImpl) AssignTeacherToDepartment(
	ctx context.Context,
	teacherID, departmentID uint,
	position, degree string,
) error {
	if teacherID == 0 || departmentID == 0 {
		return errors.New("invalid teacher or department ID")
	}
	profile := &models.TeacherProfile{
		UserID:         teacherID,
		DepartmentID:   departmentID,
		Position:       position,
		AcademicDegree: degree,
	}
	return m.teachRepo.Create(ctx, profile)
}

// GetDepartmentTeachers возвращает профили преподавателей кафедры
func (m *DepartmentManagerImpl) GetDepartmentTeachers(ctx context.Context, departmentID uint) ([]models.TeacherProfile, error) {
	return m.teachRepo.GetByDepartment(ctx, departmentID)
}
