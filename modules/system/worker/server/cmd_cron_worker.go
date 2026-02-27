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

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/hibiken/asynq"
)

type cCmdCronWorker struct {
	Type string
}

func init() {
	// 自动注册Worker
	worker.Register(NewCmdCronWorker())
}

// NewCmdCronWorker 创建命令执行Worker
func NewCmdCronWorker() *cCmdCronWorker {
	return &cCmdCronWorker{
		Type: consts.CMD_CRON,
	}
}

func (s *cCmdCronWorker) GetType() string {
	return s.Type
}

// Execute 执行任务
func (s *cCmdCronWorker) Execute(ctx context.Context, t *asynq.Task) (err error) {
	data, err := glob2.GetParamters[cron.CmdCronData](ctx, t)
	if err != nil {
		return err
	}
	glob2.WithWorkLog().Infof(ctx, `type:%s, jsonData:%+v`, t.Type(), data)
	r, err := gproc.ShellExec(gctx.New(), data.Cmd)
	if err != nil {
		return err
	}
	glob2.WithWorkLog().Infof(ctx, `type:%s, response:%+v`, t.Type(), r)

	return
}
