// Copyright 2019 Tokenomy Technologies Pte. Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v2

import "github.com/tokenomy/tokenomy-go"

//
// MarketSummaries containsall pair's latest tick value and prices in the last
// 24 hours and 7 days.
//
type MarketSummaries struct {
	Tickers   map[string]Tick              `json:"tickers"`
	Prices24h map[string]tokenomy.Rawfloat `json:"prices_24h"`
	Prices7d  map[string]tokenomy.Rawfloat `json:"prices_7d"`
}
