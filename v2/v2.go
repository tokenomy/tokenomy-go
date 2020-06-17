// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

//
// Package v2 is the official Go library for Tokenomy API v2.
// This package provide client for REST and Websocket API.
//
package v2

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/shuLhan/share/lib/math/big"
	"github.com/tokenomy/tokenomy-go"
)

const (
	DefaultAddress = "https://api.tokenomy.com"

	apiMarketDepths     = "/v2/market/depths"
	apiMarketInfo       = "/v2/market/info"
	apiMarketTradesOpen = "/v2/market/trades/open"
	apiMarketPrices     = "/v2/market/prices"
	apiMarketTicker     = "/v2/market/ticker"
	apiMarketTrades     = "/v2/market/trades"
	apiMarketSummaries  = "/v2/market/summaries"

	apiUserInfo         = "/v2/user/info"
	apiUserTrades       = "/v2/user/trades"
	apiUserTradesClosed = "/v2/user/trades/closed"
	apiUserTradesOpen   = "/v2/user/trades/open"
	apiUserTrade        = "/v2/user/trade"
	apiUserTransactions = "/v2/user/transactions"
	apiUserWithdraw     = "/v2/user/withdraw"

	apiTradeAsk       = "/v2/trade/ask"
	apiTradeBid       = "/v2/trade/bid"
	apiTradeCancelAsk = "/v2/trade/cancel/ask"
	apiTradeCancelBid = "/v2/trade/cancel/bid"

	wsPrivateEndpoint = "/v2/user/ws"

	wsPublicEndpoint = "/v2/ws"
)

func generateTradeParams(method, pairName string, amount, price *big.Rat) (
	params url.Values, wsparams *WebSocketParams, err error,
) {
	params = url.Values{}
	if len(method) == 0 {
		method = tokenomy.TradeMethodLimit
	} else {
		method = strings.ToLower(method)
		switch method {
		case tokenomy.TradeMethodMarket, tokenomy.TradeMethodLimit:
		default:
			return nil, nil, tokenomy.ErrInvalidTradeMethod
		}
	}

	if amount == nil || amount.IsLessOrEqual(0) {
		return nil, nil, tokenomy.ErrInvalidAmount
	}

	params.Set(tokenomy.ParamNameTradeMethod, method)
	params.Set(tokenomy.ParamNamePair, pairName)
	params.Set(tokenomy.ParamNameAmount, amount.String())

	wsparams = &WebSocketParams{
		Method: method,
		Pair:   pairName,
		Amount: amount,
	}

	if method == tokenomy.TradeMethodLimit {
		if price == nil || price.IsLessOrEqual(0) {
			return nil, nil, tokenomy.ErrInvalidPrice
		}
		params.Set(tokenomy.ParamNamePrice, price.String())
		wsparams.Price = price
	}

	return params, wsparams, nil
}

func timestamp() int64 {
	return time.Now().Unix()
}

func timestampAsString() string {
	return strconv.FormatInt(timestamp(), 10)
}
