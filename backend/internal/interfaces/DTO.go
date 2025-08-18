package interfaces

import (
	"time"

	"github.com/Foxpunk/courseforge/internal/models"
)

// ============================================================================
// AUTH DTOs
// ============================================================================

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginResponse struct {
	Token     string      `json:"token"`
	User      models.User `json:"user"`
	ExpiresAt int64       `json:"expires_at"`
}

type RegisterRequest struct {
	Email     string          `json:"email" validate:"required,email"`
	Password  string          `json:"password" validate:"required,min=6"`
	FirstName string          `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string          `json:"last_name" validate:"required,min=2,max=50"`
	Role      models.UserRole `json:"role" validate:"required,oneof=admin teacher student"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

type RefreshTokenRequest struct {
	Token string `json:"token" validate:"required"`
}

type RefreshTokenResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// Форма нового пароля (после клика по ссылке из письма)
type ConfirmResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

type LogoutRequest struct {
	Token string `json:"token" validate:"required"`
}

// ============================================================================
// USER DTOs
// ============================================================================

type CreateUserRequest struct {
	Email     string          `json:"email" validate:"required,email"`
	Password  string          `json:"password" validate:"required,min=6"`
	FirstName string          `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string          `json:"last_name" validate:"required,min=2,max=50"`
	Role      models.UserRole `json:"role" validate:"required,oneof=admin teacher student"`
}

type UpdateUserRequest struct {
	Email     *string          `json:"email,omitempty" validate:"omitempty,email"`
	FirstName *string          `json:"first_name,omitempty" validate:"omitempty,min=2,max=50"`
	LastName  *string          `json:"last_name,omitempty" validate:"omitempty,min=2,max=50"`
	Role      *models.UserRole `json:"role,omitempty" validate:"omitempty,oneof=admin teacher student"`
	IsActive  *bool            `json:"is_active,omitempty"`
}

type ListUsersRequest struct {
	Role   *models.UserRole `json:"role,omitempty"`
	Active *bool            `json:"active,omitempty"`
	Limit  int              `json:"limit" validate:"min=1,max=100"`
	Offset int              `json:"offset" validate:"min=0"`
}

type UserResponse struct {
	ID        uint            `json:"id"`
	Email     string          `json:"email"`
	FirstName string          `json:"first_name"`
	LastName  string          `json:"last_name"`
	Role      models.UserRole `json:"role"`
	IsActive  bool            `json:"is_active"`
	CreatedAt string          `json:"created_at"`
}

type UserListResponse struct {
	Users []UserResponse `json:"users"`
	Total int64          `json:"total"`
}

// ============================================================================
// DEPARTMENT DTOs
// ============================================================================

type CreateDepartmentRequest struct {
	DepartmentCode string `json:"department_code" validate:"required,min=2,max=10"`
	DepartmentName string `json:"department_name" validate:"required,min=3,max=200"`
	Description    string `json:"description"`
}

type UpdateDepartmentRequest struct {
	DepartmentCode *string `json:"department_code,omitempty" validate:"omitempty,min=2,max=10"`
	DepartmentName *string `json:"department_name,omitempty" validate:"omitempty,min=3,max=200"`
	Description    *string `json:"description,omitempty"`
}

type DepartmentResponse struct {
	ID             uint   `json:"id"`
	DepartmentCode string `json:"department_code"`
	DepartmentName string `json:"department_name"`
	Description    string `json:"description"`
	CreatedAt      string `json:"created_at"`
}

// ============================================================================
// STUDENT GROUP DTOs
// ============================================================================

type CreateGroupRequest struct {
	GroupCode    string `json:"group_code" validate:"required,min=2,max=20"`
	CourseYear   int    `json:"course_year" validate:"required,min=1,max=6"`
	Specialty    string `json:"specialty" validate:"required,min=3,max=100"`
	DepartmentID uint   `json:"department_id" validate:"required"`
}

type UpdateGroupRequest struct {
	GroupCode    *string `json:"group_code,omitempty" validate:"omitempty,min=2,max=20"`
	CourseYear   *int    `json:"course_year,omitempty" validate:"omitempty,min=1,max=6"`
	Specialty    *string `json:"specialty,omitempty" validate:"omitempty,min=3,max=100"`
	DepartmentID *uint   `json:"department_id,omitempty"`
}

type StudentGroupResponse struct {
	ID         uint               `json:"id"`
	GroupCode  string             `json:"group_code"`
	CourseYear int                `json:"course_year"`
	Specialty  string             `json:"specialty"`
	Department DepartmentResponse `json:"department"`
	CreatedAt  string             `json:"created_at"`
}

// ============================================================================
// STUDENT PROFILE DTOs
// ============================================================================

type CreateStudentProfileRequest struct {
	UserID        uint   `json:"user_id" validate:"required"`
	GroupID       uint   `json:"group_id" validate:"required"`
	StudentNumber string `json:"student_number" validate:"max=20"`
}

type UpdateStudentProfileRequest struct {
	GroupID       *uint   `json:"group_id,omitempty"`
	StudentNumber *string `json:"student_number,omitempty" validate:"omitempty,max=20"`
}

type StudentProfileResponse struct {
	ID            uint                 `json:"id"`
	User          UserResponse         `json:"user"`
	StudentGroup  StudentGroupResponse `json:"student_group"`
	StudentNumber string               `json:"student_number"`
	CreatedAt     string               `json:"created_at"`
}

// ============================================================================
// TEACHER PROFILE DTOs
// ============================================================================

type CreateTeacherProfileRequest struct {
	UserID         uint   `json:"user_id" validate:"required"`
	DepartmentID   uint   `json:"department_id" validate:"required"`
	Position       string `json:"position" validate:"max=100"`
	AcademicDegree string `json:"academic_degree" validate:"max=100"`
}

type UpdateTeacherProfileRequest struct {
	DepartmentID   *uint   `json:"department_id,omitempty"`
	Position       *string `json:"position,omitempty" validate:"omitempty,max=100"`
	AcademicDegree *string `json:"academic_degree,omitempty" validate:"omitempty,max=100"`
}

type TeacherProfileResponse struct {
	ID             uint               `json:"id"`
	User           UserResponse       `json:"user"`
	Department     DepartmentResponse `json:"department"`
	Position       string             `json:"position"`
	AcademicDegree string             `json:"academic_degree"`
	CreatedAt      string             `json:"created_at"`
}

// ============================================================================
// SUBJECT DTOs
// ============================================================================

type CreateSubjectRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Code        string `json:"code" validate:"required,min=2,max=20"`
	Description string `json:"description"`
	Semester    int    `json:"semester" validate:"required,min=1,max=12"`
}

type UpdateSubjectRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	Description *string `json:"description,omitempty"`
	Semester    *int    `json:"semester,omitempty" validate:"omitempty,min=1,max=12"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

type SubjectResponse struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	Code        string         `json:"code"`
	Description string         `json:"description"`
	Semester    int            `json:"semester"`
	IsActive    bool           `json:"is_active"`
	Teachers    []UserResponse `json:"teachers,omitempty"`
	CreatedAt   string         `json:"created_at"`
}

// ============================================================================
// TEACHER SUBJECT DTOs
// ============================================================================

type CreateTeacherSubjectRequest struct {
	TeacherID    uint   `json:"teacher_id" validate:"required"`
	SubjectID    uint   `json:"subject_id" validate:"required"`
	AcademicYear string `json:"academic_year" validate:"required,min=9,max=9"`
	IsLead       bool   `json:"is_lead"`
}

type UpdateTeacherSubjectRequest struct {
	IsLead *bool `json:"is_lead,omitempty"`
}

type TeacherSubjectResponse struct {
	ID           uint            `json:"id"`
	Teacher      UserResponse    `json:"teacher"`
	Subject      SubjectResponse `json:"subject"`
	AcademicYear string          `json:"academic_year"`
	IsLead       bool            `json:"is_lead"`
	CreatedAt    string          `json:"created_at"`
}

type AssignTeacherToSubjectRequest struct {
	TeacherIDs   []uint `json:"teacher_ids" validate:"required,min=1"`
	AcademicYear string `json:"academic_year" validate:"required,min=9,max=9"`
}

// ============================================================================
// COURSEWORK DTOs
// ============================================================================

type CreateCourseworkRequest struct {
	Title           string                 `json:"title" validate:"required,min=5,max=300"`
	Description     string                 `json:"description" validate:"required,min=20"`
	Requirements    string                 `json:"requirements"`
	SubjectID       uint                   `json:"subject_id" validate:"required"`
	TeacherID       uint                   `json:"teacher_id" validate:"required"`
	MaxStudents     int                    `json:"max_students" validate:"min=1,max=10"`
	DifficultyLevel models.DifficultyLevel `json:"difficulty_level" validate:"required,oneof=easy medium hard"`
}

type UpdateCourseworkRequest struct {
	Title           *string                 `json:"title,omitempty" validate:"omitempty,min=5,max=300"`
	Description     *string                 `json:"description,omitempty" validate:"omitempty,min=20"`
	Requirements    *string                 `json:"requirements,omitempty"`
	SubjectID       *uint                   `json:"subject_id,omitempty" validate:"omitempty"`
	TeacherID       *uint                   `json:"teacher_id,omitempty" validate:"omitempty"`
	MaxStudents     *int                    `json:"max_students,omitempty" validate:"omitempty,min=1,max=10"`
	DifficultyLevel *models.DifficultyLevel `json:"difficulty_level,omitempty" validate:"omitempty,oneof=easy medium hard"`
	IsAvailable     *bool                   `json:"is_available,omitempty"`
}

type ListCourseworksRequest struct {
	SubjectID  *uint                   `json:"subject_id,omitempty"`
	TeacherID  *uint                   `json:"teacher_id,omitempty"`
	Available  *bool                   `json:"available,omitempty"`
	Difficulty *models.DifficultyLevel `json:"difficulty_level,omitempty"`
	Limit      int                     `json:"limit" validate:"min=1,max=100"`
	Offset     int                     `json:"offset" validate:"min=0"`
}

type CourseworkResponse struct {
	ID              uint                   `json:"id"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	Requirements    string                 `json:"requirements"`
	MaxStudents     int                    `json:"max_students"`
	DifficultyLevel models.DifficultyLevel `json:"difficulty_level"`
	IsAvailable     bool                   `json:"is_available"`
	Subject         SubjectResponse        `json:"subject"`
	Teacher         UserResponse           `json:"teacher"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

type CourseworkListResponse struct {
	Courseworks []CourseworkResponse `json:"courseworks"`
	Total       int64                `json:"total"`
}

// ============================================================================
// STUDENT COURSEWORK DTOs
// ============================================================================

type CreateStudentCourseworkRequest struct {
	StudentID    uint `json:"student_id" validate:"required"`
	CourseworkID uint `json:"coursework_id" validate:"required"`
}

type UpdateStudentCourseworkRequest struct {
	Status   models.CourseworkStatus `json:"status" validate:"required,oneof=assigned in_progress submitted reviewed completed failed"`
	Grade    *int                    `json:"grade,omitempty" validate:"omitempty,min=2,max=5"`
	Feedback *string                 `json:"feedback,omitempty"`
}

type StudentCourseworkResponse struct {
	ID          uint                    `json:"id"`
	Student     UserResponse            `json:"student"`
	Coursework  CourseworkResponse      `json:"coursework"`
	Status      models.CourseworkStatus `json:"status"`
	Grade       *int                    `json:"grade,omitempty"`
	Feedback    *string                 `json:"feedback,omitempty"`
	AssignedAt  time.Time               `json:"assigned_at"`
	SubmittedAt *time.Time              `json:"submitted_at,omitempty"`
	CompletedAt *time.Time              `json:"completed_at,omitempty"`
	UpdatedAt   time.Time               `json:"updated_at"`
}

type AssignStudentToCourseworkRequest struct {
	StudentID    uint `json:"student_id" validate:"required"`
	CourseworkID uint `json:"coursework_id" validate:"required"`
}

type SubmitCourseworkRequest struct {
	Status models.CourseworkStatus `json:"status" validate:"required,oneof=submitted"`
}

type GradeCourseworkRequest struct {
	Grade    *int                    `json:"grade" validate:"required,min=2,max=5"`
	Feedback *string                 `json:"feedback,omitempty"`
	Status   models.CourseworkStatus `json:"status" validate:"required,oneof=reviewed completed failed"`
}

// ============================================================================
// PROGRESS REPORTS
// ============================================================================

type CourseworkProgressReport struct {
	CourseworkID     uint                       `json:"coursework_id"`
	CourseworkTitle  string                     `json:"coursework_title"`
	TotalStudents    int                        `json:"total_students"`
	AssignedStudents int                        `json:"assigned_students"`
	InProgressCount  int                        `json:"in_progress_count"`
	SubmittedCount   int                        `json:"submitted_count"`
	CompletedCount   int                        `json:"completed_count"`
	FailedCount      int                        `json:"failed_count"`
	Students         []StudentCourseworkSummary `json:"students"`
}

type StudentCourseworkSummary struct {
	StudentID   uint                    `json:"student_id"`
	StudentName string                  `json:"student_name"`
	Status      models.CourseworkStatus `json:"status"`
	AssignedAt  time.Time               `json:"assigned_at"`
	SubmittedAt *time.Time              `json:"submitted_at,omitempty"`
	CompletedAt *time.Time              `json:"completed_at,omitempty"`
	Grade       *int                    `json:"grade,omitempty"`
}
