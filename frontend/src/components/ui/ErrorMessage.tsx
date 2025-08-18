import React from 'react';

interface ErrorMessageProps {
  message: string;
  onRetry?: () => void;
}

const ErrorMessage: React.FC<ErrorMessageProps> = ({ message, onRetry }) => {
  return (
    <div className="bg-red-900/20 border border-red-500 text-red-200 px-4 py-3 rounded-lg">
      <div className="flex items-center justify-between">
        <span>{message}</span>
        {onRetry && (
          <button
            onClick={onRetry}
            className="ml-4 px-3 py-1 bg-red-600 hover:bg-red-700 rounded text-sm"
          >
            Повторить
          </button>
        )}
      </div>
    </div>
  );
};

export default ErrorMessage;