package models

import (
	"time"

	"gorm.io/gorm"
)

type StudentGroup struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time      `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	GroupCode    string     `json:"group_code" gorm:"uniqueIndex;not null;size:20"`
	CourseYear   int        `json:"course_year"`
	Specialty    string     `json:"specialty" gorm:"size:100"`
	DepartmentID uint       `json:"department_id" gorm:"not null;index"`
	Department   Department `json:"department" gorm:"-"`
}

func (StudentGroup) TableName() string {
	return "student_groups"
}
