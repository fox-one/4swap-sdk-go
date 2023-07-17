package routes

import (
	"encoding/json"
	"strings"

	"github.com/pandodao/mtg/mtgpack"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
)

type Group []Path

// Sum returns the sum of the group.
func (g Group) Sum() decimal.Decimal {
	var sum decimal.Decimal
	for _, p := range g {
		sum = sum.Add(p.Amount)
	}

	return sum
}

// String returns the string representation of the group.
// eg: 10.1:1,2,3|20.2:4,5,6
func (g Group) String() string {
	var b strings.Builder
	for idx, p := range g {
		if idx > 0 {
			b.WriteByte('|')
		}

		b.WriteString(p.String())
	}

	return b.String()
}

// ParseGroup returns a new group from the given string.
func ParseGroup(s string) (Group, error) {
	g := Group{}
	for _, p := range strings.FieldsFunc(s, func(r rune) bool { return r == '|' }) {
		path, err := ParsePath(p)
		if err != nil {
			return nil, err
		}

		g = append(g, path)
	}

	return g, nil
}

func SinglePath(amount decimal.Decimal, route Routes) Group {
	return Group{Path{Amount: amount, Routes: route}}
}

func (g Group) EncodeMtg(enc *mtgpack.Encoder) error {
	values := []any{uint8(len(g))}

	for _, p := range g {
		values = append(values, p)
	}

	return enc.EncodeValues(values...)
}

func (g *Group) DecodeMtg(dec *mtgpack.Decoder) error {
	count, err := dec.DecodeUint8()
	if err != nil {
		return err
	}

	for i := count; i > 0; i-- {
		var p Path
		if err := dec.DecodeValue(&p); err != nil {
			return err
		}

		*g = append(*g, p)
	}

	return nil
}

// sql

func (g Group) Value() (interface{}, error) {
	return g.String(), nil
}

func (g *Group) Scan(src interface{}) error {
	s := cast.ToString(src)
	if v, err := ParseGroup(s); err == nil {
		*g = v
		return nil
	}

	// decode legacy format
	var r Routes
	if err := json.Unmarshal([]byte(s), &r); err == nil {
		*g = append(*g, Path{Routes: r})
	}

	return nil
}
