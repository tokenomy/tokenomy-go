// Copyright 2020 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v2

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/shuLhan/share/lib/math/big"
	"github.com/shuLhan/share/lib/websocket"

	"github.com/tokenomy/tokenomy-go"
)

const (
	defPrivateWebSocketEndpoint = "wss://api.tokenomy.com/v2/user/ws"
)

// TradesClosedHandler define a callback when receiving trades closed
// broadcast from server.
type TradesClosedHandler func(trade *tokenomy.Trade)

//
// PrivateWebSocket define the WebSocket client for private usage.
//
type PrivateWebSocket struct {
	// HandleTradesClosed define the callback that will be called
	// automatically by client when one of the user's trades closed in the
	// market.
	HandleTradesClosed TradesClosedHandler

	conn *websocket.Client

	requestsLocker sync.Mutex
	requests       map[uint64]chan *websocket.Response
}

//
// NewPrivateWebSocket create and initialize new WebSocket connection to
// private endpoint.
//
func NewPrivateWebSocket(env *tokenomy.Environment) (
	cl *PrivateWebSocket, err error,
) {
	if env == nil {
		env = tokenomy.NewEnvironment("", "")
	}
	if len(env.Address) == 0 {
		env.Address = defPrivateWebSocketEndpoint
	}

	params := make(url.Values, 0)

	params.Set(tokenomy.ParamNameTimestamp, timestampAsString())

	payload := params.Encode()
	sign := tokenomy.Sign(payload, env.Secret)

	endpoint := env.Address + "?" + payload

	cl = &PrivateWebSocket{
		conn: &websocket.Client{
			Endpoint: endpoint,
			TLSConfig: &tls.Config{
				InsecureSkipVerify: env.IsInsecure,
			},
			Headers: make(http.Header),
		},
		requests: make(map[uint64](chan *websocket.Response)),
	}

	cl.conn.Headers.Set(tokenomy.HeaderNameKey, env.Token)
	cl.conn.Headers.Set(tokenomy.HeaderNameSign, sign)

	err = cl.conn.Connect()
	if err != nil {
		return nil, fmt.Errorf("NewPrivateWebSocket: %w", err)
	}

	cl.conn.HandleText = cl.handleText

	return cl, nil
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
func (cl *PrivateWebSocket) TradeAsk(
	method, pairName string, amount, price *big.Rat,
) (
	trade *tokenomy.TradeResponse, err error,
) {
	_, wsparams, err := generateTradeParams(method, pairName, amount, price)
	if err != nil {
		return nil, err
	}

	return cl.sendTradeRequest(http.MethodPost, apiTradeAsk, wsparams)
}

//
// TradeBid request to buy the coin on market with specific method, amount,
// and price.
// The method parameter define the mode of buy, its either "market" (default)
// or "limit".
// If the method is "market", it will only accept amount parameter, otherwise
// if the method is "limit", the amount and price must not be zero.
//
// The pairName parameter define the coin and base assets to be traded, in the
// following format: "coin_base".
//
// The amount parameter define the volume of coin we want to buy.
//
// The price parameter define the number of base that we want to buy the
// amount of coin.
//
func (cl *PrivateWebSocket) TradeBid(
	method, pairName string, amount, price *big.Rat,
) (
	trade *tokenomy.TradeResponse, err error,
) {
	_, wsparams, err := generateTradeParams(method, pairName, amount, price)
	if err != nil {
		return nil, err
	}

	return cl.sendTradeRequest(http.MethodPost, apiTradeBid, wsparams)
}

//
// TradeCancelAsk cancel the specific open sell by pair and ID.
//
func (cl *PrivateWebSocket) TradeCancelAsk(pairName string, id int64) (
	trade *tokenomy.TradeResponse, err error,
) {
	if id <= 0 {
		return nil, tokenomy.ErrInvalidTradeID
	}
	wsparams := &WebSocketParams{
		Pair:    pairName,
		TradeID: id,
	}
	return cl.sendTradeRequest(http.MethodPost, apiTradeCancelAsk, wsparams)
}

//
// TradeCancelBid cancel the specific open buy by pair and ID.
//
func (cl *PrivateWebSocket) TradeCancelBid(pairName string, id int64) (
	trade *tokenomy.TradeResponse, err error,
) {
	if id <= 0 {
		return nil, tokenomy.ErrInvalidTradeID
	}
	wsparams := &WebSocketParams{
		Pair:    pairName,
		TradeID: id,
	}
	return cl.sendTradeRequest(http.MethodPost, apiTradeCancelBid, wsparams)
}

//
// UserInfo fetch the user information and balances.
//
func (cl *PrivateWebSocket) UserInfo() (user *tokenomy.User, err error) {
	res, err := cl.send(http.MethodGet, apiUserInfo, nil)
	if err != nil {
		return nil, err
	}

	user = &tokenomy.User{}
	err = json.Unmarshal([]byte(res.Body), user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

//
// UserTradeInfo fetch a single user's trade information based on pair's name
// and trade ID.
//
func (cl *PrivateWebSocket) UserTradeInfo(pairName string, id int64) (
	trade *tokenomy.Trade, err error,
) {
	if len(pairName) == 0 {
		return nil, tokenomy.ErrInvalidPair
	}
	wsparams := &WebSocketParams{
		Pair:    pairName,
		TradeID: id,
	}

	res, err := cl.send(http.MethodGet, apiUserTrade, wsparams)
	if err != nil {
		return nil, err
	}

	trade = &tokenomy.Trade{}

	err = json.Unmarshal([]byte(res.Body), trade)
	if err != nil {
		return nil, err
	}

	return trade, nil
}

//
// UserTradesOpen fetch the user open trades based on pair's name.
//
func (cl *PrivateWebSocket) UserTradesOpen(pairName string) (
	pairTradeOpens PairTradeOpens, err error,
) {
	if len(pairName) == 0 {
		return nil, tokenomy.ErrInvalidPair
	}
	wsparams := &WebSocketParams{
		Pair: pairName,
	}

	res, err := cl.send(http.MethodGet, apiUserTradesOpen, wsparams)
	if err != nil {
		return nil, err
	}

	pairTradeOpens = make(PairTradeOpens, 0)

	err = json.Unmarshal([]byte(res.Body), &pairTradeOpens)
	if err != nil {
		return nil, err
	}

	return pairTradeOpens, nil
}

func (cl *PrivateWebSocket) send(
	method, target string, wsparams *WebSocketParams,
) (
	res *websocket.Response, err error,
) {
	var body []byte

	if wsparams != nil {
		body, err = wsparams.Pack()
		if err != nil {
			return nil, err
		}
	}

	req := &websocket.Request{
		ID:     uint64(time.Now().UnixNano()),
		Method: method,
		Target: target,
		Body:   string(body),
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	chres := cl.requestPush(req)

	err = cl.conn.SendText(payload)
	if err != nil {
		cl.requestPop(req.ID)
		return nil, err
	}

	res = <-chres

	if res.Code != http.StatusOK {
		return nil, errors.New(res.Message)
	}

	return res, nil
}

func (cl *PrivateWebSocket) sendTradeRequest(
	method, target string, wsparams *WebSocketParams,
) (
	trade *tokenomy.TradeResponse, err error,
) {
	res, err := cl.send(method, target, wsparams)
	if err != nil {
		return nil, err
	}

	trade = &tokenomy.TradeResponse{}

	err = json.Unmarshal([]byte(res.Body), trade)
	if err != nil {
		return nil, err
	}

	return trade, nil
}

func (cl *PrivateWebSocket) handleText(
	wsclient *websocket.Client, frame *websocket.Frame,
) (
	err error,
) {
	var (
		res     = &websocket.Response{}
		payload = frame.Payload()
	)

	err = json.Unmarshal(payload, res)
	if err != nil {
		log.Printf("handleText: %q: %s", payload, err.Error())
		return nil
	}

	if res.ID != 0 {
		chres := cl.requestPop(res.ID)
		if chres != nil {
			chres <- res
		}
		return nil
	}

	// Handle broadcast from server.
	switch res.Message {
	case apiUserTradesClosed:
		if cl.HandleTradesClosed == nil {
			return nil
		}

		trade := &tokenomy.Trade{}
		err = json.Unmarshal([]byte(res.Body), trade)
		if err != nil {
			log.Printf("handleText: %s %s",
				apiUserTradesClosed, err.Error())
			return nil
		}
		cl.HandleTradesClosed(trade)
	}

	return nil
}

func (cl *PrivateWebSocket) requestPush(req *websocket.Request) (
	chres chan *websocket.Response,
) {
	chres = make(chan *websocket.Response, 1)
	cl.requestsLocker.Lock()
	cl.requests[req.ID] = chres
	cl.requestsLocker.Unlock()
	return chres
}

func (cl *PrivateWebSocket) requestPop(id uint64) (
	chres chan *websocket.Response,
) {
	cl.requestsLocker.Lock()
	chres, ok := cl.requests[id]
	if ok {
		delete(cl.requests, id)
	}
	cl.requestsLocker.Unlock()
	return chres
}
