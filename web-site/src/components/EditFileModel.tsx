import React from 'react';
import { Form, Modal, Input, Select, Tag } from 'antd';
import { CheckCircleFilled, CloseCircleFilled, FolderOutlined } from '@ant-design/icons';
import { FileTreeResponse } from '@/services/api';
import { calculatePublic } from '@/utils/utils';
import { Permission } from '@/config/enums';

export interface FormValues {
  directory: string;
  permission: number;
}

interface Props {
  title: string;
  open: boolean;
  onClose: () => void;
  directoryTree: FileTreeResponse['tree'];
  onFinish: (values: FormValues) => void;
  values?: FormValues;
  type: 'directory' | 'file';
}

function EditFileModel(props: Props) {
  const [form] = Form.useForm();
  const [permission, setPermission] = React.useState<number>(props.values?.permission ?? Permission.inherit);
  const isPublic = React.useMemo(() => {
    if (!props.directoryTree) {
      return false;
    }
    return calculatePublic(props.directoryTree.public, permission);
  }, [props.directoryTree, permission]);

  const tipText = props.type === 'file' ? '文件' : '文件夹';

  return (
    <Modal title={props.title}
      open={props.open}
      onCancel={props.onClose}
      okText='确认'
      cancelText='取消'
      onOk={form.submit}>
      <Form onFinish={props.onFinish}
        onValuesChange={(values) => {
          if (typeof values.permission === 'number') {
            setPermission(values.permission);
          }
        }}
        form={form}
        initialValues={{
          permission,
          directory: props.values?.directory ?? '',
        }}>
        <Form.Item name="directory"
          rules={[{ required: true, message: `${tipText}名称不能为空` }]}>
          <Input placeholder={`请输入${tipText}名称`}
            autoFocus
            prefix={<FolderOutlined />}
          />
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

export default EditFileModel;
