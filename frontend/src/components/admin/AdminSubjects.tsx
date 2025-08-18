import React, { useState, useEffect } from 'react';
import { subjectsApi } from '../../api/subjects';
import { usersApi } from '../../api/users';
import { SubjectResponse, UserResponse } from '../../types';
import Button from '../ui/Button';
import LoadingSpinner from '../ui/LoadingSpinner';
import ErrorMessage from '../ui/ErrorMessage';

interface AdminSubjectsProps {
  onBack: () => void;
}

const AdminSubjects: React.FC<AdminSubjectsProps> = ({ onBack }) => {
  const [subjects, setSubjects] = useState<SubjectResponse[]>([]);
  const [teachers, setTeachers] = useState<UserResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showCreateForm, setShowCreateForm] = useState(false);

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      setLoading(true);
      const [subjectsData, teachersData] = await Promise.all([
        subjectsApi.getSubjects(),
        usersApi.getUsers({ role: 'teacher', limit: 100 })
      ]);
      setSubjects(subjectsData);
      setTeachers(teachersData.users);
    } catch (err: any) {
      setError('Ошибка загрузки данных');
    } finally {
      setLoading(false);
    }
  };

  const createSubject = async (data: any) => {
    try {
      await subjectsApi.createSubject(data);
      fetchData();
      setShowCreateForm(false);
    } catch (err: any) {
      setError('Ошибка создания дисциплины');
    }
  };

const assignTeacher = async (subjectId: number, teacherId: number) => {
  try {
    await subjectsApi.assignTeacher(subjectId, {
      teacher_id: teacherId, 
      academic_year: '2024-2025',
      is_lead: false
    });
    fetchData();
  } catch (err: any) {
    setError('Ошибка назначения преподавателя');
  }
};

  if (loading) return <LoadingSpinner />;

  return (
    <div>
      <div className="mb-6 flex justify-between items-center">
        <h2 className="text-2xl font-bold text-orange-500">Управление дисциплинами</h2>
        <div className="flex gap-3">
          <Button variant="secondary" onClick={onBack}>Назад</Button>
          <Button onClick={() => setShowCreateForm(true)}>Создать дисциплину</Button>
        </div>
      </div>

      {error && <ErrorMessage message={error} />}

      {showCreateForm && (
        <div className="mb-6 bg-gray-800 p-6 rounded-lg">
          <CreateSubjectForm onSubmit={createSubject} onCancel={() => setShowCreateForm(false)} />
        </div>
      )}

      <div className="space-y-4">
        {subjects.map(subject => (
          <div key={subject.id} className="bg-gray-800 p-4 rounded-lg">
            <div className="flex justify-between items-start">
              <div>
                <h3 className="text-lg font-semibold text-orange-500">{subject.name}</h3>
                <p className="text-gray-400">Код: {subject.code} | Семестр: {subject.semester}</p>
                <p className="text-sm text-gray-500">{subject.description}</p>
              </div>
              <div className="flex gap-2">
                <select 
                  onChange={(e) => e.target.value && assignTeacher(subject.id, parseInt(e.target.value))}
                  className="px-3 py-1 bg-gray-700 text-white rounded"
                  defaultValue=""
                >
                  <option value="">Назначить преподавателя</option>
                  {teachers.map(teacher => (
                    <option key={teacher.id} value={teacher.id}>
                      {teacher.first_name} {teacher.last_name}
                    </option>
                  ))}
                </select>
              </div>
            </div>
            
            {subject.teachers && subject.teachers.length > 0 && (
              <div className="mt-3">
                <p className="text-sm text-gray-400">Преподаватели:</p>
                <div className="flex flex-wrap gap-2 mt-1">
                  {subject.teachers.map(teacher => (
                    <span key={teacher.id} className="bg-blue-600 text-white px-2 py-1 rounded text-sm">
                      {teacher.first_name} {teacher.last_name}
                    </span>
                  ))}
                </div>
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
};

// Простая форма создания дисциплины
const CreateSubjectForm: React.FC<{ onSubmit: (data: any) => void; onCancel: () => void }> = ({ onSubmit, onCancel }) => {
  const [formData, setFormData] = useState({
    name: '',
    code: '',
    description: '',
    semester: 1
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(formData);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <h3 className="text-xl text-orange-500">Создать дисциплину</h3>
      
      <div className="grid grid-cols-2 gap-4">
        <input
          placeholder="Название дисциплины"
          value={formData.name}
          onChange={(e) => setFormData({...formData, name: e.target.value})}
          className="px-3 py-2 bg-gray-700 text-white rounded"
          required
        />
        <input
          placeholder="Код (например: CS101)"
          value={formData.code}
          onChange={(e) => setFormData({...formData, code: e.target.value})}
          className="px-3 py-2 bg-gray-700 text-white rounded"
          required
        />
      </div>

      <select
        value={formData.semester}
        onChange={(e) => setFormData({...formData, semester: parseInt(e.target.value)})}
        className="w-full px-3 py-2 bg-gray-700 text-white rounded"
      >
        {Array.from({length: 8}, (_, i) => (
          <option key={i+1} value={i+1}>{i+1} семестр</option>
        ))}
      </select>

      <textarea
        placeholder="Описание дисциплины"
        value={formData.description}
        onChange={(e) => setFormData({...formData, description: e.target.value})}
        className="w-full px-3 py-2 bg-gray-700 text-white rounded"
        rows={3}
      />

      <div className="flex gap-3">
        <Button type="submit">Создать</Button>
        <Button type="button" variant="secondary" onClick={onCancel}>Отмена</Button>
      </div>
    </form>
  );
};

export default AdminSubjects;