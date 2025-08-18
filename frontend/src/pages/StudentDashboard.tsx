import React, { useState } from 'react';
import { useCourseworks } from '../hooks/useCourseworks';
import CourseworkCard from '../components/coursework/CourseworkCard';
import LoadingSpinner from '../components/ui/LoadingSpinner';
import ErrorMessage from '../components/ui/ErrorMessage';

const StudentDashboard: React.FC = () => {
  const { 
    courseworks, 
    loading, 
    error, 
    assignStudent, 
    refetch, 
    userHasAssignment, 
    userAssignedCourseworkId 
  } = useCourseworks();
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  const handleAssignStudent = async (courseworkId: number) => {
    try {
      await assignStudent(courseworkId);
      setSuccessMessage('Курсовая работа успешно выбрана! Теперь вы можете приступить к выполнению.');
      setTimeout(() => setSuccessMessage(null), 5000);
    } catch (err: any) {
      console.error('Assignment failed:', err);
    }
  };

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
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-orange-500 mb-2">
          Доступные курсовые работы
        </h1>
        <p className="text-gray-400">
          Выберите курсовую работу для выполнения. Вы можете выбрать только одну курсовую работу.
        </p>
        
        {/* Уведомление об успехе */}
        {successMessage && (
          <div className="mt-4 bg-green-900/20 border border-green-500 rounded-lg p-4">
            <div className="flex items-center">
              <svg className="h-5 w-5 text-green-400 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
              </svg>
              <p className="text-green-200 font-medium">{successMessage}</p>
            </div>
          </div>
        )}

        {/* Уведомление о том что работа уже выбрана */}
        {userHasAssignment && (
          <div className="mt-4 bg-blue-900/20 border border-blue-500 rounded-lg p-4">
            <p className="text-blue-200 text-sm">
              <strong>Внимание:</strong> Вы уже выбрали курсовую работу. Другие работы недоступны для выбора.
            </p>
          </div>
        )}

        <div className="mt-4 bg-blue-900/20 border border-blue-500 rounded-lg p-4">
          <p className="text-blue-200 text-sm">
            <strong>Совет:</strong> Внимательно изучите описание и требования к курсовой работе перед выбором.
            После выбора вы сможете приступить к выполнению работы под руководством преподавателя.
          </p>
        </div>
      </div>

      {courseworks.length === 0 ? (
        <div className="text-center py-12">
          <div className="max-w-md mx-auto">
            <div className="mb-4">
              <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
            </div>
            <h3 className="text-lg font-medium text-gray-300 mb-2">
              Нет доступных курсовых работ
            </h3>
            <p className="text-gray-400">
              В данный момент преподаватели не опубликовали курсовые работы для выбора.
              Обратитесь к своему куратору или проверьте позже.
            </p>
          </div>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {courseworks.map((coursework) => (
            <CourseworkCard
              key={coursework.id}
              coursework={coursework}
              onAssign={handleAssignStudent}
              isUserAssigned={userAssignedCourseworkId === coursework.id}
              isFullyAssigned={false} // TODO: добавить логику проверки заполненности
              userHasAnyAssignment={userHasAssignment}
              onViewDetails={(coursework) => {
                console.log('View details:', coursework);
              }}
            />
          ))}
        </div>
      )}
    </div>
  );
};

export default StudentDashboard;