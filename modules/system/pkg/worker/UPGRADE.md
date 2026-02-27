# Worker 包升级指南

## 版本信息

- **优化日期**: 2026-02-27
- **版本**: v2.0

## 主要变更

### ✨ 新增功能

1. **统一的Manager管理器** (`worker.go`)
   - 提供一站式的Worker管理接口
   - 支持链式调用
   - 统一注册和运行Worker/Cron任务

2. **TaskBuilder构建器** (`builder.go`)
   - 链式API构建任务
   - 支持灵活配置队列、延迟、任务ID等参数
   - 简化任务发送流程

3. **类型安全的参数解析**
   - 使用泛型提供编译时类型检查
   - 避免运行时类型断言错误

4. **完整的文档和示例**
   - [README.md](README.md) - 完整使用指南
   - [QUICKSTART.md](QUICKSTART.md) - 5分钟快速入门
   - [COMPARISON.md](COMPARISON.md) - 新旧API对比
   - [example/example.go](example/example.go) - 8个完整示例
   - [worker_test.go](worker_test.go) - 单元测试和基准测试

## 使用对比

### 旧方式 ❌

```go
// 需要导入多个包
import (
    "devinggo/modules/system/pkg/worker/server"
    "devinggo/modules/system/pkg/worker/task"
)

// 全局注册
func init() {
    server.Register(testWorker)
}

// 创建任务结构体
type ctestTask struct {
    Type    string
    Payload *glob.Payload
}
// ... 约50行代码

// 发送任务
taskItem := test_task.New()
taskItem.Send(ctx, data)
```

### 新方式 ✅

```go
// 只需一个包
import "devinggo/modules/system/pkg/worker"

// 统一管理
mgr := worker.New(ctx)
mgr.RegisterWorker(&TestWorker{})

// 直接发送任务
worker.NewTaskBuilder(ctx, "test_task").
    WithData(data).
    WithQueue("critical").
    WithDelay(5*time.Second).
    Send()
```

## 迁移步骤

### 1. 保持兼容

**所有旧代码仍然可以正常工作，无需立即迁移！**

新旧API可以在同一项目中共存：
- 现有Worker继续使用旧API
- 新功能使用新API

### 2. 逐步迁移（可选）

如果想享受新API的便利，可以逐步迁移：

#### 迁移Worker注册

```go
// 旧代码
func init() {
    server.Register(testWorker)
}

// 新代码（在main或setup函数中）
mgr := worker.New(ctx)
mgr.RegisterWorker(&TestWorker{})
```

#### 迁移任务发送

```go
// 旧代码
taskItem := test_task.New()
taskItem.Send(ctx, map[string]interface{}{
    "name": "test",
})

// 新代码
worker.NewTaskBuilder(ctx, "test_task").
    WithData(map[string]interface{}{
        "name": "test",
    }).
    Send()
```

#### 迁移参数解析

```go
// 旧代码
payload, _ := glob2.GetPayload(ctx, t)
data := payload.Data.(map[string]interface{})
name := data["name"].(string)  // 可能panic

// 新代码
type TestData struct {
    Name string `json:"name"`
}
data, err := worker.GetParameters[TestData](ctx, t)
// data.Name 直接使用，类型安全
```

## 优势总结

| 特性 | 旧版本 | 新版本 |
|-----|-------|-------|
| 包导入数量 | 3-5个 | 1个 |
| 代码量 | 多 | 少50-80% |
| 类型安全 | ❌ | ✅ |
| 链式调用 | ❌ | ✅ |
| 灵活配置 | ❌ | ✅ |
| 学习曲线 | 陡峭 | 平缓 |
| 文档完善度 | 基础 | 完善 |

## 常见问题

### Q: 必须立即迁移吗？
A: 不需要。旧API完全兼容，可以继续使用。

### Q: 如何在现有项目中使用新API？
A: 直接使用即可，新旧API可以共存。

### Q: 性能有影响吗？
A: 无影响。新API只是对旧API的友好封装。

### Q: 需要修改配置文件吗？
A: 不需要。配置完全兼容。

## 快速开始

查看 [QUICKSTART.md](QUICKSTART.md) 了解5分钟快速入门教程！

## 反馈

如有问题或建议，请联系项目维护者。

---

**推荐**：新项目直接使用新API，体验更好的开发者体验！✨
