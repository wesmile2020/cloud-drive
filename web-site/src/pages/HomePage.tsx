import React from 'react';
import { Link, useParams } from 'react-router';
import { Breadcrumb, BreadcrumbProps, message, Upload, UploadProps } from 'antd';
import { HomeFilled, InboxOutlined } from '@ant-design/icons';

import HomeMenu from '@/components/HomeMenu';
import FileList from '@/components/FileList';

import { FileTreeResponse, getFiles } from '@/services/api';

import styles from './HomePage.module.css';
import { Permission } from '@/config/enums';
import { useUpload } from '@/hooks/useUpload';

const homeTree: FileTreeResponse['tree'] = {
  id: 0,
  userId: 0,
  name: '首页',
  parent: null,
  public: true,
  permission: Permission.public, 
};

function createBreadcrumbItems(directoryId: number, tree: FileTreeResponse['tree']) {
  const result: BreadcrumbProps['items'] = [];
  while (tree) {
    result.unshift({
      title: result.length > 0 ? (
        <Link to={`/home/${tree.id}`}>
          {tree.name}
        </Link>
      ) : tree.name,
    });
    tree = tree.parent;
  }
  result.unshift({
    title: directoryId !== 0 ? (
      <Link to="/home/0">
        <HomeFilled />
        <span className={styles.text}>首页</span>
      </Link>
    ) : (
      <>
        <HomeFilled />
        <span className={styles.text}>首页</span>
      </>
    ),
  });

  return result;
}

function HomePage() {
  const params = useParams();
  const directoryId = React.useMemo(() => {
    return Number(params.id);
  }, [params.id]);
  const [files, setFiles] = React.useState<FileTreeResponse['files']>([]);
  const [fileTree, setFileTree] = React.useState<FileTreeResponse['tree']>(null);
  const [loading, setLoading] = React.useState(false);
  const [breadcrumbs, setBreadcrumbs] = React.useState<BreadcrumbProps['items']>([]);
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const [_, upload] = useUpload();

  const fetchFiles = React.useCallback(() => {
    setLoading(true);
    getFiles(directoryId).then((res) => {
      setFiles(res.files);
      if (directoryId === 0) {
        setFileTree(homeTree)
      } else {
        setFileTree(res.tree);
      }
      const items = createBreadcrumbItems(directoryId, res.tree);
      setBreadcrumbs(items);
    
    }).finally(() => {
      setLoading(false);
    });
  }, [directoryId]);

  React.useEffect(() => {
    fetchFiles();
  }, [fetchFiles]);

  const uploadAction: UploadProps['customRequest'] = async (options) => {
    const file = options.file as File;
    await upload(file, directoryId, Permission.inherit);
    message.success('上传成功');
    fetchFiles();
  }

  return (
    <>
      <HomeMenu directoryTree={fileTree}
        afterCreate={fetchFiles}
      />
      <div className={styles.breadcrumb_wrapper}>
        <Breadcrumb className={styles.breadcrumb}
          items={breadcrumbs}
        />
      </div>
      <div className={styles.upload}>
        <Upload.Dragger showUploadList={false}
          customRequest={uploadAction}>
          <p className="ant-upload-drag-icon">
            <InboxOutlined />
          </p>
          <p className="ant-upload-text">
            点击或拖拽文件到此区域上传
          </p>
          <p className="ant-upload-hint">
            涉密文件请不要上传
          </p>
        </Upload.Dragger>
      </div>
      
      <FileList files={files}
        directoryTree={fileTree}
        loading={loading}
        afterUpdate={fetchFiles}
      />
    </> 
  );
}

export default HomePage;
