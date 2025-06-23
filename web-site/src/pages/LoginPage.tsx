import React from 'react';
import { Form, Input, Button, message } from 'antd';
import { useNavigate, Link, useSearchParams } from 'react-router';
import { UserOutlined, LockOutlined, LoginOutlined } from '@ant-design/icons';
import CardWrapper from '@/components/CardWrapper';
import LinkWrapper from '@/components/LinkWrapper';

import { login } from '@/services/api';


interface LoginParams {
  account: string;
  password: string;
}

function LoginPage() {
  const navigate = useNavigate();
  const [params] = useSearchParams();

  const registerUrl = React.useMemo(() => {
    const redirect = params.get('redirect');
    if (redirect) {
      return `/register?redirect=${encodeURIComponent(redirect)}`;
    }
    return `/register`;
  }, [params]);

  const onFinish = async (values: LoginParams) => {
    await login(values.account, values.password);
    message.success('登录成功');
    const redirect = params.get('redirect') ?? '/';
    navigate(redirect);
  };

  return (
    <CardWrapper title='云盘登录'>
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
            icon={<LoginOutlined />}
          >
            登录
          </Button>
        </Form.Item>
      </Form>
      <LinkWrapper>
        <span>还没有账号？</span>
        <Link to={registerUrl}>
          立即注册
        </Link>
      </LinkWrapper>
      <LinkWrapper>
        <Link to='/retrieve-password'
          style={{ fontSize: 12, color: '#1890ff' }}>
          忘记密码
        </Link>
      </LinkWrapper>
    </CardWrapper>
  );
}

export default LoginPage;
