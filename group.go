package fswap

import (
	"context"
)

type Group struct {
	Members    []string `json:"members,omitempty"`
	Threshold  uint8    `json:"threshold,omitempty"`
	MixAddress string   `json:"mix_address,omitempty"`
}

// ReadGroup return mtg Group info ( MTG only)
func (c *Client) ReadGroup(ctx context.Context) (*Group, error) {
	const uri = "/api/info"
	resp, err := c.request(ctx).Get(uri)
	if err != nil {
		return nil, err
	}

	var group Group
	if err := UnmarshalResponse(resp, &group); err != nil {
		return nil, err
	}

	return &group, err
}
