package gamerules

// Bonus describes a starting bonus granted to a side / hero.
// `parameters` contents are type-dependent; all observed values are JSON strings.
type Bonus struct {
	SID            string   `json:"sid"`
	ReceiverSide   int      `json:"receiverSide"`
	ReceiverFilter string   `json:"receiverFilter,omitempty"`
	Parameters     []string `json:"parameters"`
}
