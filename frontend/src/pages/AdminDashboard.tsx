import React, { useState, useEffect } from 'react';
import { subjectsApi } from '../api/subjects';
import { usersApi } from '../api/users';
import { SubjectResponse, UserResponse } from '../types';
import LoadingSpinner from '../components/ui/LoadingSpinner';
import ErrorMessage from '../components/ui/ErrorMessage';
import Button from '../components/ui/Button';
import AdminSubjects from '../components/admin/AdminSubjects';

interface DepartmentStats {
 name: string;
 teachers_count: number;
 students_count: number;
 courseworks_count: number;
}

type AdminView = 'dashboard' | 'subjects';

const AdminDashboard: React.FC = () => {
 const [currentView, setCurrentView] = useState<AdminView>('dashboard');
 const [stats, setStats] = useState<DepartmentStats[]>([]);
 const [totalStats, setTotalStats] = useState({
   total_users: 0,
   total_teachers: 0,
   total_students: 0,
   total_subjects: 0,
 });
 const [loading, setLoading] = useState(true);
 const [error, setError] = useState<string | null>(null);

 useEffect(() => {
   if (currentView === 'dashboard') {
     fetchStats();
   }
 }, [currentView]);

 const fetchStats = async () => {
   try {
     setLoading(true);
     setError(null);

     // Получаем пользователей и дисциплины
     const [usersResponse, subjects] = await Promise.all([
       usersApi.getUsers({ limit: 1000, offset: 0 }),
       subjectsApi.getSubjects(),
     ]);

     const users = usersResponse.users;
     
     // Подсчитываем общую статистику
     const teachers = users.filter(u => u.role === 'teacher');
     const students = users.filter(u => u.role === 'student');

     setTotalStats({
       total_users: users.length,
       total_teachers: teachers.length,
       total_students: students.length,
       total_subjects: subjects.length,
     });

     // Группируем данные по дисциплинам (имитация статистики по кафедрам)
     const departmentMap = new Map<string, DepartmentStats>();

     subjects.forEach(subject => {
       const deptName = `Кафедра ${subject.name.split(' ')[0]}`;
       
       if (!departmentMap.has(deptName)) {
         departmentMap.set(deptName, {
           name: deptName,
           teachers_count: 0,
           students_count: 0,
           courseworks_count: 0,
         });
       }

       const dept = departmentMap.get(deptName)!;
       dept.teachers_count += subject.teachers?.length || 0;
       // Mock данные для курсовых и студентов
       dept.courseworks_count += Math.floor(Math.random() * 5) + 1;
       dept.students_count += Math.floor(Math.random() * 20) + 5;
     });

     setStats(Array.from(departmentMap.values()));
   } catch (err: any) {
     const errorMessage = err.response?.data?.error || err.message || 'Ошибка при загрузке статистики';
     setError(errorMessage);
   } finally {
     setLoading(false);
   }
 };

 // Переключение на управление дисциплинами
 if (currentView === 'subjects') {
   return (
     <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
       <AdminSubjects onBack={() => setCurrentView('dashboard')} />
     </div>
   );
 }

 if (loading) {
   return (
     <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
       <LoadingSpinner />
     </div>
   );
 }

 if (error) {
   return (
     <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
       <ErrorMessage message={error} />
     </div>
   );
 }

 return (
   <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
     <div className="mb-8">
       <h1 className="text-3xl font-bold text-orange-500 mb-2">
         Панель администратора
       </h1>
       <p className="text-gray-400">
         Общая статистика и управление системой
       </p>
     </div>

     {/* Общая статистика */}
     <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
       <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
         <div className="flex items-center">
           <div className="flex-shrink-0">
             <svg className="h-8 w-8 text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
               <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-8.5a4 4 0 11-8 0 4 4 0 018 0z" />
             </svg>
           </div>
           <div className="ml-4">
             <div className="text-2xl font-bold text-white">{totalStats.total_users}</div>
             <div className="text-sm text-gray-400">Всего пользователей</div>
           </div>
         </div>
       </div>

       <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
         <div className="flex items-center">
           <div className="flex-shrink-0">
             <svg className="h-8 w-8 text-green-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
               <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
             </svg>
           </div>
           <div className="ml-4">
             <div className="text-2xl font-bold text-white">{totalStats.total_teachers}</div>
             <div className="text-sm text-gray-400">Преподавателей</div>
           </div>
         </div>
       </div>

       <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
         <div className="flex items-center">
           <div className="flex-shrink-0">
             <svg className="h-8 w-8 text-yellow-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
               <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 14l9-5-9-5-9 5 9 5z" />
               <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 14l6.16-3.422a12.083 12.083 0 01.665 6.479A11.952 11.952 0 0012 20.055a11.952 11.952 0 00-6.824-2.998 12.078 12.078 0 01.665-6.479L12 14z" />
             </svg>
           </div>
           <div className="ml-4">
             <div className="text-2xl font-bold text-white">{totalStats.total_students}</div>
             <div className="text-sm text-gray-400">Студентов</div>
           </div>
         </div>
       </div>

       <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
         <div className="flex items-center">
           <div className="flex-shrink-0">
             <svg className="h-8 w-8 text-orange-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
               <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.746 0 3.332.477 4.5 1.253v13C19.832 18.477 18.246 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
             </svg>
           </div>
           <div className="ml-4">
             <div className="text-2xl font-bold text-white">{totalStats.total_subjects}</div>
             <div className="text-sm text-gray-400">Дисциплин</div>
           </div>
         </div>
       </div>
     </div>

     {/* Управление */}
     <div className="mb-8">
       <h2 className="text-xl font-semibold text-white mb-4">Управление системой</h2>
       <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
         <Button variant="primary" className="w-full">
           Управление пользователями
         </Button>
         <Button 
           variant="secondary" 
           className="w-full"
           onClick={() => setCurrentView('subjects')}
         >
           Управление дисциплинами
         </Button>
         <Button variant="secondary" className="w-full">
           Генерация отчетов
         </Button>
       </div>
     </div>

     {/* Статистика по кафедрам */}
     <div>
       <h2 className="text-xl font-semibold text-white mb-6">Статистика по направлениям</h2>
       
       {stats.length === 0 ? (
         <div className="text-center py-8">
           <p className="text-gray-400">Нет данных для отображения</p>
         </div>
       ) : (
         <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
           {stats.map((stat, index) => (
             <div
               key={index}
               className="bg-gray-800 rounded-lg p-6 border border-gray-700 hover:border-orange-500 transition-colors"
             >
               <h3 className="text-xl font-semibold text-orange-500 mb-4">
                 {stat.name}
               </h3>
               
               <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
                 <div className="text-center">
                   <div className="text-2xl font-bold text-orange-500">
                     {stat.courseworks_count}
                   </div>
                   <div className="text-sm text-gray-400">Курсовых</div>
                 </div>
                 <div className="text-center">
                   <div className="text-2xl font-bold text-orange-500">
                     {stat.teachers_count}
                   </div>
                   <div className="text-sm text-gray-400">Преподов</div>
                 </div>
                 <div className="text-center">
                   <div className="text-2xl font-bold text-orange-500">
                     {stat.students_count}
                   </div>
                   <div className="text-sm text-gray-400">Студентов</div>
                 </div>
               </div>

               <div className="flex gap-2">
                 <button className="flex-1 px-3 py-2 bg-orange-600 hover:bg-orange-700 text-white text-sm rounded-md transition-colors">
                   Управление
                 </button>
                 <button className="flex-1 px-3 py-2 bg-gray-600 hover:bg-gray-700 text-white text-sm rounded-md transition-colors">
                   Отчеты
                 </button>
               </div>
             </div>
           ))}
         </div>
       )}
     </div>
   </div>
 );
};

export default AdminDashboard;