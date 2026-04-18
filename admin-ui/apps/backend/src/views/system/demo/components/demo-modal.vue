<script lang="ts" setup>
import { logger } from '#/utils/logger';
import type { DemoApi } from '#/api/system/demo';

import { nextTick } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import { MessagePlugin } from 'tdesign-vue-next';

import { useVbenForm } from '#/adapter/form';
import { saveDemo, updateDemo } from '#/api/system/demo';

import { createDemoFormDefaultValues } from '../schemas';

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
      fieldName: 'name',
      label: '名称',
      componentProps: {
        placeholder: '请输入'
      },
      rules: 'required',
    },
    {
      component: 'Input',
      fieldName: 'code',
      label: '编码',
      componentProps: {
        placeholder: '请输入'
      },
      rules: 'required',
    },
    {
      component: 'RadioGroup',
      fieldName: 'status',
      label: '状态',
      defaultValue: 1,
      componentProps: {
        options: [{"label":"正常","value":1},{"label":"已禁用","value":2}],
      },
      rules: 'required',
    },
    {
      component: 'InputNumber',
      fieldName: 'sort',
      label: '排序',
      defaultValue: 1,
      componentProps: {
        min: 0,
        max: 1000,
      },
      rules: 'required',
    },
    {
      component: 'InputNumber',
      fieldName: 'price',
      label: '数字',
      componentProps: {
        min: 0,
        max: 1000,
        placeholder: '请输入'
      },
    },
    {
      component: 'Upload',
      fieldName: 'cover',
      label: '图片',
      componentProps: {
        accept: 'image/*',
        placeholder: '请上传图片',
      },
      formItemClass: 'col-span-2',
    },
    {
      component: 'Input',
      fieldName: 'email',
      label: '邮箱',
      componentProps: {
        type: 'email',
        placeholder: '请输入邮箱',
      },
      rules: 'email',
    },
    {
      component: 'Input',
      fieldName: 'phone',
      label: '手机号',
      componentProps: {
        placeholder: '请输入手机号',
      },
    },
    {
      component: 'DatePicker',
      fieldName: 'birthday',
      label: '日期',
      componentProps: {
        placeholder: '请选择日期',
        clearable: true,
      },
    },
    {
      component: 'Textarea',
      fieldName: 'remark',
      label: '备注',
      componentProps: {
        placeholder: '请输入备注',
        autosize: { minRows: 3, maxRows: 6 },
      },
      formItemClass: 'col-span-2',
    },
  ],
});

const [Modal, modalApi] = useVbenModal({
  onConfirm: async () => {
    try {
      const { valid } = await formApi.validate();
      if (!valid) return;

      const values = await formApi.getValues<DemoApi.SubmitPayload & { id?: number }>();
      modalApi.setState({ confirmLoading: true });

      if (values.id) {
        await updateDemo(Number(values.id), values);
      } else {
        await saveDemo(values);
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

async function open(data?: Partial<DemoApi.SubmitPayload & { id?: number }>) {
  modalApi.setState({
    title: data?.id ? '编辑演示' : '新增演示',
  });
  modalApi.open();

  await formApi.resetForm();
  formApi.setValues(createDemoFormDefaultValues());
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
