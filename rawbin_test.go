package rawbin

import (
	"bytes"
	"testing"
)

// TODO: Look at generating code to generate all the test cases

type TestStruct struct {
	NumUint8  uint8
	NumUint16 uint16
	NumUint32 uint32
	NumUint64 uint64
	NumInt8   int8
	NumInt16  int16
	NumInt32  int32
	NumInt64  int64
	Boolean   bool
	Str       string
	Varint    Varint
	Varlong   Varlong
	String    String
}

func TestEncodeDecodeStruct(t *testing.T) {
	expected := TestStruct{125, 456, 1235, 12345, -12, -45, -987, -64554, true, "Some string", 123, 12345, "Another String"}
	var val TestStruct

	buffer := bytes.NewBuffer([]byte{})
	err := Encode(expected, buffer)
	if err != nil {
		t.Error(err)
	}
	reader := bytes.NewReader(buffer.Bytes())
	err = Decode(&val, reader)
	if err != nil {
		t.Error(err)
	}
	if val != expected {
		t.Errorf("Expected %v, got %v", expected, val)
	}
}
