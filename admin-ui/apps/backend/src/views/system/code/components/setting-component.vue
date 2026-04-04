<script lang="ts" setup>
import type { FieldConfigRow } from '../model';

import { computed, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import {
  Form,
  FormItem,
  Input,
  InputNumber,
  Select,
  Switch,
  Textarea,
} from 'tdesign-vue-next';

import {
  viewTypeOptions,
} from '../schemas';

const emit = defineEmits<{
  'update:modelValue': [value: FieldConfigRow];
}>();

const props = defineProps<{
  modelValue: FieldConfigRow;
}>();

// 本地编辑状态
const localRow = ref<FieldConfigRow>({ ...props.modelValue });

// 监听 props 变化
const viewType = computed(() => localRow.value.view_type || 'text');

// 根据 viewType 显示不同配置
const showNumberConfig = computed(() =>
  ['inputNumber', 'slider'].includes(viewType.value)
);
const showSwitchConfig = computed(() => viewType.value === 'switch');
const showSelectConfig = computed(() =>
  ['select', 'checkbox', 'radio', 'transfer'].includes(viewType.value)
);
const showDateConfig = computed(() =>
  ['date', 'time'].includes(viewType.value)
);
const showUploadConfig = computed(() =>
  ['upload', 'selectResource'].includes(viewType.value)
);

// 数字配置
const min = ref(0);
const max = ref(100);
const step = ref(1);
const precision = ref(0);

// Switch 配置
const checkedValue = ref('true');
const uncheckedValue = ref('false');

// Select 配置
const isMultiple = ref(false);
const optionsData = ref('');

// 日期配置
const dateType = ref('date');
const showTime = ref(false);
const isRange = ref(false);

function handleConfirm() {
  emit('update:modelValue', { ...localRow.value });
  modalApi.close();
}

const [Modal, modalApi] = useVbenModal({
  onConfirm: handleConfirm,
  class: 'w-[600px]',
  title: '字段组件配置',
});
</script>

<template>
  <Modal>
    <Form :label-width="100" colon>
      <FormItem label="字段名称">
        <Input v-model="localRow.column_name" disabled />
      </FormItem>
      <FormItem label="字段描述">
        <Input v-model="localRow.column_comment" />
      </FormItem>
      <FormItem label="控件类型">
        <Select v-model="localRow.view_type" :options="viewTypeOptions" />
      </FormItem>

      <!-- 数字类配置 -->
      <template v-if="showNumberConfig">
        <div class="grid grid-cols-2 gap-x-4">
          <FormItem label="最小值">
            <InputNumber v-model="min" />
          </FormItem>
          <FormItem label="最大值">
            <InputNumber v-model="max" />
          </FormItem>
          <FormItem label="步长">
            <InputNumber v-model="step" />
          </FormItem>
          <FormItem label="精度">
            <InputNumber v-model="precision" :min="0" :max="10" />
          </FormItem>
        </div>
      </template>

      <!-- Switch 配置 -->
      <template v-if="showSwitchConfig">
        <div class="grid grid-cols-2 gap-x-4">
          <FormItem label="选中值">
            <Input v-model="checkedValue" />
          </FormItem>
          <FormItem label="未选中值">
            <Input v-model="uncheckedValue" />
          </FormItem>
        </div>
      </template>

      <!-- Select 配置 -->
      <template v-if="showSelectConfig">
        <FormItem label="多选">
          <Switch v-model="isMultiple" />
        </FormItem>
        <FormItem label="选项数据">
          <Textarea
            v-model="optionsData"
            placeholder="请输入选项数据，JSON格式，如：[{label: '是', value: 1}]"
          />
        </FormItem>
      </template>

      <!-- 日期配置 -->
      <template v-if="showDateConfig">
        <FormItem label="选择器类型">
          <Select v-model="dateType" :options="[
            { label: '日期', value: 'date' },
            { label: '周', value: 'week' },
            { label: '月', value: 'month' },
            { label: '年', value: 'year' },
          ]" />
        </FormItem>
        <FormItem label="显示时间">
          <Switch v-model="showTime" />
        </FormItem>
        <FormItem label="范围选择">
          <Switch v-model="isRange" />
        </FormItem>
      </template>

      <!-- 上传配置 -->
      <template v-if="showUploadConfig">
        <FormItem label="返回数据类型">
          <Select v-model="localRow.dict_type" :options="[
            { label: 'URL', value: 'url' },
            { label: 'ID', value: 'id' },
          ]" />
        </FormItem>
        <FormItem label="多选">
          <Switch v-model="isMultiple" />
        </FormItem>
      </template>
    </Form>
  </Modal>
</template>
