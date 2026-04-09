import type {
  CodeColumnOptionItem,
  CodeSearchFormModel,
  CodeTableColumn,
} from './model';

export const generateTypeOptions = [
  { label: '单表CRUD', value: 'single' },
  { label: '树表CRUD', value: 'tree' },
];

export const componentTypeOptions = [
  { label: '模态框', value: 1 },
  { label: '抽屉', value: 2 },
  { label: 'Tagҳ', value: 3 },
];

export const tplTypeOptions = [
  { label: 'default', value: 'default' },
  { label: 'ruoyi', value: 'ruoyi' },
];

export const queryTypeOptions = [
  { label: '=', value: 'eq' },
  { label: '!=', value: 'neq' },
  { label: '>', value: 'gt' },
  { label: '>=', value: 'gte' },
  { label: '<', value: 'lt' },
  { label: '<=', value: 'lte' },
  { label: 'LIKE', value: 'like' },
  { label: 'IN', value: 'in' },
  { label: 'NOT IN', value: 'notin' },
  { label: 'BETWEEN', value: 'between' },
];

export const viewTypeOptions = [
  { label: '文本框', value: 'text' },
  { label: '密码框', value: 'password' },
  { label: '文本域', value: 'textarea' },
  { label: '数字输入', value: 'inputNumber' },
  { label: '开关', value: 'switch' },
  { label: '滑块', value: 'slider' },
  { label: '下拉选择', value: 'select' },
  { label: '树形选择', value: 'treeSelect' },
  { label: '单选框', value: 'radio' },
  { label: '复选框', value: 'checkbox' },
  { label: '日期选择', value: 'date' },
  { label: '时间选择', value: 'time' },
  { label: '评分', value: 'rate' },
  { label: '级联选择', value: 'cascader' },
  { label: '穿梭框', value: 'transfer' },
  { label: '用户选择', value: 'selectUser' },
  { label: '城市联动', value: 'cityLinkage' },
  { label: '上传组件', value: 'upload' },
  { label: '富文本', value: 'editor' },
  { label: '代码编辑器', value: 'codeEditor' },
];

export const menuButtonOptions = [
  { label: '新增(save)', value: 'save' },
  { label: '更新(update)', value: 'update' },
  { label: '读取(read)', value: 'read' },
  { label: '删除(delete)', value: 'delete' },
  { label: '回收站(recycle)', value: 'recycle' },
  { label: '状态切换(changeStatus)', value: 'changeStatus' },
  { label: '数字操作(numberOperation)', value: 'numberOperation' },
  { label: '导入(import)', value: 'import' },
  { label: '导出(export)', value: 'export' },
];

export function createCodeSearchForm(): CodeSearchFormModel {
  return {
    table_name: '',
    type: undefined,
  };
}

export function createCodeTableColumns(): CodeTableColumn[] {
  return [
    { colKey: 'row-select', title: '', width: 52, fixed: 'left' },
    { colKey: 'id', title: 'ID', width: 80 },
    { colKey: 'table_name', title: '表名称', minWidth: 200 },
    { colKey: 'table_comment', title: '表描述', minWidth: 200 },
    { colKey: 'type', title: '生成类型', width: 120 },
    { colKey: 'module_name', title: '所属模块', width: 150 },
    { colKey: 'menu_name', title: '菜单名称', width: 150 },
    { colKey: 'created_at', title: '创建时间', width: 180 },
    { colKey: 'action', title: '操作', width: 320, fixed: 'right' },
  ];
}

export function createCodeColumnOptions(
  columns: CodeTableColumn[],
): CodeColumnOptionItem[] {
  return columns
    .filter((column) => column.colKey !== 'row-select' && column.title)
    .map((column) => ({
      label: String(column.title),
      value: String(column.colKey),
    }));
}
