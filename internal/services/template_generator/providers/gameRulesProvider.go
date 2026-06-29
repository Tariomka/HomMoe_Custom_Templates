package providers

import (
	"strconv"
	"strings"

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
		Bonuses:                this.createBonuses(configuration.Bonuses),
		WinConditions:          this.createAdvancedWinConditions(configuration),
	}
}

func (this *GameRulesProvider) createAdvancedWinConditions(configuration config.GeneratorConfig) entities.WinConditions {
	victoryCondition := configuration.GetVictoryCondition()
	gameEndConditions := configuration.GetGameEndConditions()
	gladiatorRules := configuration.GetGladiatorArenaRules()
	tournamentRules := configuration.GetTournamentRules()

	useGladiator := gladiatorRules.Enabled || victoryCondition == winConditionValues.FinalBattle
	winConditions := entities.WinConditions{
		Classic:          true,
		Desertion:        true,
		DesertionDay:     3,
		DesertionValue:   3000,
		HeroLighting:     true,
		HeroLightingDay:  1,
		LostStartCity:    gameEndConditions.LostStartCity || victoryCondition == winConditionValues.CapitalHold,
		LostStartCityDay: helpers.Clamp(gameEndConditions.LostStartCityDay, 1, 30),
		LostStartHero:    gameEndConditions.LostStartHero || useGladiator || configuration.GameMode == gameModes.SingleHero,
		CityHold:         gameEndConditions.CityHold || victoryCondition == winConditionValues.CityHold,
		CityHoldDays:     helpers.Clamp(gameEndConditions.CityHoldDays, 1, 30),
	}
	if useGladiator {
		winConditions.GladiatorArena = true
		winConditions.GladiatorArenaRegistrationStartFight = true
		winConditions.GladiatorArenaDaysDelayStart = helpers.Clamp(gladiatorRules.DaysDelayStart, 1, 60)
		winConditions.GladiatorArenaCountDay = helpers.Clamp(gladiatorRules.CountDay, 1, 30)
		winConditions.ChampionSelectRule = championSelectRules.StartHero
	}
	if tournamentRules.Enabled || victoryCondition == winConditionValues.Tournament {
		firstDay := helpers.Clamp(tournamentRules.FirstTournamentDay, 3, 60)
		interval := helpers.Clamp(tournamentRules.Interval, 3, 30)
		pointsToWin := helpers.Clamp(tournamentRules.PointsToWin, 1, 10)
		roundCount := pointsToWin*2 - 1
		winConditions.ChampionSelectRule = championSelectRules.StartHero
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

func (this *GameRulesProvider) createBonuses(bonusEntries []config.BonusEntry) entities.BonusList {
	bonuses := entities.BonusList{}
	for _, entry := range bonusEntries {
		bonuses = append(bonuses, expandBonusEntry(entry)...)
	}
	return bonuses
}

// expandBonusEntry turns a single UI bonus preset into the one or two raw Bonus
// objects the template needs, mirroring the parallel C# editor's
// BonusEntry.ToBonuses(). Every bonus targets all sides (receiverSide -1) and
// uses the entry's receiver filter ("start_hero" / "all_heroes").
func expandBonusEntry(entry config.BonusEntry) []entities.Bonus {
	bonus := func(sid string, parameters ...string) entities.Bonus {
		return entities.Bonus{
			SID:            sid,
			ReceiverSide:   -1,
			ReceiverFilter: entry.ReceiverFilter,
			Parameters:     parameters,
		}
	}

	switch entry.PresetType {
	case config.BonusTownPortalFree:
		return []entities.Bonus{
			bonus(mapBonuses.HeroSpell, spellSids.NeutralTownPortal),
			bonus(mapBonuses.HeroStat, "magicCostSidSet", spellSids.NeutralTownPortal, "-999", "0"),
		}
	case config.BonusSpell:
		bonuses := []entities.Bonus{bonus(mapBonuses.HeroSpell, entry.Param)}
		if entry.Param2 == "1" {
			bonuses = append(bonuses, bonus(mapBonuses.HeroStat, "magicCostSidSet", entry.Param, "-999", "0"))
		}
		return bonuses
	case config.BonusUnitMultiplier:
		return []entities.Bonus{bonus(mapBonuses.HeroUnitMultiplier, entry.Param)}
	case config.BonusMovementBonus:
		return []entities.Bonus{bonus(mapBonuses.HeroStat, "movementBonus", entry.Param)}
	case config.BonusStartingItem:
		return []entities.Bonus{bonus(mapBonuses.HeroItem, entry.Param)}
	case config.BonusStartingGold:
		return []entities.Bonus{bonus(mapBonuses.Resource, "gold", entry.Param)}
	case config.BonusStartingGems:
		return []entities.Bonus{bonus(mapBonuses.Resource, "gemstones", entry.Param)}
	case config.BonusStartingCrystals:
		return []entities.Bonus{bonus(mapBonuses.Resource, "crystals", entry.Param)}
	case config.BonusStartingMercury:
		return []entities.Bonus{bonus(mapBonuses.Resource, "mercury", entry.Param)}
	case config.BonusStartingWood:
		return []entities.Bonus{bonus(mapBonuses.Resource, "wood", entry.Param)}
	case config.BonusStartingOre:
		return []entities.Bonus{bonus(mapBonuses.Resource, "ore", entry.Param)}
	}
	return nil
}

// CreateValueOverrides parses the newline-separated "sid=guardValue" overrides
// edited in the UI into ValueOverride entries (one per valid line). Blank or
// unparseable lines are skipped. Variant -1 applies the override to all variants.
func (this *GameRulesProvider) CreateValueOverrides(configuration config.GeneratorConfig) []entities.ValueOverride {
	var overrides []entities.ValueOverride
	for _, line := range strings.Split(configuration.ValueOverridesText, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		equals := strings.Index(line, "=")
		if equals <= 0 {
			continue
		}
		sid := strings.TrimSpace(line[:equals])
		if sid == "" {
			continue
		}
		guardValue, err := strconv.Atoi(strings.TrimSpace(line[equals+1:]))
		if err != nil {
			continue
		}
		overrides = append(overrides, entities.ValueOverride{
			SID:        sid,
			Variant:    -1,
			GuardValue: guardValue,
		})
	}
	return overrides
}

// CreateGlobalBans turns the newline-separated banned item / magic SIDs edited
// in the UI into a GlobalBans block, or nil when nothing is banned.
func (this *GameRulesProvider) CreateGlobalBans(configuration config.GeneratorConfig) *entities.GlobalBans {
	items := parseSidLines(configuration.BannedItems)
	magics := parseSidLines(configuration.BannedMagics)
	if len(items) == 0 && len(magics) == 0 {
		return nil
	}
	return &entities.GlobalBans{Items: items, Magics: magics}
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
