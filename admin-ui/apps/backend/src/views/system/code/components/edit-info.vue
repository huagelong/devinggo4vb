<script lang="ts" setup>
import type { FieldConfigRow } from '../model';

import { ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { message } from '#/adapter/tdesign';
import { logger } from '#/utils/logger';

import { readTable, updateCode } from '#/api/system/generate';

import {
  Button,
  Checkbox,
  CheckboxGroup,
  Form,
  FormItem,
  Input,
  Select,
  Switch,
  TabPanel,
  Table,
  Tabs,
} from 'tdesign-vue-next';

import {
  componentTypeOptions,
  generateTypeOptions,
  menuButtonOptions,
  queryTypeOptions,
  tplTypeOptions,
  viewTypeOptions,
} from '../schemas';

const emit = defineEmits<{
  success: [];
}>();

const loading = ref(false);
const submitLoading = ref(false);
const activeTab = ref('base_config');

// 表基本信息
const tableName = ref('');
const tableComment = ref('');
const remark = ref('');

// 生成配置
const moduleName = ref('');
const type = ref<'single' | 'tree'>('single');
const menuName = ref('');
const componentType = ref<number>(1);
const tplType = ref('default');

// 树表配置
const treeId = ref('');
const treeParentId = ref('');
const treeName = ref('');

// Tag页配置
const tagId = ref('');
const tagName = ref('');
const tagViewName = ref('');

// 字段配置
const fieldColumns = [
  { colKey: 'sort', title: '排序', width: 70 },
  { colKey: 'column_name', title: '字段名称', width: 150 },
  { colKey: 'column_comment', title: '字段描述', width: 150 },
  { colKey: 'column_type', title: '物理类型', width: 120 },
  { colKey: 'is_required', title: '必填', width: 70 },
  { colKey: 'is_insert', title: '新增', width: 70 },
  { colKey: 'is_edit', title: '编辑', width: 70 },
  { colKey: 'is_list', title: '列表', width: 70 },
  { colKey: 'is_query', title: '查询', width: 70 },
  { colKey: 'is_sort', title: '排序', width: 70 },
  { colKey: 'query_type', title: '查询方式', width: 100 },
  { colKey: 'view_type', title: '控件类型', width: 120 },
  { colKey: 'action', title: '操作', width: 100 },
];

const fieldList = ref<FieldConfigRow[]>([]);

// 菜单配置
const menuButtons = ref<string[]>(['save', 'update', 'read', 'delete']);

async function open(id: number) {
  loading.value = true;
  try {
    await readTable(id);
    // 模拟数据加载
    tableName.value = '';
    tableComment.value = '';
    fieldList.value = [];
    activeTab.value = 'base_config';
    modalApi.setState({ title: '编辑代码生成信息' });
    modalApi.open();
  } catch (error) {
    logger.error(error);
    message.error($t('common.tableInfoFailed'));
  } finally {
    loading.value = false;
  }
}

function handleEditField(row: FieldConfigRow) {
  void row;
  message.info($t('common.advancedConfigWIP'));
}

function handleAllChecked(checked: boolean, key: keyof FieldConfigRow) {
  if (checked) {
    fieldList.value.forEach((f) => {
      (f as any)[key] = 1;
    });
  } else {
    fieldList.value.forEach((f) => {
      (f as any)[key] = 2;
    });
  }
}

async function handleSubmit() {
  submitLoading.value = true;
  try {
    const payload: any = {
      id: 0, // 需要从列表中获取
      type: type.value,
      menu_buttons: menuButtons.value,
      fields: fieldList.value,
    };
    await updateCode(payload);
    message.success($t('common.saveSuccess'));
    emit('success');
    modalApi.close();
  } catch (error) {
    logger.error(error);
    message.error($t('common.saveFailed'));
  } finally {
    submitLoading.value = false;
  }
}

const [Modal, modalApi] = useVbenModal({
  onConfirm: handleSubmit,
  class: 'w-[1100px]',
});

defineExpose({ open });
</script>

<template>
  <Modal>
    <div class="flex flex-col gap-4">
      <Tabs v-model:value="activeTab">
        <!-- 基础配置 -->
        <TabPanel value="base_config" label="配置信息">
          <Form :label-width="120" colon>
            <div class="grid grid-cols-2 gap-x-4">
              <FormItem label="表名称">
                <Input v-model="tableName" disabled />
              </FormItem>
              <FormItem label="表描述" name="table_comment">
                <Input v-model="tableComment" placeholder="请输入表描述" />
              </FormItem>
              <FormItem label="备注">
                <Input v-model="remark" placeholder="请输入备注" />
              </FormItem>
              <FormItem label="所属模块" name="module_name">
                <Input v-model="moduleName" placeholder="请输入所属模块" />
              </FormItem>
              <FormItem label="生成类型" name="type">
                <Select v-model="type" :options="generateTypeOptions" />
              </FormItem>
              <FormItem label="菜单名称" name="menu_name">
                <Input v-model="menuName" placeholder="请输入菜单名称" />
              </FormItem>
              <FormItem label="组件样式" name="component_type">
                <Select v-model="componentType" :options="componentTypeOptions" />
              </FormItem>
              <FormItem label="模板类型" name="tpl_type">
                <Select v-model="tplType" :options="tplTypeOptions" />
              </FormItem>
            </div>

            <!-- 树表配置 -->
            <template v-if="type === 'tree'">
              <div class="mt-4 border-t pt-4">
                <div class="mb-2 text-sm font-medium text-gray-500">树表配置</div>
                <div class="grid grid-cols-3 gap-x-4">
                  <FormItem label="树主ID" name="tree_id">
                    <Input v-model="treeId" placeholder="请输入树主ID字段" />
                  </FormItem>
                  <FormItem label="树父ID" name="tree_parent_id">
                    <Input v-model="treeParentId" placeholder="请输入树父ID字段" />
                  </FormItem>
                  <FormItem label="树名称" name="tree_name">
                    <Input v-model="treeName" placeholder="请输入树名称字段" />
                  </FormItem>
                </div>
              </div>
            </template>

            <!-- Tag页配置 -->
            <template v-if="componentType === 3">
              <div class="mt-4 border-t pt-4">
                <div class="mb-2 text-sm font-medium text-gray-500">Tag页配置</div>
                <div class="grid grid-cols-3 gap-x-4">
                  <FormItem label="标签页ID" name="tag_id">
                    <Input v-model="tagId" placeholder="请输入标签页ID" />
                  </FormItem>
                  <FormItem label="标签页名称" name="tag_name">
                    <Input v-model="tagName" placeholder="请输入标签页名称" />
                  </FormItem>
                  <FormItem label="标签显示字段" name="tag_view_name">
                    <Input v-model="tagViewName" placeholder="请输入标签显示字段" />
                  </FormItem>
                </div>
              </div>
            </template>
          </Form>
        </TabPanel>

        <!-- 字段配置 -->
        <TabPanel value="field_config" label="字段配置">
          <div class="mb-2 flex items-center gap-2">
            <Checkbox
              :checked="fieldList.every((f) => f.is_required === 1)"
              @change="(checked: boolean) => handleAllChecked(checked, 'is_required')"
            >
              全选必填
            </Checkbox>
            <Checkbox
              :checked="fieldList.every((f) => f.is_insert === 1)"
              @change="(checked: boolean) => handleAllChecked(checked, 'is_insert')"
            >
              全选新增
            </Checkbox>
            <Checkbox
              :checked="fieldList.every((f) => f.is_edit === 1)"
              @change="(checked: boolean) => handleAllChecked(checked, 'is_edit')"
            >
              全选编辑
            </Checkbox>
            <Checkbox
              :checked="fieldList.every((f) => f.is_list === 1)"
              @change="(checked: boolean) => handleAllChecked(checked, 'is_list')"
            >
              全选列表
            </Checkbox>
            <Checkbox
              :checked="fieldList.every((f) => f.is_query === 1)"
              @change="(checked: boolean) => handleAllChecked(checked, 'is_query')"
            >
              全选查询
            </Checkbox>
          </div>

          <Table
            :columns="fieldColumns"
            :data="fieldList"
            row-key="column_name"
            hover
            stripe
          >
            <template #is_required="{ row }">
              <Switch v-model="row.is_required" :value="row.is_required === 1" />
            </template>
            <template #is_insert="{ row }">
              <Switch v-model="row.is_insert" :value="row.is_insert === 1" />
            </template>
            <template #is_edit="{ row }">
              <Switch v-model="row.is_edit" :value="row.is_edit === 1" />
            </template>
            <template #is_list="{ row }">
              <Switch v-model="row.is_list" :value="row.is_list === 1" />
            </template>
            <template #is_query="{ row }">
              <Switch v-model="row.is_query" :value="row.is_query === 1" />
            </template>
            <template #is_sort="{ row }">
              <Switch v-model="row.is_sort" :value="row.is_sort === 1" />
            </template>
            <template #query_type="{ row }">
              <Select
                v-model="row.query_type"
                :options="queryTypeOptions"
                size="small"
              />
            </template>
            <template #view_type="{ row }">
              <Select
                v-model="row.view_type"
                :options="viewTypeOptions"
                size="small"
              />
            </template>
            <template #action="{ row }">
              <Button size="small" theme="primary" @click="handleEditField(row)">
                配置
              </Button>
            </template>
          </Table>
        </TabPanel>

        <!-- 菜单配置 -->
        <TabPanel value="menu_config" label="菜单配置">
          <div class="py-4">
            <CheckboxGroup v-model="menuButtons" :options="menuButtonOptions" />
          </div>
        </TabPanel>
      </Tabs>
    </div>
  </Modal>
</template>
