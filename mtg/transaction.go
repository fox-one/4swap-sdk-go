package mtg

//go:generate enumer -type TransactionType -trimprefix TransactionType

type TransactionType int

const (
	_                     TransactionType = iota
	TransactionTypeAdd                    // 1 加注
	TransactionTypeRemove                 // 2 提取
	TransactionTypeSwap                   // 3 兑换
	_
	_
	TransactionTypeSwapV2 // 6 兑换 v2
)

func ParseTransactionType(t string) TransactionType {
	typ, _ := TransactionTypeString(t)
	return typ
}
