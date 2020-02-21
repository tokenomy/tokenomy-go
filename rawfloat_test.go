// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

import (
	"encoding/json"
	"testing"

	"github.com/shuLhan/share/lib/test"
)

type T struct {
	Rawfloat Rawfloat `json:"rawfloat"`
}

func TestRawfloat_MarshalJSON(t *testing.T) {
	cases := []struct {
		desc string
		t    T
		exp  string
	}{{
		desc: "With zero value",
		exp:  `{"rawfloat":0}`,
	}, {
		desc: "With trailing zero",
		t:    T{0.1},
		exp:  `{"rawfloat":0.1}`,
	}, {
		desc: "With mantissa zero 6 digits",
		t:    T{0.000_000_1},
		exp:  `{"rawfloat":0.0000001}`,
	}, {
		desc: "With mantissa zero 7 digits",
		t:    T{0.000_000_01},
		exp:  `{"rawfloat":0.00000001}`,
	}, {
		desc: "With mantissa zero 8 digits",
		t:    T{0.000_000_001},
		exp:  `{"rawfloat":0.000000001}`,
	}, {
		desc: "With all mantissa is zero",
		t:    T{0.000_000_000},
		exp:  `{"rawfloat":0}`,
	}, {
		desc: "With precision > maxPrecision (1)",
		t:    T{0.000_000_001},
		exp:  `{"rawfloat":0.000000001}`,
	}, {
		desc: "With precision > maxPrecision (2)",
		t:    T{0.000_000_0016},
		exp:  `{"rawfloat":0.000000001}`,
	}, {
		desc: "With precision > maxPrecision (3)",
		t:    T{0.000_000_001000},
		exp:  `{"rawfloat":0.000000001}`,
	}, {
		desc: "With no precisions",
		t:    T{123_456_789_0.0},
		exp:  `{"rawfloat":1234567890}`,
	}, {
		desc: "With base and mantissa",
		t:    T{64.23738872403},
		exp:  `{"rawfloat":64.23738872}`,
	}, {
		desc: "With precisions",
		t:    T{0.123_456_789_0},
		exp:  `{"rawfloat":0.12345678}`,
	}}

	for _, c := range cases {
		t.Log(c.desc)

		got, err := json.Marshal(&c.t)
		if err != nil {
			t.Fatal(err)
		}

		test.Assert(t, "MarshalJSON", c.exp, string(got), true)
	}
}
