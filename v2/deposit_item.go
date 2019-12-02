// Copyright 2019 Tokenomy Technologies Pte. Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v2

import "github.com/tokenomy/tokenomy-go"

//
// DepositItem contains the information of deposit.
//
type DepositItem struct {
	ID          int64             `json:"id"`
	Asset       string            `json:"asset,omitempty"`
	Status      string            `json:"status"`
	Amount      tokenomy.Rawfloat `json:"amount"`
	FinalAmount tokenomy.Rawfloat `json:"final_amount"`
	SuccessTime int64             `json:"success_time"`
}
