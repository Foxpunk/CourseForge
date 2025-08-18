package managers

import (
	"context"
	"errors"

	"github.com/Foxpunk/courseforge/internal/interfaces"
	"github.com/Foxpunk/courseforge/internal/models"
)

// UserManagerImpl реализует интерфейс interfaces.UserManager
type UserManagerImpl struct {
	userRepo interfaces.UserRepository
}

// NewUserManager создаёт новый UserManager
func NewUserManager(userRepo interfaces.UserRepository) interfaces.UserManager {
	return &UserManagerImpl{userRepo: userRepo}
}

// CreateUser создаёт нового пользователя
func (m *UserManagerImpl) CreateUser(ctx context.Context, req interfaces.CreateUserRequest) (*models.User, error) {
	// проверка, что email уникален
	if _, err := m.userRepo.GetByEmail(ctx, req.Email); err == nil {
		return nil, errors.New("user with this email already exists")
	}

	// хешируем пароль
	hash, err := models.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        req.Email,
		PasswordHash: hash,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         req.Role,
		IsActive:     true,
	}

	if err := m.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// GetUser возвращает пользователя по ID
func (m *UserManagerImpl) GetUser(ctx context.Context, userID uint) (*models.User, error) {
	return m.userRepo.GetByID(ctx, userID)
}

// GetUserByEmail возвращает пользователя по email
func (m *UserManagerImpl) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return m.userRepo.GetByEmail(ctx, email)
}

// UpdateUser обновляет базовые поля пользователя
func (m *UserManagerImpl) UpdateUser(ctx context.Context, userID uint, req interfaces.UpdateUserRequest) (*models.User, error) {
	user, err := m.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}

	if err := m.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// DeleteUser удаляет пользователя
func (m *UserManagerImpl) DeleteUser(ctx context.Context, userID uint) error {
	return m.userRepo.Delete(ctx, userID)
}

// ListUsers возвращает список пользователей и общее число
func (m *UserManagerImpl) ListUsers(ctx context.Context, req interfaces.ListUsersRequest) ([]models.User, int, error) {
	users, err := m.userRepo.List(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, 0, err
	}
	// общее число равняется длине возвращённого слайса
	total := len(users)
	return users, total, nil
}

// AssignRole назначает роль пользователю
func (m *UserManagerImpl) AssignRole(ctx context.Context, userID uint, role models.UserRole) error {
	return m.userRepo.UpdateRole(ctx, userID, role)
}

// ActivateUser активирует пользователя
func (m *UserManagerImpl) ActivateUser(ctx context.Context, userID uint) error {
	return m.userRepo.SetActive(ctx, userID, true)
}

// DeactivateUser деактивирует пользователя
func (m *UserManagerImpl) DeactivateUser(ctx context.Context, userID uint) error {
	return m.userRepo.SetActive(ctx, userID, false)
}

// GetTeachers возвращает всех преподавателей
func (m *UserManagerImpl) GetTeachers(ctx context.Context) ([]models.User, error) {
	return m.userRepo.GetByRole(ctx, models.RoleTeacher)
}

// GetStudents возвращает всех студентов
func (m *UserManagerImpl) GetStudents(ctx context.Context) ([]models.User, error) {
	return m.userRepo.GetByRole(ctx, models.RoleStudent)
}

// GetAdmins возвращает всех админов
func (m *UserManagerImpl) GetAdmins(ctx context.Context) ([]models.User, error) {
	return m.userRepo.GetByRole(ctx, models.RoleAdmin)
}
