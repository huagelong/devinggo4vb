<script lang="ts" setup>
import type { UploadTreeItem } from './model';

import { computed, onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { message } from '#/adapter/tdesign';
import { logger } from '#/utils/logger';
import CrudToolbar from '#/components/crud/crud-toolbar.vue';

import {
  DeleteIcon,
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
  Radio,
  RadioGroup,
  Space,
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
  fetchTableData,
  handleReset,
  handleSearch,
  selectedRowKeys,
  tableData,
} = useUploadCrud();

function handleTreeChange(value: Array<string | number>) {
  const keys = value.map((item) => String(item));
  selectedTreeKey.value = keys.length > 0 ? keys : ['all'];
  const key = selectedTreeKey.value[0];
  if (key === 'all') {
    searchForm.value.mime_type = '';
  } else {
    searchForm.value.mime_type = key ?? '';
  }
  handleSearch();
}

async function handleUpload(_file: File) {
  uploadingFiles.value++;
  try {
    // TODO: replace with actual upload API call
    await new Promise((resolve) => setTimeout(resolve, 1000));
    MessagePlugin.success($t('common.uploadSuccess'));
  } catch (error) {
    if (import.meta.env.DEV) {
      logger.error(error);
    }
    MessagePlugin.error($t('common.uploadFailed'));
  } finally {
    uploadingFiles.value--;
  }
}

function handleBatchDelete() {
  if (selectedRowKeys.value.length === 0) {
    message.warning($t('common.selectDataFirst'));
    return;
  }
  // TODO: implement batch delete
  message.info($t('common.batchDeleteFeaturePending'));
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
                 files.forEach((f: any) => {
                   handleUpload(f.raw ?? f);
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
