package fswap

import (
	"testing"

	"github.com/fox-one/4swap-sdk-go/v2/route"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_mergePaths(t *testing.T) {
	type args struct {
		paths route.Paths
		p     route.Path
	}
	tests := []struct {
		name string
		args args
		want route.Paths
	}{
		{
			name: "merge into empty paths",
			args: args{
				paths: route.Paths{},
				p:     route.Path{Weight: 100, Pairs: []uint16{1, 2, 3}},
			},
			want: route.Paths{
				{Weight: 100, Pairs: []uint16{1, 2, 3}},
			},
		},
		{
			name: "merge into same pairs paths",
			args: args{
				paths: route.Paths{
					{Weight: 80, Pairs: []uint16{1, 2, 3}},
				},
				p: route.Path{Weight: 20, Pairs: []uint16{1, 2, 3}},
			},
			want: route.Paths{
				{Weight: 100, Pairs: []uint16{1, 2, 3}},
			},
		},
		{
			name: "merge into different pairs paths",
			args: args{
				paths: route.Paths{
					{Weight: 80, Pairs: []uint16{1, 2, 3}},
				},
				p: route.Path{Weight: 20, Pairs: []uint16{1, 2, 4}},
			},
			want: route.Paths{
				{Weight: 80, Pairs: []uint16{1, 2, 3}},
				{Weight: 20, Pairs: []uint16{1, 2, 4}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, mergePaths(tt.args.paths, tt.args.p), "mergePaths(%v, %v)", tt.args.paths, tt.args.p)
		})
	}
}

func Test_equalPairs(t *testing.T) {
	type args struct {
		a []uint16
		b []uint16
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "equal",
			args: args{
				a: []uint16{1, 2, 3},
				b: []uint16{1, 2, 3},
			},
			want: true,
		},
		{
			name: "not equal",
			args: args{
				a: []uint16{1, 2, 3},
				b: []uint16{1, 2, 4},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, equalPairs(tt.args.a, tt.args.b), "equalPairs(%v, %v)", tt.args.a, tt.args.b)
		})
	}
}

func TestMergeOrders(t *testing.T) {
	orders := []*Order{
		{
			PayAssetID:  "pUSD",
			PayAmount:   decimal.NewFromFloat(100),
			FillAssetID: "USDT",
			FillAmount:  decimal.NewFromFloat(100.1),
			Paths: route.Paths{
				{Weight: 100, Pairs: []uint16{1, 2, 3}},
			},
			PriceImpact: decimal.NewFromFloat(0.001),
		},
		{
			PayAssetID:  "pUSD",
			PayAmount:   decimal.NewFromFloat(64),
			FillAssetID: "USDT",
			FillAmount:  decimal.NewFromFloat(63.9),
			Paths: route.Paths{
				{Weight: 100, Pairs: []uint16{1, 2, 3}},
			},
			PriceImpact: decimal.NewFromFloat(0.002),
		},
		{
			PayAssetID:  "pUSD",
			PayAmount:   decimal.NewFromFloat(36),
			FillAssetID: "USDT",
			FillAmount:  decimal.NewFromFloat(36.2),
			Paths: route.Paths{
				{Weight: 100, Pairs: []uint16{1, 4, 3}},
			},
			PriceImpact: decimal.NewFromFloat(0.003),
		},
	}

	m := MergeOrders(orders)
	assert.Equal(t, "pUSD", m.PayAssetID)
	assert.Equal(t, decimal.NewFromFloat(100+64+36).String(), m.PayAmount.String())
	assert.Equal(t, "USDT", m.FillAssetID)
	assert.Equal(t, decimal.NewFromFloat(100.1+63.9+36.2).String(), m.FillAmount.String())
	assert.Equal(t, "82:1,2,3;18:1,4,3", m.Paths.String())
	assert.Equal(t, decimal.NewFromFloat(0.00168).String(), m.PriceImpact.String())
}
