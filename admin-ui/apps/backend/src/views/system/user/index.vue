<script lang="ts" setup>
import { computed, onMounted, onUnmounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { $t } from '@vben/locales';
import { message } from '#/adapter/tdesign';
import { logger } from '#/utils/logger';

import {
  AddIcon,
  DeleteIcon,
  DownloadIcon,
  EditIcon,
  FullscreenExitIcon,
  FullscreenIcon,
  MoreIcon,
  RefreshIcon,
  SearchIcon,
  UploadIcon,
} from 'tdesign-icons-vue-next';
import {
  Button,
  DateRangePicker,
  Dialog,
  Dropdown,
  Form,
  FormItem,
  Input,
  Popconfirm,
  Select,
  Space,
  Switch,
  Table,
  Tooltip,
  TreeSelect,
} from 'tdesign-vue-next';

import CrudToolbar from '#/components/crud/crud-toolbar.vue';
import { getDeptTree } from '#/api/system/dept';
import { getPostList } from '#/api/system/post';
import { getRoleList } from '#/api/system/role';
import type { DeptApi } from '#/api/system/dept';
import type { PostApi } from '#/api/system/post';
import type { RoleApi } from '#/api/system/role';
import type { DictOption } from '#/composables/crud/use-dict-options';
import { useDictOptions } from '#/composables/crud/use-dict-options';

import DeptTree from './components/dept-tree.vue';
import UserModal from './components/user-modal.vue';
import type {
  UserActionDropdownItem,
  UserListItem,
  UserTableColumn,
} from './model';
import { createUserColumnOptions, createUserTableColumns, userActionDropdownOptions } from './schemas';
import { useUserActions } from './use-user-actions';
import { useUserCrud } from './use-user-crud';

const currentDeptId = ref<number | string>('');
type UserModalInstance = {
  open: (data?: Partial<UserListItem>) => void;
};

const userModalRef = ref<UserModalInstance>();
const tableContainerRef = ref<HTMLElement>();
const isFullscreen = ref(false);
const importDialogVisible = ref(false);

const roleOptions = ref<RoleApi.ListItem[]>([]);
const postOptions = ref<PostApi.ListItem[]>([]);
const deptTreeData = ref<DeptApi.TreeNode[]>([]);
const statusOptions = ref<DictOption[]>([]);
const userTypeOptions = ref<DictOption[]>([]);
const homePageOptions = ref<DictOption[]>([]);

const columns: UserTableColumn[] = createUserTableColumns();
const columnOptions = createUserColumnOptions(columns);
const allColumnKeys = columnOptions.map((item) => item.value);
const visibleColumns = ref<string[]>([...allColumnKeys]);

const displayColumns = computed({
  get: () => ['row-select', ...visibleColumns.value],
  set: (value: string[]) => {
    visibleColumns.value = value.filter((item) => item !== 'row-select');
  },
});

const {
  buildRequestParams,
  clearSelectedRowKeys,
  fetchTableData,
  handleDeptSelect,
  handlePageChange,
  handleResetWithDept,
  handleSearch,
  handleSelectChange,
  isRecycleBin,
  loading,
  pagination,
  searchForm,
  selectedRowKeys,
  tableData,
  toggleRecycleBin,
} = useUserCrud(currentDeptId);

const {
  exportLoading,
  handleActionDropdownClick,
  handleBatchDelete,
  handleBatchRecovery,
  handleClearCache,
  handleDelete,
  handleDownloadTemplate,
  handleExport,
  handleImportChange,
  handleRecovery,
  handleSetHomePage,
  handleStatusChange,
  importInputRef,
  importLoading,
  isSuperAdmin,
  selectedHomePage,
  setHomePageLoading,
  setHomePageVisible,
  templateLoading,
  triggerImport,
} = useUserActions({
  buildRequestParams,
  clearSelectedRowKeys,
  fetchTableData,
  isRecycleBin,
  selectedRowKeys,
});

void importInputRef;

const { getDictOptions } = useDictOptions();

function handleFullscreenChange() {
  isFullscreen.value = !!document.fullscreenElement;
}

function toggleFullscreen() {
  if (document.fullscreenElement) {
    document.exitFullscreen();
    return;
  }
  tableContainerRef.value?.requestFullscreen();
}

function normalizeListData<T>(data: T[] | { items?: T[] } | null | undefined): T[] {
  if (Array.isArray(data)) {
    return data;
  }
  return Array.isArray(data?.items) ? data.items : [];
}

async function fetchOptions() {
  try {
    const [roleRes, postRes, deptRes, statusDict, userTypeDict, dashboardDict] = await Promise.all([
      getRoleList().catch(() => null),
      getPostList().catch(() => null),
      getDeptTree().catch(() => null),
      getDictOptions('data_status'),
      getDictOptions('user_type'),
      getDictOptions('dashboard'),
    ]);

    roleOptions.value = normalizeListData(roleRes);
    postOptions.value = normalizeListData(postRes);
    deptTreeData.value = deptRes || [];
    statusOptions.value = statusDict || [];
    userTypeOptions.value = userTypeDict || [];
    homePageOptions.value = dashboardDict || [];
  } catch (error) {
    logger.error(error);
    message.error($t('common.filterLoadFailed'));
  }
}

function handleAdd() {
  userModalRef.value?.open();
}

function openImportDialog() {
  importDialogVisible.value = true;
}

async function handleImportChangeWithClose(event: Event) {
  const success = await handleImportChange(event);
  if (success) {
    importDialogVisible.value = false;
  }
}

function handleEdit(row: UserListItem) {
  if (isSuperAdmin(row)) {
    message.warning($t('common.superAdminCannotEdit'));
    return;
  }
  userModalRef.value?.open(row);
}

function handleSuccess() {
  void fetchTableData();
}

function handleTableSelectChange(keys: Array<number | string>) {
  handleSelectChange(keys);
}

function handleStatusSwitchChange(row: UserListItem, value: unknown) {
  void handleStatusChange(row, Boolean(value));
}

function handleActionDropdownItemClick(
  item: unknown,
  row: UserListItem,
) {
  handleActionDropdownClick(item as UserActionDropdownItem, row);
}

onMounted(() => {
  void fetchOptions();
  void fetchTableData();
  document.addEventListener('fullscreenchange', handleFullscreenChange);
});

onUnmounted(() => {
  document.removeEventListener('fullscreenchange', handleFullscreenChange);
});
</script>

<template>
  <Page auto-content-height>
    <div class="flex h-full flex-row gap-4">
      <div class="h-full rounded-md bg-background p-2">
        <DeptTree @select="handleDeptSelect" />
      </div>

      <div class="flex h-full min-w-0 flex-1 flex-col gap-3 overflow-hidden">
        <div class="rounded-md bg-white p-4">
          <Form :data="searchForm" label-width="80px" layout="inline" colon>
            <div class="grid grid-cols-3 gap-x-4 gap-y-3">
              <FormItem :label="$t('system.user.username')" name="username">
                <Input
                  v-model="searchForm.username"
                  :placeholder="$t('ui.placeholder.input', [$t('system.user.username')])"
                  clearable
                />
              </FormItem>
              <FormItem :label="$t('system.user.dept')" name="dept_ids">
                <TreeSelect
                  v-model="searchForm.dept_ids"
                  :data="deptTreeData"
                  :keys="{ value: 'id', label: 'label', children: 'children' }"
                  :multiple="true"
                  :tree-props="{ checkStrictly: true }"
                  :placeholder="$t('ui.placeholder.select', [$t('system.user.dept')])"
                  clearable
                  class="w-full"
                />
              </FormItem>
              <FormItem :label="$t('system.user.role')" name="role_id">
                <Select
                  v-model="searchForm.role_id"
                  :options="roleOptions"
                  :keys="{ value: 'id', label: 'name' }"
                  :placeholder="$t('ui.placeholder.select', [$t('system.user.role')])"
                  clearable
                  class="w-full"
                />
              </FormItem>
              <FormItem :label="$t('system.user.phone')" name="phone">
                <Input
                  v-model="searchForm.phone"
                  :placeholder="$t('ui.placeholder.input', [$t('system.user.phone')])"
                  clearable
                />
              </FormItem>
              <FormItem :label="$t('system.user.post')" name="post_id">
                <Select
                  v-model="searchForm.post_id"
                  :options="postOptions"
                  :keys="{ value: 'id', label: 'name' }"
                  :placeholder="$t('ui.placeholder.select', [$t('system.user.post')])"
                  clearable
                  class="w-full"
                />
              </FormItem>
              <FormItem :label="$t('system.user.email')" name="email">
                <Input
                  v-model="searchForm.email"
                  :placeholder="$t('ui.placeholder.input', [$t('system.user.email')])"
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
              <FormItem :label="$t('system.user.userType')" name="user_type">
                <Select
                  v-model="searchForm.user_type"
                  :options="userTypeOptions"
                  :placeholder="$t('ui.placeholder.select', [$t('system.user.userType')])"
                  clearable
                  class="w-full"
                />
              </FormItem>
              <FormItem :label="$t('system.user.registerTime')" name="created_at" class="col-span-1">
                <DateRangePicker
                  v-model="searchForm.created_at"
                  :placeholder="[$t('common.startTime'), $t('common.endTime')]"
                  clearable
                  class="w-full"
                />
              </FormItem>
            </div>
            <div class="flex justify-end gap-2 pt-2">
              <Button theme="default" @click="handleResetWithDept">{{ $t('common.reset') }}</Button>
              <Button theme="primary" @click="handleSearch">
                <template #icon><SearchIcon /></template>
                {{ $t('common.query') }}
              </Button>
            </div>
          </Form>
        </div>

        <div ref="tableContainerRef" class="flex min-h-0 flex-1 flex-col rounded-md bg-white p-4">
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
                <Button variant="outline" :loading="importLoading" @click="openImportDialog">
                  <template #icon><UploadIcon /></template>
                  {{ $t('common.import') }}
                </Button>
                <Button variant="outline" :loading="exportLoading" @click="handleExport">
                  <template #icon><DownloadIcon /></template>
                  {{ $t('common.export') }}
                </Button>
              </template>
              <template v-else>
                <Button theme="success" @click="handleBatchRecovery">{{ $t('common.recovery') }}</Button>
                <Button theme="danger" @click="handleBatchDelete">{{ $t('common.permanentDelete') }}</Button>
              </template>
            </Space>

            <div class="flex items-center gap-2">
              <Tooltip :content="isFullscreen ? $t('common.exitFullscreen') : $t('common.fullscreen')">
                <Button shape="square" variant="outline" @click="toggleFullscreen">
                  <template #icon>
                    <FullscreenExitIcon v-if="isFullscreen" />
                    <FullscreenIcon v-else />
                  </template>
                </Button>
              </Tooltip>

              <CrudToolbar
                v-model="visibleColumns"
                :column-options="columnOptions"
                :is-recycle-bin="isRecycleBin"
                @refresh="fetchTableData"
                @toggle-recycle="toggleRecycleBin"
              />
            </div>
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
            <template #avatar="{ row }">
              <img
                v-if="row?.avatar"
                :src="row?.avatar"
                class="mx-auto h-8 w-8 rounded-full object-cover"
                alt="avatar"
              />
              <span v-else class="text-gray-400">-</span>
            </template>

            <template #status="{ row }">
              <Switch
                :disabled="isRecycleBin || isSuperAdmin(row)"
                :value="row?.status === 1"
                @change="(value: unknown) => handleStatusSwitchChange(row, value)"
              />
            </template>

            <template #action="{ row }">
              <div class="flex items-center justify-center gap-1">
                <template v-if="!isRecycleBin">
                  <template v-if="isSuperAdmin(row)">
                    <Button
                      size="small"
                      theme="default"
                      variant="outline"
                      @click="handleClearCache(row)"
                    >
                      <template #icon><RefreshIcon /></template>
                      {{ $t('system.user.clearCache') }}
                    </Button>
                  </template>
                  <template v-else>
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
                      :content="$t('system.user.confirmDelete')"
                      @confirm="handleDelete(row)"
                    >
                      <Button size="small" theme="danger" variant="outline">
                        <template #icon><DeleteIcon /></template>
                        {{ $t('common.delete') }}
                      </Button>
                    </Popconfirm>
                    <Dropdown
                      :options="userActionDropdownOptions"
                      trigger="click"
                      @click="(item: unknown) => handleActionDropdownItemClick(item, row)"
                    >
                      <Button size="small" theme="default" variant="outline">
                        <template #icon><MoreIcon /></template>
                        {{ $t('common.more') }}
                      </Button>
                    </Dropdown>
                  </template>
                </template>

                <template v-else>
                  <Popconfirm
                    :content="$t('system.user.confirmRecovery')"
                    @confirm="handleRecovery(row)"
                  >
                    <Button size="small" theme="primary" variant="outline">
                      {{ $t('common.recovery') }}
                    </Button>
                  </Popconfirm>
                  <Popconfirm
                    :content="$t('system.user.confirmPermanentDelete')"
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
    </div>

    <input
      ref="importInputRef"
      type="file"
      accept=".xlsx,.xls,.csv"
      class="hidden"
      @change="handleImportChangeWithClose"
    />

    <Dialog
      v-model:visible="importDialogVisible"
      width="420px"
      :header="$t('common.import')"
      destroy-on-close
      :close-on-overlay-click="true"
    >
      <div class="flex flex-col gap-4">
        <p class="text-sm text-text-2">
          {{ $t('common.importDialogDescription') }}
        </p>
        <div class="flex flex-col gap-3">
          <Button
            variant="outline"
            :loading="templateLoading"
            @click="handleDownloadTemplate"
          >
            <template #icon><DownloadIcon /></template>
            {{ $t('common.importTemplate') }}
          </Button>
          <Button
            theme="primary"
            :loading="importLoading"
            @click="triggerImport"
          >
            <template #icon><UploadIcon /></template>
            {{ $t('common.import') }}
          </Button>
        </div>
      </div>
      <template #footer>
        <div class="flex justify-end">
          <Button theme="default" @click="importDialogVisible = false">
            {{ $t('common.cancel') }}
          </Button>
        </div>
      </template>
    </Dialog>

    <Dialog
      v-model:visible="setHomePageVisible"
      width="520px"
      :header="$t('system.user.setHomePage')"
      destroy-on-close
    >
      <Form label-width="130px" layout="inline">
        <FormItem :label="$t('system.user.selectHomePage')">
          <Select
            v-model="selectedHomePage"
            :options="homePageOptions"
            :placeholder="$t('ui.placeholder.select', [$t('system.user.selectHomePage')])"
            clearable
            class="w-full"
          />
        </FormItem>
      </Form>
      <template #footer>
        <div class="flex justify-end gap-2">
          <Button theme="default" @click="setHomePageVisible = false">{{ $t('common.cancel') }}</Button>
          <Button theme="primary" :loading="setHomePageLoading" @click="handleSetHomePage">
            {{ $t('common.save') }}
          </Button>
        </div>
      </template>
    </Dialog>

    <UserModal ref="userModalRef" @success="handleSuccess" />
  </Page>
</template>
