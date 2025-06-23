import { message } from 'antd';
import { FileTreeResponse, createDirectory } from '@/services/api';
import EditFileModel, { FormValues } from './EditFileModel';

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
    <EditFileModel title="新建文件夹"
      type="directory"
      open={props.open}
      onClose={props.onClose}
      directoryTree={props.directoryTree}
      onFinish={onFinish}
    />
  );
}

export default CreateDirectory;
