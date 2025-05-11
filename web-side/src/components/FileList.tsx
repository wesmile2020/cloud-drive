import React from 'react';
import { Avatar, Table, TableColumnType, Tag } from 'antd';
import { format } from 'date-fns';
import { Link } from 'react-router';
import { FileTreeResponse } from '@/services/api';

import styles from './FileList.module.css';
import { formatSize } from '@/utils/utils';
import { FileFilled, FolderFilled } from '@ant-design/icons';
import { useUserInfo } from '@/hooks/useUserInfo';
import { Permission } from '@/config/enums';

interface Props {
  files: FileTreeResponse['files'];
  loading: boolean;
}

const PermissionText = {
  [Permission.private]: '私有',
  [Permission.inherit]: '继承',
  [Permission.public]: '公开',
};

function FileList(props: Props) {
  const { files, loading } = props;
  const domRef = React.useRef<HTMLDivElement>(null);
  const [scrollY, setScrollY] = React.useState(0);
  const [userInfo] = useUserInfo();

  const columns: TableColumnType<FileTreeResponse['files'][0]>[] = [
    {
      title: '名称',
      dataIndex: 'name',
      render(name: string, record) {
        if (record.isDirectory) {
          return (
            <Link to={`/home/${record.id}`}>
              <FolderFilled /> {name}
            </Link>
          );
        } else {
          return (
            <>
              <FileFilled /> {name}
            </>
          );
        }
      }
    },
    {
      title: '日期',
      dataIndex:'timestamp',
      render(timestamp: number) {
        return format(timestamp * 1000, 'yyyy-MM-dd HH:mm:ss');
      }
    },
    {
      title: '用户',
      dataIndex: 'user',
      render(user: FileTreeResponse['files'][0]['user']) {
        return (
          <Avatar className={styles.avatar}
            size="small">
            {user.name.slice(0, 1)}
          </Avatar>
        );
      }
    },
    {
      title: '大小',
      dataIndex: 'size',
      render: (size: number, record) => {
        if (record.isDirectory) {
          return '-' 
        }
        return formatSize(size);
      }
    },
    {
      title: '标签',
      dataIndex: 'permission',
      render: (permission: Permission, record) => {
        return (
          <>
            <Tag color={record.public ? 'green' : 'gold'}>
              {PermissionText[permission]}
            </Tag>
            {userInfo?.id === record.user.id && <Tag color='blue'>自有</Tag>}
          </>
        );
      }
    }
  ];

  function initScroll() {
    if (!domRef.current) {
      return;
    }
    requestAnimationFrame(() => {
      const { clientHeight } = domRef.current!;
      setScrollY(clientHeight - 55);
    });
  }

  React.useEffect(() => {
    initScroll();
    window.addEventListener('resize', initScroll);
    return () => {
      window.removeEventListener('resize', initScroll)
    };
  }, [])

  return (
    <div className={styles.file_list}
      ref={domRef}>
      <Table className={styles.table}
        pagination={false}
        loading={loading}
        columns={columns}
        dataSource={files}
        rowKey="id"
        scroll={{ y: scrollY }}
        virtual={true}
      />
    </div>
    
  );
}

export default FileList;
