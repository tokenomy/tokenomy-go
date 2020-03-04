// Copyright 2019 Tokenomy Technologies Ltd. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package tokenomy

import (
	"strconv"
	"strings"
)

const (
	// Maximum precision when converting the Rawfloat value to string.
	maxPrecision = 8
)

//
// Rawfloat represent internal float64 with custom String and marshaling for
// JSON.
// The rules for converting the float to string are,
//
//	(1) If the value is zero it should return "0", not "0.000000"
//	(2) If the value does not have mantissa, it should return only the
//	base value without precision.  For example 123.00 must be printed as
//	"123".
//	(3) If the value contains more than 8 digits mantissa, then
//	(3.1) if one of last 8 digits in mantissa is not zero, then the
//	printed value should be limited to 8 digits only.  For example,
//	0.000_000_016 should be printed as "0.00000001" without rounding.
//	(3.2) if all of last 8 digits mantissa is zero, the printed value
//	should be limited to the last non-zero.  For example 0.000_000_0001
//	should be printed as "0.0000000001" not "0".
//
type Rawfloat float64

//
// ParseRawfloat convert the string `s` into Rawfloat with precision set to 64
// bit.
//
func ParseRawfloat(s string) (Rawfloat, error) {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return Rawfloat(v), nil
}

//
// MarshalJSON convert the Rawfloat value into specific format limited to 8
// digits precision with the following exceptions: zero value will be returned
// as "0", and trailing zero digits at precision will be removed.
//
func (f Rawfloat) MarshalJSON() ([]byte, error) {
	return []byte(f.String()), nil
}

//
// String convert the Rawfloat to string.
//
func (f Rawfloat) String() (s string) {
	// Rule (1).
	if f == 0 {
		return "0"
	}

	// Rule (2), automatically handled by FormatFloat.
	s = strconv.FormatFloat(float64(f), 'f', -1, 64)

	decimalIndex := strings.IndexByte(s, '.')

	if decimalIndex < 0 {
		return s
	}

	s = strings.TrimRight(s, "0")

	lastZero := 0
	for x := decimalIndex + 1; x < len(s); x++ {
		if s[x] != '0' {
			lastZero = x
			break
		}
	}

	decimalPrecision := decimalIndex + maxPrecision + 1
	if lastZero < decimalPrecision {
		// Rule (3.1).
		if len(s) > decimalPrecision {
			s = s[:decimalPrecision]
		}
	} else {
		// Rule (3.2)
		s = s[:lastZero+1]
	}

	s = strings.TrimRight(s, ".")

	return s
}
