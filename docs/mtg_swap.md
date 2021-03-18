# 4swap mtg 兑换指南

4swap mtg 兑换步骤：

1. 获取多签组收款地址
2. 下单转账给多签收款地址
3. 查询订单，收款

## 获取 4swap mtg 多签收款地址

> 多签组信息不会频繁变动，建议保存在配置文件
> [go sdk](https://github.com/fox-one/4swap-sdk-go/blob/master/group.go#L17)

```http request
GET /api/info
```

**Response:**

> members 是多签组成员 mixin id 列表；
> threshold 是多签组最少签名数

```json5
{
  "ts": 1616045842815,
  "data": {
    "members": [
      "a753e0eb-3010-4c4a-a7b2-a7bda4063f62",
      "099627f8-4031-42e3-a846-006ee598c56e",
      "aefbfd62-727d-4424-89db-ae41f75d2e04",
      "d68ca71f-0e2c-458a-bb9c-1d6c2eed2497",
      "e4bc0740-f8fe-418c-ae1b-32d9926f5863"
    ],
    "public_key": "dt351xp3KjNlVCMqBYUeUSF45upCEiReSZAqcjcP/Lc=",
    "threshold": 3
  }
}
```

## 构造下单 memo

可以直接使用 ```4swap-go-sdk``` [mtg.SwapAction](https://github.com/fox-one/4swap-sdk-go/blob/master/mtg/action.go#L48) 方法生成

```http request
POST /api/actions
```

**Body:**

```json5
{
  "action": "{user_id},{follow_id},3,{asset_id},{route_id},{minimum_fill}"
}
```

### body.action 介绍

1. user_id 是收款用户的 mixin_id (uuid)
2. follow_id 是查询订单用的自定义 id (uuid)
3. asset_id 是要买的币的 asset id
4. route_id 是自定义路径 id，为空的话引擎会自动选择最佳路径
5. minimum_fill 是做少买入量，如果因为深度变化导致无法买入至少这个数量的币，则兑换失败退款


**Response:**

```json5
{
  "follow_id": "follow id", // 和 body.action 一致
  "action": "memo", // 下单转账需要的 memo
}
```


## 转账下单

使用 mixin api ```POST /transactions``` 付款给 4swap mtg 多签，付款的币和数量就是下单的币和数量，memo 为上面创建的 memo

> ```POST /transactions``` 的使用参考 [mixin-sdk-go](https://github.com/fox-one/mixin-sdk-go/blob/faab649ffba80acf12948d5bb2205e149d5ace7b/transaction_raw.go#L41) 


## 查询订单

> 需要在 Header 带上 Authorization:Bearer token，token 为订单收款人的 Mixin Authorization Token，生成方式参考
> [4swap-go-sdk auth.go](https://github.com/fox-one/4swap-sdk-go/blob/master/auth.go#L12)

```http request
GET /api/orders/{follow_id}
```

**Response:**

```json5
{
  "data": {
    "id": "87ae5014-d20f-4cf1-b530-8771137e4e0e",
    "created_at": "2020-09-15T03:35:34Z",
    "user_id": "8017d200-7870-4b82-b53f-74bae1d2dad7",
    "state": "Done", // 订单状态 Trading/Rejected/Done
    "pay_asset_id": "6cfe566e-4aad-470b-8c9a-2fd35b49c68d",
    "fill_asset_id": "c6d0c728-2624-429b-8e0d-d9d19b6592fa",
    "pay_amount": "1", // 付款的币的数量
    "funds": "1", // deprecated, same as pay_amount
    "fill_amount": "00025725", // 买到的币的数量
    "amount": "0.00025725", // deprecated, same as fill_amount
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
