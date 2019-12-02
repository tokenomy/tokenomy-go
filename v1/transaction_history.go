// Copyright 2019 Tokenomy Technologies Pte. Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v1

type transHistoryResponse struct {
	Success   int
	Error     string
	ErrorCode string
	Return    *TransactionHistory
}

//
// TransDeposit represent deposit transaction by user for all status.
//
type TransDeposit struct {
	DepositID   string `json:"deposit_id"`
	Status      string
	Asset       string
	Amount      string
	FinalAmount string
	SuccessTime string
}

//
// TransWithdraw represent withdraw transaction by user for all status.
//
type TransWithdraw struct {
	WithdrawID  string `json:"withdraw_id"`
	Status      string
	Asset       string
	Amount      string
	Fee         string
	FinalAmount string
	SubmitTime  string
	SuccessTime string
}

//
// TransactionHistory contains all user's deposits and withdraws.
//
type TransactionHistory struct {
	Deposit  map[string][]TransDeposit
	Withdraw map[string][]TransWithdraw
}
