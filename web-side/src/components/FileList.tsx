import React from 'react';
import { Avatar, Table, TableColumnType } from 'antd';
import { format } from 'date-fns';
import { Link } from 'react-router';
import { FileTreeResponse } from '@/services/api';

import styles from './FileList.module.css';
import { formatSize } from '@/utils/utils';

interface Props {
  files: FileTreeResponse['files'];
  loading: boolean;
}

function FileList(props: Props) {
  const { files, loading } = props;
  const domRef = React.useRef<HTMLDivElement>(null);
  const [scrollY, setScrollY] = React.useState(0);

  const columns: TableColumnType<FileTreeResponse['files'][0]>[] = [
    {
      title: '名称',
      dataIndex: 'name',
      render(name: string, record) {
        if (record.isDirectory) {
          return <Link to={`/home/${record.id}`}>{name}</Link>;
        } else {
          return name;
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
          <Avatar className={styles.avatar}>
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
      dataIndex: 'public',
      render: (publicFlag: boolean) => {
        return publicFlag ? '公开' : '私密';
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
