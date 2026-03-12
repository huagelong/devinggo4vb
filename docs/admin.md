# 管理后台前端项目需求分析

## 一、项目概述

基于后端 `modules/system/api/system` 的26个API功能模块，构建一个商业化管理后台前端。

**目标目录**: `admin-ui/`
**后端API**: `modules/system/api/system/`

---

## 二、功能模块清单

| 序号 | 模块名称 | 后端API文件 | 功能描述 | 优先级 |
|------|---------|------------|---------|-------|
| 1 | 登录认证 | login.go | 登录、登出、刷新Token | P0 |
| 2 | 用户管理 | user.go | 管理员CRUD、在线用户、回收站 | P0 |
| 3 | 角色管理 | role.go | 角色CRUD、权限分配、数据权限 | P0 |
| 4 | 菜单管理 | menu.go | 菜单树CRUD | P0 |
| 5 | 部门管理 | dept.go | 部门树CRUD、部门领导 | P0 |
| 6 | 岗位管理 | post.go | 岗位CRUD | P1 |
| 7 | 接口管理 | api.go | API接口CRUD | P1 |
| 8 | 应用管理 | app.go | 应用CRUD、API绑定 | P1 |
| 9 | 数据字典 | dict.go | 字典类型和数据CRUD | P1 |
| 10 | 系统配置 | config.go | 配置组和配置项管理 | P1 |
| 11 | 定时任务 | crontab.go | 定时任务CRUD、执行日志 | P1 |
| 12 | 文件上传 | upload.go | 文件上传、分片上传、下载 | P2 |
| 13 | 附件管理 | attachment.go | 附件CRUD | P2 |
| 14 | 通知管理 | notice.go | 通知CRUD | P2 |
| 15 | 消息中心 | message.go | 消息接收、已读状态 | P2 |
| 16 | 日志管理 | logs.go | 登录日志、操作日志、API日志 | P2 |
| 17 | 缓存管理 | cache.go | 缓存监控、查看、清理 | P2 |
| 18 | 仪表板 | dashboard.go | 统计数据、图表 | P0 |
| 19 | 数据维护 | data_maintain.go | 数据表维护 | P3 |
| 20 | 系统模块 | system_modules.go | 模块管理 | P3 |
| 21-26 | Pusher相关 | pusher_*.go | WebSocket认证、频道、事件 | P1 |

---

## 三、技术栈

| 类别 | 技术 | 版本 |
|------|------|------|
| 框架 | Vite + React 18 | - |
| UI库 | Ant Design | 6.3.2 |
| 组件库 | ProComponents | 2.8.0 |
| 路由 | TanStack Router | 1.80.0 |
| 状态管理 | Zustand | 5.0.0 |
| HTTP | Axios | 1.7.0 |
| WebSocket | Pusher-js | 8.4.0 |
| 国际化 | i18next + react-i18next | - |
| 代码编辑 | @uiw/react-codemirror | 4.23.0 |
| 富文本 | react-quill | 2.0.0 |
| 代码规范 | ESLint + Prettier | - |
| 样式 | TailwindCSS + Less | - |
| 包管理 | Yarn | - |

---

## 四、布局设计

### 布局结构
```
┌─────────────────────────────────────────────────────┐
│  Logo    |  面包屑          搜索    通知    用户头像  │ Header
├─────────┬───────────────────────────────────────────┤
│         │                                           │
│  菜单   │           内容区域                         │
│         │                                           │
│         │                                           │
│  侧边栏  │                                           │
│         │                                           │
└─────────┴───────────────────────────────────────────┘
```

- 使用 `ProLayout` 实现左侧两栏布局
- 左侧菜单从后端API动态获取
- 支持折叠/展开、面包屑、页头
- 主题切换（亮色/暗色）

---

## 五、API接口清单

### 1. 登录认证 (login.go)
| 接口 | 方法 | 功能 |
|------|------|------|
| /login | POST | 用户登录 |
| /logout | POST | 退出登录 |
| /refresh | POST | 刷新Token |

### 2. 用户管理 (user.go)
| 接口 | 方法 | 功能 |
|------|------|------|
| /getInfo | GET | 获取登录管理员信息 |
| /user/updateInfo | POST | 更新管理员信息 |
| /user/modifyPassword | POST | 修改密码 |
| /user/index | GET | 管理员信息列表 |
| /user/recycle | GET | 回收站管理员信息列表 |
| /user/save | POST | 新增管理员 |
| /user/read/{Id} | GET | 获取管理员详情 |
| /user/update/{Id} | PUT | 更新管理员 |
| /user/delete | DELETE | 删除管理员 |
| /user/realDelete | DELETE | 真实删除 |
| /user/recovery | PUT | 恢复回收站数据 |
| /user/changeStatus | PUT | 更改状态 |
| /onlineUser/index | GET | 获取在线用户列表 |
| /onlineUser/kick | POST | 强退用户 |

### 3. 角色管理 (role.go)
| 接口 | 方法 | 功能 |
|------|------|------|
| /role/index | GET | 角色列表 |
| /role/save | POST | 新增角色 |
| /role/update/{Id} | PUT | 更新角色 |
| /role/delete | DELETE | 删除角色 |
| /role/menuPermission/{Id} | PUT | 更新用户菜单权限 |
| /role/dataPermission/{Id} | PUT | 更新用户数据权限 |
| /role/getMenuByRole/{Id} | GET | 通过角色获取菜单 |
| /role/getDeptByRole/{Id} | GET | 通过角色获取部门 |

### 4. 菜单管理 (menu.go)
| 接口 | 方法 | 功能 |
|------|------|------|
| /menu/index | GET | 菜单树列表 |
| /menu/tree | GET | 前端选择树 |
| /menu/save | POST | 新增菜单 |
| /menu/update/{Id} | PUT | 更新菜单 |
| /menu/delete | DELETE | 删除菜单 |

### 5. 部门管理 (dept.go)
| 接口 | 方法 | 功能 |
|------|------|------|
| /dept/index | GET | 部门树列表 |
| /dept/tree | GET | 前端选择树 |
| /dept/save | POST | 新增部门 |
| /dept/update/{Id} | PUT | 更新部门 |
| /dept/delete | DELETE | 删除部门 |
| /dept/addLeader | POST | 新增部门领导 |
| /dept/delLeader | DELETE | 删除部门领导 |

### 6. 其他模块接口
- 岗位、接口、应用、字典、配置、定时任务、上传、附件、通知、消息、日志、缓存、仪表板等

---

## 六、WebSocket (Pusher)

### 后端实现位置
- `modules/system/pkg/websocket/`

### 支持的频道类型
| 类型 | 前缀 | 说明 |
|------|------|------|
| Public | 无前缀 | 公开频道 |
| Private | `private-*` | 私有频道 |
| Presence | `presence-*` | 在线状态频道 |
| Encrypted | `private-encrypted-*` | 加密频道 |

### 事件类型
- `pusher:connection_established` - 连接建立
- `pusher:ping` / `pusher:pong` - 心跳
- `pusher:subscribe` / `pusher:unsubscribe` - 订阅/取消订阅
- `pusher:member_added` / `pusher:member_removed` - 成员加入/离开

### 认证接口
- `/pusher/auth` - 频道认证
- `/pusher/auth/batch` - 批量频道认证
- `/pusher/user-auth` - 用户认证

---

## 七、项目目录结构

```
admin-ui/
├── public/
│   └── favicon.svg
├── src/
│   ├── assets/              # 静态资源
│   │   └── styles/
│   ├── components/          # 公共组件
│   │   ├── common/          # 通用组件
│   │   ├── form/            # 表单组件
│   │   └── layout/          # 布局组件
│   ├── configs/             # 配置
│   │   └── routes/          # 路由配置
│   ├── hooks/               # 自定义Hooks
│   ├── i18n/                # 国际化
│   │   └── locales/
│   ├── pages/               # 页面
│   │   ├── login/
│   │   ├── dashboard/
│   │   └── system/          # 系统管理模块
│   ├── services/            # API服务
│   ├── stores/              # Zustand状态
│   ├── types/               # TypeScript类型
│   ├── utils/               # 工具函数
│   ├── App.tsx
│   ├── main.tsx
│   └── router.tsx
├── .eslintrc.cjs
├── .prettierrc
├── tailwind.config.js
├── tsconfig.json
├── vite.config.ts
└── package.json
```

---

## 八、UI/UX要求

1. **商业化设计**: 使用ProComponents开箱即用的商业组件
2. **左侧两栏布局**: ProLayout实现
3. **国际化支持**: 中文/英文切换
4. **响应式设计**: 适配桌面端
5. **优雅交互**: 平滑过渡动画

---

## 九、开发规范

1. 代码注释详细，使用中文
2. 代码结构清晰，易于阅读
3. 不修改后端代码
4. 使用TypeScript类型检查
5. 遵循ESLint规则
