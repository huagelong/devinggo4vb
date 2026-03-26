import type { Ref } from 'vue';

import { useCrudPage } from '#/composables/crud/use-crud-page';

import { getRecycleUserList, getUserList } from '#/api/system/user';

import type { UserListItem } from './model';
import { createUserSearchForm } from './schemas';

export function useUserCrud(currentDeptId: Ref<number | string>) {
  const crud = useCrudPage<UserListItem, ReturnType<typeof createUserSearchForm>>({
    defaultSearchForm: createUserSearchForm,
    fetchList: (params, context) =>
      context.isRecycleBin ? getRecycleUserList(params) : getUserList(params),
    buildParams: (form) => {
      const params: Record<string, any> = {};

      if (form.username) params.username = form.username;
      if (form.role_id !== undefined) params.role_id = form.role_id;
      if (form.phone) params.phone = form.phone;
      if (form.post_id !== undefined) params.post_id = form.post_id;
      if (form.email) params.email = form.email;
      if (form.status !== undefined) params.status = form.status;
      if (form.user_type) params.user_type = form.user_type;
      if (form.created_at?.length === 2 && form.created_at[0]) {
        params.created_at = form.created_at;
      }

      if (currentDeptId.value) {
        params.dept_id = currentDeptId.value;
      } else if (form.dept_id !== undefined) {
        params.dept_id = form.dept_id;
      }

      return params;
    },
    resolveTotal: (response) =>
      Number(response?.pageInfo?.total || response?.total || 0),
  });

  function handleDeptSelect(deptId: number | string) {
    currentDeptId.value = deptId;
    crud.pagination.current = 1;
    crud.fetchTableData();
  }

  function handleResetWithDept() {
    currentDeptId.value = '';
    crud.handleReset();
  }

  return {
    ...crud,
    handleDeptSelect,
    handleResetWithDept,
  };
}
