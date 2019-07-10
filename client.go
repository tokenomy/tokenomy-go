// Copyright 2019 Tokenomy Technologies Pte. Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

//
// client represent common fields and methods across different API versions.
//
type client struct {
	conn        *http.Client
	env         *environment
	baseHost    string
	privatePath string
	token       string
	secret      string
}

//
// callPrivate call private API with specific method and parameters.
// On success it will return response body.
// On fail it will return an empty body with an error.
//
func (cl *client) callPrivate(method string, params url.Values) (
	body []byte, err error,
) {
	if len(cl.token) == 0 {
		return nil, ErrUnauthenticated
	}

	req, err := cl.newPrivateRequest(method, params)
	if err != nil {
		err = fmt.Errorf("callPrivate: " + err.Error())
		return nil, err
	}

	if cl.env.debug >= 2 {
		fmt.Printf("<<< callPrivate: request: %+v\n", req)
	}

	res, err := cl.conn.Do(req)
	if err != nil {
		err = fmt.Errorf("callPrivate: " + err.Error())
		return nil, err
	}

	body, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		err = fmt.Errorf("callPrivate: " + err.Error())
		return nil, err
	}

	return body, nil
}

//
// callPublic call public API with specific path.
// On success, it will return HTTP response body.
// On fail, it will return an empty body and an error.
//
func (cl *client) callPublic(publicPath string) (body []byte, err error) {
	req := &http.Request{
		Method: http.MethodPost,
		Header: http.Header{
			"Content-Type": []string{
				"application/x-www-form-urlencoded",
			},
		},
	}

	req.URL, err = url.Parse(cl.baseHost + publicPath)
	if err != nil {
		return nil, fmt.Errorf("callPublic: " + err.Error())
	}

	res, err := cl.conn.Do(req)
	if err != nil {
		err = fmt.Errorf("callPublic: " + err.Error())
		return nil, err
	}

	body, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		err = fmt.Errorf("callPublic: " + err.Error())
		return nil, err
	}

	return body, nil
}

//
// cancelOrder cancel buy or sell order type on spesific pair by order ID.
//
func (cl *client) cancelOrder(tipe, pairName string, orderID int64) (
	body []byte, err error,
) {
	orderIDStr := strconv.FormatInt(orderID, 10)

	params := map[string][]string{
		"type":     {tipe},
		"pair":     {pairName},
		"order_id": {orderIDStr},
	}

	body, err = cl.callPrivate(apiTradeCancelOrder, params)
	if err != nil {
		return nil, err
	}

	return body, nil
}

//
// encodeToString return the hex encoding of data hashed with client's API
// secret key.
//
func (cl *client) encodeToString(in []byte) string {
	mac := hmac.New(sha512.New, []byte(cl.secret))

	_, _ = mac.Write(in)

	return hex.EncodeToString(mac.Sum(nil))
}

//
// newPrivateRequest create a new HTTP request for private API with specific
// "method" and custom parameters "params" as form values in POST body.
//
func (cl *client) newPrivateRequest(apiMethod string, params url.Values) (
	req *http.Request, err error,
) {
	q := url.Values{
		"nonce": []string{
			timestampAsString(),
		},
		"method": []string{
			apiMethod,
		},
	}

	vparams := map[string][]string(params)
	for k, v := range vparams {
		if len(v) > 0 {
			q.Set(k, v[0])
		}
	}

	reqBody := q.Encode()

	if cl.env.debug >= 2 {
		fmt.Printf("<<< newPrivateRequest: request body: %s\n", reqBody)
	}

	sign := cl.encodeToString([]byte(reqBody))

	req = &http.Request{
		Method: http.MethodPost,
		Header: http.Header{
			"Content-Type": []string{
				"application/x-www-form-urlencoded",
			},
			"Key": []string{
				cl.token,
			},
			"Sign": []string{
				sign,
			},
		},
		Body: ioutil.NopCloser(strings.NewReader(reqBody)),
	}

	req.URL, err = url.Parse(cl.baseHost + cl.privatePath)
	if err != nil {
		err = fmt.Errorf("newPrivateRequest: " + err.Error())
		return nil, err
	}

	return req, nil
}

func (cl *client) trade(method, tipe, pair string, amount, price float64) (
	body []byte, err error,
) {
	assetBase := strings.Split(pair, "_")
	if len(assetBase) != 2 {
		return nil, fmt.Errorf("trade: invalid pair name %q", pair)
	}

	amountStr := strconv.FormatFloat(amount, 'f', 8, 64)
	priceStr := strconv.FormatFloat(price, 'f', 8, 64)

	params := map[string][]string{
		"order_method": {method},
		"pair":         {pair},
		"price":        {priceStr},
		"type":         {tipe},
	}

	params[assetBase[0]] = []string{amountStr}

	body, err = cl.callPrivate(apiTrade, params)
	if err != nil {
		return nil, err
	}

	return body, nil
}
