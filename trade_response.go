// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

//
// TradeResponse contains the recorded order information in the market and
// part of full amount of trade that have been completed (matched) as a list
// of deals.
//
type TradeResponse struct {
	Order *Order `json:"order"`
	User  User   `json:"user"`
	// Deals contains matched orders, only available on v2.
	Deals []TradePrice `json:"deals,omitempty"`
}
