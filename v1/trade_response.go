// Copyright 2019 Tokenomy Technologies Pte. Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v1

import (
	"encoding/json"
	"strconv"
	"strings"
)

//
// tradeResponse is data returned from buying, selling, or canceling order.
//
type tradeResponse struct {
	OrderID   int64
	Pair      string
	Success   int64
	Error     string
	ErrorCode string
	Receive   float64
	Spend     float64
	Remain    float64
	Balance   map[string]float64
	IsError   bool
}

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
			tres.Balance, err = jsonToMapStringFloat64(v.(map[string]interface{}))
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
			case strings.HasPrefix(k, "spend_"):
				tres.Spend, err = strconv.ParseFloat(v.(string), 64)
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
