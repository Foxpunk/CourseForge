import React, { useState } from 'react';
import { useCourseworks } from '../hooks/useCourseworks';
import CourseworkCard from '../components/coursework/CourseworkCard';
import CreateCourseworkForm from '../components/coursework/CreateCourseworkForm';
import LoadingSpinner from '../components/ui/LoadingSpinner';
import ErrorMessage from '../components/ui/ErrorMessage';
import Button from '../components/ui/Button';

const TeacherDashboard: React.FC = () => {
  const { courseworks, loading, error, createCoursework, updateCoursework, deleteCoursework, refetch } = useCourseworks();
  const [showCreateForm, setShowCreateForm] = useState(false);

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
        <ErrorMessage message={error} onRetry={refetch} />
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-8 flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-orange-500 mb-2">
            Мои курсовые работы
          </h1>
          <p className="text-gray-400">
            Создавайте и управляйте курсовыми работами для ваших дисциплин
          </p>
        </div>
        <Button
          variant="primary"
          onClick={() => setShowCreateForm(true)}
        >
          Создать курсовую
        </Button>
      </div>

      {showCreateForm && (
        <div className="mb-8">
          <CreateCourseworkForm
            onSubmit={async (data) => {
              await createCoursework(data);
              setShowCreateForm(false);
            }}
            onCancel={() => setShowCreateForm(false)}
          />
        </div>
      )}

      {courseworks.length === 0 ? (
        <div className="text-center py-12">
          <div className="max-w-md mx-auto">
            <div className="mb-4">
              <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.746 0 3.332.477 4.5 1.253v13C19.832 18.477 18.246 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
              </svg>
            </div>
            <h3 className="text-lg font-medium text-gray-300 mb-2">
              Нет созданных курсовых работ
            </h3>
            <p className="text-gray-400 mb-4">
              Создайте свою первую курсовую работу для студентов
            </p>
            <Button
              variant="primary"
              onClick={() => setShowCreateForm(true)}
            >
              Создать курсовую
            </Button>
          </div>
        </div>
      ) : (
        <>
          <div className="mb-6 bg-gray-800 rounded-lg p-4">
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4 text-center">
              <div>
                <div className="text-2xl font-bold text-orange-500">{courseworks.length}</div>
                <div className="text-sm text-gray-400">Всего курсовых</div>
              </div>
              <div>
                <div className="text-2xl font-bold text-green-400">
                  {courseworks.filter(c => c.is_available).length}
                </div>
                <div className="text-sm text-gray-400">Доступны</div>
              </div>
              <div>
                <div className="text-2xl font-bold text-yellow-400">
                  {courseworks.filter(c => !c.is_available).length}
                </div>
                <div className="text-sm text-gray-400">Недоступны</div>
              </div>
              <div>
                <div className="text-2xl font-bold text-blue-400">
                  {courseworks.reduce((acc, c) => acc + c.max_students, 0)}
                </div>
                <div className="text-sm text-gray-400">Мест всего</div>
              </div>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {courseworks.map((coursework) => (
              <CourseworkCard
                key={coursework.id}
                coursework={coursework}
                onEdit={(coursework) => {
                  console.log('Edit coursework:', coursework);
                  // Здесь можно открыть форму редактирования
                }}
                onDelete={deleteCoursework}
                onViewDetails={(coursework) => {
                  console.log('View details:', coursework);
                  // Здесь можно открыть модальное окно с деталями и списком студентов
                }}
              />
            ))}
          </div>
        </>
      )}
    </div>
  );
};

export default TeacherDashboard;