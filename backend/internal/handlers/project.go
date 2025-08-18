package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/Foxpunk/courseforge/internal/interfaces"
	"github.com/Foxpunk/courseforge/internal/models"
)

// ProjectHandler управляет CRUD операциями с проектами (курсовыми работами)
type ProjectHandler struct {
	courseworkManager        interfaces.CourseworkManager
	studentCourseworkManager interfaces.StudentCourseworkManager
	validator                *validator.Validate
}

// NewProjectHandler создаёт новый ProjectHandler
func NewProjectHandler(
	cwManager interfaces.CourseworkManager,
	scManager interfaces.StudentCourseworkManager,
) *ProjectHandler {
	return &ProjectHandler{
		courseworkManager:        cwManager,
		studentCourseworkManager: scManager,
		validator:                validator.New(),
	}
}

// CreateProject создаёт новый проект (преподаватель)
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req interfaces.CreateCourseworkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Валидация
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем текущего пользователя
	user := h.getCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Проверяем, что teacher_id соответствует текущему пользователю (для безопасности)
	if !user.IsAdmin() && req.TeacherID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "can only create projects for yourself"})
		return
	}

	coursework, err := h.courseworkManager.CreateCoursework(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Формируем ответ
	response := h.buildCourseworkResponse(coursework)
	c.JSON(http.StatusCreated, response)
}

// GetProjects возвращает список проектов с фильтрацией
func (h *ProjectHandler) GetProjects(c *gin.Context) {
	// Парсим параметры запроса
	var req interfaces.ListCourseworksRequest

	// Параметры фильтрации
	if subjectID := c.Query("subject_id"); subjectID != "" {
		if id, err := strconv.ParseUint(subjectID, 10, 32); err == nil {
			subjID := uint(id)
			req.SubjectID = &subjID
		}
	}

	if teacherID := c.Query("teacher_id"); teacherID != "" {
		if id, err := strconv.ParseUint(teacherID, 10, 32); err == nil {
			tchrID := uint(id)
			req.TeacherID = &tchrID
		}
	}

	if available := c.Query("available"); available != "" {
		if avl, err := strconv.ParseBool(available); err == nil {
			req.Available = &avl
		}
	}

	if difficulty := c.Query("difficulty"); difficulty != "" {
		diff := models.DifficultyLevel(difficulty)
		req.Difficulty = &diff
	}

	// Пагинация
	req.Limit = 20 // по умолчанию
	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
			req.Limit = l
		}
	}

	req.Offset = 0 // по умолчанию
	if offset := c.Query("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil && o >= 0 {
			req.Offset = o
		}
	}

	// Валидация
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем список
	courseworks, total, err := h.courseworkManager.ListCourseworks(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Формируем ответ
	response := interfaces.CourseworkListResponse{
		Courseworks: make([]interfaces.CourseworkResponse, len(courseworks)),
		Total:       int64(total),
	}

	for i, cw := range courseworks {
		response.Courseworks[i] = h.buildCourseworkResponse(&cw)
	}

	c.JSON(http.StatusOK, response)
}

// GetProject возвращает детали конкретного проекта
func (h *ProjectHandler) GetProject(c *gin.Context) {
	idParam := c.Param("id")
	cwID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	coursework, err := h.courseworkManager.GetCoursework(c.Request.Context(), uint(cwID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	response := h.buildCourseworkResponse(coursework)
	c.JSON(http.StatusOK, response)
}

// UpdateProject обновляет проект (только владелец или админ)
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	idParam := c.Param("id")
	cwID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	var req interfaces.UpdateCourseworkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Валидация
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем права доступа
	user := h.getCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Получаем текущий проект для проверки владельца
	currentCw, err := h.courseworkManager.GetCoursework(c.Request.Context(), uint(cwID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	// Проверяем права (только владелец или админ может редактировать)
	if !user.IsAdmin() && currentCw.TeacherID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	coursework, err := h.courseworkManager.UpdateCoursework(c.Request.Context(), uint(cwID), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := h.buildCourseworkResponse(coursework)
	c.JSON(http.StatusOK, response)
}

// DeleteProject удаляет проект (только владелец или админ)
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	idParam := c.Param("id")
	cwID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	// Проверяем права доступа
	user := h.getCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Получаем текущий проект для проверки владельца
	currentCw, err := h.courseworkManager.GetCoursework(c.Request.Context(), uint(cwID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	// Проверяем права (только владелец или админ может удалять)
	if !user.IsAdmin() && currentCw.TeacherID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	err = h.courseworkManager.DeleteCoursework(c.Request.Context(), uint(cwID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// AssignStudent назначает студента на проект
func (h *ProjectHandler) AssignStudent(c *gin.Context) {
	idParam := c.Param("id")
	cwID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	var req interfaces.AssignStudentToCourseworkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Валидация
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем права доступа
	user := h.getCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Студент может назначить только себя
	if user.IsStudent() && req.StudentID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "students can only assign themselves"})
		return
	}

	// Используем ID из URL, а не из тела запроса для безопасности
	assignment, err := h.studentCourseworkManager.AssignStudentToCoursework(
		c.Request.Context(),
		req.StudentID,
		uint(cwID),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Формируем ответ
	response := h.buildStudentCourseworkResponse(assignment)
	c.JSON(http.StatusCreated, response)
}

// GetAvailableProjects возвращает доступные проекты для студента
func (h *ProjectHandler) GetAvailableProjects(c *gin.Context) {
	user := h.getCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	courseworks, err := h.courseworkManager.GetAvailableCourseworks(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]interfaces.CourseworkResponse, len(courseworks))
	for i, cw := range courseworks {
		response[i] = h.buildCourseworkResponse(&cw)
	}

	c.JSON(http.StatusOK, response)
}

// SetProjectAvailability изменяет доступность проекта (только владелец или админ)
func (h *ProjectHandler) SetProjectAvailability(c *gin.Context) {
	idParam := c.Param("id")
	cwID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	var req struct {
		Available bool `json:"available"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем права доступа
	user := h.getCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Получаем текущий проект для проверки владельца
	currentCw, err := h.courseworkManager.GetCoursework(c.Request.Context(), uint(cwID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	// Проверяем права (только владелец или админ может изменять доступность)
	if !user.IsAdmin() && currentCw.TeacherID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	err = h.courseworkManager.SetCourseworkAvailability(c.Request.Context(), uint(cwID), req.Available)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Вспомогательные методы

// getCurrentUser извлекает текущего пользователя из контекста
func (h *ProjectHandler) getCurrentUser(c *gin.Context) *models.User {
	raw, exists := c.Get("user")
	if !exists {
		return nil
	}
	user, ok := raw.(*models.User)
	if !ok {
		return nil
	}
	return user
}

// buildCourseworkResponse создаёт ответ для курсовой работы
func (h *ProjectHandler) buildCourseworkResponse(cw *models.Coursework) interfaces.CourseworkResponse {
	return interfaces.CourseworkResponse{
		ID:              cw.ID,
		Title:           cw.Title,
		Description:     cw.Description,
		Requirements:    cw.Requirements,
		MaxStudents:     cw.MaxStudents,
		DifficultyLevel: cw.DifficultyLevel,
		IsAvailable:     cw.IsAvailable,
		Subject: interfaces.SubjectResponse{
			ID:          cw.Subject.ID,
			Name:        cw.Subject.Name,
			Code:        cw.Subject.Code,
			Description: cw.Subject.Description,
			Semester:    cw.Subject.Semester,
			IsActive:    cw.Subject.IsActive,
			Teachers:    []interfaces.UserResponse{}, // TODO: загрузить при необходимости
			CreatedAt:   cw.Subject.CreatedAt.Format(time.RFC3339),
		},
		Teacher: interfaces.UserResponse{
			ID:        cw.Teacher.ID,
			Email:     cw.Teacher.Email,
			FirstName: cw.Teacher.FirstName,
			LastName:  cw.Teacher.LastName,
			Role:      cw.Teacher.Role,
			IsActive:  cw.Teacher.IsActive,
			CreatedAt: cw.Teacher.CreatedAt.Format(time.RFC3339),
		},
		CreatedAt: cw.CreatedAt,
		UpdatedAt: cw.UpdatedAt,
	}
}

// buildStudentCourseworkResponse создаёт ответ для назначения студента
func (h *ProjectHandler) buildStudentCourseworkResponse(sc *models.StudentCoursework) interfaces.StudentCourseworkResponse {
	return interfaces.StudentCourseworkResponse{
		ID: sc.ID,
		Student: interfaces.UserResponse{
			ID:        sc.Student.ID,
			Email:     sc.Student.Email,
			FirstName: sc.Student.FirstName,
			LastName:  sc.Student.LastName,
			Role:      sc.Student.Role,
			IsActive:  sc.Student.IsActive,
			CreatedAt: sc.Student.CreatedAt.Format(time.RFC3339),
		},
		Coursework:  h.buildCourseworkResponse(&sc.Coursework),
		Status:      sc.Status,
		Grade:       sc.Grade,
		Feedback:    &sc.Feedback,
		AssignedAt:  sc.CreatedAt,
		SubmittedAt: sc.SubmittedAt,
		CompletedAt: sc.CompletedAt,
		UpdatedAt:   sc.UpdatedAt,
	}
}
