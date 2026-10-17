package guiHandler_test

import (
	"slices"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenTheSnapshotOfTheBaselineLayoutIsConverted_NoZoneEditCarriesOverrides(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	baseline := newDeterministicState()
	template := generateTemplateFor(t, handler, baseline)
	request := dtos.ManualEditConvertRequestDto{
		Baseline: toDto(baseline),
		ManualEdits: editor_state_model.ManualEditSettingsFromSnapshot(
			template.Variants[0].Zones, template.Variants[0].Connections),
	}

	// Act
	converted, err := handler.ConvertManualEdits(request)

	// Assert
	require.NoError(t, err)
	assert.Empty(t, zoneEditsWithOverrides(converted.ManualEdits))
}

func TestWhenTheSnapshotOfTheBaselineLayoutIsConverted_NoConnectionEditCarriesOverrides(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	baseline := newDeterministicState()
	template := generateTemplateFor(t, handler, baseline)
	request := dtos.ManualEditConvertRequestDto{
		Baseline: toDto(baseline),
		ManualEdits: editor_state_model.ManualEditSettingsFromSnapshot(
			template.Variants[0].Zones, template.Variants[0].Connections),
	}

	// Act
	converted, err := handler.ConvertManualEdits(request)

	// Assert
	require.NoError(t, err)
	assert.Empty(t, connectionEditsWithOverrides(converted.ManualEdits))
}

func TestWhenTheSnapshotHoldsAnEditedZone_ConversionRecordsTheEditAsAnOverride(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	baseline := newDeterministicState()
	template := generateTemplateFor(t, handler, baseline)
	zones := template.Variants[0].Zones
	editedZones := slices.Clone(zones)
	editedZones[0] = editedZones[0].Clone()
	editedZones[0].Size = zones[0].Size + 1
	request := dtos.ManualEditConvertRequestDto{
		Baseline:    toDto(baseline),
		ManualEdits: editor_state_model.ManualEditSettingsFromSnapshot(editedZones, template.Variants[0].Connections),
	}

	// Act
	converted, err := handler.ConvertManualEdits(request)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, new(editedZones[0].Size), converted.ManualEdits.ManualZones[0].FieldOverrides.Size)
}

func TestWhenTheBaselineIsRejected_ConvertReturnsTheError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	baseline := editor_state_model.NewDefaultEditorStateModel()
	baseline.TemplateName = ""

	// Act
	_, err := handler.ConvertManualEdits(dtos.ManualEditConvertRequestDto{Baseline: toDto(baseline)})

	// Assert
	assert.Error(t, err)
}
