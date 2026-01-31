package teavalid

import (
	"fmt"
	"regexp"
	"slices"
	"unicode/utf8"

	"github.com/tea-frame-go/tea/teatypes"
	"github.com/tea-frame-go/tea/teatypes/teaconstraints"
)

type rule[Value any] interface {
	Validate(k teatypes.ToStringer, value Value) (err Error)
}

func Nop[Value any]() nop[Value] {
	return nop[Value]{}
}

type nop[Value any] struct{}

func (nop[Value]) Validate(k teatypes.ToStringer, value Value) (err Error) {
	return
}

func Between[Value teaconstraints.Integer](min Value, max Value) betweenRule[Value] {
	return betweenRule[Value]{min, max}
}

type betweenRule[Value teaconstraints.Integer] struct {
	min Value
	max Value
}

func (rule betweenRule[Value]) Validate(k teatypes.ToStringer, value Value) (err Error) {
	if value < rule.min || value > rule.max {
		err = NewError(k.ToString(), fmt.Sprintf(
			"value `%d` must be between `%d` and `%d`",
			value, rule.min, rule.max,
		))
		return
	}
	return
}

func Min[Value teaconstraints.Integer](min Value) minRule[Value] {
	return minRule[Value]{min}
}

type minRule[Value teaconstraints.Integer] struct {
	min Value
}

func (rule minRule[Value]) Validate(k teatypes.ToStringer, value Value) (err Error) {
	if value < rule.min {
		err = NewError(k.ToString(), fmt.Sprintf(
			"value `%d` min must be `%d`",
			value, rule.min,
		))
		return
	}
	return
}

func Max[Value teaconstraints.Integer](max Value) maxRule[Value] {
	return maxRule[Value]{max}
}

type maxRule[Value teaconstraints.Integer] struct {
	max Value
}

func (rule maxRule[Value]) Validate(k teatypes.ToStringer, value Value) (err Error) {
	if value > rule.max {
		err = NewError(k.ToString(), fmt.Sprintf(
			"value `%d` max must be `%d`",
			value, rule.max,
		))
		return
	}
	return
}

func In[Value teaconstraints.Integer | teaconstraints.String](in ...Value) inRule[Value] {
	return inRule[Value]{in}
}

type inRule[Value teaconstraints.Integer | teaconstraints.String] struct {
	in []Value
}

func (rule inRule[Value]) Validate(k teatypes.ToStringer, value Value) (err Error) {
	if slices.Contains(rule.in, value) {
		return
	}
	err = NewError(k.ToString(), fmt.Sprintf(
		"value `%v` must be in `%v`",
		value, rule.in,
	))
	return
}

func Length[Value teaconstraints.String](length int) lengthRule[Value] {
	return lengthRule[Value]{length}
}

type lengthRule[Value teaconstraints.String] struct {
	length int
}

func (rule lengthRule[Value]) Validate(k teatypes.ToStringer, value Value) (err Error) {
	length := utf8.RuneCountInString(string(value))
	if length != rule.length {
		err = NewError(k.ToString(), fmt.Sprintf(
			"value `%s` length must be `%d`",
			value, rule.length,
		))
		return
	}
	return
}

func LengthBetween[Value teaconstraints.String](min int, max int) lengthBetweenRule[Value] {
	return lengthBetweenRule[Value]{min, max}
}

type lengthBetweenRule[Value teaconstraints.String] struct {
	min int
	max int
}

func (rule lengthBetweenRule[Value]) Validate(k teatypes.ToStringer, value Value) (err Error) {
	length := utf8.RuneCountInString(string(value))
	if length < rule.min || length > rule.max {
		err = NewError(k.ToString(), fmt.Sprintf(
			"value `%s` length must be between `%d` and `%d`",
			value, rule.min, rule.max,
		))
		return
	}
	return
}

func LengthMin[Value teaconstraints.String](min int) lengthMinRule[Value] {
	return lengthMinRule[Value]{min}
}

type lengthMinRule[Value teaconstraints.String] struct {
	min int
}

func (rule lengthMinRule[Value]) Validate(k teatypes.ToStringer, value Value) (err Error) {
	length := utf8.RuneCountInString(string(value))
	if length < rule.min {
		err = NewError(k.ToString(), fmt.Sprintf(
			"value `%s` length min must be `%d`",
			value, rule.min,
		))
		return
	}
	return
}

func LengthMax[Value teaconstraints.String](max int) lengthMaxRule[Value] {
	return lengthMaxRule[Value]{max}
}

type lengthMaxRule[Value teaconstraints.String] struct {
	max int
}

func (rule lengthMaxRule[Value]) Validate(k teatypes.ToStringer, value Value) (err Error) {
	length := utf8.RuneCountInString(string(value))
	if length > rule.max {
		err = NewError(k.ToString(), fmt.Sprintf(
			"value `%s` length max must be `%d`",
			value, rule.max,
		))
		return
	}
	return
}

func Regexp[Value teaconstraints.String](regexp *regexp.Regexp) regexpRule[Value] {
	return regexpRule[Value]{regexp}
}

type regexpRule[Value teaconstraints.String] struct {
	regexp *regexp.Regexp
}

func (rule regexpRule[Value]) Validate(k teatypes.ToStringer, value Value) (err Error) {
	if rule.regexp.MatchString(string(value)) {
		return
	}
	err = NewError(k.ToString(), fmt.Sprintf(
		"value `%s` regexp must be `%s`",
		value, rule.regexp.String(),
	))
	return
}

func SliceLength[Value ~[]Elem, Elem any](length int) sliceLengthRule[Value, Elem] {
	return sliceLengthRule[Value, Elem]{length}
}

type sliceLengthRule[Value ~[]Elem, Elem any] struct {
	length int
}

func (rule sliceLengthRule[Value, Elem]) Validate(k teatypes.ToStringer, value Value) (err Error) {
	length := len(value)
	if length != rule.length {
		err = NewError(k.ToString(), fmt.Sprintf(
			"value `%v` length must be `%d`",
			value, rule.length,
		))
		return
	}
	return
}

func SliceLengthBetween[Value ~[]Elem, Elem any](min int, max int) sliceLengthBetweenRule[Value, Elem] {
	return sliceLengthBetweenRule[Value, Elem]{min, max}
}

type sliceLengthBetweenRule[Value ~[]Elem, Elem any] struct {
	min int
	max int
}

func (rule sliceLengthBetweenRule[Value, Elem]) Validate(k teatypes.ToStringer, value Value) (err Error) {
	length := len(value)
	if length < rule.min || length > rule.max {
		err = NewError(k.ToString(), fmt.Sprintf(
			"value `%v` length must be between `%d` and `%d`",
			value, rule.min, rule.max,
		))
		return
	}
	return
}

func SliceLengthMin[Value ~[]Elem, Elem any](min int) sliceLengthMinRule[Value, Elem] {
	return sliceLengthMinRule[Value, Elem]{min}
}

type sliceLengthMinRule[Value ~[]Elem, Elem any] struct {
	min int
}

func (rule sliceLengthMinRule[Value, Elem]) Validate(k teatypes.ToStringer, value Value) (err Error) {
	length := len(value)
	if length < rule.min {
		err = NewError(k.ToString(), fmt.Sprintf(
			"value `%v` length min must be `%d`",
			value, rule.min,
		))
		return
	}
	return
}

func SliceLengthMax[Value ~[]Elem, Elem any](max int) sliceLengthMaxRule[Value, Elem] {
	return sliceLengthMaxRule[Value, Elem]{max}
}

type sliceLengthMaxRule[Value ~[]Elem, Elem any] struct {
	max int
}

func (rule sliceLengthMaxRule[Value, Elem]) Validate(k teatypes.ToStringer, value Value) (err Error) {
	length := len(value)
	if length > rule.max {
		err = NewError(k.ToString(), fmt.Sprintf(
			"value `%v` length max must be `%d`",
			value, rule.max,
		))
		return
	}
	return
}
