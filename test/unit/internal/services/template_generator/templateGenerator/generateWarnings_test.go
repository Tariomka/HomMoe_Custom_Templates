package templateGenerator_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

// Value-override text is parsed once at load now, so generation has nothing
// left to warn about.
func TestWhenValueOverridesAreConfigured_ReturnsNoWarnings(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.ValueOverrides = []config.ValueOverrideEntry{{Sid: "watchtower", GuardValue: 25000}}
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	_, warnings := generateTemplate(generator)

	// Assert
	assert.Empty(t, warnings)
}

func TestWhenValueOverridesAreConfigured_TheTemplateCarriesThem(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.ValueOverrides = []config.ValueOverrideEntry{{Sid: "watchtower", GuardValue: 25000}}
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	template, _ := generateTemplate(generator)

	// Assert
	assert.Equal(
		t,
		[]any{"watchtower", -1, 25000},
		[]any{
			template.ValueOverrides[0].SID,
			template.ValueOverrides[0].Variant,
			template.ValueOverrides[0].GuardValue,
		},
	)
}
