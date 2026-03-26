import type { Ref } from 'vue';

import { ref } from 'vue';

import { DialogPlugin, MessagePlugin } from 'tdesign-vue-next';

import {
  changeUserStatus,
  clearUserCache,
  deleteUser,
  downloadUserImportTemplate,
  exportUserList,
  importUserFile,
  realDeleteUser,
  recoveryUser,
  resetPassword,
  setHomePage,
} from '#/api/system/user';

interface UseUserActionsOptions {
  buildRequestParams: (includePagination?: boolean) => Record<string, any>;
  clearSelectedRowKeys: () => void;
  fetchTableData: () => void;
  isRecycleBin: Ref<boolean>;
  selectedRowKeys: Ref<Array<number | string>>;
}

export function useUserActions(options: UseUserActionsOptions) {
  const importInputRef = ref<HTMLInputElement>();
  const importLoading = ref(false);
  const exportLoading = ref(false);
  const templateLoading = ref(false);

  const setHomePageVisible = ref(false);
  const setHomePageLoading = ref(false);
  const selectedHomePage = ref('');
  const selectedHomePageUserId = ref<null | number>(null);

  function isSuperAdmin(row: any) {
    return Number(row?.id) === 1;
  }

  function toIds(keys: Array<number | string>) {
    return keys.map((key) => Number(key));
  }

  function getFileNameFromDisposition(disposition?: string) {
    if (!disposition) return '';

    const utf8Match = disposition.match(/filename\*=UTF-8''([^;]+)/i);
    if (utf8Match?.[1]) {
      return decodeURIComponent(utf8Match[1]);
    }

    const asciiMatch = disposition.match(/filename="?([^"]+)"?/i);
    return asciiMatch?.[1] ?? '';
  }

  function saveBlobFile(blob: Blob, fileName: string) {
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = fileName;
    link.style.display = 'none';
    document.body.append(link);
    link.click();
    link.remove();
    URL.revokeObjectURL(url);
  }

  async function handleDelete(row: any) {
    if (isSuperAdmin(row)) {
      MessagePlugin.warning('超级管理员不可删除');
      return;
    }

    try {
      await (options.isRecycleBin.value
        ? realDeleteUser([row.id])
        : deleteUser([row.id]));
      MessagePlugin.success('删除成功');
      options.fetchTableData();
    } catch (error) {
      console.error(error);
    }
  }

  async function handleBatchDelete() {
    if (options.selectedRowKeys.value.length === 0) {
      MessagePlugin.warning('请选择需要操作的数据');
      return;
    }

    const ids = toIds(options.selectedRowKeys.value);
    if (ids.some((id) => id === 1)) {
      MessagePlugin.warning('超级管理员不可删除');
      return;
    }

    try {
      await (options.isRecycleBin.value ? realDeleteUser(ids) : deleteUser(ids));
      MessagePlugin.success('操作成功');
      options.clearSelectedRowKeys();
      options.fetchTableData();
    } catch (error) {
      console.error(error);
    }
  }

  async function handleRecovery(row: any) {
    try {
      await recoveryUser([row.id]);
      MessagePlugin.success('恢复成功');
      options.fetchTableData();
    } catch (error) {
      console.error(error);
    }
  }

  async function handleBatchRecovery() {
    if (options.selectedRowKeys.value.length === 0) {
      MessagePlugin.warning('请选择需要操作的数据');
      return;
    }

    const ids = toIds(options.selectedRowKeys.value);
    if (ids.some((id) => id === 1)) {
      MessagePlugin.warning('超级管理员不可恢复');
      return;
    }

    try {
      await recoveryUser(ids);
      MessagePlugin.success('恢复成功');
      options.clearSelectedRowKeys();
      options.fetchTableData();
    } catch (error) {
      console.error(error);
    }
  }

  async function handleStatusChange(row: any, checked: boolean) {
    if (isSuperAdmin(row)) {
      MessagePlugin.warning('超级管理员不可禁用');
      return;
    }

    const status = checked ? 1 : 2;
    try {
      await changeUserStatus({ id: row.id, status });
      MessagePlugin.success('更新状态成功');
      options.fetchTableData();
    } catch (error) {
      console.error(error);
    }
  }

  async function handleResetPassword(row: any) {
    try {
      await resetPassword({ id: row.id });
      MessagePlugin.success('密码重置成功');
    } catch (error) {
      console.error(error);
    }
  }

  async function handleClearCache(row: any) {
    try {
      await clearUserCache({ id: row.id });
      MessagePlugin.success('清除缓存成功');
    } catch (error) {
      console.error(error);
    }
  }

  function triggerImport() {
    importInputRef.value?.click();
  }

  async function handleImportChange(event: Event) {
    const input = event.target as HTMLInputElement;
    const file = input.files?.[0];
    if (!file) return;

    importLoading.value = true;
    try {
      await importUserFile(file);
      MessagePlugin.success('导入成功');
      options.fetchTableData();
    } catch (error) {
      console.error(error);
    } finally {
      importLoading.value = false;
      input.value = '';
    }
  }

  async function handleExport() {
    exportLoading.value = true;
    try {
      const response: any = await exportUserList(options.buildRequestParams(false));
      const disposition =
        response?.headers?.get?.('content-disposition') ||
        response?.headers?.['content-disposition'];
      const fileName = getFileNameFromDisposition(disposition) || '用户列表.xlsx';
      saveBlobFile(response?.data, fileName);
      MessagePlugin.success('导出成功');
    } catch (error) {
      console.error(error);
    } finally {
      exportLoading.value = false;
    }
  }

  async function handleDownloadTemplate() {
    templateLoading.value = true;
    try {
      const response: any = await downloadUserImportTemplate();
      const disposition =
        response?.headers?.get?.('content-disposition') ||
        response?.headers?.['content-disposition'];
      const fileName = getFileNameFromDisposition(disposition) || '用户导入模板.xlsx';
      saveBlobFile(response?.data, fileName);
      MessagePlugin.success('模板下载成功');
    } catch (error) {
      console.error(error);
    } finally {
      templateLoading.value = false;
    }
  }

  function openSetHomePageDialog(row: any) {
    if (isSuperAdmin(row)) {
      MessagePlugin.warning('超级管理员不可设置首页');
      return;
    }
    selectedHomePageUserId.value = Number(row.id);
    selectedHomePage.value = row.dashboard || '';
    setHomePageVisible.value = true;
  }

  async function handleSetHomePage() {
    if (!selectedHomePageUserId.value) {
      MessagePlugin.warning('用户信息无效');
      return;
    }

    if (!selectedHomePage.value) {
      MessagePlugin.warning('请选择用户首页');
      return;
    }

    setHomePageLoading.value = true;
    try {
      await setHomePage({
        dashboard: selectedHomePage.value,
        id: selectedHomePageUserId.value,
      });
      MessagePlugin.success('设置首页成功');
      setHomePageVisible.value = false;
      options.fetchTableData();
    } catch (error) {
      console.error(error);
    } finally {
      setHomePageLoading.value = false;
    }
  }

  function handleActionDropdownClick(data: any, row: any) {
    if (data.value === 'reset_password') {
      const dialog = DialogPlugin.confirm({
        body: '确认重置该用户密码吗？',
        header: '提示',
        onClose: () => dialog.hide(),
        onConfirm: () => {
          handleResetPassword(row);
          dialog.hide();
        },
      });
      return;
    }

    if (data.value === 'clear_cache') {
      const dialog = DialogPlugin.confirm({
        body: '确认更新该用户缓存吗？',
        header: '提示',
        onClose: () => dialog.hide(),
        onConfirm: () => {
          handleClearCache(row);
          dialog.hide();
        },
      });
      return;
    }

    if (data.value === 'set_homepage') {
      openSetHomePageDialog(row);
    }
  }

  return {
    exportLoading,
    handleActionDropdownClick,
    handleBatchDelete,
    handleBatchRecovery,
    handleClearCache,
    handleDelete,
    handleDownloadTemplate,
    handleExport,
    handleImportChange,
    handleRecovery,
    handleSetHomePage,
    handleStatusChange,
    importInputRef,
    importLoading,
    isSuperAdmin,
    selectedHomePage,
    setHomePageLoading,
    setHomePageVisible,
    templateLoading,
    triggerImport,
  };
}
