# Getting Started



## Import 4swap SDK in your project

```golang
import (
	fswap "github.com/fox-one/4swap-sdk-go"
	mtg "github.com/fox-one/4swap-sdk-go/mtg"
  "github.com/fox-one/mixin-sdk-go"
	"github.com/shopspring/decimal"
)
```

## Get the multisig group information

4swap is a decentralized application based on [MTG](https://developers.mixin.one/document/mainnet/mtg).

When you perform essential operations, such as swapping crypto, adding liquidity, removing liquidity, you must create multisig transactions and interact with Mixin Network.

The participants of each multisig are also the members of MTG. So please read them first and save them for later using.

```golang
  ctx := context.TODO()

  // use the 4swap's MTG api endpoint
  fswap.UseEndpoint(fswap.MtgEndpoint)

  // read the mtg group
  // the group information would change frequently
  // it's recommended to save it for later using
  group, err := fswap.ReadGroup(ctx)
  if err != nil {
    return err
  }
  ...
```

## Get all tradable pairs

To get all supported by 4swap is easy:

```golang
	pairs, err := fswap.ListPairs(ctx)
	if err != nil {
		return err
	}
  ...
```

## Calculate the best route to trade

Before swapping crypto, we need to specify the swapping route.

At present, you may let's 4swap to decide the route, but it has the performance issues and may slow down the bot. Because of that, it is recommended to calculate the a swapping route yourself.

To calculate route is easy. Sort the pairs according to the liquidity and call `Route` or `ReverseRoute` methods, which will return an `order` object that includes the result of calculation.

```golang
	// sort first
	sort.Slice(pairs, func(i, j int) bool {
		aLiquidity := pairs[i].BaseValue.Add(pairs[i].QuoteValue)
		bLiquidity := pairs[j].BaseValue.Add(pairs[j].QuoteValue)
		return aLiquidity.GreaterThan(bLiquidity)
	})

  // calculate the route
  // InputAssetID - the id of the asset you want to paid
  // OutputAssetID - the id of the asset you trade for
  // InputAmount - the amount to calucate the route, for example, 1000
	preOrder, err := fswap.Route(pairs, InputAssetID, OutputAssetID, InputAmount)
	if err != nil {
		return err
	}

  // you can read the best route from Order.RouteAssets, which is an array of asset_id
  log.Printf("Route: %v", preOrder.RouteAssets)
	log.Printf("Price: %v", preOrder.FillAmount.Div(InputAmount))
  ...
```

## Construct a real order

All required information about an order are store in the transaction memo, in JSON format:

```json
{
  "action": "1,{receiver_id},{follow_id},{asset_id},{slippage},{timeout}"
}
```

in which,

- `{receiver_id}` is the id of user who will receive the crypto from swapping
- `{follow_id}` is a UUID to trace the order
- `{asset_id}` is the asset's ID you are swapping for
- `{slippage}` is the slippage ratio, e.g. 0.01 = 1%
- `{timeout}` is the timeout in sec

If you are using 4swap SDK, you can also use the method `mtg.SwapAction` to simplify the process:

```golang
    // the ID to trace the orders at 4swap
    followID, _ := uuid.NewV4()

    // build a swap action, specified the parameters
    action := mtg.SwapAction(
        receiverID,
        followID.String(),
        OutputAssetID,
        preOrder.Routes,
        // the minimum amount of asset you will get.
        // you may want to change this value to a number which less than preOrder.FillAmount
        preOrder.FillAmount.Div(decimal.NewFromFloat(0.005)),
    )

    // generate the memo
    memo, err := action.Encode(group.PublicKey)
    if err != nil {
        return err
    }
    log.Println("memo", memo)
    ...

```

## Place an order programmatically

If you want the bot to place an order, send a multisig transaction from the bot.

This is a common scene for arbitrage bot. Please make sure the bot have enough crypto in the bot's wallet.

We need to use [mixin-sdk-go](https://github.com/fox-one/mixin-sdk-go) client to create and send the transaction to the kenerl nodes.

```golang
    // send a transaction to a multi-sign address which specified by `OpponentMultisig`
    // the OpponentMultisig.Receivers are the MTG group members of 4swap
    tx, err := client.Transaction(ctx, &mixin.TransferInput{
        AssetID: payAssetID,
        Amount:  decimal.RequireFromString(amount),
        TraceID: mixin.RandomTraceID(),
        Memo:    memo,
        OpponentMultisig: struct {
            Receivers []string `json:"receivers,omitempty"`
            Threshold uint8    `json:"threshold,omitempty"`
        }{
            Receivers: group.Members,
            Threshold: uint8(group.Threshold),
        },
    }, *pin)
```

## Place an order via Mixin Messenger

If you want to place an order via Mixin Messenger, generate a payment scheme to invoke Messenger from the webview.

This is a common scene for a webapp which provide swapping service to users, like [4swap's webpage](https://app.4swap.org).

We need to post `https://api.mixin.one/payments` to get a payment object which contains `code_id` to create the scheme:

```typescript
  function getPayments(asset_id, amount, memo, receivers, threshold): Promise<any> {
    const params = {
      asset_id,
      amount,
      memo,
      trace_id: uuid(),
      opponent_multisig: { receivers, threshold },
    };
    // use your http request lib here
    return http.post("/payments", { data: params });
  }

  ...

  const resp = await getPayments(
    asset_id,  // the input asset id
    value,     // the input amount
    memo,      // create by `SwapAction`
    members,   // read from the mulitsig group
    threshold, // read from the mulitsig group
  );

  window.location.href = `https://mixin.one/codes/${resp.code_id}`;
```


