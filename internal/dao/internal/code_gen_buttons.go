// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CodeGenButtonsDao is the data access object for table code_gen_buttons.
type CodeGenButtonsDao struct {
	table    string                    // table is the underlying table name of the DAO.
	group    string                    // group is the database configuration group name of the current DAO.
	columns  CodeGenButtonsColumns    // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler       // handlers for customized model modification.
}

// CodeGenButtonsColumns defines and stores column names for table code_gen_buttons.
type CodeGenButtonsColumns struct {
	Id           string
	TableId      string
	ButtonCode   string
	ButtonName   string
	ButtonComment string
	IsShow       string
	Sort         string
	CreatedAt    string
	UpdatedAt    string
}

// codeGenButtonsColumns holds the columns for table code_gen_buttons.
var codeGenButtonsColumns = CodeGenButtonsColumns{
	Id:           "id",
	TableId:      "table_id",
	ButtonCode:   "button_code",
	ButtonName:   "button_name",
	ButtonComment: "button_comment",
	IsShow:       "is_show",
	Sort:         "sort",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
}

// NewCodeGenButtonsDao creates and returns a new DAO object for table data access.
func NewCodeGenButtonsDao(handlers ...gdb.ModelHandler) *CodeGenButtonsDao {
	return &CodeGenButtonsDao{
		group:    "default",
		table:    "code_gen_buttons",
		columns:  codeGenButtonsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CodeGenButtonsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CodeGenButtonsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CodeGenButtonsDao) Columns() CodeGenButtonsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CodeGenButtonsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CodeGenButtonsDao) Ctx(ctx context.Context) *gdb.Model {
	return g.DB(dao.group).Model(dao.table).Safe().Ctx(ctx).Hook(dao.handlers...)
}
