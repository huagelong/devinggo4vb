import { requestClient } from '#/api/request';

/**
 * 仪表板统计数据结构
 */
export interface DashboardStatistics {
  userStats: { total: number; new: number };
  attachmentStats: { total: number; new: number };
  loginStats: { total: number; new: number };
  operationStats: { total: number; new: number };
}

/**
 * 获取仪表板统计数据
 */
export async function getDashboardStatisticsApi() {
  return requestClient.get<DashboardStatistics>('/system/dashboard/statistics');
}

/**
 * 登录图表数据结构
 */
export interface DashboardLoginChart {
  xAxis: string[];
  chartsData: number[];
}

/**
 * 获取仪表板登录图表数据
 */
export async function getDashboardLoginChartApi(days: number = 10) {
  return requestClient.get<DashboardLoginChart>('/system/dashboard/loginChart', {
    params: { days },
  });
}
