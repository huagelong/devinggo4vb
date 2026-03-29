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
- 统一”刷新 / 回收站切换 / 列显隐”工具栏。
- 通过 `v-model` 控制表格显示列。

---

## 2. 页面模板类型

根据页面形态不同，分为三类模板：

| 模板类型 | 适用页面 | 关联文档 |
|---------|---------|---------|
| 标准列表 CRUD | post / notice / apiGroup | `ADMIN-Page-Template.md` |
| 树形 CRUD | menu / dept | `ADMIN-Tree-CRUD-Template.md` |
| 主从联动 | dict | `ADMIN-MasterDetail-Template.md` |
| 扩展动作 | role | `ADMIN-Extended-Action-Template.md` |

---

## 3. 页面接入示例

### 3.1 `system/user`

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

### 3.2 `system/post`

- 已从占位页升级为可用 CRUD 模板页：
1. `index.vue`
2. `components/post-modal.vue`
3. `model.ts`
4. `schemas.ts`
5. `use-post-crud.ts`

- API 已补齐 `list/recycle/save/update/delete/realDelete/recovery/changeStatus/sort`。

### 3.3 `system/menu`（树形 CRUD）

- 参考 `ADMIN-Tree-CRUD-Template.md`
- 目录结构：
  1. `index.vue` - 树形表格页
  2. `components/menu-modal.vue` - 新增/编辑弹窗
  3. `model.ts` - 类型定义（MenuTreeItem）
  4. `schemas.ts` - 列定义、表单默认值
  5. `use-menu-page.ts` - 列表逻辑

### 3.4 `system/dept`（树形 CRUD + 领导列表）

- 参考 `ADMIN-Tree-CRUD-Template.md`
- 额外组件：`dept-leader-modal.vue`（领导列表管理）

### 3.5 `system/dict`（主从联动）

- 参考 `ADMIN-MasterDetail-Template.md`
- 目录结构：
  1. `index.vue` - 主表（字典类型）列表
  2. `components/dict-type-modal.vue` - 主表弹窗
  3. `components/dict-data-panel.vue` - 从表（字典数据）面板
  4. `components/dict-data-form-modal.vue` - 从表弹窗
  5. `use-dict-type-crud.ts` - 主表 CRUD Hook

### 3.6 `system/role`（扩展动作 + 业务保护）

- 参考 `ADMIN-Extended-Action-Template.md`
- 目录结构：
  1. `index.vue` - 列表页
  2. `components/role-modal.vue` - 新增/编辑弹窗
  3. `components/role-menu-permission-modal.vue` - 菜单权限
  4. `components/role-data-permission-modal.vue` - 数据权限
  5. `use-role-crud.ts` - 主表 CRUD Hook

---

## 4. 新增页面建议流程

1. **判断页面模板类型**：标准列表 / 树形 CRUD / 主从联动 / 扩展动作
2. 新建模块目录：`views/system/<module>/`
3. 创建对应文件：
   - 标准列表：`model.ts + schemas.ts + use-<module>-crud.ts + components/<module>-modal.vue`
   - 树形 CRUD：`model.ts + schemas.ts + use-<module>-page.ts + components/<module>-modal.vue`
   - 主从联动：参考 dict 目录结构
   - 扩展动作：参考 role 目录结构
4. 在 `index.vue` 只保留页面编排和模块业务动作。
5. 统一复用：
   - `use-crud-page.ts`（标准列表/主从联动）
   - `use-dict-options.ts`
   - `crud-toolbar.vue`
6. 树形页面使用独立 hook（如 `use-menu-page.ts`）管理树形数据

