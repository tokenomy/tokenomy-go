// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v2

import "github.com/tokenomy/tokenomy-go"

//
// WithdrawItem contains the information of single withdraw transaction.
//
type WithdrawItem struct {
	ID          int64             `json:"id,omitempty"`
	RequestID   string            `json:"request_id,omitempty"`
	RequesterIP string            `json:"requester_ip,omitempty"`
	Asset       string            `json:"asset,omitempty"`
	Status      string            `json:"status,omitempty"`
	Address     string            `json:"address,omitempty"`
	Memo        string            `json:"memo,omitempty"`
	Amount      tokenomy.Rawfloat `json:"amount,omitempty"`
	Fee         tokenomy.Rawfloat `json:"fee,omitempty"`
	FinalAmount tokenomy.Rawfloat `json:"final_amount,omitempty"`
	SubmitTime  int64             `json:"submit_time,omitempty"`
	SuccessTime int64             `json:"success_time,omitempty"`
}
