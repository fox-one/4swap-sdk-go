# 4swap 程序交易接入文档

1. api base 是 **https://f1-uniswap-api.firesbox.com/api**

2. 4swap 机器人的 mixin id 是 a753e0eb-3010-4c4a-a7b2-a7bda4063f62

3. 需要授权的接口，请在 header Authorization 带上 [mixin authorization token](https://developers.mixin.one/api/a-beginning/authentication-token)，要求签名的 url 是 **/me**，然后 scope 是 **FULL**。

## 下单 Swap 转账

转账 memo 需要带上要买的币，以及路径和最小买入数量

```json5
{
    "t": "Add", // action type
    "a": "66152c0b-3355-38ef-9ec5-cae97e29472a", // fill asset id,
    "r": "xxx", // route id ，可以由 pre order 得到，不传的话则由引擎自动选择最优路径
    "m": "0.001" // 最小买入量，不传的话则不限制
}
```

将上述 payload 先 json encode 或者 msgpack encode 然后再 base64 encode 得到转账 memo

## API

### 读取交易对详情

```http request
GET /pairs/{base_asset_id}/{quote_asset_id}
```

> url path 里面 asset id 的顺序不影响结果

**Response:**

```json5
{
  "data": {
    "base_asset_id": "66152c0b-3355-38ef-9ec5-cae97e29472a",
    "quote_asset_id": "965e5c6e-434c-3fa9-b780-c50f43cd955c",
    "base_amount": "1714432671.4860945", // 流动池里面 base asset 的数量
    "quote_amount": "1428792957.89508717", // 流动池里面 quote asset 的数量
    "fee_percent": "0.003", // 手续费比例
    "liquidity": "1441664019.41091031" // 总的流动性份额
  }
}
```

### 读取全部交易对

```http request
GET /pairs
```

**Response:**

```json5
{
  "data": {
    "pairs": [
      {
        "base_asset_id": "965e5c6e-434c-3fa9-b780-c50f43cd955c",
        "quote_asset_id": "f5ef6b5d-cc5a-3d90-b2c0-a2fd386e7a3c",
        "base_amount": "21416602.86502364",
        "quote_amount": "0.33923379",
        "fee_percent": "0.003",
        "liquidity": "2695.28124196"
      },
      {
        "base_asset_id": "3edb734c-6d6f-32ff-ab03-4eb43640c758",
        "quote_asset_id": "965e5c6e-434c-3fa9-b780-c50f43cd955c",
        "base_amount": "18.67863207",
        "quote_amount": "426.53119598",
        "fee_percent": "0.003",
        "liquidity": "89.25629191"
      },
      {
        "base_asset_id": "4d8c508b-91c5-375b-92b0-ee702ed2dac5",
        "quote_asset_id": "965e5c6e-434c-3fa9-b780-c50f43cd955c",
        "base_amount": "0.12647595",
        "quote_amount": "89167738.70230537",
        "fee_percent": "0.003",
        "liquidity": "3354.1133805"
      }
    ]
  }
}
```

### 读取全部 assets

```http request
GET /assets
```

**Response:**

```json5
{
  "data": {
    "assets": [
      {
        "id": "c6d0c728-2624-429b-8e0d-d9d19b6592fa",
        "name": "Bitcoin",
        "symbol": "BTC",
        "logo": "https://mixin-images.zeromesh.net/HvYGJsV5TGeZ-X9Ek3FEQohQZ3fE9LBEBGcOcn4c4BNHovP4fW4YB97Dg5LcXoQ1hUjMEgjbl1DPlKg1TW7kK6XP=s128",
        "chain_id": "c6d0c728-2624-429b-8e0d-d9d19b6592fa",
        "chain": {
          "id": "c6d0c728-2624-429b-8e0d-d9d19b6592fa",
          "name": "Bitcoin",
          "symbol": "BTC",
          "logo": "https://mixin-images.zeromesh.net/HvYGJsV5TGeZ-X9Ek3FEQohQZ3fE9LBEBGcOcn4c4BNHovP4fW4YB97Dg5LcXoQ1hUjMEgjbl1DPlKg1TW7kK6XP=s128",
          "chain_id": "c6d0c728-2624-429b-8e0d-d9d19b6592fa",
          "price": "10741.16"
        },
        "price": "10741.16"
      },
      {
        "id": "43d61dcd-e413-450d-80b8-101d5e903357",
        "name": "Ether",
        "symbol": "ETH",
        "logo": "https://mixin-images.zeromesh.net/zVDjOxNTQvVsA8h2B4ZVxuHoCF3DJszufYKWpd9duXUSbSapoZadC7_13cnWBqg0EmwmRcKGbJaUpA8wFfpgZA=s128",
        "chain_id": "43d61dcd-e413-450d-80b8-101d5e903357",
        "chain": {
          "id": "43d61dcd-e413-450d-80b8-101d5e903357",
          "name": "Ether",
          "symbol": "ETH",
          "logo": "https://mixin-images.zeromesh.net/zVDjOxNTQvVsA8h2B4ZVxuHoCF3DJszufYKWpd9duXUSbSapoZadC7_13cnWBqg0EmwmRcKGbJaUpA8wFfpgZA=s128",
          "chain_id": "43d61dcd-e413-450d-80b8-101d5e903357",
          "price": "375.23"
        },
        "price": "375.23"
      },
      {
        "id": "fd11b6e3-0b87-41f1-a41f-f0e9b49e5bf0",
        "name": "Bitcoin Cash",
        "symbol": "BCH",
        "logo": "https://mixin-images.zeromesh.net/tqt14x8iwkiCR_vIKIw6gAAVO8XpZH7ku7ZJYB5ArMRA6grN9M1oCI7kKt2QqBODJwr17sZxDCDTjXHOgIixzv6X=s128",
        "chain_id": "fd11b6e3-0b87-41f1-a41f-f0e9b49e5bf0",
        "chain": {
          "id": "fd11b6e3-0b87-41f1-a41f-f0e9b49e5bf0",
          "name": "Bitcoin Cash",
          "symbol": "BCH",
          "logo": "https://mixin-images.zeromesh.net/tqt14x8iwkiCR_vIKIw6gAAVO8XpZH7ku7ZJYB5ArMRA6grN9M1oCI7kKt2QqBODJwr17sZxDCDTjXHOgIixzv6X=s128",
          "chain_id": "fd11b6e3-0b87-41f1-a41f-f0e9b49e5bf0",
          "price": "238.83"
        },
        "price": "238.83"
      },
      {
        "id": "574388fd-b93f-4034-a682-01c2bc095d17",
        "name": "Bitcoin SV",
        "symbol": "BSV",
        "logo": "https://mixin-images.zeromesh.net/1iUl5doLjMSv-ElcVCI4YgD1uIayDbZcQP0WjFEajoY1-qQZmVEl5GgUCtsp8CP0aj96a5Rwi-weQ5YA64lyQzU=s128",
        "chain_id": "574388fd-b93f-4034-a682-01c2bc095d17",
        "chain": {
          "id": "574388fd-b93f-4034-a682-01c2bc095d17",
          "name": "Bitcoin SV",
          "symbol": "BSV",
          "logo": "https://mixin-images.zeromesh.net/1iUl5doLjMSv-ElcVCI4YgD1uIayDbZcQP0WjFEajoY1-qQZmVEl5GgUCtsp8CP0aj96a5Rwi-weQ5YA64lyQzU=s128",
          "chain_id": "574388fd-b93f-4034-a682-01c2bc095d17",
          "price": "167.78"
        },
        "price": "167.78"
      }
    ]
  }
}
```

### 注入流动性 (Authorization Token Required)

```http request
POST /pairs/{base_asset_id}/{quote_asset_id}/deposit
```

> 两种币的数量尽量按照池子流动性比例注入。
> api 会返回两个 transfer，包含了付款需要的所有信息

**Params:**

```json5
{
    "base_amount":"1", // 注入 base asset 的数量
    "quote_amount":"1.00185986" // 注入 quote asset 的数量
}
```

**Response:**

```json5
{
  "data": {
    "id": "9d6ac652-fcb8-42b5-a248-de70ad78b660",
    "created_at": "2020-09-15T03:20:43.047872131Z",
    "state": "Pending",
    "user_id": "8017d200-7870-4b82-b53f-74bae1d2dad7",
    "transfers": [
      {
        "trace_id": "34c41b49-5a13-5ccd-8a9c-a2906ffac571",
        "opponent_id": "a753e0eb-3010-4c4a-a7b2-a7bda4063f62",
        "asset_id": "4d8c508b-91c5-375b-92b0-ee702ed2dac5",
        "amount": "1",
        "memo": "eyJ0IjoiQWRkIiwiZCI6IjlkNmFjNjUyLWZjYjgtNDJiNS1hMjQ4LWRlNzBhZDc4YjY2MCJ9"
      },
      {
        "trace_id": "430886c6-e980-53b6-9c45-a306319f6908",
        "opponent_id": "a753e0eb-3010-4c4a-a7b2-a7bda4063f62",
        "asset_id": "815b0b1a-2764-3736-8faa-42d694fa620a",
        "amount": "1.00185986",
        "memo": "eyJ0IjoiQWRkIiwiZCI6IjlkNmFjNjUyLWZjYjgtNDJiNS1hMjQ4LWRlNzBhZDc4YjY2MCJ9"
      }
    ]
  }
}
```

### Swap 预下单

```http request
POST /orders/pre
```

**Params:**

```json5
{
    "pay_asset_id": "c6d0c728-2624-429b-8e0d-d9d19b6592fa",
    "fill_asset_id": "6cfe566e-4aad-470b-8c9a-2fd35b49c68d",
    "funds": "1", // 下单的币的数量，和 amount 二选一
    "amount": "0.01" // 要买的币的数量，和 funds 二选一
}
```

**Response:**

```json5
{
  "data": {
    "created_at": "0001-01-01T00:00:00Z",
    "state": "Done",
    "pay_asset_id": "f5ef6b5d-cc5a-3d90-b2c0-a2fd386e7a3c",
    "fill_asset_id": "4d8c508b-91c5-375b-92b0-ee702ed2dac5",
    "funds": "990",
    "amount": "0.02442823", // 预估能买到的币的数量
    "min_amount": "0",
    "routes": "d6TR", // route id，代表 route 路径
    "route_assets": [
      "f5ef6b5d-cc5a-3d90-b2c0-a2fd386e7a3c",
      "965e5c6e-434c-3fa9-b780-c50f43cd955c",
      "4d8c508b-91c5-375b-92b0-ee702ed2dac5"
    ],
    "transactions": [
      {
        "created_at": "0001-01-01T00:00:00Z",
        "type": "Swap",
        "base_asset_id": "965e5c6e-434c-3fa9-b780-c50f43cd955c",
        "quote_asset_id": "f5ef6b5d-cc5a-3d90-b2c0-a2fd386e7a3c",
        "base_amount": "-21409244.69027989",
        "quote_amount": "990",
        "fee_asset_id": "f5ef6b5d-cc5a-3d90-b2c0-a2fd386e7a3c",
        "fee_amount": "2.97",
        "pay_asset_id": "f5ef6b5d-cc5a-3d90-b2c0-a2fd386e7a3c",
        "filled_asset_id": "965e5c6e-434c-3fa9-b780-c50f43cd955c",
        "funds": "990",
        "amount": "21409244.69027989"
      },
      {
        "created_at": "0001-01-01T00:00:00Z",
        "type": "Swap",
        "base_asset_id": "4d8c508b-91c5-375b-92b0-ee702ed2dac5",
        "quote_asset_id": "965e5c6e-434c-3fa9-b780-c50f43cd955c",
        "base_amount": "-0.02442823",
        "quote_amount": "21409244.69027989",
        "fee_asset_id": "965e5c6e-434c-3fa9-b780-c50f43cd955c",
        "fee_amount": "64227.73407083",
        "pay_asset_id": "965e5c6e-434c-3fa9-b780-c50f43cd955c",
        "filled_asset_id": "4d8c508b-91c5-375b-92b0-ee702ed2dac5",
        "funds": "21409244.69027989",
        "amount": "0.02442823"
      }
    ],
    "route_price": "0.089547130029244",
    "price_impact": "0.999724"
  }
}
```

### 查询 Swap 订单 (Authorization Token Required)

```http request
GET /orders/{order_id}
```

> order id 即下单转账的 trace id

**Response:**

```json5
{
  "data": {
    "id": "87ae5014-d20f-4cf1-b530-8771137e4e0e",
    "created_at": "2020-09-15T03:35:34Z",
    "user_id": "8017d200-7870-4b82-b53f-74bae1d2dad7",
    "state": "Done", // 订单状态 Trading Rejected Done
    "pay_asset_id": "6cfe566e-4aad-470b-8c9a-2fd35b49c68d",
    "fill_asset_id": "c6d0c728-2624-429b-8e0d-d9d19b6592fa",
    "funds": "1", // 付款的币的数量
    "amount": "0.00025725", // 买到的币的数量
    "min_amount": "0.0002521",
    "routes": "1bv",
    "route_assets": [
      "6cfe566e-4aad-470b-8c9a-2fd35b49c68d",
      "c6d0c728-2624-429b-8e0d-d9d19b6592fa"
    ],
    "transactions": [
      {
        "id": "87ae5014-d20f-4cf1-b530-8771137e4e0e",
        "created_at": "2020-09-15T03:35:34Z",
        "user_id": "8017d200-7870-4b82-b53f-74bae1d2dad7",
        "type": "Swap",
        "base_asset_id": "6cfe566e-4aad-470b-8c9a-2fd35b49c68d",
        "quote_asset_id": "c6d0c728-2624-429b-8e0d-d9d19b6592fa",
        "base_amount": "1",
        "quote_amount": "-0.00025725",
        "fee_asset_id": "6cfe566e-4aad-470b-8c9a-2fd35b49c68d",
        "fee_amount": "0.003",
        "pay_asset_id": "6cfe566e-4aad-470b-8c9a-2fd35b49c68d",
        "filled_asset_id": "c6d0c728-2624-429b-8e0d-d9d19b6592fa",
        "funds": "1",
        "amount": "0.00025725"
      }
    ]
  }
}
```
