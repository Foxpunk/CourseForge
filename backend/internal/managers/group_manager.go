package managers

import (
	"context"
	"errors"

	"github.com/Foxpunk/courseforge/internal/interfaces"
	"github.com/Foxpunk/courseforge/internal/models"
)

// GroupManagerImpl реализует интерфейс interfaces.GroupManager
type GroupManagerImpl struct {
	groupRepo   interfaces.StudentGroupRepository
	profileRepo interfaces.StudentProfileRepository
}

// NewGroupManager создаёт новый GroupManager
func NewGroupManager(
	groupRepo interfaces.StudentGroupRepository,
	profileRepo interfaces.StudentProfileRepository,
) interfaces.GroupManager {
	return &GroupManagerImpl{
		groupRepo:   groupRepo,
		profileRepo: profileRepo,
	}
}

// CreateGroup создаёт новую группу студентов
func (m *GroupManagerImpl) CreateGroup(ctx context.Context, req interfaces.CreateGroupRequest) (*models.StudentGroup, error) {
	// проверка уникальности кода
	if _, err := m.groupRepo.GetByCode(ctx, req.GroupCode); err == nil {
		return nil, errors.New("group code already exists")
	}
	group := &models.StudentGroup{
		GroupCode:    req.GroupCode,
		CourseYear:   req.CourseYear,
		Specialty:    req.Specialty,
		DepartmentID: req.DepartmentID,
	}
	if err := m.groupRepo.Create(ctx, group); err != nil {
		return nil, err
	}
	return group, nil
}

// GetGroup возвращает группу по ID
func (m *GroupManagerImpl) GetGroup(ctx context.Context, groupID uint) (*models.StudentGroup, error) {
	return m.groupRepo.GetByID(ctx, groupID)
}

// UpdateGroup обновляет данные группы
func (m *GroupManagerImpl) UpdateGroup(ctx context.Context, groupID uint, req interfaces.UpdateGroupRequest) (*models.StudentGroup, error) {
	group, err := m.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if req.GroupCode != nil {
		group.GroupCode = *req.GroupCode
	}
	if req.CourseYear != nil {
		group.CourseYear = *req.CourseYear
	}
	if req.Specialty != nil {
		group.Specialty = *req.Specialty
	}
	if req.DepartmentID != nil {
		group.DepartmentID = *req.DepartmentID
	}
	if err := m.groupRepo.Update(ctx, group); err != nil {
		return nil, err
	}
	return group, nil
}

// DeleteGroup удаляет группу по ID
func (m *GroupManagerImpl) DeleteGroup(ctx context.Context, groupID uint) error {
	return m.groupRepo.Delete(ctx, groupID)
}

// ListGroups возвращает все группы
func (m *GroupManagerImpl) ListGroups(ctx context.Context) ([]models.StudentGroup, error) {
	return m.groupRepo.List(ctx)
}

// GetGroupsByDepartment возвращает группы указанной кафедры
func (m *GroupManagerImpl) GetGroupsByDepartment(ctx context.Context, departmentID uint) ([]models.StudentGroup, error) {
	return m.groupRepo.GetByDepartment(ctx, departmentID)
}

// AssignStudentToGroup назначает студента в группу
func (m *GroupManagerImpl) AssignStudentToGroup(ctx context.Context, studentID, groupID uint, studentNumber string) error {
	if studentID == 0 || groupID == 0 {
		return errors.New("invalid student or group ID")
	}
	profile := &models.StudentProfile{
		UserID:        studentID,
		GroupID:       groupID,
		StudentNumber: studentNumber,
	}
	return m.profileRepo.Create(ctx, profile)
}

// GetGroupStudents возвращает профили студентов группы
func (m *GroupManagerImpl) GetGroupStudents(ctx context.Context, groupID uint) ([]models.StudentProfile, error) {
	return m.profileRepo.GetByGroup(ctx, groupID)
}

// RemoveStudentFromGroup удаляет студента из группы
func (m *GroupManagerImpl) RemoveStudentFromGroup(ctx context.Context, studentID uint) error {
	profile, err := m.profileRepo.GetByUserID(ctx, studentID)
	if err != nil {
		return err
	}
	return m.profileRepo.Delete(ctx, profile.ID)
}
