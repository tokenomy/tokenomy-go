// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

import "github.com/shuLhan/share/lib/math/big"

//
// Trade contains information about trade bid or ask, either open or
// closed.
//
type Trade struct {
	ID     int64    `json:"id,omitempty"`
	Pair   string   `json:"pair,omitempty"`
	Type   string   `json:"type,omitempty"`   // Its either "sell" or "buy".
	Method string   `json:"method,omitempty"` // Its either "limit" or "market".
	Status string   `json:"status,omitempty"` // Status for closed trade.
	Price  *big.Rat `json:"price,omitempty"`
	Fee    *big.Rat `json:"fee,omitempty"`

	BaseAsset  string   `json:"base_asset,omitempty"`
	BaseAmount *big.Rat `json:"base_amount,omitempty"`
	BaseRemain *big.Rat `json:"base_remain,omitempty"`
	BaseFilled *big.Rat `json:"base_filled,omitempty"`

	CoinAsset  string   `json:"coin_asset,omitempty"`
	CoinAmount *big.Rat `json:"coin_amount,omitempty"`
	CoinRemain *big.Rat `json:"coin_remain,omitempty"`
	CoinFilled *big.Rat `json:"coin_filled,omitempty"`

	SubmitTime int64 `json:"submit_time"`
	FinishTime int64 `json:"finish_time,omitempty"`
}
