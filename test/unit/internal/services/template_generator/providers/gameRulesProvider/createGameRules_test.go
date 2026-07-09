package gameRulesProvider_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

// ── Default configuration golden ─────────────────────────────────────

func TestWhenDefaultConfiguration_ReturnsExpectedGameRules(t *testing.T) {
	// Arrange
	expected := entities.GameRules{
		HeroCountMin:           4,
		HeroCountMax:           8,
		HeroCountIncrement:     1,
		HeroHireBan:            false,
		EncounterHoles:         false,
		FactionLawsExpModifier: 1,
		AstrologyExpModifier:   1,
		Bonuses:                entities.BonusList{},
		WinConditions: entities.WinConditions{
			Classic:          true,
			Desertion:        true,
			DesertionDay:     3,
			DesertionValue:   3000,
			HeroLighting:     true,
			HeroLightingDay:  1,
			LostStartCity:    false,
			LostStartCityDay: 3,
			LostStartHero:    false,
			CityHold:         false,
			CityHoldDays:     6,
		},
	}

	// Act
	actual := createGameRules(nil)

	// Assert
	assert.Equal(t, expected, actual)
}

// ── Hero settings ────────────────────────────────────────────────────

func TestWhenGameModeIsClassic_DisablesHeroHireBan(t *testing.T) {
	// Arrange & Act
	actual := createGameRules(func(configuration *config.GeneratorConfig) {
		configuration.GameMode = "Classic"
	})

	// Assert
	assert.False(t, actual.HeroHireBan)
}

func TestWhenGameModeIsSingleHero_EnablesHeroHireBan(t *testing.T) {
	// Arrange & Act
	actual := createGameRules(func(configuration *config.GeneratorConfig) {
		configuration.GameMode = "SingleHero"
	})

	// Assert
	assert.True(t, actual.HeroHireBan)
}

func TestWhenGameModeIsSingleHero_ForcesSingleHeroCounts(t *testing.T) {
	// Arrange
	expectedHeroSettings := config.HeroSettings{
		HeroCountMin:       1,
		HeroCountMax:       1,
		HeroCountIncrement: 1,
	}

	// Act
	actual := createGameRules(func(configuration *config.GeneratorConfig) {
		configuration.GameMode = "SingleHero"
		configuration.HeroSettings = config.HeroSettings{
			HeroCountMin:       gofakeit.Number(2, 5),
			HeroCountMax:       gofakeit.Number(6, 10),
			HeroCountIncrement: gofakeit.Number(2, 3),
		}
	})

	// Assert
	actualHeroSettings := config.HeroSettings{
		HeroCountMin:       actual.HeroCountMin,
		HeroCountMax:       actual.HeroCountMax,
		HeroCountIncrement: actual.HeroCountIncrement,
	}
	assert.Equal(t, expectedHeroSettings, actualHeroSettings)
}

func TestWhenGameModeIsClassic_PropagatesConfiguredHeroCounts(t *testing.T) {
	// Arrange
	expectedHeroSettings := config.HeroSettings{
		HeroCountMin:       gofakeit.Number(1, 5),
		HeroCountMax:       gofakeit.Number(5, 10),
		HeroCountIncrement: gofakeit.Number(1, 3),
	}

	// Act
	actual := createGameRules(func(configuration *config.GeneratorConfig) {
		configuration.HeroSettings = expectedHeroSettings
	})

	// Assert
	actualHeroSettings := config.HeroSettings{
		HeroCountMin:       actual.HeroCountMin,
		HeroCountMax:       actual.HeroCountMax,
		HeroCountIncrement: actual.HeroCountIncrement,
	}
	assert.Equal(t, expectedHeroSettings, actualHeroSettings)
}

// ── Experience modifiers ─────────────────────────────────────────────

func TestWhenFactionLawsExpPercentWithinRange_SetsPercentDividedByHundred(t *testing.T) {
	// Arrange
	expectedPercent := gofakeit.IntRange(25, 200)

	// Act
	actual := createGameRules(func(configuration *config.GeneratorConfig) {
		configuration.FactionLawsExpPercent = expectedPercent
	})

	// Assert
	assert.Equal(t, float64(expectedPercent)/100, actual.FactionLawsExpModifier)
}

func TestWhenFactionLawsExpPercentBelowMinimum_ClampsModifierToQuarter(t *testing.T) {
	// Arrange & Act
	actual := createGameRules(func(configuration *config.GeneratorConfig) {
		configuration.FactionLawsExpPercent = gofakeit.IntRange(-1000, 24)
	})

	// Assert
	assert.Equal(t, 0.25, actual.FactionLawsExpModifier)
}

func TestWhenFactionLawsExpPercentAboveMaximum_ClampsModifierToTwo(t *testing.T) {
	// Arrange & Act
	actual := createGameRules(func(configuration *config.GeneratorConfig) {
		configuration.FactionLawsExpPercent = gofakeit.IntRange(201, 9999)
	})

	// Assert
	assert.Equal(t, 2.0, actual.FactionLawsExpModifier)
}

func TestWhenAstrologyExpPercentWithinRange_SetsPercentDividedByHundred(t *testing.T) {
	// Arrange
	expectedPercent := gofakeit.IntRange(25, 200)

	// Act
	actual := createGameRules(func(configuration *config.GeneratorConfig) {
		configuration.AstrologyExpPercent = expectedPercent
	})

	// Assert
	assert.Equal(t, float64(expectedPercent)/100, actual.AstrologyExpModifier)
}

func TestWhenAstrologyExpPercentBelowMinimum_ClampsModifierToQuarter(t *testing.T) {
	// Arrange & Act
	actual := createGameRules(func(configuration *config.GeneratorConfig) {
		configuration.AstrologyExpPercent = gofakeit.IntRange(-1000, 24)
	})

	// Assert
	assert.Equal(t, 0.25, actual.AstrologyExpModifier)
}

func TestWhenAstrologyExpPercentAboveMaximum_ClampsModifierToTwo(t *testing.T) {
	// Arrange & Act
	actual := createGameRules(func(configuration *config.GeneratorConfig) {
		configuration.AstrologyExpPercent = gofakeit.IntRange(201, 9999)
	})

	// Assert
	assert.Equal(t, 2.0, actual.AstrologyExpModifier)
}

// ── Win conditions from victory-condition presets ────────────────────

func TestWhenVictoryConditionThree_EnablesLostStartCity(t *testing.T) {
	// Arrange & Act
	actual := createGameRules(func(configuration *config.GeneratorConfig) {
		configuration.GameEndConditions = &config.GameEndConditions{
			VictoryCondition: "win_condition_3",
			LostStartCityDay: 5,
			CityHoldDays:     6,
		}
	})

	// Assert
	assert.True(t, actual.WinConditions.LostStartCity)
}

func TestWhenVictoryConditionFive_EnablesCityHold(t *testing.T) {
	// Arrange & Act
	actual := createGameRules(func(configuration *config.GeneratorConfig) {
		configuration.GameEndConditions = &config.GameEndConditions{
			VictoryCondition: "win_condition_5",
			LostStartCityDay: 3,
			CityHoldDays:     6,
		}
	})

	// Assert
	assert.True(t, actual.WinConditions.CityHold)
}

func TestWhenVictoryConditionFour_EnablesGladiatorArena(t *testing.T) {
	// Arrange & Act
	actual := createGameRules(func(configuration *config.GeneratorConfig) {
		configuration.GameEndConditions = &config.GameEndConditions{
			VictoryCondition: "win_condition_4",
			LostStartCityDay: 3,
			CityHoldDays:     6,
		}
	})

	// Assert
	assert.True(t, actual.WinConditions.GladiatorArena)
}

func TestWhenVictoryConditionSix_EnablesTournament(t *testing.T) {
	// Arrange & Act
	actual := createGameRules(func(configuration *config.GeneratorConfig) {
		configuration.GameEndConditions = &config.GameEndConditions{
			VictoryCondition: "win_condition_6",
			LostStartCityDay: 3,
			CityHoldDays:     6,
		}
	})

	// Assert
	assert.True(t, actual.WinConditions.Tournament)
}

func TestWhenCityHoldConditionEnabled_EnablesCityHold(t *testing.T) {
	// Arrange & Act
	actual := createGameRules(func(configuration *config.GeneratorConfig) {
		configuration.GameEndConditions = &config.GameEndConditions{
			VictoryCondition: "win_condition_1",
			CityHold:         true,
			CityHoldDays:     gofakeit.Number(1, 10),
			LostStartCityDay: 3,
		}
	})

	// Assert
	assert.True(t, actual.WinConditions.CityHold)
}

func TestWhenLostStartCityDayAboveMaximum_ClampsToThirty(t *testing.T) {
	// Arrange & Act
	actual := createGameRules(func(configuration *config.GeneratorConfig) {
		configuration.GameEndConditions = &config.GameEndConditions{
			VictoryCondition: "win_condition_3",
			LostStartCityDay: gofakeit.Number(31, 999),
			CityHoldDays:     6,
		}
	})

	// Assert
	assert.Equal(t, 30, actual.WinConditions.LostStartCityDay)
}

func TestWhenCityHoldDaysAboveMaximum_ClampsToThirty(t *testing.T) {
	// Arrange & Act
	actual := createGameRules(func(configuration *config.GeneratorConfig) {
		configuration.GameEndConditions = &config.GameEndConditions{
			VictoryCondition: "win_condition_5",
			LostStartCityDay: 3,
			CityHoldDays:     gofakeit.Number(31, 999),
		}
	})

	// Assert
	assert.Equal(t, 30, actual.WinConditions.CityHoldDays)
}

// ── Gladiator arena ──────────────────────────────────────────────────

func TestWhenGladiatorArenaEnabled_EnablesGladiatorArenaWinCondition(t *testing.T) {
	// Arrange & Act
	actual := createGameRules(func(configuration *config.GeneratorConfig) {
		configuration.GladiatorArenaRules = &config.GladiatorArenaRules{
			Enabled:        true,
			DaysDelayStart: 30,
			CountDay:       3,
		}
	})

	// Assert
	assert.True(t, actual.WinConditions.GladiatorArena)
}

func TestWhenGladiatorArenaEnabled_SetsConfiguredDaysDelayStart(t *testing.T) {
	// Arrange
	expectedDaysDelayStart := gofakeit.Number(1, 60)

	// Act
	actual := createGameRules(func(configuration *config.GeneratorConfig) {
		configuration.GladiatorArenaRules = &config.GladiatorArenaRules{
			Enabled:        true,
			DaysDelayStart: expectedDaysDelayStart,
			CountDay:       3,
		}
	})

	// Assert
	assert.Equal(t, expectedDaysDelayStart, actual.WinConditions.GladiatorArenaDaysDelayStart)
}

func TestWhenGladiatorArenaEnabled_SetsConfiguredCountDay(t *testing.T) {
	// Arrange
	expectedCountDay := gofakeit.Number(1, 30)

	// Act
	actual := createGameRules(func(configuration *config.GeneratorConfig) {
		configuration.GladiatorArenaRules = &config.GladiatorArenaRules{
			Enabled:        true,
			DaysDelayStart: 30,
			CountDay:       expectedCountDay,
		}
	})

	// Assert
	assert.Equal(t, expectedCountDay, actual.WinConditions.GladiatorArenaCountDay)
}

func TestWhenGladiatorArenaEnabled_SetsChampionSelectRuleToStartHero(t *testing.T) {
	// Arrange & Act
	actual := createGameRules(func(configuration *config.GeneratorConfig) {
		configuration.GladiatorArenaRules = &config.GladiatorArenaRules{
			Enabled:        true,
			DaysDelayStart: 30,
			CountDay:       3,
		}
	})

	// Assert
	assert.Equal(t, "StartHero", actual.WinConditions.ChampionSelectRule)
}

func TestWhenGladiatorArenaEnabled_ForcesLostStartHero(t *testing.T) {
	// Arrange & Act
	actual := createGameRules(func(configuration *config.GeneratorConfig) {
		configuration.GladiatorArenaRules = &config.GladiatorArenaRules{
			Enabled:        true,
			DaysDelayStart: 30,
			CountDay:       3,
		}
	})

	// Assert
	assert.True(t, actual.WinConditions.LostStartHero)
}

func TestWhenGameModeIsSingleHero_ForcesLostStartHero(t *testing.T) {
	// Arrange & Act
	actual := createGameRules(func(configuration *config.GeneratorConfig) {
		configuration.GameMode = "SingleHero"
	})

	// Assert
	assert.True(t, actual.WinConditions.LostStartHero)
}

func TestWhenGladiatorDaysDelayStartAboveMaximum_ClampsToSixty(t *testing.T) {
	// Arrange & Act
	actual := createGameRules(func(configuration *config.GeneratorConfig) {
		configuration.GladiatorArenaRules = &config.GladiatorArenaRules{
			Enabled:        true,
			DaysDelayStart: gofakeit.Number(61, 999),
			CountDay:       3,
		}
	})

	// Assert
	assert.Equal(t, 60, actual.WinConditions.GladiatorArenaDaysDelayStart)
}

// ── Tournament ───────────────────────────────────────────────────────

func TestWhenTournamentEnabled_EnablesTournamentWinCondition(t *testing.T) {
	// Arrange & Act
	actual := createGameRules(func(configuration *config.GeneratorConfig) {
		configuration.TournamentRules = &config.TournamentRules{
			Enabled:            true,
			FirstTournamentDay: 10,
			Interval:           5,
			PointsToWin:        2,
		}
	})

	// Assert
	assert.True(t, actual.WinConditions.Tournament)
}

func TestWhenTournamentEnabled_EnablesTournamentSaveArmy(t *testing.T) {
	// Arrange & Act
	actual := createGameRules(func(configuration *config.GeneratorConfig) {
		configuration.TournamentRules = &config.TournamentRules{
			Enabled:            true,
			FirstTournamentDay: 10,
			Interval:           5,
			PointsToWin:        2,
		}
	})

	// Assert
	assert.True(t, actual.WinConditions.TournamentSaveArmy)
}

func TestWhenTournamentEnabled_SetsConfiguredPointsToWin(t *testing.T) {
	// Arrange
	expectedPointsToWin := gofakeit.Number(1, 10)

	// Act
	actual := createGameRules(func(configuration *config.GeneratorConfig) {
		configuration.TournamentRules = &config.TournamentRules{
			Enabled:            true,
			FirstTournamentDay: 10,
			Interval:           5,
			PointsToWin:        expectedPointsToWin,
		}
	})

	// Assert
	assert.Equal(t, expectedPointsToWin, actual.WinConditions.TournamentPointsToWin)
}

func TestWhenTournamentEnabled_CreatesAnnounceDayPerRound(t *testing.T) {
	// Arrange
	pointsToWin := gofakeit.Number(1, 5)
	expectedRoundCount := pointsToWin*2 - 1

	// Act
	actual := createGameRules(func(configuration *config.GeneratorConfig) {
		configuration.TournamentRules = &config.TournamentRules{
			Enabled:            true,
			FirstTournamentDay: 10,
			Interval:           5,
			PointsToWin:        pointsToWin,
		}
	})

	// Assert
	assert.Len(t, actual.WinConditions.TournamentAnnounceDays, expectedRoundCount)
}

func TestWhenTournamentEnabled_CreatesBattleDayOffsetPerRound(t *testing.T) {
	// Arrange
	pointsToWin := gofakeit.Number(1, 5)
	expectedRoundCount := pointsToWin*2 - 1

	// Act
	actual := createGameRules(func(configuration *config.GeneratorConfig) {
		configuration.TournamentRules = &config.TournamentRules{
			Enabled:            true,
			FirstTournamentDay: 10,
			Interval:           5,
			PointsToWin:        pointsToWin,
		}
	})

	// Assert
	assert.Len(t, actual.WinConditions.TournamentDays, expectedRoundCount)
}

func TestWhenTournamentScheduleComputed_SetsExpectedAnnounceDays(t *testing.T) {
	// Arrange: first battle on day 10, then every 5 days; announcements follow
	// the day after each previous battle: rounds land on days 10, 15, 20.
	// Act
	actual := createGameRules(func(configuration *config.GeneratorConfig) {
		configuration.TournamentRules = &config.TournamentRules{
			Enabled:            true,
			FirstTournamentDay: 10,
			Interval:           5,
			PointsToWin:        2,
		}
	})

	// Assert
	assert.Equal(t, []int{1, 11, 16}, actual.WinConditions.TournamentAnnounceDays)
}

func TestWhenTournamentScheduleComputed_SetsExpectedBattleDayOffsets(t *testing.T) {
	// Arrange: offsets are relative to each announcement day (first-day - 1,
	// then interval - 1 per subsequent round).
	// Act
	actual := createGameRules(func(configuration *config.GeneratorConfig) {
		configuration.TournamentRules = &config.TournamentRules{
			Enabled:            true,
			FirstTournamentDay: 10,
			Interval:           5,
			PointsToWin:        2,
		}
	})

	// Assert
	assert.Equal(t, []int{9, 4, 4}, actual.WinConditions.TournamentDays)
}

func TestWhenTournamentPointsToWinAboveMaximum_ClampsToTen(t *testing.T) {
	// Arrange & Act
	actual := createGameRules(func(configuration *config.GeneratorConfig) {
		configuration.TournamentRules = &config.TournamentRules{
			Enabled:            true,
			FirstTournamentDay: 10,
			Interval:           5,
			PointsToWin:        gofakeit.Number(11, 999),
		}
	})

	// Assert
	assert.Equal(t, 10, actual.WinConditions.TournamentPointsToWin)
}

// ── Bonuses ──────────────────────────────────────────────────────────

func TestWhenTownPortalFreeBonusConfigured_ExpandsToSpellAndFreeCostBonuses(t *testing.T) {
	// Arrange
	entry := config.BonusEntry{PresetType: config.BonusTownPortalFree, ReceiverFilter: "start_hero"}

	// Act
	actual := bonusesFor(entry)

	// Assert
	assert.Equal(t, entities.BonusList{
		{
			SID:            "add_bonus_hero_spell",
			ReceiverSide:   -1,
			ReceiverFilter: "start_hero",
			Parameters:     []string{"neutral_magic_town_portal"},
		},
		{
			SID:            "add_bonus_hero_stat",
			ReceiverSide:   -1,
			ReceiverFilter: "start_hero",
			Parameters:     []string{"magicCostSidSet", "neutral_magic_town_portal", "-999", "0"},
		},
	}, actual)
}

func TestWhenFreeSpellBonusConfigured_ExpandsToSpellAndCostOverrideBonuses(t *testing.T) {
	// Arrange
	entry := config.BonusEntry{
		PresetType:     config.BonusSpell,
		ReceiverFilter: "all_heroes",
		Param:          "magic_fireball",
		Param2:         "1",
	}

	// Act
	actual := bonusesFor(entry)

	// Assert
	assert.Equal(t, entities.BonusList{
		{
			SID:            "add_bonus_hero_spell",
			ReceiverSide:   -1,
			ReceiverFilter: "all_heroes",
			Parameters:     []string{"magic_fireball"},
		},
		{
			SID:            "add_bonus_hero_stat",
			ReceiverSide:   -1,
			ReceiverFilter: "all_heroes",
			Parameters:     []string{"magicCostSidSet", "magic_fireball", "-999", "0"},
		},
	}, actual)
}

func TestWhenPaidSpellBonusConfigured_ProducesOnlySpellBonus(t *testing.T) {
	// Arrange
	entry := config.BonusEntry{
		PresetType:     config.BonusSpell,
		ReceiverFilter: "all_heroes",
		Param:          "magic_fireball",
		Param2:         "0",
	}

	// Act
	actual := bonusesFor(entry)

	// Assert
	assert.Equal(t, entities.BonusList{
		{
			SID:            "add_bonus_hero_spell",
			ReceiverSide:   -1,
			ReceiverFilter: "all_heroes",
			Parameters:     []string{"magic_fireball"},
		},
	}, actual)
}

func TestWhenSingleBonusPresetConfigured_ProducesExpectedBonus(t *testing.T) {
	cases := []struct {
		name               string
		preset             config.BonusPresetType
		param              string
		expectedSID        string
		expectedParameters []string
	}{
		{
			"WhenUnitMultiplierConfigured_ProducesHeroUnitMultiplierBonus",
			config.BonusUnitMultiplier,
			"2",
			"add_bonus_hero_unit_multipler",
			[]string{"2"},
		},
		{
			"WhenMovementBonusConfigured_ProducesHeroStatBonus",
			config.BonusMovementBonus,
			"300",
			"add_bonus_hero_stat",
			[]string{"movementBonus", "300"},
		},
		{
			"WhenStartingItemConfigured_ProducesHeroItemBonus",
			config.BonusStartingItem,
			"some_item",
			"add_bonus_hero_item",
			[]string{"some_item"},
		},
		{
			"WhenStartingGoldConfigured_ProducesResourceBonus",
			config.BonusStartingGold,
			"1000",
			"add_bonus_res",
			[]string{"gold", "1000"},
		},
		{
			"WhenStartingGemsConfigured_ProducesResourceBonus",
			config.BonusStartingGems,
			"5",
			"add_bonus_res",
			[]string{"gemstones", "5"},
		},
		{
			"WhenStartingCrystalsConfigured_ProducesResourceBonus",
			config.BonusStartingCrystals,
			"5",
			"add_bonus_res",
			[]string{"crystals", "5"},
		},
		{
			"WhenStartingMercuryConfigured_ProducesResourceBonus",
			config.BonusStartingMercury,
			"5",
			"add_bonus_res",
			[]string{"mercury", "5"},
		},
		{
			"WhenStartingWoodConfigured_ProducesResourceBonus",
			config.BonusStartingWood,
			"10",
			"add_bonus_res",
			[]string{"wood", "10"},
		},
		{
			"WhenStartingOreConfigured_ProducesResourceBonus",
			config.BonusStartingOre,
			"10",
			"add_bonus_res",
			[]string{"ore", "10"},
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			// Arrange
			entry := config.BonusEntry{
				PresetType:     testCase.preset,
				ReceiverFilter: "start_hero",
				Param:          testCase.param,
			}

			// Act
			actual := bonusesFor(entry)

			// Assert
			assert.Equal(t, entities.BonusList{
				{
					SID:            testCase.expectedSID,
					ReceiverSide:   -1,
					ReceiverFilter: "start_hero",
					Parameters:     testCase.expectedParameters,
				},
			}, actual)
		})
	}
}

func TestWhenMultipleBonusEntriesConfigured_ConcatenatesExpansions(t *testing.T) {
	// Arrange
	entries := []config.BonusEntry{
		{PresetType: config.BonusStartingGold, ReceiverFilter: "start_hero", Param: "1000"},
		{PresetType: config.BonusTownPortalFree, ReceiverFilter: "start_hero"},
	}

	// Act
	actual := bonusesFor(entries...)

	// Assert
	assert.Len(t, actual, 3) // 1 resource + 2 town-portal bonuses.
}

func TestWhenNoBonusEntriesConfigured_ProducesEmptyBonusList(t *testing.T) {
	// Arrange & Act
	actual := bonusesFor()

	// Assert
	assert.Empty(t, actual)
}

// Functional-equivalence check against a real game template: Blitz grants the
// free Town Portal, whose raw bonuses must match our expansion exactly.
func TestWhenTownPortalFreeExpansionComparedToBlitzTemplate_MatchesExactly(t *testing.T) {
	// Arrange
	blitz := loadExampleTemplate(t, "Blitz.rmg.json")

	// Act
	actual := bonusesFor(config.BonusEntry{PresetType: config.BonusTownPortalFree, ReceiverFilter: "start_hero"})

	// Assert
	assert.Equal(t, blitz.GameRules.Bonuses, actual)
}
