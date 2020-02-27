// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v1

import (
	"encoding/json"
	"strconv"
	"strings"

	libjson "github.com/shuLhan/share/lib/json"
)

type cancelOrderResponse struct {
	Success int          `json:"success"`
	Error   string       `json:"error,omitempty"`
	Return  *CancelOrder `json:"return"`
}

//
// CancelOrder contains a success response from calling a "cancelOrder"
// method.
//
type CancelOrder struct {
	OrderID  int64              `json:"order_id"`
	Type     string             `json:"type"`
	Pair     string             `json:"pair"`
	Balances map[string]float64 `json:"balance"`
}

func (cancelOrder *CancelOrder) UnmarshalJSON(b []byte) (err error) {
	var kv map[string]interface{}

	err = json.Unmarshal(b, &kv)
	if err != nil {
		return err
	}

	for k, v := range kv {
		valstr, _ := v.(string)

		k = strings.ToLower(k)
		switch k {
		case fieldNameOrderID:
			cancelOrder.OrderID, err = strconv.ParseInt(valstr, 10, 64)
		case fieldNameType:
			cancelOrder.Type = valstr
		case fieldNamePair:
			cancelOrder.Pair = valstr
		case fieldNameBalance:
			balances := v.(map[string]interface{})
			cancelOrder.Balances, err = libjson.ToMapStringFloat64(balances)
		}
		if err != nil {
			return err
		}
	}

	return nil
}
