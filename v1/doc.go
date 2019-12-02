// Copyright 2019 Tokenomy Technologies Pte. Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

//
// Package v1 provide a library for accesing Tokenomy API v1 (see
// https://exchange.tokenomy.com/help/api for API documentation).
//
// Tokenomy provide public and private APIs.
// The public APIs can be accessed directly by creating new client with empty
// API and secret key parameters.
// The private APIs can only be accessed by using API and secret keys (API
// credential).
//
// An API credential can be retrieved manually by logging in into your
// Tokenomy Exchange account (https://exchange.tokenomy.com) and open the
// "Trade API" menu section or https://exchange.tokenomy.com/trade_api.
// Please keep these credentials safe and do not reveal to any external party.
//
// Beside passing the API and secret keys to NewClient or Authenticate, this
// library also read API and secret keys from environment variables
// "TOKENOMY_KEY" for API key and "TOKENOMY_SECRET" for secret key.
//
package v1
