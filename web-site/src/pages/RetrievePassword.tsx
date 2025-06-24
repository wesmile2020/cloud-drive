import { startTransition, useActionState, useEffect, useRef, useState } from 'react';
import { Button, Form, Input, message } from 'antd';
import { LockOutlined, MailOutlined, SendOutlined, VerifiedOutlined } from '@ant-design/icons';
import { Link } from 'react-router';

import CardWrapper from '@/components/CardWrapper';
import LinkWrapper from '@/components/LinkWrapper';
import { getVerifyCode, retrievePassword } from '@/services/api';

interface FormParams {
  email: string;
  code: string;
  password: string;
  confirmPassword: string;
}

const COUNT_DOWN_TIME_KEY = 'countDownTime';

function RetrievePassword() {
  const [countDown, setCountDown] = useState(0);
  const intervalRef = useRef<number>(0);

  const [form] = Form.useForm();

  function startCountDown(rest: number) {
    if (intervalRef.current > 0) {
      window.clearInterval(intervalRef.current);
      intervalRef.current = 0;
    }
    setCountDown(rest);
   
    intervalRef.current = window.setInterval(() => {
      setCountDown(countDown => {
        if (countDown <= 1) {
          window.clearInterval(intervalRef.current);
          localStorage.removeItem(COUNT_DOWN_TIME_KEY);
          intervalRef.current = 0;
          return 0;
        }
        return countDown - 1;
      });
    }, 1000);
  }

  const [, sendVerifyCode, verifyCodeLoading] = useActionState(async () => {
    try {
      const { email } = await form.validateFields(['email']);
      await getVerifyCode(email);
      message.success('发送验证码成功');
      startCountDown(60);
      localStorage.setItem(COUNT_DOWN_TIME_KEY, Date.now().toString());
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
    } catch (error) {
      // 发送验证码失败
    }
  }, null);

  async function onFinish(params: FormParams) {
    await retrievePassword(params.email, params.code, params.password);
    form.resetFields();
    message.success('密码重置成功');
  }

  useEffect(() => {
    const time = localStorage.getItem(COUNT_DOWN_TIME_KEY);
    if (time) {
      const diff = Date.now() - parseInt(time);
      if (diff < 60 * 1000) {
        startCountDown(60 - Math.floor(diff / 1000));
      }
    }
  }, []);

  return (
    <CardWrapper title='找回密码'>
      <Form form={form}
        onFinish={onFinish}>
        <Form.Item
          name="email"
          rules={[{ required: true, type: 'email', message: '请输入邮箱' }]}
        >
          <Input placeholder="请输入邮箱"
            prefix={<MailOutlined />}
          />
        </Form.Item>
        <Form.Item
          name="code"
          rules={[{ required: true, message: '请输入验证码' }]}
        >
          <Input placeholder='请输入验证码'
            prefix={<VerifiedOutlined />}
            suffix={
              <Button
                variant="link"
                color="primary"
                loading={verifyCodeLoading}
                disabled={countDown > 0}
                onClick={() => startTransition(() => sendVerifyCode())}
              >
                {countDown > 0 ? `${countDown}s` : '获取验证码'}
              </Button>
            }
          />
        </Form.Item>
        <Form.Item
          name="password"
          rules={[{ required: true, message: '请输入新密码最少6位', min: 6 }]}
        >
          <Input.Password
            placeholder="请输入新密码"
            prefix={<LockOutlined />}
          />
        </Form.Item>
        <Form.Item
          name="confirmPassword"
          rules={[
            { required: true, message: '请确认新密码' },
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
            placeholder="请确认新密码"
            prefix={<LockOutlined />}
          />
        </Form.Item>
        <Form.Item>
          <Button
            type="primary"
            htmlType="submit"
            block
            icon={<SendOutlined />}
          >
            确认
          </Button>
        </Form.Item>
      </Form>
      <LinkWrapper>
        <span>已完成？</span>
        <Link to="/login">去登陆</Link>
      </LinkWrapper>
    </CardWrapper>
  );
}


export default RetrievePassword;
