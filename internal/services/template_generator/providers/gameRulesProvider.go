package providers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

type GameRulesProvider struct{}

func NewGameRulesProvider() *GameRulesProvider {
	return &GameRulesProvider{}
}

func (this *GameRulesProvider) CreateGameRules(configuration config.GeneratorConfig) entities.GameRules {
	heroSettings := configuration.GetHeroSettings()

	return entities.GameRules{
		HeroCountMin:           heroSettings.HeroCountMin,
		HeroCountMax:           heroSettings.HeroCountMax,
		HeroCountIncrement:     heroSettings.HeroCountIncrement,
		HeroHireBan:            configuration.IsSingleHeroMode(),
		EncounterHoles:         false,
		FactionLawsExpModifier: percentToModifier(configuration.FactionLawsExpPercent),
		AstrologyExpModifier:   percentToModifier(configuration.AstrologyExpPercent),
		Bonuses: entities.BonusList{
			{SID: "add_bonus_hero_stat", ReceiverSide: -1, ReceiverFilter: "all_heroes", Parameters: []string{"movementBonus", "0"}},
		},
		WinConditions: this.createAdvancedWinConditions(configuration),
	}
}

func (this *GameRulesProvider) createAdvancedWinConditions(configuration config.GeneratorConfig) entities.WinConditions {
	victoryCondition := configuration.GetVictoryCondition()
	gameEndConditions := configuration.GetGameEndConditions()
	gladiatorRules := configuration.GetGladiatorArenaRules()
	tournamentRules := configuration.GetTournamentRules()

	useGladiator := gladiatorRules.Enabled || victoryCondition == "win_condition_4"
	winConditions := entities.WinConditions{
		Classic:          true,
		Desertion:        true,
		DesertionDay:     3,
		DesertionValue:   3000,
		HeroLighting:     true,
		HeroLightingDay:  1,
		LostStartCity:    gameEndConditions.LostStartCity || victoryCondition == "win_condition_3",
		LostStartCityDay: helpers.Clamp(gameEndConditions.LostStartCityDay, 1, 30),
		LostStartHero:    gameEndConditions.LostStartHero || useGladiator || configuration.GameMode == "SingleHero",
		CityHold:         gameEndConditions.CityHold || victoryCondition == "win_condition_5",
		CityHoldDays:     helpers.Clamp(gameEndConditions.CityHoldDays, 1, 30),
	}
	if useGladiator {
		winConditions.GladiatorArena = true
		winConditions.GladiatorArenaRegistrationStartFight = true
		winConditions.GladiatorArenaDaysDelayStart = helpers.Clamp(gladiatorRules.DaysDelayStart, 1, 60)
		winConditions.GladiatorArenaCountDay = helpers.Clamp(gladiatorRules.CountDay, 1, 30)
		winConditions.ChampionSelectRule = "StartHero"
	}
	if tournamentRules.Enabled || victoryCondition == "win_condition_6" {
		firstDay := helpers.Clamp(tournamentRules.FirstTournamentDay, 3, 60)
		interval := helpers.Clamp(tournamentRules.Interval, 3, 30)
		pointsToWin := helpers.Clamp(tournamentRules.PointsToWin, 1, 10)
		roundCount := pointsToWin*2 - 1
		winConditions.ChampionSelectRule = "StartHero"
		winConditions.Tournament = true
		winConditions.TournamentSaveArmy = true
		winConditions.TournamentPointsToWin = pointsToWin

		var announceDays, battleOffsets []int
		prevBattle := 0
		for i := range roundCount {
			announce := 1
			if i > 0 {
				announce = prevBattle + 1
			}
			offset := firstDay - 1
			if i > 0 {
				offset = interval - 1
			}
			announceDays = append(announceDays, announce)
			battleOffsets = append(battleOffsets, offset)
			prevBattle = announce + offset
		}
		winConditions.TournamentAnnounceDays = announceDays
		winConditions.TournamentDays = battleOffsets
	}
	return winConditions
}

func percentToModifier(percent int) float64 {
	return helpers.RoundWithPrecision(float64(helpers.Clamp(percent, 25, 200))/100.0, 2)
}
