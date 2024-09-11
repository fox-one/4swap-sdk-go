package fswap

import (
	"time"

	"github.com/fox-one/mixin-sdk-go/v2/mixinnet"
	"github.com/pandodao/mtg/protocol"
	"github.com/shopspring/decimal"
)

type Transfer struct {
	ID        string                    `json:"id,omitempty"`
	CreatedAt time.Time                 `json:"created_at,omitempty"`
	AssetID   string                    `json:"asset_id,omitempty"`
	Amount    decimal.Decimal           `json:"amount,omitempty"`
	Memo      string                    `json:"memo,omitempty"`
	Receiver  protocol.MultisigReceiver `json:"receiver,omitempty"`
	TxHash    mixinnet.Hash             `json:"tx_hash,omitempty"`
}
