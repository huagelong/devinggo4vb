// Package cmd
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package cmd

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"

	"devinggo/hack/generator/internal/generator"
)

var (
	// CrudGenerate CRUD代码生成命令
	CrudGenerate = &gcmd.Command{
		Name:        "crud:generate",
		Brief:       "生成CRUD代码（API、Model、Controller、Logic）",
		Description: "基于数据库表的Entity模型，生成完整的CRUD代码，包括API定义、请求/响应模型、Controller和Logic层。",
		Arguments: []gcmd.Argument{
			{
				Name:  "module",
				Short: "m",
				Brief: "模块名称（例如：system）",
				IsArg: false,
			},
			{
				Name:  "table",
				Short: "t",
				Brief: "数据库表名（例如：system_user）",
				IsArg: false,
			},
			{
				Name:  "name",
				Short: "n",
				Brief: "资源中文名称（例如：用户）",
				IsArg: false,
			},
		},
		Func: runCrudGenerate,
		Examples: `
生成system模块的system_user表的CRUD代码：
  go run main.go crud:generate -m=system -t=system_user -n=用户

生成system模块的system_dept表的CRUD代码：
  go run main.go crud:generate -m=system -t=system_dept -n=部门
`,
	}
)

// runCrudGenerate 执行CRUD代码生成
func runCrudGenerate(ctx context.Context, parser *gcmd.Parser) (err error) {
	// 获取参数
	moduleName := parser.GetOpt("module", "").String()
	tableName := parser.GetOpt("table", "").String()
	chineseName := parser.GetOpt("name", "").String()

	// 参数验证
	if moduleName == "" {
		return fmt.Errorf("模块名不能为空，请使用 -m 或 --module 指定")
	}

	if tableName == "" {
		return fmt.Errorf("表名不能为空，请使用 -t 或 --table 指定")
	}

	if chineseName == "" {
		return fmt.Errorf("中文名不能为空，请使用 -n 或 --name 指定")
	}

	// 打印生成信息
	g.Log().Infof(ctx, "开始生成CRUD代码...")
	g.Log().Infof(ctx, "  模块名: %s", moduleName)
	g.Log().Infof(ctx, "  表名: %s", tableName)
	g.Log().Infof(ctx, "  中文名: %s", chineseName)
	fmt.Println()

	// 创建生成器
	gen, err := generator.NewCRUDGenerator(moduleName, tableName, chineseName)
	if err != nil {
		return fmt.Errorf("创建CRUD生成器失败：%v", err)
	}

	// 执行生成
	if err := gen.Generate(); err != nil {
		return fmt.Errorf("生成CRUD代码失败：%v", err)
	}

	// 打印成功信息
	fmt.Println()
	g.Log().Infof(ctx, "CRUD代码生成成功！")
	g.Log().Info(ctx, "生成的文件：")
	g.Log().Infof(ctx, "  1. modules/%s/api/%s/%s.go", moduleName, moduleName, gen.VarName)
	g.Log().Infof(ctx, "  2. modules/%s/model/req/%s.go", moduleName, tableName)
	g.Log().Infof(ctx, "  3. modules/%s/model/res/%s.go", moduleName, tableName)
	g.Log().Infof(ctx, "  4. modules/%s/controller/%s/%s.go", moduleName, moduleName, gen.VarName)
	g.Log().Infof(ctx, "  5. modules/%s/logic/%s/%s.go", moduleName, moduleName, tableName)
	fmt.Println()
	g.Log().Info(ctx, "下一步操作：")
	g.Log().Info(ctx, "  1. 运行 gf gen service 生成Service接口")
	g.Log().Info(ctx, "  2. 在Router中注册Controller")
	g.Log().Info(ctx, "  3. 根据需要调整生成的代码")

	return nil
}
