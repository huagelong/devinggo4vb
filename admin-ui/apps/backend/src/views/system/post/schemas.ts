import type {
  PostColumnOptionItem,
  PostFormModel,
  PostSearchFormModel,
  PostTableColumn,
} from './model';

export function createPostSearchForm(): PostSearchFormModel {
  return {
    code: '',
    created_at: [],
    name: '',
    status: undefined,
  };
}

export function createPostFormDefaultValues(): PostFormModel {
  return {
    code: '',
    name: '',
    remark: '',
    sort: 1,
    status: 1,
  };
}

export function createPostTableColumns(): PostTableColumn[] {
  return [
    {
      align: 'center',
      colKey: 'row-select',
      type: 'multiple',
      width: 52,
    },
    { align: 'center', colKey: 'id', title: 'ID', width: 80 },
    { align: 'center', colKey: 'name', minWidth: 140, title: '岗位名称' },
    { align: 'center', colKey: 'code', minWidth: 140, title: '岗位标识' },
    { align: 'center', colKey: 'sort', title: '排序', width: 140 },
    { align: 'center', colKey: 'status', title: '״̬', width: 120 },
    { align: 'center', colKey: 'remark', minWidth: 160, title: '备注' },
    { align: 'center', colKey: 'created_at', title: '创建时间', width: 180 },
    {
      align: 'center',
      colKey: 'action',
      fixed: 'right',
      title: '操作',
      width: 220,
    },
  ];
}

export function createPostColumnOptions(
  columns: PostTableColumn[],
): PostColumnOptionItem[] {
  return columns
    .filter((column) => column.colKey !== 'row-select' && column.title)
    .map((column) => ({
      label: String(column.title),
      value: String(column.colKey),
    }));
}
