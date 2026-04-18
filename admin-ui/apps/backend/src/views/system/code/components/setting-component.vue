<script lang="ts" setup>
import type { FieldConfigRow } from '../model';

import { computed, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';
import { $t } from '@vben/locales';

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
  title: $t('system.code.setting.title'),
});
</script>

<template>
  <Modal>
    <Form :label-width="100" colon>
      <FormItem :label="$t('system.code.field.name')">
        <Input v-model="localRow.column_name" disabled />
      </FormItem>
      <FormItem :label="$t('system.code.field.comment')">
        <Input v-model="localRow.column_comment" />
      </FormItem>
      <FormItem :label="$t('system.code.setting.controlType')">
        <Select v-model="localRow.view_type" :options="viewTypeOptions" />
      </FormItem>

      <!-- 数字类配置 -->
      <template v-if="showNumberConfig">
        <div class="grid grid-cols-2 gap-x-4">
          <FormItem :label="$t('system.code.setting.minValue')">
            <InputNumber v-model="min" />
          </FormItem>
          <FormItem :label="$t('system.code.setting.maxValue')">
            <InputNumber v-model="max" />
          </FormItem>
          <FormItem :label="$t('system.code.setting.step')">
            <InputNumber v-model="step" />
          </FormItem>
          <FormItem :label="$t('system.code.setting.precision')">
            <InputNumber v-model="precision" :min="0" :max="10" />
          </FormItem>
        </div>
      </template>

      <!-- Switch 配置 -->
      <template v-if="showSwitchConfig">
        <div class="grid grid-cols-2 gap-x-4">
          <FormItem :label="$t('system.code.setting.checkedValue')">
            <Input v-model="checkedValue" />
          </FormItem>
          <FormItem :label="$t('system.code.setting.uncheckedValue')">
            <Input v-model="uncheckedValue" />
          </FormItem>
        </div>
      </template>

      <!-- Select 配置 -->
      <template v-if="showSelectConfig">
        <FormItem :label="$t('system.code.setting.multiple')">
          <Switch v-model="isMultiple" />
        </FormItem>
        <FormItem :label="$t('system.code.setting.optionData')">
          <Textarea
            v-model="optionsData"
            :placeholder="$t('system.code.setting.optionDataPlaceholder')"
          />
        </FormItem>
      </template>

      <!-- 日期配置 -->
      <template v-if="showDateConfig">
        <FormItem :label="$t('system.code.setting.pickerType')">
          <Select v-model="dateType" :options="[
            { label: $t('system.code.setting.pickerDate'), value: 'date' },
            { label: $t('system.code.setting.pickerWeek'), value: 'week' },
            { label: $t('system.code.setting.pickerMonth'), value: 'month' },
            { label: $t('system.code.setting.pickerYear'), value: 'year' },
          ]" />
        </FormItem>
        <FormItem :label="$t('system.code.setting.showTime')">
          <Switch v-model="showTime" />
        </FormItem>
        <FormItem :label="$t('system.code.setting.rangePicker')">
          <Switch v-model="isRange" />
        </FormItem>
      </template>

      <!-- 上传配置 -->
      <template v-if="showUploadConfig">
        <FormItem :label="$t('system.code.setting.returnDataType')">
          <Select v-model="localRow.dict_type" :options="[
            { label: 'URL', value: 'url' },
            { label: 'ID', value: 'id' },
          ]" />
        </FormItem>
        <FormItem :label="$t('system.code.setting.multiple')">
          <Switch v-model="isMultiple" />
        </FormItem>
      </template>
    </Form>
  </Modal>
</template>
