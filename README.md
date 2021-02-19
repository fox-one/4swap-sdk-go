## 4swap

4swap (Mixin ID: 7000103537) 是 Fox.ONE 团队基于 Mixin Network 开发的 Swap 协议交易所，支持 Mixin 钱包里面任意币之间的兑换。

### Authorization

支持两种 Token

1. 4swap 机器人 Oauth 授权得到的 token
2. 机器人自己签出来的 **/me** 的 [Authentication Token](https://developers.mixin.one/api/a-beginning/authentication-token)，scp 必须是 **FULL**。

### Example

```golang
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
    
    // 使用 mixin-sdk-go 或者 bot-api-client-go 转给给 4swap 多签
    
    // 查询订单
    ctx = fswap.WithToken(ctx, token)
    order, err := fswap.ReadOrder(ctx, followID.String())
    if err != nil {
        t.Fatal(err)
    }
    
    t.Log(order.State)
}
```
