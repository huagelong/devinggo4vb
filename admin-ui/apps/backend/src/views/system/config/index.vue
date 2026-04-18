<script lang="ts" setup>
import type { ConfigApi } from '#/api/system/config';
import type { ConfigFormModel, ConfigGroup } from './model';

import { computed, onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';

import { message } from '#/adapter/tdesign';
import { logger } from '#/utils/logger';

import { AddIcon, InfoCircleIcon, SettingIcon } from 'tdesign-icons-vue-next';
import {
  Button,
  Form,
  FormItem,
  MessagePlugin,
  Space,
  Tabs,
  TabPanel,
} from 'tdesign-vue-next';

import {
  deleteConfigGroup,
  getConfigGroupList,
  getConfigList,
  updateConfigByKeys,
} from '#/api/system/config';

import ConfigFieldRenderer from './components/config-field-renderer.vue';
import ConfigFormModal from './components/config-form-modal.vue';
import ConfigGroupModal from './components/config-group-modal.vue';
import ConfigManageModal from './components/config-manage-modal.vue';
import type { ConfigFieldMeta } from './model';

const configGroupModalRef = ref<{ open: () => void }>();
const configFormModalRef = ref<{ open: () => void }>();
const configManageModalRef = ref<{ open: (groupId: number) => void }>();

const groupList = ref<ConfigGroup[]>([]);
const activeGroupKey = ref<number>();
const groupLoading = ref(false);
const configFieldsMap = reactive<Record<number, ConfigFieldMeta[]>>({});
const configFormMap = reactive<Record<number, ConfigFormModel>>({});

const hasGroups = computed(() => groupList.value.length > 0);

function normalizeOptions(data: unknown): ConfigFieldMeta['config_select_data'] {
  if (Array.isArray(data)) {
    return data as ConfigFieldMeta['config_select_data'];
  }
  if (typeof data === 'string' && data.trim()) {
    try {
      const parsed = JSON.parse(data);
      if (Array.isArray(parsed)) {
        return parsed as ConfigFieldMeta['config_select_data'];
      }
    } catch {
      return undefined;
    }
  }
  return undefined;
}

function normalizeSwitchValues(value: unknown) {
  if (typeof value === 'string') {
    return { checked: 'true', unchecked: 'false' };
  }
  if (typeof value === 'number') {
    return { checked: 1, unchecked: 0 };
  }
  return { checked: true, unchecked: false };
}

function parseCheckboxValue(rawValue: unknown): string[] {
  if (Array.isArray(rawValue)) {
    return rawValue.map((item) => String(item));
  }
  if (typeof rawValue === 'string') {
    try {
      const parsed = JSON.parse(rawValue);
      if (Array.isArray(parsed)) {
        return parsed.map((item) => String(item));
      }
    } catch {
      if (rawValue.includes(',')) {
        return rawValue.split(',').map((value) => value.trim());
      }
    }
  }
  return [];
}

function parseKeyValueInput(rawValue: unknown) {
  if (Array.isArray(rawValue)) {
    return rawValue;
  }
  if (typeof rawValue === 'string' && rawValue.trim()) {
    try {
      const parsed = JSON.parse(rawValue);
      if (Array.isArray(parsed)) {
        return parsed;
      }
    } catch {
      return [{ key: '', value: rawValue }];
    }
  }
  return [{ key: '', value: '' }];
}

async function fetchGroups() {
  groupLoading.value = true;
  try {
    const list = await getConfigGroupList();
    groupList.value = list;
    if (list.length > 0) {
      const firstGroup = list[0]!;
      activeGroupKey.value = firstGroup.id;
      await fetchGroupConfigs(firstGroup.id);
    } else {
      activeGroupKey.value = undefined;
    }
  } catch (error) {
    logger.error(error);
    message.error('配置分组加载失败，请稍后重试');
  } finally {
    groupLoading.value = false;
  }
}

async function fetchGroupConfigs(groupId: number) {
  try {
    const response = await getConfigList({
      group_id: groupId,
      orderBy: 'sort',
      orderType: 'asc',
    });
    const items = response.items ?? response.data ?? [];
    const fields: ConfigFieldMeta[] = [];
    const form: ConfigFormModel = {};
    items.forEach((item) => {
      const field: ConfigFieldMeta = {
        id: item.id,
        input_type: item.input_type,
        key: item.key,
        label: item.name,
        remark: item.remark,
        sort: item.sort,
        config_select_data: normalizeOptions(item.config_select_data),
      };
      let value: ConfigFormModel[string];
      if (item.input_type === 'switch') {
        field.switchValues = normalizeSwitchValues(item.value);
      }
      if (item.input_type === 'checkbox') {
        value = parseCheckboxValue(item.value);
      } else if (item.input_type === 'key-value') {
        value = parseKeyValueInput(item.value);
      } else {
        value = item.value as ConfigFormModel[string];
      }
      form[item.key] = value;
      fields.push(field);
    });
    configFieldsMap[groupId] = fields;
    configFormMap[groupId] = form;
  } catch (error) {
    logger.error(error);
    message.error('配置数据加载失败，请稍后重试');
  }
}

function handleTabChange(value: string | number) {
  const id = Number(value);
  if (Number.isNaN(id)) return;
  activeGroupKey.value = id;
  if (!configFieldsMap[id]) {
    void fetchGroupConfigs(id);
  }
}

function formatSubmitValue(
  field: ConfigFieldMeta,
  value: unknown,
): ConfigApi.UpdateByKeysPayload[string] {
  if (field.input_type === 'key-value') {
    try {
      return JSON.stringify(value ?? []);
    } catch {
      return JSON.stringify([]);
    }
  }
  return value as ConfigApi.UpdateByKeysPayload[string];
}

async function handleSubmit(groupId: number) {
  const currentForm = configFormMap[groupId];
  if (!currentForm) return;
  const fields = configFieldsMap[groupId] ?? [];
  const payload: ConfigApi.UpdateByKeysPayload = {};
  fields.forEach((field) => {
    payload[field.key] = formatSubmitValue(field, currentForm[field.key]);
  });

  try {
    await updateConfigByKeys(payload);
    MessagePlugin.success('配置更新成功');
    await fetchGroupConfigs(groupId);
  } catch (error) {
    logger.error(error);
    MessagePlugin.error('配置更新失败，请稍后重试');
  }
}

async function handleDeleteGroup(groupId: number) {
  const group = groupList.value.find((item) => item.id === groupId);
  if (!group) return;
  if (groupId === 1 || groupId === 2) {
    MessagePlugin.info('系统默认分组不可删除');
    return;
  }
  try {
    await deleteConfigGroup({ id: groupId });
    MessagePlugin.success('分组删除成功');
    await fetchGroups();
  } catch (error) {
    logger.error(error);
    MessagePlugin.error('分组删除失败，请稍后重试');
  }
}

onMounted(() => {
  void fetchGroups();
});
</script>

<template>
  <Page auto-content-height>
    <div class="flex h-full flex-col gap-4">
      <div class="flex flex-col gap-4 lg:flex-row">
        <div class="rounded-md bg-white p-4 lg:w-1/2">
          <div class="mb-3 flex items-center justify-between">
            <Space>
              <Button
                theme="primary"
                variant="outline"
                @click="configGroupModalRef?.open()"
              >
                <template #icon><AddIcon /></template>
                新增分组
              </Button>
              <Button
                theme="default"
                variant="outline"
                :disabled="!activeGroupKey"
                @click="
                  activeGroupKey &&
                  configManageModalRef?.open(Number(activeGroupKey))
                "
              >
                <template #icon><SettingIcon /></template>
                管理配置
              </Button>
              <Button
                v-if="activeGroupKey && activeGroupKey > 2"
                theme="danger"
                variant="outline"
                @click="handleDeleteGroup(activeGroupKey)"
              >
                删除分组
              </Button>
            </Space>
          </div>
          <div v-if="hasGroups" class="min-h-[400px]">
            <Tabs
              v-model="activeGroupKey"
              :loading="groupLoading"
              class="config-tabs"
              @change="handleTabChange"
            >
              <TabPanel
                v-for="group in groupList"
                :key="group.id"
                :value="group.id"
                :label="group.name"
              >
                <Form
                  v-if="configFormMap[group.id]"
                  :data="configFormMap[group.id]!"
                  label-width="120px"
                  colon
                >
                  <div class="flex flex-col gap-4">
                    <FormItem
                      v-for="field in configFieldsMap[group.id] ?? []"
                      :key="field.id"
                      :label="field.label"
                      :name="field.key"
                      class="border-b border-gray-50 pb-3"
                    >
                      <ConfigFieldRenderer
                        v-model="configFormMap[group.id]![field.key]"
                        :field="field"
                      />
                      <div
                        v-if="field.remark"
                        class="mt-1 text-xs text-gray-500"
                      >
                        {{ field.remark }}
                      </div>
                    </FormItem>
                  </div>
                  <div class="mt-4 flex justify-end">
                    <Button theme="primary" @click="handleSubmit(group.id)">
                      保存配置
                    </Button>
                  </div>
                </Form>
                <div
                  v-else
                  class="flex min-h-[200px] items-center justify-center text-gray-500"
                >
                  <InfoCircleIcon class="mr-2" /> 配置加载中...
                </div>
              </TabPanel>
            </Tabs>
          </div>
          <div
            v-else
            class="flex min-h-[200px] items-center justify-center text-gray-500"
          >
            暂无配置分组，先新增一个吧。
          </div>
        </div>

        <div class="rounded-md bg-white p-4 lg:w-1/2">
          <div class="mb-4 flex items-center justify-between">
            <h3 class="text-lg font-semibold">新增配置</h3>
            <Button theme="primary" @click="configFormModalRef?.open()">
              <template #icon><AddIcon /></template>
              新增
            </Button>
          </div>
          <p class="text-sm text-gray-500">
            新增配置后可在左侧对应分组中实时看到并维护值。
          </p>
        </div>
      </div>
    </div>

    <ConfigGroupModal ref="configGroupModalRef" @success="fetchGroups" />
    <ConfigFormModal ref="configFormModalRef" @success="fetchGroups" />
    <ConfigManageModal ref="configManageModalRef" />
  </Page>
</template>
