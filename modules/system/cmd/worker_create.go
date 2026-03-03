// Package cmd
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
)

var (
	CreateWorker = &gcmd.Command{
		Name:  "worker:create",
		Brief: "创建 Task 或 Cron",
		Description: `
		创建一个新的 Task 或 Cron，包含必要的文件和代码
		用法: go run main.go worker:create -name 任务名称 -module 模块名称 [-type task|cron|both]
		
		示例:
		  创建 Task: go run main.go worker:create -name send_email -module system -type task
		  创建 Cron: go run main.go worker:create -name clean_logs -module system -type cron
		  同时创建: go run main.go worker:create -name process_data -module system -type both
		`,
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			CmdInit(ctx, parser)

			// 获取参数
			nameOpt := gcmd.GetOpt("name")
			if nameOpt == nil {
				return gerror.New("任务名称必须输入，使用 -name 参数指定")
			}
			name := nameOpt.String()
			if name == "" {
				return gerror.New("任务名称不能为空")
			}

			moduleOpt := gcmd.GetOpt("module")
			if moduleOpt == nil {
				return gerror.New("模块名称必须输入，使用 -module 参数指定")
			}
			moduleName := moduleOpt.String()
			if moduleName == "" {
				return gerror.New("模块名称不能为空")
			}

			// 获取类型，默认为 both
			typeOpt := gcmd.GetOpt("type")
			workerType := "both"
			if typeOpt != nil && typeOpt.String() != "" {
				workerType = typeOpt.String()
			}

			// 验证类型
			if workerType != "task" && workerType != "cron" && workerType != "both" {
				return gerror.New("类型必须是 task、cron 或 both")
			}

			// 验证模块是否存在
			modulePath := fmt.Sprintf("./modules/%s", moduleName)
			if !gfile.Exists(modulePath) {
				return gerror.Newf("模块 '%s' 不存在", moduleName)
			}

			g.Log().Infof(ctx, "开始创建 Worker: %s (模块: %s, 类型: %s)", name, moduleName, workerType)

			// 创建必要的目录
			if err := createWorkerDirectories(ctx, moduleName, workerType); err != nil {
				return err
			}

			// 创建或更新常量文件
			if err := updateConstFile(ctx, moduleName, name, workerType); err != nil {
				return err
			}

			// 根据类型创建文件
			if workerType == "cron" || workerType == "both" {
				if err := createCronFile(ctx, moduleName, name); err != nil {
					return err
				}
			}

			if workerType == "task" || workerType == "both" {
				if err := createTaskFile(ctx, moduleName, name); err != nil {
					return err
				}
			}

			successMsg := fmt.Sprintf("Worker '%s' 创建成功!", name)
			g.Log().Info(ctx, successMsg)
			fmt.Printf("\n%s\n", successMsg)

			// 输出创建的文件列表
			fmt.Println("创建的文件:")
			if workerType == "cron" || workerType == "both" {
				fmt.Printf("  - modules/%s/worker/cron/%s_cron.go\n", moduleName, name)
			}
			if workerType == "task" || workerType == "both" {
				fmt.Printf("  - modules/%s/worker/server/%s_worker.go\n", moduleName, name)
			}
			fmt.Printf("  - modules/%s/worker/consts/const.go (已更新)\n", moduleName)
			fmt.Println()

			return nil
		},
	}
)

// createWorkerDirectories 创建必要的目录
func createWorkerDirectories(ctx context.Context, moduleName, workerType string) error {
	dirs := []string{
		fmt.Sprintf("./modules/%s/worker", moduleName),
		fmt.Sprintf("./modules/%s/worker/consts", moduleName),
	}

	if workerType == "cron" || workerType == "both" {
		dirs = append(dirs, fmt.Sprintf("./modules/%s/worker/cron", moduleName))
	}

	if workerType == "task" || workerType == "both" {
		dirs = append(dirs, fmt.Sprintf("./modules/%s/worker/server", moduleName))
	}

	for _, dir := range dirs {
		if !gfile.Exists(dir) {
			if err := gfile.Mkdir(dir); err != nil {
				return gerror.Wrapf(err, "创建目录 '%s' 失败", dir)
			}
			g.Log().Debugf(ctx, "创建目录: %s", dir)
		}
	}
	return nil
}

// updateConstFile 创建或更新常量文件
func updateConstFile(ctx context.Context, moduleName, name, workerType string) error {
	constPath := fmt.Sprintf("./modules/%s/worker/consts/const.go", moduleName)

	// 常量名称（转为大写加下划线格式）
	constName := strings.ToUpper(gstr.CaseSnake(name))

	// 检查文件是否存在
	if !gfile.Exists(constPath) {
		// 创建新的常量文件
		content := fmt.Sprintf(`// Package consts
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package consts

var (
`)

		if workerType == "cron" || workerType == "both" {
			content += fmt.Sprintf("\t%s_CRON = \"%s_cron\" // %s定时任务\n", constName, name, name)
		}
		if workerType == "task" || workerType == "both" {
			content += fmt.Sprintf("\t%s_TASK = \"%s_task\" // %s任务\n", constName, name, name)
		}

		content += ")\n"

		if err := gfile.PutContents(constPath, content); err != nil {
			return gerror.Wrapf(err, "创建常量文件失败")
		}
		g.Log().Debugf(ctx, "创建常量文件: %s", constPath)
	} else {
		// 读取现有文件内容
		content := gfile.GetContents(constPath)

		// 检查常量是否已存在
		if workerType == "cron" || workerType == "both" {
			cronConstName := fmt.Sprintf("%s_CRON", constName)
			if strings.Contains(content, cronConstName) {
				g.Log().Warningf(ctx, "常量 %s 已存在，跳过", cronConstName)
			} else {
				// 在最后一个 ) 之前插入新常量
				newConst := fmt.Sprintf("\t%s = \"%s_cron\" // %s定时任务\n", cronConstName, name, name)
				// 找到最后一个 )
				lastParen := strings.LastIndex(content, ")")
				if lastParen != -1 {
					content = content[:lastParen] + newConst + content[lastParen:]
				}
			}
		}

		if workerType == "task" || workerType == "both" {
			taskConstName := fmt.Sprintf("%s_TASK", constName)
			if strings.Contains(content, taskConstName) {
				g.Log().Warningf(ctx, "常量 %s 已存在，跳过", taskConstName)
			} else {
				// 在最后一个 ) 之前插入新常量
				newConst := fmt.Sprintf("\t%s = \"%s_task\" // %s任务\n", taskConstName, name, name)
				// 找到最后一个 )
				lastParen := strings.LastIndex(content, ")")
				if lastParen != -1 {
					content = content[:lastParen] + newConst + content[lastParen:]
				}
			}
		}

		// 写回文件
		if err := gfile.PutContents(constPath, content); err != nil {
			return gerror.Wrapf(err, "更新常量文件失败")
		}
		g.Log().Debugf(ctx, "更新常量文件: %s", constPath)
	}

	return nil
}

// createCronFile 创建 Cron 文件
func createCronFile(ctx context.Context, moduleName, name string) error {
	cronPath := fmt.Sprintf("./modules/%s/worker/cron/%s_cron.go", moduleName, name)

	// 检查文件是否已存在
	if gfile.Exists(cronPath) {
		return gerror.Newf("Cron 文件 '%s' 已存在", cronPath)
	}

	// 常量名称
	constName := strings.ToUpper(gstr.CaseSnake(name))
	// 数据结构名称（驼峰命名）
	structName := gstr.CaseCamel(name) + "CronData"

	content := fmt.Sprintf(`// Package cron
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package cron

import (
	"context"
	"devinggo/modules/%s/worker/consts"
	"devinggo/modules/system/pkg/worker"
	glob2 "devinggo/modules/system/pkg/worker/glob"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

// %s %s定时任务数据结构
type %s struct {
	// TODO: 在这里添加你的参数字段
	// 例如: Name string `+"`json:\"name\"`"+`
}

func init() {
	// 注册定时任务
	worker.RegisterCronFunc(consts.%s_CRON, "%s定时任务", handle%sParams)
}

// handle%sParams 处理定时任务参数
func handle%sParams(ctx context.Context, payload *glob2.Payload, params *gjson.Json) {
	if g.IsEmpty(params) {
		return
	}
	data := new(%s)
	if err := params.Scan(data); err != nil {
		glob2.WithWorkLog().Errorf(ctx, "[%%s] cron SetParams failed:%%v", consts.%s_CRON, err)
		return
	}
	payload.Data = data
}
`, moduleName, structName, name, structName, constName, name, structName, structName, structName, structName, constName)

	if err := gfile.PutContents(cronPath, content); err != nil {
		return gerror.Wrapf(err, "创建 Cron 文件失败")
	}

	g.Log().Debugf(ctx, "创建 Cron 文件: %s", cronPath)
	return nil
}

// createTaskFile 创建 Task 文件
func createTaskFile(ctx context.Context, moduleName, name string) error {
	taskPath := fmt.Sprintf("./modules/%s/worker/server/%s_worker.go", moduleName, name)

	// 检查文件是否已存在
	if gfile.Exists(taskPath) {
		return gerror.Newf("Task 文件 '%s' 已存在", taskPath)
	}

	// 常量名称
	constName := strings.ToUpper(gstr.CaseSnake(name))
	// 函数名称（驼峰命名）
	funcName := gstr.CaseCamel(name)
	// 数据结构名称（如果有对应的 cron，则引用 cron 包）
	structName := gstr.CaseCamel(name) + "CronData"

	content := fmt.Sprintf(`// Package server
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package server

import (
	"context"
	"devinggo/modules/%s/worker/consts"
	"devinggo/modules/system/pkg/worker"
	glob2 "devinggo/modules/system/pkg/worker/glob"

	"github.com/hibiken/asynq"
)

// %sData %s任务数据结构
// 如果有对应的 Cron 任务，可以直接使用 cron 包中的数据结构
// 例如: type %sData = cron.%s
type %sData struct {
	// TODO: 在这里添加你的参数字段
	// 例如: Name string `+"`json:\"name\"`"+`
}

func init() {
	// 注册任务处理器
	worker.RegisterWorkerFunc(consts.%s_TASK, execute%sWorker)
}

// execute%sWorker 执行%s任务
func execute%sWorker(ctx context.Context, t *asynq.Task) error {
	// 获取任务参数
	data, err := glob2.GetParamters[%sData](ctx, t)
	if err != nil {
		glob2.WithWorkLog().Errorf(ctx, "[%%s] 获取参数失败:%%v", consts.%s_TASK, err)
		return err
	}
	
	// TODO: 在这里实现你的业务逻辑
	glob2.WithWorkLog().Infof(ctx, "[%%s] 开始执行任务, 参数: %%+v", consts.%s_TASK, data)
	
	// 示例: 处理业务逻辑
	// ...
	
	glob2.WithWorkLog().Infof(ctx, "[%%s] 任务执行完成", consts.%s_TASK)
	return nil
}
`, moduleName, funcName, name, funcName, structName, funcName, constName, funcName, funcName, name, funcName, funcName, constName, constName, constName)

	if err := gfile.PutContents(taskPath, content); err != nil {
		return gerror.Wrapf(err, "创建 Task 文件失败")
	}

	g.Log().Debugf(ctx, "创建 Task 文件: %s", taskPath)
	return nil
}
