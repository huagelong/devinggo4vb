<script lang="ts" setup>
import { logger } from '#/utils/logger';
import type { DeptApi } from '#/api/system/dept';
import type { PostApi } from '#/api/system/post';
import type { RoleApi } from '#/api/system/role';
import type { UserApi } from '#/api/system/user';

import { defineComponent, h, markRaw, nextTick } from 'vue';

import { useVbenModal } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { MessagePlugin } from 'tdesign-vue-next';

import { useVbenForm } from '#/adapter/form';
import { uploadImageFileApi } from '#/api/system/upload';
import { getDeptTree } from '#/api/system/dept';
import { getPostList } from '#/api/system/post';
import { getRoleList } from '#/api/system/role';
import { getUserDetail, saveUser, updateUser } from '#/api/system/user';

const emit = defineEmits(['success']);

interface AvatarUploadResponse {
  url?: string;
}

type UserModalOpenData = Partial<UserApi.Detail>;

function extractEntityIds(list?: UserApi.RelatedEntity[]) {
  return list?.map((item) => item.id) ?? [];
}

function normalizeListData<T>(data: T[] | { items?: T[] } | null | undefined): T[] {
  if (Array.isArray(data)) {
    return data;
  }
  return Array.isArray(data?.items) ? data.items : [];
}

const AvatarUpload = defineComponent({
  props: {
    modelValue: { type: String, default: '' },
    value: { type: String, default: '' },
  },
  emits: ['update:modelValue', 'update:value', 'change'],
  setup(props, { emit: emitInner }) {
    const getAvatarValue = () => props.modelValue || props.value;

    function handleClick() {
      const input = document.createElement('input');
      input.type = 'file';
      input.accept = 'image/*';
      input.onchange = async (event: Event) => {
        const file = (event.target as HTMLInputElement).files?.[0];
        if (!file) return;

        try {
          const res = (await uploadImageFileApi(file)) as AvatarUploadResponse;
          if (res?.url) {
            emitInner('update:modelValue', res.url);
            emitInner('update:value', res.url);
            emitInner('change', res.url);
          }
        } catch {
          MessagePlugin.error($t('common.avatarUploadFailed2'));
        }
      };
      input.click();
    }

    return () =>
      h('div', { class: 'flex items-center gap-3' }, [
        h(
          'div',
          {
            onClick: handleClick,
            class:
              'relative flex h-16 w-16 cursor-pointer items-center justify-center overflow-hidden rounded-full border-2 border-dashed border-gray-300 bg-gray-50 hover:border-blue-400',
          },
          [
            getAvatarValue()
              ? h('img', {
                  src: getAvatarValue(),
                  class: 'h-full w-full rounded-full object-cover',
                  alt: '头像',
                })
              : h(
                  'span',
                  { class: 'text-3xl leading-none text-gray-400' },
                  '+',
                ),
          ],
        ),
        h('div', { class: 'flex flex-col gap-1' }, [
          h('span', { class: 'text-sm text-gray-500' }, '点击上传头像'),
          h(
            'span',
            { class: 'text-xs text-gray-400' },
            '支持 JPG、PNG 等图片格式',
          ),
        ]),
      ]);
  },
});

const [Form, formApi] = useVbenForm({
  showDefaultActions: false,
  commonConfig: {
    labelWidth: 80,
  },
  wrapperClass: 'grid-cols-2',
  schema: [
    {
      fieldName: 'id',
      label: 'ID',
      component: 'Input',
      dependencies: {
        show: false,
        triggerFields: [''],
      },
    },
    {
      fieldName: 'avatar',
      label: '头像',
      component: markRaw(AvatarUpload),
      formItemClass: 'col-span-2',
    },
    {
      fieldName: 'username',
      label: '账户',
      component: 'Input',
      rules: 'required',
      componentProps: { placeholder: '请输入账户' },
    },
    {
      fieldName: 'dept_ids',
      label: '所属部门',
      component: 'TreeSelect',
      componentProps: {
        data: [],
        keys: { label: 'label', value: 'value', children: 'children' },
        multiple: true,
        placeholder: '请选择所属部门',
      },
      rules: 'required',
    },
    {
      fieldName: 'password',
      label: '密码',
      component: 'Input',
      dependencies: {
        show: (values) => !values?.id,
        triggerFields: ['id'],
      },
      rules: 'required',
      componentProps: {
        placeholder: '请输入密码',
        type: 'password',
      },
    },
    {
      fieldName: 'nickname',
      label: '昵称',
      component: 'Input',
      componentProps: { placeholder: '请输入昵称' },
    },
    {
      fieldName: 'role_ids',
      label: '角色',
      component: 'Select',
      componentProps: {
        options: [],
        keys: { label: 'name', value: 'id' },
        multiple: true,
        placeholder: '请选择角色',
      },
      rules: 'required',
    },
    {
      fieldName: 'phone',
      label: '手机',
      component: 'Input',
      componentProps: { placeholder: '请输入手机' },
    },
    {
      fieldName: 'post_ids',
      label: '岗位',
      component: 'Select',
      componentProps: {
        options: [],
        keys: { label: 'name', value: 'id' },
        multiple: true,
        placeholder: '请选择岗位',
      },
    },
    {
      fieldName: 'email',
      label: '邮箱',
      component: 'Input',
      componentProps: { placeholder: '请输入邮箱' },
    },
    {
      fieldName: 'status',
      label: '状态',
      component: 'RadioGroup',
      defaultValue: 1,
      componentProps: {
        options: [
          { label: '正常', value: 1 },
          { label: '停用', value: 2 },
        ],
      },
      formItemClass: 'col-span-2',
    },
    {
      fieldName: 'user_type',
      label: '用户类型',
      component: 'Select',
      defaultValue: '100',
      componentProps: {
        options: [{ label: '系统用户', value: '100' }],
        placeholder: '请选择用户类型',
      },
      rules: 'required',
      formItemClass: 'col-span-2',
    },
    {
      fieldName: 'remark',
      label: '备注',
      component: 'Textarea',
      formItemClass: 'col-span-2',
      componentProps: { placeholder: '请输入备注' },
    },
  ],
});

const [Modal, modalApi] = useVbenModal({
  onConfirm: async () => {
    try {
      const { valid } = await formApi.validate();
      if (!valid) return;

      const values = await formApi.getValues<UserApi.SubmitPayload>();
      modalApi.setState({ confirmLoading: true });

      values.id ? await updateUser(values.id, values) : await saveUser(values);

      MessagePlugin.success(values.id ? '更新成功' : '新增成功');
      emit('success');
      modalApi.close();
    } catch (error) {
      logger.error(error);
    } finally {
      modalApi.setState({ confirmLoading: false });
    }
  },
  class: 'w-[800px]',
});

async function open(data?: UserModalOpenData) {
  modalApi.setState({ title: data?.id ? '编辑管理员' : '新增管理员' });
  modalApi.open();

  const [roleRes, postRes, deptRes] = await Promise.all([
    getRoleList().catch(() => null),
    getPostList().catch(() => null),
    getDeptTree().catch(() => null),
  ]);

  const roleOptions = normalizeListData<RoleApi.ListItem>(roleRes);
  const postOptions = normalizeListData<PostApi.ListItem>(postRes);
  const deptData = Array.isArray(deptRes) ? (deptRes as DeptApi.TreeNode[]) : [];

  formApi.updateSchema([
    {
      fieldName: 'role_ids',
      componentProps: { options: roleOptions },
    },
    {
      fieldName: 'post_ids',
      componentProps: { options: postOptions },
    },
    {
      fieldName: 'dept_ids',
      componentProps: { data: deptData },
    },
  ]);

  await formApi.resetForm();

  if (data?.id) {
    const detail = await getUserDetail(data.id).catch(() => null);
    const detailValues = detail
      ? {
          ...detail,
          dept_ids:
            detail.dept_ids ??
            extractEntityIds(detail.deptList) ??
            data.dept_ids ??
            [],
          role_ids:
            detail.role_ids ??
            extractEntityIds(detail.roleList) ??
            data.role_ids ??
            [],
          post_ids:
            detail.post_ids ??
            extractEntityIds(detail.postList) ??
            data.post_ids ??
            [],
        }
      : data;
    formApi.setValues(detailValues);
  } else if (data) {
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
