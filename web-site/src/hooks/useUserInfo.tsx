import { useState, useRef, useCallback, createContext, useEffect, useContext, PropsWithChildren } from 'react';
import { useLocation, useNavigate } from 'react-router';
import { getUserInfo, UserInfo } from '@/services/api';

type UserInfoContextType = [UserInfo | null, (toLogin?: boolean, reload?: boolean) => Promise<void>];

const UserInfoContext = createContext<UserInfoContextType>([null, () => Promise.resolve()]);

export function UserInfoProvider({ children }: PropsWithChildren) {
  const [userInfo, setUserInfo] = useState<UserInfo | null>(null);
  const fetchPromise = useRef<Promise<UserInfo> | null>(null);
  const location = useLocation();
  const navigate = useNavigate();

  const fetchUserInfo = useCallback(async (toLogin = true, reload = false) => {
    try {
      if (!fetchPromise.current || reload) {
        fetchPromise.current = getUserInfo();
      }
      const data = await fetchPromise.current;
      setUserInfo(data);
    } catch (error) {
      const code = (error as Error).message;
      if (code === 'Unauthorized' && toLogin) {
        navigate(`/login?redirect=${encodeURIComponent(location.pathname)}`);
      }
    }
  }, [location.pathname, navigate]);

  useEffect(() => {
    fetchUserInfo(false);
  }, [fetchUserInfo]);

  return (
    <UserInfoContext value={[userInfo, fetchUserInfo]}>
      {children}
    </UserInfoContext>
  );
}

// eslint-disable-next-line react-refresh/only-export-components
export const useUserInfo = () => useContext(UserInfoContext);
