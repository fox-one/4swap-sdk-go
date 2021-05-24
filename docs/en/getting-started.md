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

```golang
	pairs, err := fswap.ListPairs(ctx)
	if err != nil {
		return err
	}
  ...
```

## Calculate the best route to trade

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

```golang
    // the ID to trace the orders at 4swap
    followID, _ := uuid.NewV4()

    // build a swap action, specified the swapping parameters
    action := mtg.SwapAction(
        // the user ID to receive the money
        receiverID,
        // an UUID get trace the order
        followID.String(),
        // the asset's ID you are swapping for.
        OutputAssetID,
        // use `order.routes` to specified a route
        preOrder.Routes,
        // the minimum amount of asset you will get.
        // you may want to change this value to a specified number which less than preOrder.FillAmount
        preOrder.FillAmount.Div(decimal.NewFromFloat(0.005)),
    )

    // the action will be sent to 4swap in the memo
    memo, err := action.Encode(group.PublicKey)
    if err != nil {
        return err
    }
    log.Println("memo", memo)
    ...

```

## place a real order programmatically

If you have enough crypto in the bot's wallet, you can let the bot place an order.

A multisig transaction should be create and send to the kenerl nodes.

We need to use [mixin-sdk-go](https://github.com/fox-one/mixin-sdk-go) client to send the transaction.

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
