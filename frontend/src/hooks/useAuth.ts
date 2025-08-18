import { useState, useEffect, createContext, useContext } from 'react';
import { UserResponse, LoginResponse } from '../types';
import { authApi } from '../api/auth';

interface AuthContextType {
  user: UserResponse | null;
  isAuthenticated: boolean;
  login: (authData: LoginResponse) => void;
  logout: () => void;
  loading: boolean;
  refreshProfile: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

export const useAuthProvider = () => {
  const [user, setUser] = useState<UserResponse | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const initAuth = async () => {
      const storedUser = localStorage.getItem('user');
      const storedToken = localStorage.getItem('token');

      if (storedUser && storedToken) {
        try {
          // Используем сохраненные данные вместо API запроса при инициализации
          const parsedUser = JSON.parse(storedUser);
          setUser(parsedUser);
          console.log('User loaded from localStorage:', parsedUser);
        } catch (error) {
          console.error('Error parsing stored user:', error);
          localStorage.removeItem('user');
          localStorage.removeItem('token');
        }
      }
      setLoading(false);
    };

    initAuth();
  }, []);

  const login = (authData: LoginResponse) => {
    console.log('Login called with:', authData);
    setUser(authData.user);
    localStorage.setItem('user', JSON.stringify(authData.user));
    localStorage.setItem('token', authData.token);
  };

  const logout = async () => {
    try {
      await authApi.logout();
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      setUser(null);
      localStorage.removeItem('user');
      localStorage.removeItem('token');
    }
  };

  const refreshProfile = async () => {
    try {
      const currentUser = await authApi.getProfile();
      setUser(currentUser);
      localStorage.setItem('user', JSON.stringify(currentUser));
    } catch (error) {
      console.error('Error refreshing profile:', error);
      throw error;
    }
  };

  return {
    user,
    isAuthenticated: !!user,
    login,
    logout,
    loading,
    refreshProfile,
  };
};

export { AuthContext };