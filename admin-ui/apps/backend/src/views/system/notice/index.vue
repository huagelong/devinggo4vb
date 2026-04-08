<script lang="ts" setup>
import { computed, onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { message } from '#/adapter/tdesign';

import {
  AddIcon,
  DeleteIcon,
  EditIcon,
  SearchIcon,
} from 'tdesign-icons-vue-next';
import {
  Button,
  DateRangePicker,
  Form,
  FormItem,
  Input,
  Popconfirm,
  Select,
  Space,
  Table,
} from 'tdesign-vue-next';

import CrudToolbar from '#/components/crud/crud-toolbar.vue';
import {
  deleteNotice,
  realDeleteNotice,
  recoveryNotice,
} from '#/api/system/notice';
import type { DictOption } from '#/composables/crud/use-dict-options';
import { useDictOptions } from '#/composables/crud/use-dict-options';

import NoticeModal from './components/notice-modal.vue';
import type { NoticeListItem, NoticeTableColumn } from './model';
import {
  createNoticeColumnOptions,
  createNoticeTableColumns,
} from './schemas';
import { useNoticeCrud } from './use-notice-crud';

defineOptions({ name: 'SystemNotice' });

type NoticeModalInstance = {
  open: (data?: NoticeListItem) => void;
};

const noticeModalRef = ref<NoticeModalInstance>();
const fallbackNoticeTypeOptions: DictOption[] = [
  { label: '通知', value: 1 },
  { label: '公告', value: 2 },
];

function normalizeNoticeTypeOptions(options: DictOption[]) {
  return (options || []).map((item) => {
    const numericValue = Number(item.value);
    return Number.isNaN(numericValue) ? { ...item } : { ...item, value: numericValue };
  });
}

const noticeTypeOptions = ref<DictOption[]>([]);

const columns: NoticeTableColumn[] = createNoticeTableColumns();
const columnOptions = createNoticeColumnOptions(columns);
const allColumnKeys = columnOptions.map((item) => item.value);
const visibleColumns = ref<string[]>([...allColumnKeys]);
const displayColumns = computed({
  get: () => ['row-select', ...visibleColumns.value],
  set: (value: string[]) => {
    visibleColumns.value = value.filter((item) => item !== 'row-select');
  },
});

const noticeTypeMap = computed(() => {
  const map = new Map<string, string>();
  noticeTypeOptions.value.forEach((item) => {
    map.set(String(item.value), item.label ?? '');
  });
  return map;
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
} = useNoticeCrud();

const { getDictOptions } = useDictOptions();

async function fetchNoticeTypeDict() {
  const options = await getDictOptions('backend_notice_type');
  noticeTypeOptions.value = normalizeNoticeTypeOptions(
    options.length > 0 ? options : fallbackNoticeTypeOptions,
  );
}

function resolveNoticeTypeLabel(value?: number | string) {
  return noticeTypeMap.value.get(String(value ?? '')) || '-';
}

function normalizeId(value: number | string) {
  const numericValue = Number(value);
  return Number.isNaN(numericValue) ? undefined : numericValue;
}

function toIds(keys: Array<number | string>) {
  return keys
    .map((key) => normalizeId(key))
    .filter((id): id is number => typeof id === 'number');
}

function handleAdd() {
  noticeModalRef.value?.open();
}

function handleEdit(row: NoticeListItem) {
  noticeModalRef.value?.open(row);
}

async function handleDelete(row: NoticeListItem) {
  try {
    await (isRecycleBin.value ? realDeleteNotice([row.id]) : deleteNotice([row.id]));
    message.success('操作成功');
    await fetchTableData();
  } catch (error) {
    console.error(error);
    message.error(isRecycleBin.value ? '彻底删除失败，请稍后重试' : '删除失败，请稍后重试');
  }
}

async function handleBatchDelete() {
  if (selectedRowKeys.value.length === 0) {
    message.warning('请选择需要操作的数据');
    return;
  }
  const ids = toIds(selectedRowKeys.value);
  if (ids.length === 0) {
    message.warning('所选数据格式异常，请重试');
    return;
  }
  try {
    await (isRecycleBin.value ? realDeleteNotice(ids) : deleteNotice(ids));
    message.success('操作成功');
    clearSelectedRowKeys();
    await fetchTableData();
  } catch (error) {
    console.error(error);
    message.error(isRecycleBin.value ? '批量彻底删除失败，请稍后重试' : '批量删除失败，请稍后重试');
  }
}

async function handleRecovery(row: NoticeListItem) {
  try {
    await recoveryNotice([row.id]);
    message.success('恢复成功');
    await fetchTableData();
  } catch (error) {
    console.error(error);
    message.error('恢复失败，请稍后重试');
  }
}

async function handleBatchRecovery() {
  if (selectedRowKeys.value.length === 0) {
    message.warning('请选择需要操作的数据');
    return;
  }
  const ids = toIds(selectedRowKeys.value);
  if (ids.length === 0) {
    message.warning('所选数据格式异常，请重试');
    return;
  }
  try {
    await recoveryNotice(ids);
    message.success('恢复成功');
    clearSelectedRowKeys();
    await fetchTableData();
  } catch (error) {
    console.error(error);
    message.error('批量恢复失败，请稍后重试');
  }
}

function handleSuccess() {
  void fetchTableData();
}

function handleTableSelectChange(keys: Array<number | string>) {
  handleSelectChange(keys);
}

onMounted(() => {
  void fetchNoticeTypeDict();
  void fetchTableData();
});
</script>

<template>
  <Page auto-content-height>
    <div class="flex h-full flex-col gap-3">
      <div class="rounded-md bg-white p-4">
        <Form :data="searchForm" label-width="80px" colon>
          <div class="grid grid-cols-4 gap-x-4">
            <FormItem label="公告标题" name="title">
              <Input
                v-model="searchForm.title"
                placeholder="请输入公告标题"
                clearable
              />
            </FormItem>
            <FormItem label="公告类型" name="type">
              <Select
                v-model="searchForm.type"
                :options="noticeTypeOptions"
                placeholder="请选择公告类型"
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
          <template #type="{ row }">
            <span>{{ resolveNoticeTypeLabel(row.type) }}</span>
          </template>

          <template #remark="{ row }">
            <span class="block max-w-[320px] truncate">
              {{ row.remark || '-' }}
            </span>
          </template>

          <template #content="{ row }">
            <span
              class="block max-w-[320px] truncate text-gray-600 text-sm"
              v-html="row.content || '-'"
            />
          </template>

          <template #created_at="{ row }">
            <span>{{ row.created_at || '-' }}</span>
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
                  content="确认删除该公告吗？"
                  @confirm="handleDelete(row)"
                >
                  <Button size="small" theme="danger" variant="outline">
                    <template #icon><DeleteIcon /></template>
                    删除
                  </Button>
                </Popconfirm>
              </template>
              <template v-else>
                <Popconfirm
                  content="确认恢复该公告吗？"
                  @confirm="handleRecovery(row)"
                >
                  <Button size="small" theme="success" variant="outline">
                    恢复
                  </Button>
                </Popconfirm>
                <Popconfirm
                  content="确认彻底删除该公告吗？"
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

    <NoticeModal ref="noticeModalRef" @success="handleSuccess" />
  </Page>
</template>
