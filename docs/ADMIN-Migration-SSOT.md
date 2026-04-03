# Admin 前端迁移统一执行文档（SSOT）

更新时间：2026-04-03
适用目标：将 docs/old_admin 的前端能力重构到 admin-ui/apps/backend

---

## 1. 文档定位

本文件是唯一执行文档（Single Source of Truth）。

- 不再并行维护多份周计划与阶段清单。
- old_admin 仅作为功能对照基线，不作为实现规范。
- 实现规范统一采用 admin-ui 当前技术栈与目录结构。

---

## 2. 事实源优先级（冲突裁决规则）

当不同来源内容冲突时，按以下优先级裁决：

1. 当前代码现状（admin-ui/apps/backend + modules/system/api/system）
2. 本文档
3. docs/old_admin 代码
4. 历史计划文档（PLAN/TODO/模板说明）

说明：
- old_admin 用于回答“旧系统有哪些功能”。
- 后端 API 用于回答“现在到底该对接什么接口”。
- admin-ui 现状用于回答“哪些已经完成，哪些只差闭环”。

---

## 3. 迁移目标与边界

### 3.1 迁移目标

在不回退旧技术方案的前提下，完成以下目标：

- 功能覆盖达到 old_admin 的系统管理能力。
- 页面实现遵循 admin-ui 的模板化结构（列表 CRUD / 树形 CRUD / 主从联动 / 扩展动作）。
- 类型、错误提示、权限与回收站交互保持一致。

### 3.2 边界

- in：system 业务页、dashboard 相关页、对应 API 封装、必要共享类型与轻量工具。
- out：回退 ma-crud 黑盒封装、直接复制 old_admin 架构。
- 差异保留：old_admin 的 system/module 不作为当前迁移目标，统一采用 systemModules 能力。

---

## 4. old_admin -> admin-ui 模块映射（当前状态）

状态说明：
- 已迁移：页面与 API 已有可用实现。
- 部分迁移：已有 API 或部分页面，仍需闭环。
- 未迁移：未发现对应页面实现。

| old_admin 能力 | 当前 admin-ui 对应 | 状态 | 备注 |
|---|---|---|---|
| login | dashboard/_core 登录体系 | 已迁移 | 保持现有实现 |
| dashboard | views/dashboard/analytics | 已迁移 | 已对接 dashboard 统计/图表 |
| user | views/system/user | 已迁移 | 继续做联调回归 |
| role | views/system/role | 已迁移 | 包含权限扩展动作 |
| menu | views/system/menu | 已迁移 | 树形 CRUD |
| dept | views/system/dept | 已迁移 | 树形 CRUD |
| post | views/system/post | 已迁移 | 标准 CRUD |
| dict | views/system/dict | 已迁移 | 主从联动 |
| config | views/system/config | 已迁移 | 标准 CRUD |
| crontab | views/system/crontab | 已迁移 | 含日志面板 |
| attachment | views/system/attachment | 已迁移 | 标准 CRUD |
| notice | views/system/notice | 已迁移 | 标准 CRUD |
| api/apiGroup/app/appGroup | views/system 对应目录 | 已迁移 | 标准 CRUD |
| monitor/onlineUser, cache | views/system/monitor/* | 已迁移 | 已有 monitor API |
| systemModules | views/system/systemModules | 已迁移 | 作为模块管理目标实现 |
| code | views/system/code | 已迁移 | 已存在生成器页面 |
| logs/loginLog, operLog, apiLog | 无对应 system/logs 页面 | 未迁移 | 仅有 API 封装，需补页面 |
| dataMaintain | 无对应页面 | 未迁移 | 后端接口已存在 |
| queueMessage | views/dashboard/message | 部分迁移 | 页面可用，类型需收敛 any |
| upload 专项能力 | profile 等局部上传 | 部分迁移 | 缺独立上传管理能力整合 |
| pusher 实时能力 | 暂无明确接入 | 未迁移 | 后端接口存在，前端待建设 |

---

## 5. 迁移执行计划（按闭环优先）

## T0：基线稳定（1 周）

目标：把现有已迁移页面收敛为可复制模板，避免继续扩散技术债。

- [ ] 完成 user 与 post 冒烟联调闭环（查询/分页/增改删/状态/回收站）。
- [ ] 清理公共层高频 any（优先 use-crud-page、dict options、message API）。
- [ ] 统一失败提示策略（用户可见提示 + 控制台日志并存）。
- [ ] 统一分页参数约定（page/pageSize），避免 page_size 与 pageSize 混用。

完成标准：
- [ ] typecheck 通过
- [ ] user/post 作为模板页稳定复用

## T1：补齐日志域（1 周）

目标：完成 old_admin 的 logs 三页迁移闭环。

- [ ] 新建 system/logs/loginLog 页面
- [ ] 新建 system/logs/operLog 页面
- [ ] 新建 system/logs/apiLog 页面
- [ ] 接入现有 log API（列表、删除）
- [ ] 对齐统一列表模板（搜索、分页、批量、错误提示、空态）

完成标准：
- [ ] logs 三页可在联调环境可用
- [ ] 与旧版功能对齐到可回归状态

## T2：补齐数据维护域（0.5~1 周）

目标：完成 dataMaintain 能力迁移。

- [ ] 新增 api/system/data-maintain.ts（或并入规范命名文件）
- [ ] 新建 system/dataMaintain 页面
- [ ] 落地表列表、字段详情、优化、碎片整理操作
- [ ] 补齐风险操作确认与反馈文案

完成标准：
- [ ] dataMaintain 页面核心流程可用

## T3：消息与上传能力收口（1 周）

目标：把“部分迁移”能力升级为稳定能力。

- [ ] dashboard/message 去除高频 any，补齐最小类型
- [ ] 梳理 upload 与 attachment 的职责边界
- [ ] 统一上传/下载反馈、文件名处理、异常提示

完成标准：
- [ ] 消息中心和附件/上传链路具备稳定用户体验

## T4：实时能力建设（1~2 周）

目标：接入 pusher 相关功能，补齐 old_admin 之后的新需求能力。

- [ ] 设计实时能力最小范围（认证、连接、订阅、消息消费）
- [ ] 封装 websocket/pusher 客户端层
- [ ] 在消息中心或监控域完成首个落地点
- [ ] 增加异常重连、鉴权失败、连接状态提示

完成标准：
- [ ] 至少 1 条核心实时链路在联调可稳定运行

---

## 6. 统一开发约束

- 页面结构遵循现有模板，不引入黑盒大封装。
- API 文件只做接口与类型，不写页面业务流程。
- 新增能力优先复用现有 composables 与组件。
- 所有关键动作必须提供用户可见反馈。
- 危险动作统一二次确认。

---

## 7. 验收清单（全局）

每个迁移模块至少满足：

- [ ] 查询、重置、分页、刷新可用
- [ ] 新增、编辑、删除（若有）可用
- [ ] 回收站能力（若后端提供）可用
- [ ] 权限按钮显隐正确
- [ ] 错误提示清晰、空状态明确
- [ ] typecheck 通过

全量验收通过条件：

- [ ] old_admin 基线中约定迁移范围内能力全部在 admin-ui 可访问
- [ ] 登录 -> 菜单 -> 系统管理主链路无阻塞缺陷
- [ ] 关键页面联调冒烟全部通过

---

## 8. 周推进记录（只更新本节）

> 用于持续推进，避免再拆分到多份 TODO。

### 2026-04-03

- 统一目标确认：以 old_admin 为功能基线，重构到 admin-ui。
- 当前最大缺口：logs 三页、dataMaintain、pusher 接入。
- 已完成：logs 三页页面骨架与列表删除能力落地（loginLog/operLog/apiLog）。
- 已完成：日志 API 端点切换为 `/system/logs/*` 路由组。
- 已完成：登录状态值语义对齐（1=成功，2=失败，其他=未知）。
- 已完成：操作/接口日志字段展示补齐（请求数据、响应数据）。
- 已完成：日志删除按钮权限点接入（system:loginLog:delete / system:operLog:delete / system:apiLog:delete）。
- 已完成：dataMaintain 首版骨架（api + model + schemas + use-crud + index）。
- 约束确认：当前后端仅开放 `/system/dataMaintain/index`，旧版 `detailed/optimize/fragment` 暂未提供。
- 下一步执行：联调日志页查询字段与响应字段，补菜单与权限点校验。
