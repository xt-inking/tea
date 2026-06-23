package teasql

import (
	"github.com/tea-frame-go/tea/internal/bufferpool"
)

type identifier interface {
	Quote(sqlBuilder sqlBuilder) string
}

type TableName string

func (tableName TableName) Quote(sqlBuilder sqlBuilder) string {
	return sqlBuilder.QuoteIdentifier(string(tableName))
}

type FieldName string

func (fieldName FieldName) Quote(sqlBuilder sqlBuilder) string {
	return sqlBuilder.QuoteIdentifier(string(fieldName))
}

type Raw string

func (raw Raw) Quote(sqlBuilder sqlBuilder) string {
	return string(raw)
}

var bufPool = bufferpool.New()
