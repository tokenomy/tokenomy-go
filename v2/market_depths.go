package v2

//
// MarketDepths contains list of depth on open ask and bid orders.
//
type MarketDepths struct {
	Asks []Depth `json:"asks"`
	Bids []Depth `json:"bids"`
}
