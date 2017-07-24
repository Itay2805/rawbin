package rawbin

import (
	"bytes"
	"errors"
)

// ErrNotPointer ...
var ErrNotPointer = errors.New("given value is not a pointer")

// Marshaller ...
type Marshaller interface {
	RawMarshal(buffer *bytes.Buffer) error
}

// Unmarshaller ...
type Unmarshaller interface {
	RawUnmarshal(reader *bytes.Reader) error
}
