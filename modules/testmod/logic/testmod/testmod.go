// Package testmod
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package testmod

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

type sTestmodLogic struct{}

func New() *sTestmodLogic {
	return &sTestmodLogic{}
}

// Test 测试方法
func (s *sTestmodLogic) Test(ctx context.Context) string {
	g.Log().Info(ctx, "testmod Test method called")
	return "testmod module test success"
}
