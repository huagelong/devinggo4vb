<script lang="ts" setup>
import type { RoleApi } from '#/api/system/role';
import type { DictOption } from '#/composables/crud/use-dict-options';

import { computed, onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';

import { message } from '#/adapter/tdesign';
import { logger } from '#/utils/logger';
import {
  changeRoleStatus,
  deleteRole,
  realDeleteRole,
  recoveryRole,
  updateRoleNumber,
} from '#/api/system/role';
import CrudToolbar from '#/components/crud/crud-toolbar.vue';
import { useDictOptions } from '#/composables/crud/use-dict-options';

import { AddIcon, DeleteIcon, EditIcon, SearchIcon } from 'tdesign-icons-vue-next';
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

import RoleDataPermissionModal from './components/role-data-permission-modal.vue';
import RoleMenuPermissionModal from './components/role-menu-permission-modal.vue';
import RoleModal from './components/role-modal.vue';
import type { RoleListItem, RoleTableColumn } from './model';
import {
  createRoleColumnOptions,
  createRoleTableColumns,
  getRoleDataScopeLabel,
} from './schemas';
import { useRoleCrud } from './use-role-crud';

defineOptions({ name: 'SystemRole' });

type RoleModalInstance = {
  open: (data?: Partial<RoleApi.SubmitPayload>) => void;
};

type RolePermissionModalInstance = {
  open: (data: RoleApi.ListItem) => void;
};

const roleModalRef = ref<RoleModalInstance>();
const roleMenuPermissionModalRef = ref<RolePermissionModalInstance>();
const roleDataPermissionModalRef = ref<RolePermissionModalInstance>();
const statusOptions = ref<DictOption[]>([]);

const columns: RoleTableColumn[] = createRoleTableColumns();
const columnOptions = createRoleColumnOptions(columns);
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
} = useRoleCrud();

const { getDictOptions } = useDictOptions();

function toIds(keys: Array<number | string>) {
  return keys.map((key) => Number(key));
}

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

function handleAdd() {
  roleModalRef.value?.open();
}

function handleEdit(row: RoleListItem) {
  if (row.id === 1 || row.code === 'superAdmin') {
    message.error('超级管理员角色不可编辑');
    return;
  }
  roleModalRef.value?.open({
    ...row,
    data_scope: Number(row.data_scope ?? 1),
    status: Number(row.status ?? 1),
  });
}

function handleMenuPermission(row: RoleListItem) {
  roleMenuPermissionModalRef.value?.open(row);
}

function handleDataPermission(row: RoleListItem) {
  roleDataPermissionModalRef.value?.open(row);
}

async function handleDelete(row: RoleListItem) {
  if (row.id === 1 || row.code === 'superAdmin') {
    message.error('超级管理员角色不可删除');
    return;
  }
  try {
    await (isRecycleBin.value ? realDeleteRole([row.id]) : deleteRole([row.id]));
    message.success('操作成功');
    await fetchTableData();
  } catch (error) {
    logger.error(error);
    message.error('删除失败，请稍后重试');
  }
}

async function handleBatchDelete() {
  if (selectedRowKeys.value.length === 0) {
    message.warning('请选择需要操作的数据');
    return;
  }

  const ids = toIds(selectedRowKeys.value);
  if (ids.includes(1)) {
    message.error('超级管理员角色不可删除');
    return;
  }
  try {
    await (isRecycleBin.value ? realDeleteRole(ids) : deleteRole(ids));
    message.success('操作成功');
    clearSelectedRowKeys();
    await fetchTableData();
  } catch (error) {
    logger.error(error);
    message.error('批量删除失败，请稍后重试');
  }
}

async function handleRecovery(row: RoleListItem) {
  try {
    await recoveryRole([row.id]);
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
  try {
    await recoveryRole(ids);
    message.success('恢复成功');
    clearSelectedRowKeys();
    await fetchTableData();
  } catch (error) {
    logger.error(error);
    message.error('批量恢复失败，请稍后重试');
  }
}

async function handleStatusChange(row: RoleListItem, checked: boolean) {
  if (row.code === 'superAdmin') {
    message.info('超级管理员角色不能禁用');
    return;
  }
  const status = checked ? 1 : 2;
  try {
    await changeRoleStatus({ id: row.id, status });
    message.success('状态更新成功');
    await fetchTableData();
  } catch (error) {
    logger.error(error);
    message.error('状态更新失败，请稍后重试');
  }
}

async function handleSortChange(value: number | string, row: RoleListItem) {
  const numberValue = Number(value);
  if (Number.isNaN(numberValue)) return;
  if (row.id === 1) {
    message.info('超级管理员不可修改');
    return;
  }

  try {
    await updateRoleNumber({
      id: Number(row.id),
      numberName: 'sort',
      numberValue,
    });
    message.success('排序更新成功');
    await fetchTableData();
  } catch (error) {
    logger.error(error);
    message.error('排序更新失败，请稍后重试');
  }
}

function handleSuccess() {
  void fetchTableData();
}

function handleTableSelectChange(keys: Array<number | string>) {
  handleSelectChange(keys);
}

function handleStatusSwitchChange(row: RoleListItem, value: unknown) {
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
            <FormItem label="角色名称" name="name">
              <Input
                v-model="searchForm.name"
                placeholder="请输入角色名称"
                clearable
              />
            </FormItem>
            <FormItem label="角色标识" name="code">
              <Input
                v-model="searchForm.code"
                placeholder="请输入角色标识"
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
          <template #data_scope="{ row }">
            <span>{{ getRoleDataScopeLabel(row.data_scope) }}</span>
          </template>

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
              :disabled="isRecycleBin || row.code === 'superAdmin'"
              :value="Number(row.status) === 1"
              @change="(value: unknown) => handleStatusSwitchChange(row, value)"
            />
          </template>

          <template #action="{ row }">
            <div class="flex items-center justify-center gap-1">
              <template v-if="!isRecycleBin">
                <template v-if="row.code !== 'superAdmin'">
                  <Button
                    size="small"
                    theme="primary"
                    variant="outline"
                    @click="handleEdit(row)"
                  >
                    <template #icon><EditIcon /></template>
                    编辑
                  </Button>
                  <Button
                    size="small"
                    theme="default"
                    variant="outline"
                    @click="handleMenuPermission(row)"
                  >
                    菜单权限
                  </Button>
                  <Button
                    size="small"
                    theme="default"
                    variant="outline"
                    @click="handleDataPermission(row)"
                  >
                    数据权限
                  </Button>
                  <Popconfirm
                    content="确认删除该角色吗？"
                    @confirm="handleDelete(row)"
                  >
                    <Button size="small" theme="danger" variant="outline">
                      <template #icon><DeleteIcon /></template>
                      删除
                    </Button>
                  </Popconfirm>
                </template>
              </template>

              <template v-else>
                <Popconfirm
                  content="确认恢复该角色吗？"
                  @confirm="handleRecovery(row)"
                >
                  <Button size="small" theme="primary" variant="outline">
                    恢复
                  </Button>
                </Popconfirm>
                <Popconfirm
                  content="确认彻底删除该角色吗？"
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

    <RoleModal ref="roleModalRef" @success="handleSuccess" />
    <RoleMenuPermissionModal
      ref="roleMenuPermissionModalRef"
      @success="handleSuccess"
    />
    <RoleDataPermissionModal
      ref="roleDataPermissionModalRef"
      @success="handleSuccess"
    />
  </Page>
</template>
