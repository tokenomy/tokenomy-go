// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v2

import "github.com/tokenomy/tokenomy-go"

//
// MarketPrices contains mapping between pair and its latest price in the
// market.
//
type MarketPrices map[string]tokenomy.Rawfloat
