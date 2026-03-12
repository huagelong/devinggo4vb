# Admin UI 开发计划 (Todo List)

该计划基于 `docs/admin.md` 的需求，按照功能优先级分类，旨在有条理、分阶段地推进商业化管理后台前端的开发任务。

## 阶段一：项目初始化与基础设施建设 (P0)
- [x] 初始化 `admin-ui` 目录，并确定具体技术栈核心版本：
  - **核心框架**: React `^18.3.1` + Vite `^5.2.0` + TypeScript `^5.4.5`
  - **路由与状态**: TanStack Router `^1.80.0` + Zustand `^5.0.0`
  - **网络与实时**: Axios `^1.7.0` + Pusher-js `^8.3.0`
  - **UI 与组件库**: Ant Design `^6.3.2` + ProComponents `^2.8.0` 
  - **CSS 与样式辅助**: TailwindCSS `^3.4.3` + Less `^4.2.0`
  - **国际化 (18n)**: i18next `^23.11.0` + react-i18next `^14.1.0`
  - **编辑器插件**: `@uiw/react-codemirror` `^4.23.0` + `react-quill` `^2.0.0`
  - **构建与质量规范**: ESLint `^9.0.0` + Prettier `^3.2.0` + Yarn `1.22.x` (或 4.x Berry)
- [x] 按照确定的版本清单安装核心依赖。
- [x] 配置代码规范环境 (`ESLint`、`Prettier` 相关规则)。
- [x] 搭建项目标准目录结构 (`assets`, `components`, `configs`, `hooks`, `pages`, `services`, `stores`, `types`, `utils`)。
- [x] 配置 `Axios` 拦截器 (统一处理 Token 携带、过期刷新、请求响应错误提示)。
- [x] 实现基础系统布局 (`ProLayout`)：侧边双栏菜单、顶级导航、用户头像、面包屑。
- [x] 配置初始化多语言环境 (`i18next`)，支持中英文切换。

## 阶段二：认证与看板核心模块 (P0)
- [x] **登录认证 (login.go)** 
  - [x] 登录页面开发及交互。
  - [x] 集成 `/login`、`/logout`、`/refresh` 接口。
  - [x] 使用 `Zustand` 管理全局应用鉴权状态及持久化。
- [x] **仪表板 (dashboard.go)**
  - [x] 开发工作台响应式图表与数据看板。
  - [x] 数据集成： `/dashboard/statistics`与 `/dashboard/loginChart`。

## 阶段三：RBAC 核心权限体系 (P0)
- [x] **菜单管理 (menu.go)**： 树状表格及 CRUD 处理。
- [x] **部门管理 (dept.go)**： 树状表格关联 CRUD，支持添加/分离部门负责人。
- [x] **角色管理 (role.go)**： 角色基础 CRUD，分配菜单使用树结构钩子、管理数据权限配置项。
- [x] **用户管理 (user.go)**： 人员列表展现、分配角色/部门视图、在线人员与强退操作融合、管理被删至回收站的用户。

## 阶段四：扩展系统管理功能 (P1)
- [ ] **岗位管理 (post.go)**：列表展示、绑定及 CRUD，支持回收站。
- [ ] **接口管理 (api.go)**：服务端点查看及其 CRUD 管理。
- [ ] **应用管理 (app.go)**：接入方应用列表及 API 端点映射关联绑定。
- [ ] **数据字典 (dict.go)**：关联型 CRUD。左侧列表查询类型 (`dictType`) ，右边明细操作 (`dataDict`)，以及数据字典缓存控制。
- [ ] **系统配置 (config.go)**：通过不同的配置分组（如基础设置、存储设置等），提供通用键值对读取及批量保存能力。
- [ ] **定时任务 (crontab.go)**：定时任务表达式编排及在线执行日志记录查看。
- [ ] **Pusher 服务扩展 (pusher_*.go)**：
  - [ ] 新建 `utils/pusher.ts` 服务端通道认证 (`pusher-js` auth endpoint 联调)。
  - [ ] 维护 WebSocket 连接、事件下发调试以及在线态 (Presence) 展示页。

## 阶段五：组件化及次要功能 (P2)
- [ ] **文件引擎与附件 (upload.go / attachment.go)**：
  - [ ] 基于 Ant Design 封装上传组件适配分片和常规图片上传。
  - [ ] 管理展示本地/云端的所有附件清单与销毁支持。
- [ ] **消息与通知中心 (notice.go / message.go)**：全局角标提醒与消息列表展示 (接合 WebSocket)。
- [ ] **日志管理 (logs.go)**：按表查询系统操作日志、请求日志报表视图。
- [ ] **缓存管理 (cache.go)**：简易监控看板：统计Redis情况与根据Key精准清理缓存记录。

## 阶段六：维护与配置 (P3)
- [ ] **数据维护 (data_maintain.go)**：快捷的单表维护能力面板。
- [ ] **系统模块 (system_modules.go)**：后端架构系统内部模块管理视窗。

## 注意事项
- 不得直接修改后端 `modules/system/api/` 的源码，接口变更由后端推动并同步对应 Swagger/定义。
- 表格数据尽量全系通过封装好的 `@ant-design/pro-table` / `pro-form` 控制布局减少手写样板代码。
- 采用明确的 TS 类型接口声明 (基于后端 `xx.go` 的响应参数)。