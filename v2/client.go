// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v2

import (
	"encoding/json"
	"fmt"
	stdhttp "net/http"
	"net/url"
	"strconv"

	"github.com/shuLhan/share/lib/http"
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

	cl = &Client{
		conn: http.NewClient(env.Address, nil, env.IsInsecure),
		env:  env,
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

	if len(pairName) == 0 {
		return nil, tokenomy.ErrInvalidPair
	}

	_, resBody, err := cl.conn.Get(nil, apiMarketDepths, params)
	if err != nil {
		return nil, fmt.Errorf("MarketDepths: %w", err)
	}

	depths = &MarketDepths{}
	res := &Response{
		Data: depths,
	}

	err = json.Unmarshal(resBody, res)
	if err != nil {
		return nil, err
	}

	return depths, nil
}

//
// MarketInfo return information about all the pair in the platform.
//
func (cl *Client) MarketInfo() (marketInfos []MarketInfo, err error) {
	_, resBody, err := cl.conn.Get(nil, apiMarketInfo, nil)
	if err != nil {
		return nil, fmt.Errorf("MarketInfo: %w", err)
	}

	marketInfos = make([]MarketInfo, 0)
	res := &Response{
		Data: marketInfos,
	}

	err = json.Unmarshal(resBody, res)
	if err != nil {
		return nil, err
	}

	return marketInfos, nil
}

//
// MarketTradesOpen return list of all open trades in the market, specific to
// pair's name, grouped by ask and bid.
//
func (cl *Client) MarketTradesOpen(pairName string) (openTrades *TradesOpen, err error) {
	params := url.Values{
		tokenomy.ParamNamePair: []string{pairName},
	}

	_, resBody, err := cl.conn.Get(nil, apiMarketTradesOpen, params)
	if err != nil {
		return nil, fmt.Errorf("MarketTradesOpen: %w", err)
	}

	openTrades = &TradesOpen{}
	res := &Response{
		Data: openTrades,
	}

	err = json.Unmarshal(resBody, res)
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

	_, resBody, err := cl.conn.Get(nil, apiMarketPrices, params)
	if err != nil {
		return nil, fmt.Errorf("MarketPrices: %w", err)
	}

	marketPrices = make(MarketPrices)
	res := &Response{
		Data: marketPrices,
	}

	err = json.Unmarshal(resBody, res)
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

	_, resBody, err := cl.conn.Get(nil, apiMarketTicker, params)
	if err != nil {
		return nil, fmt.Errorf("MarketTicker: %w", err)
	}

	tick = &Tick{}
	res := &Response{
		Data: tick,
	}

	err = json.Unmarshal(resBody, res)
	if err != nil {
		return nil, err
	}

	return tick, nil
}

//
// MarketTrades return list of all closed trades in the market, specific to
// pair's name, grouped by ask and bid.
//
func (cl *Client) MarketTrades(pairName string) (
	tradePrices *tokenomy.MarketTradePrices, err error,
) {
	params := url.Values{
		tokenomy.ParamNamePair: []string{pairName},
	}

	_, resBody, err := cl.conn.Get(nil, apiMarketTrades, params)
	if err != nil {
		return nil, fmt.Errorf("MarketTrades: %w", err)
	}

	tradePrices = &tokenomy.MarketTradePrices{}
	res := &Response{
		Data: tradePrices,
	}

	err = json.Unmarshal(resBody, res)
	if err != nil {
		return nil, err
	}

	return tradePrices, nil
}

//
// MarketSummaries return the summaries (ticker) of all pairs.
//
func (cl *Client) MarketSummaries() (summaries *MarketSummaries, err error) {
	_, resBody, err := cl.conn.Get(nil, apiMarketSummaries, nil)
	if err != nil {
		return nil, fmt.Errorf("MarketSummaries: %w", err)
	}

	summaries = &MarketSummaries{}
	res := &Response{
		Data: summaries,
	}

	err = json.Unmarshal(resBody, res)
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

	b, err := cl.doSecureRequest(stdhttp.MethodGet, apiUserInfo, params)
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

	b, err := cl.doSecureRequest(stdhttp.MethodGet, apiUserTrades, params)
	if err != nil {
		return nil, fmt.Errorf("UserTrades: %w", err)
	}

	res := &Response{
		Data: &trades,
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
	params := url.Values{
		tokenomy.ParamNamePair: []string{pairName},
	}

	if offset > 0 {
		params.Set(tokenomy.ParamNameOffset, strconv.FormatInt(offset, 10))
	}
	if limit <= 0 || limit > tokenomy.DefaultLimit {
		limit = tokenomy.DefaultLimit
	}
	params.Set(tokenomy.ParamNameLimit, strconv.FormatInt(limit, 10))

	b, err := cl.doSecureRequest(stdhttp.MethodGet, apiUserTradesClosed, params)
	if err != nil {
		return nil, fmt.Errorf("UserTradesClosed: %w", err)
	}

	res := &Response{
		Data: &trades,
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
func (cl *Client) UserTradesOpen(pairName string) (
	pairTradesOpen PairTradesOpen, err error,
) {
	params := url.Values{
		tokenomy.ParamNamePair: []string{pairName},
	}

	b, err := cl.doSecureRequest(stdhttp.MethodGet, apiUserTradesOpen, params)
	if err != nil {
		return nil, fmt.Errorf("UserTradesOpen: %w", err)
	}

	pairTradesOpen = make(PairTradesOpen)
	res := &Response{
		Data: &pairTradesOpen,
	}

	err = json.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}

	return pairTradesOpen, nil
}

//
// UserTradeInfo fetch a single user's trade information based on pair's name
// and trade ID.
//
// This method require authentication.
//
func (cl *Client) UserTradeInfo(pairName string, id int64) (
	trade *tokenomy.Trade, err error,
) {
	params := url.Values{
		tokenomy.ParamNamePair:    []string{pairName},
		tokenomy.ParamNameTradeID: []string{strconv.FormatInt(id, 10)},
	}

	b, err := cl.doSecureRequest(stdhttp.MethodGet, apiUserTrade, params)
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

	b, err := cl.doSecureRequest(stdhttp.MethodGet, apiUserTransactions, params)
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
// The method parameter define the mode of sell, its either "market" (default)
// or "limit".
// If the method is "market", it will only accept amount parameter, otherwise
// if the method is "limit", the amount and price must not be zero.
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
	params, _, err := generateTradeParams(method, pairName, amount, price)
	if err != nil {
		return nil, err
	}

	b, err := cl.doSecureRequest(stdhttp.MethodPost, api, params)
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
// TradeCancel cancel the open trade using ID and pair information in Trade.
//
func (cl *Client) TradeCancel(trade *tokenomy.Trade) (*tokenomy.Trade, error) {
	var (
		tradeResponse *tokenomy.TradeResponse
		err           error
	)

	switch trade.Type {
	case tokenomy.TradeTypeAsk:
		tradeResponse, err = cl.TradeCancelAsk(trade.Pair, trade.ID)
	case tokenomy.TradeTypeBid:
		tradeResponse, err = cl.TradeCancelBid(trade.Pair, trade.ID)
	default:
		return nil, tokenomy.ErrInvalidTradeType
	}
	if err != nil {
		return nil, err
	}

	return tradeResponse.Trade, nil
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

	params.Set(tokenomy.ParamNamePair, pairName)

	if id <= 0 {
		return nil, tokenomy.ErrInvalidTradeID
	}
	params.Set(tokenomy.ParamNameTradeID, strconv.FormatInt(id, 10))

	b, err := cl.doSecureRequest(stdhttp.MethodDelete, api, params)
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

func (cl *Client) doSecureRequest(httpMethod, path string, params url.Values) (
	resBody []byte, err error,
) {
	params.Set(tokenomy.ParamNameTimestamp, timestampAsString())

	payload := params.Encode()
	sign := tokenomy.Sign(payload, cl.env.Secret)

	headers := stdhttp.Header{
		tokenomy.HeaderNameKey:  []string{cl.env.Token},
		tokenomy.HeaderNameSign: []string{sign},
	}

	var httpres *stdhttp.Response

	switch httpMethod {
	case stdhttp.MethodGet:
		httpres, resBody, err = cl.conn.Get(headers, path, params)
	case stdhttp.MethodDelete:
		httpres, resBody, err = cl.conn.Delete(headers, path, params)
	case stdhttp.MethodPost:
		httpres, resBody, err = cl.conn.PostForm(headers, path, params)
	}

	if httpres.StatusCode >= 400 {
		res := &Response{}

		err = json.Unmarshal(resBody, res)
		if err != nil {
			return resBody, err
		}

		res.Code = httpres.StatusCode

		return nil, &res.E
	}

	return resBody, nil
}
