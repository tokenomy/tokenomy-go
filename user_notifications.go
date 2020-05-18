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
	Deposit  bool `json:"deposit"`
	Login    bool `json:"login"`
	Trade    bool `json:"trade"`
	Withdraw bool `json:"withdraw"`
}
