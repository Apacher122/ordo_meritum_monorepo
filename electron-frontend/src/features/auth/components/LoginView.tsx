import '@/assets/styles/views/LoginView.css';

import React from 'react';
import { useLoginForm } from '../hooks/useLoginForm';

const LoginInput = ({ label, type, value, onChange }: any) => (
  <div>
    <label className="block text-sm font-medium">{label}</label>
    <input
      type={type}
      value={value}
      onChange={e => onChange(e.target.value)}
      className="w-full border rounded p-2 mt-1"
      required
    />
  </div>
);

export const LoginView: React.FC = () => {
  const {
    mode,
    email,
    password,
    error,
    loading,
    setEmail,
    setPassword,
    handleSubmit,
    toggleMode,
  } = useLoginForm();

  return (
    <div className="flex items-center justify-center h-screen bg-gray-100">
      <div className="bg-white shadow-lg rounded-lg p-6 w-96">
        <h2 className="text-2xl font-bold mb-4 text-center">
          {mode === 'login' ? 'Login' : 'Register'}
        </h2>

        {error && (
          <div className="bg-red-100 text-red-700 p-2 rounded mb-4">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <LoginInput label="Email" type="email" value={email} onChange={setEmail} />
          <LoginInput
            label="Password"
            type="password"
            value={password}
            onChange={setPassword}
          />
          <button
            type="submit"
            className="w-full bg-blue-600 text-white rounded p-2 hover:bg-blue-700"
            disabled={loading}
          >
            {loading ? 'Loading...' : mode === 'login' ? 'Log In' : 'Register'}
          </button>
        </form>

        <p className="text-sm mt-4 text-center">
          {mode === 'login'
            ? "Don't have an account?"
            : 'Already have an account?'}{' '}
          <button
            onClick={toggleMode}
            className="text-blue-600 hover:underline"
          >
            {mode === 'login' ? 'Register' : 'Login'}
          </button>
        </p>
      </div>
    </div>
  );
};