// Copyright 2019 Tokenomy Technologies Pte. Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

//
// TradePrice contains the information about completed trade.
//
type TradePrice struct {
	ID         int64    `json:"id"`
	TradeTime  int64    `json:"trade_time"`
	Amount     Rawfloat `json:"amount"`
	AmountCoin Rawfloat `json:"amount_coin"`
	Price      Rawfloat `json:"price"`
}
