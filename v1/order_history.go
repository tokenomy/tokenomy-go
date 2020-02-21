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

type orderHistoryResponse struct {
	Success   int
	Error     string
	ErrorCode string
	Return    *listOrders
}

type listOrders struct {
	Orders []*OrderHistory
}

//
// OrderHistory contains a history of order of user.
//
type OrderHistory struct {
	PairName    string
	OrderID     int64
	Type        string
	Price       float64
	SubmitAt    time.Time
	FinishAt    time.Time
	Method      string
	Status      string
	OrderPrice  float64
	AssetRemain float64
}

func (oh *OrderHistory) UnmarshalJSON(b []byte) (err error) {
	var kv map[string]interface{}

	err = json.Unmarshal(b, &kv)
	if err != nil {
		return err
	}

	for k, v := range kv {
		var i64 int64
		valstr := v.(string)

		k = strings.ToLower(k)
		switch k {
		case fieldNameOrderID:
			oh.OrderID, err = strconv.ParseInt(valstr, 10, 64)
		case fieldNameType:
			oh.Type = valstr
		case fieldNamePrice:
			oh.Price, err = strconv.ParseFloat(valstr, 64)
		case fieldNameSubmitTime:
			i64, err = strconv.ParseInt(valstr, 10, 64)
			if err != nil {
				return err
			}
			oh.SubmitAt = time.Unix(i64, 0)
		case fieldNameFinishTime:
			i64, err = strconv.ParseInt(valstr, 10, 64)
			if err != nil {
				return err
			}
			oh.FinishAt = time.Unix(i64, 0)
		case fieldNameMethod:
			oh.Method = valstr
		case fieldNameStatus:
			oh.Status = valstr
		default:
			switch {
			case strings.HasPrefix(k, "order_"):
				oh.OrderPrice, err = strconv.ParseFloat(valstr, 64)
			case strings.HasPrefix(k, "remain_"):
				oh.AssetRemain, err = strconv.ParseFloat(valstr, 64)
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}
