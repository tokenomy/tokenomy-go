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

//
// WebSocketPublic define a WebSocket client for public APIs.
//
type WebSocketPublic struct {
	env  *tokenomy.Environment
	conn *websocket.Client

	requestsLocker sync.Mutex
	requests       map[uint64]chan *websocket.Response
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

	wsparams := &WebSocketParams{
		Pair: pair,
	}

	res, err := cl.send(http.MethodGet, apiMarketDepths, wsparams)
	if err != nil {
		return nil, err
	}

	resBody, err := base64.StdEncoding.DecodeString(res.Body)
	if err != nil {
		return nil, err
	}

	depths = &MarketDepths{}

	err = json.Unmarshal(resBody, depths)
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

	wsparams := &WebSocketParams{
		Pair: pair,
	}

	res, err := cl.send(http.MethodGet, apiMarketTicker, wsparams)
	if err != nil {
		return nil, err
	}

	resBody, err := base64.StdEncoding.DecodeString(res.Body)
	if err != nil {
		return nil, err
	}

	tick = &Tick{}

	err = json.Unmarshal(resBody, tick)
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

	wsparams := &WebSocketParams{
		Pair:   pair,
		Offset: offset,
		Limit:  limit,
	}

	res, err := cl.send(http.MethodGet, apiMarketTrades, wsparams)
	if err != nil {
		return nil, err
	}

	resBody, err := base64.StdEncoding.DecodeString(res.Body)
	if err != nil {
		return nil, err
	}

	marketTrades = &tokenomy.MarketTrades{}

	err = json.Unmarshal(resBody, marketTrades)
	if err != nil {
		return nil, err
	}

	return marketTrades, nil
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

	if res.ID != 0 {
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
