import * as pages from './pages';

import { Navigate, Outlet, Route, Routes } from 'react-router-dom';

import { MainShell } from '@/components/Layouts/MainShell';
import { useAuth } from '@/features/auth/providers/AuthProvider';

const ProtectedRoute = () => {
  const { user } = useAuth();
  return user ? <Outlet /> : <Navigate to="/login" />;
};

export const AppRoutes = () => {
  return (
    <Routes>
      <Route path="/login" element={<pages.LoginView />} />
      <Route element={<ProtectedRoute />}>
        <Route path="/" element={<MainShell />}>
          <Route index element={<Navigate to="/info" />} />
          
          <Route path="info" element={<pages.InfoPage />} />
          <Route path="applications" element={<pages.ApplicationListPage />} />
          <Route path="user-info" element={<pages.UserInfoPage />} />
          
          <Route path="documents" element={<pages.DocumentPage />} />
          <Route path="settings" element={<pages.SettingsPage />} />
        </Route>
      </Route>
    </Routes>
  );
};