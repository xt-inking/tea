package teasql

import (
	"github.com/tea-frame-go/tea/internal/bufferpool"
)

type whereBuilder struct {
	WhereBuilder         WhereBuilderAnd
	args                 []any
	placeholderGenerator PlaceholderGenerator
}

type whereCondition interface {
	buildHead(buf *bufferpool.Buffer, sqlBuilder sqlBuilder)
	buildBody(buf *bufferpool.Buffer, sqlBuilder sqlBuilder)
	logicalOperator(buf *bufferpool.Buffer)
}

type WhereBuilderAnd = WhereBuilder[logicalOperatorAnd]

type WhereBuilderOr = WhereBuilder[logicalOperatorOr]

type WhereBuilder[LogicalOperator logicalOperator] struct {
	conditions []whereCondition
	root       *whereBuilder
}

func (w WhereBuilder[LogicalOperator]) build(buf *bufferpool.Buffer, sqlBuilder sqlBuilder) {
	buf.WriteByte('(')
	w.conditions[0].buildHead(buf, sqlBuilder)
	for i := 1; i < len(w.conditions); i++ {
		w.conditions[i].buildBody(buf, sqlBuilder)
	}
	buf.WriteByte(')')
}

func (w WhereBuilder[LogicalOperator]) buildHead(buf *bufferpool.Buffer, sqlBuilder sqlBuilder) {
	if len(w.conditions) == 0 {
		return
	}
	w.build(buf, sqlBuilder)
}

func (w WhereBuilder[LogicalOperator]) buildBody(buf *bufferpool.Buffer, sqlBuilder sqlBuilder) {
	if len(w.conditions) == 0 {
		return
	}
	w.logicalOperator(buf)
	w.build(buf, sqlBuilder)
}

func (w WhereBuilder[LogicalOperator]) logicalOperator(buf *bufferpool.Buffer) {
	LogicalOperator{}.logicalOperator(buf)
}

func (w *WhereBuilder[LogicalOperator]) WhereBuilder(f func(w *WhereBuilderAnd)) *WhereBuilder[LogicalOperator] {
	whereBuilder := WhereBuilderAnd{
		conditions: []whereCondition{},
		root:       w.root,
	}
	f(&whereBuilder)
	w.conditions = append(w.conditions, whereBuilder)
	return w
}

func (w *WhereBuilder[LogicalOperator]) Where(field identifier, comparisonOperator string, arg any) *WhereBuilder[LogicalOperator] {
	whereConditionStruct := whereConditionStructAnd{
		field:              field,
		comparisonOperator: " " + comparisonOperator + " " + w.root.placeholderGenerator.GeneratePlaceholder(),
	}
	w.conditions = append(w.conditions, whereConditionStruct)
	w.root.args = append(w.root.args, arg)
	return w
}

func (w *WhereBuilder[LogicalOperator]) WhereBetween(field identifier, args [2]any) *WhereBuilder[LogicalOperator] {
	whereConditionStruct := whereConditionStructAnd{
		field:              field,
		comparisonOperator: " BETWEEN " + w.root.placeholderGenerator.GeneratePlaceholder() + " AND " + w.root.placeholderGenerator.GeneratePlaceholder(),
	}
	w.conditions = append(w.conditions, whereConditionStruct)
	w.root.args = append(w.root.args, args[0], args[1])
	return w
}

func (w *WhereBuilder[LogicalOperator]) WhereNotBetween(field identifier, args [2]any) *WhereBuilder[LogicalOperator] {
	whereConditionStruct := whereConditionStructAnd{
		field:              field,
		comparisonOperator: " NOT BETWEEN " + w.root.placeholderGenerator.GeneratePlaceholder() + " AND " + w.root.placeholderGenerator.GeneratePlaceholder(),
	}
	w.conditions = append(w.conditions, whereConditionStruct)
	w.root.args = append(w.root.args, args[0], args[1])
	return w
}

func (w *WhereBuilder[LogicalOperator]) WhereIn(field identifier, args []any) *WhereBuilder[LogicalOperator] {
	if len(args) == 0 {
		whereConditionStruct := whereConditionStructAnd{
			field:              Raw(""),
			comparisonOperator: "FALSE",
		}
		w.conditions = append(w.conditions, whereConditionStruct)
		return w
	}
	buf := bufferpool.NewBuffer(bufPool)
	defer buf.Free(bufPool)
	buf.WriteString(" IN (" + w.root.placeholderGenerator.GeneratePlaceholder())
	for i := 1; i < len(args); i++ {
		buf.WriteString("," + w.root.placeholderGenerator.GeneratePlaceholder())
	}
	buf.WriteByte(')')
	whereConditionStruct := whereConditionStructAnd{
		field:              field,
		comparisonOperator: buf.String(),
	}
	w.conditions = append(w.conditions, whereConditionStruct)
	w.root.args = append(w.root.args, args...)
	return w
}

func (w *WhereBuilder[LogicalOperator]) WhereNotIn(field identifier, args []any) *WhereBuilder[LogicalOperator] {
	if len(args) == 0 {
		whereConditionStruct := whereConditionStructAnd{
			field:              Raw(""),
			comparisonOperator: "TRUE",
		}
		w.conditions = append(w.conditions, whereConditionStruct)
		return w
	}
	buf := bufferpool.NewBuffer(bufPool)
	defer buf.Free(bufPool)
	buf.WriteString(" NOT IN (" + w.root.placeholderGenerator.GeneratePlaceholder())
	for i := 1; i < len(args); i++ {
		buf.WriteString("," + w.root.placeholderGenerator.GeneratePlaceholder())
	}
	buf.WriteByte(')')
	whereConditionStruct := whereConditionStructAnd{
		field:              field,
		comparisonOperator: buf.String(),
	}
	w.conditions = append(w.conditions, whereConditionStruct)
	w.root.args = append(w.root.args, args...)
	return w
}

func (w *WhereBuilder[LogicalOperator]) WhereLike(field identifier, arg string) *WhereBuilder[LogicalOperator] {
	whereConditionStruct := whereConditionStructAnd{
		field:              field,
		comparisonOperator: " LIKE " + w.root.placeholderGenerator.GeneratePlaceholder(),
	}
	w.conditions = append(w.conditions, whereConditionStruct)
	w.root.args = append(w.root.args, arg)
	return w
}

func (w *WhereBuilder[LogicalOperator]) WhereNotLike(field identifier, arg string) *WhereBuilder[LogicalOperator] {
	whereConditionStruct := whereConditionStructAnd{
		field:              field,
		comparisonOperator: " NOT LIKE " + w.root.placeholderGenerator.GeneratePlaceholder(),
	}
	w.conditions = append(w.conditions, whereConditionStruct)
	w.root.args = append(w.root.args, arg)
	return w
}

func (w *WhereBuilder[LogicalOperator]) WhereNull(field identifier) *WhereBuilder[LogicalOperator] {
	whereConditionStruct := whereConditionStructAnd{
		field:              field,
		comparisonOperator: " IS NULL",
	}
	w.conditions = append(w.conditions, whereConditionStruct)
	return w
}

func (w *WhereBuilder[LogicalOperator]) WhereNotNull(field identifier) *WhereBuilder[LogicalOperator] {
	whereConditionStruct := whereConditionStructAnd{
		field:              field,
		comparisonOperator: " IS NOT NULL",
	}
	w.conditions = append(w.conditions, whereConditionStruct)
	return w
}

func (w *WhereBuilder[LogicalOperator]) Or() *orWhereBuilder[LogicalOperator] {
	return &orWhereBuilder[LogicalOperator]{
		w: w,
	}
}

type orWhereBuilder[LogicalOperator logicalOperator] struct {
	w *WhereBuilder[LogicalOperator]
}

func (w *orWhereBuilder[LogicalOperator]) WhereBuilder(f func(w *WhereBuilderOr)) *WhereBuilder[LogicalOperator] {
	whereBuilderOr := WhereBuilderOr{
		conditions: []whereCondition{},
		root:       w.w.root,
	}
	f(&whereBuilderOr)
	w.w.conditions = append(w.w.conditions, whereBuilderOr)
	return w.w
}

func (w *orWhereBuilder[LogicalOperator]) Where(field identifier, comparisonOperator string, arg any) *WhereBuilder[LogicalOperator] {
	whereConditionStructOr := whereConditionStructOr{
		field:              field,
		comparisonOperator: " " + comparisonOperator + " " + w.w.root.placeholderGenerator.GeneratePlaceholder(),
	}
	w.w.conditions = append(w.w.conditions, whereConditionStructOr)
	w.w.root.args = append(w.w.root.args, arg)
	return w.w
}

func (w *orWhereBuilder[LogicalOperator]) WhereBetween(field identifier, args [2]any) *WhereBuilder[LogicalOperator] {
	whereConditionStructOr := whereConditionStructOr{
		field:              field,
		comparisonOperator: " BETWEEN " + w.w.root.placeholderGenerator.GeneratePlaceholder() + " AND " + w.w.root.placeholderGenerator.GeneratePlaceholder(),
	}
	w.w.conditions = append(w.w.conditions, whereConditionStructOr)
	w.w.root.args = append(w.w.root.args, args[0], args[1])
	return w.w
}

func (w *orWhereBuilder[LogicalOperator]) WhereNotBetween(field identifier, args [2]any) *WhereBuilder[LogicalOperator] {
	whereConditionStructOr := whereConditionStructOr{
		field:              field,
		comparisonOperator: " NOT BETWEEN " + w.w.root.placeholderGenerator.GeneratePlaceholder() + " AND " + w.w.root.placeholderGenerator.GeneratePlaceholder(),
	}
	w.w.conditions = append(w.w.conditions, whereConditionStructOr)
	w.w.root.args = append(w.w.root.args, args[0], args[1])
	return w.w
}

func (w *orWhereBuilder[LogicalOperator]) WhereIn(field identifier, args []any) *WhereBuilder[LogicalOperator] {
	if len(args) == 0 {
		whereConditionStructOr := whereConditionStructOr{
			field:              Raw(""),
			comparisonOperator: "FALSE",
		}
		w.w.conditions = append(w.w.conditions, whereConditionStructOr)
		return w.w
	}
	buf := bufferpool.NewBuffer(bufPool)
	defer buf.Free(bufPool)
	buf.WriteString(" IN (" + w.w.root.placeholderGenerator.GeneratePlaceholder())
	for i := 1; i < len(args); i++ {
		buf.WriteString("," + w.w.root.placeholderGenerator.GeneratePlaceholder())
	}
	buf.WriteByte(')')
	whereConditionStructOr := whereConditionStructOr{
		field:              field,
		comparisonOperator: buf.String(),
	}
	w.w.conditions = append(w.w.conditions, whereConditionStructOr)
	w.w.root.args = append(w.w.root.args, args...)
	return w.w
}

func (w *orWhereBuilder[LogicalOperator]) WhereNotIn(field identifier, args []any) *WhereBuilder[LogicalOperator] {
	if len(args) == 0 {
		whereConditionStructOr := whereConditionStructOr{
			field:              Raw(""),
			comparisonOperator: "TRUE",
		}
		w.w.conditions = append(w.w.conditions, whereConditionStructOr)
		return w.w
	}
	buf := bufferpool.NewBuffer(bufPool)
	defer buf.Free(bufPool)
	buf.WriteString(" NOT IN (" + w.w.root.placeholderGenerator.GeneratePlaceholder())
	for i := 1; i < len(args); i++ {
		buf.WriteString("," + w.w.root.placeholderGenerator.GeneratePlaceholder())
	}
	buf.WriteByte(')')
	whereConditionStructOr := whereConditionStructOr{
		field:              field,
		comparisonOperator: buf.String(),
	}
	w.w.conditions = append(w.w.conditions, whereConditionStructOr)
	w.w.root.args = append(w.w.root.args, args...)
	return w.w
}

func (w *orWhereBuilder[LogicalOperator]) WhereLike(field identifier, arg string) *WhereBuilder[LogicalOperator] {
	whereConditionStructOr := whereConditionStructOr{
		field:              field,
		comparisonOperator: " LIKE " + w.w.root.placeholderGenerator.GeneratePlaceholder(),
	}
	w.w.conditions = append(w.w.conditions, whereConditionStructOr)
	w.w.root.args = append(w.w.root.args, arg)
	return w.w
}

func (w *orWhereBuilder[LogicalOperator]) WhereNotLike(field identifier, arg string) *WhereBuilder[LogicalOperator] {
	whereConditionStructOr := whereConditionStructOr{
		field:              field,
		comparisonOperator: " NOT LIKE " + w.w.root.placeholderGenerator.GeneratePlaceholder(),
	}
	w.w.conditions = append(w.w.conditions, whereConditionStructOr)
	w.w.root.args = append(w.w.root.args, arg)
	return w.w
}

func (w *orWhereBuilder[LogicalOperator]) WhereNull(field identifier) *WhereBuilder[LogicalOperator] {
	whereConditionStructOr := whereConditionStructOr{
		field:              field,
		comparisonOperator: " IS NULL",
	}
	w.w.conditions = append(w.w.conditions, whereConditionStructOr)
	return w.w
}

func (w *orWhereBuilder[LogicalOperator]) WhereNotNull(field identifier) *WhereBuilder[LogicalOperator] {
	whereConditionStructOr := whereConditionStructOr{
		field:              field,
		comparisonOperator: " IS NOT NULL",
	}
	w.w.conditions = append(w.w.conditions, whereConditionStructOr)
	return w.w
}

type whereConditionStructAnd = whereConditionStruct[logicalOperatorAnd]

type whereConditionStructOr = whereConditionStruct[logicalOperatorOr]

type whereConditionStruct[LogicalOperator logicalOperator] struct {
	field              identifier
	comparisonOperator string
}

func (w whereConditionStruct[LogicalOperator]) build(buf *bufferpool.Buffer, sqlBuilder sqlBuilder) {
	buf.WriteString(w.field.Quote(sqlBuilder))
	buf.WriteString(w.comparisonOperator)
}

func (w whereConditionStruct[LogicalOperator]) buildHead(buf *bufferpool.Buffer, sqlBuilder sqlBuilder) {
	w.build(buf, sqlBuilder)
}

func (w whereConditionStruct[LogicalOperator]) buildBody(buf *bufferpool.Buffer, sqlBuilder sqlBuilder) {
	w.logicalOperator(buf)
	w.build(buf, sqlBuilder)
}

func (w whereConditionStruct[LogicalOperator]) logicalOperator(buf *bufferpool.Buffer) {
	LogicalOperator{}.logicalOperator(buf)
}

type logicalOperator interface {
	~struct{}
	logicalOperator(buf *bufferpool.Buffer)
}

type logicalOperatorAnd struct{}

func (logicalOperatorAnd) logicalOperator(buf *bufferpool.Buffer) {
	buf.WriteString(" AND ")
}

type logicalOperatorOr struct{}

func (logicalOperatorOr) logicalOperator(buf *bufferpool.Buffer) {
	buf.WriteString(" OR ")
}
