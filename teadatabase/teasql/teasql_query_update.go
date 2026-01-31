package teasql

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
)

func (db *DB) Update() *updateHandler {
	placeholderGenerator := db.sqlBuilder.GetPlaceholderGenerator()
	updateHandler := &updateHandler{
		db:                   db,
		placeholderGenerator: placeholderGenerator,
		table:                "",
		where: whereBuilder{
			WhereBuilder: WhereBuilderAnd{
				conditions: []whereCondition{},
				root:       nil,
			},
			args:                 []any{},
			placeholderGenerator: placeholderGenerator,
		},
		orderBy: []string{},
		limit:   0,
	}
	updateHandler.where.WhereBuilder.root = &updateHandler.where
	return updateHandler
}

type updateHandler struct {
	db                   *DB
	placeholderGenerator PlaceholderGenerator
	table                TableName
	where                whereBuilder
	orderBy              []string
	limit                int
}

func (handler *updateHandler) Table(table TableName) *updateHandler {
	handler.table = table
	return handler
}

func (handler *updateHandler) WhereBuilder(f func(w *WhereBuilderAnd)) *updateHandler {
	f(&handler.where.WhereBuilder)
	return handler
}

func (handler *updateHandler) OrderBy(orderBy ...string) *updateHandler {
	handler.orderBy = append(handler.orderBy, orderBy...)
	return handler
}

func (handler *updateHandler) OrderByAsc(field FieldName) *updateHandler {
	handler.orderBy = append(handler.orderBy, field.Quote(handler.db.sqlBuilder)+" ASC")
	return handler
}

func (handler *updateHandler) OrderByDesc(field FieldName) *updateHandler {
	handler.orderBy = append(handler.orderBy, field.Quote(handler.db.sqlBuilder)+" DESC")
	return handler
}

func (handler *updateHandler) Limit(limit int) *updateHandler {
	handler.limit = limit
	return handler
}

// todo))
func (handler *updateHandler) Update(ctx context.Context, fields []FieldName, values []any) (sql.Result, error) {
	var sb strings.Builder
	args := handler.query(&sb, fields, values)
	res, err := handler.db.Exec(ctx, sb.String(), args...)
	return res, err
}

func (handler *updateHandler) query(sb *strings.Builder, fields []FieldName, values []any) []any {
	args := make([]any, len(fields)+len(handler.where.args))
	index := 0
	sb.WriteString("UPDATE ")
	sb.WriteString(handler.table.Quote(handler.db.sqlBuilder))
	sb.WriteString(" SET ")
	if len(fields) > 0 {
		sb.WriteString(fields[0].Quote(handler.db.sqlBuilder))
		sb.WriteByte('=')
		if v, ok := values[0].(Raw); ok {
			sb.WriteString(v.Quote(handler.db.sqlBuilder))
		} else {
			sb.WriteString(handler.placeholderGenerator.GeneratePlaceholder())
			args[index] = values[0]
			index++
		}
		for i := 1; i < len(fields); i++ {
			sb.WriteByte(',')
			sb.WriteString(fields[i].Quote(handler.db.sqlBuilder))
			sb.WriteByte('=')
			if v, ok := values[i].(Raw); ok {
				sb.WriteString(v.Quote(handler.db.sqlBuilder))
			} else {
				sb.WriteString(handler.placeholderGenerator.GeneratePlaceholder())
				args[index] = values[i]
				index++
			}
		}
	}
	if len(handler.where.WhereBuilder.conditions) > 0 {
		sb.WriteString(" WHERE ")
		handler.where.WhereBuilder.conditions[0].buildHead(sb, handler.db.sqlBuilder)
		for i := 1; i < len(handler.where.WhereBuilder.conditions); i++ {
			handler.where.WhereBuilder.conditions[i].buildBody(sb, handler.db.sqlBuilder)
		}
		copy(args[index:], handler.where.args)
		index += len(handler.where.args)
	}
	if len(handler.orderBy) > 0 {
		sb.WriteString(" ORDER BY ")
		sb.WriteString(handler.orderBy[0])
		for i := 1; i < len(handler.orderBy); i++ {
			sb.WriteByte(',')
			sb.WriteString(handler.orderBy[i])
		}
	}
	if handler.limit > 0 {
		sb.WriteString(" LIMIT ")
		sb.WriteString(strconv.Itoa(handler.limit))
	}
	args = args[:index]
	return args
}
