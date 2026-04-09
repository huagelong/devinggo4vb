<script lang="ts" setup>
import type { DictApi } from '#/api/system/dict';

import { reactive, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import { MessagePlugin } from 'tdesign-vue-next';
import { AddIcon, DeleteIcon, SearchIcon } from 'tdesign-icons-vue-next';
import {
  Button,
  DateRangePicker,
  Form,
  FormItem,
  Input,
  InputNumber,
  Popconfirm,
  Select,
  Space,
  Switch,
  Table,
} from 'tdesign-vue-next';

import {
  changeDictDataStatus,
  deleteDictData,
  getDictDataPageList,
  getRecycleDictDataList,
  realDeleteDictData,
  recoveryDictData,
  updateDictDataNumber,
} from '#/api/system/dict';

import DictDataFormModal from './dict-data-form-modal.vue';
import type { DictDataListItem } from '../model';
import { createDictDataSearchForm, createDictDataTableColumns } from '../schemas';

type DictDataFormInstance = {
  open: (options: {
    data?: DictApi.DictDataSubmitPayload;
    typeInfo: { id: number; code: string; name?: string };
  }) => void;
};

const dictDataModalRef = ref<DictDataFormInstance>();

const searchForm = reactive(createDictDataSearchForm());
const tableData = ref<DictDataListItem[]>([]);
const loading = ref(false);
const isRecycleBin = ref(false);
const selectedRowKeys = ref<Array<number | string>>([]);
const currentDict = ref<{ id: number; code: string; name?: string }>();
const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0,
  showJumper: true,
  showPageSize: true,
});

const columns = createDictDataTableColumns();

function toIds(keys: Array<number | string>) {
  return keys.map((key) => Number(key));
}

function buildParams() {
  const params: DictApi.DictDataQuery = {
    page: pagination.current,
    pageSize: pagination.pageSize,
  };
  if (!currentDict.value) return params;
  params.type_id = currentDict.value.id;
  params.code = currentDict.value.code;
  if (searchForm.label) params.label = searchForm.label;
  if (searchForm.value) params.value = searchForm.value;
  if (searchForm.status !== undefined) params.status = searchForm.status;
  if (searchForm.created_at?.length === 2 && searchForm.created_at[0]) {
    params.created_at = searchForm.created_at;
  }
  return params;
}

async function fetchTableData() {
  if (!currentDict.value) return;
  loading.value = true;
  try {
    const params = buildParams();
    const response = isRecycleBin.value
      ? await getRecycleDictDataList(params)
      : await getDictDataPageList(params);
    tableData.value = response.items ?? [];
    pagination.total = Number(response?.pageInfo?.total || response?.total || 0);
  } catch (error) {
    console.error(error);
    MessagePlugin.error('字典数据加载失败，请稍后重试');
  } finally {
    loading.value = false;
  }
}

function handleSearch() {
  pagination.current = 1;
  void fetchTableData();
}

function handleReset() {
  Object.assign(searchForm, createDictDataSearchForm(), {
    type_id: currentDict.value?.id,
    code: currentDict.value?.code,
  });
  pagination.current = 1;
  void fetchTableData();
}

function handlePageChange(pageInfo: { current: number; pageSize: number }) {
  pagination.current = pageInfo.current;
  pagination.pageSize = pageInfo.pageSize;
  void fetchTableData();
}

function handleSelectChange(keys: Array<number | string>) {
  selectedRowKeys.value = keys;
}

function clearSelectedRowKeys() {
  selectedRowKeys.value = [];
}

function handleToggleRecycleBin(next?: boolean) {
  isRecycleBin.value = typeof next === 'boolean' ? next : !isRecycleBin.value;
  pagination.current = 1;
  clearSelectedRowKeys();
  void fetchTableData();
}

function handleAdd() {
  if (!currentDict.value) return;
  dictDataModalRef.value?.open({
    typeInfo: currentDict.value,
  });
}

function handleEdit(row: DictDataListItem) {
  if (!currentDict.value) return;
  dictDataModalRef.value?.open({
    data: {
      ...row,
      sort: Number(row.sort ?? 1),
      status: Number(row.status ?? 1),
      type_id: currentDict.value.id,
      code: currentDict.value.code,
    },
    typeInfo: currentDict.value,
  });
}

async function handleDelete(row: DictDataListItem) {
  try {
    await (isRecycleBin.value ? realDeleteDictData([row.id]) : deleteDictData([row.id]));
    MessagePlugin.success('操作成功');
    await fetchTableData();
  } catch (error) {
    console.error(error);
    MessagePlugin.error('删除失败，请稍后重试');
  }
}

async function handleBatchDelete() {
  if (selectedRowKeys.value.length === 0) {
    MessagePlugin.warning('请选择需要操作的数据');
    return;
  }
  const ids = toIds(selectedRowKeys.value);
  try {
    await (isRecycleBin.value ? realDeleteDictData(ids) : deleteDictData(ids));
    MessagePlugin.success('操作成功');
    clearSelectedRowKeys();
    await fetchTableData();
  } catch (error) {
    console.error(error);
    MessagePlugin.error('批量删除失败，请稍后重试');
  }
}

async function handleRecovery(row: DictDataListItem) {
  try {
    await recoveryDictData([row.id]);
    MessagePlugin.success('恢复成功');
    await fetchTableData();
  } catch (error) {
    console.error(error);
    MessagePlugin.error('恢复失败，请稍后重试');
  }
}

async function handleBatchRecovery() {
  if (selectedRowKeys.value.length === 0) {
    MessagePlugin.warning('请选择需要操作的数据');
    return;
  }
  const ids = toIds(selectedRowKeys.value);
  try {
    await recoveryDictData(ids);
    MessagePlugin.success('恢复成功');
    clearSelectedRowKeys();
    await fetchTableData();
  } catch (error) {
    console.error(error);
    MessagePlugin.error('批量恢复失败，请稍后重试');
  }
}

async function handleStatusChange(row: DictDataListItem, checked: boolean) {
  try {
    await changeDictDataStatus({ id: row.id, status: checked ? 1 : 2 });
    MessagePlugin.success('状态更新成功');
    await fetchTableData();
  } catch (error) {
    console.error(error);
    MessagePlugin.error('状态更新失败，请稍后重试');
  }
}

function handleStatusSwitchChange(row: DictDataListItem, value: unknown) {
  void handleStatusChange(row, Boolean(value));
}

async function handleSortChange(value: number | string, row: DictDataListItem) {
  const numberValue = Number(value);
  if (Number.isNaN(numberValue)) return;
  try {
    await updateDictDataNumber({
      id: Number(row.id),
      numberName: 'sort',
      numberValue,
    });
    MessagePlugin.success('排序更新成功');
    await fetchTableData();
  } catch (error) {
    console.error(error);
    MessagePlugin.error('排序更新失败，请稍后重试');
  }
}

const [PanelModal, panelModalApi] = useVbenModal({
  footer: false,
  class: 'w-[1200px]',
  onClosed: () => {
    currentDict.value = undefined;
    tableData.value = [];
    clearSelectedRowKeys();
  },
});

async function open(row: DictApi.DictTypeItem) {
  currentDict.value = { id: row.id, code: row.code, name: row.name };
  Object.assign(searchForm, createDictDataSearchForm(), {
    type_id: row.id,
    code: row.code,
  });
  pagination.current = 1;
  pagination.total = 0;
  selectedRowKeys.value = [];
  isRecycleBin.value = false;

  panelModalApi.setState({
    title: `维护「${row.name}」字典数据`,
  });
  panelModalApi.open();
  await fetchTableData();
}

defineExpose({
  open,
});
</script>

<template>
  <PanelModal>
    <div class="flex flex-col gap-4">
      <div class="rounded-md bg-gray-50 p-3 text-sm text-gray-700">
        当前字典：{{ currentDict ? currentDict.name : '-' }}（{{
          currentDict ? currentDict.code : '-'
        }}）
      </div>

      <div class="rounded-md bg-white p-4">
        <Form :data="searchForm" label-width="90px" colon>
          <div class="grid grid-cols-4 gap-x-4">
            <FormItem label="字典标签" name="label">
              <Input
                v-model="searchForm.label"
                placeholder="请输入字典标签"
                clearable
              />
            </FormItem>
            <FormItem label="字典键值" name="value">
              <Input
                v-model="searchForm.value"
                placeholder="请输入字典键值"
                clearable
              />
            </FormItem>
            <FormItem label="״̬" name="status">
              <Select
                v-model="searchForm.status"
                :options="[
                  { label: '正常', value: 1 },
                  { label: '停用', value: 2 },
                ]"
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

      <div class="rounded-md bg-white p-4">
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

          <Button variant="outline" @click="handleToggleRecycleBin()">
            {{ isRecycleBin ? '返回列表' : '查看回收站' }}
          </Button>
        </div>

        <Table
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
          <template #sort="{ row }">
            <InputNumber
              :value="row.sort"
              :min="0"
              :max="1000"
              size="small"
              :disabled="isRecycleBin"
              @change="(value: number | string) => handleSortChange(value, row)"
            />
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
                <Button size="small" theme="primary" variant="outline" @click="handleEdit(row)">
                  编辑
                </Button>
                <Popconfirm content="确认删除该数据吗？" @confirm="handleDelete(row)">
                  <Button size="small" theme="danger" variant="outline">删除</Button>
                </Popconfirm>
              </template>
              <template v-else>
                <Popconfirm content="确认恢复该数据吗？" @confirm="handleRecovery(row)">
                  <Button size="small" theme="primary" variant="outline">恢复</Button>
                </Popconfirm>
                <Popconfirm content="确认彻底删除该数据吗？" @confirm="handleDelete(row)">
                  <Button size="small" theme="danger" variant="outline">彻底删除</Button>
                </Popconfirm>
              </template>
            </div>
          </template>
        </Table>
      </div>

      <DictDataFormModal ref="dictDataModalRef" @success="fetchTableData" />
    </div>
  </PanelModal>
</template>
