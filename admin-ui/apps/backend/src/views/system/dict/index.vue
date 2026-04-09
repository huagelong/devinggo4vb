<script lang="ts" setup>
import type { DictApi } from '#/api/system/dict';
import type { DictOption } from '#/composables/crud/use-dict-options';

import { computed, onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';

import { message } from '#/adapter/tdesign';
import {
  changeDictTypeStatus,
  deleteDictType,
  realDeleteDictType,
  recoveryDictType,
} from '#/api/system/dict';
import CrudToolbar from '#/components/crud/crud-toolbar.vue';
import { useDictOptions } from '#/composables/crud/use-dict-options';

import {
  AddIcon,
  DeleteIcon,
  EditIcon,
  SearchIcon,
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

import DictDataPanel from './components/dict-data-panel.vue';
import DictTypeModal from './components/dict-type-modal.vue';
import type { DictTableColumn, DictTypeListItem } from './model';
import {
  createDictTypeColumnOptions,
  createDictTypeTableColumns,
} from './schemas';
import { useDictTypeCrud } from './use-dict-type-crud';

defineOptions({ name: 'SystemDict' });

type DictTypeModalInstance = {
  open: (data?: DictApi.DictTypeSubmitPayload) => void;
};

type DictDataPanelInstance = {
  open: (row: DictApi.DictTypeItem) => void;
};

const dictTypeModalRef = ref<DictTypeModalInstance>();
const dictDataPanelRef = ref<DictDataPanelInstance>();
const statusOptions = ref<DictOption[]>([]);

const columns: DictTableColumn[] = createDictTypeTableColumns();
const columnOptions = createDictTypeColumnOptions(columns);
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
} = useDictTypeCrud();

const { getDictOptions } = useDictOptions();

function toIds(keys: Array<number | string>) {
  return keys.map((key) => Number(key));
}

async function fetchStatusOptions() {
  statusOptions.value =
    (await getDictOptions('data_status')) || [
      { label: '正常', value: 1 },
      { label: '停用', value: 2 },
    ];
}

function handleAdd() {
  dictTypeModalRef.value?.open();
}

function handleEdit(row: DictTypeListItem) {
  dictTypeModalRef.value?.open({
    ...row,
    status: Number(row.status ?? 1),
  });
}

function handleOpenData(row: DictTypeListItem) {
  dictDataPanelRef.value?.open({
    code: row.code,
    id: row.id,
    name: row.name,
  } as DictApi.DictTypeItem);
}

async function handleDelete(row: DictTypeListItem) {
  try {
    await (isRecycleBin.value ? realDeleteDictType([row.id]) : deleteDictType([row.id]));
    message.success('操作成功');
    await fetchTableData();
  } catch (error) {
    console.error(error);
    message.error('删除失败，请稍后重试');
  }
}

async function handleBatchDelete() {
  if (selectedRowKeys.value.length === 0) {
    message.warning('请选择需要操作的数据');
    return;
  }
  const ids = toIds(selectedRowKeys.value);
  try {
    await (isRecycleBin.value ? realDeleteDictType(ids) : deleteDictType(ids));
    message.success('操作成功');
    clearSelectedRowKeys();
    await fetchTableData();
  } catch (error) {
    console.error(error);
    message.error('批量删除失败，请稍后重试');
  }
}

async function handleRecovery(row: DictTypeListItem) {
  try {
    await recoveryDictType([row.id]);
    message.success('恢复成功');
    await fetchTableData();
  } catch (error) {
    console.error(error);
    message.error('恢复失败，请稍后重试');
  }
}

async function handleBatchRecovery() {
  if (selectedRowKeys.value.length === 0) {
    message.warning('请选择需要操作的数据');
    return;
  }
  const ids = toIds(selectedRowKeys.value);
  try {
    await recoveryDictType(ids);
    message.success('恢复成功');
    clearSelectedRowKeys();
    await fetchTableData();
  } catch (error) {
    console.error(error);
    message.error('批量恢复失败，请稍后重试');
  }
}

async function handleStatusChange(row: DictTypeListItem, checked: boolean) {
  try {
    await changeDictTypeStatus({ id: row.id, status: checked ? 1 : 2 });
    message.success('状态更新成功');
    await fetchTableData();
  } catch (error) {
    console.error(error);
    message.error('状态更新失败，请稍后重试');
  }
}

function handleStatusSwitchChange(row: DictTypeListItem, value: unknown) {
  void handleStatusChange(row, Boolean(value));
}

function handleSuccess() {
  void fetchTableData();
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
            <FormItem label="字典名称" name="name">
              <Input
                v-model="searchForm.name"
                placeholder="请输入字典名称"
                clearable
              />
            </FormItem>
            <FormItem label="字典标识" name="code">
              <Input
                v-model="searchForm.code"
                placeholder="请输入字典标识"
                clearable
              />
            </FormItem>
            <FormItem label="״̬" name="status">
              <Select
                v-model="searchForm.status"
                :options="statusOptions"
                placeholder="请选择状态"
                clearable
                class="w-full"
              />
            </FormItem>
            <FormItem label="创建时间" name="created_at">
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
          table-layout="fixed"
          @page-change="handlePageChange"
          @select-change="handleSelectChange"
        >
          <template #code="{ row }">
            <Button v-if="!isRecycleBin" theme="primary" variant="text" @click="handleOpenData(row)">
              {{ row.code }}
            </Button>
            <span v-else>{{ row.code }}</span>
          </template>

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
                  theme="default"
                  variant="outline"
                  @click="handleOpenData(row)"
                >
                  字典数据
                </Button>
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
                  content="确认删除该字典吗？"
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
                  content="确认恢复该字典吗？"
                  @confirm="handleRecovery(row)"
                >
                  <Button size="small" theme="primary" variant="outline">
                    恢复
                  </Button>
                </Popconfirm>
                <Popconfirm
                  content="确认彻底删除该字典吗？"
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

    <DictTypeModal ref="dictTypeModalRef" @success="handleSuccess" />
    <DictDataPanel ref="dictDataPanelRef" />
  </Page>
</template>
