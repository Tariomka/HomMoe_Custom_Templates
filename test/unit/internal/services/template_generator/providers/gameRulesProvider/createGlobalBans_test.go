package gameRulesProvider_test

import (
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenItemsAndMagicsBanned_ReturnsBothLists(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.BannedItems = "voodoosh_doll_artifact\nflag_of_truce_artifact"
	configuration.BannedMagics = "magic_armageddon"

	// Act
	actual := providers.NewGameRulesProvider().CreateGlobalBans(*configuration)

	// Assert
	assert.Equal(t, &entities.GlobalBans{
		Items:  []string{"voodoosh_doll_artifact", "flag_of_truce_artifact"},
		Magics: []string{"magic_armageddon"},
	}, actual)
}

func TestWhenNothingBanned_ReturnsNil(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()

	// Act
	actual := providers.NewGameRulesProvider().CreateGlobalBans(*configuration)

	// Assert
	assert.Nil(t, actual)
}

func TestWhenOnlyMagicsBanned_ReturnsBansWithNilItems(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.BannedMagics = "magic_armageddon"

	// Act
	actual := providers.NewGameRulesProvider().CreateGlobalBans(*configuration)

	// Assert
	assert.Equal(t, &entities.GlobalBans{Items: nil, Magics: []string{"magic_armageddon"}}, actual)
}

func TestWhenBannedLinesContainWhitespace_TrimsAndSkipsBlankLines(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.BannedItems = "  voodoosh_doll_artifact  \n\n\t\nflag_of_truce_artifact"

	// Act
	actual := providers.NewGameRulesProvider().CreateGlobalBans(*configuration)

	// Assert
	assert.Equal(t, []string{"voodoosh_doll_artifact", "flag_of_truce_artifact"}, actual.Items)
}

// Functional-equivalence check: feeding the real Blitz item bans back through
// the parser must reproduce them exactly.
func TestWhenBlitzItemBansParsed_ReproducesBlitzItems(t *testing.T) {
	// Arrange
	blitz := loadExampleTemplate(t, "Blitz.rmg.json")
	require.NotNil(t, blitz.GlobalBans)
	require.NotEmpty(t, blitz.GlobalBans.Items)
	configuration := config.NewGeneratorConfig()
	configuration.BannedItems = strings.Join(blitz.GlobalBans.Items, "\n")

	// Act
	actual := providers.NewGameRulesProvider().CreateGlobalBans(*configuration)

	// Assert
	assert.Equal(t, blitz.GlobalBans.Items, actual.Items)
}
