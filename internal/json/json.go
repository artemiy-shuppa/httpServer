package json

import (
	"encoding/json"
	"io"
)

func Parse[T any](r io.Reader, t T) error {
	dec := json.NewDecoder(r)
	return dec.Decode(t)
}
func Format[T any](w io.Writer, t T) error {
	enc := json.NewEncoder(w)
	return enc.Encode(t)
}
