// Copyright 2019 Tokenomy Technologies Pte. Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

//
// WithdrawResponse contains status, balance, and fee information from
// withdrawal request.
//
type WithdrawResponse struct {
	Success        int
	Error          string
	ErrorCode      string
	ID             string `json:"withdraw_id"`
	Txid           string
	Status         string
	Currency       string `json:"withdraw_currency"`
	Address        string `json:"withdraw_address"`
	Amount         string `json:"withdraw_amount"`
	Fee            string
	AmountAfterFee string
	SubmitTime     string
	Memo           string `json:"withdraw_memo"`
}
