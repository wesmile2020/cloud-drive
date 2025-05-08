import { Link, useLocation } from 'react-router';
import { Outlet } from 'react-router';
import React from 'react';
import { HomeOutlined, UserOutlined } from '@ant-design/icons';
import { Layout, Menu } from 'antd';
import UserAvatar from '@/components/UserAvatar';

import styles from './HomeLayout.module.css';

const menus = [
  {
    key: '/home',
    label: <Link to="/home/0">首页</Link>,
    icon: <HomeOutlined />
  },
  {
    key: '/user',
    label: <Link to="/user">用户中心</Link>,
    icon: <UserOutlined />
  },
]

function HomeLayout() {
  const location = useLocation();
  const [selectedKeys, setSelectedKeys] = React.useState<string[]>([])
  React.useEffect(() => {
    if (/\/home/.test(location.pathname)) {
      setSelectedKeys(['/home']);
      return;
    }
    for (const item of menus) {
      if (item.key === location.pathname) {
        setSelectedKeys([item.key]);
        return;
      }
    }

  }, [location.pathname]);

  return (
    <Layout className={styles.home_layout}>
      <Layout.Sider theme='light'
        collapsible>
        <UserAvatar />
        <Menu theme='light'
          mode='inline'
          items={menus}
          selectedKeys={selectedKeys}  
        />
         
      </Layout.Sider>
      <Layout.Content className={styles.layout_content}>
        <Outlet />
      </Layout.Content>

    </Layout>
  );
}

export default HomeLayout;
