import { getUserInfo, UserInfo } from '@/services/api';
import React from 'react';
import { useLocation, useNavigate } from 'react-router';

type UserInfoContextType = [UserInfo | null, (toLogin?: boolean, reload?: boolean) => Promise<void>];

const UserInfoContext = React.createContext<UserInfoContextType>([null, () => Promise.resolve()]);

export function UserInfoProvider({ children }: React.PropsWithChildren) {
  const [userInfo, setUserInfo] = React.useState<UserInfo | null>(null);
  const fetchPromise = React.useRef<Promise<UserInfo> | null>(null);
  const location = useLocation();
  const navigate = useNavigate();

  const fetchUserInfo = React.useCallback(async (toLogin = true, reload = false) => {
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

  React.useEffect(() => {
    fetchUserInfo(false);
  }, [fetchUserInfo]);

  return (
    <UserInfoContext value={[userInfo, fetchUserInfo]}>
      {children}
    </UserInfoContext>
  );
}

// eslint-disable-next-line react-refresh/only-export-components
export const useUserInfo = () => React.useContext(UserInfoContext);
