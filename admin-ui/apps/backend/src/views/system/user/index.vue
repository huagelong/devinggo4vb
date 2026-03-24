<script lang="ts" setup>
import { computed, onMounted, reactive, ref } from 'vue';

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
  RollbackIcon,
  SearchIcon,
  SettingIcon,
  UploadIcon,
} from 'tdesign-icons-vue-next';
import {
  Button,
  Checkbox,
  CheckboxGroup,
  DateRangePicker,
  DialogPlugin,
  Dropdown,
  Form,
  FormItem,
  Input,
  MessagePlugin,
  Popconfirm,
  Popup,
  Select,
  Space,
  Switch,
  Table,
  Tooltip,
  TreeSelect,
} from 'tdesign-vue-next';

import { getDeptTree } from '#/api/system/dept';
import { getPostList } from '#/api/system/post';
import { getRoleList } from '#/api/system/role';
import {
  changeUserStatus,
  clearUserCache,
  deleteUser,
  getRecycleUserList,
  getUserList,
  realDeleteUser,
  recoveryUser,
  resetPassword,
} from '#/api/system/user';

import DeptTree from './dept-tree.vue';
import UserModal from './user-modal.vue';

const currentDeptId = ref<number | string>('');
const isRecycleBin = ref(false);
const userModalRef = ref();
const isFullscreen = ref(false);
const tableContainerRef = ref<HTMLElement>();

function toggleFullscreen() {
  if (document.fullscreenElement) {
    document.exitFullscreen();
    isFullscreen.value = false;
  } else {
    tableContainerRef.value?.requestFullscreen();
    isFullscreen.value = true;
  }
}

document.addEventListener('fullscreenchange', () => {
  isFullscreen.value = !!document.fullscreenElement;
});

const searchForm = reactive({
  username: '',
  dept_id: undefined as number | undefined,
  role_id: undefined as number | undefined,
  phone: '',
  post_id: undefined as number | undefined,
  email: '',
  status: undefined as number | undefined,
  user_type: undefined as string | undefined,
  created_at: [] as string[],
});

const tableData = ref<any[]>([]);
const loading = ref(false);
const selectedRowKeys = ref<Array<number | string>>([]);
const selectOnRowClick = ref(false);

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0,
  showJumper: true,
  showPageSize: true,
  pageSizeOptions: [10, 20, 50, 100],
});

const roleOptions = ref<any[]>([]);
const postOptions = ref<any[]>([]);
const deptTreeData = ref<any[]>([]);

async function fetchOptions() {
  try {
    const [roleRes, postRes, deptRes] = await Promise.all([
      getRoleList().catch(() => null),
      getPostList().catch(() => null),
      getDeptTree().catch(() => null),
    ]);
    roleOptions.value = roleRes?.items || roleRes || [];
    postOptions.value = postRes?.items || postRes || [];
    deptTreeData.value = deptRes || [];
  } catch (error) {
    console.error(error);
  }
}

async function fetchTableData() {
  loading.value = true;
  try {
    const params: any = {
      page: pagination.current,
      limit: pagination.pageSize,
    };
    if (searchForm.username) params.username = searchForm.username;
    if (searchForm.role_id !== undefined) params.role_id = searchForm.role_id;
    if (searchForm.phone) params.phone = searchForm.phone;
    if (searchForm.post_id !== undefined) params.post_id = searchForm.post_id;
    if (searchForm.email) params.email = searchForm.email;
    if (searchForm.status !== undefined) params.status = searchForm.status;
    if (searchForm.user_type) params.user_type = searchForm.user_type;
    if (searchForm.created_at?.length === 2 && searchForm.created_at[0]) {
      params.created_at = searchForm.created_at;
    }
    if (currentDeptId.value) {
      params.dept_id = currentDeptId.value;
    } else if (searchForm.dept_id !== undefined) {
      params.dept_id = searchForm.dept_id;
    }

    const res = isRecycleBin.value
      ? await getRecycleUserList(params)
      : await getUserList(params);

    tableData.value = res?.items || [];
    pagination.total = res?.total || 0;
  } catch (error) {
    console.error(error);
  } finally {
    loading.value = false;
  }
}

function handleSearch() {
  pagination.current = 1;
  fetchTableData();
}

function handleReset() {
  searchForm.username = '';
  searchForm.dept_id = undefined;
  searchForm.role_id = undefined;
  searchForm.phone = '';
  searchForm.post_id = undefined;
  searchForm.email = '';
  searchForm.status = undefined;
  searchForm.user_type = undefined;
  searchForm.created_at = [];
  currentDeptId.value = '';
  pagination.current = 1;
  fetchTableData();
}

function handlePageChange(pageInfo: { current: number; pageSize: number }) {
  pagination.current = pageInfo.current;
  pagination.pageSize = pageInfo.pageSize;
  fetchTableData();
}

function handleSelectChange(value: Array<number | string>, ctx: any) {
  selectedRowKeys.value = value;
  console.log(value, ctx);
}

function handleDeptSelect(deptId: number | string) {
  currentDeptId.value = deptId;
  pagination.current = 1;
  fetchTableData();
}

function handleAdd() {
  userModalRef.value?.open();
}

function handleEdit(row: any) {
  userModalRef.value?.open(row);
}

async function handleDelete(row: any) {
  try {
    await (isRecycleBin.value ? realDeleteUser([row.id]) : deleteUser([row.id]));
    MessagePlugin.success('删除成功');
    fetchTableData();
  } catch (error) {
    console.error(error);
  }
}

async function handleBatchDelete() {
  if (selectedRowKeys.value.length === 0) {
    MessagePlugin.warning('请选择需要操作的数据');
    return;
  }
  try {
    await (isRecycleBin.value
      ? realDeleteUser(selectedRowKeys.value as number[])
      : deleteUser(selectedRowKeys.value as number[]));
    MessagePlugin.success('操作成功');
    selectedRowKeys.value = [];
    fetchTableData();
  } catch (error) {
    console.error(error);
  }
}

async function handleRecovery(row: any) {
  try {
    await recoveryUser([row.id]);
    MessagePlugin.success('恢复成功');
    fetchTableData();
  } catch (error) {
    console.error(error);
  }
}

async function handleBatchRecovery() {
  if (selectedRowKeys.value.length === 0) {
    MessagePlugin.warning('请选择需要操作的数据');
    return;
  }
  try {
    await recoveryUser(selectedRowKeys.value as number[]);
    MessagePlugin.success('恢复成功');
    selectedRowKeys.value = [];
    fetchTableData();
  } catch (error) {
    console.error(error);
  }
}

async function handleStatusChange(row: any, val: boolean) {
  const status = val ? 1 : 2;
  try {
    await changeUserStatus({ id: row.id, status });
    MessagePlugin.success('更新状态成功');
    fetchTableData();
  } catch (error) {
    console.error(error);
  }
}

async function handleResetPassword(row: any) {
  try {
    await resetPassword({ id: row.id });
    MessagePlugin.success('密码重置成功');
  } catch (error) {
    console.error(error);
  }
}

async function handleClearCache(row: any) {
  try {
    await clearUserCache({ id: row.id });
    MessagePlugin.success('清除缓存成功');
  } catch (error) {
    console.error(error);
  }
}

function handleSuccess() {
  fetchTableData();
}

function toggleRecycleBin() {
  isRecycleBin.value = !isRecycleBin.value;
  selectedRowKeys.value = [];
  pagination.current = 1;
  fetchTableData();
}

const actionDropdownOptions = [
  { content: '重置密码', value: 'reset_password' },
  { content: '更新缓存', value: 'clear_cache' },
];

function handleActionDropdownClick(data: any, row: any) {
  if (data.value === 'reset_password') {
    const dialog = DialogPlugin.confirm({
      header: '提示',
      body: '确认重置该用户密码吗？',
      onConfirm: () => {
        handleResetPassword(row);
        dialog.hide();
      },
      onClose: () => dialog.hide(),
    });
  } else if (data.value === 'clear_cache') {
    const dialog = DialogPlugin.confirm({
      header: '提示',
      body: '确认更新该用户缓存吗？',
      onConfirm: () => {
        handleClearCache(row);
        dialog.hide();
      },
      onClose: () => dialog.hide(),
    });
  }
}

const columns: any[] = [
  {
    colKey: 'row-select',
    type: 'multiple',
    width: 50,
    align: 'center',
  },
  { colKey: 'avatar', title: '头像', width: 80, align: 'center' },
  { colKey: 'username', title: '账户', minWidth: 100, align: 'center' },
  { colKey: 'dept_name', title: '所属部门', minWidth: 100, align: 'center' },
  { colKey: 'nickname', title: '昵称', minWidth: 100, align: 'center' },
  { colKey: 'role_name', title: '角色', minWidth: 100, align: 'center' },
  { colKey: 'phone', title: '手机', minWidth: 120, align: 'center' },
  { colKey: 'post_name', title: '岗位', minWidth: 100, align: 'center' },
  { colKey: 'email', title: '邮箱', minWidth: 150, align: 'center' },
  { colKey: 'status', title: '状态', width: 100, align: 'center' },
  { colKey: 'user_type', title: '用户类型', width: 100, align: 'center' },
  { colKey: 'created_at', title: '注册时间', width: 160, align: 'center' },
  { colKey: 'action', title: '操作', width: 220, fixed: 'right', align: 'center' },
];

const statusOptions = [
  { label: '正常', value: 1 },
  { label: '停用', value: 2 },
];

const userTypeOptions = [{ label: '系统用户', value: '100' }];

const _displayColumns = ref<string[]>(
  columns
    .filter((c) => c.colKey !== 'row-select')
    .map((c) => c.colKey as string),
);

// 始终在最前面保留 row-select 列
const displayColumns = computed({
  get: () => ['row-select', ..._displayColumns.value],
  set: (val: string[]) => {
    _displayColumns.value = val.filter((k) => k !== 'row-select');
  },
});

const columnOptions = columns
  .filter((c) => c.colKey !== 'row-select' && c.title)
  .map((c) => ({ label: c.title as string, value: c.colKey as string }));

const allColumnKeys = columnOptions.map((c) => c.value);
const isAllSelected = computed(() => _displayColumns.value.length === allColumnKeys.length);
const isIndeterminate = computed(
  () => _displayColumns.value.length > 0 && _displayColumns.value.length < allColumnKeys.length,
);

function toggleSelectAll(checked: boolean) {
  _displayColumns.value = checked ? [...allColumnKeys] : [];
}

onMounted(() => {
  fetchOptions();
  fetchTableData();
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
              <Button theme="default" @click="handleReset">重置</Button>
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
                <Button variant="outline">
                  <template #icon><UploadIcon /></template>
                  导入
                </Button>
                <Button variant="outline">
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
              <Tooltip content="刷新">
                <Button shape="square" variant="outline" @click="fetchTableData">
                  <template #icon><RefreshIcon /></template>
                </Button>
              </Tooltip>
              <Tooltip :content="isFullscreen ? '退出全屏' : '全屏'">
                <Button shape="square" variant="outline" @click="toggleFullscreen">
                  <template #icon>
                    <FullscreenExitIcon v-if="isFullscreen" />
                    <FullscreenIcon v-else />
                  </template>
                </Button>
              </Tooltip>
              <Tooltip :content="isRecycleBin ? '返回列表' : '显示回收站'">
                <Button shape="square" variant="outline" @click="toggleRecycleBin">
                  <template #icon>
                    <RollbackIcon v-if="isRecycleBin" />
                    <DeleteIcon v-else />
                  </template>
                </Button>
              </Tooltip>
              <Tooltip content="列配置">
                <Popup placement="bottom-right" trigger="click">
                  <Button shape="square" variant="outline">
                    <template #icon><SettingIcon /></template>
                  </Button>
                  <template #content>
                    <div class="min-w-[120px] p-3">
                      <Checkbox
                        :checked="isAllSelected"
                        :indeterminate="isIndeterminate"
                        @change="(val: any) => toggleSelectAll(val as boolean)"
                      >
                        全选
                      </Checkbox>
                      <t-divider style="margin: 8px 0" />
                      <CheckboxGroup
                        v-model="_displayColumns"
                        :options="columnOptions"
                        layout="vertical"
                      />
                    </div>
                  </template>
                </Popup>
              </Tooltip>
            </div>
          </div>

          <Table
            v-model:display-columns="displayColumns"
            :columns="columns"
            :data="tableData"
            :loading="loading"
            :pagination="pagination"
            :selected-row-keys="selectedRowKeys"
            :select-on-row-click="selectOnRowClick"
            row-key="id"
            hover
            stripe
            @page-change="handlePageChange"
            @select-change="handleSelectChange"
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
                :disabled="isRecycleBin"
                :value="row.status === 1"
                @change="(val: any) => handleStatusChange(row, val as boolean)"
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
                    content="确认删除该用户吗？"
                    @confirm="handleDelete(row)"
                  >
                    <Button size="small" theme="danger" variant="outline">
                      <template #icon><DeleteIcon /></template>
                      删除
                    </Button>
                  </Popconfirm>
                  <Dropdown
                    :options="actionDropdownOptions"
                    trigger="click"
                    @click="(item: any) => handleActionDropdownClick(item, row)"
                  >
                    <Button size="small" theme="default" variant="outline">
                      <template #icon><MoreIcon /></template>
                      更多
                    </Button>
                  </Dropdown>
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

    <UserModal ref="userModalRef" @success="handleSuccess" />
  </Page>
</template>


