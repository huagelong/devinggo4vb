// Package server
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package server

import (
	"context"
	"devinggo/modules/system/worker/cron"
	"devinggo/modules/system/worker/consts"
	"devinggo/modules/system/pkg/worker"
	glob2 "devinggo/modules/system/pkg/worker/glob"

	"github.com/hibiken/asynq"
)

// 复用 Cron 的数据结构
type ProcessOrderData = cron.ProcessOrderData

func init() {
	// 注册任务处理器
	worker.RegisterWorkerFunc(consts.PROCESS_ORDER_TASK, executeProcessOrder)
}

// executeProcessOrder 执行处理订单任务
func executeProcessOrder(ctx context.Context, t *asynq.Task) error {
	// 获取任务参数
	data, err := glob2.GetParamters[ProcessOrderData](ctx, t)
	if err != nil {
		glob2.WithWorkLog().Errorf(ctx, "[%s] 获取参数失败: %v", consts.PROCESS_ORDER_TASK, err)
		return err
	}
	
	// 可以在这里添加参数验证
	// if err := g.Validator().Data(data).Run(ctx); err != nil {
	//     glob2.WithWorkLog().Errorf(ctx, "[%s] 参数验证失败: %v", consts.PROCESS_ORDER_TASK, err)
	//     return err
	// }
	
	glob2.WithWorkLog().Infof(ctx, "[%s] 开始执行任务, 参数: %+v", consts.PROCESS_ORDER_TASK, data)
	
	// ========================================
	// TODO: 在这里实现你的业务逻辑
	// ========================================
	// 示例：
	// 1. 数据库操作
	// err = dao.User.Ctx(ctx).Where("id", data.UserId).Update(g.Map{...})
	// 
	// 2. 调用外部服务
	// result, err := service.External().DoSomething(ctx, data)
	// 
	// 3. 发送通知
	// err = service.Notification().Send(ctx, data.Email, "通知内容")
	
	glob2.WithWorkLog().Infof(ctx, "[%s] 任务执行完成", consts.PROCESS_ORDER_TASK)
	return nil
}
