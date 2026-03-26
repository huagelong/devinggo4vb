# Admin 公共 CRUD 能力说明（Week 2）

更新时间：2026-03-25  
适用目录：`admin-ui/apps/backend/src/views/system/*`

---

## 1. 已落地公共能力

1. `src/composables/crud/use-crud-page.ts`
- 统一列表页的 `searchForm / tableData / loading / pagination / selectedRowKeys`。
- 内置 `handleSearch / handleReset / handlePageChange / toggleRecycleBin / fetchTableData`。
- 支持自定义 `buildParams`，并统一分页参数为 `page/pageSize`。

2. `src/composables/crud/use-dict-options.ts`
- 统一字典拉取与缓存，避免页面重复请求字典。
- 支持 `getDictOptions(code)` 和 `getMultipleDictOptions(codes)`。

3. `src/components/crud/crud-toolbar.vue`
- 统一“刷新 / 回收站切换 / 列显隐”工具栏。
- 通过 `v-model` 控制表格显示列。

---

## 2. 页面接入示例

### 2.1 `system/user`

- 页面结构已按模板拆分：
1. `index.vue`
2. `components/dept-tree.vue`
3. `components/user-modal.vue`
4. `model.ts`
5. `schemas.ts`
6. `use-user-crud.ts`
7. `use-user-actions.ts`

- `index.vue` 负责页面编排；
- 查询/分页/回收站由 `use-user-crud.ts` 承担；
- 行操作、导入导出、设置首页由 `use-user-actions.ts` 承担。

### 2.2 `system/post`

- 已从占位页升级为可用 CRUD 模板页：
1. `index.vue`
2. `components/post-modal.vue`
3. `model.ts`
4. `schemas.ts`
5. `use-post-crud.ts`

- API 已补齐 `list/recycle/save/update/delete/realDelete/recovery/changeStatus/sort`。

---

## 3. 新增页面建议流程

1. 新建模块目录：`views/system/<module>/`。
2. 创建 `model.ts + schemas.ts + use-<module>-crud.ts + components/<module>-modal.vue`。
3. 在 `index.vue` 只保留页面编排和模块业务动作。
4. 统一复用：
- `use-crud-page.ts`
- `use-dict-options.ts`
- `crud-toolbar.vue`

