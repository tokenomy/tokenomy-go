// Copyright 2020 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

import (
	"testing"

	"github.com/shuLhan/share/lib/test"
)

func TestSign(t *testing.T) {
	cases := []struct {
		payload string
		secret  string
		exp     string
	}{{
		payload: "timestamp=1574423788&pair=ten_btc",
		secret:  "secr3t",
		exp:     "db068236b2cbc0084946de7be9dce15f2ac271ddae83e6d9181f25b397d09f10d128f4e710dbf1aa7b15c13bb2032b9673d549829e7455fe3ef0ddb95a0dc1a5",
	}, {
		payload: "timestamp=1574423788&pair=ten_btc&trade_id=1",
		secret:  "secr3t",
		exp:     "5befb7f8236bf55c685c2b163e9f755c7dc6fd29c64cf30bfba5820917221dce5e5b8051216c0345f5cccd704f4a351a7b4374fc19959572b087ad6213760dc0",
	}}

	for _, c := range cases {
		got := Sign(c.payload, c.secret)

		test.Assert(t, "Sign", c.exp, got)
	}
}
