// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CodeGenTables is the golang structure for table code_gen_tables.
type CodeGenTables struct {
	Id            int64        `json:"id"             orm:"id"             description:"主键ID"`
	TableName     string       `json:"table_name"     orm:"table_name"     description:"表名称"`
	TableComment  string       `json:"table_comment"  orm:"table_comment"  description:"表描述"`
	Remark        string       `json:"remark"         orm:"remark"         description:"备注信息"`
	ModuleName    string       `json:"module_name"    orm:"module_name"    description:"所属模块"`
	BelongMenuId  int64        `json:"belong_menu_id" orm:"belong_menu_id" description:"所属菜单ID"`
	Type          string       `json:"type"           orm:"type"           description:"生成类型: single=单表, tree=树表"`
	MenuName      string       `json:"menu_name"       orm:"menu_name"      description:"菜单名称"`
	ComponentType int          `json:"component_type" orm:"component_type" description:"组件类型: 1=模态框, 2=抽屉, 3=Tag页"`
	TplType       string       `json:"tpl_type"       orm:"tpl_type"       description:"模板类型: default"`
	TreeId        string       `json:"tree_id"        orm:"tree_id"        description:"树表主ID字段"`
	TreeParentId  string       `json:"tree_parent_id" orm:"tree_parent_id" description:"树表父ID字段"`
	TreeName      string       `json:"tree_name"      orm:"tree_name"      description:"树表显示名称字段"`
	TagId         string       `json:"tag_id"         orm:"tag_id"         description:"Tag页ID"`
	TagName       string       `json:"tag_name"       orm:"tag_name"       description:"Tag页名称"`
	TagViewName   string       `json:"tag_view_name"  orm:"tag_view_name"  description:"Tag页显示字段"`
	GenerateMenus string       `json:"generate_menus" orm:"generate_menus" description:"生成的菜单按钮"`
	Options       string       `json:"options"        orm:"options"        description:"扩展配置"`
	CreatedAt     *gtime.Time  `json:"created_at"     orm:"created_at"     description:"创建时间"`
	UpdatedAt     *gtime.Time  `json:"updated_at"     orm:"updated_at"     description:"更新时间"`
	DeletedAt     *gtime.Time  `json:"deleted_at"     orm:"deleted_at"     description:"删除时间"`
	CreatedBy     int64        `json:"created_by"     orm:"created_by"     description:"创建者ID"`
	UpdatedBy     int64        `json:"updated_by"     orm:"updated_by"     description:"更新者ID"`
	Status        int          `json:"status"         orm:"status"         description:"状态: 1=正常, 0=停用"`
	Sort          int          `json:"sort"           orm:"sort"           description:"排序"`
}
