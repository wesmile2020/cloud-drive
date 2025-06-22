import React from 'react';
import { useUserInfo } from '@/hooks/useUserInfo';

import styles from './UserPage.module.css';

function User() {
  const [userInfo, fetchUserInfo] = useUserInfo();

  React.useEffect(() => {
    if (!userInfo) {
      fetchUserInfo(true);
    }
  }, [userInfo, fetchUserInfo]);

  return (
    <div className={styles.user_page}>

    </div>
  );
}

export default User;
