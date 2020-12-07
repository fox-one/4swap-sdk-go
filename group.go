package fswap

import (
	"context"
	"crypto/ed25519"
)

type (
	Group struct {
		Members   []string          `json:"members,omitempty"`
		Threshold uint              `json:"threshold,omitempty"`
		PublicKey ed25519.PublicKey `json:"public_key"`
	}
)

// ReadGroup return mtg Group info ( MTG only)
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
