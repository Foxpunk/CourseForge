// internal/interfaces/managers.go
package interfaces

import (
	"context"

	"github.com/Foxpunk/courseforge/internal/models"
)

// AuthManager - интерфейс для аутентификации и авторизации
type AuthManager interface {
	Register(ctx context.Context, req RegisterRequest) (*models.User, error)
	Login(ctx context.Context, email, password string) (string, *models.User, error)
	ValidateToken(ctx context.Context, token string) (*models.User, error)
	ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error
	ResetPassword(ctx context.Context, email string) error
	RefreshToken(ctx context.Context, token string) (string, error)
	GenerateToken(user *models.User) (string, error)
}

// UserManager - интерфейс для управления пользователями
type UserManager interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (*models.User, error)
	GetUser(ctx context.Context, userID uint) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, userID uint, req UpdateUserRequest) (*models.User, error)
	DeleteUser(ctx context.Context, userID uint) error
	ListUsers(ctx context.Context, req ListUsersRequest) ([]models.User, int, error)

	// Управление ролями и статусами
	AssignRole(ctx context.Context, userID uint, role models.UserRole) error
	ActivateUser(ctx context.Context, userID uint) error
	DeactivateUser(ctx context.Context, userID uint) error

	// Получение пользователей по ролям
	GetTeachers(ctx context.Context) ([]models.User, error)
	GetStudents(ctx context.Context) ([]models.User, error)
	GetAdmins(ctx context.Context) ([]models.User, error)
}

// DepartmentManager - интерфейс для управления кафедрами
type DepartmentManager interface {
	CreateDepartment(ctx context.Context, req CreateDepartmentRequest) (*models.Department, error)
	GetDepartment(ctx context.Context, departmentID uint) (*models.Department, error)
	UpdateDepartment(ctx context.Context, departmentID uint, req UpdateDepartmentRequest) (*models.Department, error)
	DeleteDepartment(ctx context.Context, departmentID uint) error
	ListDepartments(ctx context.Context) ([]models.Department, error)

	// Привязка преподавателей к кафедре
	AssignTeacherToDepartment(ctx context.Context, teacherID, departmentID uint, position, degree string) error
	GetDepartmentTeachers(ctx context.Context, departmentID uint) ([]models.TeacherProfile, error)
}

// GroupManager - интерфейс для управления группами студентов
type GroupManager interface {
	CreateGroup(ctx context.Context, req CreateGroupRequest) (*models.StudentGroup, error)
	GetGroup(ctx context.Context, groupID uint) (*models.StudentGroup, error)
	UpdateGroup(ctx context.Context, groupID uint, req UpdateGroupRequest) (*models.StudentGroup, error)
	DeleteGroup(ctx context.Context, groupID uint) error
	ListGroups(ctx context.Context) ([]models.StudentGroup, error)
	GetGroupsByDepartment(ctx context.Context, departmentID uint) ([]models.StudentGroup, error)

	// Управление студентами в группах
	AssignStudentToGroup(ctx context.Context, studentID, groupID uint, studentNumber string) error
	GetGroupStudents(ctx context.Context, groupID uint) ([]models.StudentProfile, error)
	RemoveStudentFromGroup(ctx context.Context, studentID uint) error
}

// SubjectManager - интерфейс для управления дисциплинами
type SubjectManager interface {
	CreateSubject(ctx context.Context, req CreateSubjectRequest) (*models.Subject, error)
	GetSubject(ctx context.Context, subjectID uint) (*models.Subject, error)
	UpdateSubject(ctx context.Context, subjectID uint, req UpdateSubjectRequest) (*models.Subject, error)
	DeleteSubject(ctx context.Context, subjectID uint) error
	ListSubjects(ctx context.Context) ([]models.Subject, error)
	GetSubjectsByDepartment(ctx context.Context, departmentID uint) ([]models.Subject, error)
	GetSubjectsBySemester(ctx context.Context, semester int) ([]models.Subject, error)

	// Назначение преподавателей на дисциплины
	AssignTeacherToSubject(ctx context.Context, teacherID, subjectID uint, academicYear string, isLead bool) error
	GetTeacherSubjects(ctx context.Context, teacherID uint, academicYear string) ([]models.Subject, error)
	GetSubjectTeachers(ctx context.Context, subjectID uint, academicYear string) ([]models.TeacherProfile, error)
	SetLeadTeacher(ctx context.Context, teacherID, subjectID uint, academicYear string) error
	RemoveTeacherFromSubject(ctx context.Context, teacherID, subjectID uint, academicYear string) error
}

// CourseworkManager - интерфейс для управления курсовыми работами
type CourseworkManager interface {
	CreateCoursework(ctx context.Context, req CreateCourseworkRequest) (*models.Coursework, error)
	GetCoursework(ctx context.Context, courseworkID uint) (*models.Coursework, error)
	UpdateCoursework(ctx context.Context, courseworkID uint, req UpdateCourseworkRequest) (*models.Coursework, error)
	DeleteCoursework(ctx context.Context, courseworkID uint) error
	ListCourseworks(ctx context.Context, req ListCourseworksRequest) ([]models.Coursework, int, error)

	// Получение курсовых работ по различным критериям
	GetCourseworksBySubject(ctx context.Context, subjectID uint) ([]models.Coursework, error)
	GetCourseworksByTeacher(ctx context.Context, teacherID uint) ([]models.Coursework, error)
	GetAvailableCourseworks(ctx context.Context, studentID uint) ([]models.Coursework, error) // только доступные для конкретного студента

	// Управление доступностью
	SetCourseworkAvailability(ctx context.Context, courseworkID uint, available bool) error

	// Проверка возможности назначения
	CanAssignStudentToCoursework(ctx context.Context, studentID, courseworkID uint) error
}

// StudentCourseworkManager - интерфейс для управления назначениями студентов на курсовые
type StudentCourseworkManager interface {
	AssignStudentToCoursework(ctx context.Context, studentID, courseworkID uint) (*models.StudentCoursework, error)
	GetStudentCoursework(ctx context.Context, studentID uint) (*models.StudentCoursework, error)
	// Управление статусами
	UpdateCourseworkStatus(ctx context.Context, assignmentID uint, status models.CourseworkStatus) error
	SubmitCoursework(ctx context.Context, assignmentID uint) error

	// Оценивание (для преподавателей)
	GradeCoursework(ctx context.Context, assignmentID uint, grade int, feedback string) error
	CompleteCoursework(ctx context.Context, assignmentID uint) error

	// Отчеты и статистика
	GetTeacherCourseworks(ctx context.Context, teacherID uint) ([]models.StudentCoursework, error)
	GetCourseworkProgress(ctx context.Context, courseworkID uint) (*CourseworkProgressReport, error)

	// Отмена назначения
	UnassignStudentFromCoursework(ctx context.Context, studentID uint) error
}
