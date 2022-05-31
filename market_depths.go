// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

import "github.com/shuLhan/share/lib/math/big"

// MarketDepths contains list of depth on open asks and bids.
type MarketDepths struct {
	Pair string   `json:"pair"`
	Asks []*Depth `json:"asks"`
	Bids []*Depth `json:"bids"`
}

// GetAskByPrice get the depth record from list of Asks by its price.
func (depths *MarketDepths) GetAskByPrice(price *big.Rat) (depth *Depth) {
	for _, depth = range depths.Asks {
		if depth.Price.IsEqual(price) {
			return depth
		}
	}
	return nil
}

// GetBidByPrice get the depth record from list of Bids by its price.
func (depths *MarketDepths) GetBidByPrice(price *big.Rat) (depth *Depth) {
	for _, depth = range depths.Bids {
		if depth.Price.IsEqual(price) {
			return depth
		}
	}
	return nil
}
