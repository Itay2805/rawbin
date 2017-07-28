package rawbin

import "io"

// Writer wraps io.Writer and implements io.ByteWriter for it
type Writer struct {
	baseWriter io.Writer
}

// Write for implementing io.Writer
func (w Writer) Write(p []byte) (n int, err error) {
	return w.baseWriter.Write(p)
}

// WriteByte for implementing io.ByteWriter
func (w Writer) WriteByte(c byte) error {
	_, err := w.baseWriter.Write([]byte{c})
	return err
}
