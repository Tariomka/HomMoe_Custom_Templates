package inner

// Road describes an in-zone road segment between two endpoints (typed references).
type Road struct {
	Type string `json:"type"`

	From TypedRef `json:"from"`
	To   TypedRef `json:"to"`

	Road                 *bool   `json:"road,omitempty"`
	SimTurnSquad         bool    `json:"simTurnSquad,omitempty"`
	GuardValue           int     `json:"guardValue,omitempty"`
	GuardWeeklyIncrement float64 `json:"guardWeeklyIncrement,omitempty"`
}
