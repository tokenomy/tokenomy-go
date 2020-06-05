// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v2

import (
	"github.com/tokenomy/tokenomy-go"
)

//
// TradesOpen contains the open asks and bids in the market place.
//
type TradesOpen struct {
	Asks []tokenomy.Trade `json:"asks"`
	Bids []tokenomy.Trade `json:"bids"`
}
