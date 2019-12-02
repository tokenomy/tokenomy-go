// Copyright 2019 Tokenomy Technologies Pte. Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v1

//
// OrderBook contains the open buy (bid) and sell (ask).
//
type OrderBook struct {
	Buys  []*Order `json:"buy"`
	Sells []*Order `json:"sell"`
}
