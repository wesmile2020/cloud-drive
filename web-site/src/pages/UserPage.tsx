import { useEffect, useRef, useState } from 'react';
import { Button, Drawer } from 'antd';
import { useUserInfo } from '@/hooks/useUserInfo';
import EditUserInfo from '@/components/EditUserInfo';

import styles from './UserPage.module.css';

function UserPage() {
  const [userInfo, fetchUserInfo] = useUserInfo();

  const [drawerOpen, setDrawerOpen] = useState(false);
  const userPageRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!userInfo) {
      fetchUserInfo(true);
    }
  }, [userInfo, fetchUserInfo]);

  return userInfo && (
    <div className={styles.user_page}
      ref={userPageRef}>
      <h2 className={styles.title}>
        用户信息
        <Button variant='link'
          color='primary'
          onClick={() => setDrawerOpen(true)}>
          编辑
        </Button>
      </h2>
      <div className={styles.info_item}>
        <span className={styles.key}>用户名</span>
        <span className={styles.value}>{userInfo.name}</span>
      </div>
      <div className={styles.info_item}>
        <span className={styles.key}>邮&nbsp;&nbsp;&nbsp;箱</span>
        <span className={styles.value}>{userInfo.email}</span>
      </div>
      <div className={styles.info_item}>
        <span className={styles.key}>手机号</span>
        <span className={styles.value}>{userInfo.phone}</span>
      </div>
      <Drawer open={drawerOpen}
        getContainer={false}
        title='编辑'
        onClose={() => setDrawerOpen(false)}>
        <EditUserInfo onSuccess={() => {
            setDrawerOpen(false);
          }}
        />
      </Drawer>
    </div>
  );
}

export default UserPage;
