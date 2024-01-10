package fswap

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	endpoint = "https://safe-swap-api.pando.im"
)

type Client struct {
	r *resty.Client
}

func New() *Client {
	r := resty.New().
		SetHeader("Content-Type", "application/json").
		SetBaseURL(endpoint).
		SetTimeout(10 * time.Second)

	return &Client{
		r: r,
	}
}

func (c *Client) Resty() *resty.Client {
	return c.r
}

func (c *Client) UseToken(token string) {
	c.r.SetAuthToken(token)
}

func (c *Client) UseEndpoint(endpoint string) {
	c.r.SetBaseURL(endpoint)
}

func (c *Client) request(ctx context.Context) *resty.Request {
	return c.r.R().SetContext(ctx)
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
