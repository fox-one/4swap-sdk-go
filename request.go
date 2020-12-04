package fswap

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	Endpoint = "https://f1-uniswap-api.firesbox.com"
	ClientID = "a753e0eb-3010-4c4a-a7b2-a7bda4063f62"

	MtgEndpoint = "https://swap-mtg-test-api.fox.one"
)

var httpClient = resty.New().
	SetHeader("Content-Type", "application/json").
	SetHostURL(Endpoint).
	SetTimeout(10 * time.Second).
	SetPreRequestHook(func(c *resty.Client, req *http.Request) error {
		if token, ok := TokenFrom(req.Context()); ok {
			req.Header.Set("Authorization", "Bearer "+token)
		}

		return nil
	})

func UseEndpoint(endpoint string) {
	httpClient.SetHostURL(endpoint)
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
