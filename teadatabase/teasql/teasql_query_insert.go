package teasql

import (
	"context"
	"database/sql"

	"github.com/tea-frame-go/tea/internal/bufferpool"
)

func (db *Db) Insert() *insertHandler {
	placeholderGenerator := db.sqlBuilder.GetPlaceholderGenerator()
	insertHandler := &insertHandler{
		db:                   db,
		placeholderGenerator: placeholderGenerator,
		table:                "",
	}
	return insertHandler
}

type insertHandler struct {
	db                   *Db
	placeholderGenerator PlaceholderGenerator
	table                TableName
}

func (handler *insertHandler) Table(table TableName) *insertHandler {
	handler.table = table
	return handler
}

// todo))
func (handler *insertHandler) Insert(ctx context.Context, fields []FieldName, values [][]any) (sql.Result, error) {
	buf := bufferpool.NewBuffer(bufPool)
	args := handler.query(buf, fields, values)
	res, err := handler.db.Exec(ctx, buf.StringUnsafe(), args...)
	buf.Free(bufPool)
	return res, err
}

func (handler *insertHandler) query(buf *bufferpool.Buffer, fields []FieldName, values [][]any) []any {
	length := len(values)
	if length > 0 {
		length = length * len(values[0])
	}
	args := make([]any, length)
	index := 0
	buf.WriteString("INSERT INTO ")
	buf.WriteString(handler.table.Quote(handler.db.sqlBuilder))
	if len(fields) > 0 {
		buf.WriteString(" (")
		buf.WriteString(fields[0].Quote(handler.db.sqlBuilder))
		for i := 1; i < len(fields); i++ {
			buf.WriteByte(',')
			buf.WriteString(fields[i].Quote(handler.db.sqlBuilder))
		}
		buf.WriteByte(')')
	}
	buf.WriteString(" VALUES (")
	if length > 0 {
		if v, ok := values[0][0].(Raw); ok {
			buf.WriteString(v.Quote(handler.db.sqlBuilder))
		} else {
			buf.WriteString(handler.placeholderGenerator.GeneratePlaceholder())
			args[index] = values[0][0]
			index++
		}
		for j := 1; j < len(values[0]); j++ {
			buf.WriteByte(',')
			if v, ok := values[0][j].(Raw); ok {
				buf.WriteString(v.Quote(handler.db.sqlBuilder))
			} else {
				buf.WriteString(handler.placeholderGenerator.GeneratePlaceholder())
				args[index] = values[0][j]
				index++
			}
		}
		for i := 1; i < len(values); i++ {
			buf.WriteString("),(")
			if v, ok := values[i][0].(Raw); ok {
				buf.WriteString(v.Quote(handler.db.sqlBuilder))
			} else {
				buf.WriteString(handler.placeholderGenerator.GeneratePlaceholder())
				args[index] = values[i][0]
				index++
			}
			for j := 1; j < len(values[i]); j++ {
				buf.WriteByte(',')
				if v, ok := values[i][j].(Raw); ok {
					buf.WriteString(v.Quote(handler.db.sqlBuilder))
				} else {
					buf.WriteString(handler.placeholderGenerator.GeneratePlaceholder())
					args[index] = values[i][j]
					index++
				}
			}
		}
	}
	buf.WriteByte(')')
	args = args[:index]
	return args
}
