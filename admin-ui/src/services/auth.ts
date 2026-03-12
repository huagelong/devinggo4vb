import request from '../utils/request';

// API: /system/login POST
export const login = (data: any) => {
  return request.post('/system/login', data);
};

// API: /system/logout POST
export const logout = () => {
  return request.post('/system/logout');
};

// API: /system/refresh POST
export const refreshToken = () => {
  return request.post('/system/refresh');
};

// API: /system/getInfo GET (对应 user.go 获取自身信息及菜单路由)
export const getInfo = () => {
  return request.get('/system/getInfo');
};
