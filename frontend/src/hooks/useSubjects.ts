import { useState, useEffect } from 'react';
import { SubjectResponse } from '../types';
import { subjectsApi } from '../api/subjects';

export const useSubjects = () => {
  const [subjects, setSubjects] = useState<SubjectResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchSubjects = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await subjectsApi.getSubjects();
      setSubjects(data);
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || err.message || 'Ошибка при загрузке дисциплин';
      setError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchSubjects();
  }, []);

  return {
    subjects,
    loading,
    error,
    refetch: fetchSubjects,
  };
};