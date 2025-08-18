package models

import (
	"time"

	"gorm.io/gorm"
)

// Subject представляет дисциплину в системе
type Subject struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Name        string `json:"name" gorm:"not null" validate:"required,min=3,max=100"`
	Code        string `json:"code" gorm:"uniqueIndex;not null" validate:"required,min=2,max=20"`
	Description string `json:"description" gorm:"type:text"`
	Semester    int    `json:"semester" gorm:"not null" validate:"required,min=1,max=12"`
	IsActive    bool   `json:"is_active" gorm:"default:true"`

	// Отношения
	Teachers    []User       `json:"teachers,omitempty" gorm:"many2many:teacher_subjects;"`
	Courseworks []Coursework `json:"courseworks,omitempty" gorm:"foreignKey:SubjectID"`
}

// TableName задает имя таблицы в БД
func (Subject) TableName() string {
	return "subjects"
}
