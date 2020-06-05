// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v1_test

import (
	"fmt"
	"log"
	"os"

	"github.com/shuLhan/share/lib/math/big"
	"github.com/tokenomy/tokenomy-go"
	v1 "github.com/tokenomy/tokenomy-go/v1"
)

func ExampleClient_TradeBid() {
	env := tokenomy.NewEnvironment("", "")

	cl, err := v1.NewClient(env)
	if err != nil {
		log.Fatal(err)
	}

	amount := big.NewRat(100)
	price := big.NewRat(0.00005)

	tres, err := cl.TradeBid(tokenomy.TradeMethodLimit,
		tokenomy.PairTokenomyBitcoin, amount, price)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Trade bid response: %+v\n", tres)
}

func ExampleClient_TradeBid_by_market() {
	env := tokenomy.NewEnvironment("", "")

	cl, err := v1.NewClient(env)
	if err != nil {
		log.Fatal(err)
	}

	amount := big.NewRat(0.0001)
	price := big.NewRat(0)

	tres, err := cl.TradeBid(tokenomy.TradeMethodMarket,
		tokenomy.PairTokenomyBitcoin, amount, price)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Trade bid by market response: %+v\n", tres)
}

func ExampleClient_TradeCancelBid() {
	env := tokenomy.NewEnvironment("", "")

	cl, err := v1.NewClient(env)
	if err != nil {
		log.Fatal(err)
	}

	openOrders, err := cl.UserOrdersOpen(tokenomy.PairTokenomyBitcoin)
	if err != nil {
		log.Fatal(err)
	}
	if len(openOrders) == 0 {
		fmt.Println("No open orders to cancel")
		return
	}

	var openBuy *v1.OrderHistory

	fmt.Println("Open orders:")
	for pairName, list := range openOrders {
		fmt.Printf("[%s]\n", pairName)
		for x, order := range list {
			if order.Type == tokenomy.TradeTypeBid && openBuy == nil {
				openBuy = order
			}
			fmt.Printf("  [%d] - %+v\n", x, order)
		}
	}

	if openBuy == nil {
		fmt.Println("There is no open bid to cancel.")
		return
	}

	fmt.Printf("Canceling the first open bid with ID '%d'\n", openBuy.OrderID)

	cancelRes, err := cl.TradeCancelBid(tokenomy.PairTokenomyBitcoin,
		openBuy.OrderID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Cancel bid response: %+v\n", cancelRes)
}

func ExampleClient_TradeCancelAsk() {
	env := tokenomy.NewEnvironment("", "")

	cl, err := v1.NewClient(env)
	if err != nil {
		log.Fatal(err)
	}

	openOrders, err := cl.UserOrdersOpen(tokenomy.PairTokenomyBitcoin)
	if err != nil {
		log.Fatal(err)
	}
	if len(openOrders) == 0 {
		fmt.Println("No open orders to cancel")
		return
	}

	var openSell *v1.OrderHistory

	fmt.Println("Open orders:")
	for pairName, list := range openOrders {
		fmt.Printf("[%s]\n", pairName)
		for x, order := range list {
			if order.Type == tokenomy.TradeTypeAsk && openSell == nil {
				openSell = order
			}
			fmt.Printf("  [%d] - %+v\n", x, order)
		}
	}
	if openSell == nil {
		fmt.Println("There is no open sell (ask) to cancel.")
		return
	}

	fmt.Printf("Canceling the open sell (ask) with ID '%d'\n", openSell.OrderID)

	cancelRes, err := cl.TradeCancelAsk(tokenomy.PairTokenomyBitcoin,
		openSell.OrderID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Cancel sell (ask) response: %+v\n", cancelRes)
}

func ExampleClient_UserOrder() {
	env := tokenomy.NewEnvironment("", "")
	orderID := int64(1023965)

	cl, err := v1.NewClient(env)
	if err != nil {
		log.Fatal(err)
	}

	order, err := cl.UserOrder(tokenomy.PairTokenomyBitcoin, orderID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Order detail for %d: %+v\n", orderID, order)
}

func ExampleClient_MarketTicker() {
	env := tokenomy.NewEnvironment("", "")

	cl, err := v1.NewClient(env)
	if err != nil {
		log.Fatal(err)
	}

	pair, err := cl.MarketTicker(tokenomy.PairTokenomyBitcoin)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Pair information for %q: %+v\n", tokenomy.PairTokenomyBitcoin, pair)
}

func ExampleClient_UserOrdersOpen() {
	env := tokenomy.NewEnvironment("", "")

	cl, err := v1.NewClient(env)
	if err != nil {
		log.Fatal(err)
	}

	openOrders, err := cl.UserOrdersOpen(tokenomy.PairTokenomyBitcoin)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Open orders:")
	for pairName, list := range openOrders {
		fmt.Printf("[%s]\n", pairName)
		for x, order := range list {
			fmt.Printf("  [%d] - %+v\n", x, order)
		}
	}
}

func ExampleClient_UserOrdersClosed() {
	env := tokenomy.NewEnvironment("", "")

	cl, err := v1.NewClient(env)
	if err != nil {
		log.Fatal(err)
	}

	orderHistory, err := cl.UserOrdersClosed(tokenomy.PairTokenomyBitcoin, 0, 0)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Order history:")
	for x, order := range orderHistory {
		fmt.Printf("  [%d] %+v\n", x, order)
	}
}

func ExampleClient_MarketTrades() {
	env := tokenomy.NewEnvironment("", "")

	cl, err := v1.NewClient(env)
	if err != nil {
		log.Fatal(err)
	}

	trades, err := cl.MarketTrades(tokenomy.PairTokenomyBitcoin)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Trades information for pair %q:", tokenomy.PairTokenomyBitcoin)
	for _, trade := range trades {
		fmt.Printf("  %+v\n", trade)
	}
}

func ExampleClient_UserTrades() {
	env := tokenomy.NewEnvironment("", "")

	cl, err := v1.NewClient(env)
	if err != nil {
		log.Fatal(err)
	}

	tradeHistory, err := cl.UserTrades(tokenomy.PairTokenomyBitcoin,
		1, 0, 0, "asc", nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Trade history:")
	for x, trade := range tradeHistory {
		fmt.Printf("  [%d]: %+v\n", x, trade)
	}
}

func ExampleClient_UserTransactions() {
	env := tokenomy.NewEnvironment(
		os.Getenv(tokenomy.EnvNameToken),
		os.Getenv(tokenomy.EnvNameSecret),
	)

	cl, err := v1.NewClient(env)
	if err != nil {
		log.Fatal(err)
	}

	transHistory, err := cl.UserTransactions()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Transaction history: %+v\n", transHistory)
}

func ExampleClient_MarketInfo() {
	env := tokenomy.NewEnvironment("", "")

	cl, err := v1.NewClient(env)
	if err != nil {
		log.Fatal(err)
	}

	marketInfos, err := cl.MarketInfo()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Market Info: %+v\n", marketInfos)
}

func ExampleClient_MarketOrdersOpen() {
	env := tokenomy.NewEnvironment("", "")

	cl, err := v1.NewClient(env)
	if err != nil {
		log.Fatal(err)
	}

	orderBook, err := cl.MarketOrdersOpen(tokenomy.PairTokenomyBitcoin)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Order book buys:")
	for _, buy := range orderBook.Buys {
		fmt.Printf("  %+v\n", buy)
	}

	fmt.Println("Order book sells:")
	for _, sell := range orderBook.Sells {
		fmt.Printf("  %+v\n", sell)
	}
}

func ExampleClient_TradeAsk_by_market() {
	env := tokenomy.NewEnvironment("", "")

	cl, err := v1.NewClient(env)
	if err != nil {
		log.Fatal(err)
	}

	amount := big.NewRat(20)
	price := big.NewRat(0)

	tres, err := cl.TradeAsk(tokenomy.TradeMethodMarket,
		tokenomy.PairTokenomyBitcoin, amount, price)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Trade ask by market response: %+v\n", tres)
}

func ExampleClient_MarketSummaries() {
	env := tokenomy.NewEnvironment("", "")

	cl, err := v1.NewClient(env)
	if err != nil {
		log.Fatal(err)
	}

	summary, err := cl.MarketSummaries()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Summaries:")
	fmt.Println("  Tickers:")
	for name, ticker := range summary.Pairs {
		fmt.Printf("    %s: %+v\n", name, ticker)
	}
	fmt.Printf("  Prices for the last 24 hours: %v\n", summary.PricesLast7Days)
	fmt.Printf("  Prices for the last 7 days  : %v\n", summary.PricesLast24Hours)
}

func ExampleClient_UserInfo() {
	env := tokenomy.NewEnvironment(
		os.Getenv(tokenomy.EnvNameToken),
		os.Getenv(tokenomy.EnvNameSecret),
	)

	cl, err := v1.NewClient(env)
	if err != nil {
		log.Fatal(err)
	}

	userInfo, err := cl.UserInfo()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", userInfo)
}
