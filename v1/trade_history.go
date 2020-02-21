// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v1

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type tradeHistoryResponse struct {
	Success   int
	Error     string
	ErrorCode string
	Return    *listTrade
}

// listTrade contains list of trade history.
type listTrade struct {
	Trades []TradeHistory
}

//
// TradeHistory represent one transaction either for buy or for sell.
//
type TradeHistory struct {
	TradeID           int64
	OrderID           int64
	Type              string
	CurrencyName      string
	Amount            float64
	BaseCurrency      string
	BaseCurrencyPrice float64
	Price             float64
	Time              time.Time
	TimePrint         string
}

func (tradeHistory *TradeHistory) UnmarshalJSON(b []byte) (err error) {
	var kv map[string]interface{}

	err = json.Unmarshal(b, &kv)
	if err != nil {
		return err
	}

	for k, v := range kv {
		valstr := v.(string)

		k = strings.ToLower(k)
		switch k {
		case fieldNameTradeID:
			tradeHistory.TradeID, err = strconv.ParseInt(valstr, 10, 64)
		case fieldNameOrderID:
			tradeHistory.OrderID, err = strconv.ParseInt(valstr, 10, 64)
		case fieldNameType:
			tradeHistory.Type = valstr
		case fieldNameBaseCurrencyPrice:
			tradeHistory.BaseCurrencyPrice, err = strconv.ParseFloat(valstr, 64)
		case fieldNameBaseCurrency:
			tradeHistory.BaseCurrency = valstr
		case fieldNamePrice:
			tradeHistory.Price, err = strconv.ParseFloat(valstr, 64)
		case fieldNameTradeTime:
			ts, err := strconv.ParseInt(valstr, 10, 64)
			if err != nil {
				return err
			}
			tradeHistory.Time = time.Unix(ts, 0)
		case fieldNameTradeTimePrint:
			tradeHistory.TimePrint = valstr

		// The default key is the asset name with the key is the
		// amount of order.
		default:
			tradeHistory.CurrencyName = k
			tradeHistory.Amount, err = strconv.ParseFloat(valstr, 64)
		}
		if err != nil {
			return err
		}
	}

	return nil
}
