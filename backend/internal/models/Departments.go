package models

import (
	"time"

	"gorm.io/gorm"
)

type Department struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time      `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	DepartmentCode string `json:"department_code" gorm:"uniqueIndex;not null;size:10" validate:"required"`
	DepartmentName string `json:"department_name" gorm:"not null;size:200" validate:"required"`
	Description    string `json:"description,omitempty"`

	TeacherProfiles []TeacherProfile `json:"teacher_profiles,omitempty" gorm:"-"`
	Subjects        []Subject        `json:"subjects,omitempty" gorm:"-"`
	StudentGroups   []StudentGroup   `json:"student_groups,omitempty" gorm:"-"`
}

func (Department) TableName() string {
	return "departments"
}
