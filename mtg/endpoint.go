package fswap

import (
	"context"

	"github.com/fox-one/mixin-sdk-go"
)

const (
	defaultEndpoint = "https://swap-mtg-test-api.fox.one"
)

var (
	group *Group = defaultGroup()
)

func UseEndpoint(ctx context.Context, endpoint string) error {
	httpClient.HostURL = endpoint

	g, err := ReadGroup(ctx)
	if err != nil {
		return err
	}
	group = g
	return nil
}

func defaultGroup() *Group {
	group := Group{
		Members: []string{
			"9656eacd-2fa7-4e7b-b0eb-c475c9964f78",
			"ab14736f-e595-4e65-9879-871819d390f5",
			"b856deb3-e92f-4c19-9733-ec43526f95ce",
			"229fc7ac-9d09-4a6a-af5a-78f7439dce76",
			"84a4db41-4992-4d35-aac7-987f965f0302",
		},
		Threshold: 4,
	}
	group.PublicKey, _ = mixin.KeyFromString("WE2b3mzyi23SiEKEiiHy6+72LVUG9gDSEJ0d1jU+yC0=")
	return &group
}

