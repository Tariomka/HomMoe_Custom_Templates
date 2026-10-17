package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenTheCapturedVariantIsExpanded_TheVariantZonesAreRebuilt(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	request := newSelfCapturedExpandRequest(t, handler)

	// Act
	expansion := handler.ExpandManualEdits(request)

	// Assert
	assert.Equal(t, request.Template.Variants[0].Zones, expansion.Zones)
}

func TestWhenTheCapturedVariantIsExpanded_TheVariantConnectionsAreRebuilt(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	request := newSelfCapturedExpandRequest(t, handler)

	// Act
	expansion := handler.ExpandManualEdits(request)

	// Assert
	assert.Equal(t, request.Template.Variants[0].Connections, expansion.Connections)
}

func TestWhenTheCapturedVariantIsExpanded_NoWarningsAreRaised(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	request := newSelfCapturedExpandRequest(t, handler)

	// Act
	expansion := handler.ExpandManualEdits(request)

	// Assert
	assert.Empty(t, expansion.Warnings)
}

// newSelfCapturedExpandRequest captures the deterministic default layout
// against itself and returns the request that expands it onto that layout.
func newSelfCapturedExpandRequest(
	t *testing.T,
	handler handler_interfaces.IGuiHandler,
) dtos.ManualEditExpandRequestDto {
	t.Helper()
	state := newDeterministicState()
	template := generateTemplateFor(t, handler, state)
	captured, err := handler.CaptureManualEdits(dtos.ManualEditCaptureRequestDto{
		State:       toDto(state),
		Zones:       template.Variants[0].Zones,
		Connections: template.Variants[0].Connections,
	})
	require.NoError(t, err)

	state.ManualEditSettings = captured.ManualEdits
	state.GeneratorVersion = 1
	return dtos.ManualEditExpandRequestDto{Template: template, State: toDto(state)}
}
