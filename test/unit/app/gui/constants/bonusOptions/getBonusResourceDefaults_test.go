package bonusOptions_test

import (
	"strconv"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenResourceDefaultsAreRequested_CoversEveryStartingResourcePreset(t *testing.T) {
	t.Parallel()
	// Arrange
	resourcePresets := []config.BonusPresetType{
		config.BonusStartingGold,
		config.BonusStartingGems,
		config.BonusStartingCrystals,
		config.BonusStartingMercury,
		config.BonusStartingWood,
		config.BonusStartingOre,
	}

	// Act
	defaults := constants.GetBonusResourceDefaults()

	// Assert
	for _, preset := range resourcePresets {
		require.Contains(t, defaults, preset)
	}
	assert.Len(t, defaults, len(resourcePresets))
}

func TestWhenResourceDefaultsAreRequested_EveryValueIsAPositiveInteger(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	defaults := constants.GetBonusResourceDefaults()

	// Assert
	for preset, value := range defaults {
		amount, err := strconv.Atoi(value)
		require.NoError(t, err, "preset %v has a non-numeric default %q", preset, value)
		require.Positive(t, amount, "preset %v has a non-positive default", preset)
	}
}
