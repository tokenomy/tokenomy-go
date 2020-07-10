// Copyright 2020 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

//
// Client is an interface for API v1 and v2.
//
type Client interface {
	TradeAsk(treq *TradeRequest) (*TradeResponse, error)
	TradeBid(treq *TradeRequest) (*TradeResponse, error)
	TradeCancelAsk(pairName string, orderID int64) (*TradeResponse, error)
	TradeCancelBid(pairName string, orderID int64) (*TradeResponse, error)
}
