import React from 'react';

interface ButtonProps {
  variant?: 'primary' | 'secondary' | 'danger';
  size?: 'small' | 'medium' | 'large';
  disabled?: boolean;
  loading?: boolean;
  onClick?: () => void;
  children: React.ReactNode;
  className?: string;
  type?: 'button' | 'submit' | 'reset';
}

const Button: React.FC<ButtonProps> = ({
  variant = 'primary',
  size = 'medium',
  disabled = false,
  loading = false,
  onClick,
  children,
  className = '',
  type = 'button',
}) => {
  const baseClasses = 'font-medium rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-gray-900';
  
  const variantClasses = {
    primary: 'bg-orange-600 hover:bg-orange-700 text-white focus:ring-orange-500',
    secondary: 'bg-gray-600 hover:bg-gray-700 text-white focus:ring-gray-500',
    danger: 'bg-red-600 hover:bg-red-700 text-white focus:ring-red-500',
  };

  const sizeClasses = {
    small: 'px-3 py-1.5 text-sm',
    medium: 'px-4 py-2 text-base',
    large: 'px-6 py-3 text-lg',
  };

  const disabledClasses = disabled || loading ? 'opacity-50 cursor-not-allowed' : '';

  return (
    <button
      type={type}
      onClick={onClick}
      disabled={disabled || loading}
      className={`${baseClasses} ${variantClasses[variant]} ${sizeClasses[size]} ${disabledClasses} ${className}`}
    >
      {loading ? (
        <div className="flex items-center">
          <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
          Загрузка...
        </div>
      ) : (
        children
      )}
    </button>
  );
};

export default Button;