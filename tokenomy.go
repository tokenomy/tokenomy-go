// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/shuLhan/share/lib/errors"
)

//
// List of known asset names.
// The list is updated rarely, it may contains asset that has been delisted
// or did not contains new asset in the Tokenomy server.
// Last update on 2020-06-11.
//
const (
	AssetNameAchain          = "achain"
	AssetNameBinancechain    = "bnb"
	AssetNameBitcoin         = "btc"
	AssetNameBitcoinabc      = "bchabc"
	AssetNameCybermiles      = "cmt"
	AssetNameDai             = "dai"
	AssetNameEos             = "eos"
	AssetNameEthereum        = "eth"
	AssetNameEthereumclassic = "etc"
	AssetNameHara            = "hart"
	AssetNameHonest          = "hnst"
	AssetNameIdk             = "idk"
	AssetNameInmax           = "inx"
	AssetNameLitecoin        = "ltc"
	AssetNameLoopring        = "lrc"
	AssetNameLoopringNeo     = "lrn"
	AssetNameLyfe            = "lyfe"
	AssetNameLyfebep         = "lyfebep"
	AssetNameMaker           = "mkr"
	AssetNameMonero          = "xmr"
	AssetNameNeo             = "neo"
	AssetNameOntology        = "ont"
	AssetNamePlaygame        = "pxg"
	AssetNameSiacash         = "scc"
	AssetNameSix             = "six"
	AssetNameStellar         = "xlm"
	AssetNameSwipe           = "swipe"
	AssetNameTether          = "usdt"
	AssetNameTokenomy        = "ten"
	AssetNameTron            = "trx"
	AssetNameVexanium        = "vex"
	AssetNameXanpool         = "xlp"
	AssetNameZcash           = "zec"
)

//
// List of valid pairs.
// The list is updated rarely, so it may contains pairs that has been delisted
// or did not contains new pairs in the Tokenomy server.
// Last update on 2020-06-11.
//
const (
	PairBitcoinabcBitcoin      = AssetNameBitcoinabc + "_" + AssetNameBitcoin
	PairEosBitcoin             = AssetNameEos + "_" + AssetNameBitcoin
	PairEthereumBitcoin        = AssetNameEthereum + "_" + AssetNameBitcoin
	PairEthereumclassicBitcoin = AssetNameEthereumclassic + "_" + AssetNameBitcoin
	PairHonestBitcoin          = AssetNameHonest + "_" + AssetNameBitcoin
	PairLitecoinBitcoin        = AssetNameLitecoin + "_" + AssetNameBitcoin
	PairLoopringBitcoin        = AssetNameLoopring + "_" + AssetNameBitcoin
	PairMoneroBitcoin          = AssetNameMonero + "_" + AssetNameBitcoin
	PairOntologyBitcoin        = AssetNameOntology + "_" + AssetNameBitcoin
	PairSixBitcoin             = AssetNameSix + "_" + AssetNameBitcoin
	PairStellarBitcoin         = AssetNameStellar + "_" + AssetNameBitcoin
	PairSwipeBitcoin           = AssetNameSwipe + "_" + AssetNameBitcoin
	PairTokenomyBitcoin        = AssetNameTokenomy + "_" + AssetNameBitcoin
	PairTronBitcoin            = AssetNameTron + "_" + AssetNameBitcoin
	PairVexaniumBitcoin        = AssetNameVexanium + "_" + AssetNameBitcoin
	PairXanpoolBitcoin         = AssetNameXanpool + "_" + AssetNameBitcoin
	PairZcashBitcoin           = AssetNameZcash + "_" + AssetNameBitcoin

	PairBitcoinIdk  = AssetNameBitcoin + "_" + AssetNameIdk
	PairDaiIdk      = AssetNameDai + "_" + AssetNameIdk
	PairHaraIdk     = AssetNameHara + "_" + AssetNameIdk
	PairHonestIdk   = AssetNameHonest + "_" + AssetNameIdk
	PairInmaxIdk    = AssetNameInmax + "_" + AssetNameIdk
	PairLyfebepIdk  = AssetNameLyfebep + "_" + AssetNameIdk
	PairMakerIdk    = AssetNameMaker + "_" + AssetNameIdk
	PairPlaygameIdk = AssetNamePlaygame + "_" + AssetNameIdk
	PairSiacashIdk  = AssetNameSiacash + "_" + AssetNameIdk
	PairSwipeIdk    = AssetNameSwipe + "_" + AssetNameIdk
	PairTetherIdk   = AssetNameTether + "_" + AssetNameIdk
	PairTokenomyIdk = AssetNameTokenomy + "_" + AssetNameIdk

	PairBitcoinTether  = AssetNameBitcoin + "_" + AssetNameTether
	PairEthereumTether = AssetNameEthereum + "_" + AssetNameTether
	PairIdkTether      = AssetNameIdk + "_" + AssetNameTether
	PairTokenomyTether = AssetNameTokenomy + "_" + AssetNameTether
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
// List of valid trade's status.
//
const (
	TradeStatusCancelled = "cancelled"
	TradeStatusFilled    = "filled"
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
	HeaderNameSign = "Sign"
	HeaderNameKey  = "Key"
)

// List of known HTTP header values.
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
	ParamNameNonce         = "nonce"
	ParamNameOrderID       = "order_id"
	ParamNameOrderMethod   = "order_method"
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

// DefaultDialTimeout define maximum time waiting for connection to be
// fully accepted.
var DefaultDialTimeout = 10 * time.Second

// DefaultTimeout define maximum time waiting for response in each HTTP
// or WebSocket request.
var DefaultTimeout = 16 * time.Second

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
	ErrInvalidTradeType = &errors.E{
		Code:    http.StatusBadRequest,
		Message: `invalid or empty trade type, its either "buy" or "sell"`,
		Name:    "ERR_INVALID_TRADE_TYPE",
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

	ErrWalletAddress = &errors.E{
		Code:    http.StatusBadRequest,
		Message: "invalid or empty wallet address",
		Name:    "ERR_WALLET_ADDRESS",
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
