import type {
  MenuColumnOptionItem,
  MenuFormModel,
  MenuSearchFormModel,
  MenuTableColumn,
  MenuTypeValue,
} from './model';

export const menuTypeOptions: Array<{ label: string; value: MenuTypeValue }> = [
  { label: 'Ŀ¼', value: 'M' },
  { label: '按钮', value: 'B' },
  { label: '外链', value: 'L' },
  { label: 'iFrame', value: 'I' },
];

export const menuHiddenOptions = [
  { label: '是', value: 1 },
  { label: '否', value: 2 },
];

export const restfulOptions = [
  { label: '是', value: '1' },
  { label: '否', value: '2' },
];

export const menuTypeTagMap: Record<string, { label: string; theme: 'default' | 'primary' | 'success' | 'warning' | 'danger' }> =
  {
    M: { label: 'Ŀ¼', theme: 'primary' },
    B: { label: '按钮', theme: 'warning' },
    L: { label: '外链', theme: 'success' },
    I: { label: 'iFrame', theme: 'default' },
  };

export function createMenuSearchForm(): MenuSearchFormModel {
  return {
    code: '',
    created_at: [],
    level: '',
    name: '',
    status: undefined,
  };
}

export function createMenuFormDefaultValues(): MenuFormModel {
  return {
    code: '',
    component: '',
    icon: '',
    is_hidden: 2,
    level: '',
    name: '',
    parent_id: 0,
    redirect: '',
    remark: '',
    restful: '2',
    route: '',
    sort: 1,
    status: 1,
    type: 'M',
  };
}

export function createMenuTableColumns(): MenuTableColumn[] {
  return [
    {
      align: 'center',
      colKey: 'row-select',
      type: 'multiple',
      width: 52,
    },
    { colKey: 'name', minWidth: 200, title: '菜单名称' },
    { align: 'center', colKey: 'type', title: '类型', width: 120 },
    { colKey: 'code', minWidth: 180, title: '菜单标识' },
    { colKey: 'icon', minWidth: 120, title: '图标' },
    { colKey: 'route', minWidth: 180, title: '路由地址' },
    { colKey: 'component', minWidth: 200, title: '组件路径' },
    { align: 'center', colKey: 'sort', title: '排序', width: 120 },
    { align: 'center', colKey: 'status', title: '״̬', width: 120 },
    { colKey: 'created_at', minWidth: 180, title: '创建时间' },
    {
      align: 'center',
      colKey: 'action',
      fixed: 'right',
      title: '操作',
      width: 360,
    },
  ];
}

export function createMenuColumnOptions(columns: MenuTableColumn[]): MenuColumnOptionItem[] {
  return columns
    .filter((column) => column.colKey !== 'row-select' && column.title)
    .map((column) => ({
      label: String(column.title),
      value: String(column.colKey),
    }));
}
