// Copyright 2019 Tokenomy Technologies Pte. Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

//
// List of known asset (coin or token) names.
//
const (
	CurrencyName0x                  = "zrx"
	CurrencyNameAelf                = "elf"
	CurrencyNameAeternity           = "ae"
	CurrencyNameAppCoins            = "appc"
	CurrencyNameBasicAttentionToken = "bat"
	CurrencyNameBitTorrent          = "btt"
	CurrencyNameBitcoin             = "btc"
	CurrencyNameBitcoinABC          = "bchabc"
	CurrencyNameBitcoinSV           = "bchsv"
	CurrencyNameBread               = "brd"
	CurrencyNameComet               = "cmt"
	CurrencyNameDAEX                = "dax"
	CurrencyNameEOS                 = "eos"
	CurrencyNameEthereum            = "eth"
	CurrencyNameEthereumClassic     = "etc"
	CurrencyNameGifto               = "gto"
	CurrencyNameGolem               = "gnt"
	CurrencyNameHara                = "hart"
	CurrencyNameHonest              = "hnst"
	CurrencyNameIDK                 = "idk"
	CurrencyNameInmax               = "inx"
	CurrencyNameLYFE                = "lyfe"
	CurrencyNameLitecoin            = "ltc"
	CurrencyNameLoopring            = "lrc"
	CurrencyNameLoopringNeo         = "lrn"
	CurrencyNameMidasProtocol       = "mas"
	CurrencyNameMithril             = "mith"
	CurrencyNameMonero              = "xmr"
	CurrencyNameOPCoin              = "opc"
	CurrencyNameOmiseGO             = "omg"
	CurrencyNameOntology            = "ont"
	CurrencyNamePaxos               = "pax"
	CurrencyNamePlayGame            = "pxg"
	CurrencyNamePundiX              = "npxs"
	CurrencyNameQASH                = "qash"
	CurrencyNameRaidenNetwork       = "rdn"
	CurrencyNameSIX                 = "six"
	CurrencyNameStatus              = "snt"
	CurrencyNameStellar             = "xlm"
	CurrencyNameStoriqa             = "stq"
	CurrencyNameTRON                = "trx"
	CurrencyNameTVND                = "tvnd"
	CurrencyNameTether              = "usdt"
	CurrencyNameTokenomy            = "ten"
	CurrencyNameVEN                 = "ven"
	CurrencyNameVeritaseum          = "veri"
	CurrencyNameVexanium            = "vex"
	CurrencyNameZcash               = "zec"
	CurrencyNameZiliqa              = "zil"
)

//
// List of available pair between currencies.
//
const (
	PairBitcoinabcBitcoin  = "bchabc_btc"
	PairBitcoinsvBitcoin   = "bchsv_btc"
	PairBittorrentBitcoin  = "btt_btc"
	PairEosBitcoin         = "eos_btc"
	PairEthclassicBitcoin  = "etc_btc"
	PairEthereumBitcoin    = "eth_btc"
	PairHonestBitcoin      = "hnst_btc"
	PairLitecoinBitcoin    = "ltc_btc"
	PairLoopringneoBitcoin = "lrn_btc"
	PairLyfeBitcoin        = "lyfe_btc"
	PairMoneroBitcoin      = "xmr_btc"
	PairOntologyBtc        = "ont_btc"
	PairPlaygameBtc        = "pxg_btc"
	PairPundixBitcoin      = "npxs_btc"
	PairSixBtc             = "six_btc"
	PairStellarBitcoin     = "xlm_btc"
	PairStoriqaBitcoin     = "stq_btc"
	PairTokenomyBitcoin    = "ten_btc"
	PairTronBitcoin        = "trx_btc"
	PairVexaniumBitcoin    = "vex_btc"
	PairZcashBitcoin       = "zec_btc"

	PairBitcoinIdk = "btc_idk"
	PairTetherIdk  = "usdt_idk"

	PairHaraEthereum     = "hart_eth"
	PairInmaxEthereum    = "inx_eth"
	PairPundixEthereum   = "npxs_eth"
	PairStoriqaEthereum  = "stq_eth"
	PairTokenomyEthereum = "ten_eth"
	PairTronEthereum     = "trx_eth"
	PairVexaniumEthereum = "vex_ten"

	PairBitcoinTether  = "btc_usdt"
	PairDaexTether     = "dax_usdt"
	PairEthereumTether = "eth_usdt"
	PairTokenomyTether = "ten_usdt"

	PairSixTokenomy     = "six_ten"
	PairStoriqaTokenomy = "stq_ten"
)

//
// List of public APIs.
//
const (
	apiSummaries  = "/api/summaries"
	apiTicker     = "/api/%s/ticker"
	apiTrades     = "/api/%s/trades"
	apiOrderBook  = "/api/%s/depth"
	apiMarketInfo = "/api/market_info"
)

//
// List of "method" form value for private API.
//
const (
	apiTrade                  = "trade"
	apiTradeCancelOrder       = "cancelOrder"
	apiViewGetInfo            = "getInfo"
	apiViewGetOrder           = "getOrder"
	apiViewOpenOrders         = "openOrders"
	apiViewOrderHistory       = "orderHistory"
	apiViewTradeHistory       = "tradeHistory"
	apiViewTransactionHistory = "transHistory"
	apiWithdraw               = "withdrawCoin"
)

const (
	tradeMethodLimit  = "limit"
	tradeMethodMarket = "market"
)

const (
	tradeTypeBuy  = "buy"
	tradeTypeSell = "sell"
)

const (
	responseSuccess = 1
)

// List of common JSON field names.
const (
	fieldNameAmount            = "amount"
	fieldNameBalance           = "balance"
	fieldNameBaseCurrency      = "base_currency"
	fieldNameBaseCurrencyPrice = "base_currency_price"
	fieldNameDate              = "date"
	fieldNameError             = "error"
	fieldNameErrorCode         = "error_code"
	fieldNameFinishTime        = "finish_time"
	fieldNameIsError           = "is_error"
	fieldNameMethod            = "method"
	fieldNameOrderID           = "order_id"
	fieldNamePair              = "pair"
	fieldNamePrice             = "price"
	fieldNameStatus            = "status"
	fieldNameSubmitTime        = "submit_time"
	fieldNameSuccess           = "success"
	fieldNameTID               = "tid"
	fieldNameTradeID           = "trade_id"
	fieldNameTradeTime         = "trade_time"
	fieldNameTradeTimePrint    = "trade_time_print"
	fieldNameType              = "type"
)

var (
	// ErrUnauthenticated define an error when user did not provide token
	// and secret keys when accessing private APIs.
	ErrUnauthenticated = fmt.Errorf("unauthenticated connection")

	// ErrInvalidPairName define an error if user call API with empty,
	// invalid or unknown pair's name.
	ErrInvalidPairName = fmt.Errorf("invalid or empty pair name")
)

//
// jsonToMapStringFloat64 convert the map of string-interface{} into map of
// string-float64.
//
func jsonToMapStringFloat64(in map[string]interface{}) (
	out map[string]float64, err error,
) {
	out = make(map[string]float64, len(in))

	for k, v := range in {
		f64, err := strconv.ParseFloat(v.(string), 64)
		if err != nil {
			return nil, err
		}
		if f64 == 0 {
			continue
		}
		k = strings.ToLower(k)
		out[k] = f64
	}
	return out, nil
}

//
// timestamp return current time in milliseconds as integer.
//
func timestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

//
// timestampAsString return current time in milliseconds as string.
//
func timestampAsString() string {
	return strconv.FormatInt(timestamp(), 10)
}
