// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CodeGenTablesDao is the data access object for table code_gen_tables.
type CodeGenTablesDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  CodeGenTablesColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// CodeGenTablesColumns defines and stores column names for table code_gen_tables.
type CodeGenTablesColumns struct {
	Id            string
	TableName     string
	TableComment  string
	Remark        string
	ModuleName    string
	BelongMenuId  string
	Type          string
	MenuName      string
	ComponentType string
	TplType       string
	TreeId        string
	TreeParentId  string
	TreeName      string
	TagId         string
	TagName       string
	TagViewName   string
	GenerateMenus string
	Options       string
	CreatedAt     string
	UpdatedAt     string
	DeletedAt     string
	CreatedBy     string
	UpdatedBy     string
	Status        string
	Sort          string
}

// codeGenTablesColumns holds the columns for table code_gen_tables.
var codeGenTablesColumns = CodeGenTablesColumns{
	Id:            "id",
	TableName:     "table_name",
	TableComment:  "table_comment",
	Remark:        "remark",
	ModuleName:    "module_name",
	BelongMenuId:  "belong_menu_id",
	Type:          "type",
	MenuName:      "menu_name",
	ComponentType: "component_type",
	TplType:       "tpl_type",
	TreeId:        "tree_id",
	TreeParentId:  "tree_parent_id",
	TreeName:      "tree_name",
	TagId:         "tag_id",
	TagName:       "tag_name",
	TagViewName:   "tag_view_name",
	GenerateMenus: "generate_menus",
	Options:       "options",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
	DeletedAt:     "deleted_at",
	CreatedBy:     "created_by",
	UpdatedBy:     "updated_by",
	Status:        "status",
	Sort:          "sort",
}

// NewCodeGenTablesDao creates and returns a new DAO object for table data access.
func NewCodeGenTablesDao(handlers ...gdb.ModelHandler) *CodeGenTablesDao {
	return &CodeGenTablesDao{
		group:    "default",
		table:    "code_gen_tables",
		columns:  codeGenTablesColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CodeGenTablesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CodeGenTablesDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CodeGenTablesDao) Columns() CodeGenTablesColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CodeGenTablesDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CodeGenTablesDao) Ctx(ctx context.Context) *gdb.Model {
	return g.DB(dao.group).Model(dao.table).Safe().Ctx(ctx).Hook(dao.handlers...)
}

// Columns returns all column names of the current DAO for internal usage.
func (dao *CodeGenTablesDao) Columns() CodeGenTablesColumns {
	return dao.columns
}
