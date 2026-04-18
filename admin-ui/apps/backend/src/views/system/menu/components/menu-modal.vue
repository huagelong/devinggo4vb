<script lang="ts" setup>
import { logger } from '#/utils/logger';
import type { MenuApi } from '#/api/system/menu';
import type { DictOption } from '#/composables/crud/use-dict-options';

import { nextTick, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import { MessagePlugin } from 'tdesign-vue-next';

import { useVbenForm } from '#/adapter/form';
import { getMenuTreeOptions, saveMenu, updateMenu } from '#/api/system/menu';
import { useDictOptions } from '#/composables/crud/use-dict-options';

import IconSelect from './icon-select.vue';

import {
  createMenuFormDefaultValues,
  menuHiddenOptions,
  menuTypeOptions,
  restfulOptions,
} from '../schemas';

const emit = defineEmits(['success']);

const parentMenuOptions = ref<MenuApi.TreeOptionItem[]>([]);
const statusOptions = ref<DictOption[]>([]);

const { getDictOptions } = useDictOptions();

const [Form, formApi] = useVbenForm({
  showDefaultActions: false,
  commonConfig: {
    labelWidth: 100,
  },
  schema: [
    {
      component: 'Input',
      dependencies: { show: false, triggerFields: [''] },
      fieldName: 'id',
      label: 'ID',
    },
    {
      component: 'TreeSelect',
      componentProps: {
        data: parentMenuOptions.value,
        keys: { children: 'children', label: 'label', value: 'value' },
        placeholder: '请选择上级菜单',
      },
      fieldName: 'parent_id',
      label: '上级菜单',
    },
    {
      component: 'Input',
      componentProps: { placeholder: '请输入菜单名称' },
      fieldName: 'name',
      label: '菜单名称',
      rules: 'required',
    },
    {
      component: 'RadioGroup',
      componentProps: { options: menuTypeOptions },
      defaultValue: 'M',
      fieldName: 'type',
      label: '菜单类型',
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: { placeholder: '请输入菜单标识' },
      fieldName: 'code',
      label: '菜单标识',
      rules: 'required',
    },
    {
      component: IconSelect,
      componentProps: { placeholder: '请选择或输入图标名称' },
      dependencies: {
        show: (values) => isFieldVisible('icon', values?.type),
        triggerFields: ['type'],
      },
      fieldName: 'icon',
      label: '图标',
    },
    {
      component: 'Input',
      componentProps: { placeholder: '请输入路由地址' },
      dependencies: {
        show: (values) => isFieldVisible('route', values?.type),
        triggerFields: ['type'],
      },
      fieldName: 'route',
      label: '路由地址',
    },
    {
      component: 'Input',
      componentProps: { placeholder: '请输入组件路径' },
      dependencies: {
        show: (values) => isFieldVisible('component', values?.type),
        triggerFields: ['type'],
      },
      fieldName: 'component',
      label: '组件路径',
    },
    {
      component: 'Input',
      componentProps: { placeholder: '请输入重定向地址' },
      dependencies: {
        show: (values) => isFieldVisible('redirect', values?.type),
        triggerFields: ['type'],
      },
      fieldName: 'redirect',
      label: '重定向',
    },
    {
      component: 'InputNumber',
      componentProps: { min: 0, max: 1000 },
      defaultValue: 1,
      dependencies: {
        show: (values) => isFieldVisible('sort', values?.type),
        triggerFields: ['type'],
      },
      fieldName: 'sort',
      label: '排序',
      rules: 'required',
    },
    {
      component: 'RadioGroup',
      componentProps: { options: menuHiddenOptions },
      defaultValue: 2,
      dependencies: {
        show: (values) => isFieldVisible('is_hidden', values?.type),
        triggerFields: ['type'],
      },
      fieldName: 'is_hidden',
      label: '是否隐藏',
    },
    {
      component: 'RadioGroup',
      componentProps: { options: restfulOptions },
      defaultValue: '2',
      dependencies: {
        show: (values) => isFieldVisible('restful', values?.type),
        triggerFields: ['type'],
      },
      fieldName: 'restful',
      label: '生成按钮',
    },
    {
      component: 'RadioGroup',
      componentProps: {
        options: statusOptions.value,
      },
      defaultValue: 1,
      fieldName: 'status',
      label: '状态',
      rules: 'required',
    },
    {
      component: 'Textarea',
      componentProps: { placeholder: '请输入备注' },
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

      const values = await formApi.getValues<MenuApi.SubmitPayload>();
      modalApi.setState({ confirmLoading: true });
      const payload: MenuApi.SubmitPayload = {
        ...values,
        is_hidden: Number(values.is_hidden),
        parent_id: Number(values.parent_id ?? 0),
        sort: Number(values.sort ?? 1),
        status: Number(values.status ?? 1),
      };

      if (payload.id) {
        await updateMenu(Number(payload.id), payload);
      } else {
        await saveMenu(payload);
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
  class: 'w-[720px]',
});

function isFieldVisible(field: string, type?: string) {
  if (type === 'B') {
    return false;
  }
  if (type === 'I' || type === 'L') {
    return ['icon', 'route', 'sort', 'is_hidden'].includes(field);
  }
  return true;
}

async function fetchParentOptions() {
  try {
    const tree = await getMenuTreeOptions({ onlyMenu: true }).catch(() => []);
    parentMenuOptions.value = [
      {
        id: 0,
        label: '顶级菜单',
        value: 0,
        children: tree as MenuApi.TreeOptionItem[],
      },
    ];
  } catch (error) {
    logger.error(error);
    parentMenuOptions.value = [{ id: 0, label: '顶级菜单', value: 0 }];
  }

  formApi.updateSchema([
    {
      fieldName: 'parent_id',
      componentProps: {
        data: parentMenuOptions.value,
      },
    },
  ]);
}

async function open(data?: Partial<MenuApi.SubmitPayload>) {
  modalApi.setState({
    title: data?.id ? '编辑菜单' : '新增菜单',
  });
  modalApi.open();

  statusOptions.value = (await getDictOptions('data_status')) || [
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

  await fetchParentOptions();
  await formApi.resetForm();
  formApi.setValues(createMenuFormDefaultValues());

  if (data) {
    formApi.setValues({
      ...data,
      parent_id: data.parent_id ?? 0,
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
