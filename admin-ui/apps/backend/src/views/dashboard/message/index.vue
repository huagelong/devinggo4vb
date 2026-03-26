<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { useWindowSize } from '@vueuse/core';
import {
  Button,
  DateRangePicker,
  Dialog,
  DialogPlugin,
  Input,
  Link,
  MessagePlugin,
  RadioButton,
  RadioGroup,
  Table,
} from 'tdesign-vue-next';

import {
  deleteQueueMessageApi,
  getDataDictListApi,
  getQueueMessageReceiveListApi,
  updateQueueMessageReadStatusApi,
} from '#/api/core/message';

const { height } = useWindowSize();

const currentType = ref('all'); // all 全部
const dictOptions = ref<any[]>([]);

const searchForm = reactive({
  title: '',
  created_at: [] as string[],
  read_status: 'all', // all 全部 1未读 2已读
});

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
});

const tableData = ref([]);
const tableLoading = ref(false);
const selectedRowKeys = ref<any[]>([]);

const columns: any[] = [
  { colKey: 'row-select', type: 'multiple' as const, width: 50 },
  { title: '发送人', colKey: 'send_user.nickname', width: 120 },
  { title: '消息标题', colKey: 'title', ellipsis: true },
  { title: '消息类型', colKey: 'content_type', width: 100 },
  { title: '发送时间', colKey: 'created_at', width: 180 },
  { title: '操作', colKey: 'action', width: 150, align: 'center' },
];

const loadDict = async () => {
  try {
    const res: any = await getDataDictListApi({ code: 'queue_msg_type' });
    dictOptions.value = Array.isArray(res)
      ? res
      : res?.items || res?.data || [];
  } catch (error) {
    console.error(error);
  }
};

const getTypeName = (val: string) => {
  const item = dictOptions.value.find((i) => i.key === val);
  return item ? item.title || item.key : val;
};

const fetchData = async () => {
  tableLoading.value = true;
  try {
    const params: any = {
      page: pagination.current,
      pageSize: pagination.pageSize,
      title: searchForm.title,
      read_status: searchForm.read_status,
      content_type: currentType.value,
    };
    if (
      searchForm.created_at &&
      searchForm.created_at.length === 2 &&
      searchForm.created_at[0]
    ) {
      params.created_at = searchForm.created_at;
    }
    const res: any = await getQueueMessageReceiveListApi(params);
    tableData.value = res?.items || res?.data?.items || [];
    pagination.total = res?.pageInfo?.total || res?.data?.pageInfo?.total || 0;
  } finally {
    tableLoading.value = false;
  }
};

const onPageChange = (pageInfo: any) => {
  pagination.current = pageInfo.current;
  pagination.pageSize = pageInfo.pageSize;
  fetchData();
};

const onSearch = () => {
  pagination.current = 1;
  fetchData();
};

const onReset = () => {
  searchForm.title = '';
  searchForm.created_at = [];
  onSearch();
};

const onSelectChange = (val: any[]) => {
  selectedRowKeys.value = val;
};

const handleChangeType = (val: string) => {
  currentType.value = val;
  onSearch();
};

const handleChangeStatus = (val: string) => {
  searchForm.read_status = val;
  onSearch();
};

const handleBatchDelete = () => {
  if (selectedRowKeys.value.length === 0) {
    MessagePlugin.warning('请先选择要删除的数据');
    return;
  }
  const dialog = DialogPlugin.confirm({
    header: '确认删除',
    body: '确定要删除所选数据吗？',
    onConfirm: async () => {
      await deleteQueueMessageApi({ ids: selectedRowKeys.value });
      MessagePlugin.success('删除成功');
      selectedRowKeys.value = [];
      fetchData();
      dialog.hide();
    },
  });
};

const handleDelete = (row: any) => {
  const dialog = DialogPlugin.confirm({
    header: '确认删除',
    body: '确定要删除该数据吗？',
    onConfirm: async () => {
      await deleteQueueMessageApi({ ids: [row.id] });
      MessagePlugin.success('删除成功');
      fetchData();
      dialog.hide();
    },
  });
};

const detailVisible = ref(false);
const detailData = ref<any>({});
const handleDetail = async (row: any) => {
  detailData.value = row;
  detailVisible.value = true;
  // Mark as read
  try {
    await updateQueueMessageReadStatusApi({ ids: [row.id] });
    fetchData(); // refresh list
  } catch (error) {
    console.error(error);
  }
};

onMounted(() => {
  loadDict().then(() => fetchData());
});
</script>

<template>
  <div class="flex h-full w-full p-4 gap-4 bg-gray-50">
    <!-- Left Menu -->
    <div
      class="w-48 bg-white h-full shrink-0 flex flex-col pt-4 drop-shadow-sm rounded"
    >
      <div
        class="menu-item"
        :class="{ active: currentType === 'all' }"
        @click="handleChangeType('all')"
      >
        <span class="icon i-lucide:mail"></span>
        全部
      </div>
      <div
        v-for="item in dictOptions"
        :key="item.key"
        class="menu-item"
        :class="{ active: currentType === item.key }"
        @click="handleChangeType(item.key)"
      >
        <span
          class="icon"
          :class="
            item.key === 'notice' ? 'i-lucide:bell' : 'i-lucide:file-text'
          "
        ></span>
        {{ item.title || item.key }}
      </div>
    </div>

    <!-- Right Content -->
    <div
      class="flex-1 bg-white h-full flex flex-col min-w-0 p-4 rounded drop-shadow-sm relative"
    >
      <!-- Search Form -->
      <div class="flex gap-4 items-center mb-4 flex-wrap text-sm">
        <div class="flex items-center gap-2">
          <span class="text-gray-600 whitespace-nowrap shrink-0">消息标题</span>
          <Input
            v-model="searchForm.title"
            placeholder="请输入消息标题"
            class="w-48"
            clearable
          />
        </div>
        <div class="flex items-center gap-2">
          <span class="text-gray-600 whitespace-nowrap shrink-0">发送时间</span>
          <DateRangePicker
            v-model="searchForm.created_at"
            allow-input
            clearable
            class="w-64"
          />
        </div>
        <Button theme="primary" @click="onSearch">
          <template #icon>
            <SearchIcon
              fill-color="transparent"
              stroke-color="currentColor"
              :stroke-width="2"
            />
          </template>
          搜索
        </Button>
        <Button theme="default" @click="onReset">
          <template #icon>
            <RefreshIcon
              fill-color="transparent"
              stroke-color="currentColor"
              :stroke-width="2"
            />
          </template>
          重置
        </Button>
      </div>

      <!-- Actions -->
      <div class="flex justify-between items-center mb-4">
        <div class="flex gap-3 items-center">
          <Button theme="danger" @click="handleBatchDelete">
            <template #icon>
              <DeleteIcon
                fill-color="transparent"
                stroke-color="currentColor"
                :stroke-width="2"
              />
            </template>
            删除
          </Button>
          <RadioGroup
            v-model="searchForm.read_status"
            variant="outline"
            @change="handleChangeStatus as any"
          >
            <RadioButton value="all">全部</RadioButton>
            <RadioButton value="1">未读</RadioButton>
            <RadioButton value="2">已读</RadioButton>
          </RadioGroup>
        </div>
        <div class="flex gap-2">
          <Button
            theme="default"
            variant="outline"
            shape="square"
            @click="onSearch"
          >
            <template #icon>
              <RefreshIcon
                fill-color="transparent"
                stroke-color="currentColor"
                :stroke-width="2"
              />
            </template>
          </Button>
        </div>
      </div>

      <!-- Table -->
      <div class="flex-1 overflow-hidden">
        <Table
          row-key="id"
          :data="tableData"
          :columns="columns"
          :loading="tableLoading"
          :pagination="pagination"
          :selected-row-keys="selectedRowKeys"
          @page-change="onPageChange"
          @select-change="onSelectChange"
          :max-height="height - 300"
          height="100%"
          table-layout="fixed"
        >
          <template #content_type="{ row }">
            {{ getTypeName(row.content_type) }}
          </template>
          <template #action="{ row }">
            <div class="flex gap-4 items-center justify-center">
              <Link theme="primary" hover="color" @click="handleDetail(row)">
                <span class="i-lucide:eye mr-1"></span> 详细
              </Link>
              <Link theme="danger" hover="color" @click="handleDelete(row)">
                <span class="i-lucide:trash mr-1"></span> 删除
              </Link>
            </div>
          </template>
        </Table>
      </div>
    </div>

    <!-- Detail Dialog -->
    <Dialog
      v-model:visible="detailVisible"
      header="消息详情"
      :footer="false"
      width="800px"
      placement="center"
      destroy-on-close
    >
      <div class="flex flex-col gap-4 py-4 px-2">
        <h2 class="text-2xl font-bold">{{ detailData.title }}</h2>
        <div class="flex justify-between text-gray-500 text-sm border-b pb-2">
          <span>{{ getTypeName(detailData.content_type) }}</span>
          <span>创建时间: {{ detailData.created_at }}</span>
        </div>
        <div
          class="bg-gray-50 p-4 rounded min-h-[200px]"
          v-html="detailData.content"
        ></div>
      </div>
    </Dialog>
  </div>
</template>

<style scoped>
.menu-item {
  display: flex;
  gap: 8px;
  align-items: center;
  padding: 12px 24px;
  margin: 4px 8px;
  color: #4b5563;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.3s;
}

.menu-item:hover {
  background-color: #f3f4f6;
}

.menu-item.active {
  font-weight: 500;
  color: #165dff;
  background-color: #f3f4f6;
}
</style>
