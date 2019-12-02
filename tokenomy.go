// Copyright 2019 Tokenomy Technologies Pte. Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
)

//
// List of known asset names.
//
const (
	AssetName0x                  = "zrx"
	AssetNameAelf                = "elf"
	AssetNameAeternity           = "ae"
	AssetNameAppCoins            = "appc"
	AssetNameBasicAttentionToken = "bat"
	AssetNameBitTorrent          = "btt"
	AssetNameBitcoin             = "btc"
	AssetNameBitcoinABC          = "bchabc"
	AssetNameBitcoinSV           = "bchsv"
	AssetNameBread               = "brd"
	AssetNameComet               = "cmt"
	AssetNameDAEX                = "dax"
	AssetNameEOS                 = "eos"
	AssetNameEthereum            = "eth"
	AssetNameEthereumClassic     = "etc"
	AssetNameGifto               = "gto"
	AssetNameGolem               = "gnt"
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
	EnvNameDebug  = "TOKENOMY_DEBUG"
	EnvNameHost   = "TOKENOMY_HOST"
	EnvNameKey    = "TOKENOMY_KEY"
	EnvNameSecret = "TOKENOMY_SECRET"
)

//
// List of knowns HTTP headers.
//
const (
	HeaderNameSign = "Sign"
	HeaderNameKey  = "Key"
)

//
// List of knowns parameter names.
//
const (
	ParamNameAmount        = "amount"
	ParamNameAsset         = "asset"
	ParamNameIDAfter       = "id_after"
	ParamNameIDBefore      = "id_before"
	ParamNameLimit         = "limit"
	ParamNameMethod        = "method"
	ParamNameOffset        = "offset"
	ParamNameOrderID       = "order_id"
	ParamNamePair          = "pair"
	ParamNamePrice         = "price"
	ParamNameReceiveWindow = "recv_window"
	ParamNameSortIDBy      = "sort_id_by"
	ParamNameTimeAfter     = "time_after"
	ParamNameTimeBefore    = "time_before"
	ParamNameTimestamp     = "timestamp"
	ParamNameTradeMethod   = "trade_method"
	ParamNameType          = "type"
)

// DefaultLimit define maximum number of record fetched per request.
const DefaultLimit = 1000

// DefaultSort set the order of fecthed records to latest first.
const DefaultSort = "desc"

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
