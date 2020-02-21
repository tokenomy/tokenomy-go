// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v2

//
// AssetTransactions contains list of deposit and withdraw transaction.
//
type AssetTransactions struct {
	Deposit  map[string][]DepositItem  `json:"deposit"`
	Withdraw map[string][]WithdrawItem `json:"withdraw"`
}
