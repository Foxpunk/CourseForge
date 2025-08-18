package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleTeacher UserRole = "teacher"
	RoleStudent UserRole = "student"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time      `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Email        string   `json:"email" gorm:"uniqueIndex;not null;size:255" validate:"required,email"`
	PasswordHash string   `json:"-" gorm:"column:password_hash;not null;size:255" validate:"required,min=6"`
	FirstName    string   `json:"first_name" gorm:"not null;size:100" validate:"required,min=2,max=50"`
	LastName     string   `json:"last_name" gorm:"not null;size:100" validate:"required,min=2,max=50"`
	Role         UserRole `json:"role" gorm:"not null;size:20;check:role IN ('admin','teacher','student');default:'student'" validate:"required,oneof=admin teacher student"`
	IsActive     bool     `json:"is_active" gorm:"default:true"`

	TeacherSubjects   []Subject   `json:"teacher_subjects,omitempty" gorm:"many2many:teacher_subjects;"`
	StudentCoursework *Coursework `json:"student_coursework,omitempty" gorm:"-"`
}

// TableName задает имя таблицы в БД
func (User) TableName() string {
	return "users"
}

// IsAdmin проверяет, является ли пользователь администратором
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsTeacher проверяет, является ли пользователь преподавателем
func (u *User) IsTeacher() bool {
	return u.Role == RoleTeacher
}

// IsStudent проверяет, является ли пользователь студентом
func (u *User) IsStudent() bool {
	return u.Role == RoleStudent
}

// GetFullName возвращает полное имя пользователя
func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}
