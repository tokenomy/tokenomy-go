// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v1

type openOrdersResponse struct {
	Success   int
	Error     string
	ErrorCode string
	Return    *pairOpenOrders
}

type pairOpenOrders struct {
	Orders OpenOrders
}

//
// OpenOrders contains mapping between pair names and their order history.
//
type OpenOrders map[string][]*OrderHistory
