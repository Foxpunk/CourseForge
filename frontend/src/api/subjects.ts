import { apiClient } from './client';
import { SubjectResponse } from '../types';

export const subjectsApi = {
  // Получить все дисциплины
  getSubjects: async (): Promise<SubjectResponse[]> => {
    const response = await apiClient.get<SubjectResponse[]>('/subjects');
    return response.data;
  },

  // Получить дисциплину по ID
  getSubject: async (id: number): Promise<SubjectResponse> => {
    const response = await apiClient.get<SubjectResponse>(`/subjects/${id}`);
    return response.data;
  },

  // Создать дисциплину (админ)
  createSubject: async (data: {
    name: string;
    code: string;
    description?: string;
    semester: number;
  }): Promise<SubjectResponse> => {
    const response = await apiClient.post<SubjectResponse>('/subjects', data);
    return response.data;
  },

  // Обновить дисциплину
  updateSubject: async (id: number, data: {
    name?: string;
    description?: string;
    semester?: number;
    is_active?: boolean;
  }): Promise<SubjectResponse> => {
    const response = await apiClient.put<SubjectResponse>(`/subjects/${id}`, data);
    return response.data;
  },

  // Удалить дисциплину
  deleteSubject: async (id: number): Promise<void> => {
    await apiClient.delete(`/subjects/${id}`);
  },

  // Назначить преподавателя на дисциплину
  assignTeacher: async (subjectId: number, data: {
    teacher_id: number;
    academic_year: string;
    is_lead?: boolean;
  }): Promise<void> => {
    await apiClient.post(`/subjects/${subjectId}/teachers`, data);
  },

  // Убрать преподавателя с дисциплины
  removeTeacher: async (subjectId: number, teacherId: number): Promise<void> => {
    await apiClient.delete(`/subjects/${subjectId}/teachers/${teacherId}`);
  },

  // Назначить ведущего преподавателя
  setLeadTeacher: async (subjectId: number, teacherId: number): Promise<void> => {
    await apiClient.put(`/subjects/${subjectId}/lead-teacher`, {
      teacher_id: teacherId,
    });
  },
};