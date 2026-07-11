package generatorConfig_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenTournamentRulesAreSet_ReturnsTheirCopy(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := config.TournamentRules{
		Enabled:            true,
		FirstTournamentDay: gofakeit.Number(1, 30),
		Interval:           gofakeit.Number(1, 14),
		PointsToWin:        gofakeit.Number(1, 5),
		SaveArmy:           gofakeit.Bool(),
	}
	configuration := config.NewGeneratorConfig()
	configuration.TournamentRules = &expected

	// Act
	actual := configuration.GetTournamentRules()

	// Assert
	assert.Equal(t, expected, actual)
}

func TestWhenTournamentRulesAreNil_ReturnsZeroRules(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.TournamentRules = nil

	// Act
	actual := configuration.GetTournamentRules()

	// Assert
	assert.Equal(t, config.TournamentRules{}, actual)
}
