package teatypes

import (
	"strconv"
)

type Int int

type String string

type ToStringer interface {
	ToString() string
}

func (value Int) ToString() string {
	return strconv.Itoa(int(value))
}

func (value String) ToString() string {
	return string(value)
}
