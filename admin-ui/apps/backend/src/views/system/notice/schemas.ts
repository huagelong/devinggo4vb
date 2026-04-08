import type {
  NoticeColumnOptionItem,
  NoticeSearchFormModel,
  NoticeTableColumn,
} from './model';

export function createNoticeSearchForm(): NoticeSearchFormModel {
  return {
    created_at: [],
    title: '',
    type: undefined,
  };
}

export function createNoticeTableColumns(): NoticeTableColumn[] {
  return [
    {
      align: 'center',
      colKey: 'row-select',
      type: 'multiple',
      width: 52,
    },
    { colKey: 'title', title: '公告标题', minWidth: 200 },
    { align: 'center', colKey: 'type', title: '公告类型', width: 100 },
    { colKey: 'content', title: '公告内容', minWidth: 260 },
    { colKey: 'remark', title: '备注', minWidth: 160 },
    { colKey: 'created_at', title: '创建时间', minWidth: 180 },
    {
      align: 'center',
      colKey: 'action',
      fixed: 'right',
      title: '操作',
      width: 220,
    },
  ];
}

export function createNoticeColumnOptions(
  columns: NoticeTableColumn[],
): NoticeColumnOptionItem[] {
  return columns
    .filter((column) => column.colKey !== 'row-select' && column.title)
    .map((column) => ({
      label: String(column.title),
      value: String(column.colKey),
    }));
}
