package teasql

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
)

func (db *DB) Delete() *deleteHandler {
	placeholderGenerator := db.sqlBuilder.GetPlaceholderGenerator()
	deleteHandler := &deleteHandler{
		db:                   db,
		placeholderGenerator: placeholderGenerator,
		from:                 "",
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
	deleteHandler.where.WhereBuilder.root = &deleteHandler.where
	return deleteHandler
}

type deleteHandler struct {
	db                   *DB
	placeholderGenerator PlaceholderGenerator
	from                 TableName
	where                whereBuilder
	orderBy              []string
	limit                int
}

func (handler *deleteHandler) From(table TableName) *deleteHandler {
	handler.from = table
	return handler
}

func (handler *deleteHandler) WhereBuilder(f func(w *WhereBuilderAnd)) *deleteHandler {
	f(&handler.where.WhereBuilder)
	return handler
}

func (handler *deleteHandler) OrderBy(orderBy ...string) *deleteHandler {
	handler.orderBy = append(handler.orderBy, orderBy...)
	return handler
}

func (handler *deleteHandler) OrderByAsc(field FieldName) *deleteHandler {
	handler.orderBy = append(handler.orderBy, field.Quote(handler.db.sqlBuilder)+" ASC")
	return handler
}

func (handler *deleteHandler) OrderByDesc(field FieldName) *deleteHandler {
	handler.orderBy = append(handler.orderBy, field.Quote(handler.db.sqlBuilder)+" DESC")
	return handler
}

func (handler *deleteHandler) Limit(limit int) *deleteHandler {
	handler.limit = limit
	return handler
}

func (handler *deleteHandler) Delete(ctx context.Context) (sql.Result, error) {
	var sb strings.Builder
	args := handler.query(&sb)
	res, err := handler.db.Exec(ctx, sb.String(), args...)
	return res, err
}

func (handler *deleteHandler) query(sb *strings.Builder) []any {
	args := make([]any, len(handler.where.args))
	sb.WriteString("DELETE FROM ")
	sb.WriteString(handler.from.Quote(handler.db.sqlBuilder))
	if len(handler.where.WhereBuilder.conditions) > 0 {
		sb.WriteString(" WHERE ")
		handler.where.WhereBuilder.conditions[0].buildHead(sb, handler.db.sqlBuilder)
		for i := 1; i < len(handler.where.WhereBuilder.conditions); i++ {
			handler.where.WhereBuilder.conditions[i].buildBody(sb, handler.db.sqlBuilder)
		}
		copy(args, handler.where.args)
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
	return args
}
