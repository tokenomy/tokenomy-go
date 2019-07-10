// Copyright 2019 Tokenomy Technologies Pte. Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

//
// Client represent an HTTP client for Tokenomy API v1.
//
type Client struct {
	Info *UserInfo
	client
}

//
// NewClient create and initialize new Tokenomy client for API v1.
//
// The token and secret parameters are used to authenticate the client when
// accessing private API.
//
// By default, the key and secret is read from environment variables
// "TOKENOMY_KEY" and "TOKENOMY_SECRET", the parameters will override the
// default value, if its set.
// If both environment variables and the parameters are empty, the client can
// only access the public API.
//
func NewClient(token, secret string) (cl *Client, err error) {
	cl = &Client{
		client: client{
			conn: &http.Client{},
			env:  newEnvironment(),
		},
	}

	cl.baseHost = cl.env.v1BaseHost
	cl.privatePath = cl.env.v1PrivatePath

	if len(token) == 0 {
		cl.token = cl.env.apiKey
		cl.secret = cl.env.apiSecret
	} else {
		cl.token = token
		cl.secret = secret
	}

	if len(cl.token) > 0 {
		err = cl.Authenticate(cl.token, cl.secret)
	}

	return cl, err
}

//
// Authenticate the current client's connection using token and secret keys.
//
func (cl *Client) Authenticate(token, secret string) (err error) {
	cl.token = token
	cl.secret = secret

	// Test the token and secret by requesting user information.
	_, err = cl.UserInfo()
	if err != nil {
		cl.token = cl.env.apiKey
		cl.secret = cl.env.apiSecret
		err = fmt.Errorf("Authenticate: " + err.Error())
		return err
	}

	return nil
}

//
// Buy put the buy order for the asset at specific amount and price into the
// market.
//
// On success, it will return the information about status of transaction and
// remaining balance.
// On fail, it will return nil TradeResponse and an error.
//
// This method require authentication.
//
func (cl *Client) Buy(pairName string, amount, price float64) (
	tres *TradeResponse, err error,
) {
	resBody, err := cl.trade(tradeMethodLimit, tradeTypeBuy, pairName,
		amount, price)
	if err != nil {
		return nil, err
	}

	if cl.env.debug >= 2 {
		fmt.Printf(">>> Buy: response body: %s\n", resBody)
	}

	tres = &TradeResponse{}

	err = json.Unmarshal(resBody, tres)
	if err != nil {
		return nil, err
	}

	if tres.Success != responseSuccess {
		return nil, fmt.Errorf("Buy: %s", tres.Error)
	}

	tres.Pair = pairName

	return tres, nil
}

//
// BuyByMarket buy the asset at specific amount from market using the base
// asset as currency.
//
// On success, it will return the information about status of transaction and
// remaining balances.
// On fail, it will return nil TradeResponse and an error.
//
// This method require authentication.
//
func (cl *Client) BuyByMarket(pairName string, amount float64) (
	tres *TradeResponse, err error,
) {
	resBody, err := cl.trade(tradeMethodMarket, tradeTypeBuy, pairName, amount, 0)
	if err != nil {
		return nil, err
	}

	if cl.env.debug >= 2 {
		fmt.Printf(">>> BuyByMarket: response body: %s\n", resBody)
	}

	tres = &TradeResponse{}

	err = json.Unmarshal(resBody, tres)
	if err != nil {
		return nil, err
	}

	if tres.Success == 0 {
		return nil, fmt.Errorf("%s: %s", tres.ErrorCode, tres.Error)
	}

	tres.Pair = pairName

	return tres, nil
}

//
// CancelSell cancel the sell (ask) order on specific pair name and order ID.
//
// This method require authentication.
//
func (cl *Client) CancelSell(pairName string, orderID int64) (
	cancelOrder *CancelOrder, err error,
) {
	if len(pairName) == 0 {
		return nil, ErrInvalidPairName
	}

	resBody, err := cl.cancelOrder(tradeTypeSell, pairName, orderID)
	if err != nil {
		return nil, err
	}

	if cl.env.debug >= 2 {
		fmt.Printf(">>> CancelSell: response body: %s\n", resBody)
	}

	cancelRes := &cancelOrderResponse{}

	err = json.Unmarshal(resBody, cancelRes)
	if err != nil {
		return nil, fmt.Errorf("CancelSell: " + err.Error())
	}

	if cancelRes.Success != responseSuccess {
		return nil, fmt.Errorf("CancelSell: " + cancelRes.Error)
	}

	cancelOrder = cancelRes.Return
	cancelRes.Return = nil

	return cancelOrder, nil
}

//
// CancelBuy cancel the buy (bid) order on specific pair name and order ID.
//
// This method require authentication.
//
func (cl *Client) CancelBuy(pairName string, orderID int64) (
	cancelOrder *CancelOrder, err error,
) {
	if len(pairName) == 0 {
		return nil, ErrInvalidPairName
	}

	resBody, err := cl.cancelOrder(tradeTypeBuy, pairName, orderID)
	if err != nil {
		return nil, err
	}

	if cl.env.debug >= 2 {
		fmt.Printf(">>> CancelBuy: response body: %s\n", resBody)
	}

	cancelRes := &cancelOrderResponse{}

	err = json.Unmarshal(resBody, cancelRes)
	if err != nil {
		return nil, fmt.Errorf("CancelBuy: " + err.Error())
	}

	if cancelRes.Success != responseSuccess {
		return nil, fmt.Errorf("CancelBuy: " + cancelRes.Error)
	}

	cancelOrder = cancelRes.Return
	cancelRes.Return = nil

	return cancelOrder, nil
}

//
// GetOrder get the detail of a specific open order by pair name and order ID.
//
// This method require authentication.
//
func (cl *Client) GetOrder(pairName string, orderID int64) (
	order *OrderHistory, err error,
) {
	if len(pairName) == 0 {
		return nil, ErrInvalidPairName
	}

	params := url.Values{}
	params.Set("pair", pairName)
	params.Set("order_id", strconv.FormatInt(orderID, 10))

	resBody, err := cl.callPrivate(apiViewGetOrder, params)
	if err != nil {
		return nil, err
	}

	if cl.env.debug >= 2 {
		fmt.Printf(">>> GetOrder: response body: %s\n", resBody)
	}

	gores := &getOrderResponse{}

	err = json.Unmarshal(resBody, gores)
	if err != nil {
		return nil, fmt.Errorf("GetOrder: " + err.Error())
	}

	order = gores.Return.Order
	order.PairName = pairName

	return order, nil
}

//
// GetTicker get the price summary of an individual pair.
//
func (cl *Client) GetTicker(pairName string) (pair *Pair, err error) {
	if len(pairName) == 0 {
		return nil, ErrInvalidPairName
	}

	apiPath := fmt.Sprintf(apiTicker, pairName)

	body, err := cl.callPublic(apiPath)
	if err != nil {
		return nil, fmt.Errorf("GetTicker: " + err.Error())
	}

	if cl.env.debug >= 2 {
		fmt.Printf(">>> GetTicker: response body: %s\n", body)
	}

	tickerRes := &tickerResponse{}

	err = json.Unmarshal(body, tickerRes)
	if err != nil {
		return nil, fmt.Errorf("GetTicker: " + err.Error())
	}

	pair = tickerRes.Ticker
	pair.Name = pairName

	return pair, nil
}

//
// ListTrades get the latest trades for a particular pair.
//
func (cl *Client) ListTrades(pairName string) (trades []*Trade, err error) {
	if len(pairName) == 0 {
		return nil, ErrInvalidPairName
	}

	apiPath := fmt.Sprintf(apiTrades, pairName)

	body, err := cl.callPublic(apiPath)
	if err != nil {
		return nil, fmt.Errorf("ListTrades: " + err.Error())
	}

	if cl.env.debug >= 2 {
		fmt.Printf(">>> ListTrades: response body: %s\n", body)
	}

	err = json.Unmarshal(body, &trades)
	if err != nil {
		return nil, fmt.Errorf("ListTrades: " + err.Error())
	}

	return trades, nil
}

//
// ListOpenOrders list the current user's open order (buy and sell) by pair
// name.
//
// This method require authentication.
//
func (cl *Client) ListOpenOrders(pairName string) (
	openOrders OpenOrders, err error,
) {
	if len(pairName) == 0 {
		return nil, ErrInvalidPairName
	}

	params := url.Values{}
	params.Set("pair", pairName)

	resBody, err := cl.callPrivate(apiViewOpenOrders, params)
	if err != nil {
		return nil, err
	}

	if cl.env.debug >= 2 {
		fmt.Printf(">>> ListOpenOrders: response body: %s\n", resBody)
	}

	oores := &openOrdersResponse{}

	err = json.Unmarshal(resBody, oores)
	if err != nil {
		return nil, fmt.Errorf("ListOpenOrders: " + err.Error())
	}

	if oores.Success != responseSuccess {
		return nil, fmt.Errorf("ListOpenOrders: " + oores.Error)
	}

	openOrders = oores.Return.Orders
	oores.Return = nil

	for _, list := range openOrders {
		for _, oo := range list {
			oo.PairName = pairName
		}
	}

	return openOrders, nil
}

//
// ListOrderHistory list user's closed order history (buy and sell).
//
// This method require authentication.
//
func (cl *Client) ListOrderHistory(pairName string, count, fromID int) (
	orderHistory []*OrderHistory, err error,
) {
	resBody, err := cl.callPrivate(apiViewOrderHistory, nil)
	if err != nil {
		return nil, err
	}

	if cl.env.debug >= 2 {
		fmt.Printf(">>> ListOrderHistory: response body: %s\n", resBody)
	}

	ohres := &orderHistoryResponse{}

	err = json.Unmarshal(resBody, ohres)
	if err != nil {
		return nil, fmt.Errorf("ListOrderHistory: " + err.Error())
	}

	if ohres.Success != responseSuccess {
		return nil, fmt.Errorf("ListOrderHistory: " + ohres.Error)
	}

	orderHistory = ohres.Return.Orders
	ohres.Return = nil

	for _, oh := range orderHistory {
		oh.PairName = pairName
	}

	return orderHistory, nil
}

//
// ListTradeHistory list all user's history of trade.
//
// The "count" parameter limit the number of history returned by this method,
// default to 1000.
// The "startTradeID" parameter filter the history to begin from specific
// trade ID, default to 0.
// The "endTradeID" parameter filter the history only until the specific trade
// ID.  sortOrder define the order of returned history.
// Valid values is "asc" or "desc" (default).
// The "sinceTime" parameter filter the history that start from specific time.
// The "endTime" parameter filter the history only until the specific time.
//
// This method require authentication.
//
func (cl *Client) ListTradeHistory(
	pairName string,
	count, startTradeID, endTradeID int64,
	sortOrder string,
	sinceTime *time.Time,
	endTime *time.Time,
) (trades []TradeHistory, err error) {
	if len(pairName) == 0 {
		return nil, ErrInvalidPairName
	}

	params := url.Values{}
	params.Set("pair", pairName)
	if count > 0 {
		params.Set("count", strconv.FormatInt(count, 10))
	}
	if startTradeID > 0 {
		params.Set("from_id", strconv.FormatInt(startTradeID, 10))
	}
	if endTradeID > 0 {
		params.Set("end_id", strconv.FormatInt(endTradeID, 10))
	}

	sortOrder = strings.ToLower(sortOrder)
	switch sortOrder {
	case "asc":
		params.Set("order", "asc")
	case "desc":
		params.Set("order", "desc")
	}

	if sinceTime != nil {
		params.Set("since", strconv.FormatInt(sinceTime.Unix(), 10))
	}
	if endTime != nil {
		params.Set("end", strconv.FormatInt(endTime.Unix(), 10))
	}

	resBody, err := cl.callPrivate(apiViewTradeHistory, params)
	if err != nil {
		return nil, err
	}

	if cl.env.debug >= 2 {
		fmt.Printf(">>> ListTradeHistory: response body: %s\n", resBody)
	}

	thres := &tradeHistoryResponse{}

	err = json.Unmarshal(resBody, thres)
	if err != nil {
		return nil, err
	}

	if thres.Success != responseSuccess {
		return nil, fmt.Errorf("ListTradeHistory: %s: %s",
			thres.ErrorCode, thres.Error)
	}

	trades = thres.Return.Trades

	return trades, nil
}

//
// ListTransactionHistory list all user's history of deposits and withdrawals
// from all assets.
//
// This method is require authentication.
//
func (cl *Client) ListTransactionHistory() (transHistory *TransactionHistory, err error) {
	resBody, err := cl.callPrivate(apiViewTransactionHistory, nil)
	if err != nil {
		return nil, err
	}

	if cl.env.debug >= 2 {
		fmt.Printf(">>> ListTransactionHistory: response body: %s\n", resBody)
	}

	thres := &transHistoryResponse{}

	err = json.Unmarshal(resBody, thres)
	if err != nil {
		return nil, fmt.Errorf("ListTransactionHistory: " + err.Error())
	}

	if thres.Success != responseSuccess {
		return nil, fmt.Errorf("ListTransactionHistory: %s: %s",
			thres.ErrorCode, thres.Error)
	}

	transHistory = thres.Return
	thres.Return = nil

	return transHistory, nil
}

//
// MarketInfo list of all available pairs including limit information and
// market status.
//
func (cl *Client) MarketInfo() (marketInfos []MarketInfo, err error) {
	body, err := cl.callPublic(apiMarketInfo)
	if err != nil {
		return nil, fmt.Errorf("MarketInfo: " + err.Error())
	}

	if cl.env.debug >= 2 {
		fmt.Printf(">>> MarketInfo: response body: %s\n", body)
	}

	marketInfos = make([]MarketInfo, 0)

	err = json.Unmarshal(body, &marketInfos)
	if err != nil {
		return nil, fmt.Errorf("MarketInfo: " + err.Error())
	}

	return marketInfos, nil
}

//
// OrderBook list the public open order book (buy and sell) for spesific
// pair.
//
func (cl *Client) OrderBook(pairName string) (orderBook *OrderBook, err error) {
	if len(pairName) == 0 {
		return nil, ErrInvalidPairName
	}

	apiPath := fmt.Sprintf(apiOrderBook, pairName)

	body, err := cl.callPublic(apiPath)
	if err != nil {
		return nil, fmt.Errorf("OrderBook: " + err.Error())
	}

	if cl.env.debug >= 2 {
		fmt.Printf(">>> OrderBook: response body: %s\n", body)
	}

	orderBook = &OrderBook{}

	err = json.Unmarshal(body, orderBook)
	if err != nil {
		return nil, fmt.Errorf("OrderBook: " + err.Error())
	}

	return orderBook, nil
}

//
// Sell put the sell order into open orders at specific amount and price.
//
// This method require authentication.
//
func (cl *Client) Sell(pairName string, amount, price float64) (
	tres *TradeResponse, err error,
) {
	resBody, err := cl.trade(tradeMethodLimit, tradeTypeSell, pairName,
		amount, price)
	if err != nil {
		return nil, err
	}

	if cl.env.debug >= 2 {
		fmt.Printf(">>> Sell: response body: %s\n", resBody)
	}

	tres = &TradeResponse{}

	err = json.Unmarshal(resBody, tres)
	if err != nil {
		return nil, err
	}

	if tres.Success != responseSuccess {
		return nil, fmt.Errorf("Sell: %s", tres.Error)
	}

	tres.Pair = pairName

	return tres, nil
}

//
// SellByMarket sell the asset at specific amount to the market.
//
// On success, it will return the information about status of sell transaction
// and your remaining balance.
// On fail, it will return nil TradeResponse and an error.
//
// This method require authentication.
//
func (cl *Client) SellByMarket(pair string, amount float64) (tres *TradeResponse, err error) {
	resBody, err := cl.trade(tradeMethodMarket, tradeTypeSell, pair, amount, 0)
	if err != nil {
		return nil, err
	}

	if cl.env.debug >= 2 {
		fmt.Printf(">>> SellByMarket: response body: %s\n", resBody)
	}

	tres = &TradeResponse{}

	err = json.Unmarshal(resBody, tres)
	if err != nil {
		return nil, err
	}

	return tres, nil
}

//
// Summaries retrieve the summary of all traded pairs, highest price, lowest
// price, volume, last price, token/coin name.
// This API method can also be used to discover all current traded pairs.
//
func (cl *Client) Summaries() (summary *Summary, err error) {
	body, err := cl.callPublic(apiSummaries)
	if err != nil {
		return nil, err
	}

	if cl.env.debug >= 2 {
		fmt.Printf(">>> Summaries: response body: %s\n", body)
	}

	summary = &Summary{}

	err = json.Unmarshal(body, summary)
	if err != nil {
		err = fmt.Errorf("Summaries: " + err.Error())
		return nil, err
	}

	summary.propagate()

	return summary, nil
}

//
// UserInfo fetch the user's balance and information.
//
// This method require authentication.
//
func (cl *Client) UserInfo() (*UserInfo, error) {
	resBody, err := cl.callPrivate(apiViewGetInfo, nil)
	if err != nil {
		err = fmt.Errorf("UserInfo: " + err.Error())
		return nil, err
	}

	if cl.env.debug >= 2 {
		fmt.Printf(">>> UserInfo: response body: %s\n", resBody)
	}

	userInfoResponse := &userInfoResponse{}

	err = json.Unmarshal(resBody, userInfoResponse)
	if err != nil {
		err = fmt.Errorf("UserInfo: " + err.Error())
		return nil, err
	}

	cl.Info = userInfoResponse.Return

	if cl.env.debug >= 2 {
		fmt.Printf("=== UserInfo: %+v\n", cl.Info)
	}

	return cl.Info, nil
}

//
// WithdrawCoin withdraw your assets into another address.
// This method accept withdrawing all coins except TEN.
//
// This method require the "withdraw" permission, otherwise it will return a
// “No permission” error.
//
// You also need to prepare a Callback URL, when setting up the API keys.
// Callback URL is an URL that our system will call to verify your withdrawal
// request.
// Various parameters will be sent to Callback URL so please make
// sure that this information is in your application.
// If all the data is correct, the callback URL should return HTTP response
// 200 with string “ok” (without quotes), and we will process the withdrawn in
// our system, otherwise the request will be fail.
//
func (cl *Client) WithdrawCoin(
	currencyName, address, amount, requestID, memo string,
) (
	withdrawRes *WithdrawResponse, err error,
) {
	params := map[string][]string{
		"currency":         {currencyName},
		"withdraw_address": {address},
		"withdraw_amount":  {amount},
		"withdraw_memo":    {memo},
		"request_id":       {requestID},
	}

	resBody, err := cl.callPrivate(apiWithdraw, params)
	if err != nil {
		return nil, err
	}

	if cl.env.debug >= 2 {
		fmt.Printf(">>> WithdrawCoin: response body: %s\n", resBody)
	}

	withdrawRes = &WithdrawResponse{}

	err = json.Unmarshal(resBody, withdrawRes)
	if err != nil {
		return nil, fmt.Errorf("withdrawCoin: " + err.Error())
	}

	if withdrawRes.Success != responseSuccess {
		return nil, fmt.Errorf("withdrawCoin: " + withdrawRes.Error)
	}

	return withdrawRes, nil
}
