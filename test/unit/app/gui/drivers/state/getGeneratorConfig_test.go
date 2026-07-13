package state_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenConfigIsRequested_ReturnsMappedConfig(t *testing.T) {
	t.Parallel()
	// Arrange
	state := drivers.NewUIStateWithHandler(&test_helpers.TemplateHandlerMock{})

	// Act
	configuration := state.GetGeneratorConfig()

	// Assert
	assert.NotNil(t, configuration)
}

func TestWhenConfigIsRequested_TemplateNameMatchesState(t *testing.T) {
	t.Parallel()
	// Arrange
	state := drivers.NewUIStateWithHandler(&test_helpers.TemplateHandlerMock{})

	// Act
	configuration := state.GetGeneratorConfig()

	// Assert
	assert.Equal(t, state.GetStateData().TemplateName, configuration.TemplateName)
}
