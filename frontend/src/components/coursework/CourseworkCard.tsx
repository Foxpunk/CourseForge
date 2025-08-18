import React from 'react';
import { CourseworkResponse } from '../../types';
import { useAuth } from '../../hooks/useAuth';
import Button from '../ui/Button';

interface CourseworkCardProps {
  coursework: CourseworkResponse;
  onAssign?: (courseworkId: number) => Promise<void>;
  onEdit?: (coursework: CourseworkResponse) => void;
  onDelete?: (courseworkId: number) => Promise<void>;
  onViewDetails?: (coursework: CourseworkResponse) => void;
  isUserAssigned?: boolean; // Выбрал ли текущий пользователь эту работу
  isFullyAssigned?: boolean; // Заполнены ли все места
  userHasAnyAssignment?: boolean; // Есть ли у пользователя любая назначенная работа
}

const CourseworkCard: React.FC<CourseworkCardProps> = ({
  coursework,
  onAssign,
  onEdit,
  onDelete,
  onViewDetails,
  isUserAssigned = false,
  isFullyAssigned = false,
  userHasAnyAssignment = false,
}) => {
  const { user } = useAuth();
  const [isAssigning, setIsAssigning] = React.useState(false);
  const [isDeleting, setIsDeleting] = React.useState(false);
  const [showSuccess, setShowSuccess] = React.useState(false);

  const getDifficultyColor = (difficulty: string) => {
    switch (difficulty) {
      case 'easy': return 'text-green-400';
      case 'medium': return 'text-yellow-400';
      case 'hard': return 'text-red-400';
      default: return 'text-gray-400';
    }
  };

  const getDifficultyText = (difficulty: string) => {
    switch (difficulty) {
      case 'easy': return 'Легкая';
      case 'medium': return 'Средняя';
      case 'hard': return 'Сложная';
      default: return 'Не указана';
    }
  };

  const handleAssign = async () => {
    if (!onAssign) return;
    
    setIsAssigning(true);
    try {
      await onAssign(coursework.id);
      setShowSuccess(true);
      setTimeout(() => setShowSuccess(false), 3000);
    } catch (error) {
      console.error('Error assigning coursework:', error);
    } finally {
      setIsAssigning(false);
    }
  };

  const handleDelete = async () => {
    if (!onDelete) return;
    
    if (!window.confirm('Вы уверены, что хотите удалить эту курсовую работу?')) {
      return;
    }

    setIsDeleting(true);
    try {
      await onDelete(coursework.id);
    } catch (error) {
      console.error('Error deleting coursework:', error);
    } finally {
      setIsDeleting(false);
    }
  };

  // Определяем состояние карточки
  const cardStyle = () => {
    if (isUserAssigned) {
      return 'bg-green-900/20 border-green-500'; // Выбрана пользователем
    }
    if (isFullyAssigned) {
      return 'bg-gray-900/50 border-gray-600 opacity-75'; // Уже занята
    }
    if (userHasAnyAssignment) {
      return 'bg-gray-800 border-gray-600 opacity-60'; // У пользователя есть другая работа
    }
    return 'bg-gray-800 border-gray-700 hover:border-orange-500'; // Доступна
  };

  return (
    <div className={`rounded-lg p-6 border transition-all relative ${cardStyle()}`}>
      
      {/* Уведомление об успешном выборе */}
      {showSuccess && (
        <div className="absolute top-2 right-2 bg-green-600 text-white px-3 py-1 rounded-md text-sm font-medium z-10 animate-pulse">
          ✓ Работа выбрана!
        </div>
      )}

      {/* Метки состояния */}
      {isUserAssigned && (
        <div className="absolute top-4 left-4 bg-green-600 text-white px-2 py-1 rounded text-xs font-bold">
          ВЫБРАНА ВАМИ
        </div>
      )}
      
      {isFullyAssigned && !isUserAssigned && (
        <div className="absolute top-4 left-4 bg-red-600 text-white px-2 py-1 rounded text-xs font-bold">
          ЗАНЯТА
        </div>
      )}

      <div className="flex justify-between items-start mb-3">
        <h3 className={`text-xl font-semibold line-clamp-2 ${
          isUserAssigned ? 'text-green-400' : 
          (isFullyAssigned || userHasAnyAssignment) ? 'text-gray-400' : 'text-orange-500'
        }`}>
          {coursework.title}
        </h3>
        <span className={`text-sm font-medium ${getDifficultyColor(coursework.difficulty_level)}`}>
          {getDifficultyText(coursework.difficulty_level)}
        </span>
      </div>

      <div className="space-y-2 mb-4 text-sm text-gray-300">
        <p>
          <span className="font-medium">Преподаватель:</span> {coursework.teacher.first_name} {coursework.teacher.last_name}
        </p>
        <p>
          <span className="font-medium">Дисциплина:</span> {coursework.subject.name} ({coursework.subject.code})
        </p>
        <p>
          <span className="font-medium">Семестр:</span> {coursework.subject.semester}
        </p>
        <p>
          <span className="font-medium">Доступна:</span> 
          <span className={coursework.is_available ? 'text-green-400' : 'text-red-400'}>
            {coursework.is_available ? ' Да' : ' Нет'}
          </span>
        </p>
        <p className="text-xs text-gray-500">
          Создана: {new Date(coursework.created_at).toLocaleDateString('ru-RU')}
        </p>
      </div>

      <div className="mb-4">
        <p className="text-gray-400 text-sm line-clamp-3">
          {coursework.description}
        </p>
      </div>

      {coursework.requirements && (
        <div className="mb-4">
          <p className="text-sm font-medium text-gray-300 mb-1">Требования:</p>
          <p className="text-gray-400 text-sm line-clamp-2">
            {coursework.requirements}
          </p>
        </div>
      )}

      <div className="flex flex-wrap gap-2">
        {/* Кнопки для студентов */}
        {user?.role === 'student' && (
          <>
            {isUserAssigned ? (
              // Пользователь выбрал эту работу
              <div className="bg-green-600 text-white px-3 py-1 rounded text-sm font-medium">
                ✓ Вы выбрали эту работу
              </div>
            ) : isFullyAssigned ? (
              // Работа уже занята другими
              <div className="bg-gray-600 text-gray-300 px-3 py-1 rounded text-sm">
                Выбрано
              </div>
            ) : userHasAnyAssignment ? (
              // У пользователя уже есть назначенная работа
              <Button 
                variant="secondary" 
                size="small" 
                disabled={true}
                className="opacity-50 cursor-not-allowed"
              >
                Недоступно
              </Button>
            ) : onAssign ? (
              // Доступна для выбора
              <Button 
                variant="primary" 
                size="small" 
                onClick={handleAssign}
                loading={isAssigning}
                disabled={!coursework.is_available || isAssigning}
              >
                {isAssigning ? 'Выбираю...' : 'Выбрать'}
              </Button>
            ) : null}
          </>
        )}
        
        {/* Кнопки для преподавателей */}
        {(user?.role === 'teacher' || user?.role === 'admin') && onEdit && (
          <Button 
            variant="secondary" 
            size="small" 
            onClick={() => onEdit(coursework)}
          >
            Редактировать
          </Button>
        )}
        
        {(user?.role === 'teacher' || user?.role === 'admin') && onDelete && (
          <Button 
            variant="danger" 
            size="small" 
            onClick={handleDelete}
            loading={isDeleting}
          >
            Удалить
          </Button>
        )}

        {onViewDetails && (
          <Button 
            variant="secondary" 
            size="small" 
            onClick={() => onViewDetails(coursework)}
          >
            Подробнее
          </Button>
        )}
      </div>
    </div>
  );
};

export default CourseworkCard;