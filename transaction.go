package fswap

import (
	"encoding/base64"

	"github.com/vmihailenco/msgpack/v5"
)

const (
	TransactionTypeAdd    = "Add"
	TransactionTypeRemove = "Remove"
	TransactionTypeSwap   = "Swap"
)

type TransactionAction struct {
	// Transaction type add remove swap
	Type string `json:"t,omitempty" msgpack:"t,omitempty"`
	// deposit
	Deposit string `json:"d,omitempty" msgpack:"d,omitempty"`
	// withdraw
	Pairs         []string `json:"p,omitempty" msgpack:"p,omitempty"`
	RemovePercent int64    `json:"l,omitempty" msgpack:"l,omitempty"`
	// Swap

	// 要买的币的 asset id
	AssetID string `json:"a,omitempty" msgpack:"a,omitempty"`
	// 路径 id，由 preOrder 得到，为空则由系统分配
	Routes string `json:"r,omitempty" msgpack:"r,omitempty"`
	// 最小买入数量，为空则不限制
	Minimum string `json:"m,omitempty" msgpack:"m,omitempty"`
}

func EncodeAction(action TransactionAction) string {
	b, err := msgpack.Marshal(action)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(b)
}
