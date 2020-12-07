package encoder

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
)

func Scan(body []byte, dest ...interface{}) ([]byte, error) {
	r := bytes.NewReader(body)

	for _, dp := range dest {
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		n := int(b)
		if n == 0 {
			continue
		}

		data := make([]byte, n)
		if _, err := io.ReadFull(r, data); err != nil {
			return nil, err
		}

		if err := scan(data, dp); err != nil {
			return nil, err
		}
	}

	return ioutil.ReadAll(r)
}

func scan(data []byte, dest interface{}) error {
	if u, ok := dest.(encoding.BinaryUnmarshaler); ok {
		return u.UnmarshalBinary(data)
	}

	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("cannot scan %v", v.Kind())
	}

	v = v.Elem()

	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		x, _ := binary.Varint(data)
		if v.OverflowInt(x) {
			return fmt.Errorf("cannot put %v", x)
		}

		v.SetInt(x)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		x, _ := binary.Uvarint(data)
		if v.OverflowUint(x) {
			return fmt.Errorf("cannot put %v", x)
		}

		v.SetUint(x)
	case reflect.String:
		v.SetString(string(data))
	default:
		return fmt.Errorf("mtg: cannot scan %v", dest)
	}

	return nil
}
