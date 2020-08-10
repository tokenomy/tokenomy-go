// Copyright 2020 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package v2

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/shuLhan/share/lib/websocket"
	"github.com/tokenomy/tokenomy-go"
)

const (
	maxQueue int = 256
)

//
// WebSocketPublic define a WebSocket client for public APIs.
//
type WebSocketPublic struct {
	env  *tokenomy.Environment
	conn *websocket.Client

	requestsLocker sync.Mutex
	requests       map[uint64]chan *websocket.Response
	subs           *tokenomy.PublicSubscription
	topicTrades    chan tokenomy.Trade

	// NotifTrades is a channel that will receive public order books
	// (open, closed, cancelled order) after calling SubscribeTrades
	// method.
	NotifTrades <-chan tokenomy.Trade
}

//
// NewWebSocketPublic create new WebSocket connection to public APIs.
//
func NewWebSocketPublic(env *tokenomy.Environment) (
	cl *WebSocketPublic, err error,
) {
	if env == nil {
		env = tokenomy.NewEnvironment("", "")
	}
	if len(env.Address) == 0 {
		env.Address = DefaultAddress
	}

	cl = &WebSocketPublic{
		env: env,
		conn: &websocket.Client{
			Headers: make(http.Header),
		},
		requests:    make(map[uint64](chan *websocket.Response)),
		subs:        &tokenomy.PublicSubscription{},
		topicTrades: make(chan tokenomy.Trade, maxQueue),
	}

	cl.NotifTrades = cl.topicTrades

	if env.IsInsecure {
		cl.conn.TLSConfig = &tls.Config{
			InsecureSkipVerify: env.IsInsecure,
		}
	}

	cl.conn.HandleText = cl.handleText
	cl.conn.HandleQuit = cl.handleUnexpectedQuit

	err = cl.connect()
	if err != nil {
		return nil, fmt.Errorf("NewWebSocketPublic: %w", err)
	}

	return cl, nil
}

//
// Close the connection and release all the resource.
//
func (cl *WebSocketPublic) Close() error {
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
// MarketDepths fetch list of market's depth for specific pair.
//
func (cl *WebSocketPublic) MarketDepths(pair string) (
	depths *MarketDepths, err error,
) {
	if len(pair) == 0 {
		return nil, tokenomy.ErrInvalidPair
	}

	wsparams := &tokenomy.WebSocketParams{
		TradeRequest: tokenomy.TradeRequest{
			Pair: pair,
		},
	}

	_, resbody, err := cl.send(http.MethodGet, apiMarketDepths, wsparams)
	if err != nil {
		return nil, err
	}

	depths = &MarketDepths{}

	err = json.Unmarshal(resbody, depths)
	if err != nil {
		return nil, err
	}

	return depths, nil
}

//
// MarketTicker return the ticker information on specific pair.
//
func (cl *WebSocketPublic) MarketTicker(pair string) (tick *Tick, err error) {
	if len(pair) == 0 {
		return nil, tokenomy.ErrInvalidPair
	}

	wsparams := &tokenomy.WebSocketParams{
		TradeRequest: tokenomy.TradeRequest{
			Pair: pair,
		},
	}

	_, resbody, err := cl.send(http.MethodGet, apiMarketTicker, wsparams)
	if err != nil {
		return nil, err
	}

	tick = &Tick{}

	err = json.Unmarshal(resbody, tick)
	if err != nil {
		return nil, err
	}

	return tick, nil
}

//
// MarketTrades return list of all completed trades in the market, specific to
// pair, grouped by ask and bid.
//
func (cl *WebSocketPublic) MarketTrades(pair string, offset, limit int64) (
	marketTrades *tokenomy.MarketTrades, err error,
) {
	if len(pair) == 0 {
		return nil, tokenomy.ErrInvalidPair
	}

	wsparams := &tokenomy.WebSocketParams{
		TradeRequest: tokenomy.TradeRequest{
			Pair: pair,
		},
		Offset: offset,
		Limit:  limit,
	}

	_, resbody, err := cl.send(http.MethodGet, apiMarketTrades, wsparams)
	if err != nil {
		return nil, err
	}

	marketTrades = &tokenomy.MarketTrades{}

	err = json.Unmarshal(resbody, marketTrades)
	if err != nil {
		return nil, err
	}

	return marketTrades, nil
}

//
// Subscription return the list and status of subscription.
//
func (cl *WebSocketPublic) Subscription() (*tokenomy.PublicSubscription, error) {
	_, resbody, err := cl.send(http.MethodGet, wsPublicSubscription, nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resbody, cl.subs)
	if err != nil {
		return nil, err
	}

	return cl.subs, nil
}

//
// SubscribeTrades subscribe to changes on public order books.
//
// Multiple calls on this method will not clear previously subscribed pairs.
// For example, if the first call subscribed to pair "X" and the second call
// subscribed to pair "Y", the client has two subscription: "X" and "Y", NOT
// "Y".
//
// The order books (open, closed, and/or cancelled) can be retrived from
// NotifTrades field.
//
func (cl *WebSocketPublic) SubscribeTrades(pairNames []string) (
	*tokenomy.PublicSubscription, error,
) {
	if len(pairNames) == 0 {
		return cl.subs, nil
	}

	wsparams := &tokenomy.WebSocketParams{
		PublicSubscription: tokenomy.PublicSubscription{
			Trades: pairNames,
		},
	}

	_, resbody, err := cl.send(http.MethodPost, wsPublicSubscription, wsparams)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resbody, cl.subs)
	if err != nil {
		return nil, err
	}

	return cl.subs, nil
}

//
// UnsubscribeTrades stop receiving broadcast notification on topic "trades"
// on specific pairs.
// If parameter is empty, it will unsubscribe all registered pairs.
//
// On success it will return the latest subscription.
//
func (cl *WebSocketPublic) UnsubscribeTrades(pairNames []string) (
	*tokenomy.PublicSubscription, error,
) {
	if len(pairNames) == 0 {
		pairNames = cl.subs.Trades
	}

	wsparams := &tokenomy.WebSocketParams{
		PublicSubscription: tokenomy.PublicSubscription{
			Trades: pairNames,
		},
	}

	_, resbody, err := cl.send(http.MethodDelete, wsPublicSubscription,
		wsparams)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resbody, cl.subs)
	if err != nil {
		return nil, err
	}

	return cl.subs, nil
}

func (cl *WebSocketPublic) connect() (err error) {
	cl.conn.Endpoint = cl.env.Address + wsPublicEndpoint

	err = cl.conn.Connect()
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	return nil
}

func (cl *WebSocketPublic) handleText(
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

	if res.ID == 0 {
		resbody, err := base64.StdEncoding.DecodeString(res.Body)
		if err != nil {
			log.Printf("handleText: broadcast %s: %s",
				res.Message, err)
			return nil
		}

		switch res.Message {
		case apiMarketTrades, apiMarketTradesOpen:
			trade := tokenomy.Trade{}
			err = json.Unmarshal(resbody, &trade)
			if err != nil {
				log.Printf("handleText: broadcast %s: %s",
					res.Message, err)
				return nil
			}
			cl.topicTrades <- trade
		}
	} else {
		chres := cl.requestPop(res.ID)
		if chres != nil {
			chres <- res
		}
		return nil
	}

	return nil
}

func (cl *WebSocketPublic) handleUnexpectedQuit() {
	log.Println("handleUnexpectedQuit: disconnected ...")
	for {
		err := cl.connect()
		if err != nil {
			log.Printf("connect: %s", err.Error())
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
	log.Println("handleUnexpectedQuit: reconnected ...")
}

func (cl *WebSocketPublic) requestPush(req *websocket.Request) (
	chres chan *websocket.Response,
) {
	chres = make(chan *websocket.Response, 1)
	cl.requestsLocker.Lock()
	cl.requests[req.ID] = chres
	cl.requestsLocker.Unlock()
	return chres
}

func (cl *WebSocketPublic) requestPop(id uint64) (
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

func (cl *WebSocketPublic) send(
	method, target string, wsparams *tokenomy.WebSocketParams,
) (
	res *websocket.Response, resbody []byte, err error,
) {
	var body []byte

	if wsparams != nil {
		body, err = wsparams.Pack()
		if err != nil {
			return nil, nil, err
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
		return nil, nil, err
	}

	chres := cl.requestPush(req)

	err = cl.conn.SendText(payload)
	if err != nil {
		cl.requestPop(req.ID)
		return nil, nil, err
	}

	res = <-chres

	if res.Code != http.StatusOK {
		return nil, nil, errors.New(res.Message)
	}

	resbody, err = base64.StdEncoding.DecodeString(res.Body)
	if err != nil {
		return res, resbody, err
	}

	return res, resbody, nil
}
