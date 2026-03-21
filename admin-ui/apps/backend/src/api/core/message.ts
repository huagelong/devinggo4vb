import { requestClient } from '#/api/request';

export function getQueueMessageReceiveListApi(params: any) {
  return requestClient.get('/system/queueMessage/receiveList', { params });
}

export function updateQueueMessageReadStatusApi(data: { ids: number[] | string[] }) {
  return requestClient.put('/system/queueMessage/updateReadStatus', data);
}

export function deleteQueueMessageApi(data: { ids: number[] | string[] }) {
  return requestClient.delete('/system/queueMessage/deletes', { data });
}

export function getDataDictListApi(params: { code: string }) {
  return requestClient.get('/system/dataDict/list', { params });
}
