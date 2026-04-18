<script lang="ts" setup>
import { logger } from '#/utils/logger';
import type { DeptApi } from '#/api/system/dept';
import type { UserApi } from '#/api/system/user';

import { reactive, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import {
  Button,
  Form,
  FormItem,
  Input,
  MessagePlugin,
  Popconfirm,
  Select,
  Space,
  Table,
} from 'tdesign-vue-next';

import {
  addDeptLeader,
  deleteDeptLeader,
  getDeptLeaderList,
} from '#/api/system/dept';
import { getUserList } from '#/api/system/user';

const currentDept = ref<null | { id: number; name?: string }>(null);
const loading = ref(false);
const candidateLoading = ref(false);
const leaderList = ref<DeptApi.LeaderListItem[]>([]);
const candidateUsers = ref<UserApi.ListItem[]>([]);
const selectedLeaderKeys = ref<Array<number | string>>([]);
const candidateOptions = ref<Array<{ label: string; value: number }>>([]);
const selectedCandidateIds = ref<number[]>([]);

const leaderSearchForm = reactive({
  nickname: '',
  status: undefined as number | undefined,
  username: '',
});

const candidateSearchForm = reactive({
  nickname: '',
  username: '',
});

const leaderPagination = reactive({
  current: 1,
  pageSize: 20,
  pageSizeOptions: [10, 20, 50, 100],
  showJumper: true,
  showPageSize: true,
  total: 0,
});

const statusOptions = [
  { label: '正常', value: 1 },
  { label: '停用', value: 2 },
];

const leaderColumns = [
  { colKey: 'username', title: '用户名', minWidth: 140 },
  { colKey: 'nickname', title: '用户昵称', minWidth: 140 },
  { colKey: 'phone', title: '手机', minWidth: 140 },
  { colKey: 'email', title: '邮箱', minWidth: 180 },
  { colKey: 'status', title: '状态', width: 100, align: 'center' as const },
  {
    colKey: 'leader_add_time',
    title: '设置领导时间',
    width: 180,
    align: 'center' as const,
  },
  { colKey: 'action', title: '操作', width: 120, align: 'center' as const },
];

function handleLeaderSelectChange(keys: Array<number | string>) {
  selectedLeaderKeys.value = keys;
}

async function fetchLeaderList() {
  if (!currentDept.value) return;

  loading.value = true;
  try {
    const response = await getDeptLeaderList({
      dept_id: currentDept.value.id,
      nickname: leaderSearchForm.nickname || undefined,
      page: leaderPagination.current,
      pageSize: leaderPagination.pageSize,
      status: leaderSearchForm.status,
      username: leaderSearchForm.username || undefined,
    });
    leaderList.value = response.items ?? [];
    leaderPagination.total = Number(
      response.pageInfo?.total || response.total || 0,
    );
  } catch (error) {
    logger.error(error);
    MessagePlugin.error('领导列表加载失败');
  } finally {
    loading.value = false;
  }
}

async function fetchCandidateUsers() {
  candidateLoading.value = true;
  try {
    const response = await getUserList({
      page: 1,
      pageSize: 100,
      username: candidateSearchForm.username || undefined,
    });
    candidateUsers.value = (response.items ?? []).filter((item) =>
      candidateSearchForm.nickname
        ? String(item.nickname ?? '').includes(candidateSearchForm.nickname)
        : true,
    );
    candidateOptions.value = candidateUsers.value.map((item) => ({
      label: item.nickname
        ? `${item.nickname} (${item.username})`
        : item.username,
      value: Number(item.id),
    }));
  } catch (error) {
    logger.error(error);
    MessagePlugin.error('候选用户加载失败');
  } finally {
    candidateLoading.value = false;
  }
}

async function handleAddLeaders() {
  if (!currentDept.value) return;
  if (selectedCandidateIds.value.length === 0) {
    MessagePlugin.warning('请选择需要添加的用户');
    return;
  }

  const selectedUsers = candidateUsers.value
    .filter((item) => selectedCandidateIds.value.includes(Number(item.id)))
    .map((item) => ({
      nickname: item.nickname,
      user_id: Number(item.id),
      username: item.username,
    }));

  if (selectedUsers.length === 0) {
    MessagePlugin.warning('未找到可添加的用户信息');
    return;
  }

  try {
    await addDeptLeader({
      id: currentDept.value.id,
      users: selectedUsers,
    });
    MessagePlugin.success('添加领导成功');
    selectedCandidateIds.value = [];
    await fetchLeaderList();
  } catch (error) {
    logger.error(error);
    MessagePlugin.error('添加领导失败');
  }
}

async function handleDeleteLeader(id: number) {
  if (!currentDept.value) return;
  try {
    await deleteDeptLeader({
      id: currentDept.value.id,
      ids: [id],
    });
    MessagePlugin.success('删除成功');
    await fetchLeaderList();
  } catch (error) {
    logger.error(error);
    MessagePlugin.error('删除失败');
  }
}

async function handleBatchDeleteLeaders() {
  if (!currentDept.value) return;
  if (selectedLeaderKeys.value.length === 0) {
    MessagePlugin.warning('请选择需要删除的领导');
    return;
  }
  try {
    await deleteDeptLeader({
      id: currentDept.value.id,
      ids: selectedLeaderKeys.value.map((item) => Number(item)),
    });
    MessagePlugin.success('删除成功');
    selectedLeaderKeys.value = [];
    await fetchLeaderList();
  } catch (error) {
    logger.error(error);
    MessagePlugin.error('删除失败');
  }
}

function handleLeaderSearch() {
  leaderPagination.current = 1;
  void fetchLeaderList();
}

function handleLeaderReset() {
  leaderSearchForm.nickname = '';
  leaderSearchForm.status = undefined;
  leaderSearchForm.username = '';
  leaderPagination.current = 1;
  void fetchLeaderList();
}

function handleCandidateSearch() {
  void fetchCandidateUsers();
}

function handleLeaderPageChange(pageInfo: { current: number; pageSize: number }) {
  leaderPagination.current = pageInfo.current;
  leaderPagination.pageSize = pageInfo.pageSize;
  void fetchLeaderList();
}

const [Modal, modalApi] = useVbenModal({
  footer: false,
  class: 'w-[1200px]',
});

async function open(row: { id: number; name?: string }) {
  currentDept.value = row;
  selectedLeaderKeys.value = [];
  selectedCandidateIds.value = [];
  candidateOptions.value = [];
  candidateUsers.value = [];
  candidateSearchForm.nickname = '';
  candidateSearchForm.username = '';
  leaderSearchForm.nickname = '';
  leaderSearchForm.status = undefined;
  leaderSearchForm.username = '';
  leaderPagination.current = 1;

  modalApi.setState({
    title: '部门领导列表',
  });
  modalApi.open();

  await Promise.all([fetchLeaderList(), fetchCandidateUsers()]);
}

defineExpose({
  open,
});
</script>

<template>
  <Modal>
    <div class="flex flex-col gap-4">
      <div class="rounded-md border border-yellow-200 bg-yellow-50 px-3 py-2 text-sm text-yellow-700">
        部门的领导人可以跨部门设置
      </div>

      <div class="rounded-md border border-gray-100 bg-gray-50 p-4">
        <div class="mb-3 text-sm font-medium text-gray-700">新增领导</div>
        <Form :data="candidateSearchForm" label-width="90px" colon>
          <div class="grid grid-cols-4 gap-x-4">
            <FormItem label="用户名" name="username">
              <Input
                v-model="candidateSearchForm.username"
                placeholder="请输入用户名"
                clearable
              />
            </FormItem>
            <FormItem label="用户昵称" name="nickname">
              <Input
                v-model="candidateSearchForm.nickname"
                placeholder="请输入用户昵称"
                clearable
              />
            </FormItem>
            <FormItem label="候选用户" name="candidate_ids" class="col-span-2">
              <Select
                v-model="selectedCandidateIds"
                :loading="candidateLoading"
                :options="candidateOptions"
                multiple
                clearable
                filterable
                placeholder="请选择要添加为领导的用户"
              />
            </FormItem>
          </div>
          <div class="flex justify-end gap-2 pt-2">
            <Button theme="default" @click="handleCandidateSearch">查询候选用户</Button>
            <Button theme="primary" @click="handleAddLeaders">新增领导</Button>
          </div>
        </Form>
      </div>

      <div class="rounded-md border border-gray-100 bg-white p-4">
        <Form :data="leaderSearchForm" label-width="90px" colon>
          <div class="grid grid-cols-4 gap-x-4">
            <FormItem label="用户名" name="username">
              <Input
                v-model="leaderSearchForm.username"
                placeholder="请输入用户名"
                clearable
              />
            </FormItem>
            <FormItem label="用户昵称" name="nickname">
              <Input
                v-model="leaderSearchForm.nickname"
                placeholder="请输入用户昵称"
                clearable
              />
            </FormItem>
            <FormItem label="状态" name="status">
              <Select
                v-model="leaderSearchForm.status"
                :options="statusOptions"
                placeholder="请选择状态"
                clearable
              />
            </FormItem>
          </div>
          <div class="flex justify-end gap-2 pt-2">
            <Button theme="default" @click="handleLeaderReset">重置</Button>
            <Button theme="primary" @click="handleLeaderSearch">查询</Button>
          </div>
        </Form>

        <div class="mb-3 mt-4">
          <Space>
            <Button theme="danger" variant="outline" @click="handleBatchDeleteLeaders">
              删除选中
            </Button>
          </Space>
        </div>

        <Table
          :columns="leaderColumns"
          :data="leaderList"
          :loading="loading"
          :pagination="leaderPagination"
          :selected-row-keys="selectedLeaderKeys"
          row-key="id"
          hover
          stripe
          @page-change="handleLeaderPageChange"
          @select-change="handleLeaderSelectChange"
        >
          <template #status="{ row }">
            <span>{{ Number(row.status) === 1 ? '正常' : '停用' }}</span>
          </template>

          <template #action="{ row }">
            <Popconfirm
              content="确认删除该领导吗？"
              @confirm="handleDeleteLeader(Number(row.id))"
            >
              <Button size="small" theme="danger" variant="outline">
                删除
              </Button>
            </Popconfirm>
          </template>
        </Table>
      </div>
    </div>
  </Modal>
</template>
