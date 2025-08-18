import { useState, useEffect } from 'react';
import { CourseworkResponse, CreateCourseworkRequest, UpdateCourseworkRequest } from '../types';
import { courseworksApi } from '../api/courseworks';
import { useAuth } from './useAuth';

export const useCourseworks = () => {
 const [courseworks, setCourseworks] = useState<CourseworkResponse[]>([]);
 const [loading, setLoading] = useState(true);
 const [error, setError] = useState<string | null>(null);
 const [userHasAssignment, setUserHasAssignment] = useState(false); 
 const [userAssignedCourseworkId, setUserAssignedCourseworkId] = useState<number | null>(null); 
 const { user } = useAuth();

 const fetchCourseworks = async () => {
   try {
     setLoading(true);
     setError(null);
     
     let data: CourseworkResponse[] = [];
     
     if (user?.role === 'student') {
       data = await courseworksApi.getAvailableProjects();
       
       // Проверить есть ли у студента назначенная курсовая
       try {
         const studentCoursework = await courseworksApi.getStudentCoursework();
         if (studentCoursework) {
           setUserHasAssignment(true);
           setUserAssignedCourseworkId(studentCoursework.coursework.id);
         } else {
           setUserHasAssignment(false);
           setUserAssignedCourseworkId(null);
         }
       } catch (err) {
         // Нет назначенной курсовой - это нормально
         setUserHasAssignment(false);
         setUserAssignedCourseworkId(null);
       }
     } else if (user?.role === 'teacher' || user?.role === 'admin') {
       const response = await courseworksApi.getProjects();
       data = response.courseworks;
     }
     
     setCourseworks(data);
   } catch (err: any) {
     const errorMessage = err.response?.data?.error || err.message || 'Произошла ошибка при загрузке курсовых';
     setError(errorMessage);
   } finally {
     setLoading(false);
   }
 };

   useEffect(() => {
   if (user) {
     fetchCourseworks();
   }
 }, [user]);

 const assignStudent = async (courseworkId: number) => {
   if (!user?.id) {
     throw new Error('Пользователь не авторизован');
   }
   
   try {
     await courseworksApi.assignStudent(courseworkId, user.id);
     setUserHasAssignment(true);
     setUserAssignedCourseworkId(courseworkId);
     await fetchCourseworks();
   } catch (err: any) {
     const errorMessage = err.response?.data?.error || err.message || 'Ошибка при выборе курсовой';
     setError(errorMessage);
     throw new Error(errorMessage);
   }
 };

  const createCoursework = async (data: CreateCourseworkRequest) => {
    try {
      await courseworksApi.createProject(data);
      await fetchCourseworks(); // Обновляем список после создания
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || err.message || 'Ошибка при создании курсовой';
      setError(errorMessage);
      throw new Error(errorMessage);
    }
  };

  const updateCoursework = async (id: number, data: UpdateCourseworkRequest) => {
    try {
      await courseworksApi.updateProject(id, data);
      await fetchCourseworks(); // Обновляем список после обновления
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || err.message || 'Ошибка при обновлении курсовой';
      setError(errorMessage);
      throw new Error(errorMessage);
    }
  };

  const deleteCoursework = async (id: number) => {
    try {
      await courseworksApi.deleteProject(id);
      await fetchCourseworks(); // Обновляем список после удаления
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || err.message || 'Ошибка при удалении курсовой';
      setError(errorMessage);
      throw new Error(errorMessage);
    }
  };

return {
   courseworks,
   loading,
   error,
   userHasAssignment, 
   userAssignedCourseworkId, 
   assignStudent,
   createCoursework,
   updateCoursework,
   deleteCoursework,
   refetch: fetchCourseworks,
 };
};