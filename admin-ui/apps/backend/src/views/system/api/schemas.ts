import type {
  ApiColumnFormModel,
  ApiColumnOptionItem,
  ApiColumnSearchFormModel,
  ApiColumnTableColumn,
  ApiFormModel,
  ApiSearchFormModel,
  ApiTableColumn,
} from './model';

export function createApiSearchForm(): ApiSearchFormModel {
  return {
    group_id: undefined,
    name: '',
    access_name: '',
    request_mode: undefined,
    status: undefined,
    created_at: [],
  };
}

export function createApiTableColumns(): ApiTableColumn[] {
  return [
    {
      align: 'center',
      colKey: 'row-select',
      type: 'multiple',
      width: 52,
    },
    { colKey: 'group_name', title: '所属组', minWidth: 160 },
    { colKey: 'name', title: '接口名称', minWidth: 200 },
    { colKey: 'access_name', title: '接口标识', minWidth: 220 },
    { colKey: 'request_mode', title: '请求模式', width: 120 },
    { colKey: 'auth_mode', title: '认证模式', width: 120 },
    { colKey: 'status', title: '״̬', width: 100, align: 'center' },
    { colKey: 'remark', title: '备注', minWidth: 200 },
    { colKey: 'created_at', title: '创建时间', minWidth: 180 },
    {
      align: 'center',
      colKey: 'action',
      fixed: 'right',
      title: '操作',
      width: 260,
    },
  ];
}

export function createApiColumnOptions(
  columns: ApiTableColumn[],
): ApiColumnOptionItem[] {
  return columns
    .filter((column) => column.colKey !== 'row-select' && column.title)
    .map((column) => ({
      label: String(column.title),
      value: String(column.colKey),
    }));
}

export function createApiColumnSearchForm(): ApiColumnSearchFormModel {
  return {
    name: '',
    data_type: undefined,
    status: undefined,
    is_required: undefined,
    created_at: [],
  };
}

export function createApiColumnTableColumns(): ApiColumnTableColumn[] {
  return [
    {
      align: 'center',
      colKey: 'row-select',
      type: 'multiple',
      width: 52,
    },
    { colKey: 'name', title: '字段名称', minWidth: 220 },
    { colKey: 'data_type', title: '数据类型', width: 140 },
    { colKey: 'type', title: '字段类型', width: 120 },
    { colKey: 'is_required', title: '是否必填', width: 100, align: 'center' },
    { colKey: 'status', title: '״̬', width: 100, align: 'center' },
    { colKey: 'default_value', title: '默认值', minWidth: 180 },
    { colKey: 'remark', title: '备注', minWidth: 200 },
    { colKey: 'created_at', title: '创建时间', minWidth: 180 },
    {
      align: 'center',
      colKey: 'action',
      fixed: 'right',
      title: '操作',
      width: 200,
    },
  ];
}

export function createApiColumnColumnOptions(
  columns: ApiColumnTableColumn[],
): ApiColumnOptionItem[] {
  return columns
    .filter((column) => column.colKey !== 'row-select' && column.title)
    .map((column) => ({
      label: String(column.title),
      value: String(column.colKey),
    }));
}

export function createApiFormDefaultValues(): ApiFormModel {
  return {
    group_id: undefined,
    name: '',
    access_name: '',
    request_mode: undefined,
    status: 1,
    auth_mode: 1,
    remark: '',
  };
}

export function createApiColumnFormDefaultValues(): ApiColumnFormModel {
  return {
    api_id: 0,
    data_type: undefined,
    default_value: '',
    description: '',
    is_required: 2,
    name: '',
    remark: '',
    status: 1,
    type: 1,
  };
}
