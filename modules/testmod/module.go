// Package testmod
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package testmod

import (
	"context"
	"devinggo/modules/testmod/service"

	"github.com/gogf/gf/v2/frame/g"
)

type sTestmod struct {
}

func New() *sTestmod {
	return &sTestmod{}
}

func init() {
	service.RegisterTestmod(New())
}

// GetModuleName 获取模块名称
func (s *sTestmod) GetModuleName(ctx context.Context) string {
	return "testmod"
}

// GetModuleVersion 获取模块版本
func (s *sTestmod) GetModuleVersion(ctx context.Context) string {
	return "1.0.0"
}

// Install 安装模块
func (s *sTestmod) Install(ctx context.Context) error {
	g.Log().Info(ctx, "testmod 模块安装")
	return nil
}

// Uninstall 卸载模块
func (s *sTestmod) Uninstall(ctx context.Context) error {
	g.Log().Info(ctx, "testmod 模块卸载")
	return nil
}

// Enable 启用模块
func (s *sTestmod) Enable(ctx context.Context) error {
	g.Log().Info(ctx, "testmod 模块启用")
	return nil
}

// Disable 禁用模块
func (s *sTestmod) Disable(ctx context.Context) error {
	g.Log().Info(ctx, "testmod 模块禁用")
	return nil
}
