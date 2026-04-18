<script lang="ts" setup>
import type { RoleApi } from '#/api/system/role';
import type { DictOption } from '#/composables/crud/use-dict-options';

import { computed, onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { $t } from '@vben/locales';

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
            { label: $t('common.statusEnabled'), value: 1 },
            { label: $t('common.statusDisabled'), value: 2 },
          ];
}

function handleAdd() {
  roleModalRef.value?.open();
}

function handleEdit(row: RoleListItem) {
  if (row.id === 1 || row.code === 'superAdmin') {
    message.error($t('common.superAdminRoleCannotEdit'));
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
    message.error($t('common.superAdminRoleCannotDelete'));
    return;
  }
  try {
    await (isRecycleBin.value ? realDeleteRole([row.id]) : deleteRole([row.id]));
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
  if (ids.includes(1)) {
    message.error($t('common.superAdminRoleCannotDelete'));
    return;
  }
  try {
    await (isRecycleBin.value ? realDeleteRole(ids) : deleteRole(ids));
    message.success($t('common.operationSuccess'));
    clearSelectedRowKeys();
    await fetchTableData();
  } catch (error) {
    logger.error(error);
    message.error($t('common.batchDeleteFailed'));
  }
}

async function handleRecovery(row: RoleListItem) {
  try {
    await recoveryRole([row.id]);
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
  try {
    await recoveryRole(ids);
    message.success($t('common.recoverySuccess'));
    clearSelectedRowKeys();
    await fetchTableData();
  } catch (error) {
    logger.error(error);
    message.error($t('common.batchRecoveryFailed'));
  }
}

async function handleStatusChange(row: RoleListItem, checked: boolean) {
  if (row.code === 'superAdmin') {
    message.info($t('common.superAdminCannotDisable'));
    return;
  }
  const status = checked ? 1 : 2;
  try {
    await changeRoleStatus({ id: row.id, status });
    message.success($t('common.statusUpdateSuccess'));
    await fetchTableData();
  } catch (error) {
    logger.error(error);
    message.error($t('common.statusUpdateFailed'));
  }
}

async function handleSortChange(value: number | string, row: RoleListItem) {
  const numberValue = Number(value);
  if (Number.isNaN(numberValue)) return;
  if (row.id === 1) {
    message.info($t('common.superAdminCannotModify'));
    return;
  }

  try {
    await updateRoleNumber({
      id: Number(row.id),
      numberName: 'sort',
      numberValue,
    });
    message.success($t('common.sortUpdateSuccess'));
    await fetchTableData();
  } catch (error) {
    logger.error(error);
    message.error($t('common.sortUpdateFailed'));
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
            <FormItem :label="$t('system.role.name')" name="name">
              <Input
                v-model="searchForm.name"
                :placeholder="$t('ui.placeholder.input', [$t('system.role.name')])"
                clearable
              />
            </FormItem>
            <FormItem :label="$t('system.role.code')" name="code">
              <Input
                v-model="searchForm.code"
                :placeholder="$t('ui.placeholder.input', [$t('system.role.code')])"
                clearable
              />
            </FormItem>
            <FormItem :label="$t('common.status')" name="status">
              <Select
                v-model="searchForm.status"
                :options="statusOptions"
                :placeholder="$t('ui.placeholder.select', [$t('common.status')])"
                clearable
                class="w-full"
              />
            </FormItem>
            <FormItem :label="$t('common.createTime')" name="created_at">
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
                    {{ $t('common.edit') }}
                  </Button>
                  <Button
                    size="small"
                    theme="default"
                    variant="outline"
                    @click="handleMenuPermission(row)"
                  >
                    {{ $t('system.role.menuPermission') }}
                  </Button>
                  <Button
                    size="small"
                    theme="default"
                    variant="outline"
                    @click="handleDataPermission(row)"
                  >
                    {{ $t('system.role.dataPermission') }}
                  </Button>
                  <Popconfirm
                    :content="$t('system.role.confirmDelete')"
                    @confirm="handleDelete(row)"
                  >
                    <Button size="small" theme="danger" variant="outline">
                      <template #icon><DeleteIcon /></template>
                      {{ $t('common.delete') }}
                    </Button>
                  </Popconfirm>
                </template>
              </template>

              <template v-else>
                <Popconfirm
                  :content="$t('common.confirmRecoveryRole')"
                  @confirm="handleRecovery(row)"
                >
                  <Button size="small" theme="primary" variant="outline">
                    {{ $t('common.recovery') }}
                  </Button>
                </Popconfirm>
                <Popconfirm
                  :content="$t('common.confirmPermanentDeleteRole')"
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
