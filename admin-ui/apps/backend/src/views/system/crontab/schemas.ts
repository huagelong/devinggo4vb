import type {
  CrontabColumnOptionItem,
  CrontabFormModel,
  CrontabSearchFormModel,
  CrontabTableColumn,
} from './model';

export const crontabTypeOptions = [
  { label: '定时任务', value: 1 },
  { label: 'Cron任务', value: 2 },
];

export const crontabFinallyOptions = [
  { label: '是', value: 1 },
  { label: '否', value: 2 },
];

export function createCrontabSearchForm(): CrontabSearchFormModel {
  return {
    created_at: [],
    is_finally: undefined,
    name: '',
    type: undefined,
  };
}

export function createCrontabFormDefaultValues(): CrontabFormModel {
  return {
    is_finally: 2,
    name: '',
    remark: '',
    rule: '',
    target: '',
    type: 1,
  };
}

export function createCrontabTableColumns(): CrontabTableColumn[] {
  return [
    {
      align: 'center',
      colKey: 'row-select',
      type: 'multiple',
      width: 52,
    },
    { colKey: 'id', title: 'ID', width: 80 },
    { colKey: 'name', title: '任务名称', minWidth: 160 },
    { colKey: 'type', title: '任务类型', width: 120 },
    { colKey: 'rule', title: '执行规则', minWidth: 180 },
    { colKey: 'is_finally', title: '最终执行', width: 100 },
    { colKey: 'created_at', title: '创建时间', width: 180 },
    {
      align: 'center',
      colKey: 'action',
      fixed: 'right',
      title: '操作',
      width: 360,
    },
  ];
}

export function createCrontabColumnOptions(
  columns: CrontabTableColumn[],
): CrontabColumnOptionItem[] {
  return columns
    .filter((column) => column.colKey !== 'row-select' && column.title)
    .map((column) => ({
      label: String(column.title),
      value: String(column.colKey),
    }));
}
