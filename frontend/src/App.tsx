import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthContext, useAuthProvider, useAuth } from './hooks/useAuth';
import Header from './components/layout/Header';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import StudentDashboard from './pages/StudentDashboard';
import TeacherDashboard from './pages/TeacherDashboard';
import AdminDashboard from './pages/AdminDashboard';
import ProtectedRoute from './components/auth/ProtectedRoute';
import LoadingSpinner from './components/ui/LoadingSpinner';
import './styles/globals.css';

const App: React.FC = () => {
  const auth = useAuthProvider();

  if (auth.loading) {
    return (
      <div className="min-h-screen bg-gray-900 flex items-center justify-center">
        <LoadingSpinner />
      </div>
    );
  }

  return (
    <AuthContext.Provider value={auth}>
      <Router>
        <div className="min-h-screen bg-gray-900 text-white">
          <Header />
          <main>
            <Routes>
              <Route 
                path="/login" 
                element={
                  auth.isAuthenticated ? (
                    <Navigate to="/dashboard" replace />
                  ) : (
                    <LoginPage />
                  )
                } 
              />
              <Route 
                path="/register" 
                element={
                  auth.isAuthenticated ? (
                    <Navigate to="/dashboard" replace />
                  ) : (
                    <RegisterPage />
                  )
                } 
              />
              <Route
                path="/dashboard"
                element={
                  <ProtectedRoute>
                    <DashboardRouter />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/"
                element={
                  auth.isAuthenticated ? (
                    <Navigate to="/dashboard" replace />
                  ) : (
                    <Navigate to="/login" replace />
                  )
                }
              />
            </Routes>
          </main>
        </div>
      </Router>
    </AuthContext.Provider>
  );
};

const DashboardRouter: React.FC = () => {
  const { user } = useAuth(); 

  console.log('DashboardRouter - Current user:', user);

  switch (user?.role) {
    case 'student':
      console.log('Redirecting to StudentDashboard');
      return <StudentDashboard />;
    case 'teacher':
      console.log('Redirecting to TeacherDashboard');
      return <TeacherDashboard />;
    case 'admin':
      console.log('Redirecting to AdminDashboard');
      return <AdminDashboard />;
    default:
      console.log('No matching role, user:', user);
      return <Navigate to="/login" replace />;
  }
};

export default App;