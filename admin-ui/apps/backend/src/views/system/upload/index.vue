<script lang="ts" setup>
import type { UploadListItem, UploadTreeItem } from './model';

import { computed, onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';

import { message } from '#/adapter/tdesign';
import CrudToolbar from '#/components/crud/crud-toolbar.vue';

import {
  DeleteIcon,
  LinkIcon,
  SearchIcon,
  UploadIcon,
} from 'tdesign-icons-vue-next';
import {
  Button,
  DateRangePicker,
  Form,
  FormItem,
  Input,
  MessagePlugin,
  Popconfirm,
  Radio,
  RadioGroup,
  Space,
  Tag,
  Tree,
  Upload,
} from 'tdesign-vue-next';

import type { UploadTableColumn } from './model';
import {
  createUploadColumnOptions,
  createUploadSearchForm,
  createUploadTableColumns,
  defaultUploadTreeData,
} from './schemas';
import { useUploadCrud } from './use-upload-crud';

defineOptions({ name: 'SystemUpload' });

const selectedTreeKey = ref<string[]>(['all']);
const treeData = ref<UploadTreeItem[]>(defaultUploadTreeData);
const uploadVisible = ref(false);
const uploadingFiles = ref(0);

const searchForm = ref(createUploadSearchForm());

const columns: UploadTableColumn[] = createUploadTableColumns();
const columnOptions = createUploadColumnOptions(columns);
const allColumnKeys = columnOptions.map((item) => item.value);
const visibleColumns = ref<string[]>(
  allColumnKeys.filter((key) => key !== 'storage_path')
);

const displayColumns = computed({
  get: () => ['row-select', ...visibleColumns.value],
  set: (value: string[]) => {
    visibleColumns.value = value.filter((item) => item !== 'row-select');
  },
});

const {
  clearSelectedRowKeys,
  fetchTableData,
  handleDownload,
  handlePageChange,
  handleReset,
  handleSearch,
  handleSelectChange,
  handleView,
  loading,
  pagination,
  selectedRowKeys,
  tableData,
} = useUploadCrud();

function toIds(keys: Array<number | string>) {
  return keys.map((key) => Number(key));
}

function handleTreeChange(value: Array<string | number>) {
  const keys = value.map((item) => String(item));
  selectedTreeKey.value = keys.length > 0 ? keys : ['all'];
  const key = selectedTreeKey.value[0];
  if (key === 'all') {
    searchForm.value.mime_type = '';
  } else {
    searchForm.value.mime_type = key;
  }
  handleSearch();
}

async function handleUpload(file: File) {
  uploadingFiles.value++;
  try {
    // TODO: 实际应该调用 upload API
    // const res = await uploadFileApi(file);
    // if (res?.url) {
    //   MessagePlugin.success('上传成功');
    //   await fetchTableData();
    // }

    // 临时模拟上传
    await new Promise((resolve) => setTimeout(resolve, 1000));
    MessagePlugin.success('上传成功');
  } catch (error) {
    console.error(error);
    MessagePlugin.error('上传失败');
  } finally {
    uploadingFiles.value--;
  }
}

function handleBatchDelete() {
  if (selectedRowKeys.value.length === 0) {
    message.warning('请选择需要删除的数据');
    return;
  }
  // TODO: 实现批量删除
  message.info('批量删除功能待实现');
}

async function handleDelete(row: UploadListItem) {
  try {
    // TODO: 实际应该调用删除 API
    // await deleteUploadApi([row.id]);
    message.success('删除成功');
    await fetchTableData();
  } catch (error) {
    console.error(error);
    message.error('删除失败');
  }
}

function isImageType(mimeType: string): boolean {
  return /^image\//.test(mimeType);
}

function getFileIcon(mimeType: string): string {
  if (isImageType(mimeType)) {
    return 'i-lucide:image';
  }
  if (mimeType.includes('pdf')) {
    return 'i-lucide:file-text';
  }
  if (mimeType.includes('zip') || mimeType.includes('rar')) {
    return 'i-lucide:archive';
  }
  return 'i-lucide:file';
}

onMounted(() => {
  void fetchTableData();
});
</script>

<template>
  <Page auto-content-height>
    <div class="flex h-full gap-3">
      <!-- Left Tree Slider -->
      <div class="w-48 flex-shrink-0 rounded-md bg-white p-2">
        <div class="mb-2 px-2 text-sm font-medium text-gray-500">文件类型</div>
        <Tree
          v-model="selectedTreeKey"
          :data="treeData"
          hover
          expand-all
          @change="handleTreeChange"
        />
      </div>

      <!-- Main Content -->
      <div class="flex min-h-0 flex-1 flex-col gap-3">
        <!-- Upload Area -->
        <div class="rounded-md bg-white p-4">
          <div class="flex items-center gap-4">
            <Upload
              :auto="false"
              :show-upload-progress="true"
              accept="*"
              multiple
              theme="file-input"
              @select-files="(files: any[]) => {
                files.forEach((file: any) => {
                  handleUpload(file);
                });
              }"
            >
              <template #trigger>
                <Button theme="primary">
                  <template #icon>
                    <UploadIcon />
                  </template>
                  上传文件
                </Button>
              </template>
            </Upload>
            <div v-if="uploadingFiles > 0" class="text-sm text-gray-500">
              正在上传 {{ uploadingFiles }} 个文件...
            </div>
          </div>
        </div>

        <!-- Search Form -->
        <div class="rounded-md bg-white p-4">
          <Form :data="searchForm" label-width="80px" colon>
            <div class="grid grid-cols-4 gap-x-4">
              <FormItem label="文件名" name="origin_name">
                <Input
                  v-model="searchForm.origin_name"
                  placeholder="请输入文件名"
                  clearable
                />
              </FormItem>
              <FormItem label="存储方式" name="storage_mode">
                <RadioGroup v-model="searchForm.storage_mode">
                  <Radio :value="1">本地</Radio>
                  <Radio :value="2">云存储</Radio>
                </RadioGroup>
              </FormItem>
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

        <!-- Table Area -->
        <div class="flex min-h-0 flex-1 flex-col rounded-md bg-white p-4">
          <div class="mb-3 flex items-center justify-between">
            <Space>
              <Button theme="danger" variant="outline" @click="handleBatchDelete">
                <template #icon><DeleteIcon /></template>
                删除
              </Button>
            </Space>

            <CrudToolbar
              v-model="displayColumns"
              :column-options="columnOptions"
              :is-recycle-bin="false"
              @refresh="fetchTableData"
            />
          </div>

          <div class="min-h-0 flex-1 overflow-hidden">
            <div
              v-if="tableData.length === 0"
              class="flex h-full items-center justify-center text-gray-400"
            >
              <div class="text-center">
                <div class="mb-4 text-6xl">📁</div>
                <div class="text-lg">暂无文件</div>
                <div class="text-sm">点击上方按钮上传文件</div>
              </div>
            </div>

            <div v-else class="text-center text-gray-500">
              文件上传管理功能开发中...
              <br />
              当前版本：基础框架已完成，等待后端 API 对接
            </div>
          </div>
        </div>
      </div>
    </div>
  </Page>
</template>

<style scoped>
.upload-area {
  border: 2px dashed #d1d5db;
  border-radius: 0.5rem;
  padding: 2rem;
  text-align: center;
  transition: all 0.3s;
}

.upload-area:hover {
  border-color: #3b82f6;
  background-color: #f9fafb;
}
</style>
