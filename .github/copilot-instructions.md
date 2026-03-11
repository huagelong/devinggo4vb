# DevingGo-Light 项目规范

**GoFrame v2** 后端项目，模块化插件架构，PostgreSQL + Redis + JWT + RBAC。

## 架构

请求流：`Router → Middleware → Controller → Service（接口） → Logic（实现） → DAO → PostgreSQL`

- `internal/` — 框架核心层（dao/model/service），**禁止手动修改** `internal/` 子目录和 `modules/*/service/`
- `modules/{module}/` — 业务模块：`api/` `controller/` `logic/` `model/req/` `model/res/` `router/` `worker/`
- `modules/bootstrap/` — blank import 触发各模块 `init()`

## 编码规范

**命名**：Logic 结构体 `s` 前缀（`sSystemUser`）；Controller 包级变量帕斯卡命名（`var UserController = userController{}`）

**API 层**（`api/{module}/{resource}.go`）：
```go
type IndexUserReq struct {
    g.Meta `path:"/user/index" method:"get" tags:"用户管理" summary:"用户列表" x-permission:"system:user:index"`
    model.AuthorHeader  // 所有需认证接口必须内嵌
    model.PageListReq   // 分页接口内嵌
    req.SystemUserSearch
}
type IndexUserRes struct {
    g.Meta `mime:"application/json"`
    page.PageRes
    Items []res.SystemUser `json:"items"`
}
```
> 权限标签：`x-permission:"模块:资源:操作"` / `x-exceptAuth:"true"`（免权限）/ `x-exceptLogin:"true"`（公开）

**Controller 层**：只调 Service，不写业务逻辑，签名 `(ctx, *XxxReq) (*XxxRes, error)`

**Logic 层**：`init()` 注册实现；必须实现 `Model(ctx)` 方法：
```go
func init() { service.RegisterSystemUser(NewSystemUser()) }

func (s *sSystemUser) Model(ctx context.Context) *gdb.Model {
    return dao.SystemUser.Ctx(ctx).Hook(hook.Default()).Cache(orm.SetCacheOption(ctx))
}
```

**ORM**：用 `dao.Xxx.Columns().Field` 引用列名；用 `do.Xxx{}` 写入数据（避免零值覆盖）；分页用 `orm.NewQuery(s.Model(ctx)).WithPageListReq(req).ScanAndCount(&items, &total)`

**路由**：`group.Bind(Controller).Middleware(service.Middleware().AdminAuth)`

## 常用命令

| 命令 | 说明 |
|------|------|
| `make run` | 生成代码后启动 |
| `make dao` | 重新生成 DAO/Entity/DO |
| `make service` | 重新生成 Service 接口 |
| `make gen-crud table=xxx` | 生成 CRUD 代码 |
| `make gen-module name=xxx` | 创建新模块 |

## 数据库

PostgreSQL 13+，主键 `id`（int64，雪花算法），公共字段 `created_at/updated_at/deleted_at/created_by/updated_by` 由 Hook 自动填充，迁移文件在 `resource/migrations/`。
