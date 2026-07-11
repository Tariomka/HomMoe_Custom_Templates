package template_rule

import "encoding/json"

// WinConditions enumerates every observed win-condition toggle and tuning value.
// All fields are optional in the source JSON; absent fields keep their Go zero value.
type WinConditions struct {
	Classic         bool `json:"classic,omitempty"`
	Desertion       bool `json:"desertion,omitempty"`
	DesertionDay    int  `json:"desertionDay,omitempty"`
	DesertionValue  int  `json:"desertionValue,omitempty"`
	HeroLighting    bool `json:"heroLighting,omitempty"`
	HeroLightingDay int  `json:"heroLightingDay,omitempty"`

	LostStartCity    bool `json:"lostStartCity,omitempty"`
	LostStartCityDay int  `json:"lostStartCityDay,omitempty"`
	LostStartHero    bool `json:"lostStartHero,omitempty"`

	CityHold     bool `json:"cityHold,omitempty"`
	CityHoldDays int  `json:"cityHoldDays,omitempty"`

	GladiatorArena                       bool   `json:"gladiatorArena,omitempty"`
	GladiatorArenaRegistrationStartWork  bool   `json:"gladiatorArenaRegistrationStartWork,omitempty"`
	GladiatorArenaRegistrationStartFight bool   `json:"gladiatorArenaRegistrationStartFight,omitempty"`
	GladiatorArenaDaysDelayStart         int    `json:"gladiatorArenaDaysDelayStart,omitempty"`
	GladiatorArenaCountDay               int    `json:"gladiatorArenaCountDay,omitempty"`
	ChampionSelectRule                   string `json:"championSelectRule,omitempty"`

	// Tournament-mode ("Chosen One", "Exodus", "Massacre", "Sprint") configuration.
	Tournament             bool  `json:"tournament,omitempty"`
	TournamentPointsToWin  int   `json:"tournamentPointsToWin,omitempty"`
	TournamentSaveArmy     bool  `json:"tournamentSaveArmy,omitempty"`
	TournamentDays         []int `json:"tournamentDays,omitempty"`
	TournamentAnnounceDays []int `json:"tournamentAnnounceDays,omitempty"`
}

// MergeWinConditionsIfDoesNotExist copies fields from source into the receiver, but only for fields
// where the receiver currently holds its zero value (so a nested `winConditions` block
// always wins over flat sibling keys).
func (this *WinConditions) MergeWinConditionsIfDoesNotExist(source WinConditions) {
	destinationBytes, err := json.Marshal(*this)
	if err != nil {
		return
	}

	sourceBytes, err := json.Marshal(source)
	if err != nil {
		return
	}

	var destinationMap, sourceMap map[string]json.RawMessage
	if err = json.Unmarshal(destinationBytes, &destinationMap); err != nil {
		return
	}

	if err = json.Unmarshal(sourceBytes, &sourceMap); err != nil {
		return
	}
	// Because every WinConditions field uses `omitempty`, any field present in
	// sourceMap but absent from destinationMap is a zero-valued destination field that should be
	// filled in from source.
	changed := false
	for k, v := range sourceMap {
		if _, ok := destinationMap[k]; !ok {
			destinationMap[k] = v
			changed = true
		}
	}
	if !changed {
		return
	}

	merged, err := json.Marshal(destinationMap)
	if err != nil {
		return
	}

	_ = json.Unmarshal(merged, this)
}
