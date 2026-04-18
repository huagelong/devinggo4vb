<script lang="ts" setup>
import { logger } from '#/utils/logger';
import type { ApiColumnFormModel, ApiColumnListItem, ApiColumnType } from '../model';
import type { DictOption } from '#/composables/crud/use-dict-options';
import type { ApiColumnApi } from '#/api/system/api-column';

import { nextTick, reactive, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { MessagePlugin, Select } from 'tdesign-vue-next';

import { useVbenForm } from '#/adapter/form';
import {
  saveApiColumn,
  updateApiColumn,
} from '#/api/system/api-column';
import { useDictOptions } from '#/composables/crud/use-dict-options';

import { createApiColumnFormDefaultValues } from '../schemas';

const emit = defineEmits(['success']);

interface OpenPayload {
  apiId: number;
  type: ApiColumnType;
  data?: ApiColumnListItem;
}

const dataTypeOptions = ref<DictOption[]>([]);
const statusOptions = ref<DictOption[]>([]);

const requiredOptions: DictOption[] = [
  { label: $t('common.no'), value: 1 },
  { label: $t('common.yes'), value: 2 },
];

const typeOptions: DictOption[] = [
  { label: $t('system.api.requestParams'), value: 1 },
  { label: $t('system.api.responseParams'), value: 2 },
];

const fallbackStatusOptions: DictOption[] = [
  { label: $t('common.statusEnabled'), value: 1 },
  { label: $t('common.statusDisabled'), value: 2 },
];

const modalContext = reactive({
  apiId: 0,
  type: 1 as ApiColumnType,
});

const { getDictOptions } = useDictOptions();

function createSelectProps(options: DictOption[], placeholder: string) {
  return {
    options,
    placeholder,
    clearable: true,
  };
}

const [Form, formApi] = useVbenForm({
  showDefaultActions: false,
  commonConfig: {
    labelWidth: 110,
  },
  schema: [
    {
      component: 'Input',
      dependencies: { show: false, triggerFields: [''] },
      fieldName: 'id',
      label: 'ID',
    },
    {
      component: 'Input',
      componentProps: { placeholder: $t('ui.placeholder.input') },
      fieldName: 'name',
      label: $t('system.api.fieldName'),
      rules: 'required',
    },
    {
      component: Select,
      componentProps: createSelectProps(dataTypeOptions.value, $t('ui.placeholder.select')),
      fieldName: 'data_type',
      label: $t('system.api.dataType'),
      rules: 'required',
    },
    {
      component: 'RadioGroup',
      componentProps: { options: typeOptions, disabled: true },
      fieldName: 'type',
      label: $t('system.api.fieldType'),
      rules: 'required',
    },
    {
      component: 'RadioGroup',
      componentProps: { options: statusOptions.value },
      fieldName: 'status',
      label: $t('common.status'),
      rules: 'required',
    },
    {
      component: 'RadioGroup',
      componentProps: { options: requiredOptions },
      fieldName: 'is_required',
      label: $t('system.api.isRequired'),
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: { placeholder: $t('ui.placeholder.input') },
      fieldName: 'default_value',
      label: $t('system.api.defaultValue'),
    },
    {
      component: 'Textarea',
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        autosize: { minRows: 3, maxRows: 6 },
      },
      fieldName: 'description',
      label: $t('system.api.fieldDescription'),
    },
    {
      component: 'Textarea',
      componentProps: { placeholder: $t('ui.placeholder.input') },
      fieldName: 'remark',
      label: $t('common.remark'),
    },
  ],
});

const [Modal, modalApi] = useVbenModal({
  onConfirm: async () => {
    try {
      const { valid } = await formApi.validate();
      if (!valid) return;
      const values = await formApi.getValues<ApiColumnFormModel>();
      modalApi.setState({ confirmLoading: true });
      const payload: ApiColumnApi.SubmitPayload = {
        ...values,
        api_id: modalContext.apiId,
        type: modalContext.type,
        data_type: values.data_type as string | number,
      };
      if (values.id) {
        await updateApiColumn(Number(values.id), payload);
      } else {
        await saveApiColumn(payload);
      }
      MessagePlugin.success(values.id ? $t('common.updateSuccess') : $t('common.createSuccess'));
      emit('success');
      modalApi.close();
    } catch (error) {
      logger.error(error);
    } finally {
      modalApi.setState({ confirmLoading: false });
    }
  },
  class: 'w-[600px]',
});

function updateSchemas() {
  formApi.updateSchema([
    {
      fieldName: 'data_type',
      componentProps: createSelectProps(dataTypeOptions.value, $t('ui.placeholder.select')),
    },
    {
      fieldName: 'status',
      componentProps: { options: statusOptions.value },
    },
  ]);
}

async function fetchFormOptions() {
  try {
    const [dataTypes, statuses] = await Promise.all([
      getDictOptions('api_data_type'),
      getDictOptions('data_status'),
    ]);
    dataTypeOptions.value = dataTypes;
    statusOptions.value =
      statuses && statuses.length > 0 ? statuses : fallbackStatusOptions;
  } catch (error) {
    logger.error(error);
    MessagePlugin.error($t('common.columnOptionsLoadFailed'));
    dataTypeOptions.value = [];
    statusOptions.value = fallbackStatusOptions;
  } finally {
    updateSchemas();
  }
}

async function open(payload: OpenPayload) {
  modalContext.apiId = payload.apiId;
  modalContext.type = payload.type;
  modalApi.setState({
    title: payload.data?.id ? $t('system.api.editColumnTitle') : $t('system.api.createColumnTitle'),
  });
  modalApi.open();
  await fetchFormOptions();
  await formApi.resetForm();
  const defaultValues = createApiColumnFormDefaultValues();
  defaultValues.api_id = payload.apiId;
  defaultValues.type = payload.type;
  formApi.setValues(defaultValues);
  if (payload.data) {
    formApi.setValues({
      ...payload.data,
      api_id: payload.apiId,
      type: payload.type,
    });
  }
  await nextTick();
  await formApi.resetValidate();
}

defineExpose({
  open,
});
</script>

<template>
  <Modal>
    <Form />
  </Modal>
</template>
