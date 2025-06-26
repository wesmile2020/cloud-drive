import { useState, useRef, useEffect, useMemo } from 'react';
import { Avatar, Table, TableColumnType, Tag, Button, Dropdown, MenuProps, Tooltip } from 'antd';
import { DeleteOutlined, DownloadOutlined, EditOutlined, FolderFilled, MoreOutlined } from '@ant-design/icons';
import { format } from 'date-fns';
import { Link } from 'react-router';
import { FileTreeResponse } from '@/services/api';
import { downloadFile, formatSize } from '@/utils/utils';
import { useUserInfo } from '@/hooks/useUserInfo';
import { Permission } from '@/config/enums';
import UpdateFileModel from './UpdateFileModel';
import DeleteMenuItem from './DeleteMenuItem';
import FilePrefix from './FileName';

import styles from './FileList.module.css';
import FileUploadProgress from './FileUploadProgress';
import { useUpload } from '@/hooks/useUpload';
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

type FileItem = FileTreeResponse['files'][0];

interface RecordItem extends FileItem {
  key: string;
}

function FileList(props: Props) {
  const { files, loading } = props;
  const domRef = useRef<HTMLDivElement>(null);
  const [scrollY, setScrollY] = useState(0);
  const [userInfo] = useUserInfo();
  const [editRecord, setEditRecord] = useState<RecordItem | null>(null);
  const [openRecordId, setOpenRecordId] = useState<string | null>(null);
  const [uploadList] = useUpload();

  const wrappedFiles = useMemo(() => {
    return files.map(file => ({
      ...file,
      key: file.isDirectory ? `d_${file.id}` : `f_${file.id}`,
    }));
  }, [files]);

  function getOperates(record: RecordItem): MenuProps['items'] {
    const operates = [
      {
        key: 'edit',
        label: '编辑',
        icon: <EditOutlined />,
        disabled: userInfo?.id !== record.user.id,
        onClick() {
          setEditRecord(record);
        }
      },
      {
        key: 'delete',
        danger: true,
        label: '删除',
        icon: <DeleteOutlined />,
        disabled: userInfo?.id !== record.user.id,
      }
    ];

    if (!record.isDirectory) {
      operates.unshift({
        key: 'download',
        label: '下载',
        icon: <DownloadOutlined />,
        disabled: false,
        onClick() {
          downloadFile(`/api/file/download/${record.id}`, record.name);
        }
      });
    }
    return operates;
  }

  const columns: TableColumnType<RecordItem>[] = [
    {
      title: '名称',
      dataIndex: 'name',
      ellipsis: {
        showTitle: false,
      },
      render(name: string, record) {
        return (
          <div className={styles.name}>
            {
              record.isDirectory ? (
                <Tooltip title={name}
                  placement='topLeft'>
                  <Link to={`/home/${record.id}`}>
                    <FolderFilled />&nbsp;{name}
                  </Link>
                </Tooltip>
              ) : (
                <FilePrefix name={name}/>
              )
            }
          </div>
        )
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
          <Dropdown menu={{
              items: getOperates(record),
              onClick() {
                setOpenRecordId(record.key);
              },
              _internalRenderMenuItem(originNode, menuItemProps) {
                if (menuItemProps.eventKey === 'delete') {
                  return (
                    <DeleteMenuItem record={record}
                      disabled={menuItemProps.disabled}
                      className={menuItemProps.className}
                      afterDelete={() => {
                        setOpenRecordId(null);
                        props.afterUpdate?.();
                      }}
                    />
                  );
                }
                return originNode;
              }
            }}
            open={openRecordId === record.key}
            onOpenChange={(open, info) => {
              if (info.source === 'trigger') {
                setOpenRecordId(open ? record.key : null);
              }
            }}
            trigger={['click']}
            disabled={record.isDirectory && record.user.id !== userInfo?.id}
          >
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
    setTimeout(() => {
      const { clientHeight } = domRef.current!;
      setScrollY(clientHeight - 55);
    }, 64);
  }

  useEffect(() => {
    initScroll();
    window.addEventListener('resize', initScroll);
    return () => {
      window.removeEventListener('resize', initScroll)
    };
  }, []);

  return (
    <div className={styles.file_list}
      ref={domRef}>
      <div className={styles.table_wrapper}>
        <Table pagination={false}
          loading={loading}
          columns={columns}
          dataSource={wrappedFiles}
          rowKey='key'
          scroll={{ y: scrollY }}
          virtual={true}
        />
        {
          editRecord && (
            <UpdateFileModel id={editRecord.id}
              open={Boolean(editRecord)}
              directoryTree={props.directoryTree}
              onClose={() => setEditRecord(null)}
              values={{
                directory: editRecord.name,
                permission: editRecord.permission,
              }}
              afterUpdate={props.afterUpdate}
              type={editRecord.isDirectory? 'directory' : 'file'}
            />
          )
        }
        {
          uploadList.map((item, index) => (
            <FileUploadProgress key={index}
              name={item.name}
              progress={item.progress * 100}
            />
          )) 
        }
      </div>
    </div>
  );
}

export default FileList;
