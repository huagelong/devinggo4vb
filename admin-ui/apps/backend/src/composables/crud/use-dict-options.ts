import { getDictList } from '#/api/system/dict';

export interface DictItem {
  [key: string]: any;
  key: number | string;
  title: string;
}

export interface DictOption {
  label: string;
  value: number | string;
}

interface DictOptionConfig {
  labelKey?: string;
  valueKey?: string;
}

const dictCache = new Map<string, DictItem[]>();
const dictPromiseCache = new Map<string, Promise<DictItem[]>>();

function toOptions(list: DictItem[], config?: DictOptionConfig): DictOption[] {
  const labelKey = config?.labelKey ?? 'title';
  const valueKey = config?.valueKey ?? 'key';

  return (list || []).map((item) => ({
    label: item?.[labelKey],
    value: item?.[valueKey],
  }));
}

async function fetchDictList(code: string, forceRefresh = false): Promise<DictItem[]> {
  if (!forceRefresh && dictCache.has(code)) {
    return dictCache.get(code)!;
  }

  if (!forceRefresh && dictPromiseCache.has(code)) {
    return dictPromiseCache.get(code)!;
  }

  const requestPromise = getDictList(code)
    .then((response) => {
      const list = Array.isArray(response) ? response : [];
      dictCache.set(code, list);
      return list;
    })
    .catch(() => [])
    .finally(() => {
      dictPromiseCache.delete(code);
    });

  dictPromiseCache.set(code, requestPromise);
  return requestPromise;
}

export function useDictOptions() {
  async function getDictOptions(
    code: string,
    config?: DictOptionConfig & { forceRefresh?: boolean },
  ) {
    const list = await fetchDictList(code, config?.forceRefresh);
    return toOptions(list, config);
  }

  async function getMultipleDictOptions(
    codes: string[],
    config?: DictOptionConfig & { forceRefresh?: boolean },
  ) {
    const entries = await Promise.all(
      codes.map(async (code) => [code, await getDictOptions(code, config)] as const),
    );
    return Object.fromEntries(entries) as Record<string, DictOption[]>;
  }

  function clearDictCache(code?: string) {
    if (code) {
      dictCache.delete(code);
      dictPromiseCache.delete(code);
      return;
    }

    dictCache.clear();
    dictPromiseCache.clear();
  }

  return {
    clearDictCache,
    getDictOptions,
    getMultipleDictOptions,
  };
}
