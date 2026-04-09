import type {
  DictColumnOptionItem,
  DictDataFormModel,
  DictDataSearchFormModel,
  DictTableColumn,
  DictTypeFormModel,
  DictTypeSearchFormModel,
} from './model';

export function createDictTypeSearchForm(): DictTypeSearchFormModel {
  return {
    code: '',
    created_at: [],
    name: '',
    status: undefined,
  };
}

export function createDictTypeFormDefaultValues(): DictTypeFormModel {
  return {
    code: '',
    name: '',
    remark: '',
    status: 1,
  };
}

export function createDictTypeTableColumns(): DictTableColumn[] {
  return [
    {
      align: 'center',
      colKey: 'row-select',
      type: 'multiple',
      width: 52,
    },
    { colKey: 'name', minWidth: 160, title: '字典名称' },
    { colKey: 'code', minWidth: 200, title: '字典标识' },
    { align: 'center', colKey: 'status', title: '״̬', width: 120 },
    { colKey: 'remark', minWidth: 180, title: '备注' },
    { colKey: 'created_at', minWidth: 180, title: '创建时间' },
    {
      align: 'center',
      colKey: 'action',
      fixed: 'right',
      title: '操作',
      width: 320,
    },
  ];
}

export function createDictTypeColumnOptions(columns: DictTableColumn[]): DictColumnOptionItem[] {
  return columns
    .filter((column) => column.colKey !== 'row-select' && column.title)
    .map((column) => ({
      label: String(column.title),
      value: String(column.colKey),
    }));
}

export function createDictDataSearchForm(): DictDataSearchFormModel {
  return {
    code: '',
    created_at: [],
    label: '',
    status: undefined,
    type_id: undefined,
    value: '',
  };
}

export function createDictDataFormDefaultValues(typeId?: number, code?: string): DictDataFormModel {
  return {
    code: code ?? '',
    label: '',
    remark: '',
    sort: 1,
    status: 1,
    type_id: typeId,
    value: '',
  };
}

export function createDictDataTableColumns(): DictTableColumn[] {
  return [
    {
      align: 'center',
      colKey: 'row-select',
      type: 'multiple',
      width: 52,
    },
    { colKey: 'label', minWidth: 160, title: '字典标签' },
    { colKey: 'value', minWidth: 140, title: '字典键值' },
    { align: 'center', colKey: 'sort', title: '排序', width: 120 },
    { align: 'center', colKey: 'status', title: '״̬', width: 120 },
    { colKey: 'remark', minWidth: 180, title: '备注' },
    { colKey: 'created_at', minWidth: 180, title: '创建时间' },
    {
      align: 'center',
      colKey: 'action',
      fixed: 'right',
      title: '操作',
      width: 240,
    },
  ];
}
