// Copyright 2019 Tokenomy Technologies Pte. Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v2

import "github.com/tokenomy/tokenomy-go"

//
// Tick contains the pair tick information.
//
type Tick struct {
	PairName   string            `json:"pair"`
	Bid        tokenomy.Rawfloat `json:"bid"`         // Bid contains the highest buy price in open buy orders.
	Ask        tokenomy.Rawfloat `json:"ask"`         // Ask contains the lowest sell price in open sell orders.
	High       tokenomy.Rawfloat `json:"high"`        // High contains the highest price in closed orders from last 24 hours.
	Low        tokenomy.Rawfloat `json:"low"`         // Low contains the lowest price in closed orders from last 24 hours.
	LastPrice  tokenomy.Rawfloat `json:"last_price"`  // LastPrice contains the last traded price.
	VolumeBase tokenomy.Rawfloat `json:"volume_base"` // VolumeBase contains the traded base asset.
	VolumeCoin tokenomy.Rawfloat `json:"volume_coin"` // VolumeCoin contains the traded coin asset.
}

//
// IsZero will return true if all fields' value is zero.
//
func (tick *Tick) IsZero() bool {
	if tick.Ask == 0 && tick.Bid == 0 &&
		tick.High == 0 && tick.Low == 0 &&
		tick.LastPrice == 0 &&
		tick.VolumeBase == 0 && tick.VolumeCoin == 0 {
		return true
	}
	return false
}
