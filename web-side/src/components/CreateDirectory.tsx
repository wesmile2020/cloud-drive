import React from 'react';
import { CheckCircleFilled, CloseCircleFilled, FolderOutlined } from '@ant-design/icons';
import { Form, Input, message, Modal, Select, Tag } from 'antd';
import { createDirectory, FileTreeResponse } from '@/services/api';
import { calculatePublic } from '@/utils/utils';
import { Permission } from '@/config/enums';

interface Props {
  open: boolean;
  onClose: () => void;
  directoryTree: FileTreeResponse['tree'] | null;
  afterCreate?: () => void;
}

function CreateDirectory(props: Props) {
  const [form] = Form.useForm();
  const [permission, setPermission] = React.useState<Permission>(Permission.inherit);

  async function onFinish() {
    if (!props.directoryTree) {
      message.error('请先选择一个目录');
      return;
    }
    const values = form.getFieldsValue();
    await createDirectory({
      name: values.directory,
      parentId: Number(props.directoryTree.id),
      permission: values.permission,
    });
    message.success('创建文件夹成功');
    props.onClose();
    props.afterCreate?.();
  }

  const isPublic = React.useMemo(() => {
    if (!props.directoryTree) {
      return false;
    }
    return calculatePublic(props.directoryTree.public, permission);
  }, [props.directoryTree, permission]);

  return (
    <Modal title="新建文件夹"
      open={props.open}
      onCancel={props.onClose}
      okText='确认'
      cancelText='取消'
      onOk={() => form.submit()}>
      <Form onFinish={onFinish}
        onValuesChange={(values) => setPermission(values.permission)}
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
        <Form.Item>
          <span>公开：</span>
          <Tag icon={isPublic ? <CheckCircleFilled /> : <CloseCircleFilled />}
            color={isPublic ? 'success' : 'gold'}>
            {isPublic ? '是' : '否'}
          </Tag>
        </Form.Item>
      </Form>
    </Modal>
  );
}

export default CreateDirectory;
