<script lang="ts" setup>
import { computed, onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { $t } from '@vben/locales';
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
import type { ApiListItem, ApiTableColumn } from './model';
import {
  createApiTableColumnOptions,
  createApiTableColumns,
} from './schemas';
import { useApiCrud } from './use-api-crud';

defineOptions({ name: 'SystemApi' });

type ApiModalInstance = {
  open: (data?: Partial<ApiListItem>) => void;
};

const apiModalRef = ref<ApiModalInstance>();

const columns: ApiTableColumn[] = createApiTableColumns();
const columnOptions = createApiTableColumnOptions(columns);
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
  { label: $t('system.api.simpleMode'), value: 1 },
  { label: $t('system.api.complexMode'), value: 2 },
];

const fallbackRequestModes: DictOption[] = [
  { label: 'GET', value: 'GET' },
  { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' },
  { label: 'DELETE', value: 'DELETE' },
];

const fallbackStatusOptions: DictOption[] = [
  { label: $t('common.statusEnabled'), value: 1 },
  { label: $t('common.statusDisabled'), value: 2 },
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
    message.error($t('common.filterLoadFailed'));
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
    message.success($t('common.operationSuccess'));
    await fetchTableData();
  } catch (error) {
    logger.error(error);
    message.error($t('common.deleteFailed'));
  }
}

async function handleBatchDelete() {
  if (selectedRowKeys.value.length === 0) {
    message.warning($t('common.selectDataFirst'));
    return;
  }
  const ids = toIds(selectedRowKeys.value);
  if (ids.length === 0) {
    message.warning($t('common.invalidDataFormat'));
    return;
  }
  try {
    await (isRecycleBin.value ? realDeleteApi(ids) : deleteApi(ids));
    message.success($t('common.operationSuccess'));
    clearSelectedRowKeys();
    await fetchTableData();
  } catch (error) {
    logger.error(error);
    message.error($t('common.batchDeleteFailed'));
  }
}

async function handleRecovery(row: ApiListItem) {
  try {
    await recoveryApi([row.id]);
    message.success($t('common.recoverySuccess'));
    await fetchTableData();
  } catch (error) {
    logger.error(error);
    message.error($t('common.recoveryFailed'));
  }
}

async function handleBatchRecovery() {
  if (selectedRowKeys.value.length === 0) {
    message.warning($t('common.selectDataFirst'));
    return;
  }
  const ids = toIds(selectedRowKeys.value);
  if (ids.length === 0) {
    message.warning($t('common.invalidDataFormat'));
    return;
  }
  try {
    await recoveryApi(ids);
    message.success($t('common.recoverySuccess'));
    clearSelectedRowKeys();
    await fetchTableData();
  } catch (error) {
    logger.error(error);
    message.error($t('common.batchRecoveryFailed'));
  }
}

async function handleStatusChange(row: ApiListItem, checked: boolean) {
  const status = checked ? 1 : 2;
  try {
    await changeApiStatus({ id: row.id, status });
    message.success($t('common.statusUpdateSuccess'));
    await fetchTableData();
  } catch (error) {
    logger.error(error);
    message.error($t('common.statusUpdateFailed'));
  }
}

function handleTableSelectChange(keys: Array<number | string>) {
  handleSelectChange(keys);
}

function handleStatusSwitchChange(row: ApiListItem, value: unknown) {
  void handleStatusChange(row, Boolean(value));
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
            <FormItem :label="$t('system.api.group')" name="group_id">
              <Select
                v-model="searchForm.group_id"
                :options="groupOptions"
                :placeholder="$t('ui.placeholder.select')"
                clearable
                class="w-full"
              />
            </FormItem>
            <FormItem :label="$t('system.api.name')" name="name">
              <Input
                v-model="searchForm.name"
                :placeholder="$t('ui.placeholder.input')"
                clearable
              />
            </FormItem>
            <FormItem :label="$t('system.api.code')" name="access_name">
              <Input
                v-model="searchForm.access_name"
                :placeholder="$t('ui.placeholder.input')"
                clearable
              />
            </FormItem>
            <FormItem :label="$t('system.api.requestMode')" name="request_mode">
              <Select
                v-model="searchForm.request_mode"
                :options="requestModeOptions"
                :placeholder="$t('ui.placeholder.select')"
                clearable
                class="w-full"
              />
            </FormItem>
            <FormItem :label="$t('common.status')" name="status">
              <Select
                v-model="searchForm.status"
                :options="statusOptions"
                :placeholder="$t('ui.placeholder.select')"
                clearable
                class="w-full"
              />
            </FormItem>
            <FormItem :label="$t('common.createTime')" name="created_at" class="col-span-2">
              <DateRangePicker
                v-model="searchForm.created_at"
                :placeholder="[$t('common.startTime'), $t('common.endTime')]"
                clearable
                class="w-full"
              />
            </FormItem>
          </div>
          <div class="flex justify-end gap-2 pt-2">
            <Button theme="default" @click="handleReset">{{ $t('common.reset') }}</Button>
            <Button theme="primary" @click="handleSearch">
              <template #icon><SearchIcon /></template>
              {{ $t('common.query') }}
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
                {{ $t('common.create') }}
              </Button>
              <Button theme="danger" variant="outline" @click="handleBatchDelete">
                <template #icon><DeleteIcon /></template>
                {{ $t('common.delete') }}
              </Button>
            </template>
            <template v-else>
              <Button theme="success" @click="handleBatchRecovery">{{ $t('common.recovery') }}</Button>
              <Button theme="danger" @click="handleBatchDelete">{{ $t('common.permanentDelete') }}</Button>
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
            {{ resolveGroupLabel(row?.group_id) }}
          </template>

          <template #request_mode="{ row }">
            {{ resolveRequestModeLabel(row?.request_mode) }}
          </template>

          <template #auth_mode="{ row }">
            {{ resolveAuthModeLabel(row?.auth_mode) }}
          </template>

          <template #status="{ row }">
            <div class="flex items-center justify-center">
              <Switch
                :disabled="isRecycleBin"
                :value="row?.status === 1"
                @change="(value: unknown) => handleStatusSwitchChange(row, value)"
              />
            </div>
          </template>

          <template #remark="{ row }">
            <span class="block max-w-[240px] truncate">
              {{ row?.remark || '-' }}
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
                  {{ $t('common.edit') }}
                </Button>
                <Popconfirm
                  :content="$t('system.api.confirmDelete')"
                  @confirm="handleDelete(row)"
                >
                  <Button size="small" theme="danger" variant="outline">
                    <template #icon><DeleteIcon /></template>
                    {{ $t('common.delete') }}
                  </Button>
                </Popconfirm>
              </template>

              <template v-else>
                <Popconfirm
                  :content="$t('system.api.confirmRecovery')"
                  @confirm="handleRecovery(row)"
                >
                  <Button size="small" theme="primary" variant="outline">
                    {{ $t('common.recovery') }}
                  </Button>
                </Popconfirm>
                <Popconfirm
                  :content="$t('system.api.confirmPermanentDelete')"
                  @confirm="handleDelete(row)"
                >
                  <Button size="small" theme="danger" variant="outline">
                    {{ $t('common.permanentDelete') }}
                  </Button>
                </Popconfirm>
              </template>
            </div>
          </template>
        </Table>
      </div>
    </div>

    <ApiModal ref="apiModalRef" @success="fetchTableData" />
  </Page>
</template>
