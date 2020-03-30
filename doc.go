// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

//
// Package tokenomy is the official Go module for client of Tokenomy API v1
// and v2.
//
// Documentation for API v1 is available at
// https://exchange.tokenomy.com/help/api and the
// Go doc page at
// https://pkg.go.dev/github.com/tokenomy/tokenomy-go/v1?tab=doc.
//
// Documentation for API v2 is available at
// https://exchange.tokenomy.com/help/api/v2 and the
// https://pkg.go.dev/github.com/tokenomy/tokenomy-go/v2?tab=doc.
//
// Note that, this module is in development state, still in v0, and may
// subject to changes in the future release until v1 is reached.
//
// Tokenomy provide public and private APIs.
// The public APIs can be accessed directly, without any keys or credential.
// The private APIs can only be accessed by using token and secret keys (API
// credential).
//
// An API credential can be retrieved manually by logging in into your
// Tokenomy Exchange account (https://exchange.tokenomy.com) and open the
// "Trade API" menu section (https://exchange.tokenomy.com/trade_api).
// Please keep these credentials safe and do not reveal it to any external
// party.
//
package tokenomy
