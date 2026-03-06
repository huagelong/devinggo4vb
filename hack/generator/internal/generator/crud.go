// Package generator
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/text/gstr"
)

// CRUDGenerator CRUD代码生成器
type CRUDGenerator struct {
	ModuleName  string  // 模块名（例如：system）
	TableName   string  // 表名（例如：system_api）
	EntityName  string  // 实体名（例如：SystemApi）
	VarName     string  // 变量名（例如：api）
	PackageName string  // 包名（例如：system）
	ChineseName string  // 中文名（例如：接口）
	Fields      []Field // 字段列表
	WorkDir     string  // 工作目录
}

// Field 字段信息
type Field struct {
	Name         string // 字段名（Go格式，例如：GroupId）
	ColumnName   string // 列名（数据库格式，例如：group_id）
	Type         string // 字段类型（例如：int64, string）
	JSONName     string // JSON标签名（例如：group_id）
	Comment      string // 字段注释
	IsSearchable bool   // 是否可搜索
	IsRequired   bool   // 是否必填
}

// NewCRUDGenerator 创建CRUD生成器实例
func NewCRUDGenerator(moduleName, tableName, chineseName string) (*CRUDGenerator, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("获取工作目录失败：%v", err)
	}

	// 如果当前目录是 hack/generator，则需要回到项目根目录
	// 统一使用正斜杠进行比较，兼容 Windows
	normalizedPath := filepath.ToSlash(workDir)
	if strings.HasSuffix(normalizedPath, "hack/generator") {
		workDir = filepath.Join(workDir, "..", "..")
	}

	// 解析表名，生成实体名和变量名
	// 例如：system_api -> SystemApi, api
	entityName := gstr.CaseCamel(tableName)

	// 获取资源名（去除模块前缀）
	// 例如：system_api -> api
	parts := strings.Split(tableName, "_")
	var resourceName string
	if len(parts) > 1 && parts[0] == moduleName {
		resourceName = strings.Join(parts[1:], "_")
	} else {
		resourceName = tableName
	}
	varName := gstr.CaseCamelLower(resourceName)

	return &CRUDGenerator{
		ModuleName:  moduleName,
		TableName:   tableName,
		EntityName:  entityName,
		VarName:     varName,
		PackageName: moduleName,
		ChineseName: chineseName,
		WorkDir:     workDir,
	}, nil
}

// ParseEntityFields 从Entity文件解析字段信息
func (g *CRUDGenerator) ParseEntityFields() error {
	// 读取entity文件
	entityPath := filepath.Join(g.WorkDir, "internal", "model", "entity", g.TableName+".go")
	content, err := os.ReadFile(entityPath)
	if err != nil {
		return fmt.Errorf("读取entity文件失败：%v", err)
	}

	// 解析字段（简单实现，实际应用可以用AST）
	lines := strings.Split(string(content), "\n")
	inStruct := false
	g.Fields = make([]Field, 0)

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// 找到结构体定义
		if strings.Contains(line, "type "+g.EntityName+" struct") {
			inStruct = true
			continue
		}

		if inStruct {
			// 结构体结束
			if strings.HasPrefix(line, "}") {
				break
			}

			// 跳过空行和注释
			if line == "" || strings.HasPrefix(line, "//") {
				continue
			}

			// 解析字段行
			// 格式：Id int64 `json:"id" orm:"id" description:""`
			field := g.parseFieldLine(line)
			if field != nil {
				// 过滤掉不需要的字段
				if g.shouldIncludeField(field.Name) {
					g.Fields = append(g.Fields, *field)
				}
			}
		}
	}

	if len(g.Fields) == 0 {
		return fmt.Errorf("未能解析到任何字段")
	}

	return nil
}

// parseFieldLine 解析字段行
func (g *CRUDGenerator) parseFieldLine(line string) *Field {
	// 移除开头空格
	line = strings.TrimSpace(line)

	// 查找反引号位置
	backquoteIndex := strings.Index(line, "`")
	if backquoteIndex == -1 {
		return nil
	}

	// 分割字段定义和标签
	fieldDef := strings.TrimSpace(line[:backquoteIndex])
	tags := line[backquoteIndex:]

	// 解析字段名和类型
	parts := strings.Fields(fieldDef)
	if len(parts) < 2 {
		return nil
	}

	name := parts[0]
	fieldType := parts[1]

	// 解析标签
	jsonName := g.extractTag(tags, "json")
	columnName := g.extractTag(tags, "orm")
	comment := g.extractTag(tags, "description")

	// 判断是否可搜索（字符串和整数类型）
	isSearchable := isSearchableType(fieldType)

	// 判断是否必填（通过字段名判断，实际应该从数据库schema获取）
	isRequired := !strings.Contains(fieldType, "*") && name != "Id"

	return &Field{
		Name:         name,
		ColumnName:   columnName,
		Type:         fieldType,
		JSONName:     jsonName,
		Comment:      comment,
		IsSearchable: isSearchable,
		IsRequired:   isRequired,
	}
}

// extractTag 提取标签值
func (g *CRUDGenerator) extractTag(tags, key string) string {
	// 查找key的位置
	keyPattern := key + `:"`
	start := strings.Index(tags, keyPattern)
	if start == -1 {
		return ""
	}
	start += len(keyPattern)

	// 查找结束引号
	end := strings.Index(tags[start:], `"`)
	if end == -1 {
		return ""
	}

	return tags[start : start+end]
}

// shouldIncludeField 判断是否应该包含该字段
func (g *CRUDGenerator) shouldIncludeField(fieldName string) bool {
	// 排除的字段
	excludeFields := []string{"CreatedBy", "UpdatedBy", "CreatedAt", "UpdatedAt", "DeletedAt"}
	for _, exclude := range excludeFields {
		if fieldName == exclude {
			return false
		}
	}
	return true
}

// isSearchableType 判断类型是否可搜索
func isSearchableType(fieldType string) bool {
	searchableTypes := []string{"string", "int", "int64", "int32"}
	for _, t := range searchableTypes {
		if strings.Contains(fieldType, t) {
			return true
		}
	}
	return false
}

// Generate 生成所有CRUD文件
func (g *CRUDGenerator) Generate() error {
	// 1. 解析字段
	if err := g.ParseEntityFields(); err != nil {
		return err
	}

	// 2. 生成各个文件
	generators := []struct {
		name string
		fn   func() error
	}{
		{"API定义", g.GenerateAPI},
		{"请求模型", g.GenerateReq},
		{"响应模型", g.GenerateRes},
		{"控制器", g.GenerateController},
		{"逻辑层", g.GenerateLogic},
	}

	for _, gen := range generators {
		fmt.Printf("正在生成%s...\n", gen.name)
		if err := gen.fn(); err != nil {
			return fmt.Errorf("生成%s失败：%v", gen.name, err)
		}
	}

	fmt.Printf("\n✓ CRUD代码生成完成！\n")
	return nil
}

// GenerateAPI 生成API定义文件
func (g *CRUDGenerator) GenerateAPI() error {
	template := `// Package {{.PackageName}}
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package {{.PackageName}}

import (
	"devinggo/modules/{{.ModuleName}}/model"
	"devinggo/modules/{{.ModuleName}}/model/page"
	"devinggo/modules/{{.ModuleName}}/model/req"
	"devinggo/modules/{{.ModuleName}}/model/res"

	"github.com/gogf/gf/v2/frame/g"
)

type Index{{.EntityName}}Req struct {
	g.Meta ` + "`" + `path:"/{{.VarName}}/index" method:"get" tags:"{{.ChineseName}}" summary:"{{.ChineseName}}列表." x-permission:"{{.ModuleName}}:{{.VarName}}:index"` + "`" + `
	model.AuthorHeader
	model.PageListReq
	req.{{.EntityName}}Search
}

type Index{{.EntityName}}Res struct {
	g.Meta ` + "`" + `mime:"application/json"` + "`" + `
	page.PageRes
	Items []res.{{.EntityName}} ` + "`" + `json:"items"  dc:"{{.ChineseName}} list"` + "`" + `
}

type List{{.EntityName}}Req struct {
	g.Meta ` + "`" + `path:"/{{.VarName}}/list" method:"get" tags:"{{.ChineseName}}" summary:"列表，无分页.." x-permission:"{{.ModuleName}}:{{.VarName}}:list"` + "`" + `
	model.AuthorHeader
}

type List{{.EntityName}}Res struct {
	g.Meta ` + "`" + `mime:"application/json"` + "`" + `
	Data   []res.{{.EntityName}} ` + "`" + `json:"data"  dc:"list"` + "`" + `
}

type Recycle{{.EntityName}}Req struct {
	g.Meta ` + "`" + `path:"/{{.VarName}}/recycle" method:"get" tags:"{{.ChineseName}}" summary:"回收站{{.ChineseName}}列表." x-permission:"{{.ModuleName}}:{{.VarName}}:recycle"` + "`" + `
	model.AuthorHeader
	model.PageListReq
	req.{{.EntityName}}Search
}

type Recycle{{.EntityName}}Res struct {
	g.Meta ` + "`" + `mime:"application/json"` + "`" + `
	page.PageRes
	Items []res.{{.EntityName}} ` + "`" + `json:"items"  dc:"{{.ChineseName}} list"` + "`" + `
}

type Save{{.EntityName}}Req struct {
	g.Meta ` + "`" + `path:"/{{.VarName}}/save" method:"post" tags:"{{.ChineseName}}" summary:"新增{{.ChineseName}}." x-permission:"{{.ModuleName}}:{{.VarName}}:save"` + "`" + `
	model.AuthorHeader
	req.{{.EntityName}}Save
}

type Save{{.EntityName}}Res struct {
	g.Meta ` + "`" + `mime:"application/json"` + "`" + `
	Id     int64 ` + "`" + `json:"id" dc:"{{.ChineseName}} id"` + "`" + `
}

type Read{{.EntityName}}Req struct {
	g.Meta ` + "`" + `path:"/{{.VarName}}/read/{Id}" method:"get" tags:"{{.ChineseName}}" summary:"读取{{.ChineseName}}." x-permission:"{{.ModuleName}}:{{.VarName}}:read"` + "`" + `
	model.AuthorHeader
	Id int64 ` + "`" + `json:"id" dc:"{{.ChineseName}} id" v:"required|min:1#{{.ChineseName}}Id不能为空"` + "`" + `
}

type Read{{.EntityName}}Res struct {
	g.Meta ` + "`" + `mime:"application/json"` + "`" + `
	Data   res.{{.EntityName}} ` + "`" + `json:"data" dc:"{{.ChineseName}}信息"` + "`" + `
}

type Update{{.EntityName}}Req struct {
	g.Meta ` + "`" + `path:"/{{.VarName}}/update/{Id}" method:"put" tags:"{{.ChineseName}}" summary:"更新{{.ChineseName}}." x-permission:"{{.ModuleName}}:{{.VarName}}:update"` + "`" + `
	model.AuthorHeader
	req.{{.EntityName}}Update
}

type Update{{.EntityName}}Res struct {
	g.Meta ` + "`" + `mime:"application/json"` + "`" + `
}

type Delete{{.EntityName}}Req struct {
	g.Meta ` + "`" + `path:"/{{.VarName}}/delete" method:"delete" tags:"{{.ChineseName}}" summary:"删除{{.ChineseName}}" x-permission:"{{.ModuleName}}:{{.VarName}}:delete"` + "`" + `
	model.AuthorHeader
	Ids []int64 ` + "`" + `json:"ids" dc:"ids" v:"min-length:1#{{.ChineseName}}Id不能为空"` + "`" + `
}

type Delete{{.EntityName}}Res struct {
	g.Meta ` + "`" + `mime:"application/json"` + "`" + `
}

type RealDelete{{.EntityName}}Req struct {
	g.Meta ` + "`" + `path:"/{{.VarName}}/realDelete" method:"delete" tags:"{{.ChineseName}}" summary:"单个或批量真实删除{{.ChineseName}} （清空回收站）." x-permission:"{{.ModuleName}}:{{.VarName}}:realDelete"` + "`" + `
	model.AuthorHeader
	Ids []int64 ` + "`" + `json:"ids" dc:"ids" v:"min-length:1#{{.ChineseName}}Id不能为空"` + "`" + `
}

type RealDelete{{.EntityName}}Res struct {
	g.Meta ` + "`" + `mime:"application/json"` + "`" + `
}

type Recovery{{.EntityName}}Req struct {
	g.Meta ` + "`" + `path:"/{{.VarName}}/recovery" method:"put" tags:"{{.ChineseName}}" summary:"单个或批量恢复在回收站的{{.ChineseName}}." x-permission:"{{.ModuleName}}:{{.VarName}}:recovery"` + "`" + `
	model.AuthorHeader
	Ids []int64 ` + "`" + `json:"ids" dc:"ids" v:"min-length:1#{{.ChineseName}}Id不能为空"` + "`" + `
}

type Recovery{{.EntityName}}Res struct {
	g.Meta ` + "`" + `mime:"application/json"` + "`" + `
}

type ChangeStatus{{.EntityName}}Req struct {
	g.Meta ` + "`" + `path:"/{{.VarName}}/changeStatus" method:"put" tags:"{{.ChineseName}}" summary:"更改状态" x-permission:"{{.ModuleName}}:{{.VarName}}:update"` + "`" + `
	model.AuthorHeader
	Id     int64 ` + "`" + `json:"id" dc:"ids" v:"min:1#Id不能为空"` + "`" + `
	Status int   ` + "`" + `json:"status" dc:"status" v:"min:1#状态不能为空"` + "`" + `
}

type ChangeStatus{{.EntityName}}Res struct {
	g.Meta ` + "`" + `mime:"application/json"` + "`" + `
}
`

	outputPath := filepath.Join(g.WorkDir, "modules", g.ModuleName, "api", g.PackageName, g.VarName+".go")
	return g.renderAndSaveTemplate(template, outputPath)
}

// GenerateReq 生成请求模型文件
func (g *CRUDGenerator) GenerateReq() error {
	// 构建Search字段
	var searchFields strings.Builder
	for _, field := range g.Fields {
		if field.IsSearchable && field.Name != "Id" {
			searchFields.WriteString(fmt.Sprintf("\t%s %s `json:\"%s\"`\n",
				field.Name, field.Type, field.JSONName))
		}
	}

	// 构建Save字段
	var saveFields strings.Builder
	for _, field := range g.Fields {
		if field.Name == "Id" {
			continue
		}
		required := ""
		if field.IsRequired {
			required = ` v:"required"`
		}
		saveFields.WriteString(fmt.Sprintf("\t%s %s `json:\"%s\"%s description:\"%s\"`\n",
			field.Name, field.Type, field.JSONName, required, field.Comment))
	}

	// 构建Update字段（包含Id）
	var updateFields strings.Builder
	for _, field := range g.Fields {
		required := ""
		if field.Name == "Id" {
			required = ` v:"required"`
		} else if field.IsRequired {
			required = ` v:"required"`
		}
		updateFields.WriteString(fmt.Sprintf("\t%s %s `json:\"%s\"%s description:\"%s\"`\n",
			field.Name, field.Type, field.JSONName, required, field.Comment))
	}

	template := `// Package req
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package req

type {{.EntityName}}Search struct {
{{.SearchFields}}}

type {{.EntityName}}Save struct {
{{.SaveFields}}}

type {{.EntityName}}Update struct {
{{.UpdateFields}}}
`

	data := map[string]string{
		"EntityName":   g.EntityName,
		"SearchFields": searchFields.String(),
		"SaveFields":   saveFields.String(),
		"UpdateFields": updateFields.String(),
	}

	outputPath := filepath.Join(g.WorkDir, "modules", g.ModuleName, "model", "req", g.TableName+".go")
	return g.renderAndSaveTemplateWithData(template, outputPath, data)
}

// GenerateRes 生成响应模型文件
func (g *CRUDGenerator) GenerateRes() error {
	// 构建字段（包含所有字段，包括时间戳）
	var fields strings.Builder
	// 先添加ID
	fields.WriteString(fmt.Sprintf("\tId %s `json:\"%s\" description:\"%s\"`\n",
		"int64", "id", "主键"))

	// 添加其他字段
	for _, field := range g.Fields {
		if field.Name == "Id" {
			continue
		}
		fields.WriteString(fmt.Sprintf("\t%s %s `json:\"%s\" description:\"%s\"`\n",
			field.Name, field.Type, field.JSONName, field.Comment))
	}

	// 添加时间戳字段
	timestampFields := []string{
		"CreatedBy   int64       `json:\"created_by\" description:\"创建者\"`",
		"UpdatedBy   int64       `json:\"updated_by\" description:\"更新者\"`",
		"CreatedAt   *gtime.Time `json:\"created_at\" description:\"创建时间\"`",
		"UpdatedAt   *gtime.Time `json:\"updated_at\" description:\"更新时间\"`",
	}
	for _, ts := range timestampFields {
		fields.WriteString("\t" + ts + "\n")
	}

	template := `// Package res
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package res

import "github.com/gogf/gf/v2/os/gtime"

type {{.EntityName}} struct {
{{.Fields}}}
`

	data := map[string]string{
		"EntityName": g.EntityName,
		"Fields":     fields.String(),
	}

	outputPath := filepath.Join(g.WorkDir, "modules", g.ModuleName, "model", "res", g.TableName+".go")
	return g.renderAndSaveTemplateWithData(template, outputPath, data)
}

// GenerateController 生成控制器文件
func (g *CRUDGenerator) GenerateController() error {
	template := `// Package {{.PackageName}}
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package {{.PackageName}}

import (
	"context"
	"devinggo/internal/dao"
	"devinggo/modules/{{.ModuleName}}/api/{{.PackageName}}"
	"devinggo/modules/{{.ModuleName}}/controller/base"
	"devinggo/modules/{{.ModuleName}}/model/res"
	"devinggo/modules/{{.ModuleName}}/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

var (
	{{.EntityName}}Controller = {{.VarName}}Controller{}
)

type {{.VarName}}Controller struct {
	base.BaseController
}

func (c *{{.VarName}}Controller) Index(ctx context.Context, in *{{.PackageName}}.Index{{.EntityName}}Req) (out *{{.PackageName}}.Index{{.EntityName}}Res, err error) {
	out = &{{.PackageName}}.Index{{.EntityName}}Res{}
	items, totalCount, err := service.{{.EntityName}}().GetPageListForSearch(ctx, &in.PageListReq, &in.{{.EntityName}}Search)
	if err != nil {
		return
	}

	if !g.IsEmpty(items) {
		for _, item := range items {
			out.Items = append(out.Items, *item)
		}
	} else {
		out.Items = make([]res.{{.EntityName}}, 0)
	}
	out.PageRes.Pack(in, totalCount)
	return
}

func (c *{{.VarName}}Controller) List(ctx context.Context, in *{{.PackageName}}.List{{.EntityName}}Req) (out *{{.PackageName}}.List{{.EntityName}}Res, err error) {
	out = &{{.PackageName}}.List{{.EntityName}}Res{}
	items, err := service.{{.EntityName}}().GetList(ctx, &{{.PackageName}}.{{.EntityName}}Search{})
	if err != nil {
		return
	}

	if !g.IsEmpty(items) {
		for _, item := range items {
			out.Data = append(out.Data, *item)
		}
	} else {
		out.Data = make([]res.{{.EntityName}}, 0)
	}
	return
}

func (c *{{.VarName}}Controller) Recycle(ctx context.Context, in *{{.PackageName}}.Recycle{{.EntityName}}Req) (out *{{.PackageName}}.Recycle{{.EntityName}}Res, err error) {
	out = &{{.PackageName}}.Recycle{{.EntityName}}Res{}
	in.Recycle = true
	items, totalCount, err := service.{{.EntityName}}().GetPageListForSearch(ctx, &in.PageListReq, &in.{{.EntityName}}Search)
	if err != nil {
		return
	}

	if !g.IsEmpty(items) {
		for _, item := range items {
			out.Items = append(out.Items, *item)
		}
	} else {
		out.Items = make([]res.{{.EntityName}}, 0)
	}
	out.PageRes.Pack(in, totalCount)
	return
}

func (c *{{.VarName}}Controller) Save(ctx context.Context, in *{{.PackageName}}.Save{{.EntityName}}Req) (out *{{.PackageName}}.Save{{.EntityName}}Res, err error) {
	out = &{{.PackageName}}.Save{{.EntityName}}Res{}
	id, err := service.{{.EntityName}}().Save(ctx, &in.{{.EntityName}}Save)
	if err != nil {
		return
	}
	out.Id = id
	return
}

func (c *{{.VarName}}Controller) Read(ctx context.Context, in *{{.PackageName}}.Read{{.EntityName}}Req) (out *{{.PackageName}}.Read{{.EntityName}}Res, err error) {
	out = &{{.PackageName}}.Read{{.EntityName}}Res{}
	info, err := service.{{.EntityName}}().GetById(ctx, in.Id)
	if err != nil {
		return
	}
	out.Data = *info
	return
}

func (c *{{.VarName}}Controller) Update(ctx context.Context, in *{{.PackageName}}.Update{{.EntityName}}Req) (out *{{.PackageName}}.Update{{.EntityName}}Res, err error) {
	err = dao.{{.EntityName}}.Transaction(ctx, func(ctx context.Context, tx gdb.TX) (err error) {
		out = &{{.PackageName}}.Update{{.EntityName}}Res{}
		err = service.{{.EntityName}}().Update(ctx, &in.{{.EntityName}}Update)
		if err != nil {
			return
		}
		return
	})
	if err != nil {
		return
	}
	return
}

func (c *{{.VarName}}Controller) Delete(ctx context.Context, in *{{.PackageName}}.Delete{{.EntityName}}Req) (out *{{.PackageName}}.Delete{{.EntityName}}Res, err error) {
	out = &{{.PackageName}}.Delete{{.EntityName}}Res{}
	err = service.{{.EntityName}}().Delete(ctx, in.Ids)
	if err != nil {
		return
	}
	return
}

func (c *{{.VarName}}Controller) RealDelete(ctx context.Context, in *{{.PackageName}}.RealDelete{{.EntityName}}Req) (out *{{.PackageName}}.RealDelete{{.EntityName}}Res, err error) {
	err = dao.{{.EntityName}}.Transaction(ctx, func(ctx context.Context, tx gdb.TX) (err error) {
		out = &{{.PackageName}}.RealDelete{{.EntityName}}Res{}
		err = service.{{.EntityName}}().RealDelete(ctx, in.Ids)
		if err != nil {
			return
		}
		return
	})
	if err != nil {
		return
	}
	return
}

func (c *{{.VarName}}Controller) Recovery(ctx context.Context, in *{{.PackageName}}.Recovery{{.EntityName}}Req) (out *{{.PackageName}}.Recovery{{.EntityName}}Res, err error) {
	err = dao.{{.EntityName}}.Transaction(ctx, func(ctx context.Context, tx gdb.TX) (err error) {
		out = &{{.PackageName}}.Recovery{{.EntityName}}Res{}
		err = service.{{.EntityName}}().Recovery(ctx, in.Ids)
		if err != nil {
			return
		}
		return
	})
	if err != nil {
		return
	}
	return
}

func (c *{{.VarName}}Controller) ChangeStatus(ctx context.Context, in *{{.PackageName}}.ChangeStatus{{.EntityName}}Req) (out *{{.PackageName}}.ChangeStatus{{.EntityName}}Res, err error) {
	out = &{{.PackageName}}.ChangeStatus{{.EntityName}}Res{}
	err = service.{{.EntityName}}().ChangeStatus(ctx, in.Id, in.Status)
	if err != nil {
		return
	}
	return
}
`

	outputPath := filepath.Join(g.WorkDir, "modules", g.ModuleName, "controller", g.PackageName, g.VarName+".go")
	return g.renderAndSaveTemplate(template, outputPath)
}

// GenerateLogic 生成逻辑层文件
func (g *CRUDGenerator) GenerateLogic() error {
	// 构建搜索条件
	var searchConditions strings.Builder
	for _, field := range g.Fields {
		if field.IsSearchable && field.Name != "Id" {
			searchConditions.WriteString(fmt.Sprintf(`
	if !g.IsEmpty(in.%s) {
		m = m.Where("%s", in.%s)
	}
`, field.Name, field.ColumnName, field.Name))
		}
	}

	// 构建保存字段（使用DO）
	var saveDoFields strings.Builder
	for _, field := range g.Fields {
		if field.Name == "Id" {
			continue
		}
		saveDoFields.WriteString(fmt.Sprintf("\t\t%s: in.%s,\n", field.Name, field.Name))
	}

	// 构建更新字段（使用DO）
	var updateDoFields strings.Builder
	for _, field := range g.Fields {
		if field.Name == "Id" {
			continue
		}
		updateDoFields.WriteString(fmt.Sprintf("\t\t%s: in.%s,\n", field.Name, field.Name))
	}

	template := `// Package {{.PackageName}}
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package {{.PackageName}}

import (
	"context"
	"devinggo/internal/dao"
	"devinggo/internal/model/do"
	"devinggo/internal/model/entity"
	"devinggo/modules/{{.ModuleName}}/logic/base"
	"devinggo/modules/{{.ModuleName}}/model"
	"devinggo/modules/{{.ModuleName}}/model/req"
	"devinggo/modules/{{.ModuleName}}/model/res"
	"devinggo/modules/{{.ModuleName}}/pkg/handler"
	"devinggo/modules/{{.ModuleName}}/pkg/hook"
	"devinggo/modules/{{.ModuleName}}/pkg/orm"
	"devinggo/modules/{{.ModuleName}}/pkg/utils"
	"devinggo/modules/{{.ModuleName}}/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type s{{.EntityName}} struct {
	base.BaseService
}

func init() {
	service.Register{{.EntityName}}(New{{.EntityName}}())
}

func New{{.EntityName}}() *s{{.EntityName}} {
	return &s{{.EntityName}}{}
}

func (s *s{{.EntityName}}) Model(ctx context.Context) *gdb.Model {
	return dao.{{.EntityName}}.Ctx(ctx).Hook(hook.Default()).Cache(orm.SetCacheOption(ctx)).OnConflict("id")
}

func (s *s{{.EntityName}}) GetPageListForSearch(ctx context.Context, req *model.PageListReq, in *req.{{.EntityName}}Search) (rs []*res.{{.EntityName}}, total int, err error) {
	m := s.handleSearch(ctx, in)
	var entity []*entity.{{.EntityName}}
	err = orm.NewQuery(m).WithPageListReq(req).ScanAndCount(&entity, &total)
	if utils.IsError(err) {
		return nil, 0, err
	}
	rs = make([]*res.{{.EntityName}}, 0)
	if !g.IsEmpty(entity) {
		if err = gconv.Structs(entity, &rs); err != nil {
			return nil, 0, err
		}
	}
	return
}

func (s *s{{.EntityName}}) GetList(ctx context.Context, in *req.{{.EntityName}}Search) (out []*res.{{.EntityName}}, err error) {
	inReq := &model.ListReq{
		OrderBy:   dao.{{.EntityName}}.Table() + ".created_at",
		OrderType: "desc",
	}
	m := s.handleSearch(ctx, in).Handler(handler.FilterAuth)
	m = orm.NewQuery(m).WithListReq(inReq).Build()
	err = m.Scan(&out)
	if utils.IsError(err) {
		return
	}
	return
}

func (s *s{{.EntityName}}) handleSearch(ctx context.Context, in *req.{{.EntityName}}Search) (m *gdb.Model) {
	m = s.Model(ctx){{.SearchConditions}}
	return
}

func (s *s{{.EntityName}}) Save(ctx context.Context, in *req.{{.EntityName}}Save) (id int64, err error) {
	saveData := do.{{.EntityName}}{
{{.SaveDoFields}}	}
	id, err = orm.NewQuery(s.Model(ctx)).Data(&saveData).InsertAndGetId()
	if utils.IsError(err) {
		return 0, err
	}
	return
}

func (s *s{{.EntityName}}) GetById(ctx context.Context, id int64) (out *res.{{.EntityName}}, err error) {
	var entity *entity.{{.EntityName}}
	err = s.Model(ctx).Where("id", id).Scan(&entity)
	if utils.IsError(err) {
		return nil, err
	}
	out = &res.{{.EntityName}}{}
	if err = gconv.Struct(entity, out); err != nil {
		return nil, err
	}
	return
}

func (s *s{{.EntityName}}) Update(ctx context.Context, in *req.{{.EntityName}}Update) (err error) {
	updateData := do.{{.EntityName}}{
{{.UpdateDoFields}}	}
	_, err = orm.NewQuery(s.Model(ctx)).Data(&updateData).Where("id", in.Id).Update()
	if utils.IsError(err) {
		return err
	}
	return
}

func (s *s{{.EntityName}}) Delete(ctx context.Context, ids []int64) (err error) {
	_, err = s.Model(ctx).WhereIn("id", ids).Delete()
	if utils.IsError(err) {
		return err
	}
	return
}

func (s *s{{.EntityName}}) RealDelete(ctx context.Context, ids []int64) (err error) {
	_, err = s.Model(ctx).Unscoped().WhereIn("id", ids).Delete()
	if utils.IsError(err) {
		return err
	}
	return
}

func (s *s{{.EntityName}}) Recovery(ctx context.Context, ids []int64) (err error) {
	_, err = s.Model(ctx).Unscoped().WhereIn("id", ids).Update(g.Map{
		"deleted_at": nil,
	})
	if utils.IsError(err) {
		return err
	}
	return
}

func (s *s{{.EntityName}}) ChangeStatus(ctx context.Context, id int64, status int) (err error) {
	_, err = s.Model(ctx).Where("id", id).Update(g.Map{
		"status": status,
	})
	if utils.IsError(err) {
		return err
	}
	return
}
`

	data := map[string]string{
		"PackageName":      g.PackageName,
		"ModuleName":       g.ModuleName,
		"EntityName":       g.EntityName,
		"SearchConditions": searchConditions.String(),
		"SaveDoFields":     saveDoFields.String(),
		"UpdateDoFields":   updateDoFields.String(),
	}

	outputPath := filepath.Join(g.WorkDir, "modules", g.ModuleName, "logic", g.PackageName, g.TableName+".go")
	return g.renderAndSaveTemplateWithData(template, outputPath, data)
}

// renderAndSaveTemplate 渲染并保存模板（使用g自身属性）
func (g *CRUDGenerator) renderAndSaveTemplate(template string, outputPath string) error {
	data := map[string]string{
		"ModuleName":  g.ModuleName,
		"TableName":   g.TableName,
		"EntityName":  g.EntityName,
		"VarName":     g.VarName,
		"PackageName": g.PackageName,
		"ChineseName": g.ChineseName,
	}
	return g.renderAndSaveTemplateWithData(template, outputPath, data)
}

// renderAndSaveTemplateWithData 渲染并保存模板（使用自定义数据）
func (g *CRUDGenerator) renderAndSaveTemplateWithData(template string, outputPath string, data map[string]string) error {
	// 简单的字符串替换渲染
	result := template
	for key, value := range data {
		placeholder := "{{." + key + "}}"
		result = strings.ReplaceAll(result, placeholder, value)
	}

	// 确保目录存在
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败：%v", err)
	}

	// 写入文件
	if err := os.WriteFile(outputPath, []byte(result), 0644); err != nil {
		return fmt.Errorf("写入文件失败：%v", err)
	}

	fmt.Printf("  ✓ 已生成：%s\n", outputPath)
	return nil
}
