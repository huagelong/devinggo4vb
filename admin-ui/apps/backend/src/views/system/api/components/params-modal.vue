<script lang="ts" setup>
import type { ApiColumnListItem, ApiColumnType, ApiListItem } from '../model';
import type { DictOption } from '#/composables/crud/use-dict-options';

import { computed, nextTick, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';
import { $t } from '@vben/locales';
import { message } from '#/adapter/tdesign';
import { logger } from '#/utils/logger';

import {
  AddIcon,
  DeleteIcon,
  DownloadIcon,
  SearchIcon,
  UploadIcon,
} from 'tdesign-icons-vue-next';
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

import type { ApiColumnApi } from '#/api/system/api-column';
import {
  changeApiColumnStatus,
  deleteApiColumn,
  downloadApiColumnTemplate,
  exportApiColumnList,
  importApiColumnFile,
  realDeleteApiColumn,
  recoveryApiColumn,
} from '#/api/system/api-column';
import { useDictOptions } from '#/composables/crud/use-dict-options';
import { downloadResponseBlob } from '#/utils/download';

import ApiColumnModal from './api-column-modal.vue';
import { useApiColumnCrud } from '../use-api-column-crud';
import {
  createApiColumnColumnOptions,
  createApiColumnTableColumns,
} from '../schemas';

const typeLabelMap: Record<ApiColumnType, string> = {
  1: '请求参数',
  2: '响应参数',
};

type ColumnModalInstance = {
  open: (payload: { apiId: number; type: ApiColumnType; data?: ApiColumnListItem }) => void;
};

const apiColumnModalRef = ref<ColumnModalInstance>();
const currentApi = ref<ApiListItem | null>(null);
const currentType = ref<ApiColumnType>(1);

const importInputRef = ref<HTMLInputElement>();
const importLoading = ref(false);
const exportLoading = ref(false);
const templateLoading = ref(false);

const dataTypeOptions = ref<DictOption[]>([]);
const statusOptions = ref<DictOption[]>([]);
const requiredOptions: DictOption[] = [
  { label: '否', value: 1 },
  { label: '是', value: 2 },
];

const { getDictOptions } = useDictOptions();

const columns = createApiColumnTableColumns();
const columnOptions = createApiColumnColumnOptions(columns);
const allColumnKeys = columnOptions.map((item) => item.value);
const visibleColumns = ref<string[]>([...allColumnKeys]);

const displayColumns = computed({
  get: () => ['row-select', ...visibleColumns.value],
  set: (value: string[]) => {
    visibleColumns.value = value.filter((item) => item !== 'row-select');
  },
});

const modalTitle = computed(() => {
  if (!currentApi.value) return '参数管理';
  return `${currentApi.value.name} - ${typeLabelMap[currentType.value]}`;
});

const {
  buildRequestParams,
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
  currentApiId,
  setContext,
} = useApiColumnCrud();

const [Modal, modalApi] = useVbenModal({
  class: 'w-[1100px]',
  showFooter: false,
} as any);

async function fetchFilterOptions() {
  try {
    const [dataTypes, statuses] = await Promise.all([
      getDictOptions('api_data_type'),
      getDictOptions('data_status'),
    ]);
    dataTypeOptions.value = dataTypes;
    statusOptions.value = statuses;
  } catch (error) {
    logger.error(error);
    message.error($t('common.filterLoadFailed'));
  }
}

function toIds(keys: Array<number | string>) {
  return keys.map((key) => Number(key)).filter((id) => !Number.isNaN(id));
}

function openModal(row: ApiListItem, type: ApiColumnType) {
  currentApi.value = row;
  currentType.value = type;
  setContext(row.id, type);
  modalApi.setState({ title: modalTitle.value });
  modalApi.open();
  void nextTick(async () => {
    await fetchFilterOptions();
    clearSelectedRowKeys();
    if (isRecycleBin.value) {
      toggleRecycleBin(false);
    } else {
      handleReset();
    }
  });
}

function handleToggleRecycleBin() {
  toggleRecycleBin(!isRecycleBin.value);
}

function handleAdd() {
  if (!currentApi.value) return;
  apiColumnModalRef.value?.open({
    apiId: currentApi.value.id,
    type: currentType.value,
  });
}

function handleEdit(row: ApiColumnListItem) {
  if (!currentApi.value) return;
  apiColumnModalRef.value?.open({
    apiId: currentApi.value.id,
    type: currentType.value,
    data: row,
  });
}

async function handleDelete(row: ApiColumnListItem) {
  try {
    await (isRecycleBin.value
      ? realDeleteApiColumn([row.id])
      : deleteApiColumn([row.id]));
    message.success($t('common.operationSuccess'));
    await fetchTableData();
  } catch (error) {
    logger.error(error);
    message.error(isRecycleBin.value ? '彻底删除失败，请稍后重试' : '删除失败，请稍后重试');
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
    await (isRecycleBin.value ? realDeleteApiColumn(ids) : deleteApiColumn(ids));
    message.success($t('common.operationSuccess'));
    clearSelectedRowKeys();
    await fetchTableData();
  } catch (error) {
    logger.error(error);
    message.error($t('common.batchDeleteFailed'));
  }
}

async function handleRecovery(row: ApiColumnListItem) {
  try {
    await recoveryApiColumn([row.id]);
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
    await recoveryApiColumn(ids);
    message.success($t('common.recoverySuccess'));
    clearSelectedRowKeys();
    await fetchTableData();
  } catch (error) {
    logger.error(error);
    message.error($t('common.batchRecoveryFailed'));
  }
}

async function handleStatusChange(row: ApiColumnListItem, checked: boolean) {
  const status = checked ? 1 : 2;
  try {
    await changeApiColumnStatus({ id: row.id, status });
    message.success($t('common.statusUpdateSuccess'));
    await fetchTableData();
  } catch (error) {
    logger.error(error);
    message.error($t('common.statusUpdateFailed'));
  }
}

function triggerImport() {
  importInputRef.value?.click();
}

async function handleImportChange(event: Event) {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  if (!file || !currentApiId.value) return;
  importLoading.value = true;
  try {
    await importApiColumnFile(file, {
      api_id: Number(currentApiId.value),
      type: currentType.value,
    });
    message.success($t('common.importSuccess'));
    await fetchTableData();
  } catch (error) {
    logger.error(error);
    message.error($t('common.importFailed'));
  } finally {
    importLoading.value = false;
    input.value = '';
  }
}

async function handleExport() {
  if (!currentApiId.value) {
    message.warning($t('common.invalidApiInfo'));
    return;
  }
  exportLoading.value = true;
  try {
    const response = await exportApiColumnList(
      buildRequestParams(false) as ApiColumnApi.ListQuery,
    );
    downloadResponseBlob(response, '接口参数列表.xlsx');
    message.success('导出成功');
  } catch (error) {
    logger.error(error);
    message.error('导出失败，请稍后重试');
  } finally {
    exportLoading.value = false;
  }
}

async function handleDownloadTemplate() {
  templateLoading.value = true;
  try {
    const response = await downloadApiColumnTemplate();
    downloadResponseBlob(response, '接口参数导入模板.xlsx');
    message.success('模板下载成功');
  } catch (error) {
    logger.error(error);
    message.error('模板下载失败，请稍后重试');
  } finally {
    templateLoading.value = false;
  }
}

function handleTableSelectChange(keys: Array<number | string>) {
  handleSelectChange(keys);
}

function handleStatusSwitchChange(row: ApiColumnListItem, value: unknown) {
  void handleStatusChange(row, Boolean(value));
}

function resolveTypeLabel(value?: number | string) {
  const numberValue = Number(value);
  if (Number.isNaN(numberValue)) {
    return typeLabelMap[currentType.value];
  }
  return typeLabelMap[numberValue as ApiColumnType] ?? typeLabelMap[currentType.value];
}

defineExpose({
  open: openModal,
});
</script>

<template>
  <Modal :title="modalTitle">
    <div class="flex flex-col gap-3">
      <div class="rounded-md bg-white p-4">
        <Form :data="searchForm" label-width="90px" colon>
          <div class="grid grid-cols-4 gap-x-4">
            <FormItem label="字段名称" name="name">
              <Input
                v-model="searchForm.name"
                placeholder="请输入字段名称"
                clearable
              />
            </FormItem>
            <FormItem label="数据类型" name="data_type">
              <Select
                v-model="searchForm.data_type"
                :options="dataTypeOptions"
                placeholder="请选择数据类型"
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
            <FormItem label="是否必填" name="is_required">
              <Select
                v-model="searchForm.is_required"
                :options="requiredOptions"
                placeholder="请选择"
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
              <Button variant="outline" :loading="importLoading" @click="triggerImport">
                <template #icon><UploadIcon /></template>
                导入
              </Button>
              <Button
                variant="outline"
                :loading="templateLoading"
                @click="handleDownloadTemplate"
              >
                <template #icon><DownloadIcon /></template>
                导入模板
              </Button>
              <Button variant="outline" :loading="exportLoading" @click="handleExport">
                <template #icon><DownloadIcon /></template>
                导出
              </Button>
            </template>
            <template v-else>
              <Button theme="success" @click="handleBatchRecovery">恢复</Button>
              <Button theme="danger" @click="handleBatchDelete">彻底删除</Button>
            </template>
          </Space>

          <Space>
            <Button variant="outline" @click="handleToggleRecycleBin">
              {{ isRecycleBin ? '返回列表' : '回收站' }}
            </Button>
          </Space>
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
          table-layout="auto"
          @page-change="handlePageChange"
          @select-change="handleTableSelectChange"
        >
          <template #type="{ row }">
            {{ resolveTypeLabel(row.type) }}
          </template>

          <template #is_required="{ row }">
            {{ row.is_required === 2 ? '是' : '否' }}
          </template>

          <template #status="{ row }">
            <Switch
              :disabled="isRecycleBin"
              :value="row.status === 1"
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
                  编辑
                </Button>
                <Popconfirm
                  content="确认删除该参数吗？"
                  @confirm="handleDelete(row)"
                >
                  <Button size="small" theme="danger" variant="outline">
                    删除
                  </Button>
                </Popconfirm>
              </template>
              <template v-else>
                <Popconfirm
                  content="确认恢复该参数吗？"
                  @confirm="handleRecovery(row)"
                >
                  <Button size="small" theme="primary" variant="outline">
                    恢复
                  </Button>
                </Popconfirm>
                <Popconfirm
                  content="确认彻底删除该参数吗？"
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

    <input
      ref="importInputRef"
      type="file"
      accept=".xlsx,.xls,.csv"
      class="hidden"
      @change="handleImportChange"
    />

    <ApiColumnModal ref="apiColumnModalRef" @success="fetchTableData" />
  </Modal>
</template>
