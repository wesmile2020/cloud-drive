import { useMemo } from 'react';
import { Form, Input, Button, message, Flex, Checkbox } from 'antd';
import { useNavigate, Link, useSearchParams } from 'react-router';
import { UserOutlined, LockOutlined, LoginOutlined } from '@ant-design/icons';
import CardWrapper from '@/components/CardWrapper';
import LinkWrapper from '@/components/LinkWrapper';

import { login } from '@/services/api';


interface LoginParams {
  account: string;
  password: string;
  remember: boolean;
}

const ACCOUNT_KEY = 'user_account';
const PASSWORD_KEY = 'user_password';
const REMEMBER_KEY = 'user_remember';

function LoginPage() {
  const navigate = useNavigate();
  const [params] = useSearchParams();

  const registerUrl = useMemo(() => {
    const redirect = params.get('redirect');
    if (redirect) {
      return `/register?redirect=${encodeURIComponent(redirect)}`;
    }
    return `/register`;
  }, [params]);

  const onFinish = async (values: LoginParams) => {
    await login(values.account, values.password);
    if (values.remember) {
      localStorage.setItem(ACCOUNT_KEY, values.account);
      localStorage.setItem(PASSWORD_KEY, values.password);
      localStorage.setItem(REMEMBER_KEY, 'true');
    } else {
      localStorage.removeItem(ACCOUNT_KEY);
      localStorage.removeItem(PASSWORD_KEY);
      localStorage.removeItem(REMEMBER_KEY);
    }
    message.success('登录成功');
    
    const redirect = params.get('redirect') ?? '/';
    navigate(redirect);
  };

  return (
    <CardWrapper title='云盘登录'>
      <Form
        onFinish={onFinish}
        layout="vertical"
        initialValues={{
          account: localStorage.getItem(ACCOUNT_KEY) ?? '',
          password: localStorage.getItem(PASSWORD_KEY) ?? '',
          remember: localStorage.getItem(REMEMBER_KEY) === 'true',
        }}
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
          <Flex justify="space-between" align="center">
            <Form.Item name="remember"
              valuePropName="checked"
              noStyle
            >
              <Checkbox>记住登录</Checkbox>
            </Form.Item>
            <Link to='/retrieve-password'>
              忘记密码
            </Link>
          </Flex>
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
    </CardWrapper>
  );
}

export default LoginPage;
