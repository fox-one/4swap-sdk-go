# 4swap SDK go

## 4swap

[4swap](https://4swap.org) is a decentralized protocol implement for automated liquidity provision on [Mixin Network](https://mixin.one)

### Authorization

4swap supports two kinds of access tokens:

1. the access token that complete the OAuth flow at 4swap's webpage: https://app.4swap.org
2. the access token that generated by your own Mixn Application. The token should sign the URL **/me** and the scope should be "FULL". Please read this [document](https://developers.mixin.one/api/a-beginning/authentication-token) for more details.

### Example

```golang
func TestMtgSwap(t *testing.T) {
    const (
        btc   = "c6d0c728-2624-429b-8e0d-d9d19b6592fa"
        xin   = "c94ac88f-4671-3976-b60a-09064f1811e8"
        token = "your authorization token"
    )
    
    ctx := context.Background()
    
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
        "", // leave the routes field as empty to let the engine decide the route. 
        decimal.NewFromFloat(0.1),
    )
    
    memo, err := action.Encode(group.PublicKey)
    if err != nil {
        t.Fatal(err)
    }
    
    t.Log(memo)
    
    // use mixin-sdk-go or bot-api-client-go to transfer to 4swap's multisig address
    
    // query the order.
    ctx = fswap.WithToken(ctx, token)
    order, err := fswap.ReadOrder(ctx, followID.String())
    if err != nil {
        t.Fatal(err)
    }
    
    t.Log(order.State)
}
```
