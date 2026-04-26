// Package res
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package res

import "github.com/gogf/gf/v2/os/gtime"

type SystemApiColumn struct {
	Id           int64       `json:"id" description:"主键"`
	ApiId        int64       `json:"api_id" description:"接口ID"`
	Name         string      `json:"name" description:"字段名称"`
	DataType     string      `json:"data_type" description:"数据类型"`
	Type         int         `json:"type" description:"字段类型(1请求参数 2响应参数)"`
	IsRequired   int         `json:"is_required" description:"是否必填(1否 2是)"`
	DefaultValue string      `json:"default_value" description:"默认值"`
	Description  string      `json:"description" description:"字段说明"`
	Status       int         `json:"status" description:"状态(1正常 2停用)"`
	CreatedBy    int64       `json:"created_by" description:"创建者"`
	UpdatedBy    int64       `json:"updated_by" description:"更新者"`
	CreatedAt    *gtime.Time `json:"created_at" description:"创建时间"`
	UpdatedAt    *gtime.Time `json:"updated_at" description:"更新时间"`
	Remark       string      `json:"remark" description:"备注"`
}
