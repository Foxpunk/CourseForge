import React from 'react';
import { useAuth } from '../../hooks/useAuth';
import Button from '../ui/Button';

const Header: React.FC = () => {
  const { user, logout } = useAuth();
  const [isLoggingOut, setIsLoggingOut] = React.useState(false);

  const getRoleText = (role: string) => {
    switch (role) {
      case 'student': return 'Студент';
      case 'teacher': return 'Преподаватель';
      case 'admin': return 'Администратор';
      default: return 'Пользователь';
    }
  };

  const getRoleColor = (role: string) => {
    switch (role) {
      case 'student': return 'bg-blue-600';
      case 'teacher': return 'bg-green-600';
      case 'admin': return 'bg-red-600';
      default: return 'bg-gray-600';
    }
  };

  const handleLogout = async () => {
    setIsLoggingOut(true);
    try {
      await logout();
    } finally {
      setIsLoggingOut(false);
    }
  };

  return (
    <header className="bg-gray-800 border-b-2 border-orange-500 sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          <div className="flex items-center">
            <div className="flex-shrink-0">
              <h1 className="text-2xl font-bold text-orange-500">CourseForge</h1>
            </div>
            <div className="hidden md:block ml-6">
              <p className="text-sm text-gray-400">
                Система управления курсовыми проектами
              </p>
            </div>
          </div>
          
          {user && (
            <div className="flex items-center space-x-4">
              <div className="flex items-center space-x-3">
                <span className={`${getRoleColor(user.role)} text-white px-3 py-1 rounded-full text-sm font-medium`}>
                  {getRoleText(user.role)}
                </span>
                <div className="text-right">
                  <div className="text-gray-300 font-medium">
                    {user.first_name} {user.last_name}
                  </div>
                  <div className="text-xs text-gray-500">
                    {user.email}
                  </div>
                </div>
              </div>
              <Button 
                variant="secondary" 
                size="small" 
                onClick={handleLogout}
                loading={isLoggingOut}
              >
                Выход
              </Button>
            </div>
          )}
        </div>
      </div>
    </header>
  );
};

export default Header;