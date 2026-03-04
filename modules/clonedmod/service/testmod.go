// Package service
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package service

import (
	"context"
)

type IClonedmod interface {
	GetModuleName(ctx context.Context) string
	GetModuleVersion(ctx context.Context) string
	Install(ctx context.Context) error
	Uninstall(ctx context.Context) error
	Enable(ctx context.Context) error
	Disable(ctx context.Context) error
}

var localClonedmod IClonedmod

func Clonedmod() IClonedmod {
	if localClonedmod == nil {
		panic("implement not found for interface IClonedmod, forgot register?")
	}
	return localClonedmod
}

func RegisterClonedmod(i IClonedmod) {
	localClonedmod = i
}
