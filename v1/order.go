// Copyright 2019 Tokenomy Technologies Pte. Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v1

import (
	"encoding/json"
	"fmt"
	"strconv"
)

//
// Order contains the number of amount and price of open order.
//
type Order struct {
	Amount float64
	Price  float64
}

func (order *Order) UnmarshalJSON(b []byte) (err error) {
	var vals []string

	err = json.Unmarshal(b, &vals)
	if err != nil {
		return err
	}

	if len(vals) != 2 {
		return fmt.Errorf("order: UnmarshalJSON: invalid length of order")
	}

	order.Price, err = strconv.ParseFloat(vals[0], 64)
	if err != nil {
		return err
	}

	order.Amount, err = strconv.ParseFloat(vals[1], 64)
	if err != nil {
		return err
	}

	return nil
}
