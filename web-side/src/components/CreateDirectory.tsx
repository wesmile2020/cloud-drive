import { message } from 'antd';
import { FileTreeResponse, createDirectory } from '@/services/api';
import BaseDirectory, { FormValues } from './BaseDirectory';

interface Props {
  open: boolean;
  onClose: () => void;
  directoryTree: FileTreeResponse['tree'];
  afterCreate?: () => void;
}

function CreateDirectory(props: Props) {
  async function onFinish(values: FormValues) {
    if (!props.directoryTree) {
      message.error('请先选择一个目录');
      return;
    }
  
    await createDirectory({
      name: values.directory,
      parentId: Number(props.directoryTree.id),
      permission: values.permission,
    });
    message.success('创建文件夹成功');
    props.onClose();
    props.afterCreate?.();
  }

  return (
    <BaseDirectory title="新建文件夹"
      open={props.open}
      onClose={props.onClose}
      directoryTree={props.directoryTree}
      onFinish={onFinish}
    />
  );
}

export default CreateDirectory;
