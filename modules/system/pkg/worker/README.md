# Worker 包使用指南

Worker包是一个基于Asynq的任务队列管理系统，提供了统一、优雅的API来管理定时任务、异步任务和Worker服务器。

## 主要特性

- 🎯 **统一的管理接口**：通过 `Manager` 统一管理所有Worker相关功能
- 🔨 **Builder模式**：使用链式调用构建任务
- 📦 **类型安全**：使用泛型提供类型安全的参数解析
- 🔄 **定时任务**：支持Cron表达式的定时任务
- ⚡ **异步任务**：支持延迟执行、指定时间执行等多种模式
- 🪝 **中间件支持**：内置日志中间件，支持自定义中间件

## 快速开始

### 1. 创建Worker管理器

```go
import "devinggo/modules/system/pkg/worker"

// 在应用启动时创建管理器
mgr := worker.New(ctx)
```

### 2. 注册Worker处理器

```go
// 实现Worker接口
type MyWorker struct{}

func (w *MyWorker) GetType() string {
    return "my_task"
}

func (w *MyWorker) Execute(ctx context.Context, t *asynq.Task) error {
    // 解析任务参数
    data, err := worker.GetParameters[MyTaskData](ctx, t)
    if err != nil {
        return err
    }
    
    // 执行任务逻辑
    worker.GetLogger().Infof(ctx, "处理任务: %+v", data)
    return nil
}

// 注册Worker
mgr.RegisterWorker(&MyWorker{})
```

### 3. 发送任务（使用Builder模式）

```go
// 方式一：使用TaskBuilder（推荐）
err := worker.NewTaskBuilder(ctx, "my_task").
    WithData(map[string]interface{}{
        "name": "test",
        "value": 123,
    }).
    WithQueue("critical").      // 设置队列优先级
    WithDelay(5*time.Second).   // 延迟5秒执行
    Send()

// 方式二：立即执行任务
err := worker.NewTaskBuilder(ctx, "my_task").
    WithData(myData).
    Send()

// 方式三：在指定时间执行 
err := worker.NewTaskBuilder(ctx, "my_task").
    WithData(myData).
    WithProcessAt(time.Now().Add(1*time.Hour)).  // 1小时后执行
    Send()

// 方式四：设置任务保留时间（用于唯一性控制）
err := worker.NewTaskBuilder(ctx, "my_task").
    WithData(myData).
    WithRetention(24*time.Hour).  // 24小时内唯一
    WithTaskID("unique_task_id").  // 自定义任务ID
    Send()
```

### 4. 启动Worker服务器

```go
// 启动Worker服务器处理任务
go func() {
    if err := mgr.RunServer(); err != nil {
        log.Fatal(err)
    }
}()
```

## 高级用法

### 定时任务（Cron）

```go
// 实现CronTask接口
type MyCronTask struct {
    Type        string
    Description string
    Payload     *glob.Payload
}

func (c *MyCronTask) GetType() string {
    return "my_cron_task"
}

func (c *MyCronTask) GetCronTask() *asynq.Task {
    return task.GetTask(c)
}

func (c *MyCronTask) GetPayload() *glob.Payload {
    return c.Payload
}

func (c *MyCronTask) GetDescription() string {
    return c.Description
}

func (c *MyCronTask) SetParams(ctx context.Context, params *gjson.Json) {
    // 从数据库配置中设置参数
    if !g.IsEmpty(params) {
        data := new(MyData)
        params.Scan(data)
        c.Payload.Data = data
    }
}

// 注册定时任务
mgr.RegisterCronTask(&MyCronTask{
    Type:        "my_cron_task",
    Description: "我的定时任务",
    Payload:     &glob.Payload{},
})

// 启动定时任务调度器
go func() {
    if err := mgr.RunCron(); err != nil {
        log.Fatal(err)
    }
}()
```

### 类型安全的参数解析

```go
type UserTaskData struct {
    UserID   int64  `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

func (w *UserWorker) Execute(ctx context.Context, t *asynq.Task) error {
    // 使用泛型进行类型安全的参数解析
    data, err := worker.GetParameters[UserTaskData](ctx, t)
    if err != nil {
        return err
    }
    
    // data的类型是UserTaskData，可以直接使用
    worker.GetLogger().Infof(ctx, "处理用户: %s (ID: %d)", data.Username, data.UserID)
    return nil
}
```

### 链式调用示例

```go
// 完整的链式调用示例
func SendNotification(ctx context.Context, userID int64, message string) error {
    return worker.NewTaskBuilder(ctx, "send_notification").
        WithData(map[string]interface{}{
            "user_id": userID,
            "message": message,
            "timestamp": time.Now().Unix(),
        }).
        WithQueue("critical").              // 高优先级队列
        WithTaskID(fmt.Sprintf("notify_%d_%d", userID, time.Now().Unix())).
        WithDelay(1*time.Minute).           // 延迟1分钟发送
        Send()
}
```

## 完整示例

### 示例：邮件发送任务

```go
package email

import (
    "context"
    "devinggo/modules/system/pkg/worker"
    "github.com/hibiken/asynq"
    "time"
)

// 1. 定义数据结构
type EmailData struct {
    To      string `json:"to"`
    Subject string `json:"subject"`
    Body    string `json:"body"`
}

// 2. 实现Worker
type EmailWorker struct{}

func (w *EmailWorker) GetType() string {
    return "send_email"
}

func (w *EmailWorker) Execute(ctx context.Context, t *asynq.Task) error {
    data, err := worker.GetParameters[EmailData](ctx, t)
    if err != nil {
        return err
    }
    
    worker.GetLogger().Infof(ctx, "发送邮件到: %s", data.To)
    
    // TODO: 实际的邮件发送逻辑
    // sendEmail(data.To, data.Subject, data.Body)
    
    return nil
}

// 3. 注册Worker
func init() {
    mgr := worker.New(context.Background())
    mgr.RegisterWorker(&EmailWorker{})
}

// 4. 发送任务的便捷方法
func SendEmail(ctx context.Context, to, subject, body string) error {
    return worker.NewTaskBuilder(ctx, "send_email").
        WithData(EmailData{
            To:      to,
            Subject: subject,
            Body:    body,
        }).
        WithQueue("default").
        WithDelay(0).  // 立即执行
        Send()
}

// 5. 发送延迟邮件
func SendDelayedEmail(ctx context.Context, to, subject, body string, delay time.Duration) error {
    return worker.NewTaskBuilder(ctx, "send_email").
        WithData(EmailData{
            To:      to,
            Subject: subject,
            Body:    body,
        }).
        WithQueue("default").
        WithDelay(delay).  // 延迟执行
        Send()
}
```

## 队列优先级配置

在配置文件中设置队列优先级：

```yaml
worker:
  queues:
    critical: 6  # 最高优先级
    default: 3   # 默认优先级
    low: 1       # 低优先级
  concurrency: 10  # 并发数
```

## API参考

### Manager方法

- `New(ctx)` - 创建Worker管理器
- `RegisterWorker(worker)` - 注册Worker处理器
- `RegisterCronTask(task)` - 注册定时任务
- `RunServer()` - 启动Worker服务器
- `RunCron()` - 启动定时任务调度器
- `SendTask(task)` - 发送任务
- `SendSimpleTask(task)` - 发送简单任务
- `GetClient()` - 获取Asynq客户端
- `GetServer()` - 获取Asynq服务器

### TaskBuilder方法

- `NewTaskBuilder(ctx, taskType)` - 创建任务构建器
- `WithData(data)` - 设置任务数据
- `WithQueue(name)` - 设置队列名称
- `WithTaskID(id)` - 设置任务ID
- `WithDelay(duration)` - 设置延迟执行
- `WithProcessAt(time)` - 设置指定时间执行
- `WithRetention(duration)` - 设置任务保留时间
- `WithCrontabID(id)` - 设置定时任务ID
- `Send()` - 发送任务
- `Build()` - 构建asynq.Task对象

### 辅助函数

- `GetPayload(ctx, task)` - 获取任务Payload
- `GetParameters[T](ctx, task)` - 类型安全地获取任务参数
- `GetLogger()` - 获取日志记录器

## 迁移指南

### 从旧版本迁移

**旧的方式：**
```go
// 需要导入多个包
import (
    "devinggo/modules/system/pkg/worker/server"
    "devinggo/modules/system/pkg/worker/task"
    "devinggo/modules/system/pkg/worker/cron"
)

// 注册
server.Register(myWorker)
cron.Register(myCronTask)

// 发送任务
taskItem := mytask.New()
task.NewTask(ctx, taskItem)
```

**新的方式：**
```go
// 只需要导入一个包
import "devinggo/modules/system/pkg/worker"

// 统一管理
mgr := worker.New(ctx)
mgr.RegisterWorker(myWorker).RegisterCronTask(myCronTask)

// 使用Builder发送任务
worker.NewTaskBuilder(ctx, "my_task").
    WithData(data).
    Send()
```

## 最佳实践

1. **使用TaskBuilder**：优先使用Builder模式创建任务，代码更清晰
2. **类型安全**：使用泛型的`GetParameters[T]`解析参数
3. **错误处理**：Worker的Execute方法应该返回错误，系统会自动重试
4. **日志记录**：使用`worker.GetLogger()`记录日志，自动添加上下文
5. **队列优先级**：根据任务重要性选择合适的队列
6. **任务幂等性**：确保任务可以安全地重试
7. **超时控制**：在Execute中使用context的超时控制

## 注意事项

- Worker服务器和Cron调度器应该在不同的goroutine中启动
- 确保Redis连接配置正确
- 生产环境建议设置合适的并发数和队列优先级
- 定时任务的配置从数据库读取，需要先在数据库中配置
