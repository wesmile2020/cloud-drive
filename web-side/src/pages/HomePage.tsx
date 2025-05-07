import React from 'react';
import { Link, useParams } from 'react-router';
import { Breadcrumb, BreadcrumbProps } from 'antd';

import HomeMenu from '@/components/HomeMenu';
import FileList from '@/components/FileList';

import { FileTreeResponse, getFiles } from '@/services/api';

import styles from './HomePage.module.css';
import { HomeFilled, HomeOutlined } from '@ant-design/icons';

function createBreadcrumbItems(tree: FileTreeResponse['tree'] | null) {
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
    title: result.length > 0 ? (
      <Link to="/home/0">
        <HomeFilled />
        <span>首页</span>
      </Link>
    ) : (
      <>
        <HomeOutlined />
        <span>首页</span>
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
  const [loading, setLoading] = React.useState(false);
  const [breadcrumbs, setBreadcrumbs] = React.useState<BreadcrumbProps['items']>([]);

  React.useEffect(() => {
    setLoading(true);
    getFiles(directoryId).then((res) => {
      setFiles(res.files); 
      const items = createBreadcrumbItems(res.tree);
      setBreadcrumbs(items);
    }).finally(() => {
      setLoading(false); 
    });
  }, [directoryId]);

  return (
    <>
      <HomeMenu directoryId={directoryId}/>
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
