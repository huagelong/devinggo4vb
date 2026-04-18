<script lang="ts" setup>
import { logger } from '#/utils/logger';
import type { MenuApi } from '#/api/system/menu';
import type { RoleApi } from '#/api/system/role';

import { computed, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';
import { $t } from '@vben/locales';

import {
  Checkbox,
  Form,
  FormItem,
  Input,
  Space,
  Tree,
  MessagePlugin,
} from 'tdesign-vue-next';

import { getMenuTreeOptions } from '#/api/system/menu';
import {
  getMenuByRole,
  updateRoleMenuPermission,
} from '#/api/system/role';

const emit = defineEmits(['success']);

const loading = ref(false);
const currentRole = ref<null | RoleApi.ListItem>(null);
const menuList = ref<MenuApi.TreeOptionItem[]>([]);
const checkedKeys = ref<Array<number | string>>([]);
const expandedKeys = ref<Array<number | string>>([]);
const cancelLinkage = ref(false);
const searchText = ref('');

const treeKeys = {
  value: 'id',
  label: 'label',
  children: 'children',
};

function flattenTreeIds(nodes: MenuApi.TreeOptionItem[]) {
  const ids: Array<number | string> = [];
  const travel = (items: MenuApi.TreeOptionItem[]) => {
    items.forEach((item) => {
      ids.push(item.id);
      if (item.children?.length) {
        travel(item.children);
      }
    });
  };
  travel(nodes);
  return ids;
}

const allNodeKeys = computed(() => flattenTreeIds(menuList.value));

function handleExpand(checked: boolean) {
  expandedKeys.value = checked ? [...allNodeKeys.value] : [];
}

function handleSelect(checked: boolean) {
  checkedKeys.value = checked ? [...allNodeKeys.value] : [];
}

function handleLinkage(checked: boolean) {
  cancelLinkage.value = checked;
}

function handleTreeClick(context: { node?: { value?: number | string } }) {
  const nodeValue = context?.node?.value;
  if (nodeValue === undefined) return;
  expandedKeys.value = expandedKeys.value.includes(nodeValue)
    ? expandedKeys.value.filter((item) => item !== nodeValue)
    : [...expandedKeys.value, nodeValue];
}

function filterTreeNode(node: { label?: string }) {
  const label = node?.label ?? '';
  return !searchText.value || label.includes(searchText.value);
}

function extractMenuIds(list?: RoleApi.RoleMenuPermission[]) {
  return (
    list?.flatMap((item) => item.menus?.map((menu) => Number(menu.id)) ?? []) ??
    []
  );
}

const [Modal, modalApi] = useVbenModal({
  onConfirm: async () => {
    const role = currentRole.value;
    if (!role) return;

    try {
      modalApi.setState({ confirmLoading: true });
      await updateRoleMenuPermission(role.id, {
        menu_ids: checkedKeys.value.map((item) => Number(item)),
      });
      MessagePlugin.success($t('common.menuPermissionUpdateSuccess'));
      emit('success');
      modalApi.close();
    } catch (error) {
      logger.error(error);
    } finally {
      modalApi.setState({ confirmLoading: false });
    }
  },
  class: 'w-[760px]',
});

async function open(role: RoleApi.ListItem) {
  currentRole.value = role;
  searchText.value = '';
  menuList.value = [];
  checkedKeys.value = [];
  expandedKeys.value = [];
  cancelLinkage.value = false;

  modalApi.setState({
    title: '菜单权限',
  });
  modalApi.open();

  loading.value = true;
  try {
    const [menuTree, relationData] = await Promise.all([
      getMenuTreeOptions({ scope: true }).catch(
        () => [] as MenuApi.TreeOptionItem[],
      ),
      getMenuByRole(role.id).catch(() => [] as RoleApi.RoleMenuPermission[]),
    ]);

    menuList.value = menuTree;
    checkedKeys.value = extractMenuIds(relationData);
    if (checkedKeys.value.length > 0) {
      cancelLinkage.value = true;
    }
  } finally {
    loading.value = false;
  }
}

defineExpose({
  open,
});
</script>

<template>
  <Modal>
    <Form :data="currentRole ?? {}" label-width="100px" colon>
      <FormItem label="角色名称" name="name">
        <Input :model-value="currentRole?.name ?? ''" disabled />
      </FormItem>
      <FormItem label="角色标识" name="code">
        <Input :model-value="currentRole?.code ?? ''" disabled />
      </FormItem>
      <FormItem label="搜索菜单" name="search">
        <Input v-model="searchText" placeholder="过滤菜单" clearable />
      </FormItem>
      <FormItem label="菜单列表" name="menu_ids">
        <div class="w-full">
          <Space class="mb-3">
            <Checkbox @change="handleExpand">展开/折叠</Checkbox>
            <Checkbox @change="handleSelect">全选/全不选</Checkbox>
            <Checkbox
              :checked="cancelLinkage"
              @change="handleLinkage"
            >
              关闭父子级联动
            </Checkbox>
          </Space>
          <div class="tree-container">
            <div
              v-if="loading"
              class="flex h-[320px] items-center justify-center text-sm text-gray-500"
            >
              菜单加载中...
            </div>
            <Tree
              v-else
              v-model="checkedKeys"
              v-model:expanded="expandedKeys"
              :check-strictly="cancelLinkage"
              :data="menuList"
              :filter="filterTreeNode"
              :keys="treeKeys"
              :value-mode="'all'"
              checkable
              hover
              line
              @click="handleTreeClick"
            />
          </div>
        </div>
      </FormItem>
    </Form>
  </Modal>
</template>

<style scoped>
.tree-container {
  border: 1px solid var(--td-component-border);
  border-radius: 6px;
  max-height: 360px;
  min-height: 320px;
  overflow: auto;
  padding: 8px;
}
</style>
