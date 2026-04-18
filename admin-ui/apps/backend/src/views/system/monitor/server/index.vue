<script lang="ts" setup>
import type { MonitorApi } from '#/api/system/monitor';

import { onMounted, onUnmounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { message } from '#/adapter/tdesign';
import { getServerInfo } from '#/api/system/monitor';
import { logger } from '#/utils/logger';

import {
  CpuIcon,
  DesktopIcon,
  ServerIcon,
} from 'tdesign-icons-vue-next';
import {
  Card,
  Col,
  Progress,
  Row,
  Space,
  Tag,
} from 'tdesign-vue-next';

defineOptions({ name: 'SystemServer' });

const loading = ref(false);
const serverInfo = ref<MonitorApi.ServerInfoResponse | null>(null);
const hasServerApi = ref(true);
const refreshTimer = ref<ReturnType<typeof setInterval>>();

function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return `${Number.parseFloat((bytes / k ** i).toFixed(2))} ${sizes[i]}`;
}

function formatUptime(seconds: number): string {
  const days = Math.floor(seconds / 86400);
  const hours = Math.floor((seconds % 86400) / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const parts: string[] = [];
  if (days > 0) parts.push(`${days}天`);
  if (hours > 0) parts.push(`${hours}小时`);
  if (minutes > 0) parts.push(`${minutes}分钟`);
  return parts.join(' ') || `${seconds}秒`;
}

async function fetchServerInfo() {
  if (!hasServerApi.value) return;

  loading.value = true;
  try {
    const response = await getServerInfo();
    serverInfo.value = response;
  } catch (error: unknown) {
    if ((error as { response?: { status?: number } })?.response?.status === 404) {
      hasServerApi.value = false;
      message.info($t('common.serverMonitorNotAvailable'));
    } else {
      logger.error(error);
      message.error($t('common.serverInfoFailed'));
    }
  } finally {
    loading.value = false;
  }
}

function startAutoRefresh() {
  refreshTimer.value = setInterval(() => {
    void fetchServerInfo();
  }, 10000);
}

function stopAutoRefresh() {
  if (refreshTimer.value) {
    clearInterval(refreshTimer.value);
    refreshTimer.value = undefined;
  }
}

onMounted(() => {
  void fetchServerInfo();
  startAutoRefresh();
});

onUnmounted(() => {
  stopAutoRefresh();
});
</script>

<template>
  <Page auto-content-height>
    <div class="flex flex-col gap-3">
      <!-- API Not Available Notice -->
      <Card v-if="!hasServerApi">
        <div class="flex flex-col items-center justify-center py-12 text-gray-400">
          <ServerIcon class="mb-4 text-6xl" />
          <div class="text-lg">服务器监控接口暂未开放</div>
          <div class="text-sm">后端 API (/system/server/monitor) 待实现后自动启用</div>
        </div>
      </Card>

      <!-- Server Info Content -->
      <template v-else>
        <!-- System Overview -->
        <Card title="系统概览">
          <template #actions>
            <Tag
              :theme="loading ? 'default' : 'success'"
              variant="light"
            >
              {{ loading ? '刷新中...' : '实时监控中' }}
            </Tag>
          </template>
          <Row :gutter="24">
            <Col :span="6">
              <div class="mb-2 flex items-center gap-2 text-sm text-gray-500">
                <DesktopIcon />
                操作系统
              </div>
              <div class="text-base">{{ serverInfo?.os || '-' }}</div>
            </Col>
            <Col :span="6">
              <div class="mb-2 flex items-center gap-2 text-sm text-gray-500">
                <ServerIcon />
                架构
              </div>
              <div class="text-base">{{ serverInfo?.arch || '-' }}</div>
            </Col>
            <Col :span="6">
              <div class="mb-2 text-sm text-gray-500">主机名</div>
              <div class="text-base">{{ serverInfo?.hostname || '-' }}</div>
            </Col>
            <Col :span="6">
              <div class="mb-2 text-sm text-gray-500">运行时间</div>
              <div class="text-base">
                {{ serverInfo?.uptime ? formatUptime(serverInfo.uptime) : '-' }}
              </div>
            </Col>
          </Row>
          <Row :gutter="24" class="mt-4">
            <Col :span="6">
              <div class="mb-2 text-sm text-gray-500">服务器时间</div>
              <div class="text-base">{{ serverInfo?.server_time || '-' }}</div>
            </Col>
            <Col :span="6">
              <div class="mb-2 text-sm text-gray-500">Go 版本</div>
              <div class="text-base">{{ serverInfo?.go_runtime?.go_version || '-' }}</div>
            </Col>
            <Col :span="6">
              <div class="mb-2 text-sm text-gray-500">Goroutines</div>
              <div class="text-base">{{ serverInfo?.go_runtime?.goroutines || '-' }}</div>
            </Col>
            <Col :span="6">
              <div class="mb-2 text-sm text-gray-500">GC 统计</div>
              <div class="text-base">{{ serverInfo?.go_runtime?.gc_stats || '-' }}</div>
            </Col>
          </Row>
        </Card>

        <!-- CPU & Memory -->
        <Row :gutter="16">
          <Col :span="12">
            <Card title="CPU 使用率">
              <div class="flex flex-col items-center py-4">
                <Progress
                  :percentage="serverInfo?.cpu?.usage ?? 0"
                  :size="'large'"
                  :theme="((
                    serverInfo?.cpu?.usage ?? 0
                  ) > 80
                    ? 'danger'
                    : (serverInfo?.cpu?.usage ?? 0) > 60
                      ? 'warning'
                      : 'primary') as any"
                  :label="`${(serverInfo?.cpu?.usage ?? 0).toFixed(1)}%`"
                />
                <div class="mt-4 text-sm text-gray-500">
                  <Space>
                    <span>
                      <CpuIcon />
                      核心数: {{ serverInfo?.cpu?.num ?? '-' }}
                    </span>
                    <span>型号: {{ serverInfo?.cpu?.model ?? '-' }}</span>
                  </Space>
                </div>
              </div>
            </Card>
          </Col>
          <Col :span="12">
            <Card title="内存使用率">
              <div class="flex flex-col items-center py-4">
                <Progress
                  :percentage="serverInfo?.memory?.usage ?? 0"
                  :size="'large'"
                  :theme="((
                    serverInfo?.memory?.usage ?? 0
                  ) > 80
                    ? 'danger'
                    : (serverInfo?.memory?.usage ?? 0) > 60
                      ? 'warning'
                      : 'primary') as any"
                  :label="`${(serverInfo?.memory?.usage ?? 0).toFixed(1)}%`"
                />
                <div class="mt-4 text-sm text-gray-500">
                  <Space>
                    <span>
                      已用: {{ formatBytes(serverInfo?.memory?.used ?? 0) }}
                    </span>
                    <span>
                      总量: {{ formatBytes(serverInfo?.memory?.total ?? 0) }}
                    </span>
                    <span>
                      可用: {{ formatBytes(serverInfo?.memory?.free ?? 0) }}
                    </span>
                  </Space>
                </div>
              </div>
            </Card>
          </Col>
        </Row>

        <!-- Disk Info -->
        <Card title="磁盘信息">
          <Row :gutter="16">
            <Col
              v-for="(disk, index) in serverInfo?.disks ?? []"
              :key="index"
              :span="8"
            >
              <Card :bordered="false" class="bg-gray-50">
                <div class="mb-2 text-sm font-medium text-gray-600">
                  {{ disk.mount_point }}
                </div>
                <Progress
                  :percentage="disk.usage"
                  :theme="(disk.usage > 80
                      ? 'danger'
                      : disk.usage > 60
                        ? 'warning'
                        : 'primary') as any"
                  :label="`${disk.usage.toFixed(1)}%`"
                />
                <div class="mt-2 text-xs text-gray-500">
                  <Space>
                    <span>文件系统: {{ disk.file_system }}</span>
                    <span>
                      总量: {{ formatBytes(disk.total) }}
                    </span>
                    <span>
                      已用: {{ formatBytes(disk.used) }}
                    </span>
                  </Space>
                </div>
              </Card>
            </Col>
          </Row>
          <div
            v-if="!serverInfo?.disks?.length"
            class="py-8 text-center text-gray-400"
          >
            暂无磁盘信息
          </div>
        </Card>

        <!-- Go Runtime -->
        <Card title="Go 运行时">
          <Row :gutter="24">
            <Col :span="6">
              <div class="mb-2 text-sm text-gray-500">Heap 分配</div>
              <div class="text-base">
                {{ formatBytes(serverInfo?.go_runtime?.heap_alloc ?? 0) }}
              </div>
            </Col>
            <Col :span="6">
              <div class="mb-2 text-sm text-gray-500">Heap 系统内存</div>
              <div class="text-base">
                {{ formatBytes(serverInfo?.go_runtime?.heap_sys ?? 0) }}
              </div>
            </Col>
            <Col :span="6">
              <div class="mb-2 text-sm text-gray-500">栈使用</div>
              <div class="text-base">
                {{ formatBytes(serverInfo?.go_runtime?.stack_in_use ?? 0) }}
              </div>
            </Col>
          </Row>
        </Card>
      </template>
    </div>
  </Page>
</template>
