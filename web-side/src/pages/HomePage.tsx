import React from 'react';
import { Link, useParams } from 'react-router';
import { Breadcrumb, BreadcrumbProps } from 'antd';
import { HomeFilled } from '@ant-design/icons';

import HomeMenu from '@/components/HomeMenu';
import FileList from '@/components/FileList';

import { FileTreeResponse, getFiles } from '@/services/api';

import styles from './HomePage.module.css';
import { Permission } from '@/config/enums';

const homeTree: FileTreeResponse['tree'] = {
  id: 0,
  userId: 0,
  name: '首页',
  parent: null,
  public: true,
  permission: Permission.public, 
};

function createBreadcrumbItems(directoryId: number, tree: FileTreeResponse['tree'] | null) {
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

function Home() {
  const params = useParams();
  const directoryId = React.useMemo(() => {
    return Number(params.id);
  }, [params.id]);
  const [files, setFiles] = React.useState<FileTreeResponse['files']>([]);
  const [fileTree, setFileTree] = React.useState<FileTreeResponse['tree'] | null>(null);
  const [loading, setLoading] = React.useState(false);
  const [breadcrumbs, setBreadcrumbs] = React.useState<BreadcrumbProps['items']>([]);

  const fetchFiles = React.useCallback(() => {
    setLoading(true);
    getFiles(directoryId).then((res) => {
      setFiles(res.files);
      const items = createBreadcrumbItems(directoryId, res.tree);
      setBreadcrumbs(items);
      if (directoryId === 0) {
        setFileTree(homeTree)
      } else {
        setFileTree(res.tree);
      }
    }).finally(() => {
      setLoading(false);
    });
  }, [directoryId]);

  React.useEffect(() => {
    fetchFiles();
  }, [fetchFiles]);

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
      <FileList files={files}
        loading={loading}
      />
    </> 
  );
}

export default Home;
