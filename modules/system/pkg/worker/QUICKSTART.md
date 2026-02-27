# Worker 快速入门

5分钟快速掌握Worker包的使用！

## 第一步：创建Worker

创建一个Worker来处理任务：

```go
package main

import (
    "context"
    "devinggo/modules/system/pkg/worker"
    "github.com/hibiken/asynq"
)

// 1. 定义数据结构
type EmailData struct {
    To      string `json:"to"`
    Subject string `json:"subject"`
    Body    string `json:"body"`
}

// 2. 创建Worker
type EmailWorker struct{}

func (w *EmailWorker) GetType() string {
    return "send_email"
}

func (w *EmailWorker) Execute(ctx context.Context, t *asynq.Task) error {
    // 解析任务参数
    data, err := worker.GetParameters[EmailData](ctx, t)
    if err != nil {
        return err
    }
    
    // 执行业务逻辑
    worker.GetLogger().Infof(ctx, "发送邮件到 %s: %s", data.To, data.Subject)
    
    // TODO: 实际的邮件发送逻辑
    // smtp.SendEmail(data.To, data.Subject, data.Body)
    
    return nil
}
```

## 第二步：注册Worker并启动服务器

```go
func main() {
    ctx := context.Background()
    
    // 创建Worker管理器
    mgr := worker.New(ctx)
    
    // 注册Worker
    mgr.RegisterWorker(&EmailWorker{})
    
    // 启动Worker服务器
    if err := mgr.RunServer(); err != nil {
        panic(err)
    }
}
```

## 第三步：发送任务

在你的业务代码中发送任务：

```go
func SendWelcomeEmail(ctx context.Context, email string) error {
    return worker.NewTaskBuilder(ctx, "send_email").
        WithData(EmailData{
            To:      email,
            Subject: "欢迎注册",
            Body:    "感谢您的注册！",
        }).
        Send()
}
```

完成！就是这么简单！✨

## 更多功能

### 延迟执行

```go
// 5分钟后发送
worker.NewTaskBuilder(ctx, "send_email").
    WithData(data).
    WithDelay(5 * time.Minute).
    Send()
```

### 指定时间执行

```go
// 明天上午9点发送
tomorrow9am := time.Now().Add(24*time.Hour).Truncate(24*time.Hour).Add(9*time.Hour)
worker.NewTaskBuilder(ctx, "send_email").
    WithData(data).
    WithProcessAt(tomorrow9am).
    Send()
```

### 设置优先级

```go
// 高优先级
worker.NewTaskBuilder(ctx, "send_email").
    WithData(data).
    WithQueue("critical").  // 配置文件中: critical: 6
    Send()

// 低优先级
worker.NewTaskBuilder(ctx, "send_email").
    WithData(data).
    WithQueue("low").      // 配置文件中: low: 1
    Send()
```

### 唯一任务（防重复）

```go
// 使用TaskID确保任务唯一性
worker.NewTaskBuilder(ctx, "send_email").
    WithData(data).
    WithTaskID("daily_report_user_123").      // 唯一标识
    WithRetention(24 * time.Hour).            // 24小时内只执行一次
    Send()
```

### 链式注册多个Worker

```go
mgr.RegisterWorker(&EmailWorker{}).
    RegisterWorker(&SMSWorker{}).
    RegisterWorker(&PushWorker{})
```

## 配置文件

在 `config.yaml` 中配置：

```yaml
worker:
  # Redis配置
  redis:
    address: "localhost:6379"
    pass: ""
    db: 3
    dialTimeout: "30s"
    readTimeout: "30s"
    writeTimeout: "30s"
  
  # 队列优先级
  queues:
    critical: 6  # 最高优先级
    default: 3   # 默认优先级
    low: 1       # 低优先级
  
  # 并发数
  concurrency: 10
  
  # 关闭超时
  shutdownTimeout: "10s"
  
  # 时区
  location: "Asia/Shanghai"
```

## 常见模式

### 模式1: 服务封装

```go
type NotificationService struct {
    ctx context.Context
}

func (s *NotificationService) SendEmail(to, subject, body string) error {
    return worker.NewTaskBuilder(s.ctx, "send_email").
        WithData(EmailData{To: to, Subject: subject, Body: body}).
        Send()
}

func (s *NotificationService) SendDelayedEmail(to, subject, body string, delay time.Duration) error {
    return worker.NewTaskBuilder(s.ctx, "send_email").
        WithData(EmailData{To: to, Subject: subject, Body: body}).
        WithDelay(delay).
        Send()
}
```

### 模式2: HTTP Handler中使用

```go
func RegisterHandler(c *gin.Context) {
    // ... 用户注册逻辑 ...
    
    // 异步发送欢迎邮件
    worker.NewTaskBuilder(c.Request.Context(), "send_email").
        WithData(EmailData{
            To:      user.Email,
            Subject: "欢迎注册",
            Body:    "感谢您注册我们的服务...",
        }).
        WithQueue("default").
        Send()
    
    c.JSON(200, gin.H{"message": "注册成功"})
}
```

### 模式3: 批量任务

```go
func SendBulkEmails(ctx context.Context, emails []string) {
    for _, email := range emails {
        worker.NewTaskBuilder(ctx, "send_email").
            WithData(EmailData{
                To:      email,
                Subject: "通知",
                Body:    "这是群发邮件",
            }).
            WithQueue("low").  // 批量任务用低优先级
            Send()
    }
}
```

## 下一步

- 查看 [README.md](README.md) 了解完整功能
- 查看 [example/example.go](example/example.go) 了解更多示例
- 查看 [COMPARISON.md](COMPARISON.md) 了解新旧API对比

祝你使用愉快！🎉
