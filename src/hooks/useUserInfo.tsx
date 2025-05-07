import React from 'react';

interface UserInfo {
  id: number;
  name: string;
  email: string;
  phone: string;
}


const UserInfoContext = React.createContext<[UserInfo | null, () => Promise<void>]>([null, () => Promise.resolve()]);

export function UserInfoProvider({ children }: React.PropsWithChildren) {
  const [userInfo, setUserInfo] = React.useState<UserInfo | null>(null);

  const fetchUserInfo = React.useCallback(async () => {
    const response = await fetch('/api/user/info');
    const data = await response.json();
    setUserInfo(data);
  }, []);

  React.useEffect(() => {
    fetchUserInfo();
  }, [fetchUserInfo]);

  return (
    <UserInfoContext value={[userInfo, fetchUserInfo]}>
      {children}
    </UserInfoContext>
  );
}

export const useUserInfo = () => React.useContext(UserInfoContext);
