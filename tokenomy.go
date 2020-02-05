// Copyright 2019 Tokenomy Technologies Pte. Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"net/http"

	"github.com/shuLhan/share/lib/errors"
)

//
// List of known asset names.
//
const (
	AssetName0x                  = "zrx"
	AssetNameAchain              = "achain"
	AssetNameAelf                = "elf"
	AssetNameAeternity           = "ae"
	AssetNameAppCoins            = "appc"
	AssetNameBasicAttentionToken = "bat"
	AssetNameBitTorrent          = "btt"
	AssetNameBitcoin             = "btc"
	AssetNameBitcoinABC          = "bchabc"
	AssetNameBitcoinSV           = "bchsv"
	AssetNameBinanceChain        = "bnb"
	AssetNameBread               = "brd"
	AssetNameComet               = "cmt"
	AssetNameDAEX                = "dax"
	AssetNameEOS                 = "eos"
	AssetNameEthereum            = "eth"
	AssetNameEthereumClassic     = "etc"
	AssetNameFgram               = "fgram"
	AssetNameGifto               = "gto"
	AssetNameGolem               = "gnt"
	AssetNameGram                = "gram"
	AssetNameHara                = "hart"
	AssetNameHonest              = "hnst"
	AssetNameIDK                 = "idk"
	AssetNameInmax               = "inx"
	AssetNameLYFE                = "lyfe"
	AssetNameLitecoin            = "ltc"
	AssetNameLoopring            = "lrc"
	AssetNameLoopringNeo         = "lrn"
	AssetNameMidasProtocol       = "mas"
	AssetNameMithril             = "mith"
	AssetNameMonero              = "xmr"
	AssetNameNeo                 = "neo"
	AssetNameOPCoin              = "opc"
	AssetNameOmiseGO             = "omg"
	AssetNameOntology            = "ont"
	AssetNamePaxos               = "pax"
	AssetNamePlayGame            = "pxg"
	AssetNamePundiX              = "npxs"
	AssetNameQASH                = "qash"
	AssetNameRaidenNetwork       = "rdn"
	AssetNameSIX                 = "six"
	AssetNameStatus              = "snt"
	AssetNameStellar             = "xlm"
	AssetNameStoriqa             = "stq"
	AssetNameTRON                = "trx"
	AssetNameTVND                = "tvnd"
	AssetNameTether              = "usdt"
	AssetNameTokenomy            = "ten"
	AssetNameVEN                 = "ven"
	AssetNameVeritaseum          = "veri"
	AssetNameVexanium            = "vex"
	AssetNameZcash               = "zec"
	AssetNameZiliqa              = "zil"
)

//
// List of valid pairs.
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
	PairOntologyBitcoin    = "ont_btc"
	PairPlaygameBitcoin    = "pxg_btc"
	PairPundixBitcoin      = "npxs_btc"
	PairSixBitcoin         = "six_btc"
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
// List of trade's method.
//
const (
	TradeMethodLimit  = "limit"
	TradeMethodMarket = "market"
)

//
// List of trade's type.
//
const (
	TradeTypeAsk = "sell"
	TradeTypeBid = "buy"
)

//
// List of knowns environment variables.
//
const (
	EnvNameAddress = "TOKENOMY_ADDRESS"
	EnvNameDebug   = "TOKENOMY_DEBUG"
	EnvNameToken   = "TOKENOMY_TOKEN"
	EnvNameSecret  = "TOKENOMY_SECRET"
)

//
// List of knowns HTTP headers.
//
const (
	HeaderNameContentType = "Content-Type"
	HeaderNameSign        = "Sign"
	HeaderNameKey         = "Key"
)

// List of known HTTP header values
const (
	ContentTypeForm = "application/x-www-form-urlencoded"
)

//
// List of knowns parameter names.
//
const (
	ParamNameAddress       = "address"
	ParamNameAmount        = "amount"
	ParamNameAsset         = "asset"
	ParamNameIDAfter       = "id_after"
	ParamNameIDBefore      = "id_before"
	ParamNameLimit         = "limit"
	ParamNameMemo          = "memo"
	ParamNameMethod        = "method"
	ParamNameOrderID       = "order_id"
	ParamNameOffset        = "offset"
	ParamNamePair          = "pair"
	ParamNamePrice         = "price"
	ParamNameReceiveWindow = "recv_window"
	ParamNameRequestID     = "request_id"
	ParamNameTimeAfter     = "time_after"
	ParamNameTimeBefore    = "time_before"
	ParamNameTimestamp     = "timestamp"
	ParamNameTradeID       = "trade_id"
	ParamNameTradeMethod   = "trade_method"
	ParamNameType          = "type"
)

// DefaultLimit define maximum number of record fetched per request.
const DefaultLimit = 1000

// List of valid sort.
const (
	SortAscending  = "asc"
	SortDescending = "desc"
)

// List of predefined errors.
var (
	ErrInvalidAmount = &errors.E{
		Code:    http.StatusBadRequest,
		Message: "invalid or empty amount parameter",
		Name:    "ERR_INVALID_AMOUNT",
	}
	ErrInvalidAsset = &errors.E{
		Code:    http.StatusBadRequest,
		Message: "invalid or empty asset parameter",
		Name:    "ERR_INVALID_ASSET",
	}
	ErrInvalidPair = &errors.E{
		Code:    http.StatusBadRequest,
		Message: "invalid or empty pair parameter",
		Name:    "ERR_INVALID_PAIR",
	}
	ErrInvalidPrice = &errors.E{
		Code:    http.StatusBadRequest,
		Message: "invalid or empty price parameter",
		Name:    "ERR_INVALID_PRICE",
	}
	ErrInvalidRequestID = &errors.E{
		Code:    http.StatusBadRequest,
		Message: "invalid or empty request ID",
		Name:    "ERR_INVALID_REQUEST_ID",
	}
	ErrInvalidSortBy = &errors.E{
		Code:    http.StatusBadRequest,
		Message: `invalid sort-by parameter, its either "asc" or "desc"`,
		Name:    "ERR_INVALID_SORT_BY",
	}
	ErrInvalidTradeID = &errors.E{
		Code:    http.StatusBadRequest,
		Message: "invalid trade ID",
		Name:    "ERR_INVALID_TRADE_ID",
	}
	ErrInvalidTradeMethod = &errors.E{
		Code:    http.StatusBadRequest,
		Message: `invalid or empty trade method, its either "limit" or "market"`,
		Name:    "ERR_INVALID_TRADE_METHOD",
	}

	ErrAssetKYCRequired = &errors.E{
		Code:    http.StatusForbidden,
		Message: `the traded asset require user account to finish KYC process`,
		Name:    "ERR_ASSET_KYC_REQUIRED",
	}
	ErrAssetCountryBlacklisted = &errors.E{
		Code:    http.StatusForbidden,
		Message: `the traded asset is not allowed in user account country`,
		Name:    "ERR_ASSET_COUNTRY_BLACKLISTED",
	}
	ErrAssetTermsRequired = &errors.E{
		Code:    http.StatusForbidden,
		Message: `the traded asset require user account to accept terms of sale`,
		Name:    "ERR_ASSET_TERMS_REQUIRED",
	}
)

//
// Sign sign the payload using secret and return it as encoded
// hexadecimal characters.
//
func Sign(payload, secret string) string {
	hasher := hmac.New(sha512.New, []byte(secret))

	_, err := hasher.Write([]byte(payload))
	if err != nil {
		return ""
	}

	signed := hasher.Sum(nil)

	return hex.EncodeToString(signed)
}
