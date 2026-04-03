import type {
  DataMaintainColumnOptionItem,
  DataMaintainSearchFormModel,
  DataMaintainTableColumn,
} from './model';

export function createDataMaintainSearchForm(): DataMaintainSearchFormModel {
  return {
    group_name: 'default',
    name: '',
  };
}

export function createDataMaintainTableColumns(): DataMaintainTableColumn[] {
  return [
    { colKey: 'name', title: '表名', minWidth: 220 },
    { colKey: 'comment', title: '表注释', minWidth: 220 },
    { colKey: 'engine', title: '引擎', width: 140 },
    { colKey: 'collation', title: '字符集', width: 160 },
    { colKey: 'rows', title: '行数', width: 120 },
    { colKey: 'create_time', title: '创建时间', minWidth: 180 },
  ];
}

export function createDataMaintainColumnOptions(
  columns: DataMaintainTableColumn[],
): DataMaintainColumnOptionItem[] {
  return columns
    .filter((column) => column.title)
    .map((column) => ({
      label: String(column.title),
      value: String(column.colKey),
    }));
}
