import { requestClient } from '#/api/request';

export function getPostPageList(params: Record<string, any>) {
  return requestClient.get('/system/post/index', { params });
}

export function getPostList(params?: any) {
  return requestClient.get('/system/post/list', { params });
}

export function getRecyclePostList(params: Record<string, any>) {
  return requestClient.get('/system/post/recycle', { params });
}

export function savePost(data: Record<string, any>) {
  return requestClient.post('/system/post/save', data);
}

export function updatePost(id: number, data: Record<string, any>) {
  return requestClient.put(`/system/post/update/${id}`, data);
}

export function deletePost(ids: number[]) {
  return requestClient.delete('/system/post/delete', { data: { ids } });
}

export function realDeletePost(ids: number[]) {
  return requestClient.delete('/system/post/realDelete', { data: { ids } });
}

export function recoveryPost(ids: number[]) {
  return requestClient.put('/system/post/recovery', { ids });
}

export function changePostStatus(data: { id: number; status: number }) {
  return requestClient.put('/system/post/changeStatus', data);
}

export function updatePostSort(data: {
  id: number;
  numberName: string;
  numberValue: number;
}) {
  return requestClient.put('/system/post/numberOperation', data);
}
