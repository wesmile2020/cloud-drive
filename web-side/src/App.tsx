import React from 'react';
import { HashRouter, Routes, Route, Navigate } from 'react-router';
import HomeLayout from '@/layouts/HomeLayout';
import { UserInfoProvider } from './hooks/useUserInfo';
import { UploadProvider } from './hooks/useUpload';

const LoginPage = React.lazy(() => import('./pages/LoginPage'));
const RegisterPage = React.lazy(() => import('./pages/RegisterPage'));
const UserPage = React.lazy(() => import('./pages/UserPage'));
const HomePage = React.lazy(() => import('./pages/HomePage'));

const HomeLayoutWithProvider = () => {
  return (
    <UserInfoProvider>
      <UploadProvider>
        <HomeLayout />
      </UploadProvider>
    </UserInfoProvider>
  ); 
}

function App() {
  return (
    <HashRouter>
      <Routes>
        <Route path="/" element={<HomeLayoutWithProvider />}>
          <Route path="/home/:id" element={<HomePage />} />
          <Route path="/user" element={<UserPage />} />
          <Route path="/" element={<Navigate to="/home/0" replace />} />
          <Route path="/home" element={<Navigate to="/home/0" replace />} />
        </Route>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
      </Routes>
    </HashRouter>
  );
};

export default App;
