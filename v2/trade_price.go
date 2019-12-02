// Copyright 2019 Tokenomy Technologies Pte. Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v2

import "github.com/tokenomy/tokenomy-go"

type TradePrice struct {
	ID         int64             `json:"id"`
	TradeTime  int64             `json:"trade_time"`
	Amount     tokenomy.Rawfloat `json:"amount"`
	AmountCoin tokenomy.Rawfloat `json:"amount_coin"`
	Price      tokenomy.Rawfloat `json:"price"`
}
