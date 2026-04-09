<script lang="ts" setup>
import { computed, onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { message } from '#/adapter/tdesign';

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

import CrudToolbar from '#/components/crud/crud-toolbar.vue';
import {
  changeApiGroupStatus,
  deleteApiGroup,
  realDeleteApiGroup,
  recoveryApiGroup,
} from '#/api/system/api-group';
import type { DictOption } from '#/composables/crud/use-dict-options';
import { useDictOptions } from '#/composables/crud/use-dict-options';

import ApiGroupModal from './components/api-group-modal.vue';
import type { ApiGroupListItem, ApiGroupTableColumn } from './model';
import {
  createApiGroupColumnOptions,
  createApiGroupTableColumns,
} from './schemas';
import { useApiGroupCrud } from './use-api-group-crud';

defineOptions({ name: 'SystemApiGroup' });

type ApiGroupModalInstance = {
  open: (data?: Partial<ApiGroupListItem>) => void;
};

const apiGroupModalRef = ref<ApiGroupModalInstance>();

const columns: ApiGroupTableColumn[] = createApiGroupTableColumns();
const columnOptions = createApiGroupColumnOptions(columns);
const allColumnKeys = columnOptions.map((item) => item.value);
const visibleColumns = ref<string[]>([...allColumnKeys]);

const displayColumns = computed({
  get: () => ['row-select', ...visibleColumns.value],
  set: (value: string[]) => {
    visibleColumns.value = value.filter((item) => item !== 'row-select');
  },
});

const statusOptions = ref<DictOption[]>([]);
const fallbackStatusOptions: DictOption[] = [
  { label: '正常', value: 1 },
  { label: '停用', value: 2 },
];

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
} = useApiGroupCrud();

const { getDictOptions } = useDictOptions();

function toIds(keys: Array<number | string>) {
  return keys.map((key) => Number(key)).filter((id) => !Number.isNaN(id));
}

async function fetchStatusOptions() {
  try {
    const options = await getDictOptions('data_status');
    statusOptions.value =
      options && options.length > 0 ? options : fallbackStatusOptions;
  } catch (error) {
    console.error(error);
    statusOptions.value = fallbackStatusOptions;
  }
}

function handleAdd() {
  apiGroupModalRef.value?.open();
}

function handleEdit(row: ApiGroupListItem) {
  apiGroupModalRef.value?.open(row);
}

async function handleDelete(row: ApiGroupListItem) {
  try {
    const id = Number(row.id);
    await (isRecycleBin.value ? realDeleteApiGroup([id]) : deleteApiGroup([id]));
    message.success('操作成功');
    await fetchTableData();
  } catch (error) {
    console.error(error);
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
    await (isRecycleBin.value ? realDeleteApiGroup(ids) : deleteApiGroup(ids));
    message.success('操作成功');
    clearSelectedRowKeys();
    await fetchTableData();
  } catch (error) {
    console.error(error);
    message.error(isRecycleBin.value ? '批量彻底删除失败，请稍后重试' : '批量删除失败，请稍后重试');
  }
}

async function handleRecovery(row: ApiGroupListItem) {
  try {
    await recoveryApiGroup([Number(row.id)]);
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
  if (ids.length === 0) {
    message.warning('所选数据格式异常，请重试');
    return;
  }
  try {
    await recoveryApiGroup(ids);
    message.success('恢复成功');
    clearSelectedRowKeys();
    await fetchTableData();
  } catch (error) {
    console.error(error);
    message.error('批量恢复失败，请稍后重试');
  }
}

async function handleStatusChange(row: ApiGroupListItem, checked: boolean) {
  const status = checked ? 1 : 2;
  try {
    await changeApiGroupStatus({ id: Number(row.id), status });
    message.success('状态更新成功');
    await fetchTableData();
  } catch (error) {
    console.error(error);
    message.error('状态更新失败，请稍后重试');
  }
}

function handleStatusSwitchChange(row: ApiGroupListItem, value: unknown) {
  void handleStatusChange(row, Boolean(value));
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
          <div class="grid grid-cols-3 gap-x-4">
            <FormItem label="分组名称" name="name">
              <Input
                v-model="searchForm.name"
                placeholder="请输入分组名称"
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
          @page-change="handlePageChange"
          @select-change="handleTableSelectChange"
        >
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
            <span class="block max-w-[260px] truncate">
              {{ row.remark || '-' }}
            </span>
          </template>

          <template #created_at="{ row }">
            {{ row.created_at || '-' }}
          </template>

          <template #updated_at="{ row }">
            {{ row.updated_at || '-' }}
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
                  content="确认删除该分组吗？"
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
                  content="确认恢复该分组吗？"
                  @confirm="handleRecovery(row)"
                >
                  <Button size="small" theme="primary" variant="outline">
                    恢复
                  </Button>
                </Popconfirm>
                <Popconfirm
                  content="确认彻底删除该分组吗？"
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

    <ApiGroupModal ref="apiGroupModalRef" @success="fetchTableData" />
  </Page>
</template>
