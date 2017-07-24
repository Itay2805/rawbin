package rawbin

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

// EncodeToBytes will encode the given struct or value and return it in the bytes array
func EncodeToBytes(value interface{}) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	err := EncodeToBytesBuffer(value, buffer)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// EncodeToWriter will encode the given struct value and will write it to the given writer
func EncodeToWriter(value interface{}, writer io.Writer) error {
	bytes, err := EncodeToBytes(value)
	writer.Write(bytes)
	return err
}

// EncodeToBytesBuffer will encode the given struct or value and will write it to the given buffer
func EncodeToBytesBuffer(value interface{}, buffer *bytes.Buffer) error {
	// make sure not a pointer
	reflectValue := reflect.ValueOf(value)
	if reflectValue.Kind() == reflect.Ptr {
		value = reflect.Indirect(reflectValue).Interface()
		reflectValue = reflect.ValueOf(value)
	}
	if reflectValue.Kind() == reflect.Struct {
		return encodeStruct(value, buffer)
	}
	return encodeValue(value, buffer)
}

func encodeStruct(value interface{}, buffer *bytes.Buffer) error {
	reflectValue := reflect.ValueOf(value)
	for i := 0; i < reflectValue.NumField(); i++ {
		field := reflectValue.Field(i)
		if field.CanInterface() {
			data := field.Interface()
			if field.Kind() == reflect.Struct {
				err := encodeStruct(data, buffer)
				if err != nil {
					return err
				}
			} else {
				err := encodeValue(data, buffer)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func encodeValue(value interface{}, buffer *bytes.Buffer) error {
	switch val := value.(type) {
	case Marshaller:
		return val.RawMarshal(buffer)
	case encoding.BinaryMarshaler:
		bytes, err := val.MarshalBinary()
		if err != nil {
			return err
		}
		_, err = buffer.Write(bytes)
		if err != nil {
			return err
		}
		return nil
	case int8, int16, int32, int64, uint8, uint16, uint32, uint64, float32, float64:
		return binary.Write(buffer, binary.BigEndian, val)
	case string:
		binary.Write(buffer, binary.BigEndian, int16(len(val)))
		_, err := buffer.WriteString(val)
		if err != nil {
			return err
		}
		return nil
	case bool:
		if val {
			return buffer.WriteByte(byte(1))
		}
		return buffer.WriteByte(byte(0))
	default:
		return fmt.Errorf("Invalid type %T", value)
	}
}
