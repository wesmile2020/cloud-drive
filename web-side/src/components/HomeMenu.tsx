import React from 'react';
import { FolderAddOutlined, MoreOutlined } from '@ant-design/icons';
import { FloatButton, message } from 'antd';
import CreateDirectory from './CreateDirectory';
import { useUserInfo } from '@/hooks/useUserInfo';

interface Props {
  directoryId: number;
  afterCreate?: () => void;
}

function HomeMenu(props: Props) {
  const [open, setOpen] = React.useState(false);
  const [userInfo] = useUserInfo();
  async function triggerCreateDirectory() {
    if (!userInfo) {
      message.error('请先登录');
      return; 
    }
    setOpen(true);
  }

  return (
    <>
      <FloatButton.Group trigger='hover'
        type='primary'
        icon={<MoreOutlined />}>
        <FloatButton type='primary'
          icon={<FolderAddOutlined />}
          tooltip='新建文件夹'
          onClick={triggerCreateDirectory}
        />
      </FloatButton.Group>
      <CreateDirectory open={open}
        directoryId={props.directoryId}
        afterCreate={props.afterCreate}
        onClose={() => setOpen(false)}/>
    </>
  );
}

export default HomeMenu;
