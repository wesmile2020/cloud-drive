import { Form, Input, Button, Switch, message } from 'antd';
import { useUserInfo } from '@/hooks/useUserInfo';
import { CheckOutlined, LockOutlined, MailOutlined, PhoneOutlined, UserOutlined } from '@ant-design/icons';
import { editUserInfo, updatePassword } from '@/services/api';

interface BaseInfoParams {
  name: string;
  email: string;
  phone: string;
}

interface ModifyPasswordParams {
  modifyPassword: true;
  oldPassword: string;
  newPassword: string;
}

interface NoModifyPasswordParams {
  modifyPassword: false;
}

type EditUserInfoParams = BaseInfoParams & (ModifyPasswordParams | NoModifyPasswordParams);

interface Props {
  onSuccess?: () => void;
}

function EditUserInfo(props: Props) {
  const [userInfo, fetchUserInfo] = useUserInfo();
  const [form] = Form.useForm();
  const modifyPassword = Form.useWatch('modifyPassword', form);

  async function onFinish(params: EditUserInfoParams) {
    if (params.modifyPassword) {
      await updatePassword(params.oldPassword, params.newPassword);
    }
    await editUserInfo({
      name: params.name,
      email: params.email,
      phone: params.phone,
    });
    fetchUserInfo(false, true);
    message.info('修改成功');
    props.onSuccess?.();
  }

  return userInfo && (
    <Form form={form}
      onFinish={onFinish}
      initialValues={{
        name: userInfo.name,
        email: userInfo.email,
        phone: userInfo.phone,
        modifyPassword: false,
      }}>
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
        name="modifyPassword"
        label="修改密码？"
        hasFeedback
      >
        <Switch />
      </Form.Item>
      {
        modifyPassword && (
          <>
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
          </>
          
        )
      }
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

export default EditUserInfo;
