// Copyright 2019 Tokenomy Technologies Pte. Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v2

import "github.com/tokenomy/tokenomy-go"

//
// MarketTradePrices contains list of closed ask and bid on the market.
//
type MarketTradePrices struct {
	Asks []tokenomy.TradePrice `json:"asks"`
	Bids []tokenomy.TradePrice `json:"bids"`
}
