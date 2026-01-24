import {
  createUserWithEmailAndPassword,
  getAuth,
  signInWithEmailAndPassword,
} from 'firebase/auth';

import { useNavigate } from 'react-router-dom';
import { useState } from 'react';

/**
 * Custom hook to manage the state and logic for a login/registration form.
 * It handles Firebase authentication for both signing in and creating new users.
 * @returns {UseLoginFormReturn} An object containing form state and handler functions.
 */
export const useLoginForm = () => {
  const [mode, setMode] = useState<'login' | 'register'>('login');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const auth = getAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      if (mode === 'login') {
        await signInWithEmailAndPassword(auth, email, password);
      } else {
        await createUserWithEmailAndPassword(auth, email, password);
      }
      navigate('/');
    } catch (err: any) {
      setError(err.message || 'Authentication failed. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const toggleMode = () => {
    setMode(prevMode => (prevMode === 'login' ? 'register' : 'login'));
  };

  return {
    mode,
    email,
    password,
    error,
    loading,
    setEmail,
    setPassword,
    handleSubmit,
    toggleMode,
  };
};