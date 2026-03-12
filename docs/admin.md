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
| WebSocket | Pusher-js | 8.3.0 |
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
┌─────────────────────────────────────────────────────────────────┐
│  Logo    |  面包屑          搜索    通知    用户头像              │ Header
├──────┬─────────┬─────────────────────────────────────────────────┤
│      │         │                                                 │
│ 一级 │  二级   │                   内容区域                        │
│ 菜单 │  菜单   │                                                 │
│      │         │                                                 │
│ 侧边栏│ 辅助栏  │                                                 │
│      │         │                                                 │
└──────┴─────────┴─────────────────────────────────────────────────┘
```

- 使用 `ProLayout` 实现三栏布局
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

### 6. 岗位管理 (post.go)
| 接口 | 方法 | 功能 |
|------|------|------|
| /post/index | GET | 岗位列表 |
| /post/recycle | GET | 回收站岗位列表 |
| /post/list | GET | 前端选择树 |
| /post/save | POST | 新增岗位 |
| /post/read/{Id} | GET | 获取岗位信息 |
| /post/update/{Id} | PUT | 更新岗位 |
| /post/delete | DELETE | 删除岗位 |
| /post/realDelete | DELETE | 真实删除 |
| /post/recovery | PUT | 恢复回收站数据 |
| /post/changeStatus | PUT | 更改状态 |
| /post/numberOperation | PUT | 数字运算操作 |

### 7. 接口管理 (api.go)
| 接口 | 方法 | 功能 |
|------|------|------|
| /api/index | GET | 接口列表 |
| /api/list | GET | 列表（无分页） |
| /api/recycle | GET | 回收站接口列表 |
| /api/save | POST | 新增接口 |
| /api/read/{Id} | GET | 获取接口信息 |
| /api/update/{Id} | PUT | 更新接口 |
| /api/delete | DELETE | 删除接口 |
| /api/realDelete | DELETE | 真实删除 |
| /api/recovery | PUT | 恢复回收站数据 |
| /api/changeStatus | PUT | 更改状态 |

### 8. 应用管理 (app.go)
| 接口 | 方法 | 功能 |
|------|------|------|
| /app/index | GET | 应用列表 |
| /app/recycle | GET | 回收站应用列表 |
| /app/save | POST | 新增应用 |
| /app/read/{Id} | GET | 获取应用信息 |
| /app/update/{Id} | PUT | 更新应用 |
| /app/bind/{Id} | PUT | 绑定API接口 |
| /app/delete | DELETE | 删除应用 |
| /app/realDelete | DELETE | 真实删除 |
| /app/recovery | PUT | 恢复回收站数据 |
| /app/changeStatus | PUT | 更改状态 |
| /app/getAppId | GET | 获取应用ID |
| /app/getAppSecret | GET | 获取应用秘钥 |
| /app/getApiList | GET | 获取绑定接口列表 |

### 9. 数据字典 (dict.go)
#### 字典类型 (dictType)
| 接口 | 方法 | 功能 |
|------|------|------|
| /dictType/index | GET | 字典类型列表 |
| /dictType/list | GET | 字典列表 |
| /dictType/recycle | GET | 回收站列表 |
| /dictType/save | POST | 新增字典类型 |
| /dictType/read/{Id} | GET | 获取字典类型数据 |
| /dictType/update/{Id} | PUT | 更新字典类型 |
| /dictType/delete | DELETE | 删除字典类型 |
| /dictType/realDelete | DELETE | 真实删除 |
| /dictType/recovery | PUT | 恢复回收站数据 |
| /dictType/changeStatus | PUT | 更改状态 |

#### 字典数据 (dataDict)
| 接口 | 方法 | 功能 |
|------|------|------|
| /dataDict/index | GET | 字典数据列表 |
| /dataDict/list | GET | 快捷查询字典 |
| /dataDict/lists | GET | 批量查询字典 |
| /dataDict/clearCache | POST | 清除字典缓存 |
| /dataDict/recycle | GET | 回收站列表 |
| /dataDict/save | POST | 新增字典数据 |
| /dataDict/read/{Id} | GET | 获取字典数据 |
| /dataDict/update/{Id} | PUT | 更新字典数据 |
| /dataDict/delete | DELETE | 删除字典数据 |
| /dataDict/realDelete | DELETE | 真实删除 |
| /dataDict/recovery | PUT | 恢复回收站数据 |
| /dataDict/changeStatus | PUT | 更改状态 |
| /dataDict/numberOperation | PUT | 数字运算操作 |

### 10. 系统配置 (config.go)
#### 配置组 (configGroup)
| 接口 | 方法 | 功能 |
|------|------|------|
| /setting/configGroup/index | GET | 获取系统组配置 |
| /setting/configGroup/save | POST | 保存配置组 |
| /setting/configGroup/update | POST | 更新配置组 |
| /setting/configGroup/delete | DELETE | 删除配置组 |

#### 配置项 (config)
| 接口 | 方法 | 功能 |
|------|------|------|
| /setting/config/index | GET | 获取配置列表 |
| /setting/config/save | POST | 保存配置 |
| /setting/config/update | POST | 更新配置 |
| /setting/config/updateByKeys | POST | 按keys批量更新配置 |
| /setting/config/delete | DELETE | 删除配置 |

### 11. 定时任务 (crontab.go)
| 接口 | 方法 | 功能 |
|------|------|------|
| /setting/crontab/index | GET | 定时任务列表 |
| /setting/crontab/logPageList | GET | 日志列表 |
| /setting/crontab/save | POST | 保存定时任务 |
| /setting/crontab/read/{Id} | GET | 获取定时任务详情 |
| /setting/crontab/update/{Id} | PUT | 更新定时任务 |
| /setting/crontab/delete | DELETE | 删除定时任务 |
| /setting/crontab/deleteCrontabLog | DELETE | 删除定时任务日志 |
| /setting/crontab/run | POST | 立即执行定时任务 |
| /setting/crontab/changeStatus | PUT | 更改状态 |
| /setting/crontab/getTarget | GET | 获取执行目标 |

### 12. 文件上传 (upload.go)
| 接口 | 方法 | 功能 |
|------|------|------|
| /uploadFile | POST | 上传文件 |
| /uploadImage | POST | 上传图片 |
| /chunkUpload | POST | 分片上传文件 |
| /saveNetworkImage | POST | 保存网络图片 |
| /getFileInfoById | GET | 通过ID获取文件信息 |
| /getFileInfoByHash | GET | 通过Hash获取文件信息 |
| /downloadById | GET | 根据ID下载文件 |
| /downloadByHash | GET | 根据Hash下载文件 |
| /showFile/{Hash} | GET | 输出图片/文件 |

### 13. 附件管理 (attachment.go)
| 接口 | 方法 | 功能 |
|------|------|------|
| /attachment/index | GET | 附件列表 |
| /attachment/recycle | GET | 回收站列表 |
| /attachment/delete | DELETE | 删除附件 |
| /attachment/realDelete | DELETE | 真实删除 |
| /attachment/recovery | PUT | 恢复回收站数据 |

### 14. 通知管理 (notice.go)
| 接口 | 方法 | 功能 |
|------|------|------|
| /notice/index | GET | 通知列表 |
| /notice/recycle | GET | 回收站通知列表 |
| /notice/save | POST | 新增通知 |
| /notice/read/{Id} | GET | 获取通知信息 |
| /notice/update/{Id} | PUT | 更新通知 |
| /notice/delete | DELETE | 删除通知 |
| /notice/realDelete | DELETE | 真实删除 |
| /notice/recovery | PUT | 恢复回收站数据 |

### 15. 消息中心 (message.go)
| 接口 | 方法 | 功能 |
|------|------|------|
| /queueMessage/receiveList | GET | 接收消息列表 |
| /queueMessage/updateReadStatus | PUT | 更新已读状态 |
| /queueMessage/deletes | DELETE | 删除消息 |

### 16. 日志管理 (logs.go)
| 接口 | 方法 | 功能 |
|------|------|------|
| /logs/getLoginLogPageList | GET | 登录日志列表 |
| /logs/getOperLogPageList | GET | 操作日志列表 |
| /logs/getApiLogPageList | GET | 接口日志列表 |
| /logs/deleteLoginLog | DELETE | 删除登录日志 |
| /logs/deleteOperLog | DELETE | 删除操作日志 |
| /logs/deleteApiLog | DELETE | 删除接口日志 |

### 17. 缓存管理 (cache.go)
| 接口 | 方法 | 功能 |
|------|------|------|
| /cache/monitor | GET | 缓存信息监控 |
| /cache/view | POST | 查看key内容 |
| /cache/delete | DELETE | 根据key删除缓存 |
| /cache/clear | DELETE | 清空所有缓存 |

### 18. 仪表板 (dashboard.go)
| 接口 | 方法 | 功能 |
|------|------|------|
| /dashboard/statistics | GET | 获取仪表板统计数据 |
| /dashboard/loginChart | GET | 获取登录图表数据 |

### 19. 数据维护 (data_maintain.go)
| 接口 | 方法 | 功能 |
|------|------|------|
| /dataMaintain/index | GET | 数据表维护列表 |

### 20. 系统模块 (system_modules.go)
| 接口 | 方法 | 功能 |
|------|------|------|
| /systemModules/index | GET | 模块分页列表 |
| /systemModules/list | GET | 模块列表 |
| /systemModules/save | POST | 新增模块 |
| /systemModules/read/{Id} | GET | 获取模块信息 |
| /systemModules/update/{Id} | PUT | 更新模块 |
| /systemModules/delete | DELETE | 删除模块 |
| /systemModules/recycle | GET | 回收站列表 |
| /systemModules/realDelete | DELETE | 真实删除 |
| /systemModules/recovery | PUT | 恢复回收站数据 |
| /systemModules/changeStatus | PUT | 更改状态 |

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

### 21. Pusher 频道认证 (pusher_auth.go)
| 接口 | 方法 | 功能 |
|------|------|------|
| /pusher/auth | POST/GET | 频道认证 |
| /pusher/auth/batch | POST/GET | 批量频道认证 |

### 22. Pusher 频道管理 (pusher_channel.go)
| 接口 | 方法 | 功能 |
|------|------|------|
| /apps/:app_id/channels | GET | 获取频道列表 |
| /apps/:app_id/channels/:channel_name | GET | 获取单个频道信息 |
| /apps/:app_id/channels/:channel_name/users | GET | 获取Presence频道用户列表 |

### 23. Pusher 事件推送 (pusher_events.go)
| 接口 | 方法 | 功能 |
|------|------|------|
| /apps/:app_id/events | POST | 推送事件到频道 |
| /apps/:app_id/batch_events | POST | 批量推送事件 |
| /apps/:app_id/users/:user_id/events | POST | 向特定用户发送事件 |
| /apps/:app_id/users/:user_id/terminate_connections | POST | 终止用户的所有连接 |

### 24. Pusher Webhook (pusher_webhook.go)
| 接口 | 方法 | 功能 |
|------|------|------|
| /apps/:app_id/webhooks | POST | Pusher Webhook 验证 |

### 25. Pusher 用户认证 (pusher_user_auth.go)
| 接口 | 方法 | 功能 |
|------|------|------|
| /pusher/user-auth | POST/GET | Pusher 用户身份认证 |

### 事件类型
- `pusher:connection_established` - 连接建立
- `pusher:ping` / `pusher:pong` - 心跳
- `pusher:subscribe` / `pusher:unsubscribe` - 订阅/取消订阅
- `pusher:member_added` / `pusher:member_removed` - 成员加入/离开

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
