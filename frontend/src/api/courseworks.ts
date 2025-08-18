import { apiClient } from './client';
import { 
  CourseworkResponse, 
  CourseworkListResponse,
  CreateCourseworkRequest, 
  UpdateCourseworkRequest,
  ListCourseworksRequest,
  StudentCourseworkResponse 
} from '../types';

export const courseworksApi = {
  // Получить все проекты (в зависимости от роли пользователя)
  getProjects: async (params?: ListCourseworksRequest): Promise<CourseworkListResponse> => {
    const response = await apiClient.get<CourseworkListResponse>('/courseworks', { params });
    return response.data;
  },

  // Получить доступные курсовые для студента
  getAvailableProjects: async (): Promise<CourseworkResponse[]> => {
    const response = await apiClient.get<CourseworkResponse[]>('/courseworks/available');
    return response.data;
  },

  // Получить курсовую по ID
  getProject: async (id: number): Promise<CourseworkResponse> => {
    const response = await apiClient.get<CourseworkResponse>(`/courseworks/${id}`);
    return response.data;
  },

  // Создать курсовую (преподаватель/админ)
  createProject: async (data: CreateCourseworkRequest): Promise<CourseworkResponse> => {
    const response = await apiClient.post<CourseworkResponse>('/courseworks', data);
    return response.data;
  },

  // Обновить курсовую
  updateProject: async (id: number, data: UpdateCourseworkRequest): Promise<CourseworkResponse> => {
    const response = await apiClient.put<CourseworkResponse>(`/courseworks/${id}`, data);
    return response.data;
  },

  // Удалить курсовую
  deleteProject: async (id: number): Promise<void> => {
    await apiClient.delete(`/courseworks/${id}`);
  },

  // Назначить студента на курсовую
  assignStudent: async (courseworkId: number, studentId: number): Promise<StudentCourseworkResponse> => {
    const response = await apiClient.post<StudentCourseworkResponse>(`/courseworks/${courseworkId}/assign`, {
      student_id: studentId,
      coursework_id: courseworkId 
    });
    return response.data;
  },
    getStudentCoursework: async (): Promise<StudentCourseworkResponse | null> => {
    try {
      const response = await apiClient.get<StudentCourseworkResponse>('/student-courseworks');
      return response.data;
    } catch (error: any) {
      if (error.response?.status === 404) {
        return null; 
      }
      throw error;
    }
  },

  // Изменить доступность курсовой
  setAvailability: async (id: number, isAvailable: boolean): Promise<void> => {
    await apiClient.put(`/courseworks/${id}/availability`, {
      is_available: isAvailable,
    });
  },
};