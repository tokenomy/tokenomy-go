// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

import (
	"strconv"
	"time"
)

const (
	DefaultAddress = "https://api.tokenomy.com"
)

// List of API endpoints.
const (
	APIMarketDepths     = "/v2/market/depths"
	APIMarketInfo       = "/v2/market/info"
	APIMarketTradesOpen = "/v2/market/trades/open"
	APIMarketPrices     = "/v2/market/prices"
	APIMarketTicker     = "/v2/market/ticker"
	APIMarketTrades     = "/v2/market/trades"
	APIMarketSummaries  = "/v2/market/summaries"

	APIUserInfo         = "/v2/user/info"
	APIUserTrades       = "/v2/user/trades"
	APIUserOrdersClosed = "/v2/user/orders/closed"
	APIUserOrdersOpen   = "/v2/user/orders/open"
	APIUserOrderInfo    = "/v2/user/order"
	APIUserTransactions = "/v2/user/transactions"
	APIUserWithdraw     = "/v2/user/withdraw"

	APITradeAsk       = "/v2/trade/ask"
	APITradeBid       = "/v2/trade/bid"
	APITradeCancelAll = "/v2/trade/cancel/all"
	APITradeCancelAsk = "/v2/trade/cancel/ask"
	APITradeCancelBid = "/v2/trade/cancel/bid"

	WSPrivate = "/v2/user/ws"

	WSPublic             = "/v2/ws"
	WSPublicSubscription = "/v2/ws/subscription"
)

// List of WebSocket broadcast messages.
const (
	WSMessageUserOrdersTaken = "/v2/user/orders/taken"
)

func timestamp() int64 {
	return time.Now().Unix()
}

func timestampAsString() string {
	return strconv.FormatInt(timestamp(), 10)
}
