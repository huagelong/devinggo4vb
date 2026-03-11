---
name: new-api-endpoint
description: "在 DevingGo-Light 项目中为现有模块添加单个新 API 接口。适用于：在已有 Controller/Logic 中添加一个新的业务接口（如导出、统计、审批等），不涉及新建完整 CRUD。"
---

# 添加新 API 接口

为 `${module}` 模块的 `${resource}` 资源添加 `${action}` 接口（`${description}`）。

## 三步操作

**1. API 定义** — `modules/${module}/api/${module}/${resource}.go`

```go
type ${Action}${Resource}Req struct {
    g.Meta `path:"/${resource}/${action}" method:"${method}" tags:"${tag}" summary:"${description}" x-permission:"${module}:${resource}:${action}"`
    model.AuthorHeader
    // 参数字段
}
type ${Action}${Resource}Res struct {
    g.Meta `mime:"application/json"`
    // 响应字段
}
```

> 权限标签三选一：`x-permission:"模块:资源:操作"` / `x-exceptAuth:"true"`（免权限）/ `x-exceptLogin:"true"`（公开）

**2. Controller** — `modules/${module}/controller/${module}/${resource}.go`

```go
func (c *${resource}Controller) ${Action}(ctx context.Context, in *${module}.${Action}${Resource}Req) (out *${module}.${Action}${Resource}Res, err error) {
    out = &${module}.${Action}${Resource}Res{}
    err = service.${Module}${Resource}().${Action}(ctx, in)
    return
}
```

**3. Logic** — `modules/${module}/logic/${module}/${table}.go`

```go
func (s *s${Module}${Resource}) ${Action}(ctx context.Context, req *...) error {
    // 业务逻辑
    return nil
}
```

完成后执行 `make service` 重新生成 Service 接口。
