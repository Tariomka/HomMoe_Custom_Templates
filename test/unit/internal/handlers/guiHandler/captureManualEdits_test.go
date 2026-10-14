package guiHandler_test

import (
	"slices"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenTheGeneratedVariantIsCapturedAgainstItself_NoZoneEditCarriesOverrides(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	state := newDeterministicState()
	template := generateTemplateFor(t, handler, state)
	request := dtos.ManualEditCaptureRequestDto{
		State:       toDto(state),
		Zones:       template.Variants[0].Zones,
		Connections: template.Variants[0].Connections,
	}

	// Act
	captured, err := handler.CaptureManualEdits(request)

	// Assert
	require.NoError(t, err)
	assert.Empty(t, zoneEditsWithOverrides(captured.ManualEdits))
}

func TestWhenTheGeneratedVariantIsCapturedAgainstItself_NoConnectionEditCarriesOverrides(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	state := newDeterministicState()
	template := generateTemplateFor(t, handler, state)
	request := dtos.ManualEditCaptureRequestDto{
		State:       toDto(state),
		Zones:       template.Variants[0].Zones,
		Connections: template.Variants[0].Connections,
	}

	// Act
	captured, err := handler.CaptureManualEdits(request)

	// Assert
	require.NoError(t, err)
	assert.Empty(t, connectionEditsWithOverrides(captured.ManualEdits))
}

func TestWhenTheGeneratedVariantIsCaptured_OneEditIsRecordedPerZone(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	state := newDeterministicState()
	template := generateTemplateFor(t, handler, state)
	request := dtos.ManualEditCaptureRequestDto{
		State:       toDto(state),
		Zones:       template.Variants[0].Zones,
		Connections: template.Variants[0].Connections,
	}

	// Act
	captured, err := handler.CaptureManualEdits(request)

	// Assert
	require.NoError(t, err)
	assert.Len(t, captured.ManualEdits.ManualZones, len(template.Variants[0].Zones))
}

func TestWhenAZoneIsEditedBeforeCapture_ThatZoneCarriesTheOverride(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	state := newDeterministicState()
	template := generateTemplateFor(t, handler, state)
	zones := template.Variants[0].Zones
	editedZones := slices.Clone(zones)
	editedZones[0] = editedZones[0].Clone()
	editedZones[0].Size = zones[0].Size + 1
	request := dtos.ManualEditCaptureRequestDto{
		State:       toDto(state),
		Zones:       editedZones,
		Connections: template.Variants[0].Connections,
	}

	// Act
	captured, err := handler.CaptureManualEdits(request)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, new(editedZones[0].Size), captured.ManualEdits.ManualZones[0].FieldOverrides.Size)
}

func TestWhenTheStateIsRejected_CaptureReturnsTheError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	state := editor_state_model.NewDefaultEditorStateModel()
	state.TemplateName = ""

	// Act
	_, err := handler.CaptureManualEdits(dtos.ManualEditCaptureRequestDto{State: toDto(state)})

	// Assert
	assert.Error(t, err)
}
