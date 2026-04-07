<script lang="ts" setup>
import { ref } from 'vue';

import { UploadIcon } from 'tdesign-icons-vue-next';
import { Button, Input, MessagePlugin, Space } from 'tdesign-vue-next';

import { uploadImageFileApi } from '#/api/system/upload';

const props = defineProps<{
  modelValue?: string;
  placeholder?: string;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void;
}>();

const uploading = ref(false);
const fileInputRef = ref<HTMLInputElement>();

function handleInput(value: string) {
  emit('update:modelValue', value);
}

function triggerUpload() {
  fileInputRef.value?.click();
}

async function handleFileChange(event: Event) {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  if (!file) return;
  uploading.value = true;
  try {
    const response = (await uploadImageFileApi(file)) as { url?: string };
    if (response?.url) {
      emit('update:modelValue', response.url);
      MessagePlugin.success('上传成功');
    } else {
      MessagePlugin.error('上传失败');
    }
  } catch (error) {
    console.error(error);
    MessagePlugin.error('上传失败');
  } finally {
    uploading.value = false;
    input.value = '';
  }
}
</script>

<template>
  <div class="flex flex-col gap-2">
    <Space>
      <Input
        :model-value="modelValue"
        :placeholder="placeholder ?? '请输入或上传链接'"
        class="w-80"
        @change="(value) => handleInput(value as string)"
      />
      <Button :loading="uploading" variant="outline" @click="triggerUpload">
        <template #icon><UploadIcon /></template>
        上传
      </Button>
    </Space>
    <div v-if="modelValue" class="rounded-md border border-gray-100 p-2">
      <img
        :src="modelValue"
        alt="config upload preview"
        class="h-24 max-w-full rounded-md object-contain"
      />
    </div>
    <input
      ref="fileInputRef"
      type="file"
      accept="image/*"
      class="hidden"
      @change="handleFileChange"
    />
  </div>
</template>
