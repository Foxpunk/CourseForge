package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Foxpunk/courseforge/internal/interfaces"
)

// DisciplineHandler управляет CRUD дисциплин и назначением преподавателей
type DisciplineHandler struct {
	subjectManager interfaces.SubjectManager
}

// NewDisciplineHandler создаёт новый DisciplineHandler
func NewDisciplineHandler(sm interfaces.SubjectManager) *DisciplineHandler {
	return &DisciplineHandler{subjectManager: sm}
}

// CreateDiscipline - создание дисциплины (admin)
func (h *DisciplineHandler) CreateDiscipline(c *gin.Context) {
	var req interfaces.CreateSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subj, err := h.subjectManager.CreateSubject(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, interfaces.SubjectResponse{
		ID:          subj.ID,
		Name:        subj.Name,
		Code:        subj.Code,
		Description: subj.Description,
		Semester:    subj.Semester,
		IsActive:    subj.IsActive,
		Teachers:    []interfaces.UserResponse{}, // пустой массив вместо nil
		CreatedAt:   subj.CreatedAt.Format(time.RFC3339),
	})
}

// GetDisciplines - список всех дисциплин (teacher/admin)
func (h *DisciplineHandler) GetDisciplines(c *gin.Context) {
	list, err := h.subjectManager.ListSubjects(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make([]interfaces.SubjectResponse, len(list))
	for i, subj := range list {
		// Преобразуем учителей в UserResponse
		teachers := make([]interfaces.UserResponse, len(subj.Teachers))
		for j, teacher := range subj.Teachers {
			teachers[j] = interfaces.UserResponse{
				ID:        teacher.ID,
				Email:     teacher.Email,
				FirstName: teacher.FirstName,
				LastName:  teacher.LastName,
				Role:      teacher.Role,
				IsActive:  teacher.IsActive,
				CreatedAt: teacher.CreatedAt.Format(time.RFC3339),
			}
		}

		resp[i] = interfaces.SubjectResponse{
			ID:          subj.ID,
			Name:        subj.Name,
			Code:        subj.Code,
			Description: subj.Description,
			Semester:    subj.Semester,
			IsActive:    subj.IsActive,
			Teachers:    teachers,
			CreatedAt:   subj.CreatedAt.Format(time.RFC3339),
		}
	}

	c.JSON(http.StatusOK, resp)
}

// AssignTeacher - назначить преподавателя на дисциплину (admin)
func (h *DisciplineHandler) AssignTeacher(c *gin.Context) {
	// subjectID из URL
	idParam := c.Param("id")
	subjID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subject id"})
		return
	}

	var req interfaces.CreateTeacherSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// используем subjID из URL
	err = h.subjectManager.AssignTeacherToSubject(
		c.Request.Context(),
		req.TeacherID,
		uint(subjID),
		req.AcademicYear,
		req.IsLead,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetDiscipline - получить дисциплину по ID с учителями
func (h *DisciplineHandler) GetDiscipline(c *gin.Context) {
	idParam := c.Param("id")
	subjID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subject id"})
		return
	}

	subj, err := h.subjectManager.GetSubject(c.Request.Context(), uint(subjID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subject not found"})
		return
	}

	// Преобразуем учителей в UserResponse
	teachers := make([]interfaces.UserResponse, len(subj.Teachers))
	for j, teacher := range subj.Teachers {
		teachers[j] = interfaces.UserResponse{
			ID:        teacher.ID,
			Email:     teacher.Email,
			FirstName: teacher.FirstName,
			LastName:  teacher.LastName,
			Role:      teacher.Role,
			IsActive:  teacher.IsActive,
			CreatedAt: teacher.CreatedAt.Format(time.RFC3339),
		}
	}

	c.JSON(http.StatusOK, interfaces.SubjectResponse{
		ID:          subj.ID,
		Name:        subj.Name,
		Code:        subj.Code,
		Description: subj.Description,
		Semester:    subj.Semester,
		IsActive:    subj.IsActive,
		Teachers:    teachers,
		CreatedAt:   subj.CreatedAt.Format(time.RFC3339),
	})
}

// UpdateDiscipline - обновить дисциплину (admin)
func (h *DisciplineHandler) UpdateDiscipline(c *gin.Context) {
	idParam := c.Param("id")
	subjID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subject id"})
		return
	}

	var req interfaces.UpdateSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subj, err := h.subjectManager.UpdateSubject(c.Request.Context(), uint(subjID), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, interfaces.SubjectResponse{
		ID:          subj.ID,
		Name:        subj.Name,
		Code:        subj.Code,
		Description: subj.Description,
		Semester:    subj.Semester,
		IsActive:    subj.IsActive,
		Teachers:    []interfaces.UserResponse{}, // можно загрузить при необходимости
		CreatedAt:   subj.CreatedAt.Format(time.RFC3339),
	})
}

// DeleteDiscipline - удалить дисциплину (admin)
func (h *DisciplineHandler) DeleteDiscipline(c *gin.Context) {
	idParam := c.Param("id")
	subjID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subject id"})
		return
	}

	err = h.subjectManager.DeleteSubject(c.Request.Context(), uint(subjID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// RemoveTeacher - удалить преподавателя с дисциплины (admin)
func (h *DisciplineHandler) RemoveTeacher(c *gin.Context) {
	idParam := c.Param("id")
	subjID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subject id"})
		return
	}

	teacherIDParam := c.Param("teacherId")
	teacherID, err := strconv.ParseUint(teacherIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid teacher id"})
		return
	}

	// Академический год из query параметра
	academicYear := c.Query("academic_year")
	if academicYear == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "academic_year is required"})
		return
	}

	err = h.subjectManager.RemoveTeacherFromSubject(
		c.Request.Context(),
		uint(teacherID),
		uint(subjID),
		academicYear,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// SetLeadTeacher - назначить ведущего преподавателя (admin)
func (h *DisciplineHandler) SetLeadTeacher(c *gin.Context) {
	idParam := c.Param("id")
	subjID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subject id"})
		return
	}

	var req SetLeadTeacherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.subjectManager.SetLeadTeacher(
		c.Request.Context(),
		req.TeacherID,
		uint(subjID),
		req.AcademicYear,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

type SetLeadTeacherRequest struct {
	TeacherID    uint   `json:"teacher_id" validate:"required"`
	AcademicYear string `json:"academic_year" validate:"required"`
}
