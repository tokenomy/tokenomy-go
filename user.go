// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

//
// User contains user information including profile, balances, and
// frozen balances.
//
type User struct {
	*UserAssets

	Notifications *UserNotifications `json:"notifications,omitempty"`
	Wallets       UserWallets        `json:"wallets,omitempty"`
	Email         string             `json:"email"`
	FullName      string             `json:"full_name"`
	ID            int64              `json:"id"`
}
