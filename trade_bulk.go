// Copyright 2022 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

//
// TradeBulk contains the request for bulk trading.
//
type TradeBulk struct {
	Pair      string           `json:"pair"`
	Orders    []*BulkOrderItem `json:"orders"`
	Cancel    []*BulkOrderItem `json:"cancel"`
	Timestamp int64            `json:"timestamp"`
}
