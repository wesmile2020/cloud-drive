import { Button,Dropdown, MenuProps } from 'antd';
import { DownOutlined, LogoutOutlined, MailOutlined } from '@ant-design/icons';

import styles from './LoginUser.module.css';
import { useUserInfo } from '@/hooks/useUserInfo';
import { logout } from '@/services/api';
import { useNavigate } from 'react-router';

function LoginUser() {
  const [userInfo, fetchUserOrLogin] = useUserInfo();
  const navigate = useNavigate();

  const items: MenuProps['items'] = [
    {
      key: '1',
      label: '注销',
      icon: <LogoutOutlined />,
      onClick() {
        logout().then(() => {
          navigate('/login');
        });
      }
    },
  ];

  return (
    <div className={styles.avatar_container}>
      <div className={styles.logo}>
        <img src="cloud.svg"/>
      </div>
      {userInfo ? (
        <div className={styles.user_info}>
          <Dropdown menu={{ items }}>
            <div className={styles.user_name}>
              <div className={styles.name}>
                {userInfo.name}
              </div>
              <DownOutlined />
            </div>
          </Dropdown>
          <div className={styles.user_email}>
            <MailOutlined className={styles.email_icon}/>
            {userInfo.email}
          </div>
        </div>
      ) : (
        <Button onClick={() => fetchUserOrLogin(true)}
          className={styles.login_button}
          shape='circle'
          variant='text'
          color='primary'
          size='large'>
          登录
        </Button>
      )}
    </div>
  );
}

export default LoginUser;
