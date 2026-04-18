<script lang="ts" setup>
import type { AppApi } from '#/api/system/app';
import type { DictOption } from '#/composables/crud/use-dict-options';

import { computed, onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { message } from '#/adapter/tdesign';
import { logger } from '#/utils/logger';
import {
  changeAppStatus,
  deleteApp,
  realDeleteApp,
  recoveryApp,
} from '#/api/system/app';
import CrudToolbar from '#/components/crud/crud-toolbar.vue';
import { useDictOptions } from '#/composables/crud/use-dict-options';

import AppModal from './components/app-modal.vue';
import type { AppListItem, AppTableColumn } from './model';
import {
  createAppColumnOptions,
  createAppTableColumns,
} from './schemas';
import { useAppCrud } from './use-app-crud';

defineOptions({ name: 'SystemApp' });

type AppModalInstance = {
  open: (data?: Partial<AppApi.SubmitPayload>) => void;
};

const appModalRef = ref<AppModalInstance>();
const fallbackStatusOptions: DictOption[] = [
  { label: '正常', value: 1 },
  { label: '停用', value: 2 },
];

const statusOptions = ref<DictOption[]>([]);

const columns: AppTableColumn[] = createAppTableColumns();
const columnOptions = createAppColumnOptions(columns);
const allColumnKeys = columnOptions.map((item) => item.value);
const visibleColumns = ref<string[]>([...allColumnKeys]);

const displayColumns = computed({
  get: () => ['row-select', ...visibleColumns.value],
  set: (value: string[]) => {
    visibleColumns.value = value.filter((item) => item !== 'row-select');
  },
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
} = useAppCrud();

const { getDictOptions } = useDictOptions();

function toIds(keys: Array<number | string>) {
  return keys.map((key) => Number(key)).filter((id) => !Number.isNaN(id));
}

async function fetchStatusOptions() {
  const options = await getDictOptions('data_status');
  statusOptions.value =
    options.length > 0 ? options : fallbackStatusOptions;
}

function handleAdd() {
  appModalRef.value?.open();
}

function handleEdit(row: AppListItem) {
  appModalRef.value?.open({
    ...row,
    status: Number(row.status ?? 1),
  });
}

async function handleDelete(row: AppListItem) {
  try {
    await (isRecycleBin.value ? realDeleteApp([row.id]) : deleteApp([row.id]));
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
    await (isRecycleBin.value ? realDeleteApp(ids) : deleteApp(ids));
    message.success($t('common.operationSuccess'));
    clearSelectedRowKeys();
    await fetchTableData();
  } catch (error) {
    logger.error(error);
    message.error($t('common.batchDeleteFailed'));
  }
}

async function handleRecovery(row: AppListItem) {
  try {
    await recoveryApp([row.id]);
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
    await recoveryApp(ids);
    message.success($t('common.recoverySuccess'));
    clearSelectedRowKeys();
    await fetchTableData();
  } catch (error) {
    logger.error(error);
    message.error($t('common.batchRecoveryFailed'));
  }
}

async function handleStatusChange(row: AppListItem, checked: boolean) {
  try {
    await changeAppStatus({ id: row.id, status: checked ? 1 : 2 });
    message.success($t('common.statusUpdateSuccess'));
    await fetchTableData();
  } catch (error) {
    logger.error(error);
    message.error($t('common.statusUpdateFailed'));
  }
}

function handleStatusSwitchChange(row: AppListItem, value: unknown) {
  void handleStatusChange(row, Boolean(value));
}

function handleSuccess() {
  void fetchTableData();
}

function handleTableSelectChange(keys: Array<number | string>) {
  handleSelectChange(keys);
}

onMounted(() => {
  void fetchStatusOptions();
  void fetchTableData();
});
</script>

<template>
  <Page auto-content-height>
    <div class="flex h-full flex-col gap-3">
      <div class="rounded-md bg-white p-4">
        <Form :data="searchForm" label-width="90px" colon>
          <div class="grid grid-cols-4 gap-x-4">
            <FormItem label="应用名称" name="name">
              <Input
                v-model="searchForm.name"
                placeholder="请输入应用名称"
                clearable
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
          <template #status="{ row }">
            <Switch
              :disabled="isRecycleBin"
              :value="Number(row.status) === 1"
              @change="(value: unknown) => handleStatusSwitchChange(row, value)"
            />
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
                  content="确认删除该应用吗？"
                  @confirm="handleDelete(row)"
                >
                  <Button size="small" theme="danger" variant="outline">
                    <template #icon><DeleteIcon /></template>
                    删除
                  </Button>
                </Popconfirm>
              </template>
              <template v-else>
                <Popconfirm
                  content="确认恢复该应用吗？"
                  @confirm="handleRecovery(row)"
                >
                  <Button size="small" theme="primary" variant="outline">
                    恢复
                  </Button>
                </Popconfirm>
                <Popconfirm
                  content="确认彻底删除该应用吗？"
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

    <AppModal ref="appModalRef" @success="handleSuccess" />
  </Page>
</template>
