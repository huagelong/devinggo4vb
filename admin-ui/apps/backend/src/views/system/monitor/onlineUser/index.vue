<script lang="ts" setup>
import type { MonitorApi } from '#/api/system/monitor';

import { onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';

import { message } from '#/adapter/tdesign';
import { getOnlineUserPageList, kickUser } from '#/api/system/monitor';

import { SearchIcon } from 'tdesign-icons-vue-next';
import { Button, Input, Space, Table } from 'tdesign-vue-next';

defineOptions({ name: 'SystemOnlineUser' });

const loading = ref(false);
const tableData = ref<MonitorApi.OnlineUserItem[]>([]);
const total = ref(0);
const searchUsername = ref('');

const columns = [
  { colKey: 'username', title: '用户账户', width: 180 },
  { colKey: 'nickname', title: '用户昵称', width: 180 },
  { colKey: 'app_id', title: 'App ID', width: 120 },
  { colKey: 'login_ip', title: '登录IP', width: 180 },
  { colKey: 'login_time', title: '登录时间', width: 180 },
  { colKey: 'action', title: '操作', width: 120 },
];

async function fetchOnlineUsers() {
  loading.value = true;
  try {
    const params: MonitorApi.OnlineUserQuery = {
      page: 1,
      page_size: 20,
    };
    if (searchUsername.value) {
      params.username = searchUsername.value;
    }
    const response = await getOnlineUserPageList(params);
    tableData.value = response.items || [];
    total.value = Number(response.pageInfo?.total || response.total || 0);
  } catch (error) {
    console.error(error);
    message.error('获取在线用户失败');
  } finally {
    loading.value = false;
  }
}

async function handleKick(row: MonitorApi.OnlineUserItem) {
  try {
    await kickUser({ id: row.id, app_id: row.app_id });
    message.success('强制退出成功');
    await fetchOnlineUsers();
  } catch (error) {
    console.error(error);
    message.error('强制退出失败');
  }
}

function handleSearch() {
  void fetchOnlineUsers();
}

function handleReset() {
  searchUsername.value = '';
  void fetchOnlineUsers();
}

onMounted(() => {
  void fetchOnlineUsers();
});
</script>

<template>
  <Page auto-content-height>
    <div class="flex h-full flex-col gap-3">
      <div class="rounded-md bg-white p-4">
        <div class="flex items-center gap-4">
          <Input
            v-model="searchUsername"
            placeholder="请输入用户名"
            clearable
            class="w-64"
            @enter="handleSearch"
          />
          <Space>
            <Button theme="primary" @click="handleSearch">
              <template #icon><SearchIcon /></template>
              查询
            </Button>
            <Button theme="default" @click="handleReset">重置</Button>
          </Space>
        </div>
      </div>

      <div class="flex min-h-0 flex-1 flex-col rounded-md bg-white p-4">
        <Table
          :columns="columns"
          :data="tableData"
          :loading="loading"
          row-key="id"
          hover
          stripe
        >
          <template #action="{ row }">
            <Button
              size="small"
              theme="danger"
              variant="outline"
              @click="handleKick(row)"
            >
              强制退出
            </Button>
          </template>
        </Table>
      </div>
    </div>
  </Page>
</template>
