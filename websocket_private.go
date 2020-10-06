// Copyright 2020 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/shuLhan/share/lib/websocket"
)

// OrdersClosedHandler define a callback when receiving order closed
// broadcast from server.
type OrdersClosedHandler func(trade *Trade)

//
// WebSocketPrivate define the private WebSocket client for APIv2.
//
type WebSocketPrivate struct {
	// HandleOrdersClosed define the callback that will be called
	// automatically by client when one of the user's orders closed in the
	// market.
	HandleOrdersClosed OrdersClosedHandler

	env  *Environment
	conn *websocket.Client

	requestsLocker sync.Mutex
	requests       map[uint64]chan *websocket.Response
}

//
// NewWebSocketPrivate create and initialize new WebSocket connection to
// private endpoint.
//
func NewWebSocketPrivate(env *Environment) (
	cl *WebSocketPrivate, err error,
) {
	if env == nil {
		env = NewEnvironment("", "")
	}
	if len(env.Address) == 0 {
		env.Address = DefaultAddress
	}

	cl = &WebSocketPrivate{
		env: env,
		conn: &websocket.Client{
			Headers: make(http.Header),
		},
		requests: make(map[uint64](chan *websocket.Response)),
	}
	if env.IsInsecure {
		cl.conn.TLSConfig = &tls.Config{
			InsecureSkipVerify: env.IsInsecure,
		}
	}

	cl.conn.HandleText = cl.handleText
	cl.conn.HandleQuit = cl.handleUnexpectedQuit

	err = cl.connect()
	if err != nil {
		return nil, fmt.Errorf("NewWebSocketPrivate: %w", err)
	}

	return cl, nil
}

//
// Close the connection and release all the resource.
//
func (cl *WebSocketPrivate) Close() error {
	cl.requestsLocker.Lock()
	for id, ch := range cl.requests {
		ch <- nil
		close(ch)
		delete(cl.requests, id)
	}
	cl.requestsLocker.Unlock()

	return cl.conn.Close()
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
func (cl *WebSocketPrivate) TradeAsk(treq *TradeRequest) (
	trade *TradeResponse, err error,
) {
	if treq == nil {
		return nil, nil
	}
	_, wsparams, err := treq.Pack()
	if err != nil {
		return nil, err
	}

	return cl.sendTradeRequest(http.MethodPost, APITradeAsk, wsparams)
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
func (cl *WebSocketPrivate) TradeBid(treq *TradeRequest) (
	trade *TradeResponse, err error,
) {
	if treq == nil {
		return nil, nil
	}
	_, wsparams, err := treq.Pack()
	if err != nil {
		return nil, err
	}

	return cl.sendTradeRequest(http.MethodPost, APITradeBid, wsparams)
}

//
// TradeCancel cancel the open trade using ID and pair information in Trade.
//
func (cl *WebSocketPrivate) TradeCancel(trade *Trade) (
	*Trade, error,
) {
	if trade.ID <= 0 {
		return nil, ErrInvalidTradeID
	}
	if len(trade.Pair) == 0 {
		return nil, ErrInvalidPair
	}

	var (
		tradeResponse *TradeResponse
		err           error
	)

	switch trade.Type {
	case TradeTypeAsk:
		tradeResponse, err = cl.TradeCancelAsk(trade.Pair, trade.ID)
	case TradeTypeBid:
		tradeResponse, err = cl.TradeCancelBid(trade.Pair, trade.ID)
	default:
		return nil, ErrInvalidTradeType
	}
	if err != nil {
		return nil, err
	}
	return tradeResponse.Order, nil
}

//
// TradeCancelAll cancel all user's open ask and bid orders.
//
func (cl *WebSocketPrivate) TradeCancelAll() (
	trades []Trade, err error,
) {
	wsres, err := cl.send(http.MethodDelete, APITradeCancelAll, nil)
	if err != nil {
		return nil, err
	}

	resbody, err := base64.StdEncoding.DecodeString(wsres.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resbody, &trades)
	if err != nil {
		return nil, err
	}

	return trades, nil
}

//
// TradeCancelAsk cancel the specific open sell by pair and ID.
//
func (cl *WebSocketPrivate) TradeCancelAsk(pairName string, id int64) (
	trade *TradeResponse, err error,
) {
	if id <= 0 {
		return nil, ErrInvalidTradeID
	}
	wsparams := &WebSocketParams{
		TradeRequest: TradeRequest{
			Pair: pairName,
		},
		TradeID: id,
	}
	return cl.sendTradeRequest(http.MethodDelete, APITradeCancelAsk, wsparams)
}

//
// TradeCancelBid cancel the specific open buy by pair and ID.
//
func (cl *WebSocketPrivate) TradeCancelBid(pairName string, id int64) (
	trade *TradeResponse, err error,
) {
	if id <= 0 {
		return nil, ErrInvalidTradeID
	}
	wsparams := &WebSocketParams{
		TradeRequest: TradeRequest{
			Pair: pairName,
		},
		TradeID: id,
	}
	return cl.sendTradeRequest(http.MethodDelete, APITradeCancelBid, wsparams)
}

//
// UserInfo fetch the user information and balances.
//
func (cl *WebSocketPrivate) UserInfo() (user *User, err error) {
	res, err := cl.send(http.MethodGet, APIUserInfo, nil)
	if err != nil {
		return nil, err
	}

	resb, err := base64.StdEncoding.DecodeString(res.Body)
	if err != nil {
		return nil, err
	}

	user = &User{}

	err = json.Unmarshal(resb, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

//
// UserOrderInfo fetch a single user's trade information based on pair's name
// and trade ID.
//
func (cl *WebSocketPrivate) UserOrderInfo(pairName string, id int64) (
	trade *Trade, err error,
) {
	if len(pairName) == 0 {
		return nil, ErrInvalidPair
	}
	wsparams := &WebSocketParams{
		TradeRequest: TradeRequest{
			Pair: pairName,
		},
		TradeID: id,
	}

	res, err := cl.send(http.MethodGet, APIUserOrderInfo, wsparams)
	if err != nil {
		return nil, err
	}

	resb, err := base64.StdEncoding.DecodeString(res.Body)
	if err != nil {
		return nil, err
	}

	trade = &Trade{}

	err = json.Unmarshal(resb, trade)
	if err != nil {
		return nil, err
	}

	return trade, nil
}

//
// UserOrdersOpen fetch the user open orders based on pair's name.
//
func (cl *WebSocketPrivate) UserOrdersOpen(pairName string) (
	pairTradesOpen PairTradesOpen, err error,
) {
	wsparams := &WebSocketParams{
		TradeRequest: TradeRequest{
			Pair: pairName,
		},
	}

	res, err := cl.send(http.MethodGet, APIUserOrdersOpen, wsparams)
	if err != nil {
		return nil, err
	}

	resb, err := base64.StdEncoding.DecodeString(res.Body)
	if err != nil {
		return nil, err
	}

	pairTradesOpen = make(PairTradesOpen)

	err = json.Unmarshal(resb, &pairTradesOpen)
	if err != nil {
		return nil, err
	}

	return pairTradesOpen, nil
}

func (cl *WebSocketPrivate) connect() error {
	params := make(url.Values)

	params.Set(ParamNameTimestamp, timestampAsString())

	payload := params.Encode()
	sign := Sign(payload, cl.env.Secret)

	cl.conn.Endpoint = cl.env.Address + WSPrivate + "?" + payload

	cl.conn.Headers.Set(HeaderNameKey, cl.env.Token)
	cl.conn.Headers.Set(HeaderNameSign, sign)

	err := cl.conn.Connect()
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	return nil
}

func (cl *WebSocketPrivate) send(
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
		Body:   base64.StdEncoding.EncodeToString(body),
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

func (cl *WebSocketPrivate) sendTradeRequest(
	method, target string, wsparams *WebSocketParams,
) (
	trade *TradeResponse, err error,
) {
	res, err := cl.send(method, target, wsparams)
	if err != nil {
		return nil, err
	}

	resb, err := base64.StdEncoding.DecodeString(res.Body)
	if err != nil {
		return nil, err
	}

	trade = &TradeResponse{}

	err = json.Unmarshal(resb, trade)
	if err != nil {
		return nil, err
	}

	return trade, nil
}

func (cl *WebSocketPrivate) handleText(
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
	if res.Message == APIUserOrdersClosed {
		if cl.HandleOrdersClosed == nil {
			return nil
		}

		resb, err := base64.StdEncoding.DecodeString(res.Body)
		if err != nil {
			log.Printf("handleText: %s %s",
				APIUserOrdersClosed, err.Error())
			return nil
		}

		trade := &Trade{}
		err = json.Unmarshal(resb, trade)
		if err != nil {
			log.Printf("handleText: %s %s",
				APIUserOrdersClosed, err.Error())
			return nil
		}
		cl.HandleOrdersClosed(trade)
	}

	return nil
}

func (cl *WebSocketPrivate) handleUnexpectedQuit() {
	log.Println("handleUnexpectedQuit: disconnected ...")
	for {
		err := cl.connect()
		if err != nil {
			log.Printf("Connect: %s", err.Error())
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
	log.Println("handleUnexpectedQuit: reconnected ...")
}

func (cl *WebSocketPrivate) requestPush(req *websocket.Request) (
	chres chan *websocket.Response,
) {
	chres = make(chan *websocket.Response, 1)
	cl.requestsLocker.Lock()
	cl.requests[req.ID] = chres
	cl.requestsLocker.Unlock()
	return chres
}

func (cl *WebSocketPrivate) requestPop(id uint64) (
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
