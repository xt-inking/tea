package teasql

import (
	"context"
	"slices"
	"strconv"

	"github.com/tea-frame-go/tea/internal/bufferpool"
)

func (db *Db) Select() *selectHandler {
	placeholderGenerator := db.sqlBuilder.GetPlaceholderGenerator()
	selectHandler := &selectHandler{
		db:                   db,
		placeholderGenerator: placeholderGenerator,
		distinct:             false,
		fields:               []FieldName{},
		from:                 "",
		where: whereBuilder{
			WhereBuilder: WhereBuilderAnd{
				conditions: []whereCondition{},
				root:       nil,
			},
			args:                 []any{},
			placeholderGenerator: placeholderGenerator,
		},
		groupBy: []FieldName{},
		having: whereBuilder{
			WhereBuilder: WhereBuilderAnd{
				conditions: []whereCondition{},
				root:       nil,
			},
			args:                 []any{},
			placeholderGenerator: placeholderGenerator,
		},
		orderBy: []string{},
		limit:   0,
		offset:  0,
	}
	selectHandler.where.WhereBuilder.root = &selectHandler.where
	selectHandler.having.WhereBuilder.root = &selectHandler.having
	return selectHandler
}

type selectHandler struct {
	db                   *Db
	placeholderGenerator PlaceholderGenerator
	distinct             bool
	fields               []FieldName
	from                 TableName
	where                whereBuilder
	groupBy              []FieldName
	having               whereBuilder
	orderBy              []string
	limit                int
	offset               int
}

func (handler *selectHandler) Clone() *selectHandler {
	placeholderGenerator := handler.placeholderGenerator.Clone()
	newHandler := &selectHandler{
		db:                   handler.db,
		placeholderGenerator: placeholderGenerator,
		distinct:             handler.distinct,
		fields:               slices.Clone(handler.fields),
		from:                 handler.from,
		where: whereBuilder{
			WhereBuilder: WhereBuilderAnd{
				conditions: slices.Clone(handler.where.WhereBuilder.conditions),
				root:       nil,
			},
			args:                 slices.Clone(handler.where.args),
			placeholderGenerator: placeholderGenerator,
		},
		groupBy: slices.Clone(handler.groupBy),
		having: whereBuilder{
			WhereBuilder: WhereBuilderAnd{
				conditions: slices.Clone(handler.having.WhereBuilder.conditions),
				root:       nil,
			},
			args:                 slices.Clone(handler.having.args),
			placeholderGenerator: placeholderGenerator,
		},
		orderBy: slices.Clone(handler.orderBy),
		limit:   handler.limit,
		offset:  handler.offset,
	}
	newHandler.where.WhereBuilder.root = &newHandler.where
	newHandler.having.WhereBuilder.root = &newHandler.having
	return newHandler
}

func (handler *selectHandler) Distinct() *selectHandler {
	handler.distinct = true
	return handler
}

func (handler *selectHandler) Fields(fields ...FieldName) *selectHandler {
	handler.fields = fields
	return handler
}

func (handler *selectHandler) From(table TableName) *selectHandler {
	handler.from = table
	return handler
}

func (handler *selectHandler) WhereBuilder(f func(w *WhereBuilderAnd)) *selectHandler {
	f(&handler.where.WhereBuilder)
	return handler
}

func (handler *selectHandler) GroupBy(fields ...FieldName) *selectHandler {
	handler.groupBy = fields
	return handler
}

func (handler *selectHandler) Having(f func(w *WhereBuilderAnd)) *selectHandler {
	f(&handler.having.WhereBuilder)
	return handler
}

func (handler *selectHandler) OrderBy(orderBy ...string) *selectHandler {
	handler.orderBy = append(handler.orderBy, orderBy...)
	return handler
}

func (handler *selectHandler) OrderByAsc(field FieldName) *selectHandler {
	handler.orderBy = append(handler.orderBy, field.Quote(handler.db.sqlBuilder)+" ASC")
	return handler
}

func (handler *selectHandler) OrderByDesc(field FieldName) *selectHandler {
	handler.orderBy = append(handler.orderBy, field.Quote(handler.db.sqlBuilder)+" DESC")
	return handler
}

func (handler *selectHandler) Limit(limit int) *selectHandler {
	handler.limit = limit
	return handler
}

func (handler *selectHandler) Offset(offset int) *selectHandler {
	handler.offset = offset
	return handler
}

func (handler *selectHandler) Pagination(page int, size int) *selectHandler {
	handler.limit = size
	handler.offset = (page - 1) * size
	return handler
}

func (handler *selectHandler) List(ctx context.Context) {
	buf := bufferpool.NewBuffer(bufPool)
	args := handler.query(buf)
	rows, err := handler.db.Query(ctx, buf.StringUnsafe(), args...)
	buf.Free(bufPool)
	// todo))
	_, _ = rows, err
}

func (handler *selectHandler) One(ctx context.Context) {
	buf := bufferpool.NewBuffer(bufPool)
	args := handler.query(buf)
	rows, err := handler.db.Query(ctx, buf.StringUnsafe(), args...)
	buf.Free(bufPool)
	// todo))
	_, _ = rows, err
}

func (handler *selectHandler) query(buf *bufferpool.Buffer) []any {
	args := make([]any, len(handler.where.args)+len(handler.having.args))
	buf.WriteString("SELECT ")
	if handler.distinct {
		buf.WriteString("DISTINCT ")
	}
	if len(handler.fields) > 0 {
		buf.WriteString(handler.fields[0].Quote(handler.db.sqlBuilder))
		for i := 1; i < len(handler.fields); i++ {
			buf.WriteByte(',')
			buf.WriteString(handler.fields[i].Quote(handler.db.sqlBuilder))
		}
	} else {
		buf.WriteByte('*')
	}
	buf.WriteString(" FROM ")
	buf.WriteString(handler.from.Quote(handler.db.sqlBuilder))
	if len(handler.where.WhereBuilder.conditions) > 0 {
		buf.WriteString(" WHERE ")
		handler.where.WhereBuilder.conditions[0].buildHead(buf, handler.db.sqlBuilder)
		for i := 1; i < len(handler.where.WhereBuilder.conditions); i++ {
			handler.where.WhereBuilder.conditions[i].buildBody(buf, handler.db.sqlBuilder)
		}
		copy(args, handler.where.args)
	}
	if len(handler.groupBy) > 0 {
		buf.WriteString(" GROUP BY ")
		buf.WriteString(handler.groupBy[0].Quote(handler.db.sqlBuilder))
		for i := 1; i < len(handler.groupBy); i++ {
			buf.WriteByte(',')
			buf.WriteString(handler.groupBy[i].Quote(handler.db.sqlBuilder))
		}
	}
	if len(handler.having.WhereBuilder.conditions) > 0 {
		buf.WriteString(" HAVING ")
		handler.having.WhereBuilder.conditions[0].buildHead(buf, handler.db.sqlBuilder)
		for i := 1; i < len(handler.having.WhereBuilder.conditions); i++ {
			handler.having.WhereBuilder.conditions[i].buildBody(buf, handler.db.sqlBuilder)
		}
		copy(args[len(handler.where.args):], handler.having.args)
	}
	if len(handler.orderBy) > 0 {
		buf.WriteString(" ORDER BY ")
		buf.WriteString(handler.orderBy[0])
		for i := 1; i < len(handler.orderBy); i++ {
			buf.WriteByte(',')
			buf.WriteString(handler.orderBy[i])
		}
	}
	if handler.limit != 0 {
		buf.WriteString(" LIMIT ")
		if handler.offset != 0 {
			buf.WriteString(strconv.Itoa(handler.offset))
			buf.WriteByte(',')
		}
		buf.WriteString(strconv.Itoa(handler.limit))
	} else if handler.offset != 0 {
		buf.WriteString(" OFFSET ")
		buf.WriteString(strconv.Itoa(handler.offset))
	}
	return args
}
