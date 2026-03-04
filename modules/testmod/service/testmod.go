// Package service
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package service

import (
	"context"
)

type ITestmod interface {
	GetModuleName(ctx context.Context) string
	GetModuleVersion(ctx context.Context) string
	Install(ctx context.Context) error
	Uninstall(ctx context.Context) error
	Enable(ctx context.Context) error
	Disable(ctx context.Context) error
}

var localTestmod ITestmod

func Testmod() ITestmod {
	if localTestmod == nil {
		panic("implement not found for interface ITestmod, forgot register?")
	}
	return localTestmod
}

func RegisterTestmod(i ITestmod) {
	localTestmod = i
}
