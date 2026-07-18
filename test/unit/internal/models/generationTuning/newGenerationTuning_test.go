package generationTuning_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/utils"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenAllOptionalSpawnsAreEnabled_BuildsTuningFromConfiguredPercentages(t *testing.T) {
	t.Parallel()
	// Arrange
	resourceDensity := gofakeit.Number(25, 300)
	structureDensity := gofakeit.Number(25, 300)
	neutralStrength := gofakeit.Number(25, 300)
	borderStrength := gofakeit.Number(25, 300)
	footholdCount := gofakeit.Number(1, 4)
	outpostCount := gofakeit.Number(1, 4)
	ownedCastles := gofakeit.Number(0, 4)
	mapSize := gofakeit.Number(96, 224)
	totalZoneCount := gofakeit.Number(3, 20)

	configuration := config.NewGeneratorConfig()
	configuration.MapSize = mapSize
	configuration.SpawnRemoteFootholds = true
	configuration.RemoteFootholdCount = footholdCount
	configuration.ZoneConfiguration.SpawnAbandonedOutposts = true
	configuration.ZoneConfiguration.AbandonedOutpostCount = outpostCount
	configuration.ZoneConfiguration.ResourceDensityPercent = resourceDensity
	configuration.ZoneConfiguration.StructureDensityPercent = structureDensity
	configuration.ZoneConfiguration.NeutralStackStrengthPercent = neutralStrength
	configuration.ZoneConfiguration.BorderGuardStrengthPercent = borderStrength
	configuration.ZoneConfiguration.PlayerOwnedCastles = ownedCastles

	expected := models.GenerationTuning{
		ContentScale:                   utils.ComputeContentScale(mapSize, totalZoneCount),
		ResourceDensityMultiplier:      float64(resourceDensity) / 200.0,
		StructureDensityMultiplier:     float64(structureDensity) / 100.0,
		NeutralStackStrengthMultiplier: float64(neutralStrength) / 100.0,
		BorderGuardStrengthMultiplier:  float64(borderStrength) / 100.0,
		GuardRandomization:             0.05, // advanced settings disabled -> default
		RemoteFootholdCount:            footholdCount,
		AbandonedOutpostCount:          outpostCount,
		PlayerOwnedCastles:             ownedCastles,
	}

	// Act
	tuning := models.NewGenerationTuning(configuration, totalZoneCount)

	// Assert
	assert.Equal(t, expected, tuning)
}

func TestWhenRemoteFootholdsAreDisabled_ZeroesFootholdCount(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.SpawnRemoteFootholds = false
	configuration.RemoteFootholdCount = gofakeit.Number(1, 4)

	// Act
	tuning := models.NewGenerationTuning(configuration, 5)

	// Assert
	assert.Equal(t, 0, tuning.RemoteFootholdCount)
}

func TestWhenAbandonedOutpostsAreDisabled_ZeroesOutpostCount(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.ZoneConfiguration.SpawnAbandonedOutposts = false
	configuration.ZoneConfiguration.AbandonedOutpostCount = gofakeit.Number(1, 4)

	// Act
	tuning := models.NewGenerationTuning(configuration, 5)

	// Assert
	assert.Equal(t, 0, tuning.AbandonedOutpostCount)
}

func TestWhenAdvancedSettingsAreEnabled_UsesConfiguredGuardRandomization(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.ZoneConfiguration.Advanced.Enabled = true
	configuration.ZoneConfiguration.GuardRandomization = 0.2

	// Act
	tuning := models.NewGenerationTuning(configuration, 5)

	// Assert
	assert.InDelta(t, 0.2, tuning.GuardRandomization, test_helpers.Delta)
}
