package models

// TeacherSubject связывает преподавателя и дисциплину на конкретный учебный год.
type TeacherSubject struct {
	UserID    uint `json:"user_id" gorm:"column:user_id;primaryKey"`
	SubjectID uint `json:"subject_id" gorm:"column:subject_id;primaryKey"`

	// Отношения
	Teacher User    `json:"teacher" gorm:"foreignKey:UserID;references:ID"`
	Subject Subject `json:"subject" gorm:"foreignKey:SubjectID;references:ID"`
}

// TableName задаёт имя таблицы в БД
func (TeacherSubject) TableName() string {
	return "teacher_subjects"
}
