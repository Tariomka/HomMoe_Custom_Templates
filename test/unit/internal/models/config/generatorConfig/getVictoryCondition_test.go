package generatorConfig_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/stretchr/testify/assert"
)

func TestWhenGameEndConditionsExist_ReturnsConfiguredCondition(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.GameEndConditions.VictoryCondition = "win_condition_5"

	// Act
	actual := configuration.GetVictoryCondition()

	// Assert
	assert.Equal(t, "win_condition_5", actual)
}

func TestWhenGameEndConditionsAreNil_ReturnsStandardCondition(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.GameEndConditions = nil

	// Act
	actual := configuration.GetVictoryCondition()

	// Assert
	assert.Equal(t, "win_condition_1", actual)
}
