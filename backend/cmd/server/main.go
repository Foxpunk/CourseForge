package main

import (
	"log"

	"github.com/Foxpunk/courseforge/internal/config"
	"github.com/Foxpunk/courseforge/internal/drivers"
	"github.com/Foxpunk/courseforge/internal/handlers"
	"github.com/Foxpunk/courseforge/internal/managers"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := drivers.InitDB(cfg.Database.DSN)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	userRepo := drivers.NewUserRepository(db)
	subjectRepo := drivers.NewSubjectRepository(db)
	courseworkRepo := drivers.NewCourseworkRepository(db)
	studentCourseworkRepo := drivers.NewStudentCourseworkRepository(db)
	teacherSubjectRepo := drivers.NewTeacherSubjectRepository(db)
	teacherProfileRepo := drivers.NewTeacherProfileRepository(db)
	//studentProfileRepo = drivers.NewStudentProfileRepository(db)
	//studentGroupRepo = drivers.NewStudentGroupRepository(db)
	//departamentRepo = drivers.NewDepartmentRepository(db)
	// Initialize managers
	authManager := managers.NewAuthManager(userRepo, cfg.JWT)
	userManager := managers.NewUserManager(userRepo)
	subjectManager := managers.NewSubjectManager(subjectRepo, teacherSubjectRepo, teacherProfileRepo)
	courseworkManager := managers.NewCourseworkManager(courseworkRepo, studentCourseworkRepo)
	studentCourseworkManager := managers.NewStudentCourseworkManager(studentCourseworkRepo, courseworkRepo)
	//departamentManager = managers.NewDepartmentManager(departamentRepo,teacherProfileRepo)
	// Setup router
	router := handlers.NewRouter(
		authManager,
		userManager,
		subjectManager,
		courseworkManager,
		studentCourseworkManager,
		cfg.JWT.SecretKey,
	)

	addr := cfg.GetServerAddress()
	log.Printf("starting server on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
