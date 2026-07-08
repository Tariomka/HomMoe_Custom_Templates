package generatorConfig_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenGladiatorArenaRulesAreSet_ReturnsTheirCopy(t *testing.T) {
	// Arrange
	expected := config.GladiatorArenaRules{
		Enabled:        true,
		DaysDelayStart: gofakeit.Number(1, 60),
		CountDay:       gofakeit.Number(1, 7),
	}
	configuration := config.NewGeneratorConfig()
	configuration.GladiatorArenaRules = &expected

	// Act
	actual := configuration.GetGladiatorArenaRules()

	// Assert
	assert.Equal(t, expected, actual)
}

func TestWhenGladiatorArenaRulesAreNil_ReturnsZeroRules(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.GladiatorArenaRules = nil

	// Act
	actual := configuration.GetGladiatorArenaRules()

	// Assert
	assert.Equal(t, config.GladiatorArenaRules{}, actual)
}
