import { requestClient } from '#/api/request';
import type { PageResponse } from '#/types/paging';

export namespace MonitorApi {
  export interface OnlineUserItem {
    id: number;
    username: string;
    nickname: string;
    app_id: number;
    login_ip: string;
    login_time: string;
  }

  export interface OnlineUserQuery {
    username?: string;
    app_id?: number;
    page?: number;
    page_size?: number;
  }

  export interface OnlineUserResponse extends PageResponse<OnlineUserItem> {}

  export interface KickPayload {
    id: number;
    app_id: number;
  }

  export interface CacheServerInfo {
    version?: string;
    clients?: string;
    redis_mode?: string;
    run_days?: string;
    port?: string;
    aof_enabled?: string;
    expired_keys?: string;
    sys_total_keys?: string;
    use_memory?: string;
  }

  export interface CacheInfo {
    server: CacheServerInfo;
    keys: string[];
  }

  export interface ViewCachePayload {
    key: string;
  }

  export interface ViewCacheResponse {
    data: {
      content: string;
    };
  }

  export interface DeleteCachePayload {
    key: string;
  }
}

export function getOnlineUserPageList(params: MonitorApi.OnlineUserQuery) {
  return requestClient.get<MonitorApi.OnlineUserResponse>(
    '/system/onlineUser/index',
    { params },
  );
}

export function kickUser(data: MonitorApi.KickPayload) {
  return requestClient.post<void>('/system/onlineUser/kick', data);
}

export function getCacheInfo() {
  return requestClient.get<MonitorApi.CacheInfo>('/system/cache/monitor');
}

export function viewCache(data: MonitorApi.ViewCachePayload) {
  return requestClient.post<MonitorApi.ViewCacheResponse>(
    '/system/cache/view',
    data,
  );
}

export function deleteCacheKey(data: MonitorApi.DeleteCachePayload) {
  return requestClient.delete<void>('/system/cache/delete', { data });
}

export function clearAllCache() {
  return requestClient.delete<void>('/system/cache/clear');
}
