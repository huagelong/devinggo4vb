<script lang="ts" setup>
import { logger } from '#/utils/logger';
import type { DictApi } from '#/api/system/dict';
import type { DictOption } from '#/composables/crud/use-dict-options';

import { nextTick, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { MessagePlugin } from 'tdesign-vue-next';

import { useVbenForm } from '#/adapter/form';
import { saveDictData, updateDictData } from '#/api/system/dict';
import { useDictOptions } from '#/composables/crud/use-dict-options';

import { createDictDataFormDefaultValues } from '../schemas';

const emit = defineEmits(['success']);

const statusOptions = ref<DictOption[]>([]);
const currentTypeInfo = ref<{ id: number; code: string; name?: string }>();
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
      dependencies: { show: false, triggerFields: [''] },
      fieldName: 'type_id',
      label: 'type_id',
    },
    {
      component: 'Input',
      dependencies: { show: false, triggerFields: [''] },
      fieldName: 'code',
      label: 'code',
    },
    {
      component: 'Input',
      componentProps: { placeholder: '请输入字典标签' },
      fieldName: 'label',
      label: '字典标签',
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: { placeholder: '请输入字典键值' },
      fieldName: 'value',
      label: '字典键值',
      rules: 'required',
    },
    {
      component: 'InputNumber',
      componentProps: { min: 0, max: 1000 },
      defaultValue: 1,
      fieldName: 'sort',
      label: '排序',
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

      const values = await formApi.getValues<DictApi.DictDataSubmitPayload>();
      if (!values.type_id) {
        MessagePlugin.warning($t('common.missingDictTypeInfo'));
        return;
      }
      modalApi.setState({ confirmLoading: true });

      if (values.id) {
        await updateDictData(Number(values.id), values);
      } else {
        await saveDictData(values);
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

interface OpenOptions {
  data?: DictApi.DictDataSubmitPayload;
  typeInfo: { id: number; code: string; name?: string };
}

async function open(options: OpenOptions) {
  currentTypeInfo.value = options.typeInfo;
  modalApi.setState({
    title: options.data?.id ? '编辑字典数据' : `新增「${options.typeInfo.name ?? ''}」数据`,
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
  formApi.setValues(
    createDictDataFormDefaultValues(options.typeInfo.id, options.typeInfo.code),
  );
  if (options.data) {
    formApi.setValues({
      ...options.data,
      type_id: options.typeInfo.id,
      code: options.typeInfo.code,
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
