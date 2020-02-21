#  tokenomy-go

[![GoDoc](https://godoc.org/github.com/tokenomy/tokenomy-go?status.svg)](https://godoc.org/github.com/tokenomy/tokenomy-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/tokenomy/tokenomy-go)](https://goreportcard.com/report/github.com/tokenomy/tokenomy-go)

This is the official Go module for client of Tokenomy API v1 and v2.

Documentation for API v1 is available at
[API v1 help page](https://exchange.tokenomy.com/help/api) and the
[Go doc page](https://pkg.go.dev/github.com/tokenomy/tokenomy-go/v1?tab=doc).

Documentation for API v2 is available at
[API v2 help page](https://exchange.tokenomy.com/help/api/v2) and the
[Go doc page](https://pkg.go.dev/github.com/tokenomy/tokenomy-go/v2?tab=doc).

Note that, this module is in development state and subject to changes in the
future.

Tokenomy provide public and private APIs.
The public APIs can be accessed directly, without any keys or credential.
The private APIs can only be accessed by using token and secret keys (API
credential).

An API credential can be retrieved manually by logging in into your
[Tokenomy Exchange account](https://exchange.tokenomy.com) and open the
["Trade API" menu section](https://exchange.tokenomy.com/trade_api).
Please keep these credentials safe and do not reveal it to any external party.


#  License

```
Copyright (c) 2019 Tokenomy Technologies Ltd. All rights reserved.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
