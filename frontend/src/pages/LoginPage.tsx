import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import { authApi } from '../api/auth';
import { LoginRequest } from '../types';
import Button from '../components/ui/Button';
import ErrorMessage from '../components/ui/ErrorMessage';

const LoginPage: React.FC = () => {
  const [formData, setFormData] = useState<LoginRequest>({
    email: '',
    password: '',
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
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
      const authData = await authApi.login(formData);
      login(authData);
      navigate('/dashboard');
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || err.message || 'Ошибка при входе в систему';
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
            Вход в систему
          </h2>
          <p className="mt-2 text-center text-sm text-gray-400">
            Войдите в систему управления курсовыми проектами
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
                  className="mt-1 appearance-none relative block w-full px-3 py-2 border border-gray-600 placeholder-gray-500 text-white bg-gray-700 rounded-md focus:outline-none focus:ring-orange-500 focus:border-orange-500 focus:z-10 sm:text-sm"
                  placeholder="Введите email"
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
                  className="mt-1 appearance-none relative block w-full px-3 py-2 border border-gray-600 placeholder-gray-500 text-white bg-gray-700 rounded-md focus:outline-none focus:ring-orange-500 focus:border-orange-500 focus:z-10 sm:text-sm"
                  placeholder="Введите пароль"
                />
              </div>
            </div>

            <div className="mt-6">
              <Button
                type="submit"
                loading={loading}
                className="w-full"
              >
                Войти
              </Button>
            </div>

            <div className="text-center mt-4">
              <span className="text-gray-400">Нет аккаунта? </span>
              <Link
                to="/register"
                className="font-medium text-orange-500 hover:text-orange-400"
              >
                Зарегистрироваться
              </Link>
            </div>
          </div>
        </form>
      </div>
    </div>
  );
};

export default LoginPage;