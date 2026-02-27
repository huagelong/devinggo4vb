// Package server
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package server

import (
	"context"
	"devinggo/modules/system/pkg/worker"
	glob2 "devinggo/modules/system/pkg/worker/glob"
	"devinggo/modules/system/worker/consts"
	"devinggo/modules/system/worker/cron"

	"github.com/hibiken/asynq"
)

type cTestWorker struct {
	Type string
}

func init() {
	// 自动注册Worker
	worker.Register(NewTestWorker())
}

// NewTestWorker 创建测试Worker
func NewTestWorker() *cTestWorker {
	return &cTestWorker{
		Type: consts.TEST_TASK,
	}
}

func (s *cTestWorker) GetType() string {
	return s.Type
}

// Execute 执行任务
func (s *cTestWorker) Execute(ctx context.Context, t *asynq.Task) (err error) {
	data, err := glob2.GetParamters[cron.TestCronData](ctx, t)
	if err != nil {
		return err
	}
	glob2.WithWorkLog().Infof(ctx, `jsonData:%+v`, data)
	return
}
