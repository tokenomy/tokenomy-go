// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v1

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/shuLhan/share/lib/math/big"
)

//
// Trade contains trade information for a pair.
//
type Trade struct {
	ID     int64
	Type   string
	Date   time.Time
	Amount *big.Rat
	Price  *big.Rat
}

func (trade *Trade) UnmarshalJSON(b []byte) (err error) {
	var kv map[string]string

	err = json.Unmarshal(b, &kv)
	if err != nil {
		return err
	}

	for k, v := range kv {
		k = strings.ToLower(k)

		switch k {
		case fieldNameDate:
			ts, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return err
			}
			trade.Date = time.Unix(ts, 0)
		case fieldNamePrice:
			trade.Price = big.NewRat(v)
			if trade.Price == nil {
				err = fmt.Errorf("invalid price value %q", v)
			}
		case fieldNameAmount:
			trade.Amount = big.NewRat(v)
			if trade.Amount == nil {
				err = fmt.Errorf("invalid amount value %q", v)
			}
		case fieldNameTID:
			trade.ID, err = strconv.ParseInt(v, 10, 64)
		case fieldNameType:
			trade.Type = v
		}
		if err != nil {
			return err
		}
	}

	return nil
}
