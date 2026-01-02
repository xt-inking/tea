package teasql

import (
	"context"
	"database/sql"
	"strings"
)

func (db *DB) Insert() *insertHandler {
	placeholderGenerator := db.sqlBuilder.GetPlaceholderGenerator()
	insertHandler := &insertHandler{
		db:                   db,
		placeholderGenerator: placeholderGenerator,
		table:                "",
	}
	return insertHandler
}

type insertHandler struct {
	db                   *DB
	placeholderGenerator PlaceholderGenerator
	table                TableName
}

func (handler *insertHandler) Table(table TableName) *insertHandler {
	handler.table = table
	return handler
}

// todo))
func (handler *insertHandler) Insert(ctx context.Context, fields []FieldName, values [][]any) (sql.Result, error) {
	var sb strings.Builder
	args := handler.query(&sb, fields, values)
	res, err := handler.db.Exec(ctx, sb.String(), args...)
	return res, err
}

func (handler *insertHandler) query(sb *strings.Builder, fields []FieldName, values [][]any) []any {
	length := len(values)
	if length > 0 {
		length = length * len(values[0])
	}
	args := make([]any, length)
	index := 0
	sb.WriteString("INSERT INTO ")
	sb.WriteString(handler.table.Quote(handler.db.sqlBuilder))
	if len(fields) > 0 {
		sb.WriteString(" (")
		sb.WriteString(fields[0].Quote(handler.db.sqlBuilder))
		for i := 1; i < len(fields); i++ {
			sb.WriteByte(',')
			sb.WriteString(fields[i].Quote(handler.db.sqlBuilder))
		}
		sb.WriteByte(')')
	}
	sb.WriteString(" VALUES (")
	if length > 0 {
		if v, ok := values[0][0].(Raw); ok {
			sb.WriteString(v.Quote(handler.db.sqlBuilder))
		} else {
			sb.WriteString(handler.placeholderGenerator.GeneratePlaceholder())
			args[index] = values[0][0]
			index++
		}
		for j := 1; j < len(values[0]); j++ {
			sb.WriteByte(',')
			if v, ok := values[0][j].(Raw); ok {
				sb.WriteString(v.Quote(handler.db.sqlBuilder))
			} else {
				sb.WriteString(handler.placeholderGenerator.GeneratePlaceholder())
				args[index] = values[0][j]
				index++
			}
		}
		for i := 1; i < len(values); i++ {
			sb.WriteString("),(")
			if v, ok := values[i][0].(Raw); ok {
				sb.WriteString(v.Quote(handler.db.sqlBuilder))
			} else {
				sb.WriteString(handler.placeholderGenerator.GeneratePlaceholder())
				args[index] = values[i][0]
				index++
			}
			for j := 1; j < len(values[i]); j++ {
				sb.WriteByte(',')
				if v, ok := values[i][j].(Raw); ok {
					sb.WriteString(v.Quote(handler.db.sqlBuilder))
				} else {
					sb.WriteString(handler.placeholderGenerator.GeneratePlaceholder())
					args[index] = values[i][j]
					index++
				}
			}
		}
	}
	sb.WriteByte(')')
	args = args[:index]
	return args
}
