import { FolderOutlined } from '@ant-design/icons';
import {  Form, Input, message, Modal, Select } from 'antd';
import { createDirectory, Permission } from '@/services/api';

interface Props {
  open: boolean;
  onClose: () => void;
  directoryId: number;
  afterCreate?: () => void;
}

function CreateDirectory(props: Props) {
  const [form] = Form.useForm();
  async function onFinish() {
    const values = form.getFieldsValue();
    await createDirectory({
      name: values.directory,
      parentId: Number(props.directoryId),
      permission: values.permission,
    });
    message.success('创建文件夹成功');
    props.onClose();
    props.afterCreate?.();
  }

  return (
    <Modal title="新建文件夹"
      open={props.open}
      onCancel={props.onClose}
      okText='确认'
      cancelText='取消'
      onOk={() => form.submit()}>
      <Form onFinish={onFinish}
        form={form}
        initialValues={{
          permission: Permission.inherit,
        }}>
        <Form.Item name="directory"
          rules={[{ required: true, message: '文件夹名称不能为空' }]}>
          <Input placeholder="请输入文件夹名称"
            autoFocus
            prefix={<FolderOutlined />}/>
        </Form.Item>
        <Form.Item>
          <span>权限：</span>
          <Form.Item name="permission"
            noStyle>
            <Select style={{ width: 120 }}
              options={[
                { label: '私有', value: Permission.private },
                { label: '继承', value: Permission.inherit },
                { label: '公开', value: Permission.public }
              ]}
            />      
          </Form.Item>
        </Form.Item>
      </Form>
    </Modal>
  );
}

export default CreateDirectory;
