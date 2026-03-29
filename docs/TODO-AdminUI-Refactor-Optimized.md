# Admin UI 重构优化任务列表（基于当前代码现状）

更新时间：2026-03-28  
关联文档：`docs/PLAN-AdminUI-Refactor.md`、`docs/ADMIN-Page-Template.md`、`docs/ADMIN-Common-CRUD-Guide.md`

---

## 1. 当前代码快照

已落地：

- 公共 CRUD 基础能力已存在：
  - `admin-ui/apps/backend/src/composables/crud/use-crud-page.ts`
  - `admin-ui/apps/backend/src/composables/crud/use-dict-options.ts`
  - `admin-ui/apps/backend/src/components/crud/crud-toolbar.vue`
- 标准列表页样板已落地：
  - `system/user`
  - `system/post`

已知缺口：

- `api/system` 目前只有 `user/post/role/dept/dict` 5 个文件，其中：
  - `user/post` 相对完整；
  - `role/dept/dict` 仍是“支撑样板页/字典查询”的最小实现；
  - `menu.ts`、`config.ts`、`notice.ts`、`api.ts`、`app.ts` 等尚未建立。
- `views/system` 中多数页面仍是占位页：
  - `role/menu/dept/dict/config/notice/api/apiGroup/app/appGroup/attachment/crontab/systemModules/monitor/*`
- 当前模板主要覆盖“标准列表 CRUD”，但下一批待迁页面并不都是这一类：
  - `role` = 列表 CRUD + 权限分配动作
  - `menu/dept` = 树形 CRUD
  - `dict` = 主从联动页
- 现有样板页里仍有较多 `any` 和分散的错误处理，直接复制会把技术债继续扩散。
- `user` 页文档上仍有一个未关闭项：联调环境冒烟测试。

结论：

- 现在不适合继续按“每周平推若干页面”的方式推进。
- 更合理的方式是先把当前模板补强到足以支撑下一波页面，再按页面形态分批迁移。

---

## 2. 优化后的排期原则

1. 依赖先于排期  
   页面任务开始前，先补齐对应 `api + model + 最小类型`。

2. 模板先于批量迁移  
   当前已验证的只有“列表 CRUD 模板”；树形页、主从页应先沉淀模板再批量做。

3. 一次任务必须闭环  
   每一批任务都要尽量闭环到：`API -> 类型 -> 页面 -> 冒烟/联调`。

4. 不重复计算已完成页面  
   `post` 已完成模板试点，不再作为“待迁移页面”重复统计，只作为 P0 回归对象。

5. 先做高复用，再做高定制  
   在不违反业务优先级的前提下，优先迁移能复用现有模板的页面，降低切换成本。

6. 明确范围，避免隐性漏项  
   `logs/*` 等已出现在代码目录中的页面，要么明确纳入后续批次，要么明确标记暂不交付。

---

## 3. 新任务列表

## T0：基线收口

目标：把当前样板和公共能力收成“可复制基础版”，避免后续迁移时边做边返工。

任务：

- [ ] 完成 `system/user` 联调环境冒烟测试，正式关闭 Week 1 遗留项。
- [ ] 为现有 CRUD 基础补最小共享类型：
  - 分页请求/响应类型
  - 通用 `Option` / `TreeNode` / `IdPayload` 类型
  - 列表页基础查询类型
- [ ] 收敛当前样板页中的高频 `any`：
  - `api/system/user.ts`
  - `api/system/post.ts`
  - `api/system/role.ts`
  - `api/system/dept.ts`
  - `api/system/dict.ts`
  - `use-crud-page.ts`
  - `system/user/*`
  - `system/post/*`
- [ ] 统一当前 CRUD 页的错误提示策略，避免继续直接散落 `console.error`。
- [ ] 整理导入/导出/下载文件名解析等已在 `user` 页出现的通用逻辑，沉淀为轻量工具函数，不新增黑盒大 Hook。
- [ ] 明确当前发布范围中是否包含 `logs/*`，并同步菜单暴露策略。

完成标准：

- `user/post` 两个样板页可作为后续复制起点；
- 基础类型和错误处理方式不再明显拖后腿；
- 当前范围边界清楚，不会出现“页面在菜单里但仍是未知占位”的情况。

---

## T1：补齐 3 类页面模板

目标：先把下一波页面真正需要的模板补齐，而不是直接拿平铺列表模板硬套。

任务：

- [ ] 列表 CRUD 模板 v1.1
  - 基于 `user/post` 收敛标准结构
  - 明确哪些逻辑留在页面，哪些逻辑允许放在 `use-<module>-crud.ts`
- [ ] 树形 CRUD 模板
  - 面向 `menu/dept`
  - 统一树数据加载、展开、节点操作、状态切换、排序/层级编辑的页面结构
- [ ] 主从联动模板
  - 面向 `dict`
  - 统一“左侧类型/上层列表 + 右侧数据列表/弹窗”的组织方式
- [ ] 扩展动作模板
  - 面向 `role`
  - 明确“列表 CRUD + 授权/分配类动作”的拆分方式

完成标准：

- 后续页面开发不再只有一个“通用列表页模板”可选；
- `role/menu/dept/dict` 开发前已有对应页面骨架模式可套用。

---

## T2：下一批页面 API 与类型先行

目标：先补依赖，再开页面，避免 UI 开始后被 API 不完整卡住。

任务：

- [ ] 完成 `api/system/role.ts` 最小完整集：
  - `list/recycle/save/update/delete/realDelete/recovery/changeStatus`
  - 按实际需要补权限分配相关接口
- [ ] 新增 `api/system/menu.ts`：
  - 树列表
  - 新增/编辑/删除
  - 状态/排序
- [ ] 扩展 `api/system/dept.ts`：
  - 从“仅树查询”升级为完整页面 API
- [ ] 扩展 `api/system/dict.ts`：
  - 从“字典选项查询”升级为“字典类型 + 字典数据”双 API
- [ ] 为 `role/menu/dept/dict` 建立最小 `model.ts`
- [ ] 明确这些页面的查询字段、表格字段、表单字段，不在页面内临时拍脑袋命名

完成标准：

- `role/menu/dept/dict` 进入页面开发前，接口层已具备可直接消费的最小闭环；
- 页面开发时不再反复回填 API 文件。

---

## T3：P0 核心链路闭环

目标：优先打通真正影响 RBAC 主链路的页面，但按实现相似度排序，减少切换成本。

推荐顺序：

1. `role`
2. `dept`
3. `menu`
4. `dict`
5. `user/post` 回归

任务：

- [ ] 迁移 `role`
  - 列表、增改删、状态、回收站
  - 权限按钮
  - 授权/分配动作入口
- [ ] 迁移 `dept`
  - 树形展示
  - 增改删、状态
  - 与用户/组织关系相关的基础能力
- [ ] 迁移 `menu`
  - 树形展示
  - 增改删、状态、排序
  - 动态菜单联调
- [ ] 迁移 `dict`
  - 字典类型
  - 字典数据
  - 类型与数据联动
- [ ] 完成 P0 全量回归：
  - `user/role/menu/dept/dict/post`
- [ ] 收口交互一致性：
  - 空状态
  - 错误提示
  - 按钮位置
  - 回收站切换体验

完成标准：

- P0 页面全部可进入联调环境；
- 登录 -> 菜单加载 -> 用户/角色/菜单/部门 主链路可跑通；
- `post` 只做回归，不再单列一次迁移任务。

---

## T4：标准 CRUD 页面批量推进

目标：优先消化最容易复用现有模板的一批页面，快速拉升整体完成率。

建议批次：

- [ ] `config`
- [ ] `notice`
- [ ] `api`
- [ ] `apiGroup`
- [ ] `app`
- [ ] `appGroup`
- [ ] `systemModules`

说明：

- 这批页面大多更接近现有“列表 CRUD 模板”，更适合在 P0 后连续推进。
- 若业务优先级必须保持原 P1/P2 标签，可以保留标签，但执行上仍建议按该批次连续做。

完成标准：

- 这批标准页全部具备基础查询、分页、增改删、状态/回收站等能力；
- 迁移速度明显快于 P0 批次，说明模板开始真正起效。

---

## T5：高定制页面与收口补齐

目标：处理需要额外交互能力的页面，而不是继续假装它们都只是普通 CRUD。

任务：

- [ ] `crontab`
  - 任务 CRUD
  - 执行日志查看
- [ ] `attachment`
  - 上传、预览、删除、回收站
- [ ] `monitor/onlineUser`
- [ ] `monitor/cache`
- [ ] `logs/loginLog`
- [ ] `logs/operLog`
- [ ] `logs/apiLog`

配套公共能力：

- [ ] 上传/下载相关通用处理
  - loading
  - 文件名
  - 错误提示
- [ ] 日志类页面的列表与详情查看模式
- [ ] 监控类页面的空状态与异常提示模式

完成标准：

- 高定制页面不再混在标准 CRUD 批次中导致节奏被打断；
- 上传、日志、监控这些特殊能力有自己的轻量模式可复用。

---

## T6：稳态化与提效

目标：把“本轮迁移”变成“以后还能持续开发”的能力。

任务：

- [ ] 增加 `gen:crud` 脚手架
  - 生成 `api + model + schemas + index + modal`
- [ ] 增加页面 PR 检查清单
  - 类型
  - 权限
  - 错误处理
  - 空状态
  - 交互一致性
- [ ] 输出《新增页面开发指南》
- [ ] 输出《重构完成报告》
- [ ] 做一次系统页全量回归和体积/性能检查

完成标准：

- 新增一个标准 CRUD 页面可在 0.5~1 天内交付首版；
- 团队后续新增页面时不需要重新摸索结构。

---

## 4. 对旧任务列表的直接调整建议

- `post` 从“待迁移页”调整为“已完成样板页 + P0 回归项”。
- 在 `role/menu/dept/dict` 前插入“模板补齐”和“API/类型先行”两个阶段。
- 不再单纯按周平均铺页面，而改为：
  - 基线收口
  - 模板补齐
  - API 先行
  - P0 闭环
  - 标准页批量推进
  - 高定制页收口
  - 工具化与稳态化
- 明确 `logs/*` 的去向，避免继续处于”代码里有页面、计划里没位置”的状态。

**logs/* 范围确认**：
  - 当前状态：`logs/loginLog`、`logs/operLog`、`logs/apiLog` 均为占位页
  - 归属阶段：T5 高定制页面阶段
  - 前提条件：需要后端日志 API 完整支持后再启动
  - 建议：暂不交付，待 T5 阶段统一处理

---

## 5. 建议的下一步

如果只看“立刻可执行”的任务，建议按下面顺序启动：

1. 完成 `user` 页联调冒烟并关闭遗留项。
2. 补齐共享类型，收敛 `user/post` 中最明显的 `any`。
3. 先补 `role/menu/dept/dict` 的 API 与 `model.ts`。
4. 先做 `role`，再做 `dept/menu`，最后做 `dict`。

这样推进，和当前代码结构最贴近，也最不容易返工。
