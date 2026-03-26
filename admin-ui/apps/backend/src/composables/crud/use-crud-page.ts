import { reactive, ref } from 'vue';

type RowKey = number | string;

interface CrudPageContext {
  isRecycleBin: boolean;
  pagination: {
    current: number;
    pageSize: number;
    total: number;
  };
}

interface UseCrudPageOptions<TItem, TSearchForm extends Record<string, any>> {
  buildParams?: (
    searchForm: TSearchForm,
    context: CrudPageContext,
  ) => Record<string, any>;
  defaultSearchForm: () => TSearchForm;
  fetchList: (
    params: Record<string, any>,
    context: CrudPageContext,
  ) => Promise<any>;
  pageSize?: number;
  resolveItems?: (response: any) => TItem[];
  resolveTotal?: (response: any) => number;
}

export function useCrudPage<
  TItem = Record<string, any>,
  TSearchForm extends Record<string, any> = Record<string, any>,
>(options: UseCrudPageOptions<TItem, TSearchForm>) {
  const searchForm = reactive<TSearchForm>(options.defaultSearchForm());
  const tableData = ref<TItem[]>([]);
  const loading = ref(false);
  const selectedRowKeys = ref<RowKey[]>([]);
  const isRecycleBin = ref(false);

  const pagination = reactive({
    current: 1,
    pageSize: options.pageSize ?? 20,
    pageSizeOptions: [10, 20, 50, 100],
    showJumper: true,
    showPageSize: true,
    total: 0,
  });

  const resolveItems =
    options.resolveItems ??
    ((response: any) => (Array.isArray(response?.items) ? response.items : []));
  const resolveTotal =
    options.resolveTotal ??
    ((response: any) => Number(response?.pageInfo?.total || response?.total || 0));

  function getContext(): CrudPageContext {
    return {
      isRecycleBin: isRecycleBin.value,
      pagination: {
        current: pagination.current,
        pageSize: pagination.pageSize,
        total: pagination.total,
      },
    };
  }

  function buildRequestParams(includePagination = true) {
    const context = getContext();
    const businessParams = options.buildParams
      ? options.buildParams(searchForm as TSearchForm, context)
      : { ...searchForm };
    if (!includePagination) {
      return businessParams;
    }
    return {
      page: pagination.current,
      pageSize: pagination.pageSize,
      ...businessParams,
    };
  }

  function resetSearchForm() {
    const defaults = options.defaultSearchForm();
    Object.keys(defaults).forEach((key) => {
      (searchForm as Record<string, any>)[key] =
        (defaults as Record<string, any>)[key];
    });
  }

  async function fetchTableData() {
    loading.value = true;
    try {
      const response = await options.fetchList(buildRequestParams(true), getContext());
      tableData.value = resolveItems(response);
      pagination.total = resolveTotal(response);
    } catch (error) {
      console.error(error);
    } finally {
      loading.value = false;
    }
  }

  function clearSelectedRowKeys() {
    selectedRowKeys.value = [];
  }

  function handleSelectChange(keys: RowKey[]) {
    selectedRowKeys.value = keys;
  }

  function handleSearch() {
    pagination.current = 1;
    fetchTableData();
  }

  function handleReset() {
    resetSearchForm();
    pagination.current = 1;
    fetchTableData();
  }

  function handlePageChange(pageInfo: { current: number; pageSize: number }) {
    pagination.current = pageInfo.current;
    pagination.pageSize = pageInfo.pageSize;
    fetchTableData();
  }

  function toggleRecycleBin(next?: boolean) {
    isRecycleBin.value = typeof next === 'boolean' ? next : !isRecycleBin.value;
    clearSelectedRowKeys();
    pagination.current = 1;
    fetchTableData();
  }

  return {
    buildRequestParams,
    clearSelectedRowKeys,
    fetchTableData,
    handlePageChange,
    handleReset,
    handleSearch,
    handleSelectChange,
    isRecycleBin,
    loading,
    pagination,
    searchForm,
    selectedRowKeys,
    tableData,
    toggleRecycleBin,
  };
}
