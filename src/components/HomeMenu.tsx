import React from 'react';
import { FolderAddOutlined, MoreOutlined } from '@ant-design/icons';
import { FloatButton, message } from 'antd';
import CreateDirectory from './CreateDirectory';
import { useUserInfo } from '@/hooks/useUserInfo';

interface Props {
  directoryId: number;
}

function HomeMenu(props: Props) {
  const [open, setOpen] = React.useState(false);
  const [userInfo] = useUserInfo();
  function triggerCreateDirectory() {
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
        onClose={() => setOpen(false)}/>
    </>
  );
}

export default HomeMenu;
