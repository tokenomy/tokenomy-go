// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

import "github.com/shuLhan/share/lib/math/big"

//
// MarketPrices contains mapping between pair and its latest price in the
// market.
//
type MarketPrices map[string]*big.Rat
