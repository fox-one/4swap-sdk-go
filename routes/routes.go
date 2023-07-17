package routes

import (
	"fmt"
	"math"
	"strings"

	"github.com/pandodao/mtg/mtgpack"
	"github.com/spf13/cast"
)

type Routes []int64

func (r Routes) String() string {
	var b strings.Builder
	for idx, id := range r {
		if idx > 0 {
			b.WriteByte(',')
		}

		b.WriteString(cast.ToString(id))
	}

	return b.String()
}

func ParseRoutes(s string) (Routes, error) {
	fields := strings.FieldsFunc(s, func(r rune) bool { return r == ',' })
	ids := make(Routes, len(fields))
	for idx, field := range fields {
		id, err := cast.ToInt64E(field)
		if err != nil {
			return nil, fmt.Errorf("parse route id failed: %w", err)
		}

		ids[idx] = id
	}

	return ids, nil
}

// hashid

func (r Routes) HashString() string {
	s, err := h.EncodeInt64(r)
	if err != nil {
		return ""
	}

	return s
}

func ParseHashedRoutes(v string) Routes {
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

func (r Routes) EncodeMtg(enc *mtgpack.Encoder) error {
	values := []any{uint8(len(r))}
	for _, id := range r {
		if id > int64(math.MaxUint16) {
			return fmt.Errorf("path id %d is too large", id)
		}

		values = append(values, uint16(id))
	}

	return enc.EncodeValues(values...)
}

func (r *Routes) DecodeMtg(dec *mtgpack.Decoder) error {
	count, err := dec.DecodeUint8()
	if err != nil {
		return err
	}

	for i := count; i > 0; i-- {
		var id uint16
		if err := dec.DecodeValue(&id); err != nil {
			return err
		}

		*r = append(*r, int64(id))
	}

	return nil
}

func (r Routes) MarshalBinary() (data []byte, err error) {
	return []byte(r.HashString()), nil
}

func (r *Routes) UnmarshalBinary(data []byte) error {
	*r = ParseHashedRoutes(string(data))
	return nil
}
