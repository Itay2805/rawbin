package rawbin

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

// EncodeBytes will encode the given struct or value and return it in the bytes array
func EncodeBytes(value interface{}) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	err := Encode(value, buffer)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// Encode will encode the given struct value and will write it to the given writer
func Encode(value interface{}, w io.Writer) error {
	// make sure not a pointer
	reflectValue := reflect.ValueOf(value)
	if reflectValue.Kind() == reflect.Ptr {
		value = reflect.Indirect(reflectValue).Interface()
		reflectValue = reflect.ValueOf(value)
	}

	writer := Writer{w}

	if reflectValue.Kind() == reflect.Struct {
		return encodeStruct(value, &writer)
	}
	return encodeValue(value, &writer)
}

func encodeStruct(value interface{}, writer *Writer) error {
	reflectValue := reflect.ValueOf(value)
	for i := 0; i < reflectValue.NumField(); i++ {
		field := reflectValue.Field(i)
		if field.CanInterface() {
			data := field.Interface()
			if field.Kind() == reflect.Struct {
				err := encodeStruct(data, writer)
				if err != nil {
					return err
				}
			} else {
				err := encodeValue(data, writer)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func encodeValue(value interface{}, Writer *Writer) error {
	switch val := value.(type) {
	case Marshaller:
		return val.RawMarshal(Writer)
	case encoding.BinaryMarshaler:
		bytes, err := val.MarshalBinary()
		if err != nil {
			return err
		}
		_, err = Writer.Write(bytes)
		if err != nil {
			return err
		}
		return nil
	case int8, int16, int32, int64, uint8, uint16, uint32, uint64, float32, float64:
		return binary.Write(Writer, binary.BigEndian, val)
	case string:
		binary.Write(Writer, binary.BigEndian, int16(len(val)))
		_, err := Writer.Write([]byte(val))
		if err != nil {
			return err
		}
		return nil
	case bool:
		if val {
			return Writer.WriteByte(byte(1))
		}
		return Writer.WriteByte(byte(0))
	default:
		return fmt.Errorf("Invalid type %T", value)
	}
}
