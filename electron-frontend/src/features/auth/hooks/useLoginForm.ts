import {
  createUserWithEmailAndPassword,
  getAuth,
  signInWithEmailAndPassword,
} from 'firebase/auth';

import { useNavigate } from 'react-router-dom';
import { useState } from 'react';



/**
 * A hook that provides the state and functions for handling login and registration forms.
 *
 * @returns An object with the following properties:
 *   - mode: The mode of the form ('login' or 'register')
 *   - email: The email address of the user
 *   - password: The password of the user
 *   - error: The error message if the form submission failed
 *   - loading: Whether the form is currently submitting
 *   - setEmail: A function to set the email address of the user
 *   - setPassword: A function to set the password of the user
 *   - handleSubmit: A function to handle the form submission
 *   - toggleMode: A function to toggle the mode of the form between 'login' and 'register'
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