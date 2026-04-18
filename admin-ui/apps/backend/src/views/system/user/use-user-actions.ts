import type { Ref } from 'vue';
import type { IdType } from '#/types/common';

import { ref } from 'vue';

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
import { dialog, message } from '#/adapter/tdesign';
import { downloadResponseBlob } from '#/utils/download';

import type { UserActionDropdownItem, UserListItem } from './model';

type UserActionRow = Pick<UserListItem, 'id'> & Partial<UserListItem>;

interface UseUserActionsOptions {
  buildRequestParams: (includePagination?: boolean) => Record<string, unknown>;
  clearSelectedRowKeys: () => void;
  fetchTableData: () => void;
  isRecycleBin: Ref<boolean>;
  selectedRowKeys: Ref<IdType[]>;
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

  function isSuperAdmin(row: UserActionRow) {
    return Number(row?.id) === 1;
  }

  function toIds(keys: IdType[]) {
    return keys.map((key) => Number(key));
  }

  async function handleDelete(row: UserActionRow) {
    if (isSuperAdmin(row)) {
      message.warning('超级管理员不可删除');
      return;
    }

    try {
      await (options.isRecycleBin.value
        ? realDeleteUser([row.id])
        : deleteUser([row.id]));
      message.success('删除成功');
      options.fetchTableData();
    } catch (error) {
      logger.error(error);
      message.error('删除失败，请稍后重试');
    }
  }

  async function handleBatchDelete() {
    if (options.selectedRowKeys.value.length === 0) {
      message.warning('请选择需要操作的数据');
      return;
    }

    const ids = toIds(options.selectedRowKeys.value);
    if (ids.some((id) => id === 1)) {
      message.warning('超级管理员不可删除');
      return;
    }

    try {
      await (options.isRecycleBin.value ? realDeleteUser(ids) : deleteUser(ids));
      message.success('操作成功');
      options.clearSelectedRowKeys();
      options.fetchTableData();
    } catch (error) {
      logger.error(error);
      message.error('批量删除失败，请稍后重试');
    }
  }

  async function handleRecovery(row: UserActionRow) {
    try {
      await recoveryUser([row.id]);
      message.success('恢复成功');
      options.fetchTableData();
    } catch (error) {
      logger.error(error);
      message.error('恢复失败，请稍后重试');
    }
  }

  async function handleBatchRecovery() {
    if (options.selectedRowKeys.value.length === 0) {
      message.warning('请选择需要操作的数据');
      return;
    }

    const ids = toIds(options.selectedRowKeys.value);
    if (ids.some((id) => id === 1)) {
      message.warning('超级管理员不可恢复');
      return;
    }

    try {
      await recoveryUser(ids);
      message.success('恢复成功');
      options.clearSelectedRowKeys();
      options.fetchTableData();
    } catch (error) {
      logger.error(error);
      message.error('批量恢复失败，请稍后重试');
    }
  }

  async function handleStatusChange(row: UserActionRow, checked: boolean) {
    if (isSuperAdmin(row)) {
      message.warning('超级管理员不可禁用');
      return;
    }

    const status = checked ? 1 : 2;
    try {
      await changeUserStatus({ id: row.id, status });
      message.success('更新状态成功');
      options.fetchTableData();
    } catch (error) {
      logger.error(error);
      message.error('更新状态失败，请稍后重试');
    }
  }

  async function handleResetPassword(row: UserActionRow) {
    try {
      await resetPassword({ id: row.id });
      message.success('密码重置成功');
    } catch (error) {
      logger.error(error);
      message.error('密码重置失败，请稍后重试');
    }
  }

  async function handleClearCache(row: UserActionRow) {
    try {
      await clearUserCache({ id: row.id });
      message.success('清除缓存成功');
    } catch (error) {
      logger.error(error);
      message.error('更新缓存失败，请稍后重试');
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
      message.success('导入成功');
      options.fetchTableData();
    } catch (error) {
      logger.error(error);
      message.error('导入失败，请检查文件后重试');
    } finally {
      importLoading.value = false;
      input.value = '';
    }
  }

  async function handleExport() {
    exportLoading.value = true;
    try {
      const response = await exportUserList(options.buildRequestParams(false));
      downloadResponseBlob(response, '用户列表.xlsx');
      message.success('导出成功');
    } catch (error) {
      logger.error(error);
      message.error('导出失败，请稍后重试');
    } finally {
      exportLoading.value = false;
    }
  }

  async function handleDownloadTemplate() {
    templateLoading.value = true;
    try {
      const response = await downloadUserImportTemplate();
      downloadResponseBlob(response, '用户导入模板.xlsx');
      message.success('模板下载成功');
    } catch (error) {
      logger.error(error);
      message.error('模板下载失败，请稍后重试');
    } finally {
      templateLoading.value = false;
    }
  }

  function openSetHomePageDialog(row: UserActionRow) {
    if (isSuperAdmin(row)) {
      message.warning('超级管理员不可设置首页');
      return;
    }

    selectedHomePageUserId.value = Number(row.id);
    selectedHomePage.value = row.dashboard || '';
    setHomePageVisible.value = true;
  }

  async function handleSetHomePage() {
    if (!selectedHomePageUserId.value) {
      message.warning('用户信息无效');
      return;
    }

    if (!selectedHomePage.value) {
      message.warning('请选择用户首页');
      return;
    }

    setHomePageLoading.value = true;
    try {
      await setHomePage({
        dashboard: selectedHomePage.value,
        id: selectedHomePageUserId.value,
      });
      message.success('设置首页成功');
      setHomePageVisible.value = false;
      options.fetchTableData();
    } catch (error) {
      logger.error(error);
      message.error('设置首页失败，请稍后重试');
    } finally {
      setHomePageLoading.value = false;
    }
  }

  function handleActionDropdownClick(
    data: Pick<UserActionDropdownItem, 'value'>,
    row: UserActionRow,
  ) {
    if (data.value === 'reset_password') {
      const dialogInstance = dialog.confirm({
        body: '确认重置该用户密码吗？',
        header: '提示',
        onClose: () => dialogInstance.hide(),
        onConfirm: () => {
          void handleResetPassword(row);
          dialogInstance.hide();
        },
      });
      return;
    }

    if (data.value === 'clear_cache') {
      const dialogInstance = dialog.confirm({
        body: '确认更新该用户缓存吗？',
        header: '提示',
        onClose: () => dialogInstance.hide(),
        onConfirm: () => {
          void handleClearCache(row);
          dialogInstance.hide();
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
