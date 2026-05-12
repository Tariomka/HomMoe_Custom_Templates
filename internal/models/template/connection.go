package template

// Connection links two zones together
type Connection struct {
	Name     string `json:"name"`
	FromZone string `json:"from"`
	ToZone   string `json:"to"`
	Type     string `json:"connectionType"`

	IsSimTurnSquad bool `json:"simTurnSquad,omitempty"`
	IsRoad         bool `json:"road,omitempty"`

	GuardZone            string  `json:"guardZone,omitempty"`
	IsGuardEscape        bool    `json:"guardEscape,omitempty"`
	GuardValue           int     `json:"guardValue"`
	GuardRandomization   float32 `json:"guardRandomization,omitempty"`
	GuardWeeklyIncrement float32 `json:"guardWeeklyIncrement"`

	GatePlacement string `json:"gatePlacement,omitempty"`
}
