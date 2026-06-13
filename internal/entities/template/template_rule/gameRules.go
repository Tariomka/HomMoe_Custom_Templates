package template_rule

import "encoding/json"

// GameRules describes the gameplay rules for the template.
//
// A few templates (e.g. "Symmetry", "Jebus Outcast") declare gladiator-arena /
// global-bans fields as siblings of `winConditions` inside `gameRules` instead
// of nesting them properly; those fields are mirrored here for tolerant parsing.
type GameRules struct {
	HeroCountMin       int  `json:"heroCountMin"`
	HeroCountMax       int  `json:"heroCountMax"`
	HeroCountIncrement int  `json:"heroCountIncrement"`
	HeroHireBan        bool `json:"heroHireBan"`
	EncounterHoles     bool `json:"encounterHoles"`
	TournamentRules    bool `json:"tournamentRules,omitempty"`

	Bonuses       BonusList     `json:"bonuses,omitempty"`
	WinConditions WinConditions `json:"winConditions"`

	// Mirror of gladiator-arena fields when the template author placed them as
	// siblings of `winConditions` rather than inside it.
	GladiatorArena                       bool   `json:"gladiatorArena,omitempty"`
	GladiatorArenaRegistrationStartWork  bool   `json:"gladiatorArenaRegistrationStartWork,omitempty"`
	GladiatorArenaRegistrationStartFight bool   `json:"gladiatorArenaRegistrationStartFight,omitempty"`
	GladiatorArenaDaysDelayStart         int    `json:"gladiatorArenaDaysDelayStart,omitempty"`
	GladiatorArenaCountDay               int    `json:"gladiatorArenaCountDay,omitempty"`
	ChampionSelectRule                   string `json:"championSelectRule,omitempty"`

	// Some templates declare `globalBans` inside `gameRules` instead of at root.
	GlobalBans *GlobalBans `json:"globalBans,omitempty"`

	FactionLawsExpModifier float64 `json:"factionLawsExpModifier,omitempty"`
	AstrologyExpModifier   float64 `json:"astrologyExpModifier,omitempty"`
}

// UnmarshalJSON allows GameRules to also decode templates where the WinConditions
// fields appear flat inside `gameRules` (alongside `bonuses`, `encounterHoles`,
// etc.) instead of being nested under `winConditions`. Notable example: the
// "Zookeeper" template. The nested `winConditions` block, when present, still
// takes precedence; flat sibling fields only fill values that the nested form
// omitted.
func (this *GameRules) UnmarshalJSON(data []byte) error {
	// Decode the standard GameRules layout (including any nested winConditions).
	type alias GameRules
	if err := json.Unmarshal(data, (*alias)(this)); err != nil {
		return err
	}

	// Also try to populate WinConditions from the same JSON object treated as a
	// flat WinConditions blob (Zookeeper-style flat templates).
	var flat WinConditions
	if err := json.Unmarshal(data, &flat); err != nil {
		// flat decode failure is non-fatal - the regular decode already succeeded.
		return nil
	}
	this.WinConditions.MergeWinConditionsIfDoesNotExist(flat)
	return nil
}
