// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v2

import (
	"bytes"
	"encoding/gob"
	"fmt"

	liberrors "github.com/shuLhan/share/lib/errors"
)

//
// Response contains the HTTP response from server.
//
type Response struct {
	liberrors.E
	Data interface{} `json:"data,omitempty"`
}

//
// PackGob pack the response into raw bytes using gob format.
//
func (res *Response) PackGob() (b []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(res)
	if err != nil {
		return nil, fmt.Errorf("Response.Pack: %w", err)
	}
	return buf.Bytes(), nil
}

//
// UnpackGob unpack the gob formatted raw bytes into Response.
//
func (res *Response) UnpackGob(b []byte) (err error) {
	dec := gob.NewDecoder(bytes.NewReader(b))
	res.Data = make([]byte, 0)
	err = dec.Decode(res)
	if err != nil {
		return fmt.Errorf("Response.Unpack: %w", err)
	}
	return nil
}
