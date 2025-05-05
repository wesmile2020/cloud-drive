import { Form, Input, Button, message } from 'antd';
import { useNavigate, Link } from 'react-router';
import styles from './LoginPage.module.css';
import { UserOutlined, LockOutlined, LoginOutlined } from '@ant-design/icons';

import { login } from '@/services/api';

interface LoginParams {
  account: string;
  password: string;
}

function LoginPage() {

  const onFinish = async (values: LoginParams) => {
    await login(values.account, values.password);
    message.success('登录成功');
  };

  return (
    <div className={styles.login_container}>
      <div className={styles.login_card}>
        <h1 className={styles.login_title}>云盘登录</h1>
        <Form
          onFinish={onFinish}
          layout="vertical"
        >
          <Form.Item
            name="account"
            rules={[{ required: true, message: '请输入手机号或邮箱' }]}
          >
            <Input
              placeholder="请输入手机号或邮箱"
              prefix={<UserOutlined />}
            />
          </Form.Item>
          <Form.Item
            name="password"
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password
              placeholder="请输入密码"
              prefix={<LockOutlined />}
            />
          </Form.Item>
          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              block
              className={styles.login_button}
              icon={<LoginOutlined />}
            >
              登录
            </Button>
          </Form.Item>
        </Form>
        <div className={styles.login_link}>
          <span>还没有账号？</span>
          <Link to="/register" className={styles.login_link_text}>
            立即注册
          </Link>
        </div>
      </div>
    </div>
  );
}

export default LoginPage;
