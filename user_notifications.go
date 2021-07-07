// Copyright 2020 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

//
// UserNotifications contains user's status of notification in the system.
// If its true, user will receive notification (mostly by email), otherwise no
// notification will be send to user.
//
type UserNotifications struct {
	UserID                 int64 `json:"-"`
	Deposit                bool  `json:"deposit"`
	Login                  bool  `json:"login"`
	Withdraw               bool  `json:"withdraw"`
	Order                  bool  `json:"order"`
	FixedDeposit           bool  `json:"fixed_deposit"`
	FixedDepositNearMature bool  `json:"fixed_deposit_near_mature"`
	FixedDepositReward     bool  `json:"fixed_deposit_reward"`
	DualCurrencyInvest     bool  `json:"dual_currency_invest"`
	DualCurrencyReward     bool  `json:"dual_currency_reward"`
	Staking                bool  `json:"staking"`
	Unstaking              bool  `json:"unstaking"`
	StakingReward          bool  `json:"staking_reward"`
	Created                int64 `json:"-"`
	Updated                int64 `json:"-"`
}
