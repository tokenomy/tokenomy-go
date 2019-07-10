// Copyright 2019 Tokenomy Technologies Pte. Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

import (
	"fmt"
	"os"
	"strconv"
)

//
// environment contains default and dynamic values that gathered from external
// resources, for example system environment variables.
//
type environment struct {
	//
	// debug define level of logging in our library.
	// debug value is set from environment variable "DEBUG".
	// DEBUG=1 is for logging in configuration.
	// DEBUG=2 is for logging input and output.
	//
	debug         int
	v1BaseHost    string
	v1PrivatePath string
	v2BaseHost    string
	v2PrivatePath string
	apiKey        string
	apiSecret     string
	name          string
}

func newEnvironment() (env *environment) {
	env = &environment{
		v1BaseHost:    "https://exchange.tokenomy.com",
		v1PrivatePath: "/tapi",
		v2BaseHost:    "https://api.tokenomy.com",
		v2PrivatePath: "/v2/private",
		apiKey:        os.Getenv("TOKENOMY_KEY"),
		apiSecret:     os.Getenv("TOKENOMY_SECRET"),
		name:          os.Getenv("TOKENOMY_ENV"),
	}

	v := os.Getenv("DEBUG")
	if len(v) > 0 {
		env.debug, _ = strconv.Atoi(v)
	}

	baseHost := os.Getenv("TOKENOMY_HOST")
	if len(baseHost) > 0 {
		env.v1BaseHost = "https://" + baseHost
	}

	baseHost = os.Getenv("TOKENOMY_HOSTv2")
	if len(baseHost) > 0 {
		env.v2BaseHost = "https://" + baseHost
	}

	if env.debug >= 1 {
		fmt.Printf("=== Environment: %+v", env)
	}

	return env
}
