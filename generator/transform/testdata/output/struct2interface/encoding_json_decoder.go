package struct2interface

import (
	"io"
	"encoding/json"
)

// Decoder :
type Decoder interface {
	UseNumber()
	DisallowUnknownFields()
	Decode(v interface{}) error
	Buffered() io.Reader
	Token() (json.Token, error)
	More() bool
}
