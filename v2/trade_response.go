// Copyright 2019 Tokenomy Technologies Pte. Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v2

//
// TradeResponse contains the recorded trade information in the market and
// part of full amount of trade that has completed as a list of deals.
//
type TradeResponse struct {
	Order *Order
	Deals []TradePrice
}
