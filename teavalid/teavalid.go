package teavalid

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/tea-frame-go/tea/teatypes"
	"github.com/tea-frame-go/tea/teatypes/teaconstraints"
)

type Validator[T any] interface {
	*T
	Validate() (err Error)
}

func Integer[Value teaconstraints.Integer, Required required, Rule rule[Value]](
	value Value, k string, required Required, rule Rule,
) (err Error) {
	if value == 0 {
		err = required.Validate(teatypes.String(k))
		return
	}
	err = rule.Validate(teatypes.String(k), value)
	return
}

func Decimal[Value decimal.Decimal, Required required, Rule rule[Value]](
	value Value, k string, required Required, rule Rule,
) (err Error) {
	if decimal.Decimal(value).IsZero() {
		err = required.Validate(teatypes.String(k))
		return
	}
	err = rule.Validate(teatypes.String(k), value)
	return
}

func String[Value teaconstraints.String, Required required, Rule rule[Value]](
	value Value, k string, required Required, rule Rule,
) (err Error) {
	if value == "" {
		err = required.Validate(teatypes.String(k))
		return
	}
	err = rule.Validate(teatypes.String(k), value)
	return
}

func Time[Value time.Time, Required required, Rule rule[Value]](
	value Value, k string, required Required, rule Rule,
) (err Error) {
	if time.Time(value).IsZero() {
		err = required.Validate(teatypes.String(k))
		return
	}
	err = rule.Validate(teatypes.String(k), value)
	return
}

func Option[Option teatypes.Option[Value], Value any, Required required, Rule rule[Value]](
	option Option, k string, required Required, rule Rule,
) (err Error) {
	if teatypes.Option[Value](option).IsZero() {
		err = required.Validate(teatypes.String(k))
		return
	}
	err = rule.Validate(teatypes.String(k), teatypes.Option[Value](option).Value)
	return
}

func SliceInteger[Value ~[]Elem, Elem teaconstraints.Integer, Required required, Rule rule[Value], RequiredElem required, RuleElem rule[Elem]](
	value Value, k string, required Required, rule Rule, requiredElem RequiredElem, ruleElem RuleElem,
) (err Error) {
	if value == nil {
		err = required.Validate(teatypes.String(k))
		return
	}
	err = rule.Validate(teatypes.String(k), value)
	if err != nil {
		return
	}
	for i := range value {
		if value[i] == 0 {
			err = requiredElem.Validate(teatypes.Int(i))
			if err != nil {
				err.Key(k)
				return
			}
			continue
		}
		err = ruleElem.Validate(teatypes.Int(i), value[i])
		if err != nil {
			err.Key(k)
			return
		}
		continue
	}
	return
}

func SliceDecimal[Value ~[]Elem, Elem decimal.Decimal, Required required, Rule rule[Value], RequiredElem required, RuleElem rule[Elem]](
	value Value, k string, required Required, rule Rule, requiredElem RequiredElem, ruleElem RuleElem,
) (err Error) {
	if value == nil {
		err = required.Validate(teatypes.String(k))
		return
	}
	err = rule.Validate(teatypes.String(k), value)
	if err != nil {
		return
	}
	for i := range value {
		if decimal.Decimal(value[i]).IsZero() {
			err = requiredElem.Validate(teatypes.Int(i))
			if err != nil {
				err.Key(k)
				return
			}
			continue
		}
		err = ruleElem.Validate(teatypes.Int(i), value[i])
		if err != nil {
			err.Key(k)
			return
		}
		continue
	}
	return
}

func SliceString[Value ~[]Elem, Elem teaconstraints.String, Required required, Rule rule[Value], RequiredElem required, RuleElem rule[Elem]](
	value Value, k string, required Required, rule Rule, requiredElem RequiredElem, ruleElem RuleElem,
) (err Error) {
	if value == nil {
		err = required.Validate(teatypes.String(k))
		return
	}
	err = rule.Validate(teatypes.String(k), value)
	if err != nil {
		return
	}
	for i := range value {
		if value[i] == "" {
			err = requiredElem.Validate(teatypes.Int(i))
			if err != nil {
				err.Key(k)
				return
			}
			continue
		}
		err = ruleElem.Validate(teatypes.Int(i), value[i])
		if err != nil {
			err.Key(k)
			return
		}
		continue
	}
	return
}

func SliceTime[Value ~[]Elem, Elem time.Time, Required required, Rule rule[Value], RequiredElem required, RuleElem rule[Elem]](
	value Value, k string, required Required, rule Rule, requiredElem RequiredElem, ruleElem RuleElem,
) (err Error) {
	if value == nil {
		err = required.Validate(teatypes.String(k))
		return
	}
	err = rule.Validate(teatypes.String(k), value)
	if err != nil {
		return
	}
	for i := range value {
		if time.Time(value[i]).IsZero() {
			err = requiredElem.Validate(teatypes.Int(i))
			if err != nil {
				err.Key(k)
				return
			}
			continue
		}
		err = ruleElem.Validate(teatypes.Int(i), value[i])
		if err != nil {
			err.Key(k)
			return
		}
		continue
	}
	return
}

func SliceOption[Options ~[]Elem, Elem teatypes.Option[Value], Value any, Required required, Rule rule[Options], RequiredElem required, RuleElem rule[Value]](
	options Options, k string, required Required, rule Rule, requiredElem RequiredElem, ruleElem RuleElem,
) (err Error) {
	if options == nil {
		err = required.Validate(teatypes.String(k))
		return
	}
	err = rule.Validate(teatypes.String(k), options)
	if err != nil {
		return
	}
	for i := range options {
		if teatypes.Option[Value](options[i]).IsZero() {
			err = requiredElem.Validate(teatypes.Int(i))
			if err != nil {
				err.Key(k)
				return
			}
			continue
		}
		err = ruleElem.Validate(teatypes.Int(i), teatypes.Option[Value](options[i]).Value)
		if err != nil {
			err.Key(k)
			return
		}
		continue
	}
	return
}

func SliceMap[Value ~[]Elem, Elem Validator[T], T any, Required required, Rule rule[Value], RequiredElem required](
	value Value, k string, required Required, rule Rule, requiredElem RequiredElem,
) (err Error) {
	if value == nil {
		err = required.Validate(teatypes.String(k))
		return
	}
	err = rule.Validate(teatypes.String(k), value)
	if err != nil {
		return
	}
	for i := range value {
		if value[i] == nil {
			err = requiredElem.Validate(teatypes.Int(i))
			if err != nil {
				err.Key(k)
				return
			}
			continue
		}
		err = value[i].Validate()
		if err != nil {
			err.Key(teatypes.Int(i).ToString())
			err.Key(k)
			return
		}
		continue
	}
	return
}

func Map[Value Validator[T], T any, Required required](
	value Value, k string, required Required,
) (err Error) {
	if value == nil {
		err = required.Validate(teatypes.String(k))
		return
	}
	err = value.Validate()
	if err != nil {
		err.Key(k)
		return
	}
	return
}
