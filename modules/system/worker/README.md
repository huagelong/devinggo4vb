# Worker 使用指南

本模块提供了基于新版Worker API的完整实现，支持自动加载和注册机制。

## 架构说明

### 自动加载机制

系统使用了分层的自动加载架构：

1. **模块加载层** (`modules/_/worker/system.go`)
   - 通过 `import _ "devinggo/modules/system/worker"` 自动加载system模块的worker
   - 支持多模块扩展，新模块只需添加对应的加载文件

2. **自动注册层** (`modules/system/worker/worker.go`)
   - 在包的 `init()` 函数中创建全局Manager
   - 自动注册所有Worker和Cron任务到全局Manager
   - 无需手动调用注册函数

3. **使用层** (`cmd/worker.go`)
   - 通过 `worker.GetManager()` 获取全局Manager
   - 直接使用已注册的Worker启动服务

### 优势

✅ **自动发现** - 新增Worker只需创建文件，init自动注册  
✅ **模块化** - 支持多模块各自定义Worker  
✅ **简单易用** - 启动代码只需获取Manager即可  
✅ **向后兼容** - 保持原有的自动加载机制  

## 快速开始

### 1. Worker自动启动

Worker服务会通过以下链路自动启动：

```
cmd/worker.go 
  ↓ import
modules/_/worker/system.go 
  ↓ import  
modules/system/worker (package)
  ↓ init()
自动注册所有Worker和Cron任务
```

在应用启动代码中：

```go
package main

import (
    "context"
    "devinggo/modules/system/worker"
    "log"
)

func main() {
    ctx := context.Background()
    
    // 获取全局Manager（已在init中自动注册所有Worker）
    mgr := worker.GetManager()
    
    // 启动Worker服务器（处理任务）
    go func() {
        log.Println("启动Worker服务器...")
        if err := mgr.RunServer(); err != nil {
            log.Fatalf("Worker服务器错误: %v", err)
        }
    }()
    
    // 启动Cron调度器（定时任务）
    go func() {
        log.Println("启动Cron调度器...")
        if err := mgr.RunCron(); err != nil {
            log.Fatalf("Cron调度器错误: %v", err)
        }
    }()
    
    // 应用主逻辑
    // ...
}
```

### 2. 发送任务

#### 方式一：使用便捷函数（test_task）

```go
import (
    "devinggo/modules/system/worker/task/test_task"
)

// 立即发送任务
err := test_task.Send(ctx, test_task.TestTaskData{
    Name: "测试任务",
})

// 延迟5秒发送
err := test_task.SendWithDelay(ctx, test_task.TestTaskData{
    Name: "延迟任务",
}, 5)
```

#### 方式二：使用TaskBuilder（推荐）

```go
import (
    "devinggo/modules/system/pkg/worker"
    "devinggo/modules/system/worker/consts"
    "time"
)

// 立即执行
err := worker.NewTaskBuilder(ctx, consts.TEST_TASK).
    WithData(map[string]interface{}{
        "name": "测试",
    }).
    Send()

// 延迟5分钟执行
err := worker.NewTaskBuilder(ctx, consts.TEST_TASK).
    WithData(map[string]interface{}{
        "name": "测试",
    }).
    WithDelay(5 * time.Minute).
    WithQueue("critical").
    Send()

// 在指定时间执行
tomorrow := time.Now().Add(24 * time.Hour)
err := worker.NewTaskBuilder(ctx, consts.URL_CRON).
    WithData(map[string]interface{}{
        "url": "https://api.example.com/webhook",
        "method": "POST",
    }).
    WithProcessAt(tomorrow).
    Send()
```

## 已实现的Worker和任务

### Worker处理器

1. **TestWorker** (`consts.TEST_TASK`) - 测试任务处理器
2. **TestCronWorker** (`consts.TEST_CRON`) - 测试定时任务处理器
3. **CmdCronWorker** (`consts.CMD_CRON`) - 命令执行处理器
4. **UrlCronWorker** (`consts.URL_CRON`) - HTTP请求处理器

### 定时任务

1. **TestCron** - 测试定时任务
2. **CmdCron** - 命令执行定时任务
3. **UrlCron** - HTTP请求定时任务

## 任务类型说明

### 1. 测试任务（TEST_TASK）

```go
type TestTaskData struct {
    Name string `json:"name"`
}

// 发送方式
test_task.Send(ctx, test_task.TestTaskData{
    Name: "测试名称",
})
```

### 2. URL请求任务（URL_CRON）

```go
type UrlCronData struct {
    Url         string                 `json:"url"`
    Method      string                 `json:"method"`
    Headers     map[string]string      `json:"headers"`
    Params      map[string]interface{} `json:"params"`
    Timeout     int64                  `json:"timeout"`
    Retry       int                    `json:"retry"`
    Cookies     map[string]string      `json:"cookie"`
    ContentType string                 `json:"content_type"`
    Proxy       string                 `json:"proxy"`
}

// 发送方式
worker.NewTaskBuilder(ctx, consts.URL_CRON).
    WithData(map[string]interface{}{
        "url": "https://api.example.com/webhook",
        "method": "POST",
        "headers": map[string]string{
            "Authorization": "Bearer token",
        },
        "params": map[string]interface{}{
            "key": "value",
        },
    }).
    Send()
```

### 3. 命令执行任务（CMD_CRON）

```go
type CmdCronData struct {
    Cmd string `json:"cmd"`
}

// 发送方式
worker.NewTaskBuilder(ctx, consts.CMD_CRON).
    WithData(map[string]interface{}{
        "cmd": "echo 'Hello World'",
    }).
    Send()
```

## 添加新的Worker

### 1. 创建Worker实现

在 `modules/system/worker/server/` 目录下创建新文件：

```go
package server

import (
    "context"
    "devinggo/modules/system/pkg/worker"
    "devinggo/modules/system/worker/consts"
    "github.com/hibiken/asynq"
)

type MyWorker struct {
    Type string
}

// NewMyWorker 创建Worker实例
func NewMyWorker() *MyWorker {
    return &MyWorker{
        Type: consts.MY_TASK,
    }
}

func (w *MyWorker) GetType() string {
    return w.Type
}

func (w *MyWorker) Execute(ctx context.Context, t *asynq.Task) error {
    // 解析参数
    data, err := worker.GetParameters[MyTaskData](ctx, t)
    if err != nil {
        return err
    }
    
    // 执行业务逻辑
    worker.GetLogger().Infof(ctx, "处理任务: %+v", data)
    
    return nil
}

type MyTaskData struct {
    Field1 string `json:"field1"`
    Field2 int    `json:"field2"`
}
```

### 2. 注册Worker

在 `modules/system/worker/worker.go` 的 init 函数中添加注册：

```go
func init() {
    ctx := gctx.GetInitCtx()
    globalManager = worker.New(ctx)

    globalManager.RegisterWorker(server.NewTestWorker()).
        RegisterWorker(server.NewMyWorker()).  // 添加新Worker
        // ...其他Worker
}
```

### 3. 发送任务

```go
worker.NewTaskBuilder(ctx, consts.MY_TASK).
    WithData(MyTaskData{
        Field1: "value1",
        Field2: 123,
    }).
    Send()
```

## 迁移说明

### 旧代码

```go
// 旧的init注册方式（已删除）
var testWorker = &cTestWorker{
    Type: consts.TEST_TASK,
}

func init() {
    server.Register(testWorker)
}

// 旧的任务发送方式（已简化）
taskItem := test_task.New()
taskItem.Send(ctx, data)
```

### 新代码

```go
// 新的构造函数方式
func NewTestWorker() *cTestWorker {
    return &cTestWorker{
        Type: consts.TEST_TASK,
    }
}

// 在 worker.go 的 init 函数中自动注册
// 应用启动时自动执行，无需手动调用

// 新的任务发送方式
test_task.Send(ctx, data)
// 或使用TaskBuilder
worker.NewTaskBuilder(ctx, consts.TEST_TASK).
    WithData(data).
    Send()
```

## 优势

✅ **集中管理** - 所有Worker和Cron任务在一个地方注册，清晰明了  
✅ **按需加载** - 不再依赖init函数的自动注册  
✅ **灵活配置** - 可以根据环境选择性注册Worker  
✅ **易于测试** - 可以为测试创建独立的Worker集合  
✅ **代码简洁** - 使用TaskBuilder大幅简化任务发送代码  

## 配置

在 `config.yaml` 中配置Worker参数：

```yaml
worker:
  redis:
    address: "localhost:6379"
    pass: ""
    db: 3
  queues:
    critical: 6  # 高优先级
    default: 3   # 默认优先级
    low: 1       # 低优先级
  concurrency: 10  # 并发worker数量
  shutdownTimeout: "10s"
  location: "Asia/Shanghai"  # 时区
```

## 故障排查

### Worker未执行任务

1. 检查Worker服务器是否启动：`mgr.RunServer()`
2. 检查Redis连接配置
3. 查看日志是否有错误信息

### Cron任务未执行

1. 检查Cron调度器是否启动：`mgr.RunCron()`
2. 检查数据库中的定时任务配置
3. 确认定时任务状态为启用（status=1）

### 任务发送失败

1. 检查任务类型是否已注册
2. 检查数据结构是否匹配
3. 检查Redis连接是否正常

## 更多信息

- [Worker包文档](../pkg/worker/README.md)
- [快速入门](../pkg/worker/QUICKSTART.md)
- [API对比](../pkg/worker/COMPARISON.md)
- [完整示例](../pkg/worker/example/example.go)
