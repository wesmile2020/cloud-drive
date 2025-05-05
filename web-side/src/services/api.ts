import axios from 'axios';
import { message } from 'antd';

axios.interceptors.response.use(
  (response) => {
    if (response.status === 200 && response.data.code === 200) {
      return response.data;
    } else {
      message.error(response.data.message);
      throw new Error(response.data.code);
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

export async function register(params: RegisterParams) {
  return axios.post('/api/user/register', params);
}

export async function login(account: string, password: string) {
  return axios.post('/api/user/login', { account, password });
}

export async function getUserInfo() {
  return axios.get('/api/user/info');
}

interface CreateDirectoryParams {
  name: string;
  parentId: number; 
  public: boolean;
}

export async function createDirectory(params: CreateDirectoryParams) {
  return axios.post('/api/file/directory', params);
}
