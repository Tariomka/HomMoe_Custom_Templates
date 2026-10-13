package gameRulesProvider_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenSeveralOverridesAreConfigured_EmitsOnePerEntryInOrder(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.ValueOverrides = []config.ValueOverrideEntry{
		{Sid: "watchtower", GuardValue: 25000},
		{Sid: "gold_mine", GuardValue: 12000},
		{Sid: "watchtower", GuardValue: 9000},
	}

	// Act
	actual := providers.NewGameRulesProvider().CreateValueOverrides(*configuration)

	// Assert
	assert.Equal(t, []template_model.ValueOverride{
		{SID: "watchtower", Variant: -1, GuardValue: 25000},
		{SID: "gold_mine", Variant: -1, GuardValue: 12000},
		{SID: "watchtower", Variant: -1, GuardValue: 9000},
	}, actual)
}

func TestWhenNoOverridesAreConfigured_ReturnsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()

	// Act
	actual := providers.NewGameRulesProvider().CreateValueOverrides(*configuration)

	// Assert
	assert.Nil(t, actual)
}

// Functional-equivalence check: the first override of the real Blitz template
// fed back as an entry must be reproduced (with the generator's all-variants
// marker -1).
func TestWhenBlitzOverrideIsConfigured_ReproducesBlitzSidAndGuardValue(t *testing.T) {
	t.Parallel()
	// Arrange
	blitz := loadExampleTemplate(t, "Blitz.rmg.json")
	require.NotEmpty(t, blitz.ValueOverrides)
	blitzOverride := blitz.ValueOverrides[0]
	configuration := config.NewGeneratorConfig()
	configuration.ValueOverrides = []config.ValueOverrideEntry{
		{Sid: blitzOverride.SID, GuardValue: blitzOverride.GuardValue},
	}

	// Act
	actual := providers.NewGameRulesProvider().CreateValueOverrides(*configuration)

	// Assert
	assert.Equal(t, []template_model.ValueOverride{
		{SID: blitzOverride.SID, Variant: -1, GuardValue: blitzOverride.GuardValue},
	}, actual)
}
