// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

import "github.com/shuLhan/share/lib/math/big"

//
// TradePrice contains the information about completed trade.
//
type TradePrice struct {
	ID         int64    `json:"id"`
	TradeTime  int64    `json:"trade_time"`
	Amount     *big.Rat `json:"amount"`
	AmountCoin *big.Rat `json:"amount_coin"`
	Price      *big.Rat `json:"price"`
}
