package routes

import (
	"fmt"
	"strings"

	"github.com/pandodao/mtg/mtgpack"
	"github.com/shopspring/decimal"
)

type Path struct {
	Amount decimal.Decimal `json:"amount"`
	Routes Routes          `json:"routes,omitempty"`
}

// String returns the string representation of the path.
// eg: 10.1:1,2,3
func (p Path) String() string {
	return fmt.Sprintf("%s:%s", p.Amount.String(), p.Routes.String())
}

// ParsePath returns a new path from the given string.
func ParsePath(s string) (Path, error) {
	var p Path

	a, b, ok := strings.Cut(s, ":")
	if !ok {
		return p, fmt.Errorf("invalid path string: %s", s)
	}

	var err error
	p.Amount, err = decimal.NewFromString(a)
	if err != nil {
		return p, fmt.Errorf("parse path amount failed: %w", err)
	}

	p.Routes, err = ParseRoutes(b)
	if err != nil {
		return p, fmt.Errorf("parse path routes failed: %w", err)
	}

	return p, nil
}

func (p Path) Equal(other Path) bool {
	return p.Amount.Equal(other.Amount) && p.Routes.Cmp(other.Routes) == 0
}

func (p Path) EncodeMtg(enc *mtgpack.Encoder) error {
	return enc.EncodeValues(p.Amount, p.Routes)
}

func (p *Path) DecodeMtg(dec *mtgpack.Decoder) error {
	return dec.DecodeValues(&p.Amount, &p.Routes)
}
