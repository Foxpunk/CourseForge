package models

import (
	"time"

	"gorm.io/gorm"
)

type StudentProfile struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time      `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	UserID        uint   `json:"user_id" gorm:"uniqueIndex;not null" validate:"required"`
	GroupID       uint   `json:"group_id" gorm:"not null" validate:"required"`
	StudentNumber string `json:"student_number" gorm:"size:20"`

	// Связи
	User         User         `json:"user" gorm:"-"`
	StudentGroup StudentGroup `json:"student_group" gorm:"-"`
}

func (StudentProfile) TableName() string {
	return "student_profiles"
}

type TeacherProfile struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time      `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	UserID         uint   `json:"user_id" gorm:"uniqueIndex;not null" validate:"required"`
	DepartmentID   uint   `json:"department_id" gorm:"not null" validate:"required"`
	Position       string `json:"position" gorm:"size:100"`
	AcademicDegree string `json:"academic_degree" gorm:"size:100"`

	// Связи
	User       User       `json:"user" gorm:"-"`
	Department Department `json:"department" gorm:"-"`
}

func (TeacherProfile) TableName() string {
	return "teacher_profiles"
}
