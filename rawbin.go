package rawbin

import (
	"errors"
)

// ErrNotPointer ...
var ErrNotPointer = errors.New("given value is not a pointer")

// Marshaller ...
type Marshaller interface {
	RawMarshal(writer *Writer) error
}

// Unmarshaller ...
type Unmarshaller interface {
	RawUnmarshal(reader *Reader) error
}
