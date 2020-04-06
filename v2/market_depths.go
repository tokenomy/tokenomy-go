// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v2

//
// MarketDepths contains list of depth on open asks and bids.
//
type MarketDepths struct {
	Asks []Depth `json:"asks"`
	Bids []Depth `json:"bids"`
}
