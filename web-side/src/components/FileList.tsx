import { Table, TableColumnType } from 'antd';
import { format } from 'date-fns';
import { Link } from 'react-router';
import { FileTreeResponse } from '@/services/api';

interface Props {
  files: FileTreeResponse['files'];
  loading: boolean;
}

function FileList(props: Props) {
  const { files, loading } = props;

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
      title: '大小',
      dataIndex: 'size',
      render: (size: number, record) => {
        if (record.isDirectory) {
          return '-' 
        }

      }
    },
  ];

  return (
    <Table pagination={false}
      loading={loading}
      columns={columns}
      dataSource={files}
      rowKey="id"
      scroll={{ y: 0 }}
      virtual
    />
  );
}

export default FileList;
