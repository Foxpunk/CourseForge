import { apiClient } from './client';
import { DepartmentResponse, StudentGroupResponse } from '../types';

export const departmentsApi = {
 getDepartments: async (): Promise<DepartmentResponse[]> => {
   const response = await apiClient.get<DepartmentResponse[]>('/departments');
   return response.data;
 },

 getGroupsByDepartment: async (departmentId: number): Promise<StudentGroupResponse[]> => {
   const response = await apiClient.get<StudentGroupResponse[]>(`/departments/${departmentId}/groups`);
   return response.data;
 },
};