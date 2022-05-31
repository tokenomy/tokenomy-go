// Copyright 2020 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/shuLhan/share/lib/math/big"
)

// TradeRequest contains parameters for trading.
type TradeRequest struct {
	Price  *big.Rat `json:"price"`
	Amount *big.Rat `json:"amount"`

	// Type of trade, its either "buy" or "sell".
	Type string `json:"type"`

	// Method of trading, its either "limit" or "market".
	// Default to "limit" if its empty.
	Method string `json:"method,omitempty"`

	// Pair name using "<coin>_<base>" format.
	Pair string `json:"pair"`

	// TimeInForce parameter only applicable if Method is "limit".
	// This option may change the behaviour of order "limit" processed by
	// broker.
	// Currently, the valid values are empty "" (default) or "FOK"
	// (fill-or-kill).
	//
	// If its empty, the order request processed normally as "limit"
	// request.
	//
	// If its "FOK", the order will be success only if only all of
	// requested amount is fulfilled, otherwise it will return as an error
	// ErrTradeFillOrKill.
	TimeInForce string `json:"time_in_force,omitempty"`

	// IsPostOnly parameter only applicable if Method is "limit".
	// If its true, the order will be success if only if no matching
	// trades happened, otherwise it will return an error.
	IsPostOnly bool `json:"post_only,omitempty"`
}

// Pack the TradeRequest object to be send by REST and/or WebSocket client.
func (treq *TradeRequest) Pack() (
	params url.Values, wsparams *WebSocketParams, err error,
) {
	params = url.Values{}
	if len(treq.Method) == 0 {
		treq.Method = TradeMethodLimit
	} else {
		treq.Method = strings.ToLower(treq.Method)
		switch treq.Method {
		case TradeMethodMarket, TradeMethodLimit:
		default:
			return nil, nil, ErrInvalidTradeMethod
		}
	}

	if treq.Amount == nil || treq.Amount.IsLessOrEqual(0) {
		return nil, nil, ErrInvalidAmount
	}

	params.Set(ParamNameTradeMethod, treq.Method)
	params.Set(ParamNamePair, treq.Pair)
	params.Set(ParamNameAmount, treq.Amount.String())

	wsparams = &WebSocketParams{
		TradeRequest: *treq,
	}

	if treq.Method == TradeMethodLimit {
		if treq.Price == nil || treq.Price.IsLessOrEqual(0) {
			return nil, nil, ErrInvalidPrice
		}
		params.Set(ParamNamePrice, treq.Price.String())
	}

	params.Set(ParamNamePostOnly, fmt.Sprintf("%t", treq.IsPostOnly))

	return params, wsparams, nil
}
