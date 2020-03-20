// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v1

import (
	"encoding/json"
	"fmt"

	"github.com/shuLhan/share/lib/math/big"
)

//
// Order contains the number of amount and price of open order.
//
type Order struct {
	Amount *big.Rat
	Price  *big.Rat
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

	order.Price = big.NewRat(vals[0])
	if order.Price == nil {
		return fmt.Errorf("order: invalid price value %q", vals[0])
	}

	order.Amount = big.NewRat(vals[1])
	if order.Amount == nil {
		return fmt.Errorf("order: invalid amount value %q", vals[1])
	}

	return nil
}
