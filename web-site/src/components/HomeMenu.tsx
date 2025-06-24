import { useState } from 'react';
import { FolderAddOutlined, MoreOutlined } from '@ant-design/icons';
import { FloatButton, message } from 'antd';
import CreateDirectory from './CreateDirectory';
import { useUserInfo } from '@/hooks/useUserInfo';
import { FileTreeResponse } from '@/services/api';

interface Props {
  directoryTree: FileTreeResponse['tree'];
  afterCreate?: () => void;
}

function HomeMenu(props: Props) {
  const { directoryTree } = props;

  const [open, setOpen] = useState(false);
  const [userInfo] = useUserInfo();
  async function triggerCreateDirectory() {
    if (!userInfo) {
      message.error('请先登录');
      return; 
    }
    if (!directoryTree || (directoryTree.id !== 0 && directoryTree.userId !== userInfo.id)) {
      message.error('您不能在该目录下新建文件夹');
      return; 
    }
    setOpen(true);
  }

  return (
    <>
      <FloatButton.Group trigger='click'
        type='primary'
        icon={<MoreOutlined />}>
        <FloatButton type='primary'
          icon={<FolderAddOutlined />}
          tooltip='新建文件夹'
          onClick={triggerCreateDirectory}
        />
      </FloatButton.Group>
      <CreateDirectory open={open}
        directoryTree={props.directoryTree}
        afterCreate={props.afterCreate}
        onClose={() => setOpen(false)}
      />
    </>
  );
}

export default HomeMenu;
