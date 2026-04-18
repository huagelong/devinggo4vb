<script lang="ts" setup>
import { logger } from '#/utils/logger';
import type { ConfigApi } from '#/api/system/config';

import { nextTick } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import { MessagePlugin } from 'tdesign-vue-next';

import { useVbenForm } from '#/adapter/form';

import { saveConfigGroup } from '#/api/system/config';

import { createConfigGroupFormDefaultValues } from '../schemas';

const emit = defineEmits(['success']);

const [Form, formApi] = useVbenForm({
  showDefaultActions: false,
  commonConfig: { labelWidth: 90 },
  schema: [
    {
      component: 'Input',
      dependencies: { show: false, triggerFields: [''] },
      fieldName: 'id',
      label: 'ID',
    },
    {
      component: 'Input',
      componentProps: { placeholder: '请输入分组名称' },
      fieldName: 'name',
      label: '分组名称',
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: { placeholder: '请输入分组标识' },
      fieldName: 'code',
      label: '分组标识',
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
      const values = await formApi.getValues<ConfigApi.ConfigGroupSubmitPayload>();
      modalApi.setState({ confirmLoading: true });
      await saveConfigGroup(values);
      MessagePlugin.success('保存成功');
      emit('success');
      modalApi.close();
    } catch (error) {
      logger.error(error);
    } finally {
      modalApi.setState({ confirmLoading: false });
    }
  },
  class: 'w-[520px]',
});

async function open() {
  modalApi.setState({ title: '新增配置分组' });
  modalApi.open();
  await formApi.resetForm();
  formApi.setValues(createConfigGroupFormDefaultValues());
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
