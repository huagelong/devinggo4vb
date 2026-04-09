import type {
  SystemModulesColumnOptionItem,
  SystemModulesFormModel,
  SystemModulesSearchFormModel,
  SystemModulesTableColumn,
} from './model';

export function createSystemModulesSearchForm(): SystemModulesSearchFormModel {
  return {
    created_at: [],
    name: '',
    status: undefined,
  };
}

export function createSystemModulesFormDefaultValues(): SystemModulesFormModel {
  return {
    name: '',
    remark: '',
    sort: 1,
    status: 1,
  };
}

export function createSystemModulesTableColumns(): SystemModulesTableColumn[] {
  return [
    {
      align: 'center',
      colKey: 'row-select',
      type: 'multiple',
      width: 52,
    },
    { colKey: 'id', title: 'ID', width: 80 },
    { colKey: 'name', title: '模块名称', minWidth: 160 },
    { colKey: 'sort', title: '排序', width: 120 },
    { colKey: 'status', title: '״̬', width: 120 },
    { colKey: 'created_at', title: '创建时间', width: 180 },
    {
      align: 'center',
      colKey: 'action',
      fixed: 'right',
      title: '操作',
      width: 240,
    },
  ];
}

export function createSystemModulesColumnOptions(
  columns: SystemModulesTableColumn[],
): SystemModulesColumnOptionItem[] {
  return columns
    .filter((column) => column.colKey !== 'row-select' && column.title)
    .map((column) => ({
      label: String(column.title),
      value: String(column.colKey),
    }));
}
