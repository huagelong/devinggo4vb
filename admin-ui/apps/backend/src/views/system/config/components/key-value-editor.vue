<script lang="ts" setup>
import { DeleteIcon, PlusIcon } from 'tdesign-icons-vue-next';
import { Button, Input, Space } from 'tdesign-vue-next';
import { computed } from 'vue';

interface KeyValueItem {
  key: string;
  value: string;
}

const props = defineProps<{
  modelValue?: KeyValueItem[];
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: KeyValueItem[]): void;
}>();

const items = computed({
  get: () => props.modelValue ?? [{ key: '', value: '' }],
  set: (value: KeyValueItem[]) => emit('update:modelValue', value),
});

function handleAdd() {
  items.value = [...items.value, { key: '', value: '' }];
}

function handleRemove(index: number) {
  const next = [...items.value];
  next.splice(index, 1);
  items.value = next.length > 0 ? next : [{ key: '', value: '' }];
}

function handleUpdate(index: number, field: keyof KeyValueItem, value: string) {
  const next = items.value.map((item, idx) =>
    idx === index ? { ...item, [field]: value } : item,
  );
  items.value = next;
}
</script>

<template>
  <div class="flex flex-col gap-3">
    <div
      v-for="(item, index) in items"
      :key="index"
      class="flex items-center gap-2"
    >
      <Input
        v-model="item.key"
        placeholder="键"
        class="flex-1"
        @change="(val) => handleUpdate(index, 'key', val as string)"
      />
      <Input
        v-model="item.value"
        placeholder="ֵ"
        class="flex-1"
        @change="(val) => handleUpdate(index, 'value', val as string)"
      />
      <Button
        variant="outline"
        shape="square"
        theme="danger"
        @click="handleRemove(index)"
      >
        <template #icon><DeleteIcon /></template>
      </Button>
    </div>
    <Space>
      <Button variant="outline" @click="handleAdd">
        <template #icon><PlusIcon /></template>
        新增键值
      </Button>
    </Space>
  </div>
</template>
