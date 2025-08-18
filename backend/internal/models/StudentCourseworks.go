package models

import (
	"time"

	"gorm.io/gorm"
)

// CourseworkStatus описывает статус выполнения курсовой работы
type CourseworkStatus string

const (
	StatusAssigned   CourseworkStatus = "assigned"
	StatusInProgress CourseworkStatus = "in_progress"
	StatusSubmitted  CourseworkStatus = "submitted"
	StatusReviewed   CourseworkStatus = "reviewed"
	StatusCompleted  CourseworkStatus = "completed"
	StatusFailed     CourseworkStatus = "failed"
)

// StudentCoursework представляет назначение студента на курсовую работу
type StudentCoursework struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time      `json:"assigned_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	StudentID    uint             `json:"student_id" gorm:"not null" validate:"required"`
	CourseworkID uint             `json:"coursework_id" gorm:"not null" validate:"required"`
	Status       CourseworkStatus `json:"status" gorm:"type:varchar(20);default:'assigned'" validate:"oneof=assigned in_progress submitted reviewed completed failed"`
	// Coursework
	SubmittedAt *time.Time `json:"submitted_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Grade       *int       `json:"grade,omitempty" validate:"omitempty,min=2,max=5"`
	Feedback    string     `json:"feedback,omitempty" gorm:"type:text"`

	// Связи
	Student    User       `json:"student" gorm:"foreignKey:StudentID"`
	Coursework Coursework `json:"coursework" gorm:"foreignKey:CourseworkID"`
}

// TableName задаёт имя таблицы в БД
func (StudentCoursework) TableName() string {
	return "student_courseworks"
}
