// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

import "github.com/shuLhan/share/lib/math/big"

//
// MarketTicker contains the pair tick information in the market.
//
type MarketTicker struct {
	// LowestAskPrice contains the lowest sell price in open orders.
	LowestAskPrice *big.Rat `json:"ask"`

	// HighestBidPrice contains the highest buy price in open orders.
	HighestBidPrice *big.Rat `json:"bid"`

	// HighestPrice24H contains the highest price in closed trades since
	// the last 24 hours.
	HighestPrice24H *big.Rat `json:"high"`

	// LowestPrice24H contains the lowest price in closed trades since the
	// last 24 hours.
	LowestPrice24H *big.Rat `json:"low"`

	// LastPrice contains the last traded price.
	LastPrice *big.Rat `json:"last_price"`

	// VolumeBase24H contains the total base asset has been traded since
	// the last 24 hours.
	VolumeBase24H *big.Rat `json:"volume_base"`

	// VolumeCoin24H contains the total coin asset has been traded since
	// the last 24 hours.
	VolumeCoin24H *big.Rat `json:"volume_coin"`

	PairName string `json:"pair"`
}
