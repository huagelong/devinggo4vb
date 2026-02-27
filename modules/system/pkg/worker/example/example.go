// Package example
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package example

import (
	"context"
	"devinggo/modules/system/pkg/worker"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

// ================ 示例1: 简单的Worker实现 ================

// SimpleWorker 简单Worker示例
type SimpleWorker struct{}

func (w *SimpleWorker) GetType() string {
	return "simple_task"
}

func (w *SimpleWorker) Execute(ctx context.Context, t *asynq.Task) error {
	// 获取任务数据
	payload, err := worker.GetPayload(ctx, t)
	if err != nil {
		return err
	}

	worker.GetLogger().Info(ctx, "处理简单任务:", payload.Data)
	return nil
}

// ================ 示例2: 带类型参数的Worker ================

// UserTaskData 用户任务数据
type UserTaskData struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Action   string `json:"action"`
}

// UserWorker 用户任务Worker
type UserWorker struct{}

func (w *UserWorker) GetType() string {
	return "user_task"
}

func (w *UserWorker) Execute(ctx context.Context, t *asynq.Task) error {
	// 使用泛型解析参数，类型安全
	data, err := worker.GetParameters[UserTaskData](ctx, t)
	if err != nil {
		return err
	}

	worker.GetLogger().Infof(ctx, "处理用户任务 - 用户: %s (ID: %d), 操作: %s",
		data.Username, data.UserID, data.Action)

	// 执行实际的业务逻辑
	switch data.Action {
	case "send_email":
		// 发送邮件
		return sendEmailToUser(ctx, data.UserID)
	case "update_profile":
		// 更新用户信息
		return updateUserProfile(ctx, data.UserID)
	default:
		return fmt.Errorf("未知操作: %s", data.Action)
	}
}

func sendEmailToUser(ctx context.Context, userID int64) error {
	worker.GetLogger().Infof(ctx, "发送邮件给用户: %d", userID)
	return nil
}

func updateUserProfile(ctx context.Context, userID int64) error {
	worker.GetLogger().Infof(ctx, "更新用户信息: %d", userID)
	return nil
}

// ================ 示例3: 使用Manager统一管理 ================

// SetupWorkers 设置所有Worker
func SetupWorkers(ctx context.Context) *worker.Manager {
	mgr := worker.New(ctx)

	// 链式注册多个Worker
	mgr.RegisterWorker(&SimpleWorker{}).
		RegisterWorker(&UserWorker{})

	return mgr
}

// ================ 示例4: 使用TaskBuilder发送任务 ================

// SendUserTaskExample 发送用户任务示例
func SendUserTaskExample(ctx context.Context) error {
	// 方式1: 立即执行
	err := worker.NewTaskBuilder(ctx, "user_task").
		WithData(UserTaskData{
			UserID:   12345,
			Username: "张三",
			Action:   "send_email",
		}).
		Send()
	if err != nil {
		return err
	}

	// 方式2: 延迟5分钟执行
	err = worker.NewTaskBuilder(ctx, "user_task").
		WithData(UserTaskData{
			UserID:   12345,
			Username: "张三",
			Action:   "update_profile",
		}).
		WithDelay(5 * time.Minute).
		WithQueue("default").
		Send()
	if err != nil {
		return err
	}

	// 方式3: 在指定时间执行
	tomorrow := time.Now().Add(24 * time.Hour)
	err = worker.NewTaskBuilder(ctx, "user_task").
		WithData(UserTaskData{
			UserID:   12345,
			Username: "张三",
			Action:   "send_email",
		}).
		WithProcessAt(tomorrow).
		WithQueue("critical").
		Send()
	if err != nil {
		return err
	}

	// 方式4: 唯一任务（使用TaskID和Retention）
	err = worker.NewTaskBuilder(ctx, "user_task").
		WithData(UserTaskData{
			UserID:   12345,
			Username: "张三",
			Action:   "send_email",
		}).
		WithTaskID("user_12345_send_email").
		WithRetention(24 * time.Hour). // 24小时内只执行一次
		Send()
	if err != nil {
		return err
	}

	worker.GetLogger().Info(ctx, "所有任务发送成功")
	return nil
}

// ================ 示例5: 批量发送任务 ================

// SendBatchTasks 批量发送任务
func SendBatchTasks(ctx context.Context, userIDs []int64) error {
	for _, userID := range userIDs {
		err := worker.NewTaskBuilder(ctx, "user_task").
			WithData(UserTaskData{
				UserID:   userID,
				Username: fmt.Sprintf("用户%d", userID),
				Action:   "send_email",
			}).
			WithQueue("low"). // 批量任务使用低优先级队列
			Send()

		if err != nil {
			worker.GetLogger().Warningf(ctx, "发送任务失败 - 用户ID: %d, 错误: %v", userID, err)
			continue
		}
	}

	return nil
}

// ================ 示例6: 完整的应用程序示例 ================

// RunWorkerApplication 运行Worker应用程序
func RunWorkerApplication(ctx context.Context) {
	// 1. 创建Worker管理器并注册所有Worker
	mgr := SetupWorkers(ctx)

	// 2. 启动Worker服务器（在独立的goroutine中）
	go func() {
		worker.GetLogger().Info(ctx, "启动Worker服务器...")
		if err := mgr.RunServer(); err != nil {
			worker.GetLogger().Errorf(ctx, "Worker服务器错误: %v", err)
		}
	}()

	// 3. 发送测试任务
	go func() {
		time.Sleep(2 * time.Second) // 等待服务器启动

		worker.GetLogger().Info(ctx, "发送测试任务...")
		if err := SendUserTaskExample(ctx); err != nil {
			worker.GetLogger().Errorf(ctx, "发送任务失败: %v", err)
		}
	}()

	// 4. 保持程序运行
	select {}
}

// ================ 示例7: 便捷函数封装 ================

// NotificationService 通知服务（封装Worker任务）
type NotificationService struct {
	ctx context.Context
}

func NewNotificationService(ctx context.Context) *NotificationService {
	return &NotificationService{ctx: ctx}
}

// SendEmail 发送邮件（对外暴露的简洁API）
func (s *NotificationService) SendEmail(userID int64, username string) error {
	return worker.NewTaskBuilder(s.ctx, "user_task").
		WithData(UserTaskData{
			UserID:   userID,
			Username: username,
			Action:   "send_email",
		}).
		WithQueue("critical").
		Send()
}

// SendDelayedEmail 发送延迟邮件
func (s *NotificationService) SendDelayedEmail(userID int64, username string, delay time.Duration) error {
	return worker.NewTaskBuilder(s.ctx, "user_task").
		WithData(UserTaskData{
			UserID:   userID,
			Username: username,
			Action:   "send_email",
		}).
		WithDelay(delay).
		WithQueue("default").
		Send()
}

// UpdateProfile 更新用户资料
func (s *NotificationService) UpdateProfile(userID int64, username string) error {
	return worker.NewTaskBuilder(s.ctx, "user_task").
		WithData(UserTaskData{
			UserID:   userID,
			Username: username,
			Action:   "update_profile",
		}).
		WithQueue("default").
		Send()
}

// ================ 示例8: 在HTTP Handler中使用 ================

// HTTPHandlerExample HTTP处理器使用示例
func HTTPHandlerExample(ctx context.Context) {
	// 创建通知服务
	notifySvc := NewNotificationService(ctx)

	// 在HTTP请求处理中使用
	userID := int64(12345)
	username := "张三"

	// 立即发送邮件
	if err := notifySvc.SendEmail(userID, username); err != nil {
		worker.GetLogger().Errorf(ctx, "发送邮件任务失败: %v", err)
	}

	// 1小时后发送提醒邮件
	if err := notifySvc.SendDelayedEmail(userID, username, 1*time.Hour); err != nil {
		worker.GetLogger().Errorf(ctx, "发送延迟邮件任务失败: %v", err)
	}

	worker.GetLogger().Info(ctx, "任务已加入队列")
}
