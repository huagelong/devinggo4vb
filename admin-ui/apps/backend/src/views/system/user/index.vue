<script lang="ts" setup>
import { ref } from 'vue';

import { Page } from '@vben/common-ui';

import {
  AddIcon,
  DeleteIcon,
  DownloadIcon,
  EditIcon,
  MoreIcon,
  UploadIcon,
} from 'tdesign-icons-vue-next';
import {
  Button,
  DialogPlugin,
  Dropdown,
  MessagePlugin,
  Popconfirm,
  Switch,
} from 'tdesign-vue-next';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
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

const [Grid, gridApi] = useVbenVxeGrid({
  formOptions: {
    // 强制三列布局，跟设计图保持一致
    wrapperClass: 'grid-cols-3',
    handleReset: () => {
      currentDeptId.value = '';
      gridApi.reload();
    },
    schema: [
      {
        fieldName: 'username',
        label: '账户',
        component: 'Input',
        componentProps: { placeholder: '请输入账户' },
      },
      {
        fieldName: 'dept_id',
        label: '所属部门',
        component: 'ApiTreeSelect',
        componentProps: {
          api: getDeptTree,
          labelField: 'name',
          valueField: 'id',
          childrenField: 'children',
          placeholder: '请选择所属部门',
        },
      },
      {
        fieldName: 'role_id',
        label: '角色',
        component: 'ApiSelect',
        componentProps: {
          api: async () => {
            const res = await getRoleList();
            return res?.items || res || [];
          },
          labelField: 'name',
          valueField: 'id',
          placeholder: '请选择角色',
        },
      },
      {
        fieldName: 'phone',
        label: '手机',
        component: 'Input',
        componentProps: { placeholder: '请输入手机' },
      },
      {
        fieldName: 'post_id',
        label: '岗位',
        component: 'ApiSelect',
        componentProps: {
          api: async () => {
            const res = await getPostList();
            return res?.items || res || [];
          },
          labelField: 'name',
          valueField: 'id',
          placeholder: '请选择岗位',
        },
      },
      {
        fieldName: 'email',
        label: '邮箱',
        component: 'Input',
        componentProps: { placeholder: '请输入邮箱' },
      },
      {
        fieldName: 'status',
        label: '状态',
        component: 'Select',
        componentProps: {
          placeholder: '请选择状态',
          options: [
            { label: '正常', value: 1 },
            { label: '停用', value: 2 },
          ],
        },
      },
      {
        fieldName: 'user_type',
        label: '用户类型',
        component: 'Select',
        componentProps: {
          placeholder: '请选择用户类型',
          options: [{ label: '系统用户', value: '100' }],
        },
      },
      {
        fieldName: 'created_at',
        label: '注册时间',
        component: 'RangePicker',
        componentProps: {
          placeholder: ['请选择开始时间', '请选择结束时间'],
        },
      },
    ],
  },
  gridOptions: {
    toolbarConfig: {
      custom: true,
      refresh: true,
      zoom: true,
    },
    proxyConfig: {
      ajax: {
        query: async ({ page }: any, formValues: any) => {
          const params = {
            page: page.currentPage,
            limit: page.pageSize,
            ...formValues,
          };
          if (currentDeptId.value) {
            params.dept_id = currentDeptId.value;
          }
          if (isRecycleBin.value) {
            return await getRecycleUserList(params);
          }
          return await getUserList(params);
        },
      },
    },
    columns: [
      { type: 'checkbox', width: 60, align: 'center' },
      {
        field: 'avatar',
        title: '头像',
        width: 80,
        align: 'center',
        cellRender: { name: 'CellImage' },
      },
      { field: 'username', title: '账户', minWidth: 100, align: 'center' },
      { field: 'dept_name', title: '所属部门', minWidth: 100, align: 'center' },
      { field: 'nickname', title: '昵称', minWidth: 100, align: 'center' },
      { field: 'role_name', title: '角色', minWidth: 100, align: 'center' },
      { field: 'phone', title: '手机', minWidth: 120, align: 'center' },
      { field: 'post_name', title: '岗位', minWidth: 100, align: 'center' },
      { field: 'email', title: '邮箱', minWidth: 150, align: 'center' },
      {
        field: 'status',
        title: '状态',
        width: 100,
        align: 'center',
        slots: { default: 'status' },
      },
      { field: 'user_type', title: '用户类型', width: 100, align: 'center' },
      { field: 'created_at', title: '注册时间', width: 160, align: 'center' },
      {
        field: 'action',
        title: '操作',
        width: 200,
        fixed: 'right',
        align: 'center',
        slots: { default: 'action' },
      },
    ],
  },
});

function handleDeptSelect(deptId: number | string) {
  currentDeptId.value = deptId;
  gridApi.reload();
}

function handleAdd() {
  userModalRef.value?.open();
}

function handleEdit(row: any) {
  userModalRef.value?.open(row);
}

async function handleDelete(row: any) {
  try {
    await (isRecycleBin.value
      ? realDeleteUser([row.id])
      : deleteUser([row.id]));
    MessagePlugin.success('删除成功');
    gridApi.reload();
  } catch (error) {
    console.error(error);
  }
}

async function handleBatchDelete() {
  const records = gridApi.grid.getCheckboxRecords();
  if (records.length === 0) {
    MessagePlugin.warning('请选择需要操作的数据');
    return;
  }
  const ids = records.map((item: any) => item.id);
  try {
    await (isRecycleBin.value ? realDeleteUser(ids) : deleteUser(ids));
    MessagePlugin.success('操作成功');
    gridApi.reload();
  } catch (error) {
    console.error(error);
  }
}

async function handleRecovery(row: any) {
  try {
    await recoveryUser([row.id]);
    MessagePlugin.success('恢复成功');
    gridApi.reload();
  } catch (error) {
    console.error(error);
  }
}

async function handleBatchRecovery() {
  const records = gridApi.grid.getCheckboxRecords();
  if (records.length === 0) {
    MessagePlugin.warning('请选择需要操作的数据');
    return;
  }
  const ids = records.map((item: any) => item.id);
  try {
    await recoveryUser(ids);
    MessagePlugin.success('恢复成功');
    gridApi.reload();
  } catch (error) {
    console.error(error);
  }
}

async function handleStatusChange(row: any, val: boolean) {
  const status = val ? 1 : 2;
  try {
    await changeUserStatus({ id: row.id, status });
    MessagePlugin.success('更新状态成功');
    gridApi.reload();
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
  gridApi.reload();
}

function toggleRecycleBin() {
  isRecycleBin.value = !isRecycleBin.value;
  gridApi.reload();
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
</script>

<template>
  <Page auto-content-height>
    <div class="flex h-full flex-row gap-4">
      <div class="h-full rounded-md bg-background p-2">
        <DeptTree @select="handleDeptSelect" />
      </div>
      <div class="h-full min-w-0 flex-1 overflow-hidden p-2">
        <Grid>
          <template #toolbar-tools>
            <Button v-if="!isRecycleBin" theme="primary" @click="handleAdd">
              <template #icon><AddIcon /></template>
              新增
            </Button>
            <Button
              v-if="!isRecycleBin"
              theme="danger"
              variant="outline"
              @click="handleBatchDelete"
            >
              <template #icon><DeleteIcon /></template>
              删除
            </Button>
            <Button v-if="!isRecycleBin" variant="outline">
              <template #icon><UploadIcon /></template>
              导入
            </Button>
            <Button v-if="!isRecycleBin" variant="outline">
              <template #icon><DownloadIcon /></template>
              导出
            </Button>

            <Button
              v-if="isRecycleBin"
              theme="success"
              @click="handleBatchRecovery"
            >
              恢复
            </Button>
            <Button
              v-if="isRecycleBin"
              theme="danger"
              @click="handleBatchDelete"
            >
              彻底删除
            </Button>

            <Button variant="outline" class="ml-2" @click="toggleRecycleBin">
              {{ isRecycleBin ? '返回列表' : '显示回收站' }}
            </Button>
          </template>

          <template #status="{ row }">
            <Switch
              :disabled="isRecycleBin"
              :value="row.status === 1"
              @change="(val) => handleStatusChange(row, val as boolean)"
            />
          </template>

          <template #action="{ row }">
            <div class="flex items-center justify-center gap-2">
              <template v-if="!isRecycleBin">
                <Button
                  size="small"
                  theme="primary"
                  variant="text"
                  @click="handleEdit(row)"
                >
                  <template #icon><EditIcon /></template>
                  编辑
                </Button>
                <Popconfirm
                  content="确认删除该用户吗？"
                  @confirm="handleDelete(row)"
                >
                  <Button size="small" theme="danger" variant="text">
                    <template #icon><DeleteIcon /></template>
                    删除
                  </Button>
                </Popconfirm>
                <Dropdown
                  :options="actionDropdownOptions"
                  trigger="click"
                  @click="(dropdownItem) => handleActionDropdownClick(dropdownItem, row)"
                >
                  <Button size="small" theme="default" variant="text">
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
                  <Button size="small" theme="primary" variant="text">
                    恢复
                  </Button>
                </Popconfirm>
                <Popconfirm
                  content="确认彻底删除该用户吗？"
                  @confirm="handleDelete(row)"
                >
                  <Button size="small" theme="danger" variant="text">
                    彻底删除
                  </Button>
                </Popconfirm>
              </template>
            </div>
          </template>
        </Grid>
      </div>
    </div>

    <UserModal ref="userModalRef" @success="handleSuccess" />
  </Page>
</template>

<style scoped>
.layout-container {
  overflow: hidden;
}

:deep(.vben-vxe-grid) {
  height: 100%;
}
</style>

