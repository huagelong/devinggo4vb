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

	"devinggo/hack/generator/internal/generator"
	"devinggo/hack/generator/internal/utils"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	// WorkerCreate 创建Worker任务命令
	WorkerCreate = &gcmd.Command{
		Name:  "worker:create",
		Brief: "创建 Task 或 Cron 任务",
		Description: `
创建 Worker 任务（异步任务或定时任务）

快速创建 Task 或 Cron 任务，自动生成规范代码。支持三种模式：
- task: 仅创建异步任务
- cron: 仅创建定时任务
- both: 同时创建任务和定时任务（数据结构共享）

使用示例（命令行模式）:
  # 创建发送邮件任务（默认both类型）
  go run hack/generator/main.go worker:create -name send_email -desc "发送邮件"

  # 仅创建异步任务
  go run hack/generator/main.go worker:create -name process_order -type task -desc "处理订单"

  # 仅创建定时任务
  go run hack/generator/main.go worker:create -name clean_logs -type cron -desc "清理日志"

  # 在指定模块中创建
  go run hack/generator/main.go worker:create -name notify_user -module custom -desc "用户通知"

使用示例（交互式模式）:
  go run hack/generator/main.go worker:create

命令选项:
  -name   任务名称（必填），建议使用下划线命名，如: send_email
  -module 模块名称（可选，默认: system）
  -type   任务类型（可选，默认: both）
          task: 仅创建异步任务
          cron: 仅创建定时任务
          both: 同时创建任务和定时任务
  -desc   任务描述（可选），用于生成注释

生成的文件:
  • modules/{module}/worker/cron/{name}_cron.go        (定时任务文件)
  • modules/{module}/worker/server/{name}_worker.go    (异步任务文件)
  • modules/{module}/worker/consts/const.go            (常量定义文件)

注意事项:
  1. 任务名称建议使用小写下划线格式
  2. 目标模块必须已存在
  3. 不能创建重复的任务名称
  4. 生成后需要重启 worker 服务以加载新任务
		`,
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			// 获取参数
			name := parser.GetOpt("name").String()
			moduleName := parser.GetOpt("module", "system").String()
			workerType := parser.GetOpt("type", "both").String()
			description := parser.GetOpt("desc").String()

			// 交互式模式：未提供任务名称时
			if name == "" {
				fmt.Println("\n🔧 DevingGo Worker 任务创建向导")
				fmt.Println("=" + strings.Repeat("=", 40))

				name = utils.PromptRequiredString("请输入任务名称（建议使用下划线命名，如: send_email）")

				moduleName = utils.PromptString("请输入所属模块名称", "system")

				// 选择任务类型
				typeIndex := utils.PromptSelect("请选择任务类型", []string{
					"both - 同时创建异步任务和定时任务（推荐）",
					"task - 仅创建异步任务",
					"cron - 仅创建定时任务",
				}, 0)
				workerTypes := []string{"both", "task", "cron"}
				workerType = workerTypes[typeIndex]

				description = utils.PromptString("请输入任务描述", name)
			}

			// 转换类型
			var wType generator.WorkerType
			switch workerType {
			case "task":
				wType = generator.WorkerTypeTask
			case "cron":
				wType = generator.WorkerTypeCron
			case "both":
				wType = generator.WorkerTypeBoth
			default:
				return gerror.Newf("类型必须是 task、cron 或 both，当前值: %s", workerType)
			}

			// 创建生成器
			gen := generator.NewWorkerGenerator(ctx, moduleName, name, description, wType)

			// 执行生成
			return gen.Generate()
		},
	}
)
