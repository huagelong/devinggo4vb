import type { ColumnOptionItem, UserSearchFormModel } from './model';

export function createUserSearchForm(): UserSearchFormModel {
  return {
    created_at: [],
    dept_id: undefined,
    email: '',
    phone: '',
    post_id: undefined,
    role_id: undefined,
    status: undefined,
    user_type: undefined,
    username: '',
  };
}

export function createUserTableColumns() {
  return [
    {
      align: 'center',
      colKey: 'row-select',
      type: 'multiple',
      width: 50,
    },
    { align: 'center', colKey: 'avatar', title: '头像', width: 80 },
    { align: 'center', colKey: 'username', minWidth: 100, title: '账户' },
    { align: 'center', colKey: 'dept_name', minWidth: 100, title: '所属部门' },
    { align: 'center', colKey: 'nickname', minWidth: 100, title: '昵称' },
    { align: 'center', colKey: 'role_name', minWidth: 100, title: '角色' },
    { align: 'center', colKey: 'phone', minWidth: 120, title: '手机' },
    { align: 'center', colKey: 'post_name', minWidth: 100, title: '岗位' },
    { align: 'center', colKey: 'email', minWidth: 150, title: '邮箱' },
    { align: 'center', colKey: 'status', title: '状态', width: 100 },
    { align: 'center', colKey: 'user_type', title: '用户类型', width: 100 },
    { align: 'center', colKey: 'created_at', title: '注册时间', width: 160 },
    {
      align: 'center',
      colKey: 'action',
      fixed: 'right',
      title: '操作',
      width: 220,
    },
  ];
}

export function createUserColumnOptions(columns: any[]): ColumnOptionItem[] {
  return columns
    .filter((column) => column.colKey !== 'row-select' && column.title)
    .map((column) => ({
      label: String(column.title),
      value: String(column.colKey),
    }));
}

export const userActionDropdownOptions = [
  { content: '重置密码', value: 'reset_password' },
  { content: '更新缓存', value: 'clear_cache' },
  { content: '设置首页', value: 'set_homepage' },
];
