# Admin UI 重构实施方案（细化版）

更新时间：2026-03-25  
适用范围：`docs/old_admin`（旧前端） -> `admin-ui/apps/backend`（新前端）

---

## 1. 背景与目标

当前重构处于“框架已搭好、业务页未迁完”阶段：

- 旧项目优势：配置驱动快，很多页面通过 `ma-crud + columns` 可快速出 CRUD。
- 旧项目问题：封装过深、灵活性不足，业务特殊场景改造成本高。
- 新项目优势：基于 Vben5 + TDesign + TS，分层更清晰，扩展性更好。
- 新项目问题：目前除 `system/user` 外，大量系统页还是占位状态。

目标：采用“类似 RuoYi-Vue3 的清晰分层 + 适度复用”的方式完成页面迁移，达到：

1. 页面开发效率接近旧项目；
2. 可维护性、类型安全、灵活性优于旧项目；
3. 不再引入新一代“黑盒大封装”。

---

## 2. 范围与边界

### 2.1 本次重构范围

- `admin-ui/apps/backend/src/views/system/**`
- `admin-ui/apps/backend/src/api/system/**`
- 相关公共能力：`composables`、`components`、`types`、`router/access` 兼容逻辑。

### 2.2 明确不做（已确认）

- 旧项目 `system/module` 页面：不迁移。
- 旧项目 `system/code` 页面：不迁移。

### 2.3 明确保留

- 后端模块相关能力与权限码不清理（例如 `system:systemModules:*` 等仍保留）。
- 新项目中已有 `systemModules` 页面（占位）仍按新范围保留并继续实现。

---

## 3. 现状盘点（代码现状）

### 3.1 页面现状

- 旧项目 `system` 目录：19 个主页面（含 `module/code`）。
- 新项目 `system` 目录：
  - 已实现：`user/index.vue`
  - 占位页：16 个（`api/apiGroup/app/.../systemModules`）
  - 未迁移且确认下线：`module/code`

### 3.2 API 现状

- 旧项目 `src/api/system`：27 个模块文件。
- 新项目 `src/api/system`：5 个模块文件（`user/role/dept/post/dict`）。
- 差距：大量页面迁移会被 API 层覆盖不足阻塞。

### 3.3 关键技术债

1. 分页参数不统一风险：新 `user` 页仍使用 `limit`，后端标准字段为 `pageSize`。  
2. 页面重复逻辑开始出现：`search/pagination/fetchTableData` 模式已在不同页面手写。  
3. 权限与菜单兼容风险：后端动态菜单如果返回已下线路径，需前端容错映射，避免路由组件缺失导致运行错误。  

---

## 4. 目标架构（参考 RuoYi 思路，避免重封装）

核心原则：**薄封装、强约定、可落地**。

### 4.1 分层模型

1. `api` 层（纯接口）：只做请求与类型定义，不放页面逻辑。  
2. `hooks/composables` 层（复用状态）：复用列表页通用状态、查询、分页、批量操作。  
3. `schema` 层（页面配置）：搜索项、表格列、表单 schema。  
4. `view` 层（业务编排）：页面特有交互、权限按钮、弹窗流程。  

### 4.2 推荐目录结构

```text
admin-ui/apps/backend/src/
  api/system/
    user.ts
    role.ts
    ...
  views/system/<module>/
    index.vue
    components/
      <module>-modal.vue
      ...
    model.ts
    schemas.ts
    use-<module>-crud.ts
  composables/crud/
    use-crud-page.ts
    use-dict-options.ts
  types/
    paging.ts
    common.ts
```

### 4.3 封装边界（必须遵守）

- 不再做 `ma-crud` 一体化黑盒组件。
- 允许做通用 Hook 和小组件，但页面业务流程保留在页面侧。
- 复杂业务允许“半配置 + 半手写”，不强求纯配置。

---

## 5. 统一规范（迁移前先定规矩）

### 5.1 请求/响应规范

- 分页请求统一：`page`、`pageSize`。  
- 列表响应统一消费：`items` + `pageInfo.total`。  
- 业务成功码统一按现有 `requestClient` 的 `successCode: 0`。

### 5.2 页面规范

- 每个系统页至少拆分：`index.vue + modal + schemas + api`。  
- `index.vue` 不超过 500 行，超出必须拆到 `composables/components`。  
- 权限判断统一走 `@vben/access`（组件/指令/API 任一方式，但项目内保持一致）。

### 5.3 类型规范

- 禁止新增裸 `any` 接口返回；必须定义最小类型（`ListItem/ListQuery/FormModel`）。  
- 所有 `api/system/*.ts` 输出函数均有入参类型。  

---

## 6. 页面迁移清单（不含 module/code）

### 6.1 P0（核心管理能力）

1. `menu`
2. `role`
3. `dept`
4. `dict`
5. `post`

### 6.2 P1（系统运维与配置）

1. `config`
2. `crontab`
3. `notice`
4. `attachment`
5. `monitor/onlineUser`
6. `monitor/cache`

### 6.3 P2（业务扩展）

1. `api`
2. `apiGroup`
3. `app`
4. `appGroup`
5. `systemModules`

> `user` 作为模板页继续标准化，不算新增页面数量。

---

## 7. 分阶段实施计划

## Phase 0：基线治理（1 周）

目标：先把“会反复踩坑”的基础问题一次修正。

任务：

1. 修复 `system/user` 分页参数为 `page/pageSize`。
2. 完善 `system/user` 对齐项（超级管理员保护、设置首页、导入导出、编辑回填一致性）。
3. 建立 `types/paging.ts` 与基础响应类型。
4. 增加路由组件缺失保护（针对已下线页面路径做容错映射或过滤策略）。
5. 输出《页面开发模板说明》。

验收：

- `user` 页成为后续迁移模板；
- `pnpm -F @vben/backend typecheck` 通过；
- `user` 核心流程冒烟通过（查询、分页、增删改、回收站）。

## Phase 1：公共能力抽取（1 周）

目标：提效但不黑盒。

任务：

1. 实现 `use-crud-page.ts`（查询、分页、loading、批量选中、刷新）。
2. 实现 `use-dict-options.ts`（字典缓存与 options 转换）。
3. 实现通用工具栏组件（刷新、回收站切换、列显隐）。
4. 明确标准页面模板（可复制脚手架）。

验收：

- 至少在 `user + 1` 页面落地；
- 页面代码重复度明显下降（同类逻辑复用 Hook）。

## Phase 2：P0 页面迁移（2 周）

目标：先打通 RBAC 主链路。

任务：

1. 完成 `menu/role/dept/dict/post` 页面和 API。
2. 对齐旧功能核心点：搜索、分页、状态切换、回收站（有则迁移）、权限按钮。
3. 与动态菜单/权限联调（后端 `routers + codes`）。

验收：

- P0 页面可用；
- 权限控制正确（无权限按钮不可见/不可操作）；
- 回归 `user` 无退化。

## Phase 3：P1 页面迁移（2 周）

目标：补系统配置与运维页面。

任务：

1. 完成 `config/crontab/notice/attachment/monitor/*`。
2. 补齐上传、导出、日志查看等场景组件。
3. 补页面级异常处理与空状态规范。

验收：

- P1 页面可用；
- 常见异常（接口失败、空数据、权限缺失）展示一致。

## Phase 4：P2 页面迁移（1~2 周）

目标：补全剩余业务页并收口。

任务：

1. 完成 `api/apiGroup/app/appGroup/systemModules`。
2. 全量页面走一次交互巡检。
3. 输出页面迁移完成清单与差异说明。

验收：

- 范围内页面全部可用；
- 下线页（`module/code`）不在前端实现范围，行为可控（菜单侧或路由侧不崩）。

## Phase 5：提效与稳态（1 周）

目标：把“可持续开发能力”补齐。

任务：

1. 增加 `gen:crud` 脚手架（生成 `api + model + schemas + index + modal`）。
2. 建立页面 PR 检查项（类型、权限、空状态、错误处理、交互一致性）。
3. 输出开发手册（新增页面 10 分钟起步流程）。

验收：

- 新增一个标准 CRUD 页面可在 0.5~1 天内交付首版；
- 团队按统一模板开发，不再出现多套风格。

---

## 8. 质量保障与验收标准

### 8.1 通用验收清单（每页）

1. 列表：查询、重置、分页、刷新正常。  
2. 操作：新增、编辑、删除、批量、回收站（如有）正常。  
3. 权限：无权限按钮不可见或不可用。  
4. 异常：接口失败时提示清晰、无白屏。  
5. 类型：`typecheck` 通过，无新增高风险 `any`。  

### 8.2 项目级验收

1. 覆盖范围内页面全部完成。  
2. `module/code` 不实现且不引发前端运行时崩溃。  
3. 关键链路回归通过：登录 -> 菜单加载 -> 用户/角色/菜单联动。  

---

## 9. 风险与应对

1. 后端字段差异（旧接口“看似相同”但细节不同）  
应对：先统一 API 类型层，页面只消费类型，不直接拼字段。  

2. 动态菜单返回未实现页面组件  
应对：在 access 转换阶段加入路径容错与降级组件。  

3. 迁移节奏中重复代码蔓延  
应对：Phase 1 必须先落公共 Hook，再批量迁移。  

4. 权限码历史包袱（新旧混用）  
应对：保持兼容，不在本期清理权限码，仅做前端兼容映射。  

---

## 10. 交付物清单

1. 迁移后的系统页面与 API 文件。  
2. 公共 CRUD 能力（hooks + 组件）。  
3. 页面模板与脚手架（`gen:crud`）。  
4. 验收报告（页面清单、功能对齐清单、遗留问题清单）。  

---

## 11. 本周可立即执行任务（启动包）

1. 修 `system/user` 的 `pageSize` 参数兼容问题。  
2. 将 `system/user` 升级为“标准模板页”并沉淀 `use-crud-page`。  
3. 开始迁移 `role/menu/dept`（优先 P0）。  
4. 在 `docs` 增补《页面模板开发说明》与《迁移进度表》。  

---

## 12. 备注

- 本方案已按当前确认边界调整：`module/code` 不迁移；模块能力与权限不清理。  
- 若后续恢复 `module/code`，可按本方案模板增量恢复，不影响整体架构。  

