// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v2

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	liberrors "github.com/shuLhan/share/lib/errors"
	"github.com/shuLhan/share/lib/math/big"

	"github.com/tokenomy/tokenomy-go"
)

//
// Client for Tokenomy REST API v2.
//
type Client struct {
	User *tokenomy.User
	conn *http.Client
	env  *tokenomy.Environment
}

//
// NewClient create and initialize new client for REST API v2.
//
// The Environment Address parameter define the REST API v2 address, if its
// empty it will set to value in DefaultAddress.
//
// The Environment' Token and Secret parameters are used to authenticate the
// client when accessing private API.
//
// By default, the Token and Secret is read from environment variables
// "TOKENOMY_TOKEN" and "TOKENOMY_SECRET", the parameters will override the
// default value, if its set.
// If both environment variables and the parameters are empty, the client can
// only access the public API.
//
func NewClient(env *tokenomy.Environment) (cl *Client, err error) {
	if len(env.Address) == 0 {
		env.Address = DefaultAddress
	}

	transport := http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: env.IsInsecure, //nolint: gosec
		},
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	cl = &Client{
		conn: &http.Client{
			Transport: &transport,
		},
		env: env,
	}

	if len(cl.env.Token) > 0 {
		err = cl.Authenticate()
	}

	return cl, err
}

//
// Authenticate the current client's connection using token and secret keys.
//
func (cl *Client) Authenticate() (err error) {
	// Test the token and secret keys by requesting user information.
	cl.User, err = cl.UserInfo()
	if err != nil {
		return fmt.Errorf("Authenticate: %w", err)
	}

	return nil
}

//
// MarketDepths fetch list of market's depth for specific pair.
//
func (cl *Client) MarketDepths(pairName string) (depths *MarketDepths, err error) {
	params := url.Values{
		tokenomy.ParamNamePair: []string{pairName},
	}

	b, err := cl.doGet(apiMarketDepths, params)
	if err != nil {
		return nil, fmt.Errorf("MarketDepths: %w", err)
	}

	depths = &MarketDepths{}
	res := &Response{
		Data: depths,
	}

	err = json.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}

	return depths, nil
}

//
// MarketInfo return information about all the pair in the platform.
//
func (cl *Client) MarketInfo() (marketInfos []MarketInfo, err error) {
	b, err := cl.doGet(apiMarketInfo, url.Values{})
	if err != nil {
		return nil, fmt.Errorf("MarketInfo: %w", err)
	}

	marketInfos = make([]MarketInfo, 0)
	res := &Response{
		Data: marketInfos,
	}

	err = json.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}

	return marketInfos, nil
}

//
// MarketTradesOpen return list of all open trades in the market, specific to
// pair's name, grouped by ask and bid.
//
func (cl *Client) MarketTradesOpen(pairName string) (openTrades *TradeOpens, err error) {
	params := url.Values{
		tokenomy.ParamNamePair: []string{pairName},
	}

	b, err := cl.doGet(apiMarketTradesOpen, params)
	if err != nil {
		return nil, fmt.Errorf("MarketTradesOpen: %w", err)
	}

	openTrades = &TradeOpens{}
	res := &Response{
		Data: openTrades,
	}

	err = json.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}

	return openTrades, nil
}

//
// MarketPrices return list of all latest pair's prices.
//
func (cl *Client) MarketPrices() (marketPrices MarketPrices, err error) {
	params := url.Values{}

	b, err := cl.doGet(apiMarketPrices, params)
	if err != nil {
		return nil, fmt.Errorf("MarketPrices: %w", err)
	}

	marketPrices = make(MarketPrices)
	res := &Response{
		Data: marketPrices,
	}

	err = json.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}

	return marketPrices, nil
}

//
// MarketTicker return the ticker information on specific pair.
//
func (cl *Client) MarketTicker(pairName string) (tick *Tick, err error) {
	params := url.Values{
		tokenomy.ParamNamePair: []string{pairName},
	}

	b, err := cl.doGet(apiMarketTicker, params)
	if err != nil {
		return nil, fmt.Errorf("MarketTicker: %w", err)
	}

	tick = &Tick{}
	res := &Response{
		Data: tick,
	}

	err = json.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}

	return tick, nil
}

//
// MarketTrades return list of all closed trades in the market, specific to
// pair's name, grouped by ask and bid.
//
func (cl *Client) MarketTrades(pairName string) (tradePrices *MarketTradePrices, err error) {
	params := url.Values{
		tokenomy.ParamNamePair: []string{pairName},
	}

	b, err := cl.doGet(apiMarketTrades, params)
	if err != nil {
		return nil, fmt.Errorf("MarketTrades: %w", err)
	}

	tradePrices = &MarketTradePrices{}
	res := &Response{
		Data: tradePrices,
	}

	err = json.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}

	return tradePrices, nil
}

//
// MarketSummaries return the summaries (ticker) of all pairs.
//
func (cl *Client) MarketSummaries() (summaries *MarketSummaries, err error) {
	params := url.Values{}

	b, err := cl.doGet(apiMarketSummaries, params)
	if err != nil {
		return nil, fmt.Errorf("MarketSummaries: %w", err)
	}

	summaries = &MarketSummaries{}
	res := &Response{
		Data: summaries,
	}

	err = json.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}

	return summaries, nil
}

//
// UserInfo fetch the user information and balances.
//
// This method require authentication.
//
func (cl *Client) UserInfo() (user *tokenomy.User, err error) {
	params := url.Values{}

	b, err := cl.doSecureRequest(http.MethodGet, apiUserInfo, params)
	if err != nil {
		return nil, fmt.Errorf("UserInfo: %w", err)
	}

	user = &tokenomy.User{}
	res := &Response{
		Data: user,
	}

	err = json.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}

	return user, nil
}

//
// UserTrades list the user's trade history, ordered from latest to oldest
// one.
//
// The offset parameter define the number of record to be skipped.
//
// The limit parameter define the maximum number of record fetched, if its not
// set default to DefaultLimit.
//
// The idAfter and idBefore filter the records based on ID.  The idAfter will
// only fetch  record after the value of ID, and idBefore will only fetch
// record before the value of ID.
//
// the timeAfter and timeBefore filter the records based on time when the
// trades completed.  The value of time is Unix timestamp in seconds.
//
// the sortIDBy define the order of result set, default is sorted by ID in
// "desc" (descending) order.
// Valid values are "asc" for ascending and "desc" for descending order.
//
// This method require authentication.
//
func (cl *Client) UserTrades(
	pairName string,
	offset, limit, idAfter, idBefore, timeAfter, timeBefore int64,
) (
	trades []tokenomy.Trade, err error,
) {
	if !cl.env.IsValidPairName(pairName) {
		return nil, tokenomy.ErrInvalidPair
	}

	params := url.Values{
		tokenomy.ParamNamePair: []string{pairName},
	}
	if offset > 0 {
		params.Set(tokenomy.ParamNameOffset, strconv.FormatInt(offset, 10))
	}
	if limit > 0 && limit <= tokenomy.DefaultLimit {
		params.Set(tokenomy.ParamNameLimit, strconv.FormatInt(limit, 10))
	}
	if idAfter > 0 {
		params.Set(tokenomy.ParamNameIDAfter, strconv.FormatInt(idAfter, 10))
	}
	if idBefore > 0 {
		params.Set(tokenomy.ParamNameIDBefore, strconv.FormatInt(idBefore, 10))
	}
	if timeAfter > 0 {
		params.Set(tokenomy.ParamNameTimeAfter, strconv.FormatInt(timeAfter, 10))
	}
	if timeBefore > 0 {
		params.Set(tokenomy.ParamNameTimeBefore, strconv.FormatInt(timeBefore, 10))
	}

	b, err := cl.doSecureRequest(http.MethodGet, apiUserTrades, params)
	if err != nil {
		return nil, fmt.Errorf("UserTrades: %w", err)
	}

	res := &Response{
		Data: trades,
	}

	err = json.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}

	return trades, nil
}

//
// UserTradesClosed fetch the user closed trades based on pair's name.
// The offset parameter define the beginning of record and limit parameter
// define the maximum record in result set.
//
// This method require authentication.
//
func (cl *Client) UserTradesClosed(pairName string, offset, limit int64) (
	trades []tokenomy.Trade, err error,
) {
	if !cl.env.IsValidPairName(pairName) {
		return nil, tokenomy.ErrInvalidPair
	}

	params := url.Values{
		tokenomy.ParamNamePair: []string{pairName},
	}

	if offset > 0 {
		params.Set(tokenomy.ParamNameOffset, strconv.FormatInt(offset, 10))
	}
	if limit > 0 && limit <= tokenomy.DefaultLimit {
		params.Set(tokenomy.ParamNameLimit, strconv.FormatInt(limit, 10))
	}

	b, err := cl.doSecureRequest(http.MethodGet, apiUserTradesClosed, params)
	if err != nil {
		return nil, fmt.Errorf("UserTradesClosed: %w", err)
	}

	res := &Response{
		Data: trades,
	}

	err = json.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}

	return trades, nil
}

//
// UserTradesOpen fetch the user open trades based on pair's name.
//
// This method require authentication.
//
func (cl *Client) UserTradesOpen(pairName string) (openTrades *TradeOpens, err error) {
	if !cl.env.IsValidPairName(pairName) {
		return nil, tokenomy.ErrInvalidPair
	}

	params := url.Values{
		tokenomy.ParamNamePair: []string{pairName},
	}

	b, err := cl.doSecureRequest(http.MethodGet, apiUserTradesOpen, params)
	if err != nil {
		return nil, fmt.Errorf("UserTradesOpen: %w", err)
	}

	openTrades = &TradeOpens{}
	res := &Response{
		Data: openTrades,
	}

	err = json.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}

	return openTrades, nil
}

//
// UserTrade fetch a single user's trade information based on pair's name and
// trade ID.
//
// This method require authentication.
//
func (cl *Client) UserTrade(pairName string, id int64) (
	trade *tokenomy.Trade, err error,
) {
	if !cl.env.IsValidPairName(pairName) {
		return nil, tokenomy.ErrInvalidPair
	}

	params := url.Values{
		tokenomy.ParamNamePair:    []string{pairName},
		tokenomy.ParamNameTradeID: []string{strconv.FormatInt(id, 10)},
	}

	b, err := cl.doSecureRequest(http.MethodGet, apiUserTrade, params)
	if err != nil {
		return nil, fmt.Errorf("UserTrade: %w", err)
	}

	trade = &tokenomy.Trade{}
	res := &Response{
		Data: trade,
	}

	err = json.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}

	return trade, nil
}

//
// UserTransactions fetch all user deposit and withdraw transaction history.
// If the asset name is not empty, it will fetch only the deposit and withdraw
// based on the asset name.
//
// The limit parameter define the maximum record in result set.
//
// This method require authentication.
//
func (cl *Client) UserTransactions(asset string, limit int64) (trans *AssetTransactions, err error) {
	params := url.Values{}

	if len(asset) > 0 {
		params.Set(tokenomy.ParamNameAsset, asset)
	}
	if limit > 0 && limit <= tokenomy.DefaultLimit {
		params.Set(tokenomy.ParamNameLimit, strconv.FormatInt(limit, 10))
	}

	b, err := cl.doSecureRequest(http.MethodGet, apiUserTransactions, params)
	if err != nil {
		return nil, fmt.Errorf("UserTransactions: %w", err)
	}

	trans = &AssetTransactions{}
	res := &Response{
		Data: trans,
	}

	err = json.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}

	return trans, nil
}

//
// TradeAsk request to sell the coin on market with specific method, amount,
// and price.
// The method parameter define the mode of sell, its either "market" or
// "limit", default to "market" if its empty.
// If the method is "market", it will only accept amount parameter, otherwise
// if the methid is "limit", the amount and price must not be zero.
//
// The pairName parameter define the coin and base assets to be traded, in the
// following format: "coin_base".
//
// The amount parameter define the volume of coin we want to sell.
//
// The price parameter define the number of base that we want to sell the
// amount of coin.
//
func (cl *Client) TradeAsk(method, pairName string, amount, price *big.Rat) (
	trade *tokenomy.TradeResponse, err error,
) {
	return cl.trade(apiTradeAsk, method, pairName, amount, price)
}

//
// TradeBid request to buy the coin on market with specific method, amount,
// and price.
// The method parameter define the mode of buy, its either "market" or
// "limit", default to "market" if its empty.
// If the method is "market", it will only accept amount parameter, otherwise
// if the methid is "limit", the amount and price must not be zero.
//
// The pairName parameter define the coin and base assets to be traded, in the
// following format: "coin_base".
//
// The amount parameter define the volume of coin we want to buy.
//
// The price parameter define the number of base that we want to buy the
// amount of coin.
//
func (cl *Client) TradeBid(method, pairName string, amount, price *big.Rat) (
	trade *tokenomy.TradeResponse, err error,
) {
	return cl.trade(apiTradeBid, method, pairName, amount, price)
}

func (cl *Client) trade(
	api, method, pairName string,
	amount, price *big.Rat,
) (
	trade *tokenomy.TradeResponse, err error,
) {
	params := url.Values{}

	if len(method) == 0 {
		method = tokenomy.TradeMethodMarket
	} else {
		method = strings.ToLower(method)
		switch method {
		case tokenomy.TradeMethodMarket, tokenomy.TradeMethodLimit:
		default:
			return nil, tokenomy.ErrInvalidTradeMethod
		}
	}
	params.Set(tokenomy.ParamNameTradeMethod, method)

	if !cl.env.IsValidPairName(pairName) {
		return nil, tokenomy.ErrInvalidPair
	}
	params.Set(tokenomy.ParamNamePair, pairName)

	if amount.IsLessOrEqual(0) {
		return nil, tokenomy.ErrInvalidAmount
	}
	params.Set(tokenomy.ParamNameAmount, amount.String())

	if method == tokenomy.TradeMethodLimit {
		if price.IsLessOrEqual(0) {
			return nil, tokenomy.ErrInvalidPrice
		}
		params.Set(tokenomy.ParamNamePrice, price.String())
	}

	b, err := cl.doSecureRequest(http.MethodPost, api, params)
	if err != nil {
		return nil, err
	}

	trade = &tokenomy.TradeResponse{}
	res := &Response{
		Data: trade,
	}

	err = json.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}

	trade.Trade.Pair = pairName

	return trade, nil
}

//
// TradeCancelAsk cancel the specific open sell by pair and ID.
//
func (cl *Client) TradeCancelAsk(pairName string, id int64) (
	trade *tokenomy.TradeResponse, err error,
) {
	return cl.cancel(apiTradeCancelAsk, pairName, id)
}

//
// TradeCancelBid cancel the specific open buy by pair and ID.
//
func (cl *Client) TradeCancelBid(pairName string, id int64) (
	trade *tokenomy.TradeResponse, err error,
) {
	return cl.cancel(apiTradeCancelBid, pairName, id)
}

func (cl *Client) cancel(api, pairName string, id int64) (
	trade *tokenomy.TradeResponse, err error,
) {
	params := url.Values{}

	if !cl.env.IsValidPairName(pairName) {
		return nil, tokenomy.ErrInvalidPair
	}
	params.Set(tokenomy.ParamNamePair, pairName)

	if id <= 0 {
		return nil, tokenomy.ErrInvalidTradeID
	}
	params.Set(tokenomy.ParamNameTradeID, strconv.FormatInt(id, 10))

	b, err := cl.doSecureRequest(http.MethodDelete, api, params)
	if err != nil {
		return nil, err
	}

	trade = &tokenomy.TradeResponse{}
	res := &Response{
		Data: trade,
	}

	err = json.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}

	return trade, nil
}

//
// doGet create and send request without authentication on specific
// "path" with additional "params".
//
func (cl *Client) doGet(path string, params url.Values) (b []byte, err error) {
	req, err := http.NewRequest(http.MethodGet, cl.env.Address, nil)
	if err != nil {
		return nil, err
	}

	req.URL.Path = path
	req.URL.RawQuery = params.Encode()

	return cl.send(req)
}

func (cl *Client) doSecureRequest(httpMethod, path string, params url.Values) (
	b []byte, err error,
) {
	params.Set(tokenomy.ParamNameTimestamp, timestampAsString())

	payload := params.Encode()
	sign := tokenomy.Sign(payload, cl.env.Secret)

	req, err := http.NewRequest(httpMethod, cl.env.Address, nil)
	if err != nil {
		return nil, err
	}

	req.URL.Path = path
	req.Header.Set(tokenomy.HeaderNameKey, cl.env.Token)
	req.Header.Set(tokenomy.HeaderNameSign, sign)

	switch httpMethod {
	case http.MethodGet, http.MethodDelete:
		req.URL.RawQuery = payload
	case http.MethodPost, http.MethodPut:
		req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(payload)))
		req.Header.Set(tokenomy.HeaderNameContentType, tokenomy.ContentTypeForm)
	}

	return cl.send(req)
}

func (cl *Client) send(req *http.Request) (b []byte, err error) {
	if cl.env.Debug > 0 {
		fmt.Printf(">>> send: %+v\n", req)
	}

	httpres, err := cl.conn.Do(req)
	if err != nil {
		return nil, err
	}

	defer httpres.Body.Close()

	b, err = ioutil.ReadAll(httpres.Body)
	if err != nil {
		return nil, err
	}

	if cl.env.Debug > 0 {
		fmt.Printf("<<< send: %s\n", b)
	}

	if httpres.StatusCode >= 400 {
		res := &Response{}

		err = json.Unmarshal(b, res)
		if err != nil {
			return nil, err
		}

		err = &liberrors.E{
			Code:    httpres.StatusCode,
			Message: res.Message,
			Name:    res.Name,
		}

		return nil, err
	}

	return b, nil
}
