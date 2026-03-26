import type { UserInfo } from '@vben/types';

import { requestClient } from '#/api/request';

export interface SystemUserInfoResult {
  user: {
    [key: string]: any;
    avatar: string;
    dashboard?: string;
    id: number;
    nickname: string;
    username: string;
  };
  roles: string[];
  codes: string[];
  routers: any[];
}

// 模块级缓存：避免同一事件循环内重复请求
let _getInfoPromise: null | Promise<SystemUserInfoResult> = null;

/**
 * 获取完整的系统用户信息（含角色、权限码、菜单路由）
 * 同一时刻的并发调用会复用同一个 Promise
 */
export function getSystemInfoApi(): Promise<SystemUserInfoResult> {
  if (!_getInfoPromise) {
    _getInfoPromise =
      requestClient.get<SystemUserInfoResult>('/system/getInfo');
    _getInfoPromise.finally(() => {
      _getInfoPromise = null;
    });
  }
  return _getInfoPromise;
}

/**
 * 获取用户信息（转换为前端 UserInfo 格式）
 */
export async function getUserInfoApi(): Promise<UserInfo> {
  const info = await getSystemInfoApi();
  return {
    avatar: info.user.avatar || '',
    desc: info.user.signed || '',
    homePath: info.user.dashboard || '/analytics',
    realName: info.user.nickname || info.user.username,
    roles: info.roles,
    token: '',
    userId: String(info.user.id),
    username: info.user.username,
  };
}

/**
 * 获取用户权限码
 */
export async function getAccessCodesApi(): Promise<string[]> {
  const info = await getSystemInfoApi();
  return info.codes;
}
