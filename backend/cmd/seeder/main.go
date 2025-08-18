package main

import (
	"errors" // <- добавили
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/Foxpunk/courseforge/internal/config"
	"github.com/Foxpunk/courseforge/internal/models"
)

func main() {
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatal("Invalid config:", err)
	}

	db, err := gorm.Open(sqlite.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	if err := db.AutoMigrate(
		&models.Coursework{},
		&models.Department{},
		&models.StudentProfile{},
		&models.TeacherProfile{},
		&models.StudentCoursework{},
		&models.StudentGroup{},
		&models.Subject{},
		&models.TeacherSubject{},
		&models.User{},
	); err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}

	// Seeder: admin user
	email := "admin@example.com"
	password := "admin123"

	var existing models.User
	err = db.Where("email = ?", email).First(&existing).Error
	if err == nil {
		log.Println("⚠️ admin already exists, skipping seeder.")
		return
	}
	// Теперь проверяем любым err-контейнером
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Fatal("Error checking admin existence:", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Error hashing password:", err)
	}

	admin := models.User{
		Email:        email,
		PasswordHash: string(hash),
		FirstName:    "Super",
		LastName:     "Admin",
		Role:         models.RoleAdmin,
		IsActive:     true,
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Fatal("Error creating admin user:", err)
	}
	log.Println("✅ Admin user created: email=", email, "password=", password)
}
