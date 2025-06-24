import { useMemo } from 'react';
import { Form, Input, Button, message } from 'antd';
import { Link, useSearchParams } from 'react-router';
import { UserOutlined, LockOutlined, PhoneOutlined, MailOutlined } from '@ant-design/icons';

import CardWrapper from '@/components/CardWrapper';
import LinkWrapper from '@/components/LinkWrapper';
import { register, RegisterParams } from '@/services/api'

function RegisterPage() {
  const [form] = Form.useForm();
  const [params] = useSearchParams();

  const loginUrl = useMemo(() => {
    const redirect = params.get('redirect');
    if (redirect) {
      return `/login?redirect=${encodeURIComponent(redirect)}`;
    }
    return `/login`;
  }, [params]);

  const onFinish = async (values: RegisterParams) => {
    // 验证两次输入的密码是否一致
    await register(values);
    message.success('注册成功');
    form.resetFields();
  };

  return (
    <CardWrapper title='云盘注册'>
      <Form
        form={form}
        onFinish={onFinish}
        layout="vertical"
      >
        <Form.Item
          name="name"
          rules={[{ required: true, message: '请输入用户名' }]}
        >
          <Input
            placeholder="请输入用户名"
            prefix={<UserOutlined />}
          />
        </Form.Item>
        <Form.Item
          name="phone"
          rules={[{ required: true, message: '请输入手机号' }, { pattern: /^1[3-9]\d{9}$/, message: '请输入有效的手机号' }]}
        >
          <Input
            placeholder="请输入手机号"
            prefix={<PhoneOutlined />}
          />
        </Form.Item>
        <Form.Item
          name="email"
          rules={[{ required: true, message: '请输入邮箱' }, { type: 'email', message: '请输入有效的邮箱' }]}
        >
          <Input
            placeholder="请输入邮箱"
            prefix={<MailOutlined />}
          />
        </Form.Item>
        <Form.Item
          name="password"
          rules={[{ required: true, message: '密码最少6位', min: 6 }]}
        >
          <Input.Password
            placeholder="请输入密码"
            prefix={<LockOutlined />}
          />
        </Form.Item>
        <Form.Item
          name="confirmPassword"
          dependencies={['password']}
          rules={[
            { required: true, message: '请再次输入密码', min: 6 },
            ({ getFieldValue }) => ({
              validator(_, value) {
                if (!value || getFieldValue('password') === value) {
                  return Promise.resolve();
                }
                return Promise.reject();
              },
            }),
          ]}
        >
          <Input.Password
            placeholder="请再次输入密码"
            prefix={<LockOutlined />}
          />
        </Form.Item>
        <Form.Item>
          <Button
            type="primary"
            htmlType="submit"
            block
          >
            注册
          </Button>
        </Form.Item>
      </Form>
      <LinkWrapper>
        <span>已有账号？</span>
        <Link to={loginUrl}>立即登录</Link>
      </LinkWrapper>
    </CardWrapper>
  );
};

export default RegisterPage;
