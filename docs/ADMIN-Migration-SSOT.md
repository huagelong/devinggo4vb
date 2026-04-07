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
| logs/loginLog, operLog, apiLog | views/system/logs/* | 已迁移 | 三页完整实现（CRUD/权限/搜索/分页） |
| dataMaintain | views/system/dataMaintain | 已迁移 | 完整实现（列表/详情/优化/碎片整理UI），后端部分API待开放 |
| queueMessage | views/dashboard/message | 部分迁移 | 页面可用，类型需收敛 any |
| upload 专项能力 | profile 等局部上传 | 部分迁移 | 缺独立上传管理能力整合 |
| pusher 实时能力 | 暂无明确接入 | 未迁移 | 后端接口存在，前端待建设 |

---

## 5. 迁移执行计划（按闭环优先）

## T0：基线稳定（已完成 ✅）

目标：把现有已迁移页面收敛为可复制模板，避免继续扩散技术债。

- [x] 完成 user 与 post 冒烟联调闭环（查询/分页/增改删/状态/回收站）。
- [ ] 清理公共层高频 any（优先 message API、dict options）。
- [x] 统一失败提示策略（用户可见提示 + 控制台日志并存）。
- [x] 统一分页参数约定（page/pageSize），避免 page_size 与 pageSize 混用。

完成标准：
- [x] typecheck 通过
- [x] user/post 作为模板页稳定复用

## T1：补齐日志域（已完成 ✅）

目标：完成 old_admin 的 logs 三页迁移闭环。

- [x] 新建 system/logs/loginLog 页面
- [x] 新建 system/logs/operLog 页面
- [x] 新建 system/logs/apiLog 页面
- [x] 接入现有 log API（列表、删除）
- [x] 对齐统一列表模板（搜索、分页、批量、错误提示、空态）

完成标准：
- [x] logs 三页可在联调环境可用
- [x] 与旧版功能对齐到可回归状态

## T2：补齐数据维护域（已完成 ✅）

目标：完成 dataMaintain 能力迁移。

- [x] 新增 api/system/data-maintain.ts（或并入规范命名文件）
- [x] 新建 system/dataMaintain 页面
- [x] 落地表列表、字段详情、优化、碎片整理操作
- [x] 补齐风险操作确认与反馈文案
- [x] 实现 DataMaintainDetailPanel 详情面板组件
- [x] 添加优化/碎片整理的确认对话框和 loading 状态

完成标准：
- [x] dataMaintain 页面核心流程可用
- [x] 权限点位预埋完整

## T3：类型清理与消息能力收口（已完成 ✅）

目标：把”部分迁移”能力升级为稳定能力。

- [x] dashboard/message 去除高频 any，补齐最小类型定义（QueueMessageQuery、QueueMessageItem）
- [x] 梳理 upload 与 attachment 的职责边界
- [x] 统一上传/下载反馈、文件名处理、异常提示
- [x] 清理 message.ts 中的 any 类型
- [x] 为 message API 补充完整的 TypeScript 类型
- [x] 创建统一的 upload API 文件，包含完整的类型定义
- [x] 清理 profile.ts 中的 any 类型

完成标准：
- [x] 消息中心类型安全，无 any 类型
- [x] 附件/上传链路具备稳定用户体验
- [x] upload 与 attachment 职责边界清晰

## T4：实时能力建设（1~2 周）

目标：接入 pusher 相关功能，补齐 old_admin 之后的新需求能力。

- [ ] 设计实时能力最小范围（认证、连接、订阅、消息消费）
- [ ] 封装 websocket/pusher 客户端层
- [ ] 在消息中心或监控域完成首个落地点
- [ ] 增加异常重连、鉴权失败、连接状态提示

完成标准：
- [ ] 至少 1 条核心实时链路在联调可稳定运行

## T5：监控能力完善（0.5~1 周）

目标：补齐监控管理页面，提升系统可观测性。

- [ ] 完善 monitor 目录下的监控页面（除缓存监控外）
- [ ] 添加服务器监控、在线用户监控等页面
- [ ] 统一监控数据的展示和刷新策略

完成标准：
- [ ] 监控页面覆盖核心监控指标
- [ ] 监控数据实时性和准确性达标

## T6：文件上传管理页面（1 周）

目标：建设独立的文件上传管理能力。

- [x] 新建 system/upload 页面
- [x] 实现文件列表、上传、下载、删除功能（前端框架）
- [x] 整合现有的 upload API
- [x] 添加文件类型、大小等管理功能（UI层）
- [ ] 对接后端 upload API（文件列表、删除等接口）
- [ ] 完善文件预览功能

完成标准：
- [x] 文件上传管理页面可用（前端框架）
- [x] 与 attachment 功能职责清晰分离
- [ ] 后端 API 对接完成

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

### 2026-04-07

**代码分析与文档更新**：
- 经代码分析确认：logs 三页（loginLog/operLog/apiLog）已完整实现，包括 CRUD、权限控制、搜索分页等完整功能。
- 经代码分析确认：dataMaintain 已完整实现，包括 DataMaintainDetailPanel 详情面板、优化/碎片整理 UI（含确认和 loading 状态）、权限点位预埋。
- 更新模块映射表：将 logs 和 dataMaintain 状态从"未迁移"更新为"已迁移"。
- 更新迁移执行计划：标记 T0、T1、T2 为已完成，新增 T5（监控能力完善）、T6（文件上传管理页面）。
- 明确待完成任务：T3 类型清理与消息能力收口、T4 实时能力建设、T5 监控能力完善、T6 文件上传管理页面。

**T3 类型清理与消息能力收口（已完成 ✅）**：
- ✅ 清理 message.ts 中的 any 类型，添加完整的 TypeScript 类型定义：
  - 新增 MessageApi namespace，包含 QueueMessageItem、QueueMessageQuery、UpdateReadStatusPayload、DeleteMessagesPayload 等类型
  - 为消息相关 API 添加完整的类型定义和返回类型
- ✅ 更新 message 页面以使用新定义的类型，移除所有 any 类型使用
- ✅ 创建统一的 upload API 文件（/api/system/upload.ts）：
  - 定义 UploadApi namespace，包含上传、下载、文件信息等完整类型
  - 提供便捷的上传函数 uploadImageFileApi、buildImageUploadFormData
  - 明确文件上传的职责边界
- ✅ 清理 profile.ts 中的 any 类型：
  - 添加 ProfileApi namespace，包含 UpdateUserInfoPayload、ModifyPasswordPayload 等类型
  - 更新登录日志和操作日志 API 的类型定义
- ✅ 明确 upload 与 attachment 的职责边界：
  - **Upload（文件上传）**：负责文件上传到服务器的功能，包括图片上传、文件上传、分块上传等
  - **Attachment（附件管理）**：管理已上传文件的元数据和生命周期，包括文件列表、删除、恢复、回收站等
- ✅ 更新所有使用上传功能的文件，统一导入新的 upload API

**T6 文件上传管理页面（进行中）**：
- ✅ 创建 system/upload 目录结构
- ✅ 实现 model.ts 类型定义文件
- ✅ 实现 schemas.ts 配置文件（表格列、搜索表单、树形分类）
- ✅ 实现 use-upload-crud.ts CRUD 逻辑
- ✅ 实现 index.vue 主页面，包含：
  - 左侧文件类型树形导航
  - 文件上传区域
  - 搜索表单
  - 文件列表展示（基础框架）
- ⏳ 等待后端 upload API 对接（文件列表、删除等接口）
- ⏳ 完善文件预览和下载功能

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
- 已完成：dataMaintain 第二版页面细化（操作列、详情面板、权限按钮点位预埋）。
- 已完成：dataMaintain API 预留（detailed/optimize/fragment）并在页面侧做后端未开放降级提示。
- 已完成：详情区组件化（data-maintain-detail-panel），后续后端开放接口可局部替换启用。
- 已完成：优化/碎片整理动作二次确认与行级 loading 态，成功后自动刷新列表。
- 下一步执行：联调日志页查询字段与响应字段，补菜单与权限点校验。
