import type {
  DeptColumnOptionItem,
  DeptFormModel,
  DeptSearchFormModel,
  DeptTableColumn,
} from './model';

export function createDeptSearchForm(): DeptSearchFormModel {
  return {
    created_at: [],
    leader: '',
    level: '',
    name: '',
    phone: '',
    status: undefined,
  };
}

export function createDeptFormDefaultValues(): DeptFormModel {
  return {
    leader: '',
    level: '',
    name: '',
    parent_id: 0,
    phone: '',
    remark: '',
    sort: 1,
    status: 1,
  };
}

export function createDeptTableColumns(): DeptTableColumn[] {
  return [
    {
      align: 'center',
      colKey: 'row-select',
      type: 'multiple',
      width: 52,
    },
    { colKey: 'name', minWidth: 180, title: '部门名称' },
    { align: 'center', colKey: 'leader', minWidth: 120, title: '负责人' },
    { align: 'center', colKey: 'phone', minWidth: 150, title: '手机' },
    { align: 'center', colKey: 'sort', title: '排序', width: 140 },
    { align: 'center', colKey: 'status', title: '״̬', width: 120 },
    { align: 'center', colKey: 'created_at', title: '创建时间', width: 180 },
    {
      align: 'center',
      colKey: 'action',
      fixed: 'right',
      title: '操作',
      width: 320,
    },
  ];
}

export function createDeptColumnOptions(
  columns: DeptTableColumn[],
): DeptColumnOptionItem[] {
  return columns
    .filter((column) => column.colKey !== 'row-select' && column.title)
    .map((column) => ({
      label: String(column.title),
      value: String(column.colKey),
    }));
}
