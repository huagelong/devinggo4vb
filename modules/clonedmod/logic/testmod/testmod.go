// Package clonedmod
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package clonedmod

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

type sClonedmodLogic struct{}

func New() *sClonedmodLogic {
	return &sClonedmodLogic{}
}

// Test 测试方法
func (s *sClonedmodLogic) Test(ctx context.Context) string {
	g.Log().Info(ctx, "clonedmod Test method called")
	return "clonedmod module test success"
}
