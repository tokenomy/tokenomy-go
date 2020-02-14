// Copyright 2019 Tokenomy Technologies Pte. Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v1

import (
	"encoding/json"
	"strconv"
	"strings"

	libjson "github.com/shuLhan/share/lib/json"
)

//
// tradeResponse is the data returned from buying, selling, or canceling
// order.
// This is the internal struct for v1 that will be converted to common struct
// tokenomy.TradeResponse.
//
type tradeResponse struct {
	OrderID   int64
	Pair      string
	Success   int64
	Error     string
	ErrorCode string
	Receive   float64
	Filled    float64
	Remain    float64
	Balance   map[string]float64
	IsError   bool
}

//
// UnmarshalJSON response from trading bid and ask.
//
// For example, lets say that we are trading using pair "ten_btc", with amount
// of "10" (TEN) and price (0.000_02).
// Here is an example of completed ask response,
//
//	{
//		"is_error": false,
//		"success": 1,
//		"sold_ten": "10.00000000",
//		"receive_btc": "0.00004500",
//		"remain_ten": "0.00000000",
//		"order_id": 6657492,
//		"balance": {
//			"btc": "9.99382081",
//			"frozen_btc": "0.00030000",
//			"frozen_ten": "16.50000000",
//			"ten": "8689.14604860",
//			...
//		}
//	}
//
// The "sold_<coin>" define how many amount have been filled in market.
//
// Here is an example of completed bid response,
//
//	{
//		"is_error": false,
//		"success": 1,
//		"receive_ten": "9.99999999",
//		"spend_btc": "0.00020000",
//		"remain_btc": "0.00000000",
//		"order_id": 6658409,
//		"balance": {
//			"btc": "9.99382081",
//			"frozen_btc": "0.00030000",
//			"frozen_ten": "16.50000001",
//			"ten": "8689.13354860",
//			...
//		}
//	}
//
func (tres *tradeResponse) UnmarshalJSON(b []byte) (err error) {
	var kv map[string]interface{}

	err = json.Unmarshal(b, &kv)
	if err != nil {
		return err
	}

	for k, v := range kv {
		k = strings.ToLower(k)

		switch k {
		case fieldNameOrderID:
			tres.OrderID = int64(v.(float64))
		case fieldNameIsError:
			tres.IsError = v.(bool)
		case fieldNameSuccess:
			tres.Success = int64(v.(float64))
		case fieldNameError:
			tres.Error = v.(string)
		case fieldNameErrorCode:
			tres.ErrorCode = v.(string)
		case fieldNameBalance:
			tres.Balance, err = libjson.ToMapStringFloat64(v.(map[string]interface{}))
			if err != nil {
				return err
			}
		default:
			switch {
			case strings.HasPrefix(k, "receive_"):
				tres.Receive, err = strconv.ParseFloat(v.(string), 64)
				if err != nil {
					return err
				}

			case strings.HasPrefix(k, "sold_"):
				tres.Filled, err = strconv.ParseFloat(v.(string), 64)
				if err != nil {
					return err
				}

			case strings.HasPrefix(k, "spend_"):
				tres.Filled, err = strconv.ParseFloat(v.(string), 64)
				if err != nil {
					return err
				}

			case strings.HasPrefix(k, "remain_"):
				tres.Remain, err = strconv.ParseFloat(v.(string), 64)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
