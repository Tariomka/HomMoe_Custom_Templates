package template_rule

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
)

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
	// v2 reads `omitempty` as "would encode as an empty JSON value", so false and
	// 0 are written out; OmitZeroStructFields is what makes an absent key mean an
	// unset field again.
	omitZero := json.OmitZeroStructFields(true)

	destinationBytes, err := json.Marshal(*this, omitZero)
	if err != nil {
		return
	}

	sourceBytes, err := json.Marshal(source, omitZero)
	if err != nil {
		return
	}

	var destinationMap, sourceMap map[string]jsontext.Value
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
