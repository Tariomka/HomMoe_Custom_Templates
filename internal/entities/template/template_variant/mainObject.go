package template_variant

// MainObject is a featured object placed in a zone (a player Spawn, a City, etc.).
type MainObject struct {
	Type string `json:"type"`

	Spawn string `json:"spawn,omitempty"`
	Owner string `json:"owner,omitempty"`

	RemoveGuardIfHasOwner bool `json:"removeGuardIfHasOwner,omitempty"`

	GuardChance          float64 `json:"guardChance"`
	GuardValue           int     `json:"guardValue"`
	GuardRandomization   float64 `json:"guardRandomization,omitempty"`
	GuardWeeklyIncrement float64 `json:"guardWeeklyIncrement"`

	BuildingsConstructionSid string `json:"buildingsConstructionSid,omitempty"`

	Faction  *TypedRef `json:"faction,omitempty"`
	Factions []string  `json:"factions,omitempty"`

	Placement     string   `json:"placement"`
	PlacementArgs []string `json:"placementArgs,omitempty"`

	HoldCityWinCon bool `json:"holdCityWinCon,omitempty"`

	// Marks a main object that must be picked up / captured for win conditions (Harmony).
	IsKeyObject bool `json:"isKeyObject,omitempty"`

	// Per-object random-hire / unit-increment toggles (Shamrock, Madness).
	EnableWeeklyUnitIncrement bool `json:"enableWeeklyUnitIncrement,omitempty"`
	InitialUnitIncrement      int  `json:"initialUnitIncrement,omitempty"`
}
