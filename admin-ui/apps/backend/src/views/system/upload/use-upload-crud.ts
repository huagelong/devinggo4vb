import type { UploadListItem, UploadSearchFormModel } from './model';

import { reactive, ref } from 'vue';

import { message } from '#/adapter/tdesign';

import { downloadFileApi, getFileInfoApi } from '#/api/system/upload';

export function useUploadCrud() {
  const loading = ref(false);
  const tableData = ref<UploadListItem[]>([]);
  const selectedRowKeys = ref<Array<number | string>>([]);

  const searchForm = reactive<UploadSearchFormModel>({
    origin_name: '',
    mime_type: '',
    storage_mode: undefined,
    created_at: [],
  });

  const pagination = reactive({
    current: 1,
    pageSize: 10,
    total: 0,
  });

  // 模拟数据（实际应该调用API）
  const mockData: UploadListItem[] = [];

  const fetchTableData = async () => {
    loading.value = true;
    try {
      // TODO: 实际应该调用 upload API 获取文件列表
      // const params = {
      //   page: pagination.current,
      //   pageSize: pagination.pageSize,
      //   ...searchForm,
      // };
      // const res = await getUploadListApi(params);
      // tableData.value = res?.items || [];
      // pagination.total = res?.pageInfo?.total || 0;

      // 临时使用模拟数据
      tableData.value = mockData;
      pagination.total = mockData.length;
    } catch (error) {
      console.error(error);
      message.error('获取文件列表失败');
    } finally {
      loading.value = false;
    }
  };

  const handleSearch = () => {
    pagination.current = 1;
    fetchTableData();
  };

  const handleReset = () => {
    searchForm.origin_name = '';
    searchForm.mime_type = '';
    searchForm.storage_mode = undefined;
    searchForm.created_at = [];
    handleSearch();
  };

  const handlePageChange = (pageInfo: { current: number; pageSize: number }) => {
    pagination.current = pageInfo.current;
    pagination.pageSize = pageInfo.pageSize;
    fetchTableData();
  };

  const handleSelectChange = (val: Array<number | string>) => {
    selectedRowKeys.value = val;
  };

  const clearSelectedRowKeys = () => {
    selectedRowKeys.value = [];
  };

  const handleDownload = async (row: UploadListItem) => {
    try {
      const blob = await downloadFileApi({ id: row.id });
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = row.origin_name;
      link.click();
      window.URL.revokeObjectURL(url);
      message.success('下载成功');
    } catch (error) {
      console.error(error);
      message.error('下载失败');
    }
  };

  const handleView = async (row: UploadListItem) => {
    try {
      const info = await getFileInfoApi({ id: row.id });
      // 可以打开预览对话框
      console.log('File info:', info);
    } catch (error) {
      console.error(error);
      message.error('获取文件信息失败');
    }
  };

  return {
    loading,
    tableData,
    selectedRowKeys,
    searchForm,
    pagination,
    fetchTableData,
    handleSearch,
    handleReset,
    handlePageChange,
    handleSelectChange,
    clearSelectedRowKeys,
    handleDownload,
    handleView,
  };
}
