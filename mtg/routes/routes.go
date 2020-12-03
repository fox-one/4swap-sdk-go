package routes

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/spf13/cast"
)

type Routes []int64

func (r Routes) String() string {
	s, err := h.EncodeInt64(r)
	if err != nil {
		return ""
	}

	return s
}

func NewFromString(v string) Routes {
	ids, _ := h.DecodeInt64WithError(v)
	return ids
}

func (r Routes) Cmp(other Routes) int {
	if diff := len(r) - len(other); diff < 0 {
		return 1
	} else if diff > 0 {
		return -1
	}

	for idx := range r {
		if diff := r[idx] - other[idx]; diff < 0 {
			return 1
		} else if diff > 0 {
			return -1
		}
	}

	return 0
}

// json encoding

func (r Routes) MarshalJSON() ([]byte, error) {
	s := r.String()
	b := make([]byte, 0, len(s)+2)
	buf := bytes.NewBuffer(b)
	buf.WriteByte('"')
	buf.WriteString(s)
	buf.WriteByte('"')

	return buf.Bytes(), nil
}

func (r *Routes) UnmarshalJSON(b []byte) error {
	if len(b) > 2 && b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}

	*r = NewFromString(string(b))
	return nil
}

// binary encoding

func (r Routes) MarshalBinary() (data []byte, err error) {
	return []byte(r.String()), nil
}

func (r *Routes) UnmarshalBinary(data []byte) error {
	*r = NewFromString(string(data))
	return nil
}

// sql

func (r Routes) Value() (driver.Value, error) {
	return json.Marshal([]int64(r))
}

func (r *Routes) Scan(src interface{}) error {
	v := cast.ToString(src)

	var ids []int64
	if err := json.Unmarshal([]byte(v), &ids); err == nil {
		*r = ids
	} else {
		*r = NewFromString(v)
	}

	return nil
}
