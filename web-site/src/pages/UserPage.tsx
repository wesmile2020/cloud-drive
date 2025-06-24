import { useEffect, useRef, useState } from 'react';
import { Button, Drawer } from 'antd';
import { useUserInfo } from '@/hooks/useUserInfo';
import EditUserInfo from '@/components/EditUserInfo';
import EditPassword from '@/components/EditPassword';

import styles from './UserPage.module.css';

const enum DrawerType {
  editPassword,
  editInfo,
}

function UserPage() {
  const [userInfo, fetchUserInfo] = useUserInfo();

  const [drawerOpen, setDrawerOpen] = useState(false);
  const [drawerType, setDrawerType] = useState<DrawerType>(DrawerType.editInfo);
  const userPageRef = useRef<HTMLDivElement>(null);

  function handleEditInfo() {
    setDrawerType(DrawerType.editInfo);
    setDrawerOpen(true);
  }

  function handleEditPassword() {
    setDrawerType(DrawerType.editPassword);
    setDrawerOpen(true);
  }

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
          onClick={handleEditInfo}>
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
      <div className={styles.info_item}>
        <Button variant='link'
          color='primary'
          onClick={handleEditPassword}>
          修改密码
        </Button>
      </div>
      <Drawer open={drawerOpen}
        getContainer={false}
        title={drawerType === DrawerType.editInfo ? '编辑信息' : '修改密码'}
        onClose={() => setDrawerOpen(false)}>
        {
          drawerType === DrawerType.editInfo ? <EditUserInfo /> : <EditPassword />
        }
      </Drawer>
    </div>
  );
}

export default UserPage;
