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

//
// TradeRequest contains parameters for trading.
//
type TradeRequest struct {
	// Type of trade, its either "buy" or "sell".
	Type string `json:"type"`

	// Method of trading, its either "limit" or "market".
	Method string `json:"method"`

	// Pair name using "<coin>_<base>" format.
	Pair string `json:"pair"`

	Amount *big.Rat `json:"amount"`
	Price  *big.Rat `json:"price"`

	// PostOnly parameter only applicable if Method is "limit".
	// If its true, the order will be success if only if no matching
	// trades happened, otherwise it will return an error.
	PostOnly bool `json:"post_only"`
}

//
// Pack the TradeRequest object to be send by REST and/or WebSocket client.
//
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

	params.Set(ParamNamePostOnly, fmt.Sprintf("%t", treq.PostOnly))

	return params, wsparams, nil
}
