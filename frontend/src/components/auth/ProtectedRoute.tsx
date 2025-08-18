import React from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';
import LoadingSpinner from '../ui/LoadingSpinner';

interface ProtectedRouteProps {
  children: React.ReactNode;
  requiredRole?: 'admin' | 'teacher' | 'student';
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children, requiredRole }) => {
  const { isAuthenticated, user, loading } = useAuth();

  if (loading) {
    return <LoadingSpinner />;
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  if (requiredRole && user?.role !== requiredRole) {
    return <Navigate to="/dashboard" replace />;
  }

  return <>{children}</>;
};

export default ProtectedRoute;