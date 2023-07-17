package fswap

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/fox-one/4swap-sdk-go/legacy"
	"github.com/go-resty/resty/v2"
)

const (
	LegacyEndpoint = legacy.Endpoint
	ClientID       = legacy.ClientID

	MtgEndpoint = "https://api.4swap.org"
)

var httpClient = resty.New().
	SetHeader("Content-Type", "application/json").
	SetBaseURL(MtgEndpoint).
	SetTimeout(10 * time.Second).
	SetPreRequestHook(func(c *resty.Client, req *http.Request) error {
		if token, ok := TokenFrom(req.Context()); ok {
			req.Header.Set("Authorization", "Bearer "+token)
		}

		return nil
	})

func GetClient() *http.Client {
	return httpClient.GetClient()
}

func GetRestyClient() *resty.Client {
	return httpClient
}

func UseEndpoint(endpoint string) {
	httpClient.SetBaseURL(endpoint)
}

func Request(ctx context.Context) *resty.Request {
	return httpClient.R().SetContext(ctx)
}

func DecodeResponse(resp *resty.Response) ([]byte, error) {
	var body struct {
		Error
		Data json.RawMessage `json:"data,omitempty"`
	}

	if err := json.Unmarshal(resp.Body(), &body); err != nil {
		if resp.IsError() {
			return nil, &Error{
				Code: resp.StatusCode(),
				Msg:  resp.Status(),
			}
		}

		return nil, err
	}

	if body.Error.Code > 0 {
		return nil, &body.Error
	}

	return body.Data, nil
}

func UnmarshalResponse(resp *resty.Response, v interface{}) error {
	data, err := DecodeResponse(resp)
	if err != nil {
		return err
	}

	if v != nil {
		return json.Unmarshal(data, v)
	}

	return nil
}
