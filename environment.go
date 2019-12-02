// Copyright 2019 Tokenomy Technologies Pte. Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

//
// Environment contains default and dynamic values that gathered from external
// resources, for example system environment variables.
//
type Environment struct {
	//
	// Debug define level of logging in our library.
	// debug value is set from environment variable "TOKENOMY_DEBUG".
	// TOKENOMY_DEBUG=1 is for logging in configuration.
	// TOKENOMY_DEBUG=2 is for logging input and output.
	//
	Debug int

	// Host contains the host API.
	// Its value can be overriden using environment variable
	// "TOKENOMY_HOST".
	Host      string
	APIKey    string
	SecretKey string
}

//
// NewEnvironment create and initialize environment.
//
func NewEnvironment(host string) (env *Environment) {
	log.SetFlags(0)

	env = &Environment{
		Host:      host,
		APIKey:    os.Getenv(EnvNameKey),
		SecretKey: os.Getenv(EnvNameSecret),
	}

	v := os.Getenv(EnvNameDebug)
	if len(v) > 0 {
		env.Debug, _ = strconv.Atoi(v)
	}

	v = os.Getenv(EnvNameHost)
	if len(v) > 0 {
		env.Host = v
	}

	if env.Debug >= 1 {
		fmt.Printf(">>> Environment: %+v\n", env)
	}

	return env
}
