package drivers

import (
	"log"

	"github.com/Foxpunk/courseforge/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// ⚠️ Автоматическая миграция
	err = db.AutoMigrate(
		&models.Coursework{},
		&models.Department{},
		&models.StudentProfile{},
		&models.TeacherProfile{},
		&models.StudentCoursework{},
		&models.StudentGroup{},
		&models.Subject{},
		&models.TeacherSubject{},
		&models.User{})
	if err != nil {
		return nil, err
	}

	log.Println(" Миграция прошла успешно")
	return db, nil
}
