import axios, { AxiosError, AxiosInstance, InternalAxiosRequestConfig, AxiosResponse } from 'axios';
import { message } from 'antd';
import { useAuthStore } from '../stores/authStore';

// 创建 axios 实例
const request: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '', // 去掉默认的 /api 前缀，直接使用实际路径
  timeout: 10000,
});

// 请求拦截器
request.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // 从 zustand store 获取 token
    const token = useAuthStore.getState().token;
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`; // 这里的格式根据后端定义调整
    }
    return config;
  },
  (error: AxiosError) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
request.interceptors.response.use(
  (response: AxiosResponse) => {
    // 根据后端返回的数据结构调整
    const { code, data, msg } = response.data;

    // 假设非 0 或 200 为错误状态，这里根据具体业务调整
    if (code !== 0 && code !== 200) {
      message.error(msg || 'API Request Error');
      return Promise.reject(new Error(msg || 'Error'));
    }
    
    return data;
  },
  async (error: AxiosError) => {
    if (error.response) {
      const { status } = error.response;
      if (status === 401) {
        message.error('登录回话已过期，请重新登录');
        useAuthStore.getState().logout(); // 清除 token 并可能跳转登录页
      } else if (status === 403) {
        message.error('没有权限访问该资源');
      } else {
        message.error(error.message || '网络请求错误');
      }
    } else {
      message.error('网络连接异常');
    }
    return Promise.reject(error);
  }
);

export default request;
