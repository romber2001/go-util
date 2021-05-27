package parser

import (
	"encoding/json"

	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
)

type Result struct {
	SQLType        string            `json:"sql_type"`
	DBNames        []string          `json:"db_names"`
	TableNames     []string          `json:"table_names"`
	TableComments  map[string]string `json:"table_comments"`
	ColumnNames    []string          `json:"column_names"`
	ColumnTypes    map[string]string `json:"column_types"`
	ColumnComments map[string]string `json:"column_comments"`
}

// NewResult returns a new *Result
func NewResult(sqlType string, dbNames []string, tableNames []string, tableComments map[string]string,
	columnNames []string, columnTypes map[string]string, columnComments map[string]string) *Result {
	return &Result{
		SQLType:        sqlType,
		DBNames:        dbNames,
		TableNames:     tableNames,
		TableComments:  tableComments,
		ColumnNames:    columnNames,
		ColumnTypes:    columnTypes,
		ColumnComments: columnComments,
	}
}

// NewEmptyResult returns an empty *Result
func NewEmptyResult() *Result {
	return &Result{
		SQLType:        constant.EmptyString,
		DBNames:        []string{},
		TableNames:     []string{},
		TableComments:  make(map[string]string),
		ColumnNames:    []string{},
		ColumnTypes:    make(map[string]string),
		ColumnComments: make(map[string]string),
	}
}

// GetSQLType returns the sql type
func (r *Result) GetSQLType() string {
	return r.SQLType
}

// GetDBNames returns the db names
func (r *Result) GetDBNames() []string {
	return r.DBNames
}

// GetTableNames returns the table names
func (r *Result) GetTableNames() []string {
	return r.TableNames
}

// GetTableComments returns the table comments
func (r *Result) GetTableComments() map[string]string {
	return r.TableComments
}

// GetColumnNames returns the column names
func (r *Result) GetColumnNames() []string {
	return r.ColumnNames
}

// GetColumnTypes returns the column types
func (r *Result) GetColumnTypes() map[string]string {
	return r.ColumnTypes
}

// GetColumnComments returns the column comments
func (r *Result) GetColumnComments() map[string]string {
	return r.ColumnComments
}

// SetSQLType sets the sql type
func (r *Result) SetSQLType(sqlType string) {
	r.SQLType = sqlType
}

// AddDBName adds db name to the result
func (r *Result) AddDBName(dbName string) {
	if !common.StringInSlice(r.DBNames, dbName) {
		r.DBNames = append(r.DBNames, dbName)
	}
}

// AddTableName adds table name to the result
func (r *Result) AddTableName(tableName string) {
	if !common.StringInSlice(r.TableNames, tableName) {
		r.TableNames = append(r.TableNames, tableName)
	}
}

// SetTableComment sets table comment of corresponding table
func (r *Result) SetTableComment(tableName string, tableComment string) {
	r.TableComments[tableName] = tableComment
}

// AddColumn adds column name to the result
func (r *Result) AddColumn(columnName string) {
	if !common.StringInSlice(r.ColumnNames, columnName) {
		r.ColumnNames = append(r.ColumnNames, columnName)
	}
}

// SetColumnType sets column type of corresponding column
func (r *Result) SetColumnType(columnName string, columnType string) {
	r.ColumnTypes[columnName] = columnType
}

// SetColumnComment sets column comment of corresponding column
func (r *Result) SetColumnComment(columnName string, columnComment string) {
	r.ColumnComments[columnName] = columnComment
}

// Marshal marshals result to json bytes
func (r *Result) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
