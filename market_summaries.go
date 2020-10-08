// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

import "github.com/shuLhan/share/lib/math/big"

//
// MarketSummaries contains all pair's latest tick value, last price, price
// changes, price past 24 hours ago, and price past 7 days ago.
//
type MarketSummaries struct {
	Prices        map[string]*big.Rat `json:"prices"`
	Prices24h     map[string]*big.Rat `json:"prices_24h"`
	Prices7d      map[string]*big.Rat `json:"prices_7d"`
	PricesChanges map[string]*big.Rat `json:"prices_changes"`
	Tickers       map[string]Tick     `json:"tickers"`
}
