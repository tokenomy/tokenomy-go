// Copyright 2019 Tokenomy Technologies Pte. Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

type getOrderResponse struct {
	Success int
	Error   string
	Return  *singleOrder
}

type singleOrder struct {
	Order *OrderHistory
}
