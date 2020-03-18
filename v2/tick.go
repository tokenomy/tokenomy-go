// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v2

import "github.com/shuLhan/share/lib/math/big"

//
// Tick contains the pair tick information.
//
type Tick struct {
	PairName string `json:"pair"`

	// Bid contains the highest buy price in open buy orders.
	Bid *big.Rat `json:"bid"`

	// Ask contains the lowest sell price in open sell orders.
	Ask *big.Rat `json:"ask"`

	// High contains the highest price in closed orders from last 24
	// hours.
	High *big.Rat `json:"high"`

	// Low contains the lowest price in closed orders from last 24 hours.
	Low *big.Rat `json:"low"`

	// LastPrice contains the last traded price.
	LastPrice *big.Rat `json:"last_price"`

	// VolumeBase contains the traded base asset.
	VolumeBase *big.Rat `json:"volume_base"`

	// VolumeCoin contains the traded coin asset.
	VolumeCoin *big.Rat `json:"volume_coin"`
}

//
// IsZero will return true if all fields' value is zero.
//
func (tick *Tick) IsZero() bool {
	if tick.Ask.IsZero() && tick.Bid.IsZero() &&
		tick.High.IsZero() && tick.Low.IsZero() &&
		tick.LastPrice.IsZero() &&
		tick.VolumeBase.IsZero() && tick.VolumeCoin.IsZero() {
		return true
	}
	return false
}
