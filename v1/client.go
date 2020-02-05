// Copyright 2019 Tokenomy Technologies Pte. Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v1

import (
	"crypto/hmac"
	"crypto/sha512"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/tokenomy/tokenomy-go"
)

//
// Client represent an HTTP client for Tokenomy API v1.
//
type Client struct {
	Info *UserInfo
	conn *http.Client
	env  *tokenomy.Environment
}

//
// NewClient create and initialize new Tokenomy client for API v1.
//
// The Environment's Address parameter define the REST API v2 address, if its
// empty it will set to value in DefaultAddress.
//
// The Environment's Token and Secret parameters are used to
// authenticate the client when accessing private API.
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

	var transport http.Transport = http.Transport{}

	defTransport := http.DefaultTransport.(*http.Transport)
	transport = *defTransport

	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: env.IsInsecure,
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
	// Test the token and secret by requesting user information.
	_, err = cl.UserInfo()
	if err != nil {
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
	resBody, err := cl.trade(tokenomy.TradeMethodLimit,
		tokenomy.TradeTypeBid, pairName, amount, price)
	if err != nil {
		return nil, err
	}

	if cl.env.Debug >= 2 {
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
	resBody, err := cl.trade(tokenomy.TradeMethodMarket,
		tokenomy.TradeTypeBid, pairName, amount, 0)
	if err != nil {
		return nil, err
	}

	if cl.env.Debug >= 2 {
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

	resBody, err := cl.cancelOrder(tokenomy.TradeTypeAsk, pairName, orderID)
	if err != nil {
		return nil, err
	}

	if cl.env.Debug >= 2 {
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

	resBody, err := cl.cancelOrder(tokenomy.TradeTypeBid, pairName, orderID)
	if err != nil {
		return nil, err
	}

	if cl.env.Debug >= 2 {
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

	resBody, err := cl.callPrivate(MethodUserGetOrder, params)
	if err != nil {
		return nil, err
	}

	if cl.env.Debug >= 2 {
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

	if cl.env.Debug >= 2 {
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

	if cl.env.Debug >= 2 {
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

	resBody, err := cl.callPrivate(MethodUserOpenOrders, params)
	if err != nil {
		return nil, err
	}

	if cl.env.Debug >= 2 {
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
	resBody, err := cl.callPrivate(MethodUserOrderHistory, nil)
	if err != nil {
		return nil, err
	}

	if cl.env.Debug >= 2 {
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

	resBody, err := cl.callPrivate(MethodUserTradeHistory, params)
	if err != nil {
		return nil, err
	}

	if cl.env.Debug >= 2 {
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
	resBody, err := cl.callPrivate(MethodUserTransHistory, nil)
	if err != nil {
		return nil, err
	}

	if cl.env.Debug >= 2 {
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

	if cl.env.Debug >= 2 {
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

	if cl.env.Debug >= 2 {
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
	resBody, err := cl.trade(tokenomy.TradeMethodLimit,
		tokenomy.TradeTypeAsk, pairName, amount, price)
	if err != nil {
		return nil, err
	}

	if cl.env.Debug >= 2 {
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
	resBody, err := cl.trade(tokenomy.TradeMethodMarket,
		tokenomy.TradeTypeAsk, pair, amount, 0)
	if err != nil {
		return nil, err
	}

	if cl.env.Debug >= 2 {
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

	if cl.env.Debug >= 2 {
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
	resBody, err := cl.callPrivate(MethodUserGetInfo, nil)
	if err != nil {
		err = fmt.Errorf("UserInfo: " + err.Error())
		return nil, err
	}

	if cl.env.Debug >= 2 {
		fmt.Printf(">>> UserInfo: response body: %s\n", resBody)
	}

	userInfoResponse := &userInfoResponse{}

	err = json.Unmarshal(resBody, userInfoResponse)
	if err != nil {
		err = fmt.Errorf("UserInfo: " + err.Error())
		return nil, err
	}

	cl.Info = userInfoResponse.Return

	if cl.env.Debug >= 2 {
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

	resBody, err := cl.callPrivate(MethodWithdraw, params)
	if err != nil {
		return nil, err
	}

	if cl.env.Debug >= 2 {
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

//
// callPrivate call private API with specific method and parameters.
// On success it will return response body.
// On fail it will return an empty body with an error.
//
func (cl *Client) callPrivate(method string, params url.Values) (
	body []byte, err error,
) {
	if len(cl.env.Token) == 0 {
		return nil, ErrUnauthenticated
	}

	req, err := cl.newPrivateRequest(method, params)
	if err != nil {
		err = fmt.Errorf("callPrivate: " + err.Error())
		return nil, err
	}

	if cl.env.Debug >= 2 {
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
func (cl *Client) callPublic(publicPath string) (body []byte, err error) {
	req := &http.Request{
		Method: http.MethodPost,
		Header: http.Header{
			"Content-Type": []string{
				"application/x-www-form-urlencoded",
			},
		},
	}

	req.URL, err = url.Parse(cl.env.Address + publicPath)
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
func (cl *Client) cancelOrder(tipe, pairName string, orderID int64) (
	body []byte, err error,
) {
	orderIDStr := strconv.FormatInt(orderID, 10)

	params := map[string][]string{
		"type":     {tipe},
		"pair":     {pairName},
		"order_id": {orderIDStr},
	}

	body, err = cl.callPrivate(MethodTradeCancelOrder, params)
	if err != nil {
		return nil, err
	}

	return body, nil
}

//
// encodeToString return the hex encoding of data hashed with client's API
// secret key.
//
func (cl *Client) encodeToString(in []byte) string {
	mac := hmac.New(sha512.New, []byte(cl.env.Secret))

	_, _ = mac.Write(in)

	return hex.EncodeToString(mac.Sum(nil))
}

//
// newPrivateRequest create a new HTTP request for private API with specific
// "method" and custom parameters "params" as form values in POST body.
//
func (cl *Client) newPrivateRequest(apiMethod string, params url.Values) (
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

	if cl.env.Debug >= 2 {
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
				cl.env.Token,
			},
			"Sign": []string{
				sign,
			},
		},
		Body: ioutil.NopCloser(strings.NewReader(reqBody)),
	}

	req.URL, err = url.Parse(cl.env.Address + defPrivatePath)
	if err != nil {
		err = fmt.Errorf("newPrivateRequest: " + err.Error())
		return nil, err
	}

	return req, nil
}

func (cl *Client) trade(method, tipe, pair string, amount, price float64) (
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

	body, err = cl.callPrivate(MethodTrade, params)
	if err != nil {
		return nil, err
	}

	return body, nil
}
