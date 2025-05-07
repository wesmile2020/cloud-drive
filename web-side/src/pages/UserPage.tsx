import React from 'react';
import { useUserInfo } from '@/hooks/useUserInfo';

function User() {
  const [userInfo, fetchUserInfo] = useUserInfo();

  React.useEffect(() => {
    if (!userInfo) {
      fetchUserInfo();
    }
  }, [userInfo, fetchUserInfo]);

  return (
    <div>User</div>
  );
}

export default User;
