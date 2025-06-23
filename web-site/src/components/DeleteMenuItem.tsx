import { DeleteOutlined } from '@ant-design/icons';
import { message, Popconfirm } from 'antd';
import { FileTreeResponse, deleteDirectory, deleteFile } from '@/services/api';

interface Props {
  record: FileTreeResponse['files'][0];
  afterDelete?: () => void;
  disabled?: boolean;
  className?: string;
}

function DeleteMenuItem(props: Props) {
  const { record } = props;
  async function deleteItem() {
    try {
      if (record.isDirectory) {
        await deleteDirectory(record.id);
      } else {
        await deleteFile(record.id)
      }
      message.success('删除成功');
    } catch (error) {
      console.error(error);
    }
    
    props.afterDelete?.(); 
  }

  const disabledClassName = props.disabled ? 'ant-dropdown-menu-item-disabled' : '';
  return (
    <Popconfirm title="确定删除吗？"
      okText="确定"
      cancelText="取消"
      onConfirm={deleteItem}
      disabled={props.disabled}
    >
       <li className={`ant-dropdown-menu-item ${props.className} ${disabledClassName}`}
        role="menuitem"
        tabIndex={-1}
        aria-describedby="«rh»"
        aria-disabled={false}
      >
        <DeleteOutlined className="ant-dropdown-menu-item-icon"/>
        <span className="ant-dropdown-menu-item-text">
          删除
        </span>
      </li>
    </Popconfirm>
  );
}

export default DeleteMenuItem;
