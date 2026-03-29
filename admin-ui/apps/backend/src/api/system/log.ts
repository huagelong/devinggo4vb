import { requestClient } from '#/api/request';
import type { PageResponse } from '#/types/paging';

export namespace LogApi {
  export interface LoginLogItem {
    id: number;
    username: string;
    status: number;
    ip: string;
    ip_location: string;
    os: string;
    browser: string;
    message: string;
    login_time: string;
    created_at?: string;
  }

  export interface LoginLogQuery {
    username?: string;
    status?: number;
    ip?: string;
    login_time?: string[];
    page?: number;
    page_size?: number;
  }

  export interface OperLogItem {
    id: number;
    username: string;
    service_name: string;
    router: string;
    method: string;
    ip: string;
    ip_location: string;
    response_code: number;
    request_data: string;
    response_data: string;
    created_at: string;
  }

  export interface OperLogQuery {
    username?: string;
    service_name?: string;
    ip?: string;
    created_at?: string[];
    page?: number;
    page_size?: number;
  }

  export interface ApiLogItem {
    id: number;
    api_name: string;
    access_name: string;
    response_code: number;
    access_time: string;
    ip: string;
    ip_location: string;
    request_data: string;
    response_data: string;
    created_at?: string;
  }

  export interface ApiLogQuery {
    api_name?: string;
    access_name?: string;
    ip?: string;
    access_time?: string[];
    page?: number;
    page_size?: number;
  }

  export type LoginLogResponse = PageResponse<LoginLogItem>;
  export type OperLogResponse = PageResponse<OperLogItem>;
  export type ApiLogResponse = PageResponse<ApiLogItem>;
}

// Login Log APIs
export function getLoginLogPageList(params: LogApi.LoginLogQuery) {
  return requestClient.get<LogApi.LoginLogResponse>(
    '/system/loginLog/index',
    { params },
  );
}

export function deleteLoginLog(ids: number[]) {
  return requestClient.delete<void>('/system/loginLog/delete', {
    data: { ids },
  });
}

// Operation Log APIs
export function getOperLogPageList(params: LogApi.OperLogQuery) {
  return requestClient.get<LogApi.OperLogResponse>('/system/operLog/index', {
    params,
  });
}

export function deleteOperLog(ids: number[]) {
  return requestClient.delete<void>('/system/operLog/delete', {
    data: { ids },
  });
}

// API Log APIs
export function getApiLogPageList(params: LogApi.ApiLogQuery) {
  return requestClient.get<LogApi.ApiLogResponse>('/system/apiLog/index', {
    params,
  });
}

export function deleteApiLog(ids: number[]) {
  return requestClient.delete<void>('/system/apiLog/delete', {
    data: { ids },
  });
}
