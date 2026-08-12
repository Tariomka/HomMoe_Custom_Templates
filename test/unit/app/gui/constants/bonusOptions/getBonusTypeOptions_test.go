package bonusOptions_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenBonusTypeOptionsAreRequested_EveryOptionIsLabelled(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	options := constants.GetBonusTypeOptions()

	// Assert
	require.NotEmpty(t, options)
	for _, option := range options {
		require.NotEmpty(t, option.Label, "preset %v has no label", option.PresetType)
	}
}

func TestWhenBonusTypeOptionsAreRequested_EveryPresetTypeAppearsOnce(t *testing.T) {
	t.Parallel()
	// Arrange
	options := constants.GetBonusTypeOptions()

	// Act
	seen := map[config.BonusPresetType]bool{}
	for _, option := range options {
		seen[option.PresetType] = true
	}

	// Assert
	assert.Len(t, seen, len(options))
}
