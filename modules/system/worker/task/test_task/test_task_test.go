// Package test_task
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package test_task

import (
	"context"
	"testing"
)

func TestTaskQueue(t *testing.T) {
	// 使用新的Send方法
	err := Send(context.Background(), TestTaskData{Name: "helloworld"})
	if err != nil {
		t.Errorf("发送任务失败: %v", err)
	}
}

func TestTaskQueueWithDelay(t *testing.T) {
	// 测试延迟任务
	err := SendWithDelay(context.Background(), TestTaskData{Name: "delayed task"}, 5)
	if err != nil {
		t.Errorf("发送延迟任务失败: %v", err)
	}
}
