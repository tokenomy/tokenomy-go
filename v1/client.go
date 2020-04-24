// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
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
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/shuLhan/share/lib/math/big"

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

	transport := http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout: tokenomy.DefaultDialTimeout,
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
			Timeout:   tokenomy.DefaultTimeout,
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
// TradeBid put the buy order for the specific asset at specific amount and
// price into the market.
//
// On success, it will return the information about status of transaction and
// remaining balance.
// On fail, it will return nil TradeResponse and an error.
//
// This method require authentication.
//
func (cl *Client) TradeBid(method, pairName string, amount, price *big.Rat) (
	tres *tokenomy.TradeResponse, err error,
) {
	resBody, err := cl.trade(method, tokenomy.TradeTypeBid, pairName, amount, price)
	if err != nil {
		return nil, err
	}

	if cl.env.Debug >= 2 {
		fmt.Printf(">>> TradeBid: response body: %s\n", resBody)
	}

	intRes := &tradeResponse{}

	err = json.Unmarshal(resBody, intRes)
	if err != nil {
		return nil, err
	}

	if intRes.Success != responseSuccess {
		return nil, fmt.Errorf("TradeBid: %s", intRes.Error)
	}

	tres = &tokenomy.TradeResponse{
		Trade: &tokenomy.Trade{
			ID:    intRes.OrderID,
			Pair:  pairName,
			Price: price,

			CoinAmount: amount,
			CoinFilled: big.NewRat(intRes.Receive),

			BaseAmount: big.MulRat(amount, price),
			BaseRemain: big.NewRat(intRes.Remain),
			BaseFilled: big.NewRat(intRes.Filled),
		},
		User: tokenomy.User{},
	}

	tres.User.UserAssets = convertBalance(intRes.Balances)

	return tres, nil
}

//
// TradeCancelAsk cancel the sell (ask) order on specific pair name and order
// ID.
//
// This method require authentication.
//
func (cl *Client) TradeCancelAsk(pairName string, orderID int64) (
	tres *tokenomy.TradeResponse, err error,
) {
	if len(pairName) == 0 {
		return nil, ErrInvalidPairName
	}

	resBody, err := cl.cancelOrder(tokenomy.TradeTypeAsk, pairName, orderID)
	if err != nil {
		return nil, err
	}

	if cl.env.Debug >= 2 {
		fmt.Printf(">>> TradeCancelAsk: response body: %s\n", resBody)
	}

	cancelRes := &cancelOrderResponse{}

	err = json.Unmarshal(resBody, cancelRes)
	if err != nil {
		return nil, fmt.Errorf("TradeCancelAsk: " + err.Error())
	}

	if cancelRes.Success != responseSuccess {
		return nil, fmt.Errorf("TradeCancelAsk: " + cancelRes.Error)
	}

	tres = &tokenomy.TradeResponse{
		Trade: &tokenomy.Trade{
			ID:     cancelRes.Return.OrderID,
			Pair:   pairName,
			Type:   cancelRes.Return.Type,
			Status: "cancelled",
		},
		User: tokenomy.User{},
	}

	if cancelRes.Return != nil {
		tres.User.UserAssets = convertBalance(cancelRes.Return.Balances)
	}

	return tres, nil
}

//
// TradeCancelBid cancel the buy (bid) order on specific pair name and order
// ID.
//
// This method require authentication.
//
func (cl *Client) TradeCancelBid(pairName string, orderID int64) (
	tres *tokenomy.TradeResponse, err error,
) {
	if len(pairName) == 0 {
		return nil, ErrInvalidPairName
	}

	resBody, err := cl.cancelOrder(tokenomy.TradeTypeBid, pairName, orderID)
	if err != nil {
		return nil, err
	}

	if cl.env.Debug >= 2 {
		fmt.Printf(">>> TradeCancelBid: response body: %s\n", resBody)
	}

	cancelRes := &cancelOrderResponse{}

	err = json.Unmarshal(resBody, cancelRes)
	if err != nil {
		return nil, fmt.Errorf("TradeCancelBid: " + err.Error())
	}

	if cancelRes.Success != responseSuccess {
		return nil, fmt.Errorf("TradeCancelBid: " + cancelRes.Error)
	}

	tres = &tokenomy.TradeResponse{
		Trade: &tokenomy.Trade{
			ID:     cancelRes.Return.OrderID,
			Pair:   pairName,
			Type:   cancelRes.Return.Type,
			Status: "cancelled",
		},
		User: tokenomy.User{},
	}
	if cancelRes.Return != nil {
		tres.User.UserAssets = convertBalance(cancelRes.Return.Balances)
	}

	return tres, nil
}

//
// UserOrder get the detail of a specific user's open trade by pair name and
// trade ID.
//
// This method require authentication.
//
func (cl *Client) UserOrder(pairName string, orderID int64) (
	order *OrderHistory, err error,
) {
	if len(pairName) == 0 {
		return nil, ErrInvalidPairName
	}

	params := url.Values{}
	params.Set("pair", pairName)
	params.Set("order_id", strconv.FormatInt(orderID, 10))

	resBody, err := cl.callPrivate(MethodUserOrder, params)
	if err != nil {
		return nil, err
	}

	if cl.env.Debug >= 2 {
		fmt.Printf(">>> UserOrder: response body: %s\n", resBody)
	}

	gores := &getOrderResponse{}

	err = json.Unmarshal(resBody, gores)
	if err != nil {
		return nil, fmt.Errorf("UserOrder: " + err.Error())
	}

	order = gores.Return.Order
	order.PairName = pairName

	return order, nil
}

//
// MarketTicker get the price summary of an individual pair.
//
func (cl *Client) MarketTicker(pairName string) (pair *Pair, err error) {
	if len(pairName) == 0 {
		return nil, ErrInvalidPairName
	}

	apiPath := fmt.Sprintf(apiMarketTicker, pairName)

	body, err := cl.callPublic(apiPath)
	if err != nil {
		return nil, fmt.Errorf("MarketTicker: " + err.Error())
	}

	if cl.env.Debug >= 2 {
		fmt.Printf(">>> MarketTicker: response body: %s\n", body)
	}

	tickerRes := &tickerResponse{}

	err = json.Unmarshal(body, tickerRes)
	if err != nil {
		return nil, fmt.Errorf("MarketTicker: " + err.Error())
	}

	pair = tickerRes.Ticker
	pair.Name = pairName

	return pair, nil
}

//
// MarketTrades get the latest trades for a particular pair.
//
func (cl *Client) MarketTrades(pairName string) (trades []*Trade, err error) {
	if len(pairName) == 0 {
		return nil, ErrInvalidPairName
	}

	apiPath := fmt.Sprintf(apiMarketTrades, pairName)

	body, err := cl.callPublic(apiPath)
	if err != nil {
		return nil, fmt.Errorf("MarketTrades: " + err.Error())
	}

	if cl.env.Debug >= 2 {
		fmt.Printf(">>> MarketTrades: response body: %s\n", body)
	}

	err = json.Unmarshal(body, &trades)
	if err != nil {
		return nil, fmt.Errorf("MarketTrades: " + err.Error())
	}

	return trades, nil
}

//
// UserOrdersOpen list the current user's open order (buy and sell) by pair
// name.
//
// This method require authentication.
//
func (cl *Client) UserOrdersOpen(pairName string) (
	openOrders OpenOrders, err error,
) {
	if len(pairName) == 0 {
		return nil, ErrInvalidPairName
	}

	params := url.Values{}
	params.Set("pair", pairName)

	resBody, err := cl.callPrivate(MethodUserOrdersOpen, params)
	if err != nil {
		return nil, err
	}

	if cl.env.Debug >= 2 {
		fmt.Printf(">>> UserOrdersOpen: response body: %s\n", resBody)
	}

	oores := &openOrdersResponse{}

	err = json.Unmarshal(resBody, oores)
	if err != nil {
		return nil, fmt.Errorf("UserOrdersOpen: " + err.Error())
	}

	if oores.Success != responseSuccess {
		return nil, fmt.Errorf("UserOrdersOpen: " + oores.Error)
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
// UserOrdersClosed list user's closed order history (buy and sell).
//
// This method require authentication.
//
func (cl *Client) UserOrdersClosed(pairName string, count, fromID int) (
	orderHistory []*OrderHistory, err error,
) {
	resBody, err := cl.callPrivate(MethodUserOrderHistory, nil)
	if err != nil {
		return nil, err
	}

	if cl.env.Debug >= 2 {
		fmt.Printf(">>> UserOrdersClosed: response body: %s\n", resBody)
	}

	ohres := &orderHistoryResponse{}

	err = json.Unmarshal(resBody, ohres)
	if err != nil {
		return nil, fmt.Errorf("UserOrdersClosed: " + err.Error())
	}

	if ohres.Success != responseSuccess {
		return nil, fmt.Errorf("UserOrdersClosed: " + ohres.Error)
	}

	orderHistory = ohres.Return.Orders
	ohres.Return = nil

	for _, oh := range orderHistory {
		oh.PairName = pairName
	}

	return orderHistory, nil
}

//
// UserTrades list all user's history of trade.
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
func (cl *Client) UserTrades(
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
		fmt.Printf(">>> UserTrades: response body: %s\n", resBody)
	}

	thres := &tradeHistoryResponse{}

	err = json.Unmarshal(resBody, thres)
	if err != nil {
		return nil, err
	}

	if thres.Success != responseSuccess {
		return nil, fmt.Errorf("UserTrades: %s: %s",
			thres.ErrorCode, thres.Error)
	}

	trades = thres.Return.Trades

	return trades, nil
}

//
// UserTransactions list all user's history of deposits and withdrawals
// from all assets.
//
// This method is require authentication.
//
func (cl *Client) UserTransactions() (transHistory *TransactionHistory, err error) {
	resBody, err := cl.callPrivate(MethodUserTransHistory, nil)
	if err != nil {
		return nil, err
	}

	if cl.env.Debug >= 2 {
		fmt.Printf(">>> UserTransactions: response body: %s\n", resBody)
	}

	thres := &transHistoryResponse{}

	err = json.Unmarshal(resBody, thres)
	if err != nil {
		return nil, fmt.Errorf("UserTransactions: " + err.Error())
	}

	if thres.Success != responseSuccess {
		return nil, fmt.Errorf("UserTransactions: %s: %s",
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
// MarketOrdersOpen list the public open order book (buy and sell) for
// specific pair.
//
func (cl *Client) MarketOrdersOpen(pairName string) (orderBook *OrderBook, err error) {
	if len(pairName) == 0 {
		return nil, ErrInvalidPairName
	}

	apiPath := fmt.Sprintf(apiMarketOrdersOpen, pairName)

	body, err := cl.callPublic(apiPath)
	if err != nil {
		return nil, fmt.Errorf("MarketOrdersOpen: " + err.Error())
	}

	if cl.env.Debug >= 2 {
		fmt.Printf(">>> MarketOrdersOpen: response body: %s\n", body)
	}

	orderBook = &OrderBook{}

	err = json.Unmarshal(body, orderBook)
	if err != nil {
		return nil, fmt.Errorf("MarketOrdersOpen: " + err.Error())
	}

	return orderBook, nil
}

//
// TradeAsk put the sell order for specific asset at specific amount and
// price.
//
// This method require authentication.
//
func (cl *Client) TradeAsk(method, pairName string, amount, price *big.Rat) (
	tres *tokenomy.TradeResponse, err error,
) {
	resBody, err := cl.trade(method, tokenomy.TradeTypeAsk, pairName, amount, price)
	if err != nil {
		return nil, err
	}

	if cl.env.Debug >= 2 {
		fmt.Printf(">>> TradeAsk: response body: %s\n", resBody)
	}

	intRes := &tradeResponse{}

	err = json.Unmarshal(resBody, intRes)
	if err != nil {
		return nil, err
	}

	if intRes.Success != responseSuccess {
		return nil, fmt.Errorf("TradeAsk: %s", intRes.Error)
	}

	tres = &tokenomy.TradeResponse{
		Trade: &tokenomy.Trade{
			ID:    intRes.OrderID,
			Pair:  pairName,
			Price: price,

			CoinAmount: amount,
			CoinRemain: big.NewRat(intRes.Remain),
			CoinFilled: big.NewRat(intRes.Filled),

			BaseAmount: big.NewRat(intRes.Receive),
		},
		User: tokenomy.User{},
	}

	tres.User.UserAssets = convertBalance(intRes.Balances)

	return tres, nil
}

//
// MarketSummaries retrieve the summary of all traded pairs, highest price,
// lowest price, volume, last price, token/coin name.
// This API method can also be used to discover all current traded pairs.
//
func (cl *Client) MarketSummaries() (summary *Summary, err error) {
	body, err := cl.callPublic(apiMarketSummaries)
	if err != nil {
		return nil, err
	}

	if cl.env.Debug >= 2 {
		fmt.Printf(">>> MarketSummaries: response body: %s\n", body)
	}

	summary = &Summary{}

	err = json.Unmarshal(body, summary)
	if err != nil {
		err = fmt.Errorf("MarketSummaries: " + err.Error())
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
	resBody, err := cl.callPrivate(MethodUserInfo, nil)
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
// UserWithdraw withdraw your assets into another address.
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
func (cl *Client) UserWithdraw(
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

	resBody, err := cl.callPrivate(MethodUserWithdraw, params)
	if err != nil {
		return nil, err
	}

	if cl.env.Debug >= 2 {
		fmt.Printf(">>> UserWithdraw: response body: %s\n", resBody)
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

//nolint: interfacer
func (cl *Client) trade(method, tipe, pair string, amount, price *big.Rat) (
	body []byte, err error,
) {
	assets := strings.Split(pair, "_")
	if len(assets) != 2 {
		return nil, fmt.Errorf("trade: invalid pair name %q", pair)
	}

	amountStr := amount.String()
	priceStr := price.String()

	params := map[string][]string{
		tokenomy.ParamNameOrderMethod: {method},
		tokenomy.ParamNamePair:        {pair},
		tokenomy.ParamNameType:        {tipe},
	}

	if method == tokenomy.TradeMethodLimit {
		params[tokenomy.ParamNamePrice] = []string{priceStr}
		params[assets[0]] = []string{amountStr}
	}

	body, err = cl.callPrivate(MethodTrade, params)
	if err != nil {
		return nil, err
	}

	return body, nil
}
