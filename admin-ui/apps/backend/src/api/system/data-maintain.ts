import { requestClient } from '#/api/request';
import type { PageQuery, PageResponse } from '#/types/paging';

export namespace DataMaintainApi {
  export interface ListItem {
    name: string;
    collation?: string;
    comment?: string;
    engine?: string;
    create_time?: string;
    rows?: number;
  }

  export interface ListQuery extends Partial<PageQuery> {
    group_name?: string;
    name?: string;
  }

  export type ListResponse = PageResponse<ListItem>;
}

export function getDataMaintainPageList(params: DataMaintainApi.ListQuery) {
  return requestClient.get<DataMaintainApi.ListResponse>(
    '/system/dataMaintain/index',
    { params },
  );
}
