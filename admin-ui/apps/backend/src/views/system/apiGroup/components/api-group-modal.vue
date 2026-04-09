<script lang="ts" setup>
import type { ApiGroupFormModel } from '../model';
import type { DictOption } from '#/composables/crud/use-dict-options';

import { nextTick, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import { MessagePlugin } from 'tdesign-vue-next';

import { useVbenForm } from '#/adapter/form';
import { saveApiGroup, updateApiGroup } from '#/api/system/api-group';
import { useDictOptions } from '#/composables/crud/use-dict-options';

import { createApiGroupFormDefaultValues } from '../schemas';

const emit = defineEmits(['success']);

const statusOptions = ref<DictOption[]>([]);
const fallbackStatusOptions: DictOption[] = [
  { label: '正常', value: 1 },
  { label: '停用', value: 2 },
];

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
      componentProps: { placeholder: '请输入分组名称' },
      fieldName: 'name',
      label: '分组名称',
      rules: 'required',
    },
    {
      component: 'RadioGroup',
      componentProps: { options: statusOptions.value },
      fieldName: 'status',
      label: '״̬',
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
      const values = await formApi.getValues<ApiGroupFormModel>();
      modalApi.setState({ confirmLoading: true });
      if (values.id) {
        await updateApiGroup(Number(values.id), values);
      } else {
        await saveApiGroup(values);
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
  class: 'w-[520px]',
});

async function fetchStatusOptions() {
  try {
    const options = await getDictOptions('data_status');
    statusOptions.value =
      options && options.length > 0 ? options : fallbackStatusOptions;
  } catch (error) {
    console.error(error);
    statusOptions.value = fallbackStatusOptions;
  } finally {
    formApi.updateSchema([
      {
        fieldName: 'status',
        componentProps: { options: statusOptions.value },
      },
    ]);
  }
}

async function open(data?: ApiGroupFormModel) {
  modalApi.setState({
    title: data?.id ? '编辑分组' : '新增分组',
  });
  modalApi.open();
  await fetchStatusOptions();
  await formApi.resetForm();
  formApi.setValues(createApiGroupFormDefaultValues());
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
