// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

import "github.com/shuLhan/share/lib/math/big"

//
// Tick contains the pair tick information.
//
type Tick struct {
	PairName string `json:"pair"`

	// Bid contains the highest buy price in open buy.
	Bid *big.Rat `json:"bid"`

	// Ask contains the lowest sell price in open sell.
	Ask *big.Rat `json:"ask"`

	// High contains the highest price in closed trades since the last 24
	// hours.
	High *big.Rat `json:"high"`

	// Low contains the lowest price in closed trades since the last 24
	// hours.
	Low *big.Rat `json:"low"`

	// LastPrice contains the last traded price.
	LastPrice *big.Rat `json:"last_price"`

	// VolumeBase contains the total base asset has been traded since the
	// last 24 hours.
	VolumeBase *big.Rat `json:"volume_base"`

	// VolumeCoin contains the total coin asset has been traded since the
	// last 24 hours.
	VolumeCoin *big.Rat `json:"volume_coin"`
}
