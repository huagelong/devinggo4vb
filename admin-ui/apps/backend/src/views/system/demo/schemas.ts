import type {
  DemoColumnOptionItem,
  DemoFormModel,
  DemoSearchFormModel,
  DemoTableColumn,
} from './model';

export function createDemoSearchForm(): DemoSearchFormModel {
  return {
    name: '',
    code: '',
    status: undefined,
    birthday: '',
    created_at: [],
  };
}

export function createDemoFormDefaultValues(): DemoFormModel {
  return {
    name: '',
    code: '',
    status: 1,
    sort: 1,
    price: 0,
    cover: '',
    email: '',
    phone: '',
    birthday: '',
    remark: '',
  };
}

export function createDemoTableColumns(): DemoTableColumn[] {
  return [
    {
      align: 'center',
      colKey: 'row-select',
      type: 'multiple',
      width: 52,
    },
    { align: 'center', colKey: 'id', title: 'ID', width: 80 },
    { align: 'center', colKey: 'name', minWidth: 120, title: '名称' },
    { align: 'center', colKey: 'code', minWidth: 120, title: '编码' },
    { align: 'center', colKey: 'status', minWidth: 120, title: '״̬' },
    { align: 'center', colKey: 'sort', minWidth: 140, title: '排序' },
    { align: 'center', colKey: 'price', minWidth: 120, title: '数字' },
    { align: 'center', colKey: 'email', minWidth: 120, title: '邮箱' },
    { align: 'center', colKey: 'phone', minWidth: 120, title: '手机号' },
    { align: 'center', colKey: 'birthday', minWidth: 120, title: '日期' },
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

export function createDemoColumnOptions(
  columns: DemoTableColumn[],
): DemoColumnOptionItem[] {
  return columns
    .filter((column) => column.colKey !== 'row-select' && column.title)
    .map((column) => ({
      label: String(column.title),
      value: String(column.colKey),
    }));
}
