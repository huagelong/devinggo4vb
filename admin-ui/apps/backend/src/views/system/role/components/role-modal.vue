<script lang="ts" setup>
import { logger } from '#/utils/logger';
import type { RoleApi } from '#/api/system/role';

import { nextTick, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import { MessagePlugin } from 'tdesign-vue-next';

import { useVbenForm } from '#/adapter/form';
import { saveRole, updateRole } from '#/api/system/role';

import { createRoleFormDefaultValues } from '../schemas';

const emit = defineEmits(['success']);

const baseValues = ref<RoleApi.SubmitPayload>(createRoleFormDefaultValues());

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
        placeholder: '请输入角色名称',
      },
      fieldName: 'name',
      label: '角色名称',
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入角色标识',
      },
      fieldName: 'code',
      label: '角色标识',
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
    try {
      const { valid } = await formApi.validate();
      if (!valid) return;

      const values = await formApi.getValues<Partial<RoleApi.SubmitPayload>>();
      const payload: RoleApi.SubmitPayload = {
        ...baseValues.value,
        ...values,
      };

      modalApi.setState({ confirmLoading: true });

      if (payload.id) {
        await updateRole(Number(payload.id), payload);
      } else {
        await saveRole(payload);
      }

      MessagePlugin.success(payload.id ? '更新成功' : '新增成功');
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

async function open(data?: Partial<RoleApi.SubmitPayload>) {
  const defaultValues = createRoleFormDefaultValues();
  baseValues.value = {
    ...defaultValues,
    ...data,
  };

  modalApi.setState({
    title: data?.id ? '编辑角色' : '新增角色',
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
