// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CodeGenFieldsDao is the data access object for table code_gen_fields.
type CodeGenFieldsDao struct {
	table    string                   // table is the underlying table name of the DAO.
	group    string                   // group is the database configuration group name of the current DAO.
	columns  CodeGenFieldsColumns    // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// CodeGenFieldsColumns defines and stores column names for table code_gen_fields.
type CodeGenFieldsColumns struct {
	Id           string
	TableId      string
	ColumnName   string
	ColumnComment string
	ColumnType   string
	DataType     string
	IsNullable   string
	Sort         string
	IsRequired   string
	IsInsert     string
	IsEdit       string
	IsList       string
	IsQuery      string
	IsSort       string
	QueryType    string
	ViewType     string
	DictType     string
	AllowRoles   string
	Options      string
	CreatedAt    string
	UpdatedAt    string
}

// codeGenFieldsColumns holds the columns for table code_gen_fields.
var codeGenFieldsColumns = CodeGenFieldsColumns{
	Id:           "id",
	TableId:      "table_id",
	ColumnName:   "column_name",
	ColumnComment: "column_comment",
	ColumnType:   "column_type",
	DataType:     "data_type",
	IsNullable:   "is_nullable",
	Sort:         "sort",
	IsRequired:   "is_required",
	IsInsert:     "is_insert",
	IsEdit:       "is_edit",
	IsList:       "is_list",
	IsQuery:      "is_query",
	IsSort:       "is_sort",
	QueryType:    "query_type",
	ViewType:     "view_type",
	DictType:     "dict_type",
	AllowRoles:   "allow_roles",
	Options:      "options",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
}

// NewCodeGenFieldsDao creates and returns a new DAO object for table data access.
func NewCodeGenFieldsDao(handlers ...gdb.ModelHandler) *CodeGenFieldsDao {
	return &CodeGenFieldsDao{
		group:    "default",
		table:    "code_gen_fields",
		columns:  codeGenFieldsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CodeGenFieldsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CodeGenFieldsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CodeGenFieldsDao) Columns() CodeGenFieldsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CodeGenFieldsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CodeGenFieldsDao) Ctx(ctx context.Context) *gdb.Model {
	return g.DB(dao.group).Model(dao.table).Safe().Ctx(ctx).Hook(dao.handlers...)
}
