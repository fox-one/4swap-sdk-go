package encoder

import (
	"database/sql/driver"
	"errors"
)

type RawMessage []byte

func (r RawMessage) MarshalBinary() ([]byte, error) {
	b := make([]byte, len(r))
	copy(b, r)
	return b, nil
}

func (r *RawMessage) UnmarshalBinary(data []byte) error {
	*r = append((*r)[0:0], data...)
	return nil
}

func (r RawMessage) Value() (driver.Value, error) {
	return []byte(r), nil
}

func (r *RawMessage) Scan(src interface{}) error {
	var source []byte
	switch v := src.(type) {
	case string:
		source = []byte(v)
	case []byte:
		source = v
	default:
		return errors.New("incompatible type for RawMessage")
	}

	*r = append((*r)[0:0], source...)
	return nil
}
