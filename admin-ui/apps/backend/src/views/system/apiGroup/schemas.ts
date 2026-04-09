import type {
  ApiGroupColumnOptionItem,
  ApiGroupFormModel,
  ApiGroupSearchFormModel,
  ApiGroupTableColumn,
} from './model';

export function createApiGroupSearchForm(): ApiGroupSearchFormModel {
  return {
    name: '',
    status: undefined,
    created_at: [],
  };
}

export function createApiGroupTableColumns(): ApiGroupTableColumn[] {
  return [
    {
      align: 'center',
      colKey: 'row-select',
      type: 'multiple',
      width: 52,
    },
    { colKey: 'name', title: '分组名称', minWidth: 220 },
    { colKey: 'status', title: '״̬', width: 120, align: 'center' },
    { colKey: 'remark', title: '备注', minWidth: 200 },
    { colKey: 'created_at', title: '创建时间', minWidth: 180 },
    { colKey: 'updated_at', title: '更新时间', minWidth: 180 },
    {
      align: 'center',
      colKey: 'action',
      fixed: 'right',
      title: '操作',
      width: 200,
    },
  ];
}

export function createApiGroupColumnOptions(
  columns: ApiGroupTableColumn[],
): ApiGroupColumnOptionItem[] {
  return columns
    .filter((column) => column.colKey !== 'row-select' && column.title)
    .map((column) => ({
      label: String(column.title),
      value: String(column.colKey),
    }));
}

export function createApiGroupFormDefaultValues(): ApiGroupFormModel {
  return {
    name: '',
    remark: '',
    status: 1,
  };
}
