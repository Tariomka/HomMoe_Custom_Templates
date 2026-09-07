package generatorConfig_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/stretchr/testify/assert"
)

func TestWhenConfigIsCreated_ReturnsDocumentedDefaults(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := &config.GeneratorConfig{
		TemplateName: "Custom Template",
		GameMode:     "Classic",
		PlayerCount:  2,
		MapSize:      160,
		HeroSettings: config.HeroSettings{
			HeroCountMin:       4,
			HeroCountMax:       8,
			HeroCountIncrement: 1,
		},
		SpawnRemoteFootholds:  true,
		RemoteFootholdCount:   1,
		GenerateRoads:         true,
		MaxPortalConnections:  32,
		Topology:              config.TopologyCircles,
		FactionLawsExpPercent: 100,
		AstrologyExpPercent:   100,
		ZoneConfiguration: config.ZoneConfig{
			PlayerZoneCastles:           1,
			NeutralZoneCastles:          1,
			AbandonedOutpostCount:       1,
			ResourceDensityPercent:      100,
			StructureDensityPercent:     100,
			NeutralStackStrengthPercent: 100,
			BorderGuardStrengthPercent:  100,
			PlayerZoneSize:              1.0,
			NeutralZoneSize:             1.0,
			GuardRandomization:          0.05,
			HubZoneSize:                 1.0,
			Advanced:                    config.AdvancedSettings{},
		},
		GameEndConditions: &config.GameEndConditions{
			VictoryCondition: "win_condition_1",
			LostStartCityDay: 3,
			CityHoldDays:     6,
		},
		GladiatorArenaRules: &config.GladiatorArenaRules{DaysDelayStart: 30, CountDay: 3},
		TournamentRules: &config.TournamentRules{
			FirstTournamentDay: 14,
			Interval:           7,
			PointsToWin:        2,
			SaveArmy:           true,
		},
	}

	// Act
	actual := config.NewGeneratorConfig()

	// Assert
	assert.Equal(t, expected, actual)
}
