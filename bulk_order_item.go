// Copyright 2022 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

import liberrors "github.com/shuLhan/share/lib/errors"

// BulkOrderItem represent single order in bulk trading.
type BulkOrderItem struct {
	liberrors.E
	TradeRequest

	ID    int64 `json:"id,omitempty"`
	RefID int64 `json:"ref_id,omitempty"`
}
