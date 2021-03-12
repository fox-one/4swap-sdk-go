* endpoint is **https://api.4swap.org/**

##  List cmc Pairs

```http request
GET /api/cmc/pairs
```

**Response:**

```json5
{
  "ts": 1610705153506,
  "data": {
    "43d61dcd-e413-450d-80b8-101d5e903357_c94ac88f-4671-3976-b60a-09064f1811e8": {
      "base_id": "43d61dcd-e413-450d-80b8-101d5e903357", // mixin asset id
      "base_name": "Ether",
      "base_symbol": "ETH",
      "quote_id": "c94ac88f-4671-3976-b60a-09064f1811e8", // quote asset id
      "quote_name": "Mixin",
      "quote_symbol": "XIN",
      "last_price": "7.2928686298372538",
      "base_volume": "0.20031361",  // base volume in 24h
      "quote_volume": "1.40632248"  // quote volume in 24h
    },
    "4d8c508b-91c5-375b-92b0-ee702ed2dac5_c94ac88f-4671-3976-b60a-09064f1811e8": {
      "base_id": "4d8c508b-91c5-375b-92b0-ee702ed2dac5",
      "base_name": "Tether USD",
      "base_symbol": "USDT",
      "quote_id": "c94ac88f-4671-3976-b60a-09064f1811e8",
      "quote_name": "Mixin",
      "quote_symbol": "XIN",
      "last_price": "0.005905789453768",
      "base_volume": "350.75774446",
      "quote_volume": "2.03407231"
    },
    "c6d0c728-2624-429b-8e0d-d9d19b6592fa_c94ac88f-4671-3976-b60a-09064f1811e8": {
      "base_id": "c6d0c728-2624-429b-8e0d-d9d19b6592fa",
      "base_name": "Bitcoin",
      "base_symbol": "BTC",
      "quote_id": "c94ac88f-4671-3976-b60a-09064f1811e8",
      "quote_name": "Mixin",
      "quote_symbol": "XIN",
      "last_price": "226.1132472698050761",
      "base_volume": "0.00100118",
      "quote_volume": "0.22612911"
    }
  }
}
```
