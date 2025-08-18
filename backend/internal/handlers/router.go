package handlers

import (
	"log"

	"github.com/Foxpunk/courseforge/internal/interfaces"
	"github.com/gin-gonic/gin"
)

// NewRouter собирает все маршруты и middleware
func NewRouter(
	authManager interfaces.AuthManager,
	userManager interfaces.UserManager,
	subjectManager interfaces.SubjectManager,
	courseworkManager interfaces.CourseworkManager,
	studentCourseworkManager interfaces.StudentCourseworkManager,
	jwtSecret string,
) *gin.Engine {
	// создаём gin
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// Инициализируем middleware и хендлеры
	mw := NewMiddleware(authManager)
	authH := NewAuthHandler(authManager, userManager, jwtSecret)
	userH := NewUserHandler(userManager)
	discH := NewDisciplineHandler(subjectManager)
	projH := NewProjectHandler(courseworkManager, studentCourseworkManager)

	// При необходимости включить CORS
	r.Use(mw.CORS())

	api := r.Group("/api/v1")
	api.GET("/health", func(c *gin.Context) {
		log.Println("Health check called")
		c.JSON(200, gin.H{"status": "OK", "message": "Server is running"})
	})

	api.POST("/test", func(c *gin.Context) {
		log.Println("Test POST called")
		var body map[string]interface{}
		if err := c.ShouldBindJSON(&body); err != nil {
			log.Printf("Test POST bind error: %v", err)
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		log.Printf("Test POST body: %+v", body)
		c.JSON(200, gin.H{"received": body})
	})
	// AUTH (публичные)
	auth := api.Group("/auth")
	{
		auth.POST("/login", authH.Login)
		auth.POST("/register", authH.Register)
		auth.POST("/refresh", authH.RefreshToken)
		auth.POST("/reset-password", authH.ResetPassword)
	}

	// PROFILE (требует авторизацию)
	profile := api.Group("/profile", mw.AuthMiddleware())
	{
		profile.GET("", authH.GetProfile)
		profile.POST("/change-password", authH.ChangePassword)
		profile.POST("/logout", authH.Logout)
	}

	// USERS (admin only)
	users := api.Group("/users", mw.AuthMiddleware(), mw.AdminRequired())
	{
		users.POST("", userH.CreateUser)
		users.GET("", userH.ListUsers)
		users.GET("/:id", userH.GetUser)
		users.PUT("/:id", userH.UpdateUser)
		users.DELETE("/:id", userH.DeleteUser)
	}

	// SUBJECTS / DISCIPLINE
	subj := api.Group("/subjects")
	{
		// общий доступ по авторизации
		subj.GET("", mw.AuthMiddleware(), discH.GetDisciplines)
		subj.GET("/:id", mw.AuthMiddleware(), discH.GetDiscipline)

		// admin only
		adminSubj := subj.Group("", mw.AuthMiddleware(), mw.AdminRequired())
		{
			adminSubj.POST("", discH.CreateDiscipline)
			adminSubj.PUT("/:id", discH.UpdateDiscipline)
			adminSubj.DELETE("/:id", discH.DeleteDiscipline)
			adminSubj.POST("/:id/teachers", discH.AssignTeacher)
			adminSubj.DELETE("/:id/teachers/:teacherId", discH.RemoveTeacher)
			adminSubj.PUT("/:id/lead-teacher", discH.SetLeadTeacher)
		}
	}

	// COURSEWORKS / PROJECTS
	cw := api.Group("/courseworks")
	{
		// авторизованные
		cw.GET("", mw.AuthMiddleware(), projH.GetProjects)
		cw.GET("/available", mw.AuthMiddleware(), projH.GetAvailableProjects)
		cw.GET("/:id", mw.AuthMiddleware(), projH.GetProject)

		// teacher or admin
		tAdmin := cw.Group("", mw.AuthMiddleware(), mw.TeacherOrAdminRequired())
		{
			tAdmin.POST("", projH.CreateProject)
			tAdmin.PUT("/:id", projH.UpdateProject)
			tAdmin.DELETE("/:id", projH.DeleteProject)
			tAdmin.PUT("/:id/availability", projH.SetProjectAvailability)
		}

		// student only
		stud := cw.Group("/:id/assign", mw.AuthMiddleware(), mw.StudentRequired())
		{
			stud.POST("", projH.AssignStudent)
		}
	}

	return r
}
