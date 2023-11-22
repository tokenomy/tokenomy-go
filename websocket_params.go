// Copyright 2020 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

import (
	"encoding/json"
)

// WebSocketParams contains the request parameters for WebSocket client.
type WebSocketParams struct {
	Address     string `json:"address,omitempty"`
	AddressType string `json:"address_type,omitempty"`
	Asset       string `json:"asset,omitempty"`
	IDSortBy    string `json:"id_sort_by,omitempty"`
	Memo        string `json:"memo,omitempty"`
	Network     string `json:"network,omitempty"`
	RequestID   string `json:"request_id,omitempty"`

	TradeRequest
	PublicSubscription

	IDAfter  int64 `json:"id_after,omitempty"`
	IDBefore int64 `json:"id_before,omitempty"`

	TimeAfter  int64 `json:"time_after,omitempty"`
	TimeBefore int64 `json:"time_before,omitempty"`

	TradeID int64 `json:"trade_id,omitempty"`

	Limit  int64 `json:"limit,omitempty"`
	Offset int64 `json:"offset,omitempty"`
}

// Pack the WebSocket parameters as JSON.
func (wsparams *WebSocketParams) Pack() (b []byte, err error) {
	return json.Marshal(wsparams)
}

// Unpack the parameters from JSON bytes.
func (wsparams *WebSocketParams) Unpack(b []byte) (err error) {
	return json.Unmarshal(b, wsparams)
}
