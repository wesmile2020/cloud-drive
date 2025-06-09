import { message } from 'antd';
import { FileTreeResponse, updateDirectory, updateFile } from '@/services/api';
import EditFileModel, { FormValues } from './EditFileModel';

interface Props {
  id: number;
  open: boolean;
  onClose: () => void;
  values: FormValues;
  directoryTree: FileTreeResponse['tree'];
  afterUpdate?: () => void;
  type: 'directory' | 'file';
}

function UpdateFileModel(props: Props) {
  async function onFinish(values: FormValues) {
    if (props.type === 'file') {
      await updateFile(props.id, {
        name: values.directory,
        permission: values.permission,
      })
    } else {
      await updateDirectory(props.id, {
        name: values.directory,
        permission: values.permission,
      });
    }
      
    message.success('修改文件夹成功');
    props.onClose();
    props.afterUpdate?.();
  }

  return (
    <EditFileModel title={`编辑${props.type === 'file'? '文件' : '文件夹'}`}
      open={props.open}
      onClose={props.onClose}
      directoryTree={props.directoryTree}
      onFinish={onFinish}
      values={props.values}
      type={props.type}
    />
  );
}

export default UpdateFileModel;
