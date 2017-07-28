package rawbin

import "io"

// Reader for wrapping any io.Reader and implementing io.ByteReader for it
type Reader struct {
	baseReader io.Reader
}

// Read implements io.Reader
func (r Reader) Read(p []byte) (n int, err error) {
	return r.baseReader.Read(p)
}

// ReadByte implements io.ByteReader
func (r Reader) ReadByte() (byte, error) {
	// TODO: Maybe turn this to a global so won't have to recreate every time
	single := make([]byte, 1, 1)
	_, err := r.baseReader.Read(single)
	return single[0], err
}
