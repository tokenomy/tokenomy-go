// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

import (
	"fmt"
	"os"
	"strconv"
)

// Environment contains default and dynamic values that gathered from external
// resources, for example system environment variables.
type Environment struct {
	// Address of API server, optional. It will default to DefaultAddress
	// constant on each package.
	Address string

	// Token, required, is the public part of API key.
	Token string

	// Secret, required, is the private part of API key.
	Secret string

	//
	// Debug define level of logging in our library.
	// debug value is set from environment variable "TOKENOMY_DEBUG".
	// TOKENOMY_DEBUG=1 is for logging in configuration.
	// TOKENOMY_DEBUG=2 is for logging input and output.
	//
	Debug int

	// IsInsecure, optional, allow self-signed certificate, should be use
	// for testing only.
	IsInsecure bool
}

// NewEnvironment create and initialize environment.
//
// If token and/or secret is empty it will set from environment variables
// TOKENOMY_TOKEN and TOKENOMY_SECRET.
func NewEnvironment(token, secret string) (env *Environment) {
	env = &Environment{
		Address: os.Getenv(EnvNameAddress),
		Token:   os.Getenv(EnvNameToken),
		Secret:  os.Getenv(EnvNameSecret),
	}

	if len(token) > 0 {
		env.Token = token
	}
	if len(secret) > 0 {
		env.Secret = secret
	}

	v := os.Getenv(EnvNameDebug)
	if len(v) > 0 {
		env.Debug, _ = strconv.Atoi(v)
	}

	if env.Debug >= 1 {
		fmt.Printf(">>> Environment: %+v\n", env)
	}

	return env
}
