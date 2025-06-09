import axios from 'axios';
import { message } from 'antd';
import { Permission } from '@/config/enums';

axios.interceptors.response.use(
  (response) => {
    if (response.status === 200 && response.data.code === 200) {
      return response.data.data;
    } else {
      if (response.data.code === 401) {
        throw new Error(response.data.code);
      } else {
        message.error(response.data.message);
        throw new Error(response.data.message);
      }
    }
   
  },
  (error) => {
    message.error(error.message || '服务异常');
    throw error;
  }
)

export interface RegisterParams {
  name: string;
  email: string;
  password: string;
  phone: string;
}

export function register(params: RegisterParams) {
  return axios.post('/api/user/register', params);
}

export function login(account: string, password: string) {
  return axios.post('/api/user/login', { account, password });
}

export function logout() {
  return axios.post('/api/user/logout');
}

export interface UserInfo {
  id: number;
  name: string;
  email: string;
  phone: string;
}

export function getUserInfo(): Promise<UserInfo> {
  return axios.get('/api/user/info');
}

interface CreateDirectoryParams {
  name: string;
  parentId: number; 
  permission: Permission;
}

export function createDirectory(params: CreateDirectoryParams) {
  return axios.post('/api/file/directory', params);
}

interface UpdateAPIFileParams {
  name: string;
  permission: Permission;
}

export function updateDirectory(id: number, params: UpdateAPIFileParams) {
  return axios.put(`/api/file/directory/${id}`, params); 
}

export function deleteDirectory(id: number) {
  return axios.delete(`/api/file/directory/${id}`); 
}

interface FileItem {
  user: UserInfo;
  id: number;
  name: string;
  size: number;
  fileId: string;
  parentId: number;
  permission: Permission;
  public: boolean;
  isDirectory: boolean;
  timestamp: number;
}

interface FileTree {
  id: number;
  userId: number;
  name: string;
  public: boolean;
  permission: Permission;
  parent: FileTree | null;
}

export interface FileTreeResponse {
  files: FileItem[];
  tree: FileTree | null;
}

export function getFiles(parentId: number): Promise<FileTreeResponse> {
  return axios.get(`/api/file/${parentId}`);
}


interface UploadFileParams {
  file: Blob;
  name: string;
  parentId: number;
  permission: Permission;
  total: number;
  index: number;
  fileId?: string;
  onProgress?: (progress: number) => void;
}

export function uploadFile(params: UploadFileParams): Promise<{ fileId: string }> {
  const formData = new FormData();
  formData.append('file', params.file);
  formData.append('name', params.name);
  formData.append('parentId', params.parentId.toString());
  formData.append('permission', params.permission.toString());
  formData.append('total', params.total.toString());
  formData.append('index', params.index.toString());
  if (params.fileId) {
    formData.append('fileId', params.fileId);
  }

  return axios.post('/api/file/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }, 
    onUploadProgress(progressEvent) {
      params.onProgress?.(progressEvent.progress!);
    },
  })
}

export function deleteFile(id: number) {
  return axios.delete(`/api/file/${id}`); 
}

export function updateFile(id: number, params: UpdateAPIFileParams) {
  return axios.put(`/api/file/${id}`, params);
}
