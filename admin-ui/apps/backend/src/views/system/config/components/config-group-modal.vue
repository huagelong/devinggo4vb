<script lang="ts" setup>
import { logger } from '#/utils/logger';
import type { ConfigApi } from '#/api/system/config';

import { nextTick } from 'vue';

import { useVbenModal } from '@vben/common-ui';
import { $t } from '@vben/locales';

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
      componentProps: { placeholder: $t('system.config.placeholder.enterGroupName') },
      fieldName: 'name',
      label: $t('system.config.groupName'),
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: { placeholder: $t('system.config.placeholder.enterGroupCode') },
      fieldName: 'code',
      label: $t('system.config.groupCode'),
      rules: 'required',
    },
    {
      component: 'Textarea',
      componentProps: { placeholder: $t('system.config.placeholder.enterRemark') },
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
      const values = await formApi.getValues<ConfigApi.ConfigGroupSubmitPayload>();
      modalApi.setState({ confirmLoading: true });
      await saveConfigGroup(values);
      MessagePlugin.success($t('common.saveSuccess'));
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
  modalApi.setState({ title: $t('system.config.addGroupTitle') });
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
