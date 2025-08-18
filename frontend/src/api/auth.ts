import { apiClient } from './client';
import { 
  LoginRequest, 
  LoginResponse, 
  RegisterRequest, 
  ChangePasswordRequest,
  UserResponse 
} from '../types';

export const authApi = {
  login: async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await apiClient.post<LoginResponse>('/auth/login', data);
    return response.data;
  },

  register: async (data: RegisterRequest): Promise<LoginResponse> => {
    const response = await apiClient.post<LoginResponse>('/auth/register', data);
    return response.data;
  },

  logout: async (): Promise<void> => {
    await apiClient.post('/profile/logout');
  },

  refreshToken: async (): Promise<{ token: string; expires_at: number }> => {
    const response = await apiClient.post('/auth/refresh');
    return response.data;
  },

  getProfile: async (): Promise<UserResponse> => {
    const response = await apiClient.get<UserResponse>('/profile');
    return response.data;
  },

  changePassword: async (data: ChangePasswordRequest): Promise<void> => {
    await apiClient.post('/profile/change-password', data);
  },

  resetPassword: async (email: string): Promise<void> => {
    await apiClient.post('/auth/reset-password', { email });
  },
};