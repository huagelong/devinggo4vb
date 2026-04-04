<script lang="ts" setup>
import type { PreviewCodeRow } from '../model';

import { ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';
import { message } from '#/adapter/tdesign';

import { previewCode } from '#/api/system/generate';
import { CodeIcon } from 'tdesign-icons-vue-next';

import { Button, TabPanel, Tabs } from 'tdesign-vue-next';

const emit = defineEmits<{
  success: [];
}>();

const loading = ref(false);
const previewList = ref<PreviewCodeRow[]>([]);
const activeTab = ref('0');

async function open(id: number) {
  loading.value = true;
  previewList.value = [];
  try {
    const response = await previewCode(id);
    previewList.value = response.data || [];
    if (previewList.value.length > 0) {
      activeTab.value = '0';
    }
    modalApi.setState({ title: '代码预览' });
    modalApi.open();
  } catch (error) {
    console.error(error);
    message.error('获取预览失败');
  } finally {
    loading.value = false;
  }
}

function handleCopy(code: string) {
  navigator.clipboard.writeText(code).then(() => {
    message.success('复制成功');
  }).catch(() => {
    message.error('复制失败');
  });
}

const [Modal, modalApi] = useVbenModal({
  footer: false,
  class: 'w-[1000px]',
});

defineExpose({ open });
</script>

<template>
  <Modal>
    <div class="flex flex-col gap-3">
      <div v-if="loading" class="flex items-center justify-center py-8">
        加载中...
      </div>

      <div v-else-if="previewList.length === 0" class="py-8 text-center text-gray-500">
        暂无预览数据
      </div>

      <template v-else>
        <Tabs v-model:value="activeTab">
          <TabPanel
            v-for="(item, index) in previewList"
            :key="item.name"
            :value="String(index)"
            :label="item.tab_name"
          >
            <div class="flex justify-end mb-2">
              <Button
                size="small"
                @click="handleCopy(item.code)"
              >
                <template #icon><CodeIcon /></template>
                复制代码
              </Button>
            </div>
            <pre
              class="max-h-[500px] overflow-auto rounded bg-gray-900 p-4 text-sm text-gray-100"
            >{{ item.code }}</pre>
          </TabPanel>
        </Tabs>
      </template>
    </div>
  </Modal>
</template>
