package rawbin

import "bytes"
import "encoding/binary"
import "math"

import "fmt"

// Varint is an integer between -2147483648 and 2147483647
// Variable-length data encoding a two's complement signed 32-bit integer
type Varint int32

// RawMarshal is used by rawbin to encode the data
func (v Varint) RawMarshal(buffer *bytes.Buffer) error {
	bytes := make([]byte, 5, 5)
	size := binary.PutVarint(bytes, int64(v))
	buffer.Write(bytes[:size])
	return nil
}

// RawUnmarshal is used by rawbin to decode the data
// TODO: Might need to write the reader because binary.ReadVarint is meant to read int64 and not int32 so it might read more then needed
func (v *Varint) RawUnmarshal(reader *bytes.Reader) error {
	i, err := binary.ReadVarint(reader)
	if err != nil {
		return err
	}
	if i < math.MinInt32 || i > math.MaxInt32 {
		return fmt.Errorf("tried to read varint to big (%v)", i)
	}
	*v = Varint(i)
	return nil
}

// Varlong is an integer between -9223372036854775808 and 9223372036854775807
//Variable-length data encoding a two's complement signed 64-bit integer
type Varlong int64

// RawMarshal is used by rawbin to encode the data
func (v Varlong) RawMarshal(buffer *bytes.Buffer) error {
	bytes := make([]byte, 10, 10)
	size := binary.PutVarint(bytes, int64(v))
	buffer.Write(bytes[:size])
	return nil
}

// RawUnmarshal is used by rawbin to decode the data
func (v *Varlong) RawUnmarshal(reader *bytes.Reader) error {
	i, err := binary.ReadVarint(reader)
	if err != nil {
		return err
	}

	*v = Varlong(i)
	return nil
}

// String is a sequence of Unicode scalar values
// UTF-8 string prefixed with its size in bytes as a VarInt.
// Maximum length of n characters, which varies by context; up to n × 4 bytes can be used to encode n characters and both of those limits are checked.
// Maximum n value is 32767. The + 3 is due to the max size of a valid length VarInt.
type String string

// RawMarshal is used by rawbin to encode the data
func (v String) RawMarshal(buffer *bytes.Buffer) error {
	length := Varlong(len(v))
	err := length.RawMarshal(buffer)
	if err != nil {
		return err
	}
	_, err = buffer.WriteString(string(v))
	return err
}

// RawUnmarshal is used by rawbin to decode the data
func (v *String) RawUnmarshal(reader *bytes.Reader) error {
	var length Varlong
	length.RawUnmarshal(reader)
	bytes := make([]byte, length, length)
	_, err := reader.Read(bytes)
	*v = String(bytes)
	return err
}
