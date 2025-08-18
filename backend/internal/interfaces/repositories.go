package interfaces

import (
	"context"
	"time"

	"github.com/Foxpunk/courseforge/internal/models"
)

type BaseRepository[T any] interface {
	Create(ctx context.Context, e *T) error
	GetByID(ctx context.Context, id uint) (*T, error)
	Update(ctx context.Context, e *T) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, limit, offset int) ([]T, error)
}
type Validator interface {
	Validate() error
}

// UserRepository - интерфейс для работы с пользователями
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uint) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, limit, offset int) ([]models.User, error)
	GetByRole(ctx context.Context, role models.UserRole) ([]models.User, error)
	UpdateRole(ctx context.Context, userID uint, role models.UserRole) error
	SetActive(ctx context.Context, userID uint, active bool) error
}

// DepartmentRepository - интерфейс для работы с кафедрами
type DepartmentRepository interface {
	Create(ctx context.Context, department *models.Department) error
	GetByID(ctx context.Context, id uint) (*models.Department, error)
	GetByCode(ctx context.Context, code string) (*models.Department, error)
	Update(ctx context.Context, department *models.Department) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context) ([]models.Department, error)
}

// StudentGroupRepository - интерфейс для работы с группами студентов
type StudentGroupRepository interface {
	Create(ctx context.Context, group *models.StudentGroup) error
	GetByID(ctx context.Context, id uint) (*models.StudentGroup, error)
	GetByCode(ctx context.Context, code string) (*models.StudentGroup, error)
	Update(ctx context.Context, group *models.StudentGroup) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context) ([]models.StudentGroup, error)
	GetByDepartment(ctx context.Context, departmentID uint) ([]models.StudentGroup, error)
}

// StudentProfileRepository - интерфейс для работы с профилями студентов
type StudentProfileRepository interface {
	Create(ctx context.Context, profile *models.StudentProfile) error
	GetByUserID(ctx context.Context, userID uint) (*models.StudentProfile, error)
	GetByID(ctx context.Context, id uint) (*models.StudentProfile, error)
	Update(ctx context.Context, profile *models.StudentProfile) error
	Delete(ctx context.Context, id uint) error
	GetByGroup(ctx context.Context, groupID uint) ([]models.StudentProfile, error)
}

// TeacherProfileRepository - интерфейс для работы с профилями преподавателей
type TeacherProfileRepository interface {
	Create(ctx context.Context, profile *models.TeacherProfile) error
	GetByUserID(ctx context.Context, userID uint) (*models.TeacherProfile, error)
	GetByID(ctx context.Context, id uint) (*models.TeacherProfile, error)
	Update(ctx context.Context, profile *models.TeacherProfile) error
	Delete(ctx context.Context, id uint) error
	GetByDepartment(ctx context.Context, departmentID uint) ([]models.TeacherProfile, error)
}

// SubjectRepository - интерфейс для работы с дисциплинами
type SubjectRepository interface {
	Create(ctx context.Context, subject *models.Subject) error
	GetByID(ctx context.Context, id uint) (*models.Subject, error)
	Update(ctx context.Context, subject *models.Subject) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context) ([]models.Subject, error)
	GetByDepartment(ctx context.Context, departmentID uint) ([]models.Subject, error)
	GetBySemester(ctx context.Context, semester int) ([]models.Subject, error)
	SetActive(ctx context.Context, subjectID uint, active bool) error
}

// TeacherSubjectRepository - интерфейс для назначения преподавателей на дисциплины
type TeacherSubjectRepository interface {
	Create(ctx context.Context, assignment *models.TeacherSubject) error
	GetByID(ctx context.Context, id uint) (*models.TeacherSubject, error)
	Delete(ctx context.Context, id uint) error
	GetByTeacher(ctx context.Context, teacherID uint, academicYear string) ([]models.TeacherSubject, error)
	GetBySubject(ctx context.Context, subjectID uint, academicYear string) ([]models.TeacherSubject, error)
	GetLeadTeacher(ctx context.Context, subjectID uint, academicYear string) (*models.TeacherSubject, error)
	SetLead(ctx context.Context, assignmentID uint, isLead bool) error
	DeleteByTeacherAndSubject(ctx context.Context, teacherID, subjectID uint) error
}

// CourseworkRepository - интерфейс для работы с курсовыми работами
type CourseworkRepository interface {
	Create(ctx context.Context, coursework *models.Coursework) error
	GetByID(ctx context.Context, id uint) (*models.Coursework, error)
	Update(ctx context.Context, coursework *models.Coursework) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, limit, offset int) ([]models.Coursework, error)
	GetBySubject(ctx context.Context, subjectID uint) ([]models.Coursework, error)
	GetByTeacher(ctx context.Context, teacherID uint) ([]models.Coursework, error)
	GetAvailable(ctx context.Context, subjectID uint) ([]models.Coursework, error)
	SetAvailable(ctx context.Context, courseworkID uint, available bool) error
	GetWithStudentCount(ctx context.Context, courseworkID uint) (*models.Coursework, int, error)
}

// StudentCourseworkRepository - интерфейс для назначения студентов на курсовые
type StudentCourseworkRepository interface {
	Create(ctx context.Context, assignment *models.StudentCoursework) error
	GetByID(ctx context.Context, id uint) (*models.StudentCoursework, error)
	GetByStudent(ctx context.Context, studentID uint) (*models.StudentCoursework, error)
	GetByCoursework(ctx context.Context, courseworkID uint) ([]models.StudentCoursework, error)
	Update(ctx context.Context, assignment *models.StudentCoursework) error
	Delete(ctx context.Context, id uint) error
	UpdateStatus(ctx context.Context, id uint, status models.CourseworkStatus) error
	SetGrade(ctx context.Context, id uint, grade int, feedback string) error
	SetSubmitted(ctx context.Context, id uint, submittedAt time.Time) error
	SetCompleted(ctx context.Context, id uint, completedAt time.Time) error
	GetByTeacher(ctx context.Context, teacherID uint) ([]models.StudentCoursework, error)
}
