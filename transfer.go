package fswap

import (
	"github.com/shopspring/decimal"
)

type TransferReq struct {
	AssetID string `json:"asset_id,omitempty"`
	OpponentID string `json:"opponent_id,omitempty"`
	TraceID string `json:"trace_id,omitempty"`
	Amount decimal.Decimal `json:"amount,omitempty"`
	Memo string `json:"memo,omitempty"`
}
