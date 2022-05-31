// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

// TradeResponse contains the recorded order information in the market and
// part of full amount of order that have been completed (matched) as a list
// of trades.
type TradeResponse struct {
	Order *Trade `json:"order"`
	User  User   `json:"user"`
	// Trades contains matched orders, only available on v2.
	Trades []Trade `json:"trades,omitempty"`
}
