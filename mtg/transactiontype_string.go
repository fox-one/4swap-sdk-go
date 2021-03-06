// Code generated by "stringer -type TransactionType -trimprefix TransactionType"; DO NOT EDIT.

package mtg

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TransactionTypeAdd-1]
	_ = x[TransactionTypeRemove-2]
	_ = x[TransactionTypeSwap-3]
}

const _TransactionType_name = "AddRemoveSwap"

var _TransactionType_index = [...]uint8{0, 3, 9, 13}

func (i TransactionType) String() string {
	i -= 1
	if i < 0 || i >= TransactionType(len(_TransactionType_index)-1) {
		return "TransactionType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _TransactionType_name[_TransactionType_index[i]:_TransactionType_index[i+1]]
}
