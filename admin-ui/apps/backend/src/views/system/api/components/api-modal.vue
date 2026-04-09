<script lang="ts" setup>
import type { ApiFormModel } from '../model';
import type { OptionItem, IdType } from '#/types/common';
import type { DictOption } from '#/composables/crud/use-dict-options';
import type { ApiManageApi } from '#/api/system/api';

import { nextTick, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import { MessagePlugin, Select } from 'tdesign-vue-next';

import { useVbenForm } from '#/adapter/form';
import { saveApi, updateApi } from '#/api/system/api';
import { getApiGroupList } from '#/api/system/api-group';
import { useDictOptions } from '#/composables/crud/use-dict-options';

import { createApiFormDefaultValues } from '../schemas';

type SelectOption = OptionItem<IdType>;

const emit = defineEmits(['success']);

const groupOptions = ref<SelectOption[]>([]);
const requestModeOptions = ref<DictOption[]>([]);
const statusOptions = ref<DictOption[]>([]);

const fallbackRequestModes: DictOption[] = [
  { label: 'GET', value: 'GET' },
  { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' },
  { label: 'DELETE', value: 'DELETE' },
];

const fallbackStatusOptions: DictOption[] = [
  { label: '正常', value: 1 },
  { label: '停用', value: 2 },
];

const authModeOptions: DictOption[] = [
  { label: '简易模式', value: 1 },
  { label: '复杂模式', value: 2 },
];

const { getDictOptions } = useDictOptions();

function createSelectProps(options: DictOption[] | SelectOption[], placeholder: string) {
  return {
    options,
    placeholder,
    clearable: true,
  };
}

const [Form, formApi] = useVbenForm({
  showDefaultActions: false,
  commonConfig: {
    labelWidth: 100,
  },
  schema: [
    {
      component: 'Input',
      dependencies: { show: false, triggerFields: [''] },
      fieldName: 'id',
      label: 'ID',
    },
    {
      component: Select,
      componentProps: createSelectProps(groupOptions.value, '请选择所属组'),
      fieldName: 'group_id',
      label: '所属组',
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: { placeholder: '请输入接口名称' },
      fieldName: 'name',
      label: '接口名称',
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: { placeholder: '请输入接口标识，如 system:app:getAppSecret' },
      fieldName: 'access_name',
      label: '接口标识',
      rules: 'required',
    },
    {
      component: Select,
      componentProps: createSelectProps(requestModeOptions.value, '请选择请求模式'),
      fieldName: 'request_mode',
      label: '请求模式',
      rules: 'required',
    },
    {
      component: 'RadioGroup',
      componentProps: { options: authModeOptions },
      defaultValue: 1,
      fieldName: 'auth_mode',
      label: '认证模式',
      rules: 'required',
    },
    {
      component: 'RadioGroup',
      componentProps: { options: statusOptions.value },
      defaultValue: 1,
      fieldName: 'status',
      label: '״̬',
      rules: 'required',
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
      const values = await formApi.getValues<ApiFormModel>();
      modalApi.setState({ confirmLoading: true });
      const payload: ApiManageApi.SubmitPayload = {
        ...values,
        group_id: values.group_id as IdType,
        request_mode: values.request_mode as number | string,
      };
      if (values.id) {
        await updateApi(Number(values.id), payload);
      } else {
        await saveApi(payload);
      }
      MessagePlugin.success(values.id ? '更新成功' : '新增成功');
      emit('success');
      modalApi.close();
    } catch (error) {
      console.error(error);
    } finally {
      modalApi.setState({ confirmLoading: false });
    }
  },
  class: 'w-[600px]',
});

function updateFormSchemas() {
  formApi.updateSchema([
    {
      fieldName: 'group_id',
      componentProps: createSelectProps(groupOptions.value, '请选择所属组'),
    },
    {
      fieldName: 'request_mode',
      componentProps: createSelectProps(requestModeOptions.value, '请选择请求模式'),
    },
    {
      fieldName: 'status',
      componentProps: { options: statusOptions.value },
    },
  ]);
}

async function fetchFormOptions() {
  try {
    const [groupList, requestModes, statuses] = await Promise.all([
      getApiGroupList().catch(() => []),
      getDictOptions('request_mode'),
      getDictOptions('data_status'),
    ]);
    groupOptions.value = (groupList || []).map((item) => ({
      label: item.name,
      value: item.id,
    }));
    requestModeOptions.value =
      (requestModes && requestModes.length > 0 ? requestModes : fallbackRequestModes).map(
        (item) => ({
          ...item,
          value: item.value ?? item.label,
        }),
      );
    statusOptions.value =
      statuses && statuses.length > 0 ? statuses : fallbackStatusOptions;
  } catch (error) {
    console.error(error);
    MessagePlugin.error('下拉选项加载失败，请稍后重试');
    groupOptions.value = [];
    requestModeOptions.value = fallbackRequestModes;
    statusOptions.value = fallbackStatusOptions;
  } finally {
    updateFormSchemas();
  }
}

async function open(data?: ApiFormModel) {
  modalApi.setState({
    title: data?.id ? '编辑接口' : '新增接口',
  });
  modalApi.open();
  await fetchFormOptions();
  await formApi.resetForm();
  const defaultValues = createApiFormDefaultValues();
  if (!defaultValues.request_mode && requestModeOptions.value.length > 0) {
    defaultValues.request_mode = requestModeOptions.value[0]?.value;
  }
  formApi.setValues(defaultValues);
  if (data) {
    formApi.setValues(data);
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
