package template_rule_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

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

func ToGameRulesModel(entity template.GameRules) GameRules {
	return GameRules{
		HeroCountMin:                         entity.HeroCountMin,
		HeroCountMax:                         entity.HeroCountMax,
		HeroCountIncrement:                   entity.HeroCountIncrement,
		HeroHireBan:                          entity.HeroHireBan,
		EncounterHoles:                       entity.EncounterHoles,
		TournamentRules:                      entity.TournamentRules,
		Bonuses:                              ToBonusListModel(entity.Bonuses),
		WinConditions:                        ToWinConditionsModel(entity.WinConditions),
		GladiatorArena:                       entity.GladiatorArena,
		GladiatorArenaRegistrationStartWork:  entity.GladiatorArenaRegistrationStartWork,
		GladiatorArenaRegistrationStartFight: entity.GladiatorArenaRegistrationStartFight,
		GladiatorArenaDaysDelayStart:         entity.GladiatorArenaDaysDelayStart,
		GladiatorArenaCountDay:               entity.GladiatorArenaCountDay,
		ChampionSelectRule:                   entity.ChampionSelectRule,
		GlobalBans:                           helpers.MapPointer(entity.GlobalBans, ToGlobalBansModel),
		FactionLawsExpModifier:               entity.FactionLawsExpModifier,
		AstrologyExpModifier:                 entity.AstrologyExpModifier,
	}
}

func ToGameRulesEntity(model GameRules) template.GameRules {
	return template.GameRules{
		HeroCountMin:                         model.HeroCountMin,
		HeroCountMax:                         model.HeroCountMax,
		HeroCountIncrement:                   model.HeroCountIncrement,
		HeroHireBan:                          model.HeroHireBan,
		EncounterHoles:                       model.EncounterHoles,
		TournamentRules:                      model.TournamentRules,
		Bonuses:                              ToBonusListEntity(model.Bonuses),
		WinConditions:                        ToWinConditionsEntity(model.WinConditions),
		GladiatorArena:                       model.GladiatorArena,
		GladiatorArenaRegistrationStartWork:  model.GladiatorArenaRegistrationStartWork,
		GladiatorArenaRegistrationStartFight: model.GladiatorArenaRegistrationStartFight,
		GladiatorArenaDaysDelayStart:         model.GladiatorArenaDaysDelayStart,
		GladiatorArenaCountDay:               model.GladiatorArenaCountDay,
		ChampionSelectRule:                   model.ChampionSelectRule,
		GlobalBans:                           helpers.MapPointer(model.GlobalBans, ToGlobalBansEntity),
		FactionLawsExpModifier:               model.FactionLawsExpModifier,
		AstrologyExpModifier:                 model.AstrologyExpModifier,
	}
}
