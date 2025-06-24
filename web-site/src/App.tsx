import { lazy } from 'react';
import { HashRouter, Routes, Route, Navigate } from 'react-router';
import HomeLayout from '@/layouts/HomeLayout';
import { UserInfoProvider } from './hooks/useUserInfo';
import { UploadProvider } from './hooks/useUpload';

const LoginPage = lazy(() => import('./pages/LoginPage'));
const RegisterPage = lazy(() => import('./pages/RegisterPage'));
const UserPage = lazy(() => import('./pages/UserPage'));
const HomePage = lazy(() => import('./pages/HomePage'));
const RetrievePassword = lazy(() => import('./pages/RetrievePassword'));

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
        <Route path="/retrieve-password" element={<RetrievePassword />} />
      </Routes>
    </HashRouter>
  );
};

export default App;
