# Admin 系统页开发模板（Week 1 产物）

更新时间：2026-03-25  
适用目录：`admin-ui/apps/backend/src/views/system/*`

---

## 1. 目标

统一系统管理页面的开发方式，避免两种极端：

1. 过度封装成黑盒（难排查、难定制）。
2. 完全手写无约束（重复代码多、风格不一致）。

采用“半模板、半业务”的方式：

- 模板负责：列表页通用结构、表单弹窗结构、API 分层。
- 页面负责：业务动作、权限差异、特殊交互。

---

## 2. 标准目录结构

```text
src/views/system/<module>/
  index.vue                  # 列表页（页面编排）
  components/
    <module>-modal.vue       # 新增/编辑弹窗
    ...                      # 业务子组件
  model.ts                   # 页面类型（查询、实体、表单）
  schemas.ts                 # 搜索项/表单 schema（可选）
  use-<module>-crud.ts       # 列表逻辑 Hook（可选）
```

对应 API：

```text
src/api/system/<module>.ts
```

---

## 3. index.vue 模板职责

`index.vue` 推荐只做页面编排，不做重逻辑堆积。

建议包含：

1. 查询区：`searchForm`、`handleSearch`、`handleReset`。  
2. 表格区：`tableData`、`loading`、`pagination`。  
3. 动作区：新增、删除、导入、导出、回收站切换。  
4. 弹窗区：`<ModuleModal ref="..." @success="fetchTableData" />`。  
5. 业务钩子：状态切换、特殊行保护（如超级管理员保护）。  

---

## 4. modal 组件模板职责

`<module>-modal.vue` 负责新增/编辑表单流程：

1. `useVbenForm` 定义 schema。  
2. `useVbenModal` 管理打开关闭与提交 loading。  
3. `open(data?)` 时加载字典/选项数据。  
4. 编辑场景按需调用详情接口（`read/{id}`）做回填。  

---

## 5. API 层模板职责

每个 `api/system/<module>.ts`：

1. 只处理接口，不耦合页面 UI。  
2. 提供最小完整集：`list/recycle/save/update/delete/realDelete/recovery`（按实际接口有无）。  
3. 文件导入导出场景统一提供：
- `importXxx(file: File)`  
- `exportXxx(params)`  
- `downloadXxxTemplate()`（如有）  

---

## 6. 分页与返回约定

统一使用后端标准字段：

1. 请求：`page`、`pageSize`。  
2. 响应：`items` + `pageInfo.total`。  

页面读取示例：

```ts
tableData.value = res?.items || [];
pagination.total = res?.pageInfo?.total || 0;
```

---

## 7. 权限与安全约定

1. 按钮权限统一走 `@vben/access`（组件/指令/API 任一种，项目内一致）。  
2. 危险操作统一二次确认（删除、重置密码等）。  
3. 必须加业务保护项（例如 `id=1` 超级管理员不可删不可禁用）。  

---

## 8. 页面验收最小清单

每个系统页上线前至少满足：

1. 查询、重置、分页、刷新可用。  
2. 新增、编辑、删除、批量操作可用。  
3. 回收站（若接口有）可用。  
4. 错误提示清晰，不白屏。  
5. `typecheck` 通过。  

