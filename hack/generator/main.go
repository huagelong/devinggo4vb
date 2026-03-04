package main

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gctx"
)

const (
	ToolVersion = "1.0.0"
	ToolName    = "DevingGo Code Generator"
)

var (
	Main = &gcmd.Command{
		Name:        "generator",
		Usage:       "generator [COMMAND] [OPTIONS]",
		Brief:       "DevingGo代码生成工具集",
		Description: "统一的模块管理、Worker任务和CRUD代码生成工具",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			fmt.Printf("\n%s v%s\n\n", ToolName, ToolVersion)
			fmt.Println("📦 统一的模块管理、Worker任务和CRUD代码生成工具")
			fmt.Println("\n使用方法:")
			fmt.Println("  generator [COMMAND] [OPTIONS]")
			fmt.Println("\n可用命令:")
			fmt.Println("  module:create   - 创建新模块")
			fmt.Println("  module:export   - 导出模块包")
			fmt.Println("  module:import   - 导入模块包")
			fmt.Println("  module:list     - 列出已安装模块")
			fmt.Println("  worker:create   - 创建Worker任务")
			fmt.Println("  crud:generate   - 生成CRUD代码")
			fmt.Println("\n提示: 使用 'generator [COMMAND] -h' 查看命令详细帮助")
			fmt.Println()
			return nil
		},
	}
)

func init() {
	// 后续在这里添加子命令
	// Main.AddCommand(cmdModule.Command, cmdWorker.Command, cmdCrud.Command)
}

func main() {
	Main.Run(gctx.New())
}
