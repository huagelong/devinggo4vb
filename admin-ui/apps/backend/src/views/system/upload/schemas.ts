import type { UploadListItem, UploadTableColumn, UploadTreeItem } from './model';

export const defaultUploadTreeData: UploadTreeItem[] = [
  {
    label: '全部',
    value: 'all',
  },
  {
    label: 'ͼƬ',
    value: 'image',
    children: [
      { label: 'JPEG', value: 'image/jpeg' },
      { label: 'PNG', value: 'image/png' },
      { label: 'GIF', value: 'image/gif' },
      { label: 'WebP', value: 'image/webp' },
    ],
  },
  {
    label: '文档',
    value: 'document',
    children: [
      { label: 'PDF', value: 'application/pdf' },
      { label: 'Word', value: 'application/msword' },
      { label: 'Excel', value: 'application/vnd.ms-excel' },
      { label: 'PowerPoint', value: 'application/vnd.ms-powerpoint' },
    ],
  },
  {
    label: '其他',
    value: 'other',
  },
];

export function createUploadTableColumns(): UploadTableColumn[] {
  return [
    { colKey: 'row-select', width: 50, fixed: 'left' },
    { title: 'ID', colKey: 'id', width: 80 },
    { title: '原文件名', colKey: 'origin_name', ellipsis: true },
    { title: 'MIME类型', colKey: 'mime_type', width: 150 },
    { title: '存储路径', colKey: 'storage_path', ellipsis: true },
    { title: '文件大小', colKey: 'size_info', width: 100 },
    { title: '存储方式', colKey: 'storage_mode', width: 100 },
    { title: '创建时间', colKey: 'created_at', width: 180 },
    { title: '操作', colKey: 'action', width: 200, align: 'center', fixed: 'right' },
  ];
}

export function createUploadColumnOptions(columns: UploadTableColumn[]) {
  return columns
    .filter((col) => col.colKey !== 'row-select' && col.colKey !== 'action')
    .map((col) => ({
      label: col.title || col.colKey,
      value: col.colKey,
    }));
}

export function createUploadSearchForm() {
  return {
    origin_name: '',
    mime_type: '',
    storage_mode: undefined as number | undefined,
    created_at: [] as string[],
  };
}
