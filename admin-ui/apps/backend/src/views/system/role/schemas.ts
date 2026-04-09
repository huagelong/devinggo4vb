import type {
  RoleColumnOptionItem,
  RoleFormModel,
  RoleSearchFormModel,
  RoleTableColumn,
} from './model';

export const roleDataScopeOptions = [
  { label: '全部数据权限', value: 1 },
  { label: '自定义数据权限', value: 2 },
  { label: '本部门数据权限', value: 3 },
  { label: '本部门及以下数据权限', value: 4 },
  { label: '本人数据权限', value: 5 },
  { label: '按部门过滤', value: 6 },
];

export function createRoleSearchForm(): RoleSearchFormModel {
  return {
    code: '',
    created_at: [],
    name: '',
    status: undefined,
  };
}

export function createRoleFormDefaultValues(): RoleFormModel {
  return {
    code: '',
    data_scope: 1,
    dept_ids: [],
    menu_ids: [],
    name: '',
    remark: '',
    sort: 1,
    status: 1,
  };
}

export function createRoleTableColumns(): RoleTableColumn[] {
  return [
    {
      align: 'center',
      colKey: 'row-select',
      type: 'multiple',
      width: 52,
    },
    { align: 'center', colKey: 'id', title: 'ID', width: 80 },
    { align: 'center', colKey: 'name', minWidth: 140, title: '角色名称' },
    { align: 'center', colKey: 'code', minWidth: 160, title: '角色标识' },
    { align: 'center', colKey: 'data_scope', minWidth: 160, title: '数据范围' },
    { align: 'center', colKey: 'sort', title: '排序', width: 140 },
    { align: 'center', colKey: 'status', title: '״̬', width: 120 },
    { align: 'center', colKey: 'remark', minWidth: 180, title: '备注' },
    { align: 'center', colKey: 'created_at', title: '创建时间', width: 180 },
    {
      align: 'center',
      colKey: 'action',
      fixed: 'right',
      title: '操作',
      width: 420,
    },
  ];
}

export function createRoleColumnOptions(
  columns: RoleTableColumn[],
): RoleColumnOptionItem[] {
  return columns
    .filter((column) => column.colKey !== 'row-select' && column.title)
    .map((column) => ({
      label: String(column.title),
      value: String(column.colKey),
    }));
}

export function getRoleDataScopeLabel(value?: number | string) {
  const normalizedValue = Number(value);
  return (
    roleDataScopeOptions.find((item) => Number(item.value) === normalizedValue)
      ?.label ?? '-'
  );
}
