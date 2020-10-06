// Copyright 2020 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

import (
	"os"
	"testing"
)

func TestWebSocketPrivate_UserInfo(t *testing.T) {
	if len(os.Getenv(EnvNameTestE2E)) == 0 {
		t.Skipf("%s is not set, skipping ...", EnvNameTestE2E)
	}
	env := &Environment{
		Address: os.Getenv(EnvNameAddress),
		Token:   os.Getenv(EnvNameToken),
		Secret:  os.Getenv(EnvNameSecret),
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
