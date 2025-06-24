import { Form, Input, Button, message } from 'antd';
import { useUserInfo } from '@/hooks/useUserInfo';
import { CheckOutlined, MailOutlined, PhoneOutlined, UserOutlined } from '@ant-design/icons';
import { editUserInfo } from '@/services/api';

interface EditUserInfoParams {
  name: string;
  email: string;
  phone: string;
}

interface Props {
  onSuccess?: () => void;
}

function EditUserInfo(props: Props) {
  const [userInfo, fetchUserInfo] = useUserInfo();
  const [form] = Form.useForm();

  async function onFinish(params: EditUserInfoParams) {
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
      }}
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
