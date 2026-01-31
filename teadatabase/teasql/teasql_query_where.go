package teasql

import (
	"strings"
)

type whereBuilder struct {
	WhereBuilder         WhereBuilderAnd
	args                 []any
	placeholderGenerator PlaceholderGenerator
}

type whereCondition interface {
	buildHead(sb *strings.Builder, sqlBuilder sqlBuilder)
	buildBody(sb *strings.Builder, sqlBuilder sqlBuilder)
	logicalOperator(sb *strings.Builder)
}

type WhereBuilderAnd = WhereBuilder[logicalOperatorAnd]

type WhereBuilderOr = WhereBuilder[logicalOperatorOr]

type WhereBuilder[LogicalOperator logicalOperator] struct {
	conditions []whereCondition
	root       *whereBuilder
}

func (w WhereBuilder[LogicalOperator]) build(sb *strings.Builder, sqlBuilder sqlBuilder) {
	sb.WriteByte('(')
	w.conditions[0].buildHead(sb, sqlBuilder)
	for i := 1; i < len(w.conditions); i++ {
		w.conditions[i].buildBody(sb, sqlBuilder)
	}
	sb.WriteByte(')')
}

func (w WhereBuilder[LogicalOperator]) buildHead(sb *strings.Builder, sqlBuilder sqlBuilder) {
	if len(w.conditions) == 0 {
		return
	}
	w.build(sb, sqlBuilder)
}

func (w WhereBuilder[LogicalOperator]) buildBody(sb *strings.Builder, sqlBuilder sqlBuilder) {
	if len(w.conditions) == 0 {
		return
	}
	w.logicalOperator(sb)
	w.build(sb, sqlBuilder)
}

func (w WhereBuilder[LogicalOperator]) logicalOperator(sb *strings.Builder) {
	LogicalOperator{}.logicalOperator(sb)
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
	var sb strings.Builder
	sb.WriteString(" IN (" + w.root.placeholderGenerator.GeneratePlaceholder())
	for i := 1; i < len(args); i++ {
		sb.WriteString("," + w.root.placeholderGenerator.GeneratePlaceholder())
	}
	sb.WriteByte(')')
	whereConditionStruct := whereConditionStructAnd{
		field:              field,
		comparisonOperator: sb.String(),
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
	var sb strings.Builder
	sb.WriteString(" NOT IN (" + w.root.placeholderGenerator.GeneratePlaceholder())
	for i := 1; i < len(args); i++ {
		sb.WriteString("," + w.root.placeholderGenerator.GeneratePlaceholder())
	}
	sb.WriteByte(')')
	whereConditionStruct := whereConditionStructAnd{
		field:              field,
		comparisonOperator: sb.String(),
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
	var sb strings.Builder
	sb.WriteString(" IN (" + w.w.root.placeholderGenerator.GeneratePlaceholder())
	for i := 1; i < len(args); i++ {
		sb.WriteString("," + w.w.root.placeholderGenerator.GeneratePlaceholder())
	}
	sb.WriteByte(')')
	whereConditionStructOr := whereConditionStructOr{
		field:              field,
		comparisonOperator: sb.String(),
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
	var sb strings.Builder
	sb.WriteString(" NOT IN (" + w.w.root.placeholderGenerator.GeneratePlaceholder())
	for i := 1; i < len(args); i++ {
		sb.WriteString("," + w.w.root.placeholderGenerator.GeneratePlaceholder())
	}
	sb.WriteByte(')')
	whereConditionStructOr := whereConditionStructOr{
		field:              field,
		comparisonOperator: sb.String(),
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

func (w whereConditionStruct[LogicalOperator]) build(sb *strings.Builder, sqlBuilder sqlBuilder) {
	sb.WriteString(w.field.Quote(sqlBuilder))
	sb.WriteString(w.comparisonOperator)
}

func (w whereConditionStruct[LogicalOperator]) buildHead(sb *strings.Builder, sqlBuilder sqlBuilder) {
	w.build(sb, sqlBuilder)
}

func (w whereConditionStruct[LogicalOperator]) buildBody(sb *strings.Builder, sqlBuilder sqlBuilder) {
	w.logicalOperator(sb)
	w.build(sb, sqlBuilder)
}

func (w whereConditionStruct[LogicalOperator]) logicalOperator(sb *strings.Builder) {
	LogicalOperator{}.logicalOperator(sb)
}

type logicalOperator interface {
	~struct{}
	logicalOperator(sb *strings.Builder)
}

type logicalOperatorAnd struct{}

func (logicalOperatorAnd) logicalOperator(sb *strings.Builder) {
	sb.WriteString(" AND ")
}

type logicalOperatorOr struct{}

func (logicalOperatorOr) logicalOperator(sb *strings.Builder) {
	sb.WriteString(" OR ")
}
