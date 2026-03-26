<script lang="ts" setup>
import { computed, onMounted, onUnmounted, ref } from 'vue';

import { Page } from '@vben/common-ui';

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
  MessagePlugin,
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
import { useDictOptions } from '#/composables/crud/use-dict-options';

import DeptTree from './components/dept-tree.vue';
import UserModal from './components/user-modal.vue';
import { createUserColumnOptions, createUserTableColumns, userActionDropdownOptions } from './schemas';
import { useUserActions } from './use-user-actions';
import { useUserCrud } from './use-user-crud';

const currentDeptId = ref<number | string>('');
const userModalRef = ref();
const tableContainerRef = ref<HTMLElement>();
const isFullscreen = ref(false);

const roleOptions = ref<any[]>([]);
const postOptions = ref<any[]>([]);
const deptTreeData = ref<any[]>([]);
const statusOptions = ref<any[]>([]);
const userTypeOptions = ref<any[]>([]);
const homePageOptions = ref<any[]>([]);

const columns: any[] = createUserTableColumns();
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

    roleOptions.value = roleRes?.items || roleRes || [];
    postOptions.value = postRes?.items || postRes || [];
    deptTreeData.value = deptRes || [];
    statusOptions.value = statusDict || [];
    userTypeOptions.value = userTypeDict || [];
    homePageOptions.value = dashboardDict || [];
  } catch (error) {
    console.error(error);
  }
}

function handleAdd() {
  userModalRef.value?.open();
}

function handleEdit(row: any) {
  if (isSuperAdmin(row)) {
    MessagePlugin.warning('超级管理员不可编辑');
    return;
  }
  userModalRef.value?.open(row);
}

function handleSuccess() {
  fetchTableData();
}

onMounted(() => {
  fetchOptions();
  fetchTableData();
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
          <Form :data="searchForm" label-width="80px" colon>
            <div class="grid grid-cols-3 gap-x-4">
              <FormItem label="账户" name="username">
                <Input
                  v-model="searchForm.username"
                  placeholder="请输入账户"
                  clearable
                />
              </FormItem>
              <FormItem label="所属部门" name="dept_id">
                <TreeSelect
                  v-model="searchForm.dept_id"
                  :data="deptTreeData"
                  :keys="{ value: 'id', label: 'label', children: 'children' }"
                  placeholder="请选择所属部门"
                  clearable
                  class="w-full"
                />
              </FormItem>
              <FormItem label="角色" name="role_id">
                <Select
                  v-model="searchForm.role_id"
                  :options="roleOptions"
                  :keys="{ value: 'id', label: 'name' }"
                  placeholder="请选择角色"
                  clearable
                  class="w-full"
                />
              </FormItem>
              <FormItem label="手机" name="phone">
                <Input
                  v-model="searchForm.phone"
                  placeholder="请输入手机"
                  clearable
                />
              </FormItem>
              <FormItem label="岗位" name="post_id">
                <Select
                  v-model="searchForm.post_id"
                  :options="postOptions"
                  :keys="{ value: 'id', label: 'name' }"
                  placeholder="请选择岗位"
                  clearable
                  class="w-full"
                />
              </FormItem>
              <FormItem label="邮箱" name="email">
                <Input
                  v-model="searchForm.email"
                  placeholder="请输入邮箱"
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
              <FormItem label="用户类型" name="user_type">
                <Select
                  v-model="searchForm.user_type"
                  :options="userTypeOptions"
                  placeholder="请选择用户类型"
                  clearable
                  class="w-full"
                />
              </FormItem>
              <FormItem label="注册时间" name="created_at" class="col-span-1">
                <DateRangePicker
                  v-model="searchForm.created_at"
                  :placeholder="['开始时间', '结束时间']"
                  clearable
                  class="w-full"
                />
              </FormItem>
            </div>
            <div class="flex justify-end gap-2 pt-2">
              <Button theme="default" @click="handleResetWithDept">重置</Button>
              <Button theme="primary" @click="handleSearch">
                <template #icon><SearchIcon /></template>
                查询
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

            <div class="flex items-center gap-2">
              <Tooltip :content="isFullscreen ? '退出全屏' : '全屏'">
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
            @select-change="(keys: any) => handleSelectChange(keys as Array<number | string>)"
          >
            <template #avatar="{ row }">
              <img
                v-if="row.avatar"
                :src="row.avatar"
                class="mx-auto h-8 w-8 rounded-full object-cover"
                alt="avatar"
              />
              <span v-else class="text-gray-400">-</span>
            </template>

            <template #status="{ row }">
              <Switch
                :disabled="isRecycleBin || isSuperAdmin(row)"
                :value="row.status === 1"
                @change="(value: any) => handleStatusChange(row, Boolean(value))"
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
                      更新缓存
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
                      编辑
                    </Button>
                    <Popconfirm
                      content="确认删除该用户吗？"
                      @confirm="handleDelete(row)"
                    >
                      <Button size="small" theme="danger" variant="outline">
                        <template #icon><DeleteIcon /></template>
                        删除
                      </Button>
                    </Popconfirm>
                    <Dropdown
                      :options="userActionDropdownOptions"
                      trigger="click"
                      @click="(item: any) => handleActionDropdownClick(item, row)"
                    >
                      <Button size="small" theme="default" variant="outline">
                        <template #icon><MoreIcon /></template>
                        更多
                      </Button>
                    </Dropdown>
                  </template>
                </template>

                <template v-else>
                  <Popconfirm
                    content="确认恢复该用户吗？"
                    @confirm="handleRecovery(row)"
                  >
                    <Button size="small" theme="primary" variant="outline">
                      恢复
                    </Button>
                  </Popconfirm>
                  <Popconfirm
                    content="确认彻底删除该用户吗？"
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
    </div>

    <input
      ref="importInputRef"
      type="file"
      accept=".xlsx,.xls,.csv"
      class="hidden"
      @change="handleImportChange"
    />

    <Dialog
      v-model:visible="setHomePageVisible"
      width="520px"
      header="设置用户后台首页"
      destroy-on-close
    >
      <Form label-width="90px">
        <FormItem label="用户首页">
          <Select
            v-model="selectedHomePage"
            :options="homePageOptions"
            placeholder="请选择用户首页"
            clearable
            class="w-full"
          />
        </FormItem>
      </Form>
      <template #footer>
        <div class="flex justify-end gap-2">
          <Button theme="default" @click="setHomePageVisible = false">取消</Button>
          <Button theme="primary" :loading="setHomePageLoading" @click="handleSetHomePage">
            保存
          </Button>
        </div>
      </template>
    </Dialog>

    <UserModal ref="userModalRef" @success="handleSuccess" />
  </Page>
</template>
