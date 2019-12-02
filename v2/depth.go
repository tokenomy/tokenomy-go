package v2

import "github.com/tokenomy/tokenomy-go"

//
// Depth contains total amount of remaining order grouped by price in open
// orders.
// Each depth is specific to pair.
//
type Depth struct {
	Amount tokenomy.Rawfloat `json:"amount"`
	Price  tokenomy.Rawfloat `json:"price"`
}
