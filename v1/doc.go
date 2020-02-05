// Copyright 2019 Tokenomy Technologies Pte. Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

//
// Package v1 provide a library for accesing Tokenomy API v1 (see
// https://exchange.tokenomy.com/help/api for API documentation).
//
// Tokenomy provide public and private APIs.
// The public APIs can be accessed directly by creating new client with empty
// Token and Secret.
// The private APIs can only be accessed by using API key, a pair of Token and
// Secret.
//
// An API credential can be retrieved manually by logging in into your
// Tokenomy Exchange account (https://exchange.tokenomy.com) and open the
// "Trade API" menu section or https://exchange.tokenomy.com/trade_api.
// Please keep those credentials safe and do not reveal to any external party.
//
// Beside passing the API key Token and Secret to NewEnvironment, this
// library also read Token and Secret keys from environment variables
// "TOKENOMY_TOKEN" for Token and "TOKENOMY_SECRET" for Secret.
//
package v1
