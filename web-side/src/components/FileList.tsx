import React from 'react';
import { Avatar, Table, TableColumnType, Tag, Button, Dropdown, MenuProps, Popconfirm } from 'antd';
import { DeleteOutlined, EditOutlined, FileFilled, FolderFilled, MoreOutlined } from '@ant-design/icons';
import { format } from 'date-fns';
import { Link } from 'react-router';
import { FileTreeResponse } from '@/services/api';
import { formatSize } from '@/utils/utils';
import { useUserInfo } from '@/hooks/useUserInfo';
import { Permission } from '@/config/enums';
import UpdateDirectory from './UpdateDirectory';

import styles from './FileList.module.css';

interface Props {
  files: FileTreeResponse['files'];
  directoryTree: FileTreeResponse['tree'];
  loading: boolean;
  afterUpdate?: () => void;
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
  const [operateRecord, setOperateRecord] = React.useState<FileTreeResponse['files'][0] | null>(null);

  function getOperates(record: FileTreeResponse['files'][0]): MenuProps['items'] {
    if (record.isDirectory) {
      return [
        {
          key: '1',
          label: '编辑',
          icon: <EditOutlined />,
          disabled: userInfo?.id !== record.user.id,
          onClick: () => {
            setOperateRecord(record);
          }
        },
        {
          key: '2',
          label: (
            <Popconfirm title="确认删除？">
              删除
            </Popconfirm>
          ),
          icon: <DeleteOutlined />,
          disabled: userInfo?.id !== record.user.id,
          onClick: () => {
            console.log('删除', record);  
          }
        }
      ];
    } else {
      return [

      ];
    }
  }

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
    },
    {
      title: '操作',
      dataIndex: 'id',
      render: (_: number, record) => {
        return (
          <Dropdown menu={{ items: getOperates(record)}}
            trigger={['click']}
            disabled={record.isDirectory && record.user.id !== userInfo?.id}>
            <Button shape='circle'
              icon={<MoreOutlined />}
              type='primary'
            />
          </Dropdown>
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
  }, []);

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
      {
        operateRecord && (
          <UpdateDirectory id={operateRecord.id}
            open={Boolean(operateRecord)}
            directoryTree={props.directoryTree}
            onClose={() => setOperateRecord(null)}
            values={{
              directory: operateRecord.name,
              permission: operateRecord.permission,
            }}
            afterUpdate={props.afterUpdate}
          />
        )
      }
      
    </div>
  );
}

export default FileList;
