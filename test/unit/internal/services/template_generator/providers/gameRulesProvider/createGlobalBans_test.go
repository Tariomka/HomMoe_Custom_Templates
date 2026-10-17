package gameRulesProvider_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenItemsAndMagicsBanned_ReturnsBothLists(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.BannedItems = bans("voodoosh_doll_artifact", "flag_of_truce_artifact")
	configuration.BannedMagics = bans("magic_armageddon")

	// Act
	actual := providers.NewGameRulesProvider().CreateGlobalBans(*configuration)

	// Assert
	assert.Equal(t, &template_model.GlobalBans{
		Items:  []string{"voodoosh_doll_artifact", "flag_of_truce_artifact"},
		Magics: []string{"magic_armageddon"},
	}, actual)
}

func TestWhenNothingBanned_ReturnsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()

	// Act
	actual := providers.NewGameRulesProvider().CreateGlobalBans(*configuration)

	// Assert
	assert.Nil(t, actual)
}

func TestWhenOnlyMagicsBanned_ReturnsBansWithNilItems(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.BannedMagics = bans("magic_armageddon")

	// Act
	actual := providers.NewGameRulesProvider().CreateGlobalBans(*configuration)

	// Assert
	assert.Equal(t, &template_model.GlobalBans{Items: nil, Magics: []string{"magic_armageddon"}}, actual)
}

func TestWhenOnlyItemsBanned_ReturnsBansWithNilMagics(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.BannedItems = bans("voodoosh_doll_artifact")

	// Act
	actual := providers.NewGameRulesProvider().CreateGlobalBans(*configuration)

	// Assert
	assert.Equal(t, &template_model.GlobalBans{Items: []string{"voodoosh_doll_artifact"}, Magics: nil}, actual)
}

func TestWhenBansRepeat_KeepsOrderAndDuplicates(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.BannedItems = bans("flag_of_truce_artifact", "voodoosh_doll_artifact", "flag_of_truce_artifact")

	// Act
	actual := providers.NewGameRulesProvider().CreateGlobalBans(*configuration)

	// Assert
	assert.Equal(t,
		[]string{"flag_of_truce_artifact", "voodoosh_doll_artifact", "flag_of_truce_artifact"},
		actual.Items)
}

// Functional-equivalence check: the real Blitz item bans fed back as entries
// must be reproduced exactly.
func TestWhenBlitzItemBansAreConfigured_ReproducesBlitzItems(t *testing.T) {
	t.Parallel()
	// Arrange
	blitz := loadExampleTemplate(t, "Blitz.rmg.json")
	require.NotNil(t, blitz.GlobalBans)
	require.NotEmpty(t, blitz.GlobalBans.Items)
	configuration := config.NewGeneratorConfig()
	configuration.BannedItems = bans(blitz.GlobalBans.Items...)

	// Act
	actual := providers.NewGameRulesProvider().CreateGlobalBans(*configuration)

	// Assert
	assert.Equal(t, blitz.GlobalBans.Items, actual.Items)
}

func bans(sids ...string) []config.BanEntry {
	entries := make([]config.BanEntry, 0, len(sids))
	for _, sid := range sids {
		entries = append(entries, config.BanEntry{Sid: sid})
	}
	return entries
}
