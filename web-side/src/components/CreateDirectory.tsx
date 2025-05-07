import { FolderOutlined } from '@ant-design/icons';
import { Checkbox, Form, Input, message, Modal } from 'antd';
import { createDirectory } from '@/services/api';

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
      public: values.public ?? false
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
        form={form}>
        <Form.Item name="directory"
          rules={[{ required: true, message: '文件夹名称不能为空' }]}>
          <Input placeholder="请输入文件夹名称"
            autoFocus
            prefix={<FolderOutlined />}/>
        </Form.Item>
        <Form.Item name="public"
          valuePropName='checked'>
          <Checkbox>
            是否公开
          </Checkbox>
        </Form.Item>
      </Form>
    </Modal>
  );
}

export default CreateDirectory;
