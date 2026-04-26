// Package system
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package system

import (
	"devinggo/modules/system/model"
	"devinggo/modules/system/model/page"
	"devinggo/modules/system/model/req"
	"devinggo/modules/system/model/res"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type IndexApiColumnReq struct {
	g.Meta `path:"/apiColumn/index" method:"get" tags:"接口参数" summary:"接口参数列表." x-permission:"system:api:index" `
	model.AuthorHeader
	model.PageListReq
	req.SystemApiColumnSearch
}

type IndexApiColumnRes struct {
	g.Meta `mime:"application/json"`
	page.PageRes
	Items []res.SystemApiColumn `json:"items" dc:"api column list"`
}

type RecycleApiColumnReq struct {
	g.Meta `path:"/apiColumn/recycle" method:"get" tags:"接口参数" summary:"接口参数回收站列表." x-permission:"system:api:recycle" `
	model.AuthorHeader
	model.PageListReq
	req.SystemApiColumnSearch
}

type RecycleApiColumnRes struct {
	g.Meta `mime:"application/json"`
	page.PageRes
	Items []res.SystemApiColumn `json:"items" dc:"api column list"`
}

type ReadApiColumnReq struct {
	g.Meta `path:"/apiColumn/read/{Id}" method:"post,get" tags:"接口参数" summary:"获取接口参数详情." x-permission:"system:api:read"`
	model.AuthorHeader
	Id int64 `json:"id" dc:"接口参数 id" v:"required|min:1#接口参数Id不能为空"`
}

type ReadApiColumnRes struct {
	g.Meta `mime:"application/json"`
	Data   res.SystemApiColumn `json:"data" dc:"接口参数信息"`
}

type SaveApiColumnReq struct {
	g.Meta `path:"/apiColumn/save" method:"post" tags:"接口参数" summary:"新增接口参数." x-permission:"system:api:save"`
	model.AuthorHeader
	req.SystemApiColumnSave
}

type SaveApiColumnRes struct {
	g.Meta `mime:"application/json"`
	Id     int64 `json:"id" dc:"接口参数 id"`
}

type UpdateApiColumnReq struct {
	g.Meta `path:"/apiColumn/update/{Id}" method:"put" tags:"接口参数" summary:"更新接口参数." x-permission:"system:api:update"`
	model.AuthorHeader
	req.SystemApiColumnUpdate
}

type UpdateApiColumnRes struct {
	g.Meta `mime:"application/json"`
}

type DeleteApiColumnReq struct {
	g.Meta `path:"/apiColumn/delete" method:"delete" tags:"接口参数" summary:"删除接口参数." x-permission:"system:api:delete"`
	model.AuthorHeader
	Ids []int64 `json:"ids" dc:"ids" v:"min-length:1#接口参数Id不能为空"`
}

type DeleteApiColumnRes struct {
	g.Meta `mime:"application/json"`
}

type RealDeleteApiColumnReq struct {
	g.Meta `path:"/apiColumn/realDelete" method:"delete" tags:"接口参数" summary:"彻底删除接口参数." x-permission:"system:api:realDelete"`
	model.AuthorHeader
	Ids []int64 `json:"ids" dc:"ids" v:"min-length:1#接口参数Id不能为空"`
}

type RealDeleteApiColumnRes struct {
	g.Meta `mime:"application/json"`
}

type RecoveryApiColumnReq struct {
	g.Meta `path:"/apiColumn/recovery" method:"put" tags:"接口参数" summary:"恢复接口参数." x-permission:"system:api:recovery"`
	model.AuthorHeader
	Ids []int64 `json:"ids" dc:"ids" v:"min-length:1#接口参数Id不能为空"`
}

type RecoveryApiColumnRes struct {
	g.Meta `mime:"application/json"`
}

type ChangeStatusApiColumnReq struct {
	g.Meta `path:"/apiColumn/changeStatus" method:"put" tags:"接口参数" summary:"修改接口参数状态." x-permission:"system:api:update"`
	model.AuthorHeader
	Id     int64 `json:"id" dc:"id" v:"min:1#Id不能为空"`
	Status int   `json:"status" dc:"status" v:"min:1#状态不能为空"`
}

type ChangeStatusApiColumnRes struct {
	g.Meta `mime:"application/json"`
}

type ImportApiColumnReq struct {
	g.Meta `path:"/apiColumn/import" method:"post" mime:"multipart/form-data" tags:"接口参数" summary:"导入接口参数." x-permission:"system:api:update"`
	model.AuthorHeader
	ApiId int64             `json:"api_id" v:"required|min:1#接口ID不能为空"`
	Type  int               `json:"type" v:"required|min:1#字段类型不能为空"`
	File  *ghttp.UploadFile `json:"file" type:"file" dc:"please upload file"`
}

type ImportApiColumnRes struct {
	g.Meta `mime:"application/json"`
}

type ExportApiColumnReq struct {
	g.Meta `path:"/apiColumn/export" method:"post" tags:"接口参数" summary:"导出接口参数." x-permission:"system:api:read"`
	model.AuthorHeader
	model.ListReq
	req.SystemApiColumnSearch
}

type ExportApiColumnRes struct {
	g.Meta `mime:"application/json"`
}

type DownloadTemplateApiColumnReq struct {
	g.Meta `path:"/apiColumn/downloadTemplate" method:"post,get" tags:"接口参数" summary:"下载接口参数导入模板." x-permission:"system:api:read"`
	model.AuthorHeader
}

type DownloadTemplateApiColumnRes struct {
	g.Meta `mime:"application/json"`
}
