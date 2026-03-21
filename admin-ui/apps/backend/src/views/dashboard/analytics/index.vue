<script lang="ts" setup>
import type { AnalysisOverviewItem } from '@vben/common-ui';
import type { TabOption } from '@vben/types';

import { markRaw, onMounted, ref } from 'vue';

import {
  AnalysisChartsTabs,
  AnalysisOverview,
} from '@vben/common-ui';
import {
  SvgBellIcon,
  SvgCakeIcon,
  SvgCardIcon,
  SvgDownloadIcon,
} from '@vben/icons';

import { getDashboardStatisticsApi } from '#/api/core/dashboard';

import AnalyticsTrends from './analytics-trends.vue';

const overviewItems = ref<AnalysisOverviewItem[]>([
  {
    icon: markRaw(SvgCardIcon),
    title: '用户总人数',
    totalTitle: '用户新增数',
    totalValue: 0,
    value: 0,
  },
  {
    icon: markRaw(SvgDownloadIcon),
    title: '附件总数',
    totalTitle: '附件新增数',
    totalValue: 0,
    value: 0,
  },
  {
    icon: markRaw(SvgCakeIcon),
    title: '总登录数',
    totalTitle: '新增登录数',
    totalValue: 0,
    value: 0,
  },
  {
    icon: markRaw(SvgBellIcon),
    title: '总操作数',
    totalTitle: '新增操作数',
    totalValue: 0,
    value: 0,
  },
]);

const chartTabs: TabOption[] = [
  {
    label: '登录统计',
    value: 'trends',
  },
];

async function initData() {
  try {
    const data = await getDashboardStatisticsApi();
    overviewItems.value = [
      {
        icon: markRaw(SvgCardIcon),
        title: '用户数',
        totalTitle: '总用户数',
        totalValue: data.userStats?.total || 0,
        value: data.userStats?.new || 0,
      },
      {
        icon: markRaw(SvgDownloadIcon),
        title: '附件数',
        totalTitle: '总附件数',
        totalValue: data.attachmentStats?.total || 0,
        value: data.attachmentStats?.new || 0,
      },
      {
        icon: markRaw(SvgCakeIcon),
        title: '登录数',
        totalTitle: '总登录数',
        totalValue: data.loginStats?.total || 0,
        value: data.loginStats?.new || 0,
      },
      {
        icon: markRaw(SvgBellIcon),
        title: '操作数',
        totalTitle: '总操作数',
        totalValue: data.operationStats?.total || 0,
        value: data.operationStats?.new || 0,
      },
    ];
  } catch (error) {
    console.error('Failed to load dashboard statistics', error);
  }
}

onMounted(() => {
  initData();
});
</script>

<template>
  <div class="p-5">
    <AnalysisOverview :items="overviewItems" />
    <AnalysisChartsTabs :tabs="chartTabs" class="mt-5">
      <template #trends>
        <AnalyticsTrends />
      </template>
    </AnalysisChartsTabs>
  </div>
</template>
