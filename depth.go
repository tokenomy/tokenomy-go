// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

import "github.com/shuLhan/share/lib/math/big"

// Depth contains total amount of remaining open orders grouped by price.
// Each depth is specific to pair.
type Depth struct {
	Amount    *big.Rat `json:"amount"` // DEPRECATED: replaced with total_base and total_coin.
	Price     *big.Rat `json:"price"`
	TotalBase *big.Rat `json:"total_base"`
	TotalCoin *big.Rat `json:"total_coin"`
}
