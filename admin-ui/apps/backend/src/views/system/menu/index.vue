<script lang="ts" setup>
import type { MenuApi } from '#/api/system/menu';
import type { DictOption } from '#/composables/crud/use-dict-options';

import { computed, onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { message } from '#/adapter/tdesign';
import { logger } from '#/utils/logger';
import {
  changeMenuStatus,
  deleteMenu,
  realDeleteMenu,
  recoveryMenu,
  updateMenuNumber,
} from '#/api/system/menu';
import CrudToolbar from '#/components/crud/crud-toolbar.vue';
import { useDictOptions } from '#/composables/crud/use-dict-options';

import {
  AddIcon,
  DeleteIcon,
  EditIcon,
  PlusIcon,
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
  Tag,
} from 'tdesign-vue-next';

import MenuModal from './components/menu-modal.vue';
import type { MenuTableColumn } from './model';
import {
  createMenuColumnOptions,
  createMenuTableColumns,
  menuTypeTagMap,
} from './schemas';
import { useMenuPage } from './use-menu-page';

defineOptions({ name: 'SystemMenu' });

type MenuModalInstance = {
  open: (data?: Partial<MenuApi.SubmitPayload>) => void;
};

const menuModalRef = ref<MenuModalInstance>();
const statusOptions = ref<DictOption[]>([]);
const columns: MenuTableColumn[] = createMenuTableColumns();
const columnOptions = createMenuColumnOptions(columns);
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
  handleReset,
  handleSearch,
  handleSelectChange,
  isRecycleBin,
  loading,
  searchForm,
  selectedRowKeys,
  tableData,
  toggleRecycleBin,
} = useMenuPage();

const { getDictOptions } = useDictOptions();

function toIds(keys: Array<number | string>) {
  return keys.map((key) => Number(key));
}

async function fetchStatusOptions() {
  statusOptions.value =
    (await getDictOptions('data_status')) || [
      { label: '正常', value: 1 },
      { label: '停用', value: 2 },
    ];
}

function handleAdd(parentId = 0) {
  menuModalRef.value?.open({ parent_id: parentId });
}

function handleEdit(row: MenuApi.TreeItem) {
  menuModalRef.value?.open({
    ...row,
    is_hidden: Number(row.is_hidden ?? 2),
    parent_id: Number(row.parent_id ?? 0),
    restful: row.restful ?? '2',
    sort: Number(row.sort ?? 1),
    status: Number(row.status ?? 1),
  });
}

async function handleDelete(row: MenuApi.TreeItem) {
  try {
    await (isRecycleBin.value ? realDeleteMenu([row.id]) : deleteMenu([row.id]));
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
    await (isRecycleBin.value ? realDeleteMenu(ids) : deleteMenu(ids));
    message.success($t('common.operationSuccess'));
    clearSelectedRowKeys();
    await fetchTableData();
  } catch (error) {
    logger.error(error);
    message.error($t('common.batchDeleteFailed'));
  }
}

async function handleRecovery(row: MenuApi.TreeItem) {
  try {
    await recoveryMenu([row.id]);
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
    await recoveryMenu(ids);
    message.success($t('common.recoverySuccess'));
    clearSelectedRowKeys();
    await fetchTableData();
  } catch (error) {
    logger.error(error);
    message.error($t('common.batchRecoveryFailed'));
  }
}

async function handleStatusChange(row: MenuApi.TreeItem, checked: boolean) {
  try {
    await changeMenuStatus({ id: row.id, status: checked ? 1 : 2 });
    message.success($t('common.statusUpdateSuccess'));
    await fetchTableData();
  } catch (error) {
    logger.error(error);
    message.error($t('common.statusUpdateFailed'));
  }
}

function handleStatusSwitchChange(row: MenuApi.TreeItem, value: unknown) {
  void handleStatusChange(row, Boolean(value));
}

async function handleSortChange(value: number | string, row: MenuApi.TreeItem) {
  const numberValue = Number(value);
  if (Number.isNaN(numberValue)) return;
  try {
    await updateMenuNumber({
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

onMounted(() => {
  void fetchStatusOptions();
  void fetchTableData();
});
</script>

<template>
  <Page auto-content-height>
    <div class="flex h-full flex-col gap-3">
      <div class="rounded-md bg-white p-4">
        <Form :data="searchForm" label-width="90px" colon>
          <div class="grid grid-cols-4 gap-x-4">
            <FormItem label="菜单名称" name="name">
              <Input
                v-model="searchForm.name"
                placeholder="请输入菜单名称"
                clearable
              />
            </FormItem>
            <FormItem label="菜单标识" name="code">
              <Input
                v-model="searchForm.code"
                placeholder="请输入菜单标识"
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
              <Button theme="primary" @click="handleAdd()">
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
          :selected-row-keys="selectedRowKeys"
          :tree="{ childrenKey: 'children', defaultExpandAll: true }"
          row-key="id"
          hover
          stripe
          table-layout="fixed"
          @select-change="handleSelectChange"
        >
          <template #type="{ row }">
            <Tag
              size="small"
              :theme="menuTypeTagMap[row.type]?.theme ?? 'default'"
            >
              {{ menuTypeTagMap[row.type]?.label ?? row.type }}
            </Tag>
          </template>

          <template #sort="{ row }">
            <InputNumber
              :value="row.sort"
              :min="0"
              :max="1000"
              size="small"
              :disabled="isRecycleBin"
              @change="(value: number | string) => handleSortChange(value, row)"
            />
          </template>

          <template #status="{ row }">
            <Switch
              :disabled="isRecycleBin"
              :value="Number(row.status) === 1"
              @change="(value: unknown) => handleStatusSwitchChange(row, value)"
            />
          </template>

          <template #action="{ row }">
            <div class="flex items-center justify-center gap-1">
              <template v-if="!isRecycleBin">
                <Button
                  v-if="row.type === 'M'"
                  size="small"
                  theme="default"
                  variant="outline"
                  @click="handleAdd(Number(row.id))"
                >
                  <template #icon><PlusIcon /></template>
                  新增子级
                </Button>
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
                  content="确认删除该菜单吗？"
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
                  content="确认恢复该菜单吗？"
                  @confirm="handleRecovery(row)"
                >
                  <Button size="small" theme="primary" variant="outline">
                    恢复
                  </Button>
                </Popconfirm>
                <Popconfirm
                  content="确认彻底删除该菜单吗？"
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

    <MenuModal ref="menuModalRef" @success="handleSuccess" />
  </Page>
</template>
