// Package req
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package req

type SystemApiColumnSearch struct {
	ApiId      int64    `json:"api_id"`
	Name       string   `json:"name"`
	DataType   string   `json:"data_type"`
	IsRequired int      `json:"is_required"`
	Status     int      `json:"status"`
	Type       int      `json:"type"`
	CreatedAt  []string `json:"created_at" dc:"created at"`
}

type SystemApiColumnSave struct {
	ApiId        int64  `json:"api_id" v:"required|min:1#接口ID不能为空"`
	Name         string `json:"name" v:"required#字段名称不能为空"`
	DataType     any    `json:"data_type" v:"required#数据类型不能为空"`
	IsRequired   int    `json:"is_required"`
	Status       int    `json:"status"`
	Type         int    `json:"type"`
	DefaultValue string `json:"default_value"`
	Description  string `json:"description"`
	Remark       string `json:"remark"`
}

type SystemApiColumnUpdate struct {
	Id int64 `json:"id" v:"required|min:1#ID不能为空"`
	SystemApiColumnSave
}
