<script lang="ts" setup>
import type { CrontabApi } from '#/api/system/crontab';

import { ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import { message } from '#/adapter/tdesign';
import {
  deleteCrontabLog,
  getCrontabLogPageList,
} from '#/api/system/crontab';

import {
  DeleteIcon,
  SearchIcon,
} from 'tdesign-icons-vue-next';
import {
  Button,
  DateRangePicker,
  Form,
  FormItem,
  Space,
  Table,
} from 'tdesign-vue-next';

import type { CrontabLogQuery } from '../model';

const loading = ref(false);
const tableData = ref<CrontabApi.LogItem[]>([]);
const total = ref(0);
const selectedRowKeys = ref<(number | string)[]>([]);

const searchForm = ref<CrontabLogQuery>({
  crontab_id: undefined,
  created_at: [],
});

const logColumns = [
  { colKey: 'id', title: 'ID', width: 80 },
  { colKey: 'start_time', title: '开始时间', width: 180 },
  { colKey: 'end_time', title: '结束时间', width: 180 },
  { colKey: 'duration', title: '耗时(秒)', width: 100 },
  { colKey: 'status', title: '执行结果', width: 100 },
  { colKey: 'output', title: '执行输出', width: 180 },
  { colKey: 'error', title: '异常信息', width: 180 },
  { colKey: 'created_at', title: '创建时间', width: 180 },
];

const [Modal, modalApi] = useVbenModal({
  onConfirm: () => {
    modalApi.close();
  },
  onCancel: () => {
    modalApi.close();
  },
  class: 'w-[900px]',
});

const crontabId = ref<number>();

async function fetchLogList() {
  if (!crontabId.value) return;

  loading.value = true;
  try {
    const params: CrontabApi.LogQuery = {
      page: 1,
      pageSize: 20,
      crontab_id: crontabId.value,
    };
    if (searchForm.value.created_at?.length === 2 && searchForm.value.created_at[0]) {
      params.created_at = searchForm.value.created_at;
    }

    const response = await getCrontabLogPageList(params);
    tableData.value = response.items || [];
    total.value = Number(response.pageInfo?.total || response.total || 0);
  } catch (error) {
    console.error(error);
    message.error('获取日志失败');
  } finally {
    loading.value = false;
  }
}

async function handleDeleteLog() {
  if (selectedRowKeys.value.length === 0) {
    message.warning('请选择需要删除的日志');
    return;
  }

  try {
    await deleteCrontabLog(selectedRowKeys.value.map((k) => Number(k)));
    message.success('删除成功');
    selectedRowKeys.value = [];
    await fetchLogList();
  } catch (error) {
    console.error(error);
    message.error('删除失败');
  }
}

function showErrorDetail(row: CrontabApi.LogItem) {
  modalApi.setState({ title: `异常信息 - ${row.crontab_name || row.id}` });
  // Show error in a simple message box - in real implementation could use a detail modal
  message.info(row.error || '无异常信息');
}

function handleSearch() {
  void fetchLogList();
}

function handleReset() {
  searchForm.value = {
    crontab_id: crontabId.value,
    created_at: [],
  };
  void fetchLogList();
}

function handleSelectChange(value: (number | string)[]) {
  selectedRowKeys.value = value;
}

function open(id: number) {
  crontabId.value = id;
  modalApi.setState({ title: '执行日志' });
  modalApi.open();
  searchForm.value = {
    crontab_id: id,
    created_at: [],
  };
  void fetchLogList();
}

defineExpose({
  open,
});
</script>

<template>
  <Modal>
    <div class="flex flex-col gap-3">
      <div class="rounded-md bg-white p-4">
        <Form label-width="80px" colon>
          <div class="grid grid-cols-4 gap-x-4">
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
            <Button theme="danger" variant="outline" @click="handleDeleteLog">
              <template #icon><DeleteIcon /></template>
              删除日志
            </Button>
          </Space>
        </div>

        <Table
          :columns="logColumns"
          :data="tableData"
          :loading="loading"
          :selected-row-keys="selectedRowKeys"
          row-key="id"
          hover
          stripe
          @select-change="handleSelectChange"
        >
          <template #status="{ row }">
            <span
              :class="Number(row.status) === 1 ? 'text-green-600' : 'text-red-600'"
            >
              {{ Number(row.status) === 1 ? '成功' : '失败' }}
            </span>
          </template>

          <template #output="{ row }">
            <span :title="row.output || ''">
              {{ row.output ? (row.output.length > 50 ? `${row.output.slice(0, 50)}...` : row.output) : '-' }}
            </span>
          </template>

          <template #error="{ row }">
            <span
              v-if="row.error"
              class="cursor-pointer text-red-500"
              :title="row.error"
              @click="showErrorDetail(row)"
            >
              {{ row.error.length > 30 ? `${row.error.slice(0, 30)}...` : row.error }}
            </span>
            <span v-else>-</span>
          </template>
        </Table>
      </div>
    </div>
  </Modal>
</template>
