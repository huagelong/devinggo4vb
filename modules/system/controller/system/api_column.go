// Package system
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package system

import (
	"context"
	"devinggo/modules/system/api/system"
	"devinggo/modules/system/controller/base"
	"devinggo/modules/system/model/req"
	"devinggo/modules/system/model/res"
	"devinggo/modules/system/myerror"
	"devinggo/modules/system/pkg/hook"
	"devinggo/modules/system/pkg/orm"
	"devinggo/modules/system/pkg/utils"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

const apiColumnTableName = "system_api_column"

var (
	ApiColumnController = apiColumnController{}
)

type apiColumnController struct {
	base.BaseController
}

func (c *apiColumnController) model(ctx context.Context) *gdb.Model {
	return g.DB().Model(apiColumnTableName).Ctx(ctx).Hook(hook.Default()).Cache(orm.SetCacheOption(ctx)).OnConflict("id")
}

func (c *apiColumnController) handleSearch(ctx context.Context, in *req.SystemApiColumnSearch) (m *gdb.Model) {
	m = c.model(ctx)
	if !g.IsEmpty(in.ApiId) {
		m = m.Where("api_id", in.ApiId)
	}
	if !g.IsEmpty(in.Name) {
		m = m.WhereLike("name", "%"+in.Name+"%")
	}
	if !g.IsEmpty(in.DataType) {
		m = m.Where("data_type", in.DataType)
	}
	if !g.IsEmpty(in.IsRequired) {
		m = m.Where("is_required", in.IsRequired)
	}
	if !g.IsEmpty(in.Status) {
		m = m.Where("status", in.Status)
	}
	if !g.IsEmpty(in.Type) {
		m = m.Where("type", in.Type)
	}
	if !g.IsEmpty(in.CreatedAt) {
		if len(in.CreatedAt) > 0 {
			m = m.WhereGTE("created_at", in.CreatedAt[0]+" 00:00:00")
		}
		if len(in.CreatedAt) > 1 {
			m = m.WhereLTE("created_at", in.CreatedAt[1]+" 23:59:59")
		}
	}
	return
}

func (c *apiColumnController) Index(ctx context.Context, in *system.IndexApiColumnReq) (out *system.IndexApiColumnRes, err error) {
	out = &system.IndexApiColumnRes{}
	items := make([]res.SystemApiColumn, 0)
	m := c.handleSearch(ctx, &in.SystemApiColumnSearch)
	err = orm.NewQuery(m).WithPageListReq(&in.PageListReq).ScanAndCount(&items, &out.PageInfo.TotalCount)
	if utils.IsError(err) {
		return
	}
	out.Items = items
	out.PageRes.Pack(in, out.PageInfo.TotalCount)
	return
}

func (c *apiColumnController) Recycle(ctx context.Context, in *system.RecycleApiColumnReq) (out *system.RecycleApiColumnRes, err error) {
	out = &system.RecycleApiColumnRes{}
	in.Recycle = true
	items := make([]res.SystemApiColumn, 0)
	m := c.handleSearch(ctx, &in.SystemApiColumnSearch)
	err = orm.NewQuery(m).WithPageListReq(&in.PageListReq).ScanAndCount(&items, &out.PageInfo.TotalCount)
	if utils.IsError(err) {
		return
	}
	out.Items = items
	out.PageRes.Pack(in, out.PageInfo.TotalCount)
	return
}

func (c *apiColumnController) Read(ctx context.Context, in *system.ReadApiColumnReq) (out *system.ReadApiColumnRes, err error) {
	out = &system.ReadApiColumnRes{}
	err = c.model(ctx).Where("id", in.Id).Scan(&out.Data)
	if utils.IsError(err) {
		return
	}
	return
}

func (c *apiColumnController) Save(ctx context.Context, in *system.SaveApiColumnReq) (out *system.SaveApiColumnRes, err error) {
	out = &system.SaveApiColumnRes{}
	status := in.Status
	if status == 0 {
		status = 1
	}
	isRequired := in.IsRequired
	if isRequired == 0 {
		isRequired = 2
	}
	fieldType := in.Type
	if fieldType == 0 {
		fieldType = 1
	}
	saveData := g.Map{
		"api_id":        in.ApiId,
		"name":          in.Name,
		"data_type":     gconv.String(in.DataType),
		"is_required":   isRequired,
		"status":        status,
		"type":          fieldType,
		"default_value": in.DefaultValue,
		"description":   in.Description,
		"remark":        in.Remark,
	}
	rs, err := c.model(ctx).Data(saveData).Insert()
	if utils.IsError(err) {
		return
	}
	lastInsertId, err := rs.LastInsertId()
	if err != nil {
		return
	}
	out.Id = gconv.Int64(lastInsertId)
	return
}

func (c *apiColumnController) Update(ctx context.Context, in *system.UpdateApiColumnReq) (out *system.UpdateApiColumnRes, err error) {
	out = &system.UpdateApiColumnRes{}
	updateData := g.Map{
		"api_id":        in.ApiId,
		"name":          in.Name,
		"data_type":     gconv.String(in.DataType),
		"is_required":   in.IsRequired,
		"status":        in.Status,
		"type":          in.Type,
		"default_value": in.DefaultValue,
		"description":   in.Description,
		"remark":        in.Remark,
	}
	_, err = c.model(ctx).Data(updateData).Where("id", in.Id).Update()
	if utils.IsError(err) {
		return
	}
	return
}

func (c *apiColumnController) Delete(ctx context.Context, in *system.DeleteApiColumnReq) (out *system.DeleteApiColumnRes, err error) {
	out = &system.DeleteApiColumnRes{}
	_, err = c.model(ctx).WhereIn("id", in.Ids).Delete()
	if utils.IsError(err) {
		return
	}
	return
}

func (c *apiColumnController) RealDelete(ctx context.Context, in *system.RealDeleteApiColumnReq) (out *system.RealDeleteApiColumnRes, err error) {
	out = &system.RealDeleteApiColumnRes{}
	_, err = c.model(ctx).Unscoped().WhereIn("id", in.Ids).Delete()
	if utils.IsError(err) {
		return
	}
	return
}

func (c *apiColumnController) Recovery(ctx context.Context, in *system.RecoveryApiColumnReq) (out *system.RecoveryApiColumnRes, err error) {
	out = &system.RecoveryApiColumnRes{}
	_, err = c.model(ctx).Unscoped().WhereIn("id", in.Ids).Update(g.Map{
		"deleted_at": nil,
	})
	if utils.IsError(err) {
		return
	}
	return
}

func (c *apiColumnController) ChangeStatus(ctx context.Context, in *system.ChangeStatusApiColumnReq) (out *system.ChangeStatusApiColumnRes, err error) {
	out = &system.ChangeStatusApiColumnRes{}
	_, err = c.model(ctx).Data(g.Map{"status": in.Status}).Where("id", in.Id).Update()
	if utils.IsError(err) {
		return
	}
	return
}

func (c *apiColumnController) Import(ctx context.Context, in *system.ImportApiColumnReq) (out *system.ImportApiColumnRes, err error) {
	return nil, myerror.ValidationFailed(ctx, "接口参数导入功能暂未开放")
}

func (c *apiColumnController) Export(ctx context.Context, in *system.ExportApiColumnReq) (out *system.ExportApiColumnRes, err error) {
	return nil, myerror.ValidationFailed(ctx, "接口参数导出功能暂未开放")
}

func (c *apiColumnController) DownloadTemplate(ctx context.Context, in *system.DownloadTemplateApiColumnReq) (out *system.DownloadTemplateApiColumnRes, err error) {
	return nil, myerror.ValidationFailed(ctx, "接口参数模板下载功能暂未开放")
}
