// Copyright 2021 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

//
// TradeParams represent all list of parameters for querying user's
// trades.
//
type TradeParams struct {
	Pair string

	// The Offset field define the number of rows to be skipped.
	Offset int64

	// The Limit field define the maximum number of record fetched, if its not
	// set default to DefaultLimit.
	Limit int64

	// The IDAfter filter rows with ID greater or equal than its value.
	IDAfter int64

	// The IDBefore filter rows with ID less than or equal to its value.
	IDBefore int64

	// The TimeAfter filter rows with trade's time greater or equal than its
	// value.
	TimeAfter int64

	// Then TimeBefore filter rows with trade's time less or equal than its value.
	TimeBefore int64
}
