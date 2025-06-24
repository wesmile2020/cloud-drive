import { Form, Input, Button, message } from 'antd';
import { useUserInfo } from '@/hooks/useUserInfo';
import { CheckOutlined, LockOutlined } from '@ant-design/icons';
import { updatePassword, logout } from '@/services/api';
import { useNavigate } from 'react-router';

interface ModifyPasswordParams {
  oldPassword: string;
  newPassword: string;
}

function EditPassword() {
  const [userInfo] = useUserInfo();
  const [form] = Form.useForm();
  const navigate = useNavigate();

  async function onFinish(params: ModifyPasswordParams) {
    await updatePassword(params.oldPassword, params.newPassword);
    message.info('修改成功, 请重新登录');
    logout().then(() => {
      navigate('/login');
    });
  }

  return userInfo && (
    <Form form={form}
      onFinish={onFinish}
    >
      <Form.Item
        name="oldPassword"
        rules={[{ required: true, message: '请输入旧密码' }, { min: 6, message: '旧密码长度不能小于6位' }]}
      >
        <Input.Password
          placeholder="请输入旧密码"
          prefix={<LockOutlined />}
        />
      </Form.Item>
      <Form.Item
        name="newPassword"
        rules={[{ required: true, message: '请输入新密码' }, { min: 6, message: '新密码长度不能小于6位' }]}
      >
        <Input.Password
          placeholder="请输入新密码"
          prefix={<LockOutlined />}
        />
      </Form.Item>
      <Form.Item
        name="confirmPassword"
        rules={[
          {
            required: true,
            message: '请确认新密码',
          },
          ({ getFieldValue }) => ({
            validator(_, value) {
              if (!value || getFieldValue('newPassword') === value) {
                return Promise.resolve();
              }
              return Promise.reject(new Error('两次输入密码不一致!'));
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
        <Button type="primary"
          htmlType="submit"
          block
          icon={<CheckOutlined />}
        >
          提交
        </Button>
      </Form.Item>
    </Form>
  )
}

export default EditPassword;
