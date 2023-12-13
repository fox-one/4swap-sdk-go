package fswap_test

import (
	"context"
	"log"
	"testing"

	fswap "github.com/fox-one/4swap-sdk-go/v2"
	"github.com/fox-one/4swap-sdk-go/v2/mtg"
	"github.com/fox-one/mixin-sdk-go/v2"
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

func SimpleExample() {
	const (
		btc   = "c6d0c728-2624-429b-8e0d-d9d19b6592fa"
		xin   = "c94ac88f-4671-3976-b60a-09064f1811e8"
		token = "your authorization token"
	)

	ctx := fswap.WithToken(context.Background(), token)

	// 获取交易对信息，拿到流动性
	pair, err := fswap.ReadPair(ctx, btc, xin)
	if err != nil {
		log.Panicf("read pair failed: %s", err)
	}

	// 按比例添加流动性
	baseAmount := decimal.NewFromInt(1)
	quoteAmount := pair.QuoteAmount.Div(pair.BaseAmount).Mul(baseAmount).Truncate(8)
	req := &fswap.AddDepositReq{
		BaseAssetID:  pair.BaseAssetID,
		BaseAmount:   baseAmount,
		QuoteAssetID: pair.QuoteAssetID,
		QuoteAmount:  quoteAmount,
		// 滑点设置，如果为 0 则不限制
		Slippage: decimal.NewFromFloat(0.01),
	}

	deposit, err := fswap.AddDeposit(ctx, req)
	if err != nil {
		log.Panicf("request add deposit failed: %s", err)
	}

	// 需要两笔全部转账才能完成添加流动性
	for _, transfer := range deposit.Transfers {
		log.Println("handle transfer", transfer.TraceID)
	}

	// 付款 btc 兑换 xin
	payAmount := decimal.NewFromFloat(0.01)

	// 先 pre order 看看
	pre, err := fswap.PreOrder(ctx, &fswap.PreOrderReq{
		PayAssetID:  btc,
		FillAssetID: xin,
		Funds:       payAmount,
	})
	if err != nil {
		log.Panicf("pre order failed: %s", err)
	}

	log.Printf("fill amount %s", pre.FillAmount)

	// 准备转账下单
	action := fswap.TransactionAction{
		Type:   fswap.TransactionTypeSwap,
		Routes: pre.Routes,
		// 最小买入量设置成预估买入量的 98%
		// 如果想加大兑换成功率，这个滑点可以设置得低一点
		Minimum: pre.FillAmount.Mul(decimal.NewFromFloat(0.98)).Truncate(8).String(),
	}

	memo := fswap.EncodeAction(action)
	transfer := &fswap.TransferReq{
		AssetID:    btc,
		Amount:     payAmount,
		TraceID:    "new uuid v4",
		Memo:       memo,
		OpponentID: fswap.ClientID,
	}

	log.Println("handle swap transfer", transfer.TraceID)
}

func TestMtgSwap(t *testing.T) {
	const (
		btc   = "c6d0c728-2624-429b-8e0d-d9d19b6592fa"
		xin   = "c94ac88f-4671-3976-b60a-09064f1811e8"
		token = "your authorization token"
	)

	ctx := context.Background()
	fswap.UseEndpoint(fswap.MtgEndpoint)

	group, err := fswap.ReadGroup(ctx)
	if err != nil {
		t.Fatal(err)
	}

	me, err := mixin.UserMe(ctx, token)
	if err != nil {
		t.Fatal(err)
	}

	followID, _ := uuid.NewV4()
	action := mtg.SwapAction(
		me.UserID,
		followID.String(),
		btc,
		"", // routes 为空则不指定
		decimal.NewFromFloat(0.1),
	)

	memo, err := action.Encode(group.PublicKey)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(memo)

	// 使用 mixin-sdk-go/v2 或者 bot-api-client-go 转给给 4swap 多签

	// 查询订单
	ctx = fswap.WithToken(ctx, token)
	order, err := fswap.ReadOrder(ctx, followID.String())
	if err != nil {
		t.Fatal(err)
	}

	t.Log(order.State)
}
