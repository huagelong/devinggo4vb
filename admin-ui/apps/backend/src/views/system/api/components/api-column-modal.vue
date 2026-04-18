<script lang="ts" setup>
import { logger } from '#/utils/logger';
import type { ApiColumnFormModel, ApiColumnListItem, ApiColumnType } from '../model';
import type { DictOption } from '#/composables/crud/use-dict-options';
import type { ApiColumnApi } from '#/api/system/api-column';

import { nextTick, reactive, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';

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
  { label: '否', value: 1 },
  { label: '是', value: 2 },
];

const typeOptions: DictOption[] = [
  { label: '请求参数', value: 1 },
  { label: '响应参数', value: 2 },
];

const fallbackStatusOptions: DictOption[] = [
  { label: '正常', value: 1 },
  { label: '停用', value: 2 },
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
      componentProps: { placeholder: '请输入字段名称' },
      fieldName: 'name',
      label: '字段名称',
      rules: 'required',
    },
    {
      component: Select,
      componentProps: createSelectProps(dataTypeOptions.value, '请选择数据类型'),
      fieldName: 'data_type',
      label: '数据类型',
      rules: 'required',
    },
    {
      component: 'RadioGroup',
      componentProps: { options: typeOptions, disabled: true },
      fieldName: 'type',
      label: '字段类型',
      rules: 'required',
    },
    {
      component: 'RadioGroup',
      componentProps: { options: statusOptions.value },
      fieldName: 'status',
      label: '状态',
      rules: 'required',
    },
    {
      component: 'RadioGroup',
      componentProps: { options: requiredOptions },
      fieldName: 'is_required',
      label: '是否必填',
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: { placeholder: '请输入默认值' },
      fieldName: 'default_value',
      label: '默认值',
    },
    {
      component: 'Textarea',
      componentProps: {
        placeholder: '请输入字段说明',
        autosize: { minRows: 3, maxRows: 6 },
      },
      fieldName: 'description',
      label: '字段说明',
    },
    {
      component: 'Textarea',
      componentProps: { placeholder: '请输入备注' },
      fieldName: 'remark',
      label: '备注',
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
      MessagePlugin.success(values.id ? '更新成功' : '新增成功');
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
      componentProps: createSelectProps(dataTypeOptions.value, '请选择数据类型'),
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
    MessagePlugin.error('字段选项加载失败，请稍后重试');
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
    title: payload.data?.id ? '编辑参数' : '新增参数',
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
