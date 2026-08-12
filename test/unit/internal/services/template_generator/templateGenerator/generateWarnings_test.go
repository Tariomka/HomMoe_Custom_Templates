package templateGenerator_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenValueOverridesAreAllValid_ReturnsNoWarnings(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.ValueOverridesText = "watchtower=25000"
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	_, warnings := generator.Generate()

	// Assert
	assert.Empty(t, warnings)
}

func TestWhenValueOverrideLineIsRejected_ReturnsWarning(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.ValueOverridesText = "watchtower=25000\nbad_line"
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	_, warnings := generator.Generate()

	// Assert
	assert.Equal(t, []string{"line 2: 'bad_line' is not sid=value"}, warnings)
}

func TestWhenValueOverrideLineIsRejected_StillReturnsTheTemplate(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.ValueOverridesText = "bad_line"
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	template, _ := generator.Generate()

	// Assert
	assert.NotNil(t, template)
}
