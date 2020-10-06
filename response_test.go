// Copyright 2020 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

import (
	"encoding/gob"
	"testing"

	"github.com/shuLhan/share/lib/test"
)

type T struct {
	String string
}

func TestResponse_Pack(t *testing.T) {
	gob.Register(T{})

	exp := &Response{
		Data: T{
			String: "abc",
		},
	}

	b, err := exp.PackGob()
	if err != nil {
		t.Fatal(err)
	}

	got := &Response{}
	err = got.UnpackGob(b)
	if err != nil {
		t.Fatal(err)
	}

	test.Assert(t, "Response.Pack", exp, got, true)
}
