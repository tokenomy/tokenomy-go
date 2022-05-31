// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

import (
	big "github.com/shuLhan/share/lib/math/big"
)

// UserAssets contains mapping between asset name and its value.
type UserAssets struct {
	Balances          map[string]*big.Rat `json:"balances,omitempty"`
	FrozenBalances    map[string]*big.Rat `json:"frozen_balances,omitempty"`
	BalancesInt       map[string]int64    `json:"-"`
	FrozenBalancesInt map[string]int64    `json:"-"`
}

// NewUserAssets create and initialize all fields in UserAssets.
func NewUserAssets() (assets *UserAssets) {
	return &UserAssets{
		Balances:          make(map[string]*big.Rat),
		BalancesInt:       make(map[string]int64),
		FrozenBalances:    make(map[string]*big.Rat),
		FrozenBalancesInt: make(map[string]int64),
	}
}

// Copy create a new copy of assets.
func (assets *UserAssets) Copy() *UserAssets {
	newAssets := &UserAssets{
		Balances:          make(map[string]*big.Rat, len(assets.Balances)),
		FrozenBalances:    make(map[string]*big.Rat, len(assets.FrozenBalances)),
		BalancesInt:       make(map[string]int64, len(assets.BalancesInt)),
		FrozenBalancesInt: make(map[string]int64, len(assets.FrozenBalancesInt)),
	}

	for k, v := range assets.Balances {
		newAssets.Balances[k] = big.NewRat(v)
	}
	for k, v := range assets.FrozenBalances {
		newAssets.FrozenBalances[k] = big.NewRat(v)
	}
	for k, v := range assets.BalancesInt {
		newAssets.BalancesInt[k] = v
	}
	for k, v := range assets.FrozenBalancesInt {
		newAssets.FrozenBalancesInt[k] = v
	}

	return newAssets
}
