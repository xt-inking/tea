package teajson

import (
	"encoding/json"
	"io"
)

func Decode(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}

func Encode(w io.Writer, v any) error {
	return json.NewEncoder(w).Encode(v)
}
