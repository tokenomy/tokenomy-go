// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

import "github.com/shuLhan/share/lib/math/big"

//
// Order contains common information about bid or ask order, either open or
// closed.
//
type Order struct {
	ID         int64    `json:"id"`
	Pair       string   `json:"pair,omitempty"`
	Type       string   `json:"type,omitempty"` // Type of order, either "sell" or "buy".
	Method     string   `json:"method"`         // Method of order, either "limit" or "market".
	SubmitTime int64    `json:"submit_time"`
	FinishTime int64    `json:"finish_time,omitempty"`
	Status     string   `json:"status,omitempty"` // Status for closed order.
	Price      *big.Rat `json:"price,omitempty"`

	AmountBase *big.Rat `json:"amount_base,omitempty"`
	RemainBase *big.Rat `json:"remain_base,omitempty"`
	FilledBase *big.Rat `json:"filled_base,omitempty"`

	AmountCoin *big.Rat `json:"amount_coin,omitempty"`
	RemainCoin *big.Rat `json:"remain_coin,omitempty"`
	FilledCoin *big.Rat `json:"filled_coin,omitempty"`
}
