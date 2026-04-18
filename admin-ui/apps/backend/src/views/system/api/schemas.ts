import type {
  ApiColumnFormModel,
  ApiColumnOptionItem,
  ApiColumnSearchFormModel,
  ApiColumnTableColumn,
  ApiFormModel,
  ApiSearchFormModel,
  ApiTableColumn,
} from './model';

import { $t } from '@vben/locales';

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
    { colKey: 'group_name', title: $t('system.api.group'), minWidth: 160 },
    { colKey: 'name', title: $t('system.api.name'), minWidth: 200 },
    { colKey: 'access_name', title: $t('system.api.code'), minWidth: 220 },
    { colKey: 'request_mode', title: $t('system.api.requestMode'), width: 120 },
    { colKey: 'auth_mode', title: $t('system.api.authMode'), width: 120 },
    { colKey: 'status', title: $t('common.status'), width: 100, align: 'center' },
    { colKey: 'remark', title: $t('common.remark'), minWidth: 200 },
    { colKey: 'created_at', title: $t('common.createTime'), minWidth: 180 },
    {
      align: 'center',
      colKey: 'action',
      fixed: 'right',
      title: $t('common.action'),
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
    { colKey: 'name', title: $t('system.api.fieldName'), minWidth: 220 },
    { colKey: 'data_type', title: $t('system.api.dataType'), width: 140 },
    { colKey: 'type', title: $t('system.api.fieldType'), width: 120 },
    { colKey: 'is_required', title: $t('system.api.isRequired'), width: 100, align: 'center' },
    { colKey: 'status', title: $t('common.status'), width: 100, align: 'center' },
    { colKey: 'default_value', title: $t('system.api.defaultValue'), minWidth: 180 },
    { colKey: 'remark', title: $t('common.remark'), minWidth: 200 },
    { colKey: 'created_at', title: $t('common.createTime'), minWidth: 180 },
    {
      align: 'center',
      colKey: 'action',
      fixed: 'right',
      title: $t('common.action'),
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
