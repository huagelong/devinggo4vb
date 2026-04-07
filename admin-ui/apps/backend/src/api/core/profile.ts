import type { LogApi } from '#/api/system/log';

import { requestClient } from '#/api/request';

export namespace ProfileApi {
  export interface UpdateUserInfoPayload {
    nickname?: string;
    avatar?: string;
    email?: string;
    phone?: string;
    remark?: string;
    signed?: string;
  }

  export interface ModifyPasswordPayload {
    old_password?: string;
    new_password?: string;
    confirm_password?: string;
    oldPassword?: string;
    newPassword?: string;
    newPasswordConfirmation?: string;
  }

  export interface LoginLogQuery {
    page?: number;
    pageSize?: number;
    username?: string;
    status?: number;
    ip?: string;
    login_time?: string[];
  }

  export interface OperationLogQuery {
    page?: number;
    pageSize?: number;
    username?: string;
    service_name?: string;
    ip?: string;
    status?: number;
    created_at?: string[];
  }
}

// User profile APIs
export async function updateUserInfoApi(data: ProfileApi.UpdateUserInfoPayload) {
  return requestClient.post<void>('/system/user/updateInfo', data);
}

export async function modifyPasswordApi(data: ProfileApi.ModifyPasswordPayload) {
  return requestClient.post<void>('/system/user/modifyPassword', data);
}

// Log APIs (should use LogApi from system/log instead)
export async function getLoginLogListApi(params: LogApi.LoginLogQuery) {
  return requestClient.get<LogApi.LoginLogResponse>(
    '/system/common/getLoginLogList',
    { params },
  );
}

export async function getOperationLogListApi(params: LogApi.OperLogQuery) {
  return requestClient.get<LogApi.OperLogResponse>(
    '/system/common/getOperationLogList',
    { params },
  );
}
