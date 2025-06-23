import { Button, Form, Input } from 'antd';
import { LockOutlined, MailOutlined, SendOutlined, VerifiedOutlined } from '@ant-design/icons';
import CardWrapper from '@/components/CardWrapper';
import styles from './RetrievePassword.module.css';
import LinkWrapper from '@/components/LinkWrapper';
import { Link } from 'react-router';

interface MailProps {
  onChange?: (value: string) => void;
  value?: string;
}

function MailInput(props: MailProps) {
  
  return (
    <div className={styles.mail_input}>
      <Input
        placeholder="请输入邮箱"
        prefix={<MailOutlined />}
        value={props.value}
        onChange={(e) => props.onChange?.(e.target.value)}
      />
      <Button disabled={!props.value}
        icon={<SendOutlined />}
        type='primary'>
          验证码
      </Button>
    </div>
  );
}

function RetrievePassword() {
  return (
    <CardWrapper title='找回密码'>
      <Form>
        <Form.Item
          name="account"
          rules={[{ required: true, message: '请输入邮箱' }]}
        >
          <MailInput />
        </Form.Item>
        <Form.Item
          name="code"
          rules={[{ required: true, message: '请输入验证码' }]}
        >
          <Input placeholder='请输入验证码'
            prefix={<VerifiedOutlined />}
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
