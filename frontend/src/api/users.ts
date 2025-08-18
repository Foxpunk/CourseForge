import { apiClient } from './client';
import { UserResponse } from '../types';

export const usersApi = {
  // Получить всех пользователей (админ)
  getUsers: async (params?: {
    role?: string;
    active?: boolean;
    limit?: number;
    offset?: number;
  }): Promise<{ users: UserResponse[]; total: number }> => {
    const response = await apiClient.get('/users', { params });
    return response.data;
  },

  // Получить пользователя по ID
  getUser: async (id: number): Promise<UserResponse> => {
    const response = await apiClient.get<UserResponse>(`/users/${id}`);
    return response.data;
  },

  // Создать пользователя (админ)
  createUser: async (data: {
    email: string;
    password: string;
    first_name: string;
    last_name: string;
    role: string;
  }): Promise<UserResponse> => {
    const response = await apiClient.post<UserResponse>('/users', data);
    return response.data;
  },

  // Обновить пользователя
  updateUser: async (id: number, data: {
    email?: string;
    first_name?: string;
    last_name?: string;
    role?: string;
    is_active?: boolean;
  }): Promise<UserResponse> => {
    const response = await apiClient.put<UserResponse>(`/users/${id}`, data);
    return response.data;
  },

  // Удалить пользователя
  deleteUser: async (id: number): Promise<void> => {
    await apiClient.delete(`/users/${id}`);
  },
};