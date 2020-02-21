// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v2

import "github.com/tokenomy/tokenomy-go"

//
// Trade represent a completed trade information on user.
//
type Trade struct {
	ID          int64             `json:"id"`
	TradeID     int64             `json:"trade_id"`
	TradeMethod string            `json:"trade_method"`
	CoinAsset   string            `json:"coin_asset"`
	BaseAsset   string            `json:"base_asset"`
	AmountBase  tokenomy.Rawfloat `json:"amount_base"`
	AmountCoin  tokenomy.Rawfloat `json:"amount_coin"`
	Price       tokenomy.Rawfloat `json:"price"`
	Fee         tokenomy.Rawfloat `json:"fee"`
	Time        int64             `json:"trade_time"`
}
