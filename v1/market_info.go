// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v1

import "github.com/shuLhan/share/lib/math/big"

type PairPrecision struct {
	Price  int
	Amount int
}

type PairAmountLimit struct {
	Min *big.Rat
}

type PairPriceLimit struct {
	Min *big.Rat
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
