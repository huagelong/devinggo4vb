# Admin 主从联动模板（Week 3 产物）

更新时间：2026-03-29
适用目录：`admin-ui/apps/backend/src/views/system/dict/`

---

## 1. 模板概述

主从联动模板适用于**一对多关联数据**的展示与维护，典型场景：
- 字典管理（字典类型 + 字典数据）
- 类似"分类 + 分类明细"结构

核心特征：
- 左侧/上方：主表（父数据列表）
- 右侧/下方/Modal：从表（子数据列表，受父数据约束）
- 主表选择后刷新从表数据
- 新增从表数据时自动注入关联字段

---

## 2. 标准目录结构

```text
src/views/system/<module>/
  index.vue                    # 主表列表页
  components/
    <module>-form-modal.vue    # 主表新增/编辑弹窗
    <module>-detail-panel.vue  # 从表面板（Modal 或 独立区域）
    <module>-detail-form-modal.vue  # 从表新增/编辑弹窗
  model.ts                     # 页面类型定义
  schemas.ts                   # 列定义、表单默认值工厂函数
  use-<module>-crud.ts         # 主表 CRUD Hook（使用 useCrudPage）
```

---

## 3. index.vue 主表页模板

```vue
<script setup>
// 主表使用 useCrudPage
const {
  fetchTableData,
  pagination,
  tableData,
  // ...
} = useXxxCrud();

// 从表面板 ref
const detailPanelRef = ref();

// 点击主表行打开从表
function handleOpenDetail(row: XxxListItem) {
  detailPanelRef.value?.open(row);
}
</script>

<template>
  <Page>
    <div class="flex flex-col gap-3">
      <!-- 主表搜索区 -->
      <SearchForm />

      <!-- 主表工具栏 -->
      <Toolbar>
        <Button @click="handleAdd">新增</Button>
      </Toolbar>

      <!-- 主表表格 -->
      <Table :data="tableData" :pagination="pagination">
        <!-- 主表操作列：编辑 + 打开从表 -->
        <template #action="{ row }">
          <Button @click="handleEdit(row)">编辑</Button>
          <Button @click="handleOpenDetail(row)">查看数据</Button>
        </template>
      </Table>
    </div>

    <!-- 主表弹窗 -->
    <XxxModal ref="modalRef" @success="fetchTableData" />

    <!-- 从表面板 -->
    <XxxDetailPanel ref="detailPanelRef" />
  </Page>
</template>
```

---

## 4. DetailPanel 从表面板实现

```vue
<script setup>
// 从表独立管理自己的状态
const searchForm = reactive(createDetailSearchForm());
const tableData = ref<DetailItem[]>([]);
const currentMaster = ref<{ id: number; code: string }>();
const pagination = reactive({ current: 1, pageSize: 20, total: 0 });

// 打开时接收主表数据，初始化查询
async function open(masterRow: MasterItem) {
  currentMaster.value = { id: masterRow.id, code: masterRow.code };
  Object.assign(searchForm, {
    master_id: masterRow.id,
    code: masterRow.code,
  });
  pagination.current = 1;
  await fetchTableData();
  panelModal.open();
}

// 从表 API 请求时注入 master_id
async function fetchTableData() {
  const params = buildParams();
  params.master_id = currentMaster.value.id;
  const res = await getDetailList(params);
  tableData.value = res.items ?? [];
  pagination.total = res.pageInfo?.total ?? 0;
}
</script>

<template>
  <!-- 使用 VbenModal 封装 -->
  <Modal>
    <div class="flex flex-col gap-4">
      <!-- 从表标题显示主表信息 -->
      <div class="text-sm text-gray-600">
        当前主表：{{ currentMaster?.name }}（{{ currentMaster?.code }}）
      </div>

      <!-- 从表搜索区 -->
      <DetailSearchForm />

      <!-- 从表工具栏 -->
      <DetailToolbar />

      <!-- 从表表格 -->
      <DetailTable />
    </div>

    <!-- 从表弹窗 -->
    <DetailModal />
  </Modal>
</template>
```

---

## 5. 关联数据注入方式

### 5.1 打开从表时注入

```ts
// index.vue
function handleOpenDetail(row: DictTypeItem) {
  detailPanelRef.value?.open({
    id: row.id,
    code: row.code,
    name: row.name,
  });
}
```

### 5.2 新增从表时注入

```ts
// detail-panel.vue
function handleAdd() {
  detailModalRef.value?.open({
    typeInfo: currentMaster.value,  // 注入主表信息
  });
}

// detail-form-modal.vue
async function open(options: { data?: DetailItem; typeInfo: MasterInfo }) {
  if (options.typeInfo) {
    formApi.setValues({
      master_id: options.typeInfo.id,
      code: options.typeInfo.code,
    });
  }
}
```

---

## 6. 主从 vs 普通 CRUD 差异

| 特性 | 普通 CRUD | 主从联动 |
|------|---------|---------|
| 数据结构 | 单一列表 | 主表 + 从表分离 |
| 从表约束 | 无 | 受主表 selected 约束 |
| 新增从表 | 无 | 自动注入 master_id |
| 刷新影响 | 只刷新当前表 | 可选刷新主表或只刷新从表 |
| 删除从表 | 无 | 正常 CRUD |

---

## 7. 接入示例：dict

**API**：`api/system/dict.ts` - 包含 `getDictTypePageList` 和 `getDictDataPageList`

**Model**：
```ts
// 主表类型
export type DictTypeListItem = DictApi.DictTypeItem;

// 从表类型
export type DictDataListItem = DictApi.DictDataItem;
```

**Schema**：
```ts
// 主表列
export function createDictTypeTableColumns() { ... }

// 从表列
export function createDictDataTableColumns() { ... }
```

**从表面板**：`dict-data-panel.vue`
- 使用 `useVbenModal` 封装
- `open(row)` 接收主表行数据
- 内部维护从表独立状态

---

## 8. 验收清单

- [ ] 主表列表正常加载分页
- [ ] 点击主表行正确打开从表面板
- [ ] 从表标题正确显示主表信息
- [ ] 从表查询自动注入 master_id
- [ ] 新增从表自动注入 master_id 和 code
- [ ] 编辑从表回填正确
- [ ] 删除/回收站/恢复正常
- [ ] typecheck 通过
