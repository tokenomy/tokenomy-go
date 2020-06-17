// Copyright 2020 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v2

import (
	"os"
	"testing"

	"github.com/tokenomy/tokenomy-go"
)

func TestWebSocketPrivate_UserInfo(t *testing.T) {
	if len(os.Getenv(envTestE2E)) == 0 {
		t.Skipf("%s is not set, skipping ...", envTestE2E)
	}
	env := &tokenomy.Environment{
		Address: os.Getenv(envAddress),
		Token:   os.Getenv(envToken),
		Secret:  os.Getenv(envSecret),
	}

	cl, err := NewWebSocketPrivate(env)
	if err != nil {
		t.Fatal(err)
	}

	userInfo, err := cl.UserInfo()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("UserInfo: %+v\n", userInfo)
}
