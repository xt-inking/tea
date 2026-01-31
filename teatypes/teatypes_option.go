package teatypes

import (
	"github.com/tea-frame-go/tea/teaencoding/teajson"
)

type Option[Value any] struct {
	Value Value
	Valid bool
}

func (option Option[Value]) IsZero() bool {
	return !option.Valid
}

func (option Option[Value]) MarshalJSON() ([]byte, error) {
	if !option.Valid {
		return []byte("null"), nil
	}
	return teajson.Marshal(option.Value)
}

func (option *Option[Value]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		option.Valid = false
		return nil
	}
	option.Valid = true
	return teajson.Unmarshal(data, &option.Value)
}
