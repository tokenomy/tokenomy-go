// Copyright 2020 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

//
// PublicSubscription contains list of pairs that currently subscribed for
// each topic: "depths", "ticker", and "trades".
//
type PublicSubscription struct {
	Depths    []string `json:"depths"`
	Ticker    []string `json:"ticker"`
	Trades    []string `json:"trades"`
	Summaries bool     `json:"summaries"`
}
