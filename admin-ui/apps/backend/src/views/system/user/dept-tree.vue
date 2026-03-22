<script lang="ts" setup>
import { onMounted, ref } from 'vue';

import {
  CaretDownSmallIcon,
  CaretRightSmallIcon,
  SearchIcon,
} from 'tdesign-icons-vue-next';
import { Input, Tree } from 'tdesign-vue-next';

import { getDeptTree } from '#/api/system/dept';

const emit = defineEmits(['select']);

const treeData = ref<any[]>([]);
const searchText = ref('');
const treeRef = ref();

const treeKeys = {
  value: 'id',
  label: 'label',
  children: 'children',
};

// 展开的节点
const expanded = ref<Array<number | string>>([-1]);
const isFolding = ref(false);

async function fetchDeptTree() {
  try {
    const res = await getDeptTree();
    treeData.value = [
      {
        id: -1,
        label: '所有部门',
        children: res || [],
      },
    ];
  } catch (error) {
    console.error('Failed to fetch dept tree', error);
  }
}

function handleNodeClick(context: any) {
  const { node } = context;
  emit('select', node.value === -1 ? '' : node.value);
}

function toggleExpand() {
  if (isFolding.value) {
    expanded.value = [-1]; // Root node id is usually needed, fallback to root
  } else {
    // Collect all ids
    const ids: Array<number | string> = [];
    const traverse = (nodes: any[]) => {
      nodes.forEach((node) => {
        ids.push(node.id);
        if (node.children) {
          traverse(node.children);
        }
      });
    };
    traverse(treeData.value);
    expanded.value = ids;
  }
  isFolding.value = !isFolding.value;
}

onMounted(() => {
  fetchDeptTree();
});
</script>

<template>
  <div class="flex h-full w-56 flex-col bg-white">
    <div
      class="flex items-center justify-between gap-2 border-b border-gray-100 p-3"
    >
      <Input v-model="searchText" placeholder="搜索部门" class="flex-1">
        <template #prefixIcon>
          <SearchIcon />
        </template>
      </Input>
      <div
        class="cursor-pointer whitespace-nowrap text-[14px] text-blue-500 hover:text-blue-600"
        @click="toggleExpand"
      >
        {{ isFolding ? '展开' : '折叠' }}
      </div>
    </div>
    <div class="custom-tree-wrap flex-1 overflow-auto p-2">
      <Tree
        ref="treeRef"
        v-model:expanded="expanded"
        :data="treeData"
        :filter="(node: any) => !searchText || node.label.includes(searchText)"
        :keys="treeKeys"
        hover
        activable
        expand-all
        line
        @click="handleNodeClick"
      >
        <template #icon="{ node }">
          <CaretDownSmallIcon v-if="node.expanded" />
          <CaretRightSmallIcon v-else />
        </template>
      </Tree>
    </div>
  </div>
</template>

<style scoped>
.custom-tree-wrap {
  :deep(.t-tree) {
    --td-tree-hover-color: #f3f4f6;
  }
}
</style>
