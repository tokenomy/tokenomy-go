// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

import "github.com/shuLhan/share/lib/math/big"

// WithdrawItem contains the information of single withdraw transaction.
type WithdrawItem struct {
	Amount      *big.Rat `json:"amount,omitempty"`
	Fee         *big.Rat `json:"fee,omitempty"`
	FinalAmount *big.Rat `json:"final_amount,omitempty"`

	RequestID   string `json:"request_id,omitempty"`
	RequesterIP string `json:"requester_ip,omitempty"`
	Asset       string `json:"asset,omitempty"`
	Network     string `json:"network,omitempty"`
	Status      string `json:"status,omitempty"`
	Address     string `json:"address,omitempty"`
	AddressType string `json:"address_type,omitempty"`
	Memo        string `json:"memo,omitempty"`

	ID          int64 `json:"id,omitempty"`
	SubmitTime  int64 `json:"submit_time,omitempty"`
	SuccessTime int64 `json:"success_time,omitempty"`
}
