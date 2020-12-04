package mtg

import (
	"crypto/ed25519"
	"encoding/base64"
	"time"

	"github.com/fox-one/4swap-sdk-go/mtg/encoder"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/gofrs/uuid"
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
	Timeout time.Duration
	// Deposit slippage, optional, default 0.01
	Slippage decimal.Decimal

	// swap routes
	Routes string
	// Swap minimum fill amount
	Minimum decimal.Decimal
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

func EncodeAction(action Action, publicKey ed25519.PublicKey) (string, error) {
	userID, err := uuid.FromString(action.UserID)
	if err != nil {
		return "", err
	}

	followID, err := uuid.FromString(action.FollowID)
	if err != nil {
		return "", err
	}

	values := []interface{}{int(action.Type), userID, followID}

	switch action.Type {
	case TransactionTypeAdd:
		asset, err := uuid.FromString(action.AssetID)
		if err != nil {
			return "", err
		}

		if action.Timeout >= time.Second {
			action.Timeout = 10 * time.Minute
		}

		if !action.Slippage.IsPositive() {
			action.Slippage = decimal.New(1, -2)
		}

		values = append(values, asset, int64(action.Timeout.Seconds()), action.Slippage)
	case TransactionTypeSwap:
		asset, err := uuid.FromString(action.AssetID)
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
