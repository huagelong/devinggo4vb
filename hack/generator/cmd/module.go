package cmd

import (
	"context"
	"fmt"

	"devinggo/hack/generator/internal/generator"
	"devinggo/hack/generator/internal/scanner"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	// ModuleCreate 创建模块命令
	ModuleCreate = &gcmd.Command{
		Name:  "module:create",
		Brief: "创建新模块",
		Description: `
创建一个新的模块，包含基本的目录结构和文件
用法: go run hack/generator/main.go module:create -name 模块名称
`,
		Arguments: []gcmd.Argument{
			{
				Name:   "name",
				Short:  "n",
				Brief:  "模块名称",
				IsArg:  false,
				Orphan: false,
			},
		},
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			moduleName := parser.GetOpt("name").String()
			if moduleName == "" {
				return fmt.Errorf("模块名称必须输入，使用 -name 参数指定")
			}

			g.Log().Infof(ctx, "开始创建模块: %s", moduleName)

			// 创建模块创建器
			creator := generator.NewModuleCreator(ctx)
			if err := creator.Create(moduleName); err != nil {
				return fmt.Errorf("创建模块失败: %w", err)
			}

			g.Log().Infof(ctx, "模块 '%s' 创建成功！", moduleName)
			fmt.Printf("\n✅ 模块 '%s' 创建成功！\n", moduleName)
			fmt.Printf("💡 提示: 请运行 'go run main.go migrate:up' 命令开启应用\n\n")
			return nil
		},
	}

	// ModuleClone 克隆模块命令
	ModuleClone = &gcmd.Command{
		Name:  "module:clone",
		Brief: "从现有模块克隆新模块",
		Description: `
从现有模块快速克隆一个新模块
用法: go run hack/generator/main.go module:clone -source 源模块 -target 目标模块
`,
		Arguments: []gcmd.Argument{
			{
				Name:   "source",
				Short:  "s",
				Brief:  "源模块名称",
				IsArg:  false,
				Orphan: false,
			},
			{
				Name:   "target",
				Short:  "t",
				Brief:  "目标模块名称",
				IsArg:  false,
				Orphan: false,
			},
		},
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			sourceModule := parser.GetOpt("source").String()
			targetModule := parser.GetOpt("target").String()

			if sourceModule == "" {
				return fmt.Errorf("源模块名称必须输入，使用 -source 参数指定")
			}
			if targetModule == "" {
				return fmt.Errorf("目标模块名称必须输入，使用 -target 参数指定")
			}

			// 禁止克隆系统模块
			if sourceModule == "system" {
				return fmt.Errorf("系统模块不能被克隆")
			}
			if targetModule == "system" {
				return fmt.Errorf("不能创建名为 'system' 的模块")
			}

			g.Log().Infof(ctx, "开始克隆模块: %s -> %s", sourceModule, targetModule)

			// 创建模块克隆器
			cloner := generator.NewModuleCloner(ctx)
			if err := cloner.Clone(sourceModule, targetModule); err != nil {
				return fmt.Errorf("克隆模块失败: %w", err)
			}

			g.Log().Infof(ctx, "模块 '%s' 克隆成功！", targetModule)
			fmt.Printf("\n✅ 模块 '%s' 克隆成功！\n", targetModule)
			fmt.Printf("📋 源模块: %s\n", sourceModule)
			fmt.Printf("💡 提示: 请运行 'go run main.go migrate:up' 命令开启应用\n\n")
			return nil
		},
	}

	// ModuleExport 导出模块命令
	ModuleExport = &gcmd.Command{
		Name:  "module:export",
		Brief: "导出模块为zip包",
		Description: `
将模块文件导出为zip压缩包
用法: go run hack/generator/main.go module:export -name 模块名称
`,
		Arguments: []gcmd.Argument{
			{
				Name:   "name",
				Short:  "n",
				Brief:  "模块名称",
				IsArg:  false,
				Orphan: false,
			},
		},
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			moduleName := parser.GetOpt("name").String()
			if moduleName == "" {
				return fmt.Errorf("模块名称必须输入，使用 -name 参数指定")
			}

			// 禁止导出系统模块
			if moduleName == "system" {
				return fmt.Errorf("系统模块不能导出")
			}

			g.Log().Infof(ctx, "开始导出模块: %s", moduleName)

			// 创建导出器
			exporter := generator.NewModuleExporter(ctx)
			zipFile, err := exporter.Export(moduleName)
			if err != nil {
				return fmt.Errorf("导出模块失败: %w", err)
			}

			g.Log().Infof(ctx, "模块 '%s' 导出成功: %s", moduleName, zipFile)
			fmt.Printf("\n✅ 模块 '%s' 导出成功！\n", moduleName)
			fmt.Printf("📦 文件: %s\n\n", zipFile)
			return nil
		},
	}

	// ModuleImport 导入模块命令
	ModuleImport = &gcmd.Command{
		Name:  "module:import",
		Brief: "从zip包导入模块",
		Description: `
从zip压缩包导入模块
用法: go run hack/generator/main.go module:import -file 模块zip文件路径
`,
		Arguments: []gcmd.Argument{
			{
				Name:   "file",
				Short:  "f",
				Brief:  "模块zip文件路径",
				IsArg:  false,
				Orphan: false,
			},
		},
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			zipPath := parser.GetOpt("file").String()
			if zipPath == "" {
				return fmt.Errorf("模块文件路径必须输入，使用 -file 参数指定")
			}

			g.Log().Infof(ctx, "开始导入模块: %s", zipPath)

			// 创建导入器
			importer := generator.NewModuleImporter(ctx)
			moduleName, err := importer.Import(zipPath)
			if err != nil {
				return fmt.Errorf("导入模块失败: %w", err)
			}

			g.Log().Infof(ctx, "模块 '%s' 导入成功！", moduleName)
			fmt.Printf("\n✅ 模块 '%s' 导入成功！\n", moduleName)
			fmt.Printf("💡 提示: 请运行 'go run main.go migrate:up' 命令开启应用\n\n")
			return nil
		},
	}

	// ModuleList 列出模块命令
	ModuleList = &gcmd.Command{
		Name:  "module:list",
		Brief: "列出所有已安装模块",
		Description: `
列出所有已安装的模块及其信息
用法: go run hack/generator/main.go module:list
`,
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			modules, err := scanner.ListModules()
			if err != nil {
				return fmt.Errorf("列出模块失败: %w", err)
			}

			if len(modules) == 0 {
				fmt.Println("\n暂无已安装模块\n")
				return nil
			}

			fmt.Printf("\n📦 已安装模块 (%d个)\n\n", len(modules))
			fmt.Println("┌────────────┬─────────┬──────────────┬─────────┐")
			fmt.Println("│ 模块名称   │ 版本    │ 作者         │ 许可证  │")
			fmt.Println("├────────────┼─────────┼──────────────┼─────────┤")
			for _, module := range modules {
				fmt.Printf("│ %-10s │ %-7s │ %-12s │ %-7s │\n",
					truncateString(module.Name, 10),
					truncateString(module.Version, 7),
					truncateString(module.Author, 12),
					truncateString(module.License, 7))
			}
			fmt.Println("└────────────┴─────────┴──────────────┴─────────┘")
			fmt.Println()
			return nil
		},
	}

	// ModuleValidate 验证模块命令
	ModuleValidate = &gcmd.Command{
		Name:  "module:validate",
		Brief: "验证模块完整性",
		Description: `
验证模块的完整性和配置正确性
用法: go run hack/generator/main.go module:validate -name 模块名称
`,
		Arguments: []gcmd.Argument{
			{
				Name:   "name",
				Short:  "n",
				Brief:  "模块名称",
				IsArg:  false,
				Orphan: false,
			},
		},
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			moduleName := parser.GetOpt("name").String()
			if moduleName == "" {
				return fmt.Errorf("模块名称必须输入，使用 -name 参数指定")
			}

			g.Log().Infof(ctx, "验证模块: %s", moduleName)

			warnings, errors, err := scanner.ValidateModule(moduleName)
			if err != nil {
				return fmt.Errorf("验证模块失败: %w", err)
			}

			fmt.Printf("\n🔍 验证模块 '%s'\n\n", moduleName)

			if len(errors) > 0 {
				fmt.Println("❌ 错误:")
				for _, e := range errors {
					fmt.Printf("  - %s\n", e)
				}
				fmt.Println()
			}

			if len(warnings) > 0 {
				fmt.Println("⚠️  警告:")
				for _, w := range warnings {
					fmt.Printf("  - %s\n", w)
				}
				fmt.Println()
			}

			if len(errors) == 0 && len(warnings) == 0 {
				fmt.Println("✅ 模块验证通过，没有发现问题\n")
			} else {
				fmt.Printf("📊 验证结果: %d个错误，%d个警告\n\n", len(errors), len(warnings))
			}

			return nil
		},
	}
)

// truncateString 截断字符串
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
