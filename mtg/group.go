package fswap

import (
	"context"

	"github.com/fox-one/mixin-sdk-go"
)

type (
	Group struct {
		Members   []string  `json:"members,omitempty"`
		Threshold uint      `json:"threshold,omitempty"`
		PublicKey mixin.Key `json:"public_key"`
	}
)

func ReadGroup(ctx context.Context) (*Group, error) {
	const uri = "/api/info"
	resp, err := Request(ctx).Get(uri)
	if err != nil {
		return nil, err
	}

	var group Group
	if err := UnmarshalResponse(resp, &group); err != nil {
		return nil, err
	}

	return &group, err
}
