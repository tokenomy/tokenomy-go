// Copyright 2019 Tokenomy Technologies Pte. Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

//
// Package tokenomy provide a library for accesing Tokenomy API v1 (see
// https://exchange.tokenomy.com/help/api for HTTP API documentation).
//
// Tokenomy provide public and private APIs.
// The public APIs can be accessed directly by creating new client with empty
// token and secret parameters.
// The private APIs can only be accessed by using token and secret keys (API
// credential).
//
// An API credential can be retrieved manually by logging in into your
// Tokenomy Exchange account (https://exchange.tokenomy.com) and open the
// "Trade API" menu section or https://exchange.tokenomy.com/trade_api.
// Please keep these credentials safe and do not reveal to any external party.
//
// Beside passing the token and secret to NewClient or Authenticate, this
// library also read token and secret values from environment variables
// "TOKENOMY_KEY" for token and "TOKENOMY_SECRET" for secret.
//
package tokenomy
