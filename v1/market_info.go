// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v1

type PairPrecision struct {
	Price  int
	Amount int
}

type PairAmountLimit struct {
	Min float64
}

type PairPriceLimit struct {
	Min float64
}

type PairLimits struct {
	Amount PairAmountLimit
	Price  PairPriceLimit
}

type MarketInfo struct {
	ID        string
	Symbol    string
	Base      string
	Quote     string
	Active    bool
	Precision PairPrecision
	Limits    PairLimits
}
