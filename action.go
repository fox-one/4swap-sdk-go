package fswap

import (
	"encoding/base64"
	"time"

	"github.com/fox-one/4swap-sdk-go/v2/route"
	"github.com/google/uuid"
	"github.com/pandodao/mtg/mtgpack"
	"github.com/pandodao/mtg/protocol"
	"github.com/pandodao/mtg/protocol/checksum"
	"github.com/shopspring/decimal"
)

const (
	ActionAdd    uint16 = 1
	ActionRemove uint16 = 2
	ActionSwap   uint16 = 3

	ProtocolVersion uint8 = 2
)

func BuildAdd(followID string, oppositeAsset string, slippage decimal.Decimal, expireDuration time.Duration) string {
	return buildAction(
		ActionAdd,
		followID,
		uuid.MustParse(oppositeAsset),
		slippage,
		int16(expireDuration.Seconds()),
	)
}

func BuildRemove(followID string) string {
	return buildAction(
		ActionRemove,
		followID,
	)
}

func BuildSwap(followID string, fillAsset string, paths route.Paths, minAmount decimal.Decimal) string {
	return buildAction(
		ActionSwap,
		followID,
		uuid.MustParse(fillAsset),
		paths,
		minAmount,
	)
}

func buildAction(action uint16, followID string, args ...interface{}) string {
	h := protocol.Header{
		Version:    ProtocolVersion,
		ProtocolID: protocol.ProtocolFswap,
		Action:     action,
	}

	if followID != "" {
		h.FollowID = uuid.MustParse(followID)
	}

	enc := mtgpack.NewEncoder()
	if err := enc.EncodeValue(h); err != nil {
		panic(err)
	}

	if err := enc.EncodeValues(args...); err != nil {
		panic(err)
	}

	b := enc.Bytes()
	sum := checksum.Sha256(b)
	return base64.StdEncoding.EncodeToString(append(b, sum...))
}
