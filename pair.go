// Copyright 2019 Tokenomy Technologies Pte. Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

import (
	"encoding/json"
	"fmt"
	"strings"
)

//
// Pair contains a pair information.
//
type Pair struct {
	Name         string
	HighestPrice string
	LowestPrice  string
	AssetVolume  string
	BaseVolume   string
	LastPrice    string
	BuyPrice     string
	SellPrice    string
	volumes      map[string]string
}

//
// tickerResponse is wrapper for unwrapping a single pair from GetTicker
// JSON response.
//
type tickerResponse struct {
	Ticker *Pair
}

func (pair *Pair) UnmarshalJSON(b []byte) (err error) {
	var kv map[string]interface{}

	err = json.Unmarshal(b, &kv)
	if err != nil {
		return err
	}

	for k, v := range kv {
		k = strings.ToLower(k)

		valStr, ok := v.(string)
		if !ok {
			valStr = fmt.Sprintf("%v", v)
		}

		switch k {
		case "name":
			pair.Name = valStr
		case "high":
			pair.HighestPrice = valStr
		case "low":
			pair.LowestPrice = valStr
		case "last":
			pair.LastPrice = valStr
		case "buy":
			pair.BuyPrice = valStr
		case "sell":
			pair.SellPrice = valStr
		default:
			if !strings.HasPrefix(k, "vol_") {
				continue
			}

			volName := strings.Split(k, "_")
			if len(volName) != 2 {
				continue
			}
			if pair.volumes == nil {
				pair.volumes = make(map[string]string)
			}

			pair.volumes[volName[1]] = valStr
		}
	}

	return nil
}

//
// propagate the asset and base volumes from value of "vol_XXX" key.
//
func (pair *Pair) propagate(pairName string) {
	pairName = strings.ToLower(pairName)
	assetBase := strings.Split(pairName, "_")
	if len(assetBase) != 2 {
		return
	}

	pair.AssetVolume = pair.volumes[assetBase[0]]
	pair.BaseVolume = pair.volumes[assetBase[1]]
}
