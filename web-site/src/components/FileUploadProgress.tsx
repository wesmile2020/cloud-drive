import { Progress, Spin } from 'antd';
import FileName from './FileName';
import styles from './FileUploadProgress.module.css';

interface Props {
  name: string;
  progress: number;
}

function FileUploadProgress(props: Props) {
  return (
    <div className={styles.upload_progress}>
      <Spin className={styles.spin}/>
      <FileName name={props.name}/>
      <Progress className={styles.upload_bar}
        percent={props.progress}
        size={[350, 6]}
        format={(percent = 0) => {
          if (percent === 100) {
            return '文件保存中...';
          }
          return `${Math.round(percent * 100) / 100}%`;
        }}
      />
    </div>
  )
}

export default FileUploadProgress;
