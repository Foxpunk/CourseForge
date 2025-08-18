import React, { useState } from 'react';
import { CreateCourseworkRequest, DifficultyLevel } from '../../types';
import { useSubjects } from '../../hooks/useSubjects';
import { useAuth } from '../../hooks/useAuth';
import Button from '../ui/Button';
import ErrorMessage from '../ui/ErrorMessage';

interface CreateCourseworkFormProps {
 onSubmit: (data: CreateCourseworkRequest) => Promise<void>;
 onCancel: () => void;
}

const CreateCourseworkForm: React.FC<CreateCourseworkFormProps> = ({
 onSubmit,
 onCancel,
}) => {
 const { user } = useAuth();
 const [formData, setFormData] = useState<CreateCourseworkRequest>({
   title: '',
   description: '',
   requirements: '',
   subject_id: 0,
   teacher_id: user?.id || 0,
   max_students: 1,
   difficulty_level: 'medium' as DifficultyLevel,
 });
 
 const { subjects, loading: subjectsLoading } = useSubjects();
 const [loading, setLoading] = useState(false);
 const [error, setError] = useState<string | null>(null);

 const availableSubjects = subjects.filter(subject => 
   subject.teachers?.some(teacher => teacher.id === user?.id)
 );

 const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) => {
   const { name, value } = e.target;
   setFormData(prev => ({
     ...prev,
     [name]: name === 'subject_id' || name === 'teacher_id' 
       ? parseInt(value) || 0 
       : value,
   }));
 };

 const handleSubmit = async (e: React.FormEvent) => {
   e.preventDefault();
   
   if (formData.subject_id === 0) {
     setError('Выберите дисциплину');
     return;
   }

   if (formData.teacher_id === 0) {
     setError('Не удалось определить преподавателя');
     return;
   }

   if (formData.title.length < 5) {
     setError('Название должно содержать минимум 5 символов');
     return;
   }

   if (formData.description.length < 20) {
     setError('Описание должно содержать минимум 20 символов');
     return;
   }

   setLoading(true);
   setError(null);

   try {
     await onSubmit(formData);
   } catch (err: any) {
     setError(err.message || 'Ошибка при создании курсовой');
   } finally {
     setLoading(false);
   }
 };

 if (subjectsLoading) {
   return (
     <div className="bg-gray-800 p-6 rounded-lg">
       <div className="flex justify-center py-4">
         <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-orange-500"></div>
       </div>
     </div>
   );
 }

 if (availableSubjects.length === 0) {
   return (
     <div className="bg-gray-800 p-6 rounded-lg">
       <h2 className="text-2xl font-bold text-orange-500 mb-4">
         Создать новую курсовую работу
       </h2>
       <div className="text-center py-8">
         <p className="text-gray-400 mb-4">
           У вас нет назначенных дисциплин для создания курсовых работ.
         </p>
         <Button variant="secondary" onClick={onCancel}>
           Закрыть
         </Button>
       </div>
     </div>
   );
 }

 return (
   <div className="bg-gray-800 p-6 rounded-lg">
     <h2 className="text-2xl font-bold text-orange-500 mb-6">
       Создать новую курсовую работу
     </h2>

     {error && (
       <div className="mb-4">
         <ErrorMessage message={error} />
       </div>
     )}

     <form onSubmit={handleSubmit} className="space-y-4">
       <div>
         <label htmlFor="title" className="block text-sm font-medium text-gray-300 mb-1">
           Название работы <span className="text-red-400">*</span>
         </label>
         <input
           id="title"
           name="title"
           type="text"
           required
           value={formData.title}
           onChange={handleChange}
           className="w-full px-3 py-2 border border-gray-600 rounded-md bg-gray-700 text-white focus:outline-none focus:ring-2 focus:ring-orange-500"
           placeholder="Введите название курсовой работы"
         />
         <p className="text-xs text-gray-500 mt-1">
           Минимум 5 символов ({formData.title.length}/5)
         </p>
       </div>

       <div>
         <label htmlFor="subject_id" className="block text-sm font-medium text-gray-300 mb-1">
           Дисциплина <span className="text-red-400">*</span>
         </label>
         <select
           id="subject_id"
           name="subject_id"
           required
           value={formData.subject_id}
           onChange={handleChange}
           className="w-full px-3 py-2 border border-gray-600 rounded-md bg-gray-700 text-white focus:outline-none focus:ring-2 focus:ring-orange-500"
         >
           <option value={0}>Выберите дисциплину</option>
           {availableSubjects.map((subject) => (
             <option key={subject.id} value={subject.id}>
               {subject.name} ({subject.code}) - {subject.semester} семестр
             </option>
           ))}
         </select>
       </div>

       <div>
         <label htmlFor="difficulty_level" className="block text-sm font-medium text-gray-300 mb-1">
           Сложность
         </label>
         <select
           id="difficulty_level"
           name="difficulty_level"
           value={formData.difficulty_level}
           onChange={handleChange}
           className="w-full px-3 py-2 border border-gray-600 rounded-md bg-gray-700 text-white focus:outline-none focus:ring-2 focus:ring-orange-500"
         >
           <option value="easy">Легкая</option>
           <option value="medium">Средняя</option>
           <option value="hard">Сложная</option>
         </select>
       </div>

       <div>
         <label htmlFor="description" className="block text-sm font-medium text-gray-300 mb-1">
           Описание <span className="text-red-400">*</span>
         </label>
         <textarea
           id="description"
           name="description"
           rows={4}
           required
           value={formData.description}
           onChange={handleChange}
           className="w-full px-3 py-2 border border-gray-600 rounded-md bg-gray-700 text-white focus:outline-none focus:ring-2 focus:ring-orange-500"
           placeholder="Опишите задачи и цели курсовой работы"
         />
         <p className={`text-xs mt-1 ${formData.description.length >= 20 ? 'text-green-400' : 'text-red-400'}`}>
           Минимум 20 символов ({formData.description.length}/20)
         </p>
       </div>

       <div>
         <label htmlFor="requirements" className="block text-sm font-medium text-gray-300 mb-1">
           Требования (необязательно)
         </label>
         <textarea
           id="requirements"
           name="requirements"
           rows={3}
           value={formData.requirements}
           onChange={handleChange}
           className="w-full px-3 py-2 border border-gray-600 rounded-md bg-gray-700 text-white focus:outline-none focus:ring-2 focus:ring-orange-500"
           placeholder="Дополнительные требования к выполнению"
         />
       </div>

       <div className="flex gap-4 pt-4">
         <Button
           type="submit"
           loading={loading}
           disabled={formData.subject_id === 0 || formData.title.length < 5 || formData.description.length < 20}
           className="flex-1"
         >
           Создать курсовую
         </Button>
         <Button
           type="button"
           variant="secondary"
           onClick={onCancel}
           className="flex-1"
         >
           Отмена
         </Button>
       </div>
     </form>
   </div>
 );
};

export default CreateCourseworkForm;