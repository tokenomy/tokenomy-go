// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

// TradesOpen contains the open asks and bids in the market place.
type TradesOpen struct {
	Asks []Trade `json:"asks"`
	Bids []Trade `json:"bids"`
}
