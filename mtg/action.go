package mtg

import (
	"crypto/ed25519"
	"encoding/base64"
	"time"

	"github.com/fox-one/4swap-sdk-go/v2/mtg/encoder"
	"github.com/fox-one/mixin-sdk-go/v2"
	"github.com/google/uuid"
	"github.com/pandodao/mtg/mtgpack"
	"github.com/pandodao/mtg/protocol"
	"github.com/shopspring/decimal"
)

type Action struct {
	// action type, Add, Remove, Swap
	Type TransactionType
	// user mixin id
	UserID string
	// action trace id
	FollowID string

	// AssetID is pair quote asset id if base asset will be paid, otherwise this is base asset id.
	// Ignore if type is Remove
	AssetID string

	// Deposit Timeout, optional, default is 10m
	Timeout time.Duration `json:"Timeout,omitempty"`
	// Deposit slippage, optional, default 0.01
	Slippage decimal.Decimal `json:"Slippage,omitempty"`

	// swap routes
	Routes string `json:"Routes,omitempty"`
	// Swap minimum fill amount
	Minimum decimal.Decimal `json:"Minimum,omitempty"`
}

func AddAction(userID, followID, assetID string, timeout time.Duration, slippage decimal.Decimal) Action {
	return Action{
		Type:     TransactionTypeAdd,
		UserID:   userID,
		FollowID: followID,
		AssetID:  assetID,
		Timeout:  timeout,
		Slippage: slippage,
	}
}

func SwapAction(userID, followID, assetID string, routes string, min decimal.Decimal) Action {
	return Action{
		Type:     TransactionTypeSwap,
		UserID:   userID,
		FollowID: followID,
		AssetID:  assetID,
		Routes:   routes,
		Minimum:  min,
	}
}

func RemoveAction(userID, followID string) Action {
	return Action{
		Type:     TransactionTypeRemove,
		UserID:   userID,
		FollowID: followID,
	}
}

func (action Action) Encode(publicKey ed25519.PublicKey) (string, error) {
	return EncodeAction(action, publicKey)
}

// EncodeAction encode action to 4swap memo
// deprecated, use EncodeActionV1 instead
func EncodeAction(action Action, publicKey ed25519.PublicKey) (string, error) {
	userID, err := uuid.Parse(action.UserID)
	if err != nil {
		return "", err
	}

	followID, err := uuid.Parse(action.FollowID)
	if err != nil {
		return "", err
	}

	values := []interface{}{int(action.Type), userID, followID}

	switch action.Type {
	case TransactionTypeAdd:
		asset, err := uuid.Parse(action.AssetID)
		if err != nil {
			return "", err
		}

		if action.Timeout < time.Second {
			action.Timeout = 10 * time.Minute
		}

		if !action.Slippage.IsPositive() {
			action.Slippage = decimal.New(1, -2)
		}

		values = append(values, asset, action.Slippage, int64(action.Timeout.Seconds()))
	case TransactionTypeSwap:
		asset, err := uuid.Parse(action.AssetID)
		if err != nil {
			return "", err
		}

		values = append(values, asset, action.Routes, action.Minimum)
	}

	return encodeMemo(publicKey, values...)
}

func encodeMemo(pub ed25519.PublicKey, values ...interface{}) (string, error) {
	action, err := encoder.Encode(values...)
	if err != nil {
		return "", err
	}

	key := mixin.GenerateEd25519Key()
	action, err = encoder.Encrypt(action, key, pub)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(action), nil
}

// EncodeActionV1 encode action to 4swap memo v1
func EncodeActionV1(action Action) (string, error) {
	enc := mtgpack.NewEncoder()

	h := protocol.Header{
		Version:    1,
		ProtocolID: protocol.ProtocolFswap,
		Action:     uint16(action.Type),
	}

	h.FollowID, _ = uuid.Parse(action.FollowID)

	if err := enc.EncodeValue(h); err != nil {
		return "", err
	}

	var r protocol.MultisigReceiver
	if user, err := uuid.Parse(action.UserID); err == nil {
		r.Members = append(r.Members, user)
		r.Threshold = 1
	}

	if err := enc.EncodeValue(r); err != nil {
		return "", err
	}

	switch action.Type {
	case TransactionTypeAdd:
		asset, err := uuid.Parse(action.AssetID)
		if err != nil {
			return "", err
		}

		if action.Timeout < time.Second {
			action.Timeout = 10 * time.Minute
		}

		if !action.Slippage.IsPositive() {
			action.Slippage = decimal.New(1, -2)
		}

		if err := enc.EncodeValues(asset, action.Slippage, int16(action.Timeout.Seconds())); err != nil {
			return "", err
		}
	case TransactionTypeSwap:
		asset, err := uuid.Parse(action.AssetID)
		if err != nil {
			return "", err
		}

		if err := enc.EncodeValues(asset, action.Routes, action.Minimum); err != nil {
			return "", err
		}
	}

	return base64.StdEncoding.EncodeToString(enc.Bytes()), nil
}
