package gameRulesProvider_test

import (
	"fmt"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenTextMixesValidAndJunkLines_ParsesOnlyValidLines(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.ValueOverridesText = "watchtower=25000\n\n  =5 \nbad_line\ngold_mine = 12000 \nnonnum=abc"

	// Act
	actual, _ := providers.NewGameRulesProvider().CreateValueOverrides(*configuration)

	// Assert
	assert.Equal(t, []entities.ValueOverride{
		{SID: "watchtower", Variant: -1, GuardValue: 25000},
		{SID: "gold_mine", Variant: -1, GuardValue: 12000},
	}, actual)
}

func TestWhenTextIsEmpty_ReturnsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()

	// Act
	actual, _ := providers.NewGameRulesProvider().CreateValueOverrides(*configuration)

	// Assert
	assert.Nil(t, actual)
}

func TestWhenTextIsEmpty_ReturnsNoWarnings(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()

	// Act
	_, warnings := providers.NewGameRulesProvider().CreateValueOverrides(*configuration)

	// Assert
	assert.Empty(t, warnings)
}

func TestWhenTextIsOnlyBlankLines_ReturnsNoWarnings(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.ValueOverridesText = "\n   \n\t\n"

	// Act
	_, warnings := providers.NewGameRulesProvider().CreateValueOverrides(*configuration)

	// Assert
	assert.Empty(t, warnings)
}

func TestWhenLineIsRejected_ReturnsWarningNamingTheLine(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name            string
		overridesText   string
		expectedWarning string
	}{
		{"LineHasNoSeparator_ReturnsWarning", "bad_line", "line 1: 'bad_line' is not sid=value"},
		{"LineStartsWithSeparator_ReturnsWarning", "=5", "line 1: '=5' is not sid=value"},
		{"SeparatorFollowsOnlyWhitespace_ReturnsWarningForTheTrimmedLine", "  = 5",
			"line 1: '= 5' is not sid=value"},
		{"ValueIsNotANumber_ReturnsWarning", "nonnum=abc", "line 1: 'nonnum=abc' has a non-numeric value"},
		{
			"RejectedLineFollowsBlankLines_ReturnsWarningWithTheSourceLineNumber",
			"watchtower=25000\n\nbad_line",
			"line 3: 'bad_line' is not sid=value",
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			configuration := config.NewGeneratorConfig()
			configuration.ValueOverridesText = testCase.overridesText

			// Act
			_, warnings := providers.NewGameRulesProvider().CreateValueOverrides(*configuration)

			// Assert
			assert.Equal(t, []string{testCase.expectedWarning}, warnings)
		})
	}
}

// Functional-equivalence check: re-serialising the first override of the real
// Blitz template through the parser must reproduce it (with the generator's
// all-variants marker -1).
func TestWhenBlitzOverrideLineParsed_ReproducesBlitzSidAndGuardValue(t *testing.T) {
	t.Parallel()
	// Arrange
	blitz := loadExampleTemplate(t, "Blitz.rmg.json")
	require.NotEmpty(t, blitz.ValueOverrides)
	blitzOverride := blitz.ValueOverrides[0]
	configuration := config.NewGeneratorConfig()
	configuration.ValueOverridesText = fmt.Sprintf("%s=%d", blitzOverride.SID, blitzOverride.GuardValue)

	// Act
	actual, _ := providers.NewGameRulesProvider().CreateValueOverrides(*configuration)

	// Assert
	assert.Equal(t, []entities.ValueOverride{
		{SID: blitzOverride.SID, Variant: -1, GuardValue: blitzOverride.GuardValue},
	}, actual)
}
