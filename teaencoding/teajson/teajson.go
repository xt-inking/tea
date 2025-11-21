package teajson

import (
	"io"

	"github.com/bytedance/sonic"
)

func Decode(r io.Reader, v any) error {
	return config.NewDecoder(r).Decode(v)
}

func Encode(w io.Writer, v any) error {
	return config.NewEncoder(w).Encode(v)
}

func Marshal(v any) ([]byte, error) {
	return config.Marshal(v)
}

func Unmarshal(data []byte, v any) error {
	return config.Unmarshal(data, v)
}

var config = sonic.Config{
	NoQuoteTextMarshaler:    true,
	UseNumber:               true,
	NoValidateJSONMarshaler: true,
	NoValidateJSONSkip:      true,
	NoEncoderNewline:        true,
}.Froze()
