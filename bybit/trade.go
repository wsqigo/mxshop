package main

import (
	"context"
	"fmt"
	bybit "github.com/wuhewuhe/bybit.go.api"
)

func main() {
	client := bybit.NewBybitHttpClient("8ZuF0gBseHHDPRVEt3", "5ya70v8sZax6EhpqdkxuEKedJKvpqCHhFnLp", bybit.WithBaseURL(bybit.TESTNET))
	params := map[string]interface{}{
		"category":    "linear",
		"symbol":      "BTCUSDT",
		"side":        "Buy",
		"positionIdx": 0,
		"orderType":   "Limit",
		"qty":         "0.001",
		"price":       "10000",
		"timeInForce": "GTC",
	}
	orderResult, err := client.NewTradeService(params).PlaceOrder(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(bybit.PrettyPrint(orderResult))
}
