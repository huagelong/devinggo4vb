<script lang="ts" setup>
import { logger } from '#/utils/logger';
import type { AppGroupApi } from '#/api/system/app-group';

import { nextTick, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import { MessagePlugin } from 'tdesign-vue-next';

import { useVbenForm } from '#/adapter/form';
import { saveAppGroup, updateAppGroup } from '#/api/system/app-group';

import { createAppGroupFormDefaultValues } from '../schemas';

const emit = defineEmits(['success']);

const baseValues = ref<AppGroupApi.SubmitPayload>(
  createAppGroupFormDefaultValues(),
);

const [Form, formApi] = useVbenForm({
  showDefaultActions: false,
  commonConfig: {
    labelWidth: 80,
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
        placeholder: '请输入分组名称',
      },
      fieldName: 'name',
      label: '分组名称',
      rules: 'required',
    },
    {
      component: 'InputNumber',
      componentProps: {
        max: 1000,
        min: 0,
      },
      defaultValue: 1,
      fieldName: 'sort',
      label: '排序',
      rules: 'required',
    },
    {
      component: 'RadioGroup',
      componentProps: {
        options: [
          { label: '正常', value: 1 },
          { label: '停用', value: 2 },
        ],
      },
      defaultValue: 1,
      fieldName: 'status',
      label: '状态',
      rules: 'required',
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

      const values = await formApi.getValues<Partial<AppGroupApi.SubmitPayload>>();
      const payload: AppGroupApi.SubmitPayload = {
        ...baseValues.value,
        ...values,
      };
      isEdit = !!payload.id;

      modalApi.setState({ confirmLoading: true });

      if (payload.id) {
        await updateAppGroup(Number(payload.id), payload);
      } else {
        await saveAppGroup(payload);
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
  class: 'w-[480px]',
});

async function open(data?: Partial<AppGroupApi.SubmitPayload>) {
  const defaultValues = createAppGroupFormDefaultValues();
  baseValues.value = {
    ...defaultValues,
    ...data,
  };

  modalApi.setState({
    title: data?.id ? '编辑分组' : '新增分组',
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
