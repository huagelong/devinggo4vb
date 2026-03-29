import type {
  AttachmentColumnOptionItem,
  AttachmentSearchFormModel,
  AttachmentTableColumn,
  AttachmentTreeItem,
} from './model';

export const storageModeOptions = [
  { label: '本地存储', value: 1 },
  { label: '阿里云OSS', value: 2 },
  { label: '腾讯云COS', value: 3 },
  { label: '七牛云', value: 4 },
  { label: 'FTP', value: 5 },
];

export const defaultAttachmentTreeData: AttachmentTreeItem[] = [
  { title: '所有', key: 'all' },
  { title: '图片', key: 'image' },
  { title: '视频', key: 'video' },
  { title: '音频', key: 'audio' },
  { title: '文档', key: 'document' },
  { title: '压缩包', key: 'archive' },
  { title: '其他', key: 'other' },
];

export function createAttachmentSearchForm(): AttachmentSearchFormModel {
  return {
    created_at: [],
    mime_type: undefined,
    origin_name: '',
    storage_mode: undefined,
  };
}

export function createAttachmentTableColumns(): AttachmentTableColumn[] {
  return [
    {
      align: 'center',
      colKey: 'row-select',
      type: 'multiple',
      width: 52,
    },
    { colKey: 'id', title: 'ID', width: 80 },
    {
      colKey: 'url',
      title: '预览',
      width: 80,
    },
    { colKey: 'object_name', title: '存储名称', minWidth: 200 },
    { colKey: 'origin_name', title: '原文件名', minWidth: 150 },
    { colKey: 'storage_mode', title: '存储模式', width: 120 },
    { colKey: 'mime_type', title: '资源类型', minWidth: 130 },
    { colKey: 'size_info', title: '文件大小', width: 130 },
    { colKey: 'created_at', title: '上传时间', width: 180 },
    {
      align: 'center',
      colKey: 'action',
      fixed: 'right',
      title: '操作',
      width: 120,
    },
  ];
}

export function createAttachmentColumnOptions(
  columns: AttachmentTableColumn[],
): AttachmentColumnOptionItem[] {
  return columns
    .filter((column) => column.colKey !== 'row-select' && column.title)
    .map((column) => ({
      label: String(column.title),
      value: String(column.colKey),
    }));
}
