// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v2

import "github.com/shuLhan/share/lib/math/big"

//
// MarketSummaries containsall pair's latest tick value and prices in the last
// 24 hours and 7 days.
//
type MarketSummaries struct {
	Tickers   map[string]Tick     `json:"tickers"`
	Prices24h map[string]*big.Rat `json:"prices_24h"`
	Prices7d  map[string]*big.Rat `json:"prices_7d"`
}
