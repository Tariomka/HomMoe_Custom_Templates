package gamerules

import "github.com/Tariomka/hommoe_custom_templates/internal/models/template/inner"

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
	GlobalBans *inner.GlobalBans `json:"globalBans,omitempty"`

	FactionLawsExpModifier float64 `json:"factionLawsExpModifier,omitempty"`
	AstrologyExpModifier   float64 `json:"astrologyExpModifier,omitempty"`
}
