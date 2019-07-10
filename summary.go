// Copyright 2019 Tokenomy Technologies Pte. Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

//
// Summary contains list of last pair informations, and prices of pair in the
// last 24 hours and 7 days.
//
type Summary struct {
	Pairs             map[string]*Pair
	PricesLast24Hours map[string]interface{} `json:"prices_24h"`
	PricesLast7Days   map[string]interface{} `json:"prices_7d"`
}

//
// propagate the summary's fields, initialize some its dynamic values with
// custom parsing.
//
func (summary *Summary) propagate() {
	for name, pair := range summary.Pairs {
		pair.propagate(name)
	}
}
