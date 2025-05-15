import { message } from 'antd';
import { FileTreeResponse, updateDirectory } from '@/services/api';
import BaseDirectory, { FormValues } from './BaseDirectory';

interface Props {
  id: number;
  open: boolean;
  onClose: () => void;
  values: FormValues;
  directoryTree: FileTreeResponse['tree'];
  afterUpdate?: () => void;
}

function UpdateDirectory(props: Props) {
  async function onFinish(values: FormValues) {
    await updateDirectory(props.id, {
      name: values.directory,
      permission: values.permission,
    });
    message.success('修改文件夹成功');
    props.onClose();
    props.afterUpdate?.();
  }

  return (
    <BaseDirectory title="修改文件夹"
      open={props.open}
      onClose={props.onClose}
      directoryTree={props.directoryTree}
      onFinish={onFinish}
      values={props.values}
    />
  );
}

export default UpdateDirectory;
