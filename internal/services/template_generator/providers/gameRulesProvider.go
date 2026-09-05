package providers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/provider_interfaces"
)

type GameRulesProvider struct{}

func NewGameRulesProvider() provider_interfaces.IGameRulesProvider {
	return &GameRulesProvider{}
}

func (this *GameRulesProvider) CreateGameRules(configuration config.GeneratorConfig) template_model.GameRules {
	heroSettings := configuration.GetHeroSettings()

	return template_model.GameRules{
		HeroCountMin:           heroSettings.HeroCountMin,
		HeroCountMax:           heroSettings.HeroCountMax,
		HeroCountIncrement:     heroSettings.HeroCountIncrement,
		HeroHireBan:            configuration.IsSingleHeroMode(),
		EncounterHoles:         false,
		FactionLawsExpModifier: percentToModifier(configuration.FactionLawsExpPercent),
		AstrologyExpModifier:   percentToModifier(configuration.AstrologyExpPercent),
		Bonuses:                this.createBonuses(configuration.Bonuses),
		WinConditions:          this.createAdvancedWinConditions(configuration),
	}
}

// CreateValueOverrides parses the newline-separated "sid=guardValue" overrides
// edited in the UI into ValueOverride entries (one per valid line). Blank lines
// are skipped silently; every other rejected line produces a warning so the
// edit is not discarded without telling the user. Variant -1 applies the
// override to all variants.
func (this *GameRulesProvider) CreateValueOverrides(
	configuration config.GeneratorConfig) ([]template_model.ValueOverride, []string) {
	var overrides []template_model.ValueOverride
	var warnings []string
	for index, rawLine := range strings.Split(configuration.ValueOverridesText, "\n") {
		line := strings.TrimSpace(rawLine)
		if line == "" {
			continue
		}
		override, err := this.parseValueOverride(line, index+1)
		if err != nil {
			warnings = append(warnings, err.Error())
			continue
		}
		overrides = append(overrides, override)
	}
	return overrides, warnings
}

// CreateGlobalBans turns the newline-separated banned item / magic SIDs edited
// in the UI into a GlobalBans block, or nil when nothing is banned.
func (this *GameRulesProvider) CreateGlobalBans(configuration config.GeneratorConfig) *template_model.GlobalBans {
	items := parseSidLines(configuration.BannedItems)
	magics := parseSidLines(configuration.BannedMagics)
	if len(items) == 0 && len(magics) == 0 {
		return nil
	}
	return &template_model.GlobalBans{Items: items, Magics: magics}
}

func (this *GameRulesProvider) parseValueOverride(line string, lineNumber int) (template_model.ValueOverride, error) {
	equals := strings.Index(line, "=")
	if equals <= 0 {
		return template_model.ValueOverride{}, fmt.Errorf("line %d: '%s' is not sid=value", lineNumber, line)
	}

	guardValue, err := strconv.Atoi(strings.TrimSpace(line[equals+1:]))
	if err != nil {
		return template_model.ValueOverride{}, fmt.Errorf("line %d: '%s' has a non-numeric value", lineNumber, line)
	}

	sid := strings.TrimSpace(line[:equals])
	return template_model.ValueOverride{SID: sid, Variant: -1, GuardValue: guardValue}, nil
}

func (this *GameRulesProvider) createAdvancedWinConditions(
	configuration config.GeneratorConfig) template_model.WinConditions {
	victoryCondition := configuration.GetVictoryCondition()
	gameEndConditions := configuration.GetGameEndConditions()
	gladiatorRules := configuration.GetGladiatorArenaRules()
	tournamentRules := configuration.GetTournamentRules()

	useGladiator := configuration.IsGladiatorArenaMode()

	winConditionValues := registry.GetWinningConditionValues()
	winConditions := template_model.WinConditions{
		Classic:          true,
		Desertion:        true,
		DesertionDay:     3,
		DesertionValue:   3000,
		HeroLighting:     true,
		HeroLightingDay:  1,
		LostStartCity:    gameEndConditions.LostStartCity || victoryCondition == winConditionValues.CapitalHold,
		LostStartCityDay: helpers.Clamp(gameEndConditions.LostStartCityDay, 1, 30),
		LostStartHero: gameEndConditions.LostStartHero ||
			useGladiator ||
			configuration.GameMode == registry.GetGameModeValues().SingleHero,
		CityHold:     gameEndConditions.CityHold || victoryCondition == winConditionValues.CityHold,
		CityHoldDays: helpers.Clamp(gameEndConditions.CityHoldDays, 1, 30),
	}
	if useGladiator {
		this.setGladiatorArenaRules(gladiatorRules, &winConditions)
	}
	if tournamentRules.Enabled || victoryCondition == winConditionValues.Tournament {
		this.setTournamentRules(tournamentRules, &winConditions)
	}
	return winConditions
}

func (this *GameRulesProvider) setGladiatorArenaRules(
	gladiatorRules config.GladiatorArenaRules,
	winConditions *template_model.WinConditions) {
	winConditions.GladiatorArena = true
	winConditions.GladiatorArenaRegistrationStartFight = true
	winConditions.GladiatorArenaDaysDelayStart = helpers.Clamp(gladiatorRules.DaysDelayStart, 1, 60)
	winConditions.GladiatorArenaCountDay = helpers.Clamp(gladiatorRules.CountDay, 1, 30)
	winConditions.ChampionSelectRule = registry.GetChampionSelectValues().StartHero
}

func (this *GameRulesProvider) setTournamentRules(
	tournamentRules config.TournamentRules,
	winConditions *template_model.WinConditions) {
	firstDay := helpers.Clamp(tournamentRules.FirstTournamentDay, 3, 60)
	interval := helpers.Clamp(tournamentRules.Interval, 3, 30)
	pointsToWin := helpers.Clamp(tournamentRules.PointsToWin, 1, 10)
	roundCount := pointsToWin*2 - 1
	winConditions.ChampionSelectRule = registry.GetChampionSelectValues().StartHero
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

func (this *GameRulesProvider) createBonuses(bonusEntries []config.BonusEntry) template_model.BonusList {
	bonuses := template_model.BonusList{}
	for _, entry := range bonusEntries {
		bonuses = append(bonuses, expandBonusEntry(entry)...)
	}
	return bonuses
}

func expandBonusEntry(entry config.BonusEntry) []template_model.Bonus {
	bonus := func(sid string, parameters ...string) template_model.Bonus {
		return template_model.Bonus{
			SID:            sid,
			ReceiverSide:   -1,
			ReceiverFilter: entry.ReceiverFilter,
			Parameters:     parameters,
		}
	}

	mapBonuses := registry.GetMapBonusesValues()
	switch entry.PresetType {
	case config.BonusTownPortalFree:
		highNeutralSpells := registry.GetHighNeutralSpellSidValues()
		return []template_model.Bonus{
			bonus(mapBonuses.HeroSpell, highNeutralSpells.TownPortal),
			bonus(mapBonuses.HeroStat, "magicCostSidSet", highNeutralSpells.TownPortal, "-999", "0"),
		}
	case config.BonusSpell:
		bonuses := []template_model.Bonus{bonus(mapBonuses.HeroSpell, entry.Param)}
		if entry.Param2 == "1" {
			bonuses = append(bonuses, bonus(mapBonuses.HeroStat, "magicCostSidSet", entry.Param, "-999", "0"))
		}
		return bonuses
	case config.BonusUnitMultiplier:
		return []template_model.Bonus{bonus(mapBonuses.HeroUnitMultiplier, entry.Param)}
	case config.BonusMovementBonus:
		return []template_model.Bonus{bonus(mapBonuses.HeroStat, "movementBonus", entry.Param)}
	case config.BonusStartingItem:
		return []template_model.Bonus{bonus(mapBonuses.HeroItem, entry.Param)}
	case config.BonusStartingGold:
		return []template_model.Bonus{bonus(mapBonuses.Resource, "gold", entry.Param)}
	case config.BonusStartingGems:
		return []template_model.Bonus{bonus(mapBonuses.Resource, "gemstones", entry.Param)}
	case config.BonusStartingCrystals:
		return []template_model.Bonus{bonus(mapBonuses.Resource, "crystals", entry.Param)}
	case config.BonusStartingMercury:
		return []template_model.Bonus{bonus(mapBonuses.Resource, "mercury", entry.Param)}
	case config.BonusStartingWood:
		return []template_model.Bonus{bonus(mapBonuses.Resource, "wood", entry.Param)}
	case config.BonusStartingOre:
		return []template_model.Bonus{bonus(mapBonuses.Resource, "ore", entry.Param)}
	}
	return nil
}

// parseSidLines splits a newline-separated SID list into trimmed, non-empty SIDs.
func parseSidLines(raw string) []string {
	var sids []string
	for line := range strings.SplitSeq(raw, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			sids = append(sids, line)
		}
	}
	return sids
}

func percentToModifier(percent int) float64 {
	return helpers.RoundWithPrecision(float64(helpers.Clamp(percent, 25, 200))/100.0, 2)
}
