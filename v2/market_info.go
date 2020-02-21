// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v2

import "github.com/tokenomy/tokenomy-go"

//
// MarketInfo contains the pair information to be consumed by public.
//
type MarketInfo struct {
	ID              string            `json:"id"`
	Symbol          string            `json:"symbol"`
	CoinAsset       string            `json:"coin_asset"`
	BaseAsset       string            `json:"base_asset"`
	IsActive        bool              `json:"is_active"`
	AmountPrecision int               `json:"amount_precision"`
	AmountMinimum   tokenomy.Rawfloat `json:"amount_minimum"`
	PricePrecision  int               `json:"price_precision"`
	PriceMinimum    tokenomy.Rawfloat `json:"price_minimum"`
}
