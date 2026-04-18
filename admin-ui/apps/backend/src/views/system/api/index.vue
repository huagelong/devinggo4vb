<script lang="ts" setup>
import { computed, onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { message } from '#/adapter/tdesign';
import { logger } from '#/utils/logger';

import { AddIcon, DeleteIcon, EditIcon, SearchIcon } from 'tdesign-icons-vue-next';
import {
  Button,
  DateRangePicker,
  Form,
  FormItem,
  Input,
  Popconfirm,
  Select,
  Space,
  Switch,
  Table,
} from 'tdesign-vue-next';

import CrudToolbar from '#/components/crud/crud-toolbar.vue';
import {
  changeApiStatus,
  deleteApi,
  realDeleteApi,
  recoveryApi,
} from '#/api/system/api';
import { getApiGroupList } from '#/api/system/api-group';
import type { ApiGroupApi } from '#/api/system/api-group';
import type { OptionItem, IdType } from '#/types/common';
import type { DictOption } from '#/composables/crud/use-dict-options';
import { useDictOptions } from '#/composables/crud/use-dict-options';

import ApiModal from './components/api-modal.vue';
import ParamsModal from './components/params-modal.vue';
import type { ApiListItem, ApiTableColumn, ApiColumnType } from './model';
import {
  createApiColumnOptions,
  createApiTableColumns,
} from './schemas';
import { useApiCrud } from './use-api-crud';

defineOptions({ name: 'SystemApi' });

type ApiModalInstance = {
  open: (data?: Partial<ApiListItem>) => void;
};

type ParamsModalInstance = {
  open: (row: ApiListItem, type: ApiColumnType) => void;
};

const apiModalRef = ref<ApiModalInstance>();
const paramsModalRef = ref<ParamsModalInstance>();

const columns: ApiTableColumn[] = createApiTableColumns();
const columnOptions = createApiColumnOptions(columns);
const allColumnKeys = columnOptions.map((item) => item.value);
const visibleColumns = ref<string[]>([...allColumnKeys]);

const displayColumns = computed({
  get: () => ['row-select', ...visibleColumns.value],
  set: (value: string[]) => {
    visibleColumns.value = value.filter((item) => item !== 'row-select');
  },
});

const groupOptions = ref<OptionItem<IdType>[]>([]);
const requestModeOptions = ref<DictOption[]>([]);
const statusOptions = ref<DictOption[]>([]);

const authModeOptions: DictOption[] = [
  { label: '简易模式', value: 1 },
  { label: '复杂模式', value: 2 },
];

const fallbackRequestModes: DictOption[] = [
  { label: 'GET', value: 'GET' },
  { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' },
  { label: 'DELETE', value: 'DELETE' },
];

const fallbackStatusOptions: DictOption[] = [
  { label: '正常', value: 1 },
  { label: '停用', value: 2 },
];

const groupMap = computed(() => {
  const map = new Map<string, string>();
  groupOptions.value.forEach((item: OptionItem<IdType>) => {
    map.set(String(item.value), item.label ?? '');
  });
  return map;
});

const requestModeMap = computed(() => {
  const map = new Map<string, string>();
  requestModeOptions.value.forEach((item: DictOption) => {
    map.set(String(item.value ?? item.label), item.label ?? '');
  });
  return map;
});

const authModeMap = computed(() => {
  const map = new Map<number, string>();
  authModeOptions.forEach((item: DictOption) => {
    map.set(Number(item.value), item.label ?? '');
  });
  return map;
});

const {
  clearSelectedRowKeys,
  fetchTableData,
  handlePageChange,
  handleReset,
  handleSearch,
  handleSelectChange,
  isRecycleBin,
  loading,
  pagination,
  searchForm,
  selectedRowKeys,
  tableData,
  toggleRecycleBin,
} = useApiCrud();

const { getDictOptions } = useDictOptions();

function toIds(keys: Array<number | string>) {
  return keys.map((key) => Number(key)).filter((id) => !Number.isNaN(id));
}

async function fetchFilterOptions() {
  try {
    const [groups, requestModes, statuses] = await Promise.all([
      getApiGroupList().catch(() => []),
      getDictOptions('request_mode'),
      getDictOptions('data_status'),
    ]);
    const groupList = Array.isArray(groups) ? (groups as ApiGroupApi.ListItem[]) : [];
    groupOptions.value = groupList.map((item) => ({
      label: item.name,
      value: item.id,
    }));
    requestModeOptions.value =
      requestModes && requestModes.length > 0 ? requestModes : fallbackRequestModes;
    statusOptions.value =
      statuses && statuses.length > 0 ? statuses : fallbackStatusOptions;
  } catch (error) {
    logger.error(error);
    message.error('筛选项加载失败，请稍后重试');
    groupOptions.value = [];
    requestModeOptions.value = fallbackRequestModes;
    statusOptions.value = fallbackStatusOptions;
  }
}

function handleAdd() {
  apiModalRef.value?.open();
}

function handleEdit(row: ApiListItem) {
  apiModalRef.value?.open(row);
}

async function handleDelete(row: ApiListItem) {
  try {
    await (isRecycleBin.value ? realDeleteApi([row.id]) : deleteApi([row.id]));
    message.success('操作成功');
    await fetchTableData();
  } catch (error) {
    logger.error(error);
    message.error(isRecycleBin.value ? '彻底删除失败，请稍后重试' : '删除失败，请稍后重试');
  }
}

async function handleBatchDelete() {
  if (selectedRowKeys.value.length === 0) {
    message.warning('请选择需要操作的数据');
    return;
  }
  const ids = toIds(selectedRowKeys.value);
  if (ids.length === 0) {
    message.warning('所选数据格式异常，请重试');
    return;
  }
  try {
    await (isRecycleBin.value ? realDeleteApi(ids) : deleteApi(ids));
    message.success('操作成功');
    clearSelectedRowKeys();
    await fetchTableData();
  } catch (error) {
    logger.error(error);
    message.error(isRecycleBin.value ? '批量彻底删除失败，请稍后重试' : '批量删除失败，请稍后重试');
  }
}

async function handleRecovery(row: ApiListItem) {
  try {
    await recoveryApi([row.id]);
    message.success('恢复成功');
    await fetchTableData();
  } catch (error) {
    logger.error(error);
    message.error('恢复失败，请稍后重试');
  }
}

async function handleBatchRecovery() {
  if (selectedRowKeys.value.length === 0) {
    message.warning('请选择需要操作的数据');
    return;
  }
  const ids = toIds(selectedRowKeys.value);
  if (ids.length === 0) {
    message.warning('所选数据格式异常，请重试');
    return;
  }
  try {
    await recoveryApi(ids);
    message.success('恢复成功');
    clearSelectedRowKeys();
    await fetchTableData();
  } catch (error) {
    logger.error(error);
    message.error('批量恢复失败，请稍后重试');
  }
}

async function handleStatusChange(row: ApiListItem, checked: boolean) {
  const status = checked ? 1 : 2;
  try {
    await changeApiStatus({ id: row.id, status });
    message.success('状态更新成功');
    await fetchTableData();
  } catch (error) {
    logger.error(error);
    message.error('状态更新失败，请稍后重试');
  }
}

function handleTableSelectChange(keys: Array<number | string>) {
  handleSelectChange(keys);
}

function handleStatusSwitchChange(row: ApiListItem, value: unknown) {
  void handleStatusChange(row, Boolean(value));
}

function handleManageParams(row: ApiListItem, type: ApiColumnType) {
  paramsModalRef.value?.open(row, type);
}

function resolveGroupLabel(id?: IdType) {
  if (!id) return '-';
  return groupMap.value.get(String(id)) || '-';
}

function resolveRequestModeLabel(value?: number | string) {
  if (value === undefined || value === null) return '-';
  return requestModeMap.value.get(String(value)) || String(value);
}

function resolveAuthModeLabel(value?: number | string) {
  const numberValue = Number(value);
  if (Number.isNaN(numberValue)) return '-';
  return authModeMap.value.get(numberValue) || '-';
}

onMounted(() => {
  void fetchFilterOptions();
  void fetchTableData();
});
</script>

<template>
  <Page auto-content-height>
    <div class="flex h-full flex-col gap-3">
      <div class="rounded-md bg-white p-4">
        <Form :data="searchForm" label-width="90px" colon>
          <div class="grid grid-cols-4 gap-x-4">
            <FormItem label="所属组" name="group_id">
              <Select
                v-model="searchForm.group_id"
                :options="groupOptions"
                placeholder="请选择所属组"
                clearable
                class="w-full"
              />
            </FormItem>
            <FormItem label="接口名称" name="name">
              <Input
                v-model="searchForm.name"
                placeholder="请输入接口名称"
                clearable
              />
            </FormItem>
            <FormItem label="接口标识" name="access_name">
              <Input
                v-model="searchForm.access_name"
                placeholder="请输入接口标识"
                clearable
              />
            </FormItem>
            <FormItem label="请求模式" name="request_mode">
              <Select
                v-model="searchForm.request_mode"
                :options="requestModeOptions"
                placeholder="请选择请求模式"
                clearable
                class="w-full"
              />
            </FormItem>
            <FormItem label="状态" name="status">
              <Select
                v-model="searchForm.status"
                :options="statusOptions"
                placeholder="请选择状态"
                clearable
                class="w-full"
              />
            </FormItem>
            <FormItem label="创建时间" name="created_at" class="col-span-2">
              <DateRangePicker
                v-model="searchForm.created_at"
                :placeholder="['开始时间', '结束时间']"
                clearable
                class="w-full"
              />
            </FormItem>
          </div>
          <div class="flex justify-end gap-2 pt-2">
            <Button theme="default" @click="handleReset">重置</Button>
            <Button theme="primary" @click="handleSearch">
              <template #icon><SearchIcon /></template>
              查询
            </Button>
          </div>
        </Form>
      </div>

      <div class="flex min-h-0 flex-1 flex-col rounded-md bg-white p-4">
        <div class="mb-3 flex items-center justify-between">
          <Space>
            <template v-if="!isRecycleBin">
              <Button theme="primary" @click="handleAdd">
                <template #icon><AddIcon /></template>
                新增
              </Button>
              <Button theme="danger" variant="outline" @click="handleBatchDelete">
                <template #icon><DeleteIcon /></template>
                删除
              </Button>
            </template>
            <template v-else>
              <Button theme="success" @click="handleBatchRecovery">恢复</Button>
              <Button theme="danger" @click="handleBatchDelete">彻底删除</Button>
            </template>
          </Space>

          <CrudToolbar
            v-model="visibleColumns"
            :column-options="columnOptions"
            :is-recycle-bin="isRecycleBin"
            @refresh="fetchTableData"
            @toggle-recycle="toggleRecycleBin"
          />
        </div>

        <Table
          v-model:display-columns="displayColumns"
          :columns="columns"
          :data="tableData"
          :loading="loading"
          :pagination="pagination"
          :selected-row-keys="selectedRowKeys"
          row-key="id"
          hover
          stripe
          @page-change="handlePageChange"
          @select-change="handleTableSelectChange"
        >
          <template #group_name="{ row }">
            {{ resolveGroupLabel(row.group_id) }}
          </template>

          <template #request_mode="{ row }">
            {{ resolveRequestModeLabel(row.request_mode) }}
          </template>

          <template #auth_mode="{ row }">
            {{ resolveAuthModeLabel(row.auth_mode) }}
          </template>

          <template #status="{ row }">
            <div class="flex items-center justify-center">
              <Switch
                :disabled="isRecycleBin"
                :value="row.status === 1"
                @change="(value: unknown) => handleStatusSwitchChange(row, value)"
              />
            </div>
          </template>

          <template #remark="{ row }">
            <span class="block max-w-[240px] truncate">
              {{ row.remark || '-' }}
            </span>
          </template>

          <template #action="{ row }">
            <div class="flex items-center justify-center gap-1">
              <template v-if="!isRecycleBin">
                <Button
                  size="small"
                  theme="primary"
                  variant="outline"
                  @click="handleEdit(row)"
                >
                  <template #icon><EditIcon /></template>
                  编辑
                </Button>
                <Popconfirm
                  content="确认删除该接口吗？"
                  @confirm="handleDelete(row)"
                >
                  <Button size="small" theme="danger" variant="outline">
                    <template #icon><DeleteIcon /></template>
                    删除
                  </Button>
                </Popconfirm>
                <Button
                  size="small"
                  theme="default"
                  variant="outline"
                  @click="handleManageParams(row, 1)"
                >
                  请求参数
                </Button>
                <Button
                  size="small"
                  theme="default"
                  variant="outline"
                  @click="handleManageParams(row, 2)"
                >
                  响应参数
                </Button>
              </template>

              <template v-else>
                <Popconfirm
                  content="确认恢复该接口吗？"
                  @confirm="handleRecovery(row)"
                >
                  <Button size="small" theme="primary" variant="outline">
                    恢复
                  </Button>
                </Popconfirm>
                <Popconfirm
                  content="确认彻底删除该接口吗？"
                  @confirm="handleDelete(row)"
                >
                  <Button size="small" theme="danger" variant="outline">
                    彻底删除
                  </Button>
                </Popconfirm>
              </template>
            </div>
          </template>
        </Table>
      </div>
    </div>

    <ApiModal ref="apiModalRef" @success="fetchTableData" />
    <ParamsModal ref="paramsModalRef" />
  </Page>
</template>
