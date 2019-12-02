// Copyright 2019 Tokenomy Technologies Pte. Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v2

import liberrors "github.com/shuLhan/share/lib/errors"

//
// Response contains the HTTP response from server.
//
type Response struct {
	liberrors.E
	Data interface{} `json:"data,omitempty"`
}
