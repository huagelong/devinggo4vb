<script lang="ts" setup>
import { logger } from '#/utils/logger';
import type { DeptApi } from '#/api/system/dept';
import type { UserApi } from '#/api/system/user';

import { reactive, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';
import { $t } from '@vben/locales';

import {
  Button,
  Dialog,
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
  getDeptTree,
} from '#/api/system/dept';
import { getPostList } from '#/api/system/post';
import { getRoleList } from '#/api/system/role';
import { getUserList } from '#/api/system/user';

const currentDept = ref<null | { id: number; name?: string }>(null);
const loading = ref(false);
const candidateLoading = ref(false);
const leaderList = ref<DeptApi.LeaderListItem[]>([]);
const candidateUsers = ref<UserApi.ListItem[]>([]);
const selectedLeaderKeys = ref<Array<number | string>>([]);
const selectedCandidateIds = ref<Array<number | string>>([]);
const deptOptions = ref<Array<{ label: string; value: number }>>([]);
const roleOptions = ref<Array<{ label: string; value: number }>>([]);
const postOptions = ref<Array<{ label: string; value: number }>>([]);
const candidateColumns = [
  { colKey: 'row-select', type: 'multiple', width: 50 },
  { colKey: 'username', title: $t('system.user.username'), minWidth: 120 },
  { colKey: 'nickname', title: $t('system.user.nickname'), minWidth: 120 },
  { colKey: 'phone', title: $t('system.user.phone'), minWidth: 140 },
  { colKey: 'email', title: $t('system.user.email'), minWidth: 180 },
  { colKey: 'dept_name', title: $t('system.user.dept'), minWidth: 160 },
  { colKey: 'role_name', title: $t('system.user.role'), minWidth: 140 },
  { colKey: 'post_name', title: $t('system.user.post'), minWidth: 140 },
];
const candidatePagination = reactive({
  current: 1,
  pageSize: 10,
  pageSizeOptions: [10, 20, 50, 100],
  showJumper: true,
  showPageSize: true,
  total: 0,
});

const leaderSearchForm = reactive({
  nickname: '',
  status: undefined as number | undefined,
  username: '',
});

const candidateSearchForm = reactive({
  username: '',
  nickname: '',
  phone: '',
  email: '',
  dept_id: undefined as number | undefined,
  role_id: undefined as number | undefined,
  post_id: undefined as number | undefined,
});

const leaderPagination = reactive({
  current: 1,
  pageSize: 20,
  pageSizeOptions: [10, 20, 50, 100],
  showJumper: true,
  showPageSize: true,
  total: 0,
});

const addLeaderDialogVisible = ref(false);

const statusOptions = [
  { label: $t('common.statusEnabled'), value: 1 },
  { label: $t('common.statusDisabled'), value: 2 },
];

const leaderColumns = [
  { colKey: 'row-select', type: 'multiple', width: 50 },
  { colKey: 'username', title: $t('system.dept.username'), minWidth: 140 },
  { colKey: 'nickname', title: $t('system.dept.nickname'), minWidth: 140 },
  { colKey: 'phone', title: $t('system.dept.phone'), minWidth: 140 },
  { colKey: 'email', title: $t('system.dept.email'), minWidth: 180 },
  { colKey: 'status', title: $t('common.status'), width: 100, align: 'center' as const },
  {
    colKey: 'leader_add_time',
    title: $t('system.dept.leaderAddTime'),
    width: 180,
    align: 'center' as const,
  },
  { colKey: 'action', title: $t('common.action'), width: 120, align: 'center' as const },
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
    MessagePlugin.error($t('common.leaderListLoadFailed'));
  } finally {
    loading.value = false;
  }
}

async function fetchCandidateUsers(page = 1) {
  candidateLoading.value = true;
  try {
    const params: Record<string, unknown> = {
      page,
      pageSize: candidatePagination.pageSize,
      username: candidateSearchForm.username || undefined,
      nickname: candidateSearchForm.nickname || undefined,
      phone: candidateSearchForm.phone || undefined,
      email: candidateSearchForm.email || undefined,
      dept_id: candidateSearchForm.dept_id,
      role_id: candidateSearchForm.role_id,
      post_id: candidateSearchForm.post_id,
    };
    const response = await getUserList(params as unknown as UserApi.ListQuery);
    candidateUsers.value = response.items ?? [];
    candidatePagination.total = Number(response.pageInfo?.total || response.total || 0);
  } catch (error) {
    logger.error(error);
    MessagePlugin.error($t('common.candidateUserLoadFailed'));
  } finally {
    candidateLoading.value = false;
  }
}

async function handleAddLeaders() {
  if (!currentDept.value) return;
  if (selectedCandidateIds.value.length === 0) {
    MessagePlugin.warning($t('common.selectUserToAdd'));
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
    MessagePlugin.warning($t('common.noUserFoundToAdd'));
    return;
  }

  try {
    await addDeptLeader({
      id: currentDept.value.id,
      users: selectedUsers,
    });
    MessagePlugin.success($t('common.addLeaderSuccess'));
    selectedCandidateIds.value = [];
    closeAddLeaderDialog();
    await Promise.all([fetchLeaderList(), fetchCandidateUsers(1)]);
  } catch (error) {
    logger.error(error);
    MessagePlugin.error($t('common.addLeaderFailed'));
  }
}

async function handleDeleteLeader(id: number) {
  if (!currentDept.value) return;
  try {
    await deleteDeptLeader({
      id: currentDept.value.id,
      ids: [id],
    });
    MessagePlugin.success($t('common.deleteSuccess'));
    await fetchLeaderList();
  } catch (error) {
    logger.error(error);
    MessagePlugin.error($t('common.deleteFailed'));
  }
}

async function handleBatchDeleteLeaders() {
  if (!currentDept.value) return;
  if (selectedLeaderKeys.value.length === 0) {
    MessagePlugin.warning($t('common.selectLeaderToDelete'));
    return;
  }
  try {
    await deleteDeptLeader({
      id: currentDept.value.id,
      ids: selectedLeaderKeys.value.map((item) => Number(item)),
    });
    MessagePlugin.success($t('common.deleteSuccess'));
    selectedLeaderKeys.value = [];
    await fetchLeaderList();
  } catch (error) {
    logger.error(error);
    MessagePlugin.error($t('common.deleteFailed'));
  }
}

async function fetchDeptOptions() {
  try {
    const response = await getDeptTree();
    const flatten = (nodes: Array<{ id: number; name?: string; children?: any[] }>) =>
      nodes.reduce<Array<{ label: string; value: number }>>((acc, node) => {
        acc.push({ label: node.name || '', value: node.id });
        if (node.children?.length) {
          acc.push(...flatten(node.children));
        }
        return acc;
      }, []);
    deptOptions.value = flatten(response);
  } catch (error) {
    logger.error(error);
  }
}

async function fetchRoleOptions() {
  try {
    const response = await getRoleList();
    const items = Array.isArray(response) ? response : response.items ?? [];
    roleOptions.value = items.map((item) => ({ label: item.name, value: item.id }));
  } catch (error) {
    logger.error(error);
  }
}

async function fetchPostOptions() {
  try {
    const response = await getPostList();
    const items = Array.isArray(response) ? response : response.items ?? [];
    postOptions.value = items.map((item) => ({ label: item.name, value: item.id }));
  } catch (error) {
    logger.error(error);
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
  candidatePagination.current = 1;
  void fetchCandidateUsers(1);
}

function handleCandidateReset() {
  candidateSearchForm.username = '';
  candidateSearchForm.nickname = '';
  candidateSearchForm.phone = '';
  candidateSearchForm.email = '';
  candidateSearchForm.dept_id = undefined;
  candidateSearchForm.role_id = undefined;
  candidateSearchForm.post_id = undefined;
  candidatePagination.current = 1;
  void fetchCandidateUsers(1);
}

function openAddLeaderDialog() {
  addLeaderDialogVisible.value = true;
  selectedCandidateIds.value = [];
  void fetchCandidateUsers(1);
}

function closeAddLeaderDialog() {
  addLeaderDialogVisible.value = false;
}

function handleLeaderPageChange(pageInfo: { current: number; pageSize: number }) {
  leaderPagination.current = pageInfo.current;
  leaderPagination.pageSize = pageInfo.pageSize;
  void fetchLeaderList();
}

function handleCandidatePageChange(pageInfo: { current: number; pageSize: number }) {
  candidatePagination.current = pageInfo.current;
  candidatePagination.pageSize = pageInfo.pageSize;
  void fetchCandidateUsers(pageInfo.current);
}

function handleCandidateSelectChange(keys: Array<number | string>) {
  selectedCandidateIds.value = keys;
}

const [Modal, modalApi] = useVbenModal({
  footer: false,
  class: 'w-[1400px] h-[700px]',
});

async function open(row: { id: number; name?: string }) {
  currentDept.value = row;
  selectedLeaderKeys.value = [];
  selectedCandidateIds.value = [];
  candidateUsers.value = [];
  candidateSearchForm.username = '';
  candidateSearchForm.nickname = '';
  candidateSearchForm.phone = '';
  candidateSearchForm.email = '';
  candidateSearchForm.dept_id = undefined;
  candidateSearchForm.role_id = undefined;
  candidateSearchForm.post_id = undefined;
  leaderSearchForm.nickname = '';
  leaderSearchForm.status = undefined;
  leaderSearchForm.username = '';
  leaderPagination.current = 1;

  modalApi.setState({
    title: $t('system.dept.leaderManage'),
  });
  modalApi.open();

  await Promise.all([
    fetchLeaderList(),
    fetchCandidateUsers(1),
    fetchDeptOptions(),
    fetchRoleOptions(),
    fetchPostOptions(),
  ]);
}

defineExpose({
  open,
});
</script>

<template>
  <Modal>
    <div class="flex flex-col gap-4">
      <!-- Tip Alert -->
      <div class="rounded-md border border-yellow-200 bg-yellow-50 px-3 py-2 text-sm text-yellow-700">
        {{ $t('system.dept.crossDeptLeaderTip') }}
      </div>

      <!-- Search Form -->
      <div class="rounded-md border border-gray-100 bg-white p-4">
<Form :data="leaderSearchForm" label-width="90px" colon size="small">
            <div class="grid grid-cols-3 gap-x-4">
              <FormItem :label="$t('system.dept.username')" name="username">
              <Input
                v-model="leaderSearchForm.username"
                :placeholder="$t('ui.placeholder.input', [$t('system.dept.username')])"
                clearable
                size="small"
              />
            </FormItem>
            <FormItem :label="$t('system.dept.nickname')" name="nickname">
              <Input
                v-model="leaderSearchForm.nickname"
                :placeholder="$t('ui.placeholder.input', [$t('system.dept.nickname')])"
                clearable
                size="small"
              />
            </FormItem>
            <FormItem :label="$t('common.status')" name="status">
              <Select
                v-model="leaderSearchForm.status"
                :options="statusOptions"
                :placeholder="$t('ui.placeholder.select', [$t('common.status')])"
                clearable
                size="small"
              />
            </FormItem>
          </div>

          <div class="flex justify-end gap-2 pt-3">
            <Button size="small" theme="default" @click="handleLeaderReset">{{ $t('common.reset') }}</Button>
            <Button size="small" theme="primary" @click="handleLeaderSearch">
              {{ $t('common.search') }}
            </Button>
          </div>
        </Form>
      </div>

      <!-- Action Toolbar -->
      <div class="flex items-center justify-between rounded-md border border-gray-100 bg-white px-4 py-3">
        <Space>
          <Button size="small" theme="primary" @click="openAddLeaderDialog">
            {{ $t('system.dept.addLeader') }}
          </Button>
          <Button size="small" theme="danger" variant="outline" @click="handleBatchDeleteLeaders">
            {{ $t('system.dept.deleteSelected') }}
          </Button>
        </Space>
        <Space>
          <Button size="small" theme="default" variant="outline" @click="fetchLeaderList">
            {{ $t('common.refresh') }}
          </Button>
        </Space>
      </div>

      <!-- Leader Table -->
      <div class="rounded-md border border-gray-100 bg-white p-4">
        <Table
          :columns="leaderColumns"
          :data="leaderList"
          :loading="loading"
          :pagination="leaderPagination"
          :selected-row-keys="selectedLeaderKeys"
          row-key="id"
          hover
          stripe
          size="small"
          @page-change="handleLeaderPageChange"
          @select-change="handleLeaderSelectChange"
        >
          <template #status="{ row }">
            <span>{{ Number(row.status) === 1 ? $t('common.statusEnabled') : $t('common.statusDisabled') }}</span>
          </template>

          <template #action="{ row }">
            <Popconfirm
              :content="$t('system.dept.confirmDeleteLeader')"
              @confirm="handleDeleteLeader(Number(row.id))"
            >
              <Button size="small" theme="danger" variant="outline">
                {{ $t('common.delete') }}
              </Button>
            </Popconfirm>
          </template>
        </Table>
      </div>

      <Dialog
        v-model:visible="addLeaderDialogVisible"
        width="1200px"
        placement="center"
        title="$t('system.dept.addLeader')"
        :footer="false"
      >
        <div class="space-y-4">
          <Form :data="candidateSearchForm" label-width="90px" colon size="small">
            <div class="grid grid-cols-4 gap-x-4">
              <FormItem :label="$t('system.user.username')" name="username">
                <Input
                  v-model="candidateSearchForm.username"
                  :placeholder="$t('ui.placeholder.input', [$t('system.user.username')])"
                  clearable
                  size="small"
                />
              </FormItem>
              <FormItem :label="$t('system.user.nickname')" name="nickname">
                <Input
                  v-model="candidateSearchForm.nickname"
                  :placeholder="$t('ui.placeholder.input', [$t('system.user.nickname')])"
                  clearable
                  size="small"
                />
              </FormItem>
              <FormItem :label="$t('system.user.phone')" name="phone">
                <Input
                  v-model="candidateSearchForm.phone"
                  :placeholder="$t('ui.placeholder.input', [$t('system.user.phone')])"
                  clearable
                  size="small"
                />
              </FormItem>
              <FormItem :label="$t('system.user.email')" name="email">
                <Input
                  v-model="candidateSearchForm.email"
                  :placeholder="$t('ui.placeholder.input', [$t('system.user.email')])"
                  clearable
                  size="small"
                />
              </FormItem>
            </div>
            <div class="grid grid-cols-3 gap-x-4 mt-4">
              <FormItem :label="$t('system.user.dept')" name="dept_id">
                <Select
                  v-model="candidateSearchForm.dept_id"
                  :options="deptOptions"
                  clearable
                  size="small"
                  :placeholder="$t('ui.placeholder.select', [$t('system.user.dept')])"
                />
              </FormItem>
              <FormItem :label="$t('system.user.role')" name="role_id">
                <Select
                  v-model="candidateSearchForm.role_id"
                  :options="roleOptions"
                  clearable
                  size="small"
                  :placeholder="$t('ui.placeholder.select', [$t('system.user.role')])"
                />
              </FormItem>
              <FormItem :label="$t('system.user.post')" name="post_id">
                <Select
                  v-model="candidateSearchForm.post_id"
                  :options="postOptions"
                  clearable
                  size="small"
                  :placeholder="$t('ui.placeholder.select', [$t('system.user.post')])"
                />
              </FormItem>
            </div>
            <div class="flex justify-end gap-2 pt-4">
              <Button size="small" theme="default" @click="handleCandidateReset">
                {{ $t('common.reset') }}
              </Button>
              <Button size="small" theme="primary" @click="handleCandidateSearch">
                {{ $t('common.search') }}
              </Button>
            </div>
          </Form>

          <div class="rounded-md border border-gray-100 bg-white p-3">
            <Table
              :columns="candidateColumns"
              :data="candidateUsers"
              :loading="candidateLoading"
              :pagination="candidatePagination"
              :selected-row-keys="selectedCandidateIds"
              row-key="id"
              hover
              stripe
              size="small"
              @page-change="handleCandidatePageChange"
              @select-change="handleCandidateSelectChange"
            />
          </div>

          <div class="flex justify-end gap-2">
            <Button size="small" theme="default" variant="outline" @click="closeAddLeaderDialog">
              {{ $t('common.cancel') }}
            </Button>
            <Button size="small" theme="primary" @click="handleAddLeaders">
              {{ $t('common.confirm') }}
            </Button>
          </div>
        </div>
      </Dialog>
    </div>
  </Modal>
</template>
