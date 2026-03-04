// Package clonedmod
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package clonedmod

import (
	"context"
	"devinggo/modules/clonedmod/service"

	"github.com/gogf/gf/v2/frame/g"
)

type sClonedmod struct {
}

func New() *sClonedmod {
	return &sClonedmod{}
}

func init() {
	service.RegisterClonedmod(New())
}

// GetModuleName 获取模块名称
func (s *sClonedmod) GetModuleName(ctx context.Context) string {
	return "clonedmod"
}

// GetModuleVersion 获取模块版本
func (s *sClonedmod) GetModuleVersion(ctx context.Context) string {
	return "1.0.0"
}

// Install 安装模块
func (s *sClonedmod) Install(ctx context.Context) error {
	g.Log().Info(ctx, "clonedmod 模块安装")
	return nil
}

// Uninstall 卸载模块
func (s *sClonedmod) Uninstall(ctx context.Context) error {
	g.Log().Info(ctx, "clonedmod 模块卸载")
	return nil
}

// Enable 启用模块
func (s *sClonedmod) Enable(ctx context.Context) error {
	g.Log().Info(ctx, "clonedmod 模块启用")
	return nil
}

// Disable 禁用模块
func (s *sClonedmod) Disable(ctx context.Context) error {
	g.Log().Info(ctx, "clonedmod 模块禁用")
	return nil
}
