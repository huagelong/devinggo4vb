<script lang="ts" setup>
import { logger } from '#/utils/logger';
import type { PostApi } from '#/api/system/post';

import { nextTick } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import { MessagePlugin } from 'tdesign-vue-next';

import { useVbenForm } from '#/adapter/form';
import { savePost, updatePost } from '#/api/system/post';

import { createPostFormDefaultValues } from '../schemas';

const emit = defineEmits(['success']);

const [Form, formApi] = useVbenForm({
  showDefaultActions: false,
  commonConfig: {
    labelWidth: 90,
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
        placeholder: '请输入岗位名称',
      },
      fieldName: 'name',
      label: '岗位名称',
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入岗位标识',
      },
      fieldName: 'code',
      label: '岗位标识',
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
      label: '备注',
    },
  ],
});

const [Modal, modalApi] = useVbenModal({
  onConfirm: async () => {
    try {
      const { valid } = await formApi.validate();
      if (!valid) return;

      const values = await formApi.getValues<PostApi.SubmitPayload>();
      modalApi.setState({ confirmLoading: true });

      if (values.id) {
        await updatePost(Number(values.id), values);
      } else {
        await savePost(values);
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
  class: 'w-[560px]',
});

async function open(data?: Partial<PostApi.SubmitPayload>) {
  modalApi.setState({
    title: data?.id ? '编辑岗位' : '新增岗位',
  });
  modalApi.open();

  await formApi.resetForm();
  formApi.setValues(createPostFormDefaultValues());
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
