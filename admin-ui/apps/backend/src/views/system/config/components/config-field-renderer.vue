<script lang="ts" setup>
import type { ConfigFieldMeta } from '../model';

import { computed } from 'vue';

import { $t } from '@vben/locales';

import {
  Checkbox,
  CheckboxGroup,
  Input,
  Radio,
  RadioGroup,
  Select,
  Switch,
  Textarea,
} from 'tdesign-vue-next';

import ConfigUploadInput from './config-upload-input.vue';
import KeyValueEditor from './key-value-editor.vue';

const props = defineProps<{
  field: ConfigFieldMeta;
  modelValue?: any;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: any): void;
}>();

const innerValue = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value),
});

const selectOptions = computed(
  () =>
    props.field.config_select_data?.map((item) => ({
      label: item.label,
      value: item.value,
    })) ?? [],
);

const switchValues = computed(() => {
  const defaults = props.field.switchValues ?? { checked: 1, unchecked: 0 };
  return [defaults.checked, defaults.unchecked];
});
</script>

<template>
  <div>
    <Input
      v-if="field.input_type === 'input'"
      v-model="innerValue"
      allow-clear
      :placeholder="$t('common.enterContent')"
    />
    <Textarea
      v-else-if="field.input_type === 'textarea' || field.input_type === 'editor'"
      v-model="innerValue"
      :autosize="{ minRows: 3, maxRows: 6 }"
      :placeholder="$t('common.enterContent')"
    />
    <Select
      v-else-if="field.input_type === 'select'"
      v-model="innerValue"
      :options="selectOptions"
      :placeholder="$t('ui.placeholder.select')"
      clearable
      class="w-full"
    />
    <RadioGroup
      v-else-if="field.input_type === 'radio'"
      v-model="innerValue"
      class="flex flex-wrap gap-3"
    >
      <Radio
        v-for="option in selectOptions"
        :key="option.value"
        :value="option.value"
      >
        {{ option.label }}
      </Radio>
    </RadioGroup>
    <CheckboxGroup
      v-else-if="field.input_type === 'checkbox'"
      v-model="innerValue"
      class="flex flex-wrap gap-3"
    >
      <Checkbox
        v-for="option in selectOptions"
        :key="option.value"
        :value="option.value"
      >
        {{ option.label }}
      </Checkbox>
    </CheckboxGroup>
    <Switch
      v-else-if="field.input_type === 'switch'"
      v-model="innerValue"
      :custom-value="switchValues"
    />
    <ConfigUploadInput
      v-else-if="field.input_type === 'upload'"
      v-model="innerValue"
    />
    <KeyValueEditor
      v-else-if="field.input_type === 'key-value'"
      v-model="innerValue"
    />
    <Input
      v-else
      v-model="innerValue"
      allow-clear
      :placeholder="$t('common.enterContent')"
    />
  </div>
</template>
