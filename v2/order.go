// Copyright 2019 Tokenomy Technologies Pte. Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v2

import "github.com/tokenomy/tokenomy-go"

//
// Order contains common information of bid or ask order, either open or
// closed.
//
type Order struct {
	ID         int64             `json:"id"`
	Type       string            `json:"type,omitempty"` // Type of order, either "sell" or "buy".
	Method     string            `json:"method"`         // Method of order, either "limit" or "market".
	SubmitTime int64             `json:"submit_time"`
	FinishTime int64             `json:"finish_time,omitempty"`
	Status     string            `json:"status,omitempty"` // Status for closed order.
	Price      tokenomy.Rawfloat `json:"price",omitempty`

	AmountBase tokenomy.Rawfloat `json:"amount_base,omitempty"`
	RemainBase tokenomy.Rawfloat `json:"remain_base,omitempty"`
	FilledBase tokenomy.Rawfloat `json:"filled_base,omitempty"`

	AmountCoin tokenomy.Rawfloat `json:"amount_coin,omitempty"`
	RemainCoin tokenomy.Rawfloat `json:"remain_coin,omitempty"`
	FilledCoin tokenomy.Rawfloat `json:"filled_coin,omitempty"`
}
