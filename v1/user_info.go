// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v1

import (
	"encoding/json"
	"strings"

	libjson "github.com/shuLhan/share/lib/json"
)

//
// userInfoResponse represent top level structure returned by call to
// "getInfo" API method.
//
type userInfoResponse struct {
	Success int
	Return  *UserInfo
}

//
// UserInfo contains the user information includings ID, email, name, and
// their assets: address, balance, and pending balance.
//
type UserInfo struct {
	ID               string `json:"user_id"`
	Name             string
	Email            string
	AssetAddress     map[string]string  `json:"address"`
	AssetBalance     map[string]float64 `json:"balance"`
	AssetBalanceHold map[string]float64 `json:"balance_hold"`
}

func (userInfo *UserInfo) UnmarshalJSON(b []byte) (err error) {
	var kv map[string]interface{}

	err = json.Unmarshal(b, &kv)
	if err != nil {
		return err
	}

	for k, v := range kv {
		k = strings.ToLower(k)

		switch k {
		case "user_id":
			userInfo.ID = v.(string)
		case "name":
			userInfo.Name = v.(string)
		case "email":
			userInfo.Email = v.(string)
		case "address":
			userInfo.unmarshalAddresses(v.(map[string]interface{}))
		case "balance":
			balances := v.(map[string]interface{})
			userInfo.AssetBalance, err = libjson.ToMapStringFloat64(balances)
		case "balance_hold":
			balancesHold := v.(map[string]interface{})
			userInfo.AssetBalanceHold, err = libjson.ToMapStringFloat64(balancesHold)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (userInfo *UserInfo) unmarshalAddresses(addresses map[string]interface{}) {
	userInfo.AssetAddress = make(map[string]string, len(addresses))
	for k, v := range addresses {
		k = strings.ToLower(k)
		valStr := v.(string)
		if len(valStr) > 0 {
			userInfo.AssetAddress[k] = v.(string)
		}
	}
}
