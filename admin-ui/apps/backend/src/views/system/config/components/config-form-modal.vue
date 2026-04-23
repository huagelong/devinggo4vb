<script lang="ts" setup>
import { logger } from '#/utils/logger';
import type { ConfigApi } from '#/api/system/config';
import type { DictOption } from '#/composables/crud/use-dict-options';

import { nextTick, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';
import { $t } from '@vben/locales';

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
        placeholder: $t('ui.placeholder.select'),
      },
      fieldName: 'group_id',
      label: $t('system.config.group'),
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: { placeholder: $t('system.config.placeholder.enterConfigName') },
      fieldName: 'name',
      label: $t('system.config.name'),
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: { placeholder: $t('system.config.placeholder.enterConfigCode') },
      fieldName: 'key',
      label: $t('system.config.code'),
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: { placeholder: $t('system.config.placeholder.enterConfigValue') },
      fieldName: 'value',
      label: $t('system.config.value'),
    },
    {
      component: 'InputNumber',
      componentProps: { min: 0, max: 999 },
      defaultValue: 0,
      fieldName: 'sort',
      label: $t('common.sort'),
    },
    {
      component: 'Select',
      componentProps: {
        options: inputComponentOptions,
        placeholder: $t('ui.placeholder.select'),
      },
      defaultValue: 'input',
      fieldName: 'input_type',
      label: $t('system.config.inputComponent'),
      rules: 'required',
    },
    {
      component: 'Select',
      componentProps: {
        options: statusOptions.value,
        placeholder: $t('ui.placeholder.select'),
      },
      defaultValue: 1,
      fieldName: 'status',
      label: $t('common.status'),
      rules: 'required',
    },
    {
      component: 'Textarea',
      componentProps: { placeholder: $t('system.config.placeholder.enterRemark') },
      fieldName: 'remark',
      label: $t('common.remark'),
    },
    {
      component: 'Textarea',
      componentProps: {
        autosize: { minRows: 3, maxRows: 6 },
        placeholder: $t('system.config.optionsPlaceholder'),
      },
      dependencies: {
        show: (values) =>
          ['select', 'radio', 'checkbox'].includes(values?.input_type),
        triggerFields: ['input_type'],
      },
      fieldName: 'config_select_data',
      label: $t('system.config.options'),
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

      if (values.config_select_data?.trim()) {
        try {
          const parsed = JSON.parse(values.config_select_data);
          if (!Array.isArray(parsed)) {
            MessagePlugin.error($t('common.jsonArrayRequired'));
            return;
          }
          payload.config_select_data = values.config_select_data;
        } catch {
          MessagePlugin.error($t('common.jsonArrayRequired'));
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
      MessagePlugin.success(payload.id ? $t('common.updateSuccess') : $t('common.createSuccess'));
      emit('success');
      modalApi.close();
    } catch (error) {
      logger.error(error);
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
      { label: $t('common.statusEnabled'), value: 1 },
      { label: $t('common.statusDisabled'), value: 2 },
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
    title: options?.data?.id ? $t('system.config.editTitle') : $t('system.config.addConfigTitle'),
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
