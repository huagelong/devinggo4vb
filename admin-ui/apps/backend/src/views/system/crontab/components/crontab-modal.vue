<script lang="ts" setup>
import { logger } from '#/utils/logger';
import type { CrontabApi } from '#/api/system/crontab';

import { nextTick, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import { MessagePlugin } from 'tdesign-vue-next';

import { useVbenForm } from '#/adapter/form';
import { saveCrontab, updateCrontab } from '#/api/system/crontab';

import {
  crontabFinallyOptions,
  crontabTypeOptions,
  createCrontabFormDefaultValues,
} from '../schemas';

const emit = defineEmits(['success']);

const baseValues = ref<CrontabApi.SubmitPayload>(
  createCrontabFormDefaultValues(),
);

const [Form, formApi] = useVbenForm({
  showDefaultActions: false,
  commonConfig: {
    labelWidth: 100,
  },
  schema: [
    {
      component: 'Input',
      dependencies: {
        show: false,
        triggerFields: [''],
      },
      fieldName: 'id',
      label: 'ID',
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入任务名称',
      },
      fieldName: 'name',
      label: '任务名称',
      rules: 'required',
    },
    {
      component: 'RadioGroup',
      componentProps: {
        options: crontabTypeOptions,
      },
      defaultValue: 1,
      fieldName: 'type',
      label: '任务类型',
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入执行规则，如：* * * * *',
      },
      fieldName: 'rule',
      label: '执行规则',
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入调用目标，如：App\\Task\\TestTask@run',
      },
      fieldName: 'target',
      label: '调用目标',
      rules: 'required',
    },
    {
      component: 'RadioGroup',
      componentProps: {
        options: crontabFinallyOptions,
      },
      defaultValue: 2,
      fieldName: 'is_finally',
      label: '最终执行',
    },
    {
      component: 'Textarea',
      componentProps: {
        placeholder: '请输入备注',
      },
      fieldName: 'remark',
      formItemClass: 'col-span-2',
      label: '备注',
    },
  ],
});

const [Modal, modalApi] = useVbenModal({
  onConfirm: async () => {
    let isEdit = false;
    try {
      const { valid } = await formApi.validate();
      if (!valid) return;

      const values = await formApi.getValues<Partial<CrontabApi.SubmitPayload>>();
      const payload: CrontabApi.SubmitPayload = {
        ...baseValues.value,
        ...values,
      };
      isEdit = !!payload.id;

      modalApi.setState({ confirmLoading: true });

      if (payload.id) {
        await updateCrontab(Number(payload.id), payload);
      } else {
        await saveCrontab(payload);
      }

      MessagePlugin.success(isEdit ? '更新成功' : '新增成功');
      emit('success');
      modalApi.close();
    } catch (error) {
      logger.error(error);
      MessagePlugin.error(isEdit ? '更新失败' : '新增失败');
    } finally {
      modalApi.setState({ confirmLoading: false });
    }
  },
  class: 'w-[640px]',
});

async function open(data?: Partial<CrontabApi.SubmitPayload>) {
  const defaultValues = createCrontabFormDefaultValues();
  baseValues.value = {
    ...defaultValues,
    ...data,
  };

  modalApi.setState({
    title: data?.id ? '编辑定时任务' : '新增定时任务',
  });
  modalApi.open();

  await formApi.resetForm();
  formApi.setValues(baseValues.value);
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
