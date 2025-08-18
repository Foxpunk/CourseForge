import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import { authApi } from '../api/auth';
import { RegisterRequest, UserRole } from '../types';
import Button from '../components/ui/Button';
import ErrorMessage from '../components/ui/ErrorMessage';

const RegisterPage: React.FC = () => {
  const [formData, setFormData] = useState<RegisterRequest>({
    email: '',
    password: '',
    first_name: '',
    last_name: '',
    role: 'student' as UserRole,
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    setFormData(prev => ({
      ...prev,
      [e.target.name]: e.target.value,
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const authData = await authApi.register(formData);
      login(authData);
      navigate('/dashboard');
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || err.message || 'Ошибка при регистрации';
      setError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div>
          <h2 className="mt-6 text-center text-3xl font-extrabold text-orange-500">
            Регистрация
          </h2>
          <p className="mt-2 text-center text-sm text-gray-400">
            Создайте аккаунт в системе управления курсовыми проектами
          </p>
        </div>
        <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
          <div className="bg-gray-800 p-8 rounded-lg shadow-lg">
            {error && (
              <div className="mb-4">
                <ErrorMessage message={error} />
              </div>
            )}
            
            <div className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label htmlFor="first_name" className="block text-sm font-medium text-gray-300">
                    Имя
                  </label>
                  <input
                    id="first_name"
                    name="first_name"
                    type="text"
                    required
                    value={formData.first_name}
                    onChange={handleChange}
                    className="mt-1 block w-full px-3 py-2 border border-gray-600 placeholder-gray-500 text-white bg-gray-700 rounded-md focus:outline-none focus:ring-orange-500 focus:border-orange-500 sm:text-sm"
                  />
                </div>
                <div>
                  <label htmlFor="last_name" className="block text-sm font-medium text-gray-300">
                    Фамилия
                  </label>
                  <input
                    id="last_name"
                    name="last_name"
                    type="text"
                    required
                    value={formData.last_name}
                    onChange={handleChange}
                    className="mt-1 block w-full px-3 py-2 border border-gray-600 placeholder-gray-500 text-white bg-gray-700 rounded-md focus:outline-none focus:ring-orange-500 focus:border-orange-500 sm:text-sm"
                  />
                </div>
              </div>

              <div>
                <label htmlFor="email" className="block text-sm font-medium text-gray-300">
                  Email
                </label>
                <input
                  id="email"
                  name="email"
                  type="email"
                  required
                  value={formData.email}
                  onChange={handleChange}
                  className="mt-1 block w-full px-3 py-2 border border-gray-600 placeholder-gray-500 text-white bg-gray-700 rounded-md focus:outline-none focus:ring-orange-500 focus:border-orange-500 sm:text-sm"
                />
              </div>

              <div>
                <label htmlFor="password" className="block text-sm font-medium text-gray-300">
                  Пароль
                </label>
                <input
                  id="password"
                  name="password"
                  type="password"
                  required
                  value={formData.password}
                  onChange={handleChange}
                  className="mt-1 block w-full px-3 py-2 border border-gray-600 placeholder-gray-500 text-white bg-gray-700 rounded-md focus:outline-none focus:ring-orange-500 focus:border-orange-500 sm:text-sm"
                  placeholder="Минимум 6 символов"
                />
              </div>

              <div>
                <label htmlFor="role" className="block text-sm font-medium text-gray-300">
                  Роль
                </label>
                <select
                  id="role"
                  name="role"
                  value={formData.role}
                  onChange={handleChange}
                  className="mt-1 block w-full px-3 py-2 border border-gray-600 text-white bg-gray-700 rounded-md focus:outline-none focus:ring-orange-500 focus:border-orange-500 sm:text-sm"
                >
                  <option value="student">Студент</option>
                  <option value="teacher">Преподаватель</option>
                  <option value="admin">Администратор</option>
                </select>
              </div>
            </div>

            <div className="mt-6">
              <Button
                type="submit"
                loading={loading}
                className="w-full"
              >
                Зарегистрироваться
              </Button>
            </div>

            <div className="text-center mt-4">
              <span className="text-gray-400">Есть аккаунт? </span>
              <Link
                to="/login"
                className="font-medium text-orange-500 hover:text-orange-400"
              >
                Войти
              </Link>
            </div>
          </div>
        </form>
      </div>
    </div>
  );
};

export default RegisterPage;