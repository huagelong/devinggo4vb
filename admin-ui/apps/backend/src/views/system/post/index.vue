<script lang="ts" setup>
import { computed, onMounted, ref } from 'vue';
import { Page } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { message } from '#/adapter/tdesign';
import { logger } from '#/utils/logger';

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
  InputNumber,
  Popconfirm,
  Select,
  Space,
  Switch,
  Table,
} from 'tdesign-vue-next';

import CrudToolbar from '#/components/crud/crud-toolbar.vue';
import {
  changePostStatus,
  deletePost,
  realDeletePost,
  recoveryPost,
  updatePostSort,
} from '#/api/system/post';
import type { DictOption } from '#/composables/crud/use-dict-options';
import { useDictOptions } from '#/composables/crud/use-dict-options';

import PostModal from './components/post-modal.vue';
import type { PostListItem, PostTableColumn } from './model';
import { createPostColumnOptions, createPostTableColumns } from './schemas';
import { usePostCrud } from './use-post-crud';

defineOptions({ name: 'SystemPost' });

type PostModalInstance = {
  open: (data?: Partial<PostListItem>) => void;
};

const postModalRef = ref<PostModalInstance>();
const statusOptions = ref<DictOption[]>([]);

const columns: PostTableColumn[] = createPostTableColumns();
const columnOptions = createPostColumnOptions(columns);
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
} = usePostCrud();

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
          { label: '正常', value: 1 },
          { label: '停用', value: 2 },
        ];
}

function handleAdd() {
  postModalRef.value?.open();
}

function handleEdit(row: PostListItem) {
  postModalRef.value?.open(row);
}

async function handleDelete(row: PostListItem) {
  try {
    await (isRecycleBin.value ? realDeletePost([row.id]) : deletePost([row.id]));
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
  try {
    await (isRecycleBin.value ? realDeletePost(ids) : deletePost(ids));
    message.success($t('common.operationSuccess'));
    clearSelectedRowKeys();
    await fetchTableData();
  } catch (error) {
    logger.error(error);
    message.error($t('common.batchDeleteFailed'));
  }
}

async function handleRecovery(row: PostListItem) {
  try {
    await recoveryPost([row.id]);
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
    await recoveryPost(ids);
    message.success($t('common.recoverySuccess'));
    clearSelectedRowKeys();
    await fetchTableData();
  } catch (error) {
    logger.error(error);
    message.error($t('common.batchRecoveryFailed'));
  }
}

async function handleStatusChange(row: PostListItem, checked: boolean) {
  const status = checked ? 1 : 2;
  try {
    await changePostStatus({ id: row.id, status });
    message.success($t('common.statusUpdateSuccess'));
    await fetchTableData();
  } catch (error) {
    logger.error(error);
    message.error($t('common.statusUpdateFailed'));
  }
}

async function handleSortChange(value: number | string, row: PostListItem) {
  const numberValue = Number(value);
  if (Number.isNaN(numberValue)) return;

  try {
    await updatePostSort({
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

function handleStatusSwitchChange(row: PostListItem, value: unknown) {
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
            <FormItem label="岗位名称" name="name">
              <Input
                v-model="searchForm.name"
                placeholder="请输入岗位名称"
                clearable
              />
            </FormItem>
            <FormItem label="岗位标识" name="code">
              <Input
                v-model="searchForm.code"
                placeholder="请输入岗位标识"
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
              :disabled="isRecycleBin"
              :value="row.status === 1"
              @change="(value: unknown) => handleStatusSwitchChange(row, value)"
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
                  content="确认删除该岗位吗？"
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
                  content="确认恢复该岗位吗？"
                  @confirm="handleRecovery(row)"
                >
                  <Button size="small" theme="primary" variant="outline">
                    恢复
                  </Button>
                </Popconfirm>
                <Popconfirm
                  content="确认彻底删除该岗位吗？"
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

    <PostModal ref="postModalRef" @success="handleSuccess" />
  </Page>
</template>
