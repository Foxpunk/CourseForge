// Базовые типы, соответствующие Go models
export type UserRole = 'admin' | 'teacher' | 'student';
export type DifficultyLevel = 'easy' | 'medium' | 'hard';
export type CourseworkStatus = 'assigned' | 'in_progress' | 'submitted' | 'reviewed' | 'completed' | 'failed';

// User response (соответствует UserResponse в Go)
export interface UserResponse {
  id: number;
  email: string;
  first_name: string;
  last_name: string;
  role: UserRole;
  is_active: boolean;
  created_at: string;
}

// Department response
export interface DepartmentResponse {
  id: number;
  department_code: string;
  department_name: string;
  description: string;
  created_at: string;
}

// Student group response
export interface StudentGroupResponse {
  id: number;
  group_code: string;
  course_year: number;
  specialty: string;
  department: DepartmentResponse;
  created_at: string;
}

// Subject response
export interface SubjectResponse {
  id: number;
  name: string;
  code: string;
  description: string;
  semester: number;
  is_active: boolean;
  teachers?: UserResponse[];
  created_at: string;
}

// Coursework response (соответствует CourseworkResponse в Go)
export interface CourseworkResponse {
  id: number;
  title: string;
  description: string;
  requirements: string;
  max_students: number;
  difficulty_level: DifficultyLevel;
  is_available: boolean;
  subject: SubjectResponse;
  teacher: UserResponse;
  created_at: string;
  updated_at: string;
}

// Student coursework response
export interface StudentCourseworkResponse {
  id: number;
  student: UserResponse;
  coursework: CourseworkResponse;
  status: CourseworkStatus;
  grade?: number;
  feedback?: string;
  assigned_at: string;
  submitted_at?: string;
  completed_at?: string;
  updated_at: string;
}

// AUTH DTOs (соответствуют Go DTO)
export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: UserResponse;
  expires_at: number;
}

export interface RegisterRequest {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
  role: UserRole;
}

export interface ChangePasswordRequest {
  old_password: string;
  new_password: string;
}

// COURSEWORK DTOs
export interface CreateCourseworkRequest {
  title: string;
  description: string;
  requirements?: string;
  subject_id: number;
  teacher_id: number;
  max_students: number;
  difficulty_level: DifficultyLevel;
}

export interface UpdateCourseworkRequest {
  title?: string;
  description?: string;
  requirements?: string;
  subject_id?: number;
  teacher_id?: number;
  max_students?: number;
  difficulty_level?: DifficultyLevel;
  is_available?: boolean;
}

export interface CourseworkListResponse {
  courseworks: CourseworkResponse[];
  total: number;
}

export interface ListCourseworksRequest {
  subject_id?: number;
  teacher_id?: number;
  available?: boolean;
  difficulty_level?: DifficultyLevel;
  limit?: number;
  offset?: number;
}

// Error response
export interface ErrorResponse {
  error: string;
  message?: string;
}