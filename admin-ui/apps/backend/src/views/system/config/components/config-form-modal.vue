<script lang="ts" setup>
import type { ConfigApi } from '#/api/system/config';
import type { DictOption } from '#/composables/crud/use-dict-options';

import { nextTick, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import { MessagePlugin } from 'tdesign-vue-next';

import { useVbenForm } from '#/adapter/form';
import { getConfigGroupList, saveConfig, updateConfig } from '#/api/system/config';
import { useDictOptions } from '#/composables/crud/use-dict-options';

import { createConfigFormDefaultValues } from '../schemas';
import { inputComponentOptions } from '../model';

const emit = defineEmits(['success']);

const groupOptions = ref<DictOption<number>[]>([]);
const statusOptions = ref<DictOption[]>([]);

const { getDictOptions } = useDictOptions();

const [Form, formApi] = useVbenForm({
  showDefaultActions: false,
  commonConfig: { labelWidth: 100 },
  schema: [
    {
      component: 'Input',
      dependencies: { show: false, triggerFields: [''] },
      fieldName: 'id',
      label: 'ID',
    },
    {
      component: 'Select',
      componentProps: {
        options: groupOptions.value,
        placeholder: '请选择分组',
      },
      fieldName: 'group_id',
      label: '所属分组',
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: { placeholder: '请输入配置名称' },
      fieldName: 'name',
      label: '配置名称',
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: { placeholder: '请输入配置标识' },
      fieldName: 'key',
      label: '配置标识',
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: { placeholder: '请输入配置值' },
      fieldName: 'value',
      label: '配置值',
    },
    {
      component: 'InputNumber',
      componentProps: { min: 0, max: 999 },
      defaultValue: 0,
      fieldName: 'sort',
      label: '排序',
    },
    {
      component: 'Select',
      componentProps: {
        options: inputComponentOptions,
        placeholder: '请选择组件',
      },
      defaultValue: 'input',
      fieldName: 'input_type',
      label: '输入组件',
      rules: 'required',
    },
    {
      component: 'Select',
      componentProps: {
        options: statusOptions.value,
        placeholder: '请选择状态',
      },
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
    {
      component: 'Textarea',
      componentProps: {
        autosize: { minRows: 3, maxRows: 6 },
        placeholder: '配置下拉/单选/多选数据，JSON 数组格式',
      },
      dependencies: {
        show: (values) =>
          ['select', 'radio', 'checkbox'].includes(values?.input_type),
        triggerFields: ['input_type'],
      },
      fieldName: 'config_select_data',
      label: '可选项数据',
    },
  ],
});

const [Modal, modalApi] = useVbenModal({
  onConfirm: async () => {
    try {
      const { valid } = await formApi.validate();
      if (!valid) return;

      const values = await formApi.getValues<ConfigApi.ConfigSubmitPayload & {
        config_select_data?: string;
      }>();
      const payload: ConfigApi.ConfigSubmitPayload = {
        ...values,
        group_id: Number(values.group_id),
        sort: Number(values.sort ?? 0),
        status: Number(values.status ?? 1),
      };

      if (values.config_select_data) {
        try {
          payload.config_select_data = JSON.parse(values.config_select_data);
        } catch {
          MessagePlugin.error('可选项数据需为 JSON 数组');
          return;
        }
      } else {
        payload.config_select_data = undefined;
      }

      modalApi.setState({ confirmLoading: true });
      if (payload.id) {
        await updateConfig(payload);
      } else {
        await saveConfig(payload);
      }
      MessagePlugin.success(payload.id ? '更新成功' : '新增成功');
      emit('success');
      modalApi.close();
    } catch (error) {
      console.error(error);
    } finally {
      modalApi.setState({ confirmLoading: false });
    }
  },
  class: 'w-[640px]',
});

interface OpenOptions {
  data?: ConfigApi.ConfigSubmitPayload;
}

async function fetchGroupOptions() {
  const list = await getConfigGroupList();
  groupOptions.value = list.map((item) => ({
    label: item.name,
    value: item.id,
  }));
  formApi.updateSchema([
    {
      fieldName: 'group_id',
      componentProps: { options: groupOptions.value },
    },
  ]);
}

async function fetchStatusOptions() {
  statusOptions.value =
    (await getDictOptions('data_status')) || [
      { label: '正常', value: 1 },
      { label: '停用', value: 2 },
    ];
  formApi.updateSchema([
    {
      fieldName: 'status',
      componentProps: { options: statusOptions.value },
    },
  ]);
}

async function open(options?: OpenOptions) {
  modalApi.setState({
    title: options?.data?.id ? '编辑配置' : '新增配置',
  });
  modalApi.open();

  await Promise.all([fetchGroupOptions(), fetchStatusOptions()]);
  await formApi.resetForm();
  const defaults = createConfigFormDefaultValues();
  if (options?.data) {
    const configData = { ...options.data } as ConfigApi.ConfigSubmitPayload & {
      config_select_data?: string;
    };
    if (Array.isArray(configData.config_select_data)) {
      configData.config_select_data = JSON.stringify(
        configData.config_select_data,
      );
    }
    formApi.setValues({
      ...defaults,
      ...configData,
      group_id: configData.group_id ?? defaults.group_id,
    });
  } else {
    formApi.setValues(defaults);
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
