<script lang="ts" setup>
import type { DictOption } from '#/composables/crud/use-dict-options';

import { computed, onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';

import { message } from '#/adapter/tdesign';
import {
  changeDeptStatus,
  deleteDept,
  realDeleteDept,
  recoveryDept,
  updateDeptNumber,
} from '#/api/system/dept';
import CrudToolbar from '#/components/crud/crud-toolbar.vue';
import { useDictOptions } from '#/composables/crud/use-dict-options';

import { DeleteIcon, EditIcon, PlusIcon, SearchIcon, UserIcon } from 'tdesign-icons-vue-next';
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

import DeptLeaderModal from './components/dept-leader-modal.vue';
import DeptModal from './components/dept-modal.vue';
import type { DeptTableColumn, DeptTreeItem } from './model';
import {
  createDeptColumnOptions,
  createDeptTableColumns,
} from './schemas';
import { useDeptPage } from './use-dept-page';

defineOptions({ name: 'SystemDept' });

type DeptModalInstance = {
  open: (data?: Record<string, unknown>) => void;
};

type DeptLeaderModalInstance = {
  open: (data: { id: number; name?: string }) => void;
};

const deptModalRef = ref<DeptModalInstance>();
const deptLeaderModalRef = ref<DeptLeaderModalInstance>();
const statusOptions = ref<DictOption[]>([]);

const columns: DeptTableColumn[] = createDeptTableColumns();
const columnOptions = createDeptColumnOptions(columns);
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
  handleReset,
  handleSearch,
  handleSelectChange,
  isRecycleBin,
  loading,
  searchForm,
  selectedRowKeys,
  tableData,
  toggleRecycleBin,
} = useDeptPage();

const { getDictOptions } = useDictOptions();

async function fetchStatusOptions() {
  const options = await getDictOptions('data_status');
  statusOptions.value =
    options.length > 0
      ? options
      : [
          { label: '正常', value: 1 },
          { label: '停用', value: 2 },
        ];
}

function toIds(keys: Array<number | string>) {
  return keys.map((key) => Number(key));
}

function handleAdd(parentId = 0) {
  deptModalRef.value?.open({ parent_id: parentId });
}

function handleEdit(row: DeptTreeItem) {
  deptModalRef.value?.open({
    ...row,
    status: Number(row.status ?? 1),
  });
}

function handleOpenLeaderList(row: DeptTreeItem) {
  deptLeaderModalRef.value?.open({ id: row.id, name: row.name });
}

async function handleDelete(row: DeptTreeItem) {
  try {
    await (isRecycleBin.value ? realDeleteDept([row.id]) : deleteDept([row.id]));
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
    await (isRecycleBin.value ? realDeleteDept(ids) : deleteDept(ids));
    message.success('操作成功');
    clearSelectedRowKeys();
    await fetchTableData();
  } catch (error) {
    console.error(error);
    message.error('批量删除失败，请稍后重试');
  }
}

async function handleRecovery(row: DeptTreeItem) {
  try {
    await recoveryDept([row.id]);
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
    await recoveryDept(ids);
    message.success('恢复成功');
    clearSelectedRowKeys();
    await fetchTableData();
  } catch (error) {
    console.error(error);
    message.error('批量恢复失败，请稍后重试');
  }
}

async function handleStatusChange(row: DeptTreeItem, checked: boolean) {
  try {
    await changeDeptStatus({ id: row.id, status: checked ? 1 : 2 });
    message.success('状态更新成功');
    await fetchTableData();
  } catch (error) {
    console.error(error);
    message.error('状态更新失败，请稍后重试');
  }
}

async function handleSortChange(value: number | string, row: DeptTreeItem) {
  const numberValue = Number(value);
  if (Number.isNaN(numberValue)) return;

  try {
    await updateDeptNumber({
      id: row.id,
      numberName: 'sort',
      numberValue,
    });
    message.success('排序更新成功');
    await fetchTableData();
  } catch (error) {
    console.error(error);
    message.error('排序更新失败，请稍后重试');
  }
}

function handleSuccess() {
  void fetchTableData();
}

function handleStatusSwitchChange(row: DeptTreeItem, value: unknown) {
  void handleStatusChange(row, Boolean(value));
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
        <Form :data="searchForm" label-width="80px" colon>
          <div class="grid grid-cols-4 gap-x-4">
            <FormItem label="部门名称" name="name">
              <Input
                v-model="searchForm.name"
                placeholder="请输入部门名称"
                clearable
              />
            </FormItem>
            <FormItem label="负责人" name="leader">
              <Input
                v-model="searchForm.leader"
                placeholder="请输入负责人"
                clearable
              />
            </FormItem>
            <FormItem label="手机" name="phone">
              <Input
                v-model="searchForm.phone"
                placeholder="请输入手机号"
                clearable
              />
            </FormItem>
            <FormItem label="״̬" name="status">
              <Select
                v-model="searchForm.status"
                :options="statusOptions"
                placeholder="请选择状态"
                clearable
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
              <Button theme="primary" @click="handleAdd()">
                <template #icon><PlusIcon /></template>
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
          :selected-row-keys="selectedRowKeys"
          :tree="{ childrenKey: 'children', defaultExpandAll: true }"
          row-key="id"
          hover
          stripe
          @select-change="handleSelectChange"
        >
          <template #sort="{ row }">
            <InputNumber
              :value="row.sort"
              :min="0"
              :max="1000"
              size="small"
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
                <Button
                  size="small"
                  theme="default"
                  variant="outline"
                  @click="handleOpenLeaderList(row)"
                >
                  <template #icon><UserIcon /></template>
                  领导列表
                </Button>
                <Button
                  size="small"
                  theme="primary"
                  variant="outline"
                  @click="handleAdd(row.id)"
                >
                  <template #icon><PlusIcon /></template>
                  新增
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
                  content="确认删除该部门吗？"
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
                  content="确认恢复该部门吗？"
                  @confirm="handleRecovery(row)"
                >
                  <Button size="small" theme="primary" variant="outline">
                    恢复
                  </Button>
                </Popconfirm>
                <Popconfirm
                  content="确认彻底删除该部门吗？"
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

    <DeptModal ref="deptModalRef" @success="handleSuccess" />
    <DeptLeaderModal ref="deptLeaderModalRef" />
  </Page>
</template>
