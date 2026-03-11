---
name: crud-development
description: "使用场景：在 DevingGo-Light 项目中开发新的 CRUD 功能模块。当用户想要为新数据库表创建增删改查接口、或需要手动编写 CRUD 代码（含 API/Controller/Logic/Model 等各层）时使用。包含完整工作流：生成代码 → 完善业务逻辑 → 注册路由 → 运行 make service。"
---

# CRUD 功能开发工作流

## 概述

在 DevingGo-Light 中开发新的 CRUD 功能，需要按照分层架构依次创建各层代码，或使用代码生成器自动化生成。

## 工作流步骤

### 阶段一：准备数据库表

1. **创建迁移文件**：
   ```bash
   go run main.go migrate:create -name add_your_table
   ```
   
2. **编写 SQL**（`resource/migrations/{timestamp}_add_your_table.up.sql`）：
   ```sql
   CREATE TABLE {module}_{resource} (
       id          bigint       NOT NULL,
       name        varchar(100) NOT NULL DEFAULT '',
       status      int2         NOT NULL DEFAULT 1,
       remark      varchar(500) NOT NULL DEFAULT '',
       created_at  timestamp,
       updated_at  timestamp,
       deleted_at  timestamp,
       created_by  bigint,
       updated_by  bigint,
       CONSTRAINT {module}_{resource}_pkey PRIMARY KEY (id)
   );
   COMMENT ON TABLE {module}_{resource} IS '资源说明';
   ```

3. **执行迁移**：
   ```bash
   go run main.go migrate:up
   ```

4. **重新生成 DAO/Entity/DO**：
   ```bash
   make dao
   ```

### 阶段二：生成代码骨架

**方式 A：使用代码生成器（推荐）**

```bash
make gen-crud table={module}_{resource}
# 例：make gen-crud table=system_post
```

生成器会自动创建：
- `modules/{module}/api/{module}/{varName}.go` — API Req/Res 定义
- `modules/{module}/model/req/{table}.go` — 请求搜索模型
- `modules/{module}/model/res/{table}.go` — 响应模型
- `modules/{module}/controller/{module}/{varName}.go` — 控制器
- `modules/{module}/logic/{module}/{table}.go` — 业务逻辑

**方式 B：手动编写**（参照下方模板）

### 阶段三：编写各层代码

#### 1. API 定义（`modules/{module}/api/{module}/{resource}.go`）

```go
package system

import (
    "devinggo/modules/system/model/req"
    "devinggo/modules/system/model/res"
    "devinggo/modules/system/pkg/page"
    "github.com/gogf/gf/v2/frame/g"
    "devinggo/internal/model"
)

// 列表接口
type Index{Resource}Req struct {
    g.Meta `path:"/{resource}/index" method:"get" tags:"{中文分组}" summary:"{资源}列表" x-permission:"{module}:{resource}:index"`
    model.AuthorHeader
    model.PageListReq
    req.{PascalResource}Search
}
type Index{Resource}Res struct {
    g.Meta `mime:"application/json"`
    page.PageRes
    Items []res.{PascalResource} `json:"items"`
}

// 详情接口
type Read{Resource}Req struct {
    g.Meta `path:"/{resource}/read" method:"get" tags:"{中文分组}" summary:"获取{中文名}详情" x-permission:"{module}:{resource}:read"`
    model.AuthorHeader
    Id int64 `json:"id" v:"required#id不能为空"`
}
type Read{Resource}Res struct {
    g.Meta `mime:"application/json"`
    res.{PascalResource}
}

// 创建接口
type Create{Resource}Req struct {
    g.Meta `path:"/{resource}/create" method:"post" tags:"{中文分组}" summary:"创建{中文名}" x-permission:"{module}:{resource}:create"`
    model.AuthorHeader
    req.{PascalResource}Create
}
type Create{Resource}Res struct {
    g.Meta `mime:"application/json"`
}

// 更新接口
type Update{Resource}Req struct {
    g.Meta `path:"/{resource}/update" method:"put" tags:"{中文分组}" summary:"更新{中文名}" x-permission:"{module}:{resource}:update"`
    model.AuthorHeader
    Id int64 `json:"id" v:"required#id不能为空"`
    req.{PascalResource}Update
}
type Update{Resource}Res struct {
    g.Meta `mime:"application/json"`
}

// 删除接口
type Delete{Resource}Req struct {
    g.Meta `path:"/{resource}/delete" method:"delete" tags:"{中文分组}" summary:"删除{中文名}" x-permission:"{module}:{resource}:delete"`
    model.AuthorHeader
    Ids []int64 `json:"ids" v:"required#ids不能为空"`
}
type Delete{Resource}Res struct {
    g.Meta `mime:"application/json"`
}
```

#### 2. 请求/响应模型

**`modules/{module}/model/req/{table}.go`**：
```go
package req

type {PascalResource}Search struct {
    Name   string `json:"name"   dc:"名称"`
    Status int    `json:"status" dc:"状态"`
}

type {PascalResource}Create struct {
    Name   string `json:"name"   v:"required#名称不能为空" dc:"名称"`
    Status int    `json:"status" dc:"状态 1:启用 2:禁用"`
    Remark string `json:"remark" dc:"备注"`
}

type {PascalResource}Update struct {
    Name   string `json:"name"   dc:"名称"`
    Status int    `json:"status" dc:"状态"`
    Remark string `json:"remark" dc:"备注"`
}
```

**`modules/{module}/model/res/{table}.go`**：
```go
package res

type {PascalResource} struct {
    Id        int64  `json:"id"`
    Name      string `json:"name"`
    Status    int    `json:"status"`
    Remark    string `json:"remark"`
    CreatedAt string `json:"createdAt"`
    UpdatedAt string `json:"updatedAt"`
}
```

#### 3. Controller（`modules/{module}/controller/{module}/{resource}.go`）

```go
package system

import (
    "context"
    "devinggo/internal/service"
    system "devinggo/modules/system/api/{module}"
    "devinggo/modules/system/base"
)

type {resource}Controller struct{ base.BaseController }

var {PascalResource}Controller = {resource}Controller{}

func (c *{resource}Controller) Index(ctx context.Context, in *system.Index{Resource}Req) (out *system.Index{Resource}Res, err error) {
    out = &system.Index{Resource}Res{}
    out.PageRes, out.Items, err = service.{Module}{PascalResource}().GetPageList(ctx, in)
    return
}

func (c *{resource}Controller) Read(ctx context.Context, in *system.Read{Resource}Req) (out *system.Read{Resource}Res, err error) {
    out = &system.Read{Resource}Res{}
    err = service.{Module}{PascalResource}().GetById(ctx, in.Id, &out.{PascalResource})
    return
}

func (c *{resource}Controller) Create(ctx context.Context, in *system.Create{Resource}Req) (out *system.Create{Resource}Res, err error) {
    out = &system.Create{Resource}Res{}
    err = service.{Module}{PascalResource}().Create(ctx, in)
    return
}

func (c *{resource}Controller) Update(ctx context.Context, in *system.Update{Resource}Req) (out *system.Update{Resource}Res, err error) {
    out = &system.Update{Resource}Res{}
    err = service.{Module}{PascalResource}().Update(ctx, in)
    return
}

func (c *{resource}Controller) Delete(ctx context.Context, in *system.Delete{Resource}Req) (out *system.Delete{Resource}Res, err error) {
    out = &system.Delete{Resource}Res{}
    err = service.{Module}{PascalResource}().Delete(ctx, in.Ids)
    return
}
```

#### 4. Logic（`modules/{module}/logic/{module}/{table}.go`）

```go
package system

import (
    "context"
    "devinggo/internal/dao"
    "devinggo/internal/model/do"
    "devinggo/modules/system/base"
    "devinggo/modules/system/model/res"
    "devinggo/modules/system/pkg/hook"
    "devinggo/modules/system/pkg/orm"
    "devinggo/modules/system/service"
    system "devinggo/modules/system/api/{module}"
    "github.com/gogf/gf/v2/database/gdb"
)

func init() {
    service.Register{Module}{PascalResource}(New{Module}{PascalResource}())
}

type s{Module}{PascalResource} struct{ base.BaseService }

func New{Module}{PascalResource}() *s{Module}{PascalResource} {
    return &s{Module}{PascalResource}{}
}

func (s *s{Module}{PascalResource}) Model(ctx context.Context) *gdb.Model {
    return dao.{PascalTable}.Ctx(ctx).
        Hook(hook.Default()).
        Cache(orm.SetCacheOption(ctx))
}

func (s *s{Module}{PascalResource}) GetPageList(ctx context.Context, req *system.Index{Resource}Req) (pageRes interface{}, items []res.{PascalResource}, err error) {
    var list []*res.{PascalResource}
    m := s.Model(ctx)
    // 搜索条件
    if req.Name != "" {
        m = m.WhereLike(dao.{PascalTable}.Columns().Name, "%"+req.Name+"%")
    }
    if req.Status > 0 {
        m = m.Where(dao.{PascalTable}.Columns().Status, req.Status)
    }
    pageRes, err = orm.NewQuery(m).WithPageListReq(req).ScanAndCount(&list, nil)
    return pageRes, func() []res.{PascalResource} {
        result := make([]res.{PascalResource}, 0, len(list))
        for _, v := range list { result = append(result, *v) }
        return result
    }(), err
}

func (s *s{Module}{PascalResource}) GetById(ctx context.Context, id int64, out *res.{PascalResource}) error {
    return s.Model(ctx).Where(dao.{PascalTable}.Columns().Id, id).Scan(out)
}

func (s *s{Module}{PascalResource}) Create(ctx context.Context, req *system.Create{Resource}Req) error {
    _, err := s.Model(ctx).Data(do.{PascalTable}{
        Name:   req.Name,
        Status: req.Status,
        Remark: req.Remark,
    }).Insert()
    return err
}

func (s *s{Module}{PascalResource}) Update(ctx context.Context, req *system.Update{Resource}Req) error {
    _, err := s.Model(ctx).Data(do.{PascalTable}{
        Name:   req.Name,
        Status: req.Status,
        Remark: req.Remark,
    }).Where(dao.{PascalTable}.Columns().Id, req.Id).Update()
    return err
}

func (s *s{Module}{PascalResource}) Delete(ctx context.Context, ids []int64) error {
    _, err := s.Model(ctx).WhereIn(dao.{PascalTable}.Columns().Id, ids).Delete()
    return err
}
```

### 阶段四：生成 Service 接口

每次编写完 Logic 层后，**必须**重新生成 Service 接口：

```bash
make service
```

这会更新 `modules/{module}/service/` 下的接口定义文件。

### 阶段五：注册路由

在 `modules/{module}/router/{module}/router.go` 中添加 Controller 绑定：

```go
group.Group("/{module}", func(group *ghttp.RouterGroup) {
    group.Bind(
        system.{PascalResource}Controller,  // 添加这行
    ).Middleware(service.Middleware().AdminAuth)
})
```

### 阶段六：验证

```bash
make run
# 访问 Swagger UI 确认接口已注册：http://localhost:8000/swagger
```

## 变量命名对照表

| 占位符 | 说明 | 示例 |
|--------|------|------|
| `{module}` | 模块名（小写） | `system` |
| `{resource}` | 资源名（小写驼峰） | `post` |
| `{Resource}` | 资源名（帕斯卡） | `Post` |
| `{PascalResource}` | 资源名（帕斯卡） | `Post` |
| `{Module}` | 模块名（帕斯卡） | `System` |
| `{table}` | 数据库表名（下划线） | `system_post` |
| `{PascalTable}` | 表名（帕斯卡，DAO 类名） | `SystemPost` |
| `{中文名}` | 功能中文名 | `岗位` |

## 注意事项

- 禁止手动修改 `internal/dao/internal/`、`internal/model/do/`、`internal/model/entity/`、`modules/*/service/` 下的文件
- Logic 的 `init()` 函数中必须调用 `service.Register{X}()` 才能完成依赖注入
- `do.XxxTable{}` 中 `nil` 值字段不会被写入数据库，避免零值覆盖
- 软删除由框架钩子 `hook.Default()` 自动处理，无需手动过滤
