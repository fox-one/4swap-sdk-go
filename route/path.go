package route

import (
	"fmt"
	"strings"

	"github.com/pandodao/mtg/mtgpack"
	"github.com/shopspring/decimal"
)

type Path struct {
	Weight uint8    `json:"weight,omitempty"`
	Pairs  []uint16 `json:"pairs,omitempty"`
}

func (p Path) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d:", p.Weight))
	for idx, id := range p.Pairs {
		if idx > 0 {
			b.WriteByte(',')
		}

		b.WriteString(fmt.Sprintf("%d", id))
	}

	return b.String()
}

type Paths []Path

func Single(ids ...uint16) Paths {
	return Paths{{
		Weight: 100,
		Pairs:  ids,
	}}
}

func (p Paths) String() string {
	paths := make([]string, len(p))
	for idx, path := range p {
		paths[idx] = path.String()
	}

	return strings.Join(paths, ";")
}

func (p Path) Share() decimal.Decimal {
	w := decimal.NewFromInt(int64(p.Weight))
	return w.Shift(-2)
}

func (p Path) EncodeMtg(enc *mtgpack.Encoder) error {
	var values = []any{
		p.Weight,
		uint8(len(p.Pairs)),
	}

	for _, pair := range p.Pairs {
		values = append(values, pair)
	}

	return enc.EncodeValues(values...)
}

func (p *Path) DecodeMtg(dec *mtgpack.Decoder) error {
	var weight, count uint8
	if err := dec.DecodeValues(&weight, &count); err != nil {
		return err
	}

	pairs := make([]uint16, count)
	for idx := range pairs {
		id, err := dec.DecodeUint16()
		if err != nil {
			return err
		}

		pairs[idx] = id
	}

	p.Weight = weight
	p.Pairs = pairs
	return nil
}

func (p Paths) EncodeMtg(enc *mtgpack.Encoder) error {
	var values = []any{
		uint8(len(p)),
	}

	for _, path := range p {
		values = append(values, path)
	}

	return enc.EncodeValues(values...)
}

func (p *Paths) DecodeMtg(dec *mtgpack.Decoder) error {
	count, err := dec.DecodeUint8()
	if err != nil {
		return err
	}

	paths := make(Paths, count)
	for idx := range paths {
		if err := dec.DecodeValue(&paths[idx]); err != nil {
			return err
		}
	}

	*p = append((*p)[0:0], paths...)
	return nil
}
