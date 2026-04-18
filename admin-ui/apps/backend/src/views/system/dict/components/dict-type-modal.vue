<script lang="ts" setup>
import { logger } from '#/utils/logger';
import type { DictApi } from '#/api/system/dict';
import type { DictOption } from '#/composables/crud/use-dict-options';

import { nextTick, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import { MessagePlugin } from 'tdesign-vue-next';

import { useVbenForm } from '#/adapter/form';
import { saveDictType, updateDictType } from '#/api/system/dict';
import { useDictOptions } from '#/composables/crud/use-dict-options';

import { createDictTypeFormDefaultValues } from '../schemas';

const emit = defineEmits(['success']);

const statusOptions = ref<DictOption[]>([]);
const { getDictOptions } = useDictOptions();

const [Form, formApi] = useVbenForm({
  showDefaultActions: false,
  commonConfig: {
    labelWidth: 90,
  },
  schema: [
    {
      component: 'Input',
      dependencies: { show: false, triggerFields: [''] },
      fieldName: 'id',
      label: 'ID',
    },
    {
      component: 'Input',
      componentProps: { placeholder: '请输入字典名称' },
      fieldName: 'name',
      label: '字典名称',
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: { placeholder: '请输入字典标识' },
      fieldName: 'code',
      label: '字典标识',
      rules: 'required',
    },
    {
      component: 'RadioGroup',
      componentProps: { options: statusOptions.value },
      defaultValue: 1,
      fieldName: 'status',
      label: '状态',
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

      const values = await formApi.getValues<DictApi.DictTypeSubmitPayload>();
      modalApi.setState({ confirmLoading: true });

      if (values.id) {
        await updateDictType(Number(values.id), values);
      } else {
        await saveDictType(values);
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

async function open(data?: DictApi.DictTypeSubmitPayload) {
  modalApi.setState({
    title: data?.id ? '编辑字典' : '新增字典',
  });
  modalApi.open();

  statusOptions.value =
    (await getDictOptions('data_status')) || [
      { label: '正常', value: 1 },
      { label: '停用', value: 2 },
    ];

  formApi.updateSchema([
    {
      fieldName: 'status',
      componentProps: {
        options: statusOptions.value,
      },
    },
  ]);

  await formApi.resetForm();
  formApi.setValues(createDictTypeFormDefaultValues());
  if (data) {
    formApi.setValues({
      ...data,
      status: Number(data.status ?? 1),
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
