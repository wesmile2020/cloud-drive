import React from 'react';
import { useUserInfo } from '@/hooks/useUserInfo';

import styles from './UserPage.module.css';

function UserPage() {
  const [userInfo, fetchUserInfo] = useUserInfo();

  React.useEffect(() => {
    if (!userInfo) {
      fetchUserInfo(true);
    }
  }, [userInfo, fetchUserInfo]);

  return userInfo && (
    <div className={styles.user_page}>
      <h2 className={styles.title}>用户信息</h2>
      <div className={styles.info_item}>
        <span className={styles.key}>用户名</span>
        <span className={styles.value}>{userInfo.name}</span>
      </div>
      <div className={styles.info_item}>
        <span className={styles.key}>邮箱</span>
        <span className={styles.value}>{userInfo.email}</span>
      </div>
      <div className={styles.info_item}>
        <span className={styles.key}>手机号</span>
        <span className={styles.value}>{userInfo.phone}</span>
      </div>
    </div>
  );
}

export default UserPage;
