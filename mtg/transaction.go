package mtg

import (
	"strings"
)

//go:generate stringer -type TransactionType -trimprefix TransactionType

type TransactionType int

const (
	_                     TransactionType = iota
	TransactionTypeAdd                    // 1 加注
	TransactionTypeRemove                 // 2 提取
	TransactionTypeSwap                   // 3 兑换
)

func ParseTransactionType(t string) TransactionType {
	for idx := 0; idx < len(_TransactionType_index)-1; idx++ {
		l, r := _TransactionType_index[idx], _TransactionType_index[idx+1]
		if typ := _TransactionType_name[l:r]; strings.EqualFold(typ, t) {
			return TransactionType(idx + 1)
		}
	}

	return 0
}
