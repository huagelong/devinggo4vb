<script lang="ts" setup>
import { computed, onMounted, ref } from 'vue';

import { useAccess } from '@vben/access';
import { Page } from '@vben/common-ui';

import { message } from '#/adapter/tdesign';
import CrudToolbar from '#/components/crud/crud-toolbar.vue';

import { InfoCircleFilledIcon, SearchIcon } from 'tdesign-icons-vue-next';
import { Button, Form, FormItem, Input, Popup, Space, Table } from 'tdesign-vue-next';

import type { DataMaintainTableColumn } from './model';
import {
  createDataMaintainColumnOptions,
  createDataMaintainTableColumns,
} from './schemas';
import { useDataMaintainCrud } from './use-data-maintain-crud';

defineOptions({ name: 'SystemDataMaintain' });

const { hasAccessByCodes } = useAccess();
const canView = computed(() =>
  hasAccessByCodes(['system:dataMaintain:index', 'system:dataMaintain']),
);

const columns: DataMaintainTableColumn[] = createDataMaintainTableColumns();
const columnOptions = createDataMaintainColumnOptions(columns);
const allColumnKeys = columnOptions.map((item) => item.value);
const visibleColumns = ref<string[]>([...allColumnKeys]);

const displayColumns = computed({
  get: () => [...visibleColumns.value],
  set: (value: string[]) => {
    visibleColumns.value = value;
  },
});

const {
  fetchTableData,
  handlePageChange,
  handleReset,
  handleSearch,
  loading,
  pagination,
  searchForm,
  tableData,
} = useDataMaintainCrud();

function handleUnimplementedAction(actionName: string) {
  message.info(`${actionName}接口暂未在当前后端开放`);
}

onMounted(() => {
  if (!canView.value) {
    return;
  }
  void fetchTableData();
});
</script>

<template>
  <Page auto-content-height>
    <div class="flex h-full flex-col gap-3">
      <div class="rounded-md bg-white p-4">
        <Form :data="searchForm" label-width="90px" colon>
          <div class="grid grid-cols-4 gap-x-4">
            <FormItem label="数据库组" name="group_name">
              <Input
                v-model="searchForm.group_name"
                placeholder="默认 default"
                clearable
              />
            </FormItem>
            <FormItem label="表名" name="name">
              <Input v-model="searchForm.name" placeholder="请输入表名" clearable />
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
            <Button variant="outline" @click="handleUnimplementedAction('查看字段')">
              查看字段
            </Button>
            <Button variant="outline" @click="handleUnimplementedAction('优化表')">
              优化表
            </Button>
            <Button variant="outline" @click="handleUnimplementedAction('清理碎片')">
              清理碎片
            </Button>
            <Popup placement="bottom" trigger="hover" content="首版仅接入列表能力，扩展动作待后端接口开放后补齐。">
              <InfoCircleFilledIcon class="cursor-help text-gray-400" />
            </Popup>
          </Space>

          <CrudToolbar
            v-model="displayColumns"
            :column-options="columnOptions"
            :is-recycle-bin="false"
            @refresh="fetchTableData"
          />
        </div>

        <div v-if="!canView" class="rounded-md border border-dashed border-gray-300 p-6 text-center text-gray-500">
          无权限访问数据维护列表（需要 `system:dataMaintain:index`）。
        </div>

        <div v-else class="min-h-0 flex-1">
          <Table
            row-key="name"
            hover
            stripe
            :columns="columns"
            :column-controller-visible="false"
            :display-columns="displayColumns"
            :data="tableData"
            :loading="loading"
            :pagination="pagination"
            @page-change="handlePageChange"
          >
            <template #rows="{ row }">
              {{ row.rows ?? '-' }}
            </template>
            <template #comment="{ row }">
              <span :title="row.comment || '-'">{{ row.comment || '-' }}</span>
            </template>
            <template #create_time="{ row }">
              {{ row.create_time || '-' }}
            </template>
          </Table>
        </div>
      </div>
    </div>
  </Page>
</template>
