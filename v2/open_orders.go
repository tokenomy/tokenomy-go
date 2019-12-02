// Copyright 2019 Tokenomy Technologies Pte. Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v2

//
// OpenOrders contains the open ask and bid orders in the market place.
//
type OpenOrders struct {
	Asks []Order `json:"asks"`
	Bids []Order `json:"bids"`
}
