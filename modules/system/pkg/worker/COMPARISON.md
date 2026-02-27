# Worker 包优化对比

## 优化概览

Worker包经过重新设计，提供了更加优雅、易用的API。以下是新旧版本的对比。

## 主要改进

✅ **统一的管理入口**：通过`Manager`统一管理所有功能  
✅ **Builder模式**：使用链式调用构建任务，代码更清晰  
✅ **减少包导入**：只需导入一个包即可使用所有功能  
✅ **类型安全**：使用泛型提供类型安全的参数解析  
✅ **更好的可读性**：代码意图更加清晰明了  

## API对比

### 1. 注册Worker

#### 旧版本 ❌
```go
// 需要导入多个包
import (
    "devinggo/modules/system/pkg/worker/server"
    "devinggo/modules/system/pkg/worker/cron"
)

var testWorker = &cTestWorker{
    Type: consts.TEST_TASK,
}

// 使用init函数全局注册
func init() {
    server.Register(testWorker)
}
```

#### 新版本 ✅
```go
// 只需导入一个包
import "devinggo/modules/system/pkg/worker"

// 统一管理，可以在任何地方注册
mgr := worker.New(ctx)
mgr.RegisterWorker(&TestWorker{})

// 或者链式注册多个Worker
mgr.RegisterWorker(&TestWorker{}).
    RegisterWorker(&EmailWorker{}).
    RegisterWorker(&NotificationWorker{})
```

**优点：**
- 不依赖init函数，更灵活
- 统一的注册方式
- 支持链式调用
- 减少包导入

---

### 2. 发送任务

#### 旧版本 ❌
```go
import (
    "devinggo/modules/system/pkg/worker/task"
    "devinggo/modules/system/pkg/worker/glob"
    "devinggo/modules/system/worker/consts"
)

// 需要定义完整的结构体
type ctestTask struct {
    Type    string
    Payload *glob.Payload
}

func New() *ctestTask {
    return &ctestTask{
        Type: consts.TEST_TASK,
        Payload: &glob.Payload{
            Time:      asynq.ProcessIn(0),
            TaskID:    consts.TEST_TASK,
            QueueName: "default",
        },
    }
}

func (s *ctestTask) GetType() string {
    return s.Type
}

func (s *ctestTask) Send(ctx context.Context, data interface{}) error {
    s.Payload.Data = data
    return task.NewTask(ctx, s)
}

func (s *ctestTask) GetPayload() *glob.Payload {
    return s.Payload
}

// 使用时：
taskItem := New()
err := taskItem.Send(ctx, map[string]interface{}{
    "name": "test",
})
```

#### 新版本 ✅
```go
import "devinggo/modules/system/pkg/worker"

// 直接使用Builder模式，无需定义结构体
err := worker.NewTaskBuilder(ctx, "test_task").
    WithData(map[string]interface{}{
        "name": "test",
    }).
    Send()
```

**优点：**
- 代码量减少80%
- 无需定义额外的结构体
- 链式调用，意图清晰
- 自动处理默认值

---

### 3. 延迟任务

#### 旧版本 ❌
```go
// 需要在创建时就设置好所有参数
func New() *ctestTask {
    return &ctestTask{
        Type: consts.TEST_TASK,
        Payload: &glob.Payload{
            Time:      asynq.ProcessIn(5 * time.Second), // 固定的延迟
            TaskID:    consts.TEST_TASK,
            QueueName: "critical", // 固定的队列
        },
    }
}

// 如果需要不同的参数，要创建多个New函数
func NewWithDelay(delay time.Duration) *ctestTask { /* ... */ }
func NewWithQueue(queue string) *ctestTask { /* ... */ }
```

#### 新版本 ✅
```go
import "devinggo/modules/system/pkg/worker"

// 灵活配置各种参数
err := worker.NewTaskBuilder(ctx, "test_task").
    WithData(data).
    WithDelay(5 * time.Second).      // 动态设置延迟
    WithQueue("critical").            // 动态设置队列
    WithTaskID("unique_id").          // 动态设置任务ID
    Send()

// 或者指定时间执行
err := worker.NewTaskBuilder(ctx, "test_task").
    WithData(data).
    WithProcessAt(time.Now().Add(1 * time.Hour)).  // 1小时后执行
    Send()

// 或者设置唯一任务
err := worker.NewTaskBuilder(ctx, "test_task").
    WithData(data).
    WithTaskID("unique_task").
    WithRetention(24 * time.Hour).    // 24小时内唯一
    Send()
```

**优点：**
- 参数灵活可配置
- 不需要为每种配置创建新函数
- 代码可读性强
- 易于理解和维护

---

### 4. 解析任务参数

#### 旧版本 ❌
```go
import (
    glob2 "devinggo/modules/system/pkg/worker/glob"
)

func (s *cTestWorker) Execute(ctx context.Context, t *asynq.Task) (err error) {
    // 手动解析，没有类型检查
    payload, err := glob2.GetPayload(ctx, t)
    if err != nil {
        return err
    }
    
    // 需要手动类型断言，容易出错
    data := payload.Data.(map[string]interface{})
    name := data["name"].(string)  // 可能panic
    
    glob2.WithWorkLog().Infof(ctx, "name: %s", name)
    return
}
```

#### 新版本 ✅
```go
import "devinggo/modules/system/pkg/worker"

// 定义类型
type TestData struct {
    Name  string `json:"name"`
    Value int    `json:"value"`
}

func (w *TestWorker) Execute(ctx context.Context, t *asynq.Task) error {
    // 类型安全的解析，编译时检查
    data, err := worker.GetParameters[TestData](ctx, t)
    if err != nil {
        return err
    }
    
    // data是TestData类型，可以直接使用，不会panic
    worker.GetLogger().Infof(ctx, "name: %s, value: %d", data.Name, data.Value)
    return nil
}
```

**优点：**
- 类型安全，编译时检查
- 不会发生类型断言panic
- 代码更清晰易读
- IDE自动补全支持

---

### 5. 启动Worker服务器

#### 旧版本 ❌
```go
import (
    "devinggo/modules/system/pkg/worker/server"
    "devinggo/modules/system/pkg/worker/cron"
)

// 需要分别导入和运行
go func() {
    if err := server.Run(ctx); err != nil {
        log.Fatal(err)
    }
}()

go func() {
    if err := cron.Run(ctx); err != nil {
        log.Fatal(err)
    }
}()
```

#### 新版本 ✅
```go
import "devinggo/modules/system/pkg/worker"

mgr := worker.New(ctx)

// 统一管理，代码更清晰
go func() {
    if err := mgr.RunServer(); err != nil {
        log.Fatal(err)
    }
}()

go func() {
    if err := mgr.RunCron(); err != nil {
        log.Fatal(err)
    }
}()
```

**优点：**
- 统一的接口
- 清晰的职责
- 易于测试和维护

---

### 6. 完整示例对比

#### 旧版本 ❌

```go
// ========== 文件1: task定义 ==========
package test_task

import (
    "context"
    "devinggo/modules/system/pkg/worker/glob"
    "devinggo/modules/system/pkg/worker/task"
    "devinggo/modules/system/worker/consts"
    "github.com/hibiken/asynq"
)

type ctestTask struct {
    Type    string
    Payload *glob.Payload
}

func New() *ctestTask {
    return &ctestTask{
        Type: consts.TEST_TASK,
        Payload: &glob.Payload{
            Time:   asynq.ProcessIn(0),
            TaskID: consts.TEST_TASK,
        },
    }
}

func (s *ctestTask) GetType() string {
    return s.Type
}

func (s *ctestTask) Send(ctx context.Context, data interface{}) error {
    s.Payload.Data = data
    return task.NewTask(ctx, s)
}

func (s *ctestTask) GetPayload() *glob.Payload {
    return s.Payload
}

// ========== 文件2: worker定义 ==========
package server

import (
    glob2 "devinggo/modules/system/pkg/worker/glob"
    "devinggo/modules/system/pkg/worker/server"
    "devinggo/modules/system/worker/consts"
    "context"
    "github.com/hibiken/asynq"
)

var testWorker = &cTestWorker{
    Type: consts.TEST_TASK,
}

type cTestWorker struct {
    Type string
}

func init() {
    server.Register(testWorker)
}

func (s *cTestWorker) GetType() string {
    return s.Type
}

func (s *cTestWorker) Execute(ctx context.Context, t *asynq.Task) (err error) {
    payload, err := glob2.GetPayload(ctx, t)
    if err != nil {
        return err
    }
    
    data := payload.Data.(map[string]interface{})
    name := data["name"].(string)
    
    glob2.WithWorkLog().Infof(ctx, "处理任务: %s", name)
    return
}

// ========== 文件3: 使用 ==========
package main

import (
    "devinggo/modules/system/worker/task/test_task"
    "devinggo/modules/system/pkg/worker/server"
)

func main() {
    // 启动服务器
    go server.Run(ctx)
    
    // 发送任务
    taskItem := test_task.New()
    taskItem.Send(ctx, map[string]interface{}{
        "name": "测试",
    })
}
```

**问题：**
- 需要3个文件
- 约80行代码
- 需要导入5个包
- 使用了全局变量和init函数
- 类型不安全

---

#### 新版本 ✅

```go
// ========== 一个文件搞定 ==========
package main

import (
    "context"
    "devinggo/modules/system/pkg/worker"
    "github.com/hibiken/asynq"
)

// 1. 定义数据结构
type TestData struct {
    Name string `json:"name"`
}

// 2. 定义Worker
type TestWorker struct{}

func (w *TestWorker) GetType() string {
    return "test_task"
}

func (w *TestWorker) Execute(ctx context.Context, t *asynq.Task) error {
    data, err := worker.GetParameters[TestData](ctx, t)
    if err != nil {
        return err
    }
    
    worker.GetLogger().Infof(ctx, "处理任务: %s", data.Name)
    return nil
}

// 3. 使用
func main() {
    ctx := context.Background()
    
    // 创建管理器并注册
    mgr := worker.New(ctx).RegisterWorker(&TestWorker{})
    
    // 启动服务器
    go mgr.RunServer()
    
    // 发送任务
    worker.NewTaskBuilder(ctx, "test_task").
        WithData(TestData{Name: "测试"}).
        Send()
}
```

**优点：**
- 只需1个文件
- 约35行代码（减少56%）
- 只需导入2个包
- 无全局变量，无init函数
- 完全类型安全
- 代码逻辑清晰

---

## 总结

| 特性 | 旧版本 | 新版本 |
|-----|-------|-------|
| 包导入数量 | 5+ | 1-2 |
| 代码量 | 多 | 少（减少50-80%） |
| 类型安全 | ❌ | ✅ |
| 链式调用 | ❌ | ✅ |
| Builder模式 | ❌ | ✅ |
| 统一管理 | ❌ | ✅ |
| 灵活配置 | ❌ | ✅ |
| 易于测试 | ❌ | ✅ |
| 代码可读性 | 一般 | 优秀 |
| 学习曲线 | 陡峭 | 平缓 |

## 迁移建议

1. **新项目**：直接使用新版本API
2. **现有项目**：
   - 新功能使用新API
   - 旧代码可以逐步迁移
   - 新旧API可以共存

## 兼容性

新版本完全兼容旧版本，所有旧的API仍然可以正常使用。你可以：
- 在同一个项目中同时使用新旧API
- 逐步将旧代码迁移到新API
- 新功能使用新API，旧功能保持不变

建议：**新代码使用新API，享受更好的开发体验！**
