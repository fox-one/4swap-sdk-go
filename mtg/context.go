package fswap

import (
	"context"
)

type contextKey int

const (
	tokenKey contextKey = iota
)

func WithToken(ctx context.Context,token string) context.Context {
	return context.WithValue(ctx,tokenKey,token)
}

func TokenFrom(ctx context.Context) (string,bool) {
	v := ctx.Value(tokenKey)
	token,ok := v.(string)
	return token,ok
}
