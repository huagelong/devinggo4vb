---
name: new-module
description: "使用场景：在 DevingGo-Light 项目中创建全新的业务模块。当用户需要添加新的功能模块（如电商模块、CMS模块、工单模块等）而不是在现有 system 模块中添加功能时使用。涵盖完整的模块搭建流程：目录结构、模块注册、路由配置、bootstrap 绑定。"
---

# 新模块创建工作流

## 概述

DevingGo-Light 采用插件化模块架构，每个模块通过 `blank import` 自动注册到系统中，无需修改核心代码。

## 快捷命令

```bash
# 自动创建模块目录结构（推荐）
make gen-module name={module_name}

# 克隆已有模块（适合快速开始）
make clone-module name={new_module} source=system
```

## 手动创建步骤

### 第一步：目录结构

创建以下目录结构（以 `cms` 模块为例）：

```
modules/cms/
├── module.go              # 模块入口（必须）
├── api/
│   └── cms/              # API Req/Res 定义
├── controller/
│   └── cms/              # HTTP 控制器
├── logic/
│   └── cms/              # 业务逻辑
├── service/              # Service 接口（gf gen service 生成）
├── model/
│   ├── req/              # 请求模型
│   └── res/              # 响应模型
├── router/
│   └── cms/
│       └── router.go     # 路由注册
├── base/                # Controller/Logic 基类
├── consts/              # 常量定义
├── codes/               # 响应码
└── myerror/             # 自定义错误
```

### 第二步：模块入口文件

**`modules/cms/module.go`**：
```go
package cms

import (
    "context"
    "devinggo/modules"
    router "devinggo/modules/cms/router/cms"
    "github.com/gogf/gf/v2/net/ghttp"
)

func init() {
    m := &cmsModule{}
    m.Name = "cms"
    modules.Register(m)
}

type cmsModule struct {
    modules.BaseModule
}

func (m *cmsModule) Start(ctx context.Context, s *ghttp.Server) error {
    s.Group("/", func(group *ghttp.RouterGroup) {
        router.BindController(group)
    })
    return nil
}
```

### 第三步：路由注册

**`modules/cms/router/cms/router.go`**：
```go
package cms

import (
    "devinggo/internal/service"
    cms "devinggo/modules/cms/controller/cms"
    "github.com/gogf/gf/v2/net/ghttp"
)

func BindController(group *ghttp.RouterGroup) {
    // 公开接口（无需认证）
    // group.Bind(cms.PublicController)

    // 需要认证的接口
    group.Group("/cms", func(group *ghttp.RouterGroup) {
        group.Bind(
            // cms.ArticleController,
        ).Middleware(service.Middleware().AdminAuth)
    })
}
```

### 第四步：Bootstrap 注册

在 `modules/bootstrap/` 目录下注册模块：

**`modules/bootstrap/modules/modules.go`**（在已有文件中添加 import）：
```go
import (
    _ "devinggo/modules/system"
    _ "devinggo/modules/api"
    _ "devinggo/modules/cms"   // 添加新模块
)
```

**`modules/bootstrap/logic/logic.go`**（在已有文件中添加 import）：
```go
import (
    _ "devinggo/modules/system/logic/system"
    _ "devinggo/modules/cms/logic/cms"   // 添加新模块
)
```

### 第五步：配置 gf gen service

在 `hack/config.yaml` 的 `gfcli.gen.service` 部分添加新模块配置：

```yaml
gfcli:
  gen:
    service:
      - srcFolder: "modules/system/logic"
        dstFolder: "modules/system/service"
        watchFile: "modules/system/logic"
        # ...
      - srcFolder: "modules/cms/logic"      # 新增
        dstFolder: "modules/cms/service"    # 新增
        watchFile: "modules/cms/logic"      # 新增
```

### 第六步：生成 Service 接口

编写 Logic 后，生成 Service 接口：

```bash
make service
```

### 第七步：基础设施文件（选择性创建）

**基类 Controller（`modules/cms/base/controller.go`）**：
```go
package base

import "devinggo/modules/system/base"

// 直接复用 system 模块的基类，或自定义
type BaseController struct {
    base.BaseController
}
```

**基类 Service（`modules/cms/base/service.go`）**：
```go
package base

import "devinggo/modules/system/base"

type BaseService struct {
    base.BaseService
}
```

**自定义错误码（`modules/cms/codes/codes.go`）**：
```go
package codes

import "github.com/gogf/gf/v2/errors/gcode"

var (
    CodeCmsArticleNotFound = gcode.New(10001, "文章不存在", nil)
)
```

### 第八步：验证模块加载

```bash
make run
# 检查启动日志中是否有 "cms module started" 或类似信息
# 访问 http://localhost:8000/swagger 确认路由已注册
```

## 注意事项

- 每个模块的 `module.go` 中 `m.Name` 必须唯一
- Bootstrap 中的 `blank import` 顺序不影响功能，但建议按字母排序
- 模块间不应直接调用对方的 Logic，应通过 `internal/service` 接口通信
- 新模块的 `service/` 目录由 `gf gen service` 自动生成，**禁止手动修改**
- 模块的路由前缀建议与模块名一致（如 `/cms/...`）以避免冲突
