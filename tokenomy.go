// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/shuLhan/share/lib/errors"
	liberrors "github.com/shuLhan/share/lib/errors"
)

// DefaultAddress contains the official API domain.
const DefaultAddress = "https://api.tokenomy.com"

// List of API endpoints.
const (
	APIMarketDepths     = "/v2/market/depths"
	APIMarketInfo       = "/v2/market/info"
	APIMarketTradesOpen = "/v2/market/trades/open"
	APIMarketPrices     = "/v2/market/prices"
	APIMarketTicker     = "/v2/market/ticker"
	APIMarketTrades     = "/v2/market/trades"
	APIMarketSummaries  = "/v2/market/summaries"

	APIUserInfo         = "/v2/user/info"
	APIUserTrades       = "/v2/user/trades"
	APIUserOrdersClosed = "/v2/user/orders/closed"
	APIUserOrdersOpen   = "/v2/user/orders/open"
	APIUserOrderInfo    = "/v2/user/order"
	APIUserTransactions = "/v2/user/transactions"
	APIUserWithdraw     = "/v2/user/withdraw"

	APITradeAsk       = "/v2/trade/ask"
	APITradeBid       = "/v2/trade/bid"
	APITradeBulk      = "/v2/trade/bulk"
	APITradeCancelAll = "/v2/trade/cancel/all"
	APITradeCancelAsk = "/v2/trade/cancel/ask"
	APITradeCancelBid = "/v2/trade/cancel/bid"

	WSPrivate = "/v2/user/ws"

	WSPublic             = "/v2/ws"
	WSPublicSubscription = "/v2/ws/subscription"
)

// List of WebSocket broadcast messages.
const (
	WSMessageUserOrdersTaken = "/v2/user/orders/taken"
)

// List of known asset names.
// The list is updated rarely, it may contains asset that has been delisted
// or did not contains new asset in the Tokenomy platform.
const (
	AssetNameAchain          = "achain"
	AssetNameBalancer        = "bal"
	AssetNameBinancechain    = "bnb"
	AssetNameBitcoin         = "btc"
	AssetNameBitcoinCash     = "bch"
	AssetNameCardano         = "ada"
	AssetNameChainlink       = "link"
	AssetNameCompound        = "comp"
	AssetNameDai             = "dai"
	AssetNameEos             = "eos"
	AssetNameEthereum        = "eth"
	AssetNameEthereumclassic = "etc"
	AssetNameHara            = "hart"
	AssetNameIdk             = "idk"
	AssetNameInmax           = "inx"
	AssetNameLitecoin        = "ltc"
	AssetNameNeo             = "neo"
	AssetNamePolkadot        = "dot"
	AssetNameSolana          = "sol"
	AssetNameStellar         = "xlm"
	AssetNameTether          = "usdt"
	AssetNameTezos           = "xtz"
	AssetNameTokenomy        = "ten"
	AssetNameUsdc            = "usdc"
	AssetNameVexanium        = "vex"
	AssetNameXanpool         = "xlp"
	AssetNameAvalanche       = "avax"
)

// List of valid pairs.
// The list is updated rarely, so it may contains pairs that has been delisted
// or did not contains new pairs in the Tokenomy platform.
const (
	PairBitcoinCashBitcoin = AssetNameBitcoinCash + `_` + AssetNameBitcoin // bch_btc
	PairEosBitcoin         = AssetNameEos + `_` + AssetNameBitcoin         // eos_btc
	PairEthereumBitcoin    = AssetNameEthereum + `_` + AssetNameBitcoin    // eth_btc
	PairLitecoinBitcoin    = AssetNameLitecoin + `_` + AssetNameBitcoin    // ltc_btc
	PairPolkadotBitcoin    = AssetNamePolkadot + `_` + AssetNameBitcoin    // dot_btc
	PairSolanaBitcoin      = AssetNameSolana + `_` + AssetNameBitcoin      // sol_btc
	PairStellarBitcoin     = AssetNameStellar + `_` + AssetNameBitcoin     // xlm_btc
	PairTokenomyBitcoin    = AssetNameTokenomy + `_` + AssetNameBitcoin    // ten_btc
	PairUsdcBitcoin        = AssetNameUsdc + `_` + AssetNameBitcoin        // usdc_btc
	PairVexaniumBitcoin    = AssetNameVexanium + `_` + AssetNameBitcoin    // vex_btc

	PairBitcoinIdk   = AssetNameBitcoin + `_` + AssetNameIdk   // btc_idk
	PairCardanoIdk   = AssetNameCardano + `_` + AssetNameIdk   // ada_idk
	PairChainlinkIdk = AssetNameChainlink + `_` + AssetNameIdk // link_idk
	PairCompoundIdk  = AssetNameCompound + `_` + AssetNameIdk  // comp_idk
	PairDaiIdk       = AssetNameDai + `_` + AssetNameIdk       // dai_idk
	PairEthereumIdk  = AssetNameEthereum + `_` + AssetNameIdk  // eth_idk
	PairPolkadotIdk  = AssetNamePolkadot + `_` + AssetNameIdk  // dot_idk
	PairSolanaIdk    = AssetNameSolana + `_` + AssetNameIdk    // sol_idk
	PairTetherIdk    = AssetNameTether + `_` + AssetNameIdk    // usdt_idk
	PairTezosIdk     = AssetNameTezos + `_` + AssetNameIdk     // xtz_idk
	PairTokenomyIdk  = AssetNameTokenomy + `_` + AssetNameIdk  // ten_idk

	PairCardanoTether  = AssetNameCardano + `_` + AssetNameTether  // ada_usdt
	PairBitcoinTether  = AssetNameBitcoin + `_` + AssetNameTether  // btc_usdt
	PairEthereumTether = AssetNameEthereum + `_` + AssetNameTether // eth_usdt
	PairIdkTether      = AssetNameIdk + `_` + AssetNameTether      // idk_usdt
	PairPolkadotTether = AssetNamePolkadot + `_` + AssetNameTether // dot_usdt
	PairSolanaTether   = AssetNameSolana + `_` + AssetNameTether   // sol_usdt
	PairTokenomyTether = AssetNameTokenomy + `_` + AssetNameTether // ten_usdt
	PairTezosTether    = AssetNameTezos + `_` + AssetNameTether    // xtz_usdt
)

// List of trade's method.
const (
	TradeMethodLimit  = "limit"
	TradeMethodMarket = "market"
)

// List of trade's type.
const (
	TradeTypeAsk = "sell"
	TradeTypeBid = "buy"
)

// List of valid values for TradeRequest.TimeInForce.
const (
	TimeInForceFOK = "FOK" // Fill-or-Kill.
)

// List of valid trade's status.
const (
	TradeStatusCancelled = "cancelled"
	TradeStatusFilled    = "filled"
)

// List of knowns environment variables.
const (
	EnvNameAddress = "TOKENOMY_ADDRESS"
	EnvNameDebug   = "TOKENOMY_DEBUG"
	EnvNameToken   = "TOKENOMY_TOKEN"
	EnvNameSecret  = "TOKENOMY_SECRET"
	EnvNameTestE2E = "TOKENOMY_TEST_E2E"
)

// List of knowns HTTP headers.
const (
	HeaderNameSign = "Sign"
	HeaderNameKey  = "Key"
)

// List of knowns parameter names.
const (
	ParamNameAddress       = "address"
	ParamNameAmount        = "amount"
	ParamNameAsset         = "asset"
	ParamNameIDAfter       = "id_after"
	ParamNameIDBefore      = "id_before"
	ParamNameLimit         = "limit"
	ParamNameMemo          = "memo"
	ParamNameMethod        = "method"
	ParamNameNetwork       = "network"
	ParamNameNonce         = "nonce"
	ParamNameOffset        = "offset"
	ParamNameOrderID       = "order_id"
	ParamNameOrderMethod   = "order_method"
	ParamNamePair          = "pair"
	ParamNamePostOnly      = "post_only"
	ParamNamePrice         = "price"
	ParamNameReceiveWindow = "recv_window"
	ParamNameRequestID     = "request_id"
	ParamNameSort          = "sort"
	ParamNameTimeAfter     = "time_after"
	ParamNameTimeBefore    = "time_before"
	ParamNameTimeInForce   = "time_in_force"
	ParamNameTimestamp     = "timestamp"
	ParamNameTradeID       = "trade_id"
	ParamNameTradeMethod   = "trade_method"
	ParamNameType          = "type"
)

// DefaultLimit define maximum number of record fetched per request.
const DefaultLimit = 100

// List of valid sort values.
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

	ErrTradeFillOrKill = &liberrors.E{
		Code:    http.StatusUnprocessableEntity,
		Message: "not enough amount in the market to process fill-or-kill order",
		Name:    "ERR_TRADE_FILL_OR_KILL",
	}

	ErrWalletAddress = &errors.E{
		Code:    http.StatusBadRequest,
		Message: "invalid or empty wallet address",
		Name:    "ERR_WALLET_ADDRESS",
	}
)

// Sign the payload using secret and return it as encoded hexadecimal
// characters.
func Sign(payload, secret string) string {
	hasher := hmac.New(sha512.New, []byte(secret))

	_, err := hasher.Write([]byte(payload))
	if err != nil {
		log.Fatal("Sign: ", err.Error())
	}

	signed := hasher.Sum(nil)

	return hex.EncodeToString(signed)
}

func timestamp() int64 {
	return time.Now().Unix()
}

func timestampAsString() string {
	return strconv.FormatInt(timestamp(), 10)
}
