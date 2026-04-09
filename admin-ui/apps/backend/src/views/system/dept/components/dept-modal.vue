<script lang="ts" setup>
import type { DeptApi } from '#/api/system/dept';

import { nextTick } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import { MessagePlugin } from 'tdesign-vue-next';

import { useVbenForm } from '#/adapter/form';
import { getDeptTree, saveDept, updateDept } from '#/api/system/dept';

import { createDeptFormDefaultValues } from '../schemas';

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
      component: 'TreeSelect',
      componentProps: {
        data: [],
        keys: { label: 'label', value: 'value', children: 'children' },
        placeholder: '请选择上级部门',
      },
      fieldName: 'parent_id',
      label: '上级部门',
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入部门名称',
      },
      fieldName: 'name',
      label: '部门名称',
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入负责人',
      },
      fieldName: 'leader',
      label: '负责人',
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: '请输入手机号',
      },
      fieldName: 'phone',
      label: '手机',
      rules: 'pattern:^1[3-9]\\d{9}$#请输入正确的手机号',
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
      label: '״̬',
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

      const values = await formApi.getValues<DeptApi.SubmitPayload>();
      modalApi.setState({ confirmLoading: true });

      if (values.id) {
        await updateDept(Number(values.id), values);
      } else {
        await saveDept(values);
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
  class: 'w-[640px]',
});

async function open(data?: Partial<DeptApi.SubmitPayload>) {
  modalApi.setState({
    title: data?.id ? '编辑部门' : '新增部门',
  });
  modalApi.open();

  const deptTree = await getDeptTree().catch(() => [] as DeptApi.TreeNode[]);

  formApi.updateSchema([
    {
      fieldName: 'parent_id',
      componentProps: {
        data: deptTree,
      },
    },
  ]);

  await formApi.resetForm();
  formApi.setValues(createDeptFormDefaultValues());
  if (data) {
    formApi.setValues({
      ...data,
      status:
        data.status !== undefined ? Number(data.status) : createDeptFormDefaultValues().status,
    });
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
