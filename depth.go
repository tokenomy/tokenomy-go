// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

import "github.com/shuLhan/share/lib/math/big"

//
// Depth contains total amount of remaining trade grouped by price in open
// trades.
// Each depth is specific to pair.
//
type Depth struct {
	Amount *big.Rat `json:"amount"`
	Price  *big.Rat `json:"price"`
}
