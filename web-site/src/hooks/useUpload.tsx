import React from 'react';
import { Permission } from '@/config/enums';
import { uploadFile } from '@/services/api';
import { Schedule } from '@/utils/Schedule';
import { message } from 'antd';


interface UploadItem {
  name: string;
  directoryId: number;
  progress: number;
}

type UploadFunction = (file: File, parentId: number, permission: Permission) => Promise<unknown>;

const UploadContext = React.createContext<[UploadItem[], UploadFunction]>([[], () => Promise.resolve()]);

const SLICE_SIZE = 1024 * 1024 * 30; // 30MB per slice;

const schedule = new Schedule(2);

export function UploadProvider(props: React.PropsWithChildren) {
  const [uploadingFiles, setUploadingFiles] = React.useState<UploadItem[]>([]);

  const upload = React.useCallback((file: File, parentId: number, permission: Permission) => {
    const uploadItem: UploadItem = {
      name: file.name,
      progress: 0,
      directoryId: parentId,
    };
    setUploadingFiles((prev) => [...prev, uploadItem]);
    let hasUploaded = 0;
    const uploadSlice = async (start: number, end: number, fileId?: string) => {
      const slice = file.slice(start, end);
      return uploadFile({
        file: slice,
        name: file.name,
        parentId,
        permission,
        total: file.size,
        index: start,
        fileId,
        onProgress: (progress) => {
          hasUploaded += progress * slice.size;
          uploadItem.progress = hasUploaded / file.size;
          setUploadingFiles((prev) => [...prev]);
        }
      }).catch((err) => {
        if (err.message === 'Unauthorized') {
          message.error('请先登录，在上传文件');
        } else {
          message.error('上传失败，请重试');
        }
        return Promise.reject(err);
      });
    }

    // 先上传第一个切片获取fileId，再上传后续切片
    const promises: Promise<unknown>[] = [];

    const firstSize = Math.min(SLICE_SIZE, file.size);
    const firstPromise = uploadSlice(0, firstSize).then((data) => {
      const { fileId } = data;
      const totalSlice = Math.ceil((file.size - firstSize) / SLICE_SIZE);
      for (let i = 0; i < totalSlice; i += 1) {
        const start = firstSize + i * SLICE_SIZE;
        const end = Math.min(start + SLICE_SIZE, file.size);
        const slicePromise = new Promise((resolve, reject) => {
          schedule.add(() => {
            return uploadSlice(start, end, fileId).then(resolve).catch(reject);
          });
        });
        promises.push(slicePromise);
      }
    });
    return new Promise((resolve, reject) => {
      firstPromise.then(() => {
        Promise.all(promises).then(resolve).catch(reject);
      }).catch(reject);
    }).finally(() => {
      setUploadingFiles((prev) => prev.filter((item) => item !== uploadItem));
    });
  }, [setUploadingFiles]);

  return (
    <UploadContext value={[uploadingFiles, upload]}>
      {props.children}
    </UploadContext>
  );
}

// eslint-disable-next-line react-refresh/only-export-components
export const useUpload = () => React.useContext(UploadContext);
