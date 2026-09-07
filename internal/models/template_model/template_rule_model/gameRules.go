package template_rule_model

import "github.com/Tariomka/hommoe_custom_templates/internal/helpers"

type GameRules struct {
	HeroCountMin       int
	HeroCountMax       int
	HeroCountIncrement int
	HeroHireBan        bool
	EncounterHoles     bool
	TournamentRules    bool

	Bonuses       BonusList
	WinConditions WinConditions

	GladiatorArena                       bool
	GladiatorArenaRegistrationStartWork  bool
	GladiatorArenaRegistrationStartFight bool
	GladiatorArenaDaysDelayStart         int
	GladiatorArenaCountDay               int
	ChampionSelectRule                   string

	GlobalBans *GlobalBans

	FactionLawsExpModifier float64
	AstrologyExpModifier   float64
}

func (this GameRules) Clone() GameRules {
	clone := this
	clone.Bonuses = this.Bonuses.Clone()
	clone.WinConditions = this.WinConditions.Clone()
	clone.GlobalBans = helpers.MapPointer(this.GlobalBans, GlobalBans.Clone)
	return clone
}
