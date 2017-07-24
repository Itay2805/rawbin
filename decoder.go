package rawbin

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
)

// DecodeFromBytes ...
func DecodeFromBytes(value interface{}, data []byte) error {
	return DecodeFromBytesReader(value, bytes.NewReader(data))
}

// DecodeFromReader ...
func DecodeFromReader(value interface{}, reader io.Reader) error {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	return DecodeFromBytesReader(value, bytes.NewReader(b))
}

// DecodeFromBytesReader ...
func DecodeFromBytesReader(value interface{}, reader *bytes.Reader) error {
	reflectValue := reflect.ValueOf(value)
	if reflectValue.Kind() != reflect.Ptr {
		return ErrNotPointer
	}
	if reflect.Indirect(reflectValue).Kind() == reflect.Struct {
		return decodeStruct(value, reader)
	}
	return decodeValue(value, reader)
}

func decodeStruct(value interface{}, reader *bytes.Reader) error {
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

// Always big endian for now

// supports:
//		BinerUnmarshaller (for custom types)
//			built in supported types:
//				binner.Varint
// 				binner.Varlong
//				binner.String (varint + data)
//		int(8-64) & uint(8-64)
// 		string (int16 + data)
//		bool (true - 1, false - !1)

func decodeValue(value interface{}, reader *bytes.Reader) error {
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
