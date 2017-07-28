package rawbin

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

// DecodeBytes is used for decoding data which is stored in a byte array
func DecodeBytes(value interface{}, data []byte) error {
	return Decode(value, bytes.NewReader(data))
}

// Decode is used to decode data from a reader
func Decode(value interface{}, r io.Reader) error {
	reflectValue := reflect.ValueOf(value)
	if reflectValue.Kind() != reflect.Ptr {
		return ErrNotPointer
	}

	reader := Reader{r}

	if reflect.Indirect(reflectValue).Kind() == reflect.Struct {
		return decodeStruct(value, &reader)
	}
	return decodeValue(value, &reader)
}

func decodeStruct(value interface{}, reader *Reader) error {
	reflectValue := reflect.Indirect(reflect.ValueOf(value))
	for i := 0; i < reflectValue.NumField(); i++ {
		field := reflectValue.Field(i)
		if field.CanInterface() && field.CanSet() {
			data := reflect.New(field.Type()).Interface()
			if field.Kind() == reflect.Struct {
				err := decodeStruct(data, reader)
				if err != nil {
					return err
				}
				field.Set(reflect.Indirect(reflect.ValueOf(data)))
			} else {
				err := decodeValue(data, reader)
				if err != nil {
					return err
				}
				field.Set(reflect.Indirect(reflect.ValueOf(data)))
			}
		} else {
			fmt.Printf("Can not set field named %v (Maybe not exported?)\n", field)
		}
	}
	return nil
}

func decodeValue(value interface{}, reader *Reader) error {
	switch val := value.(type) {
	case Unmarshaller:
		return val.RawUnmarshal(reader)
	case *int8, *int16, *int32, *int64, *uint8, *uint16, *uint32, *uint64, *float32, *float64:
		return binary.Read(reader, binary.BigEndian, value)
	case *string:
		var length int16
		err := binary.Read(reader, binary.BigEndian, &length)
		stringData := make([]byte, length, length)
		_, err = reader.Read(stringData)
		if err != nil {
			return err
		}
		*val = string(stringData)
		return nil
	case *bool:
		b, err := reader.ReadByte()
		if err != nil {
			return err
		}
		if b == 1 {
			*val = true
		} else {
			*val = false
		}
		return nil
	default:
		return fmt.Errorf("invalid type %T", value)
	}
}
