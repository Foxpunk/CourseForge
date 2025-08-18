package models

import (
	"time"

	"gorm.io/gorm"
)

// DifficultyLevel перечисление уровня сложности проекта
type DifficultyLevel string

const (
	Easy   DifficultyLevel = "easy"
	Medium DifficultyLevel = "medium"
	Hard   DifficultyLevel = "hard"
)

// Coursework представляет курсовой проект (тему)
type Coursework struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time      `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Title           string          `json:"title" gorm:"size:300;not null" validate:"required,min=5,max=300"`
	Description     string          `json:"description" gorm:"type:text;not null" validate:"required,min=20"`
	Requirements    string          `json:"requirements" gorm:"type:text"`
	SubjectID       uint            `json:"subject_id" gorm:"not null" validate:"required"`
	TeacherID       uint            `json:"teacher_id" gorm:"not null" validate:"required"`
	MaxStudents     int             `json:"max_students" gorm:"default:1" validate:"min=1,max=10"`
	DifficultyLevel DifficultyLevel `json:"difficulty_level" gorm:"type:varchar(20);check:difficulty_level IN ('easy','medium','hard')"`
	IsAvailable     bool            `json:"is_available" gorm:"default:true"`

	// Связи
	Subject Subject `json:"subject" gorm:"foreignKey:SubjectID"`
	Teacher User    `json:"teacher" gorm:"foreignKey:TeacherID"`
}

func (Coursework) TableName() string {
	return "courseworks"
}
