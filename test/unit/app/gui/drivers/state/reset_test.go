package state_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWhenStateIsReset_UnsavedFlagIsCleared(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newDirtyGeneratedState()

	// Act
	state.Reset()

	// Assert
	assert.False(t, state.IsUnsaved())
}

func TestWhenStateIsReset_LastTemplateIsCleared(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newDirtyGeneratedState()

	// Act
	state.Reset()

	// Assert
	assert.Nil(t, state.GetLastTemplate())
}

func TestWhenStateIsReset_StateDataReturnsToDefault(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newDirtyGeneratedState()

	// Act
	state.Reset()

	// Assert
	assert.Equal(t, editor_state_model.NewDefaultEditorStateModel(), state.GetStateData())
}

func TestWhenStateIsReset_StatusReportsNewFile(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newDirtyGeneratedState()

	// Act
	state.Reset()

	// Assert
	message, isError := state.GetStatus()
	assert.Equal(t, []any{"New settings file.", false}, []any{message, isError})
}

// A baseline left over from a failed capture would convert the next
// document's edits against the previous document's settings.
func TestWhenStateIsResetAfterAFailedCapture_TheNextConversionUsesTheCurrentSettings(t *testing.T) {
	t.Parallel()
	// Arrange
	handlerMock := &test_helpers.TemplateHandlerMock{}
	template := test_helpers.GetDefaultTemplateModel()
	handlerMock.On("GenerateTemplate", mock.Anything).Return(dtos.TemplateLoadDto{Template: &template}, nil)
	handlerMock.On("UpdateTemplate", mock.Anything).Return(dtos.TemplateLoadDto{Template: &template}, nil)
	state := drivers.NewUIState(
		handlerMock,
		test_helpers.NewFileSystemHandler(),
		test_helpers.NewRegenerationHandler(),
		false)
	state.UpdateState(func(settings *editor_state_model.EditorState) { settings.ResourceDensityPercent = 150 })
	state.Generate()
	handlerMock.CaptureManualEditsFunc = func(dtos.ManualEditCaptureRequestDto) (dtos.ManualEditCaptureDto, error) {
		return dtos.ManualEditCaptureDto{}, gofakeit.ErrorValidation()
	}
	variant := template.Variants[0]
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{Zones: variant.Zones, Connections: variant.Connections})
	state.Reset()
	pendingEdits := editor_state_model.ManualEditSettingsFromSnapshot(variant.Zones, variant.Connections)
	state.UpdateState(func(settings *editor_state_model.EditorState) {
		settings.ManualEditSettings = pendingEdits
		settings.GeneratorVersion = 0
	})
	currentState := state.GetStateData()
	expectedBaseline := currentState.CloneWithoutManualEdits()
	var baseline editor_state_model.EditorState
	handlerMock.ConvertManualEditsFunc = func(
		request dtos.ManualEditConvertRequestDto) (dtos.ManualEditCaptureDto, error) {
		baseline = request.Baseline.EditorState
		return dtos.ManualEditCaptureDto{ManualEdits: request.ManualEdits}, nil
	}

	// Act
	state.Generate()

	// Assert
	assert.Equal(t, expectedBaseline, baseline)
}

// newDirtyGeneratedState returns a State with a generated template and an
// unsaved change, so Reset has something to clear in every dimension.
func newDirtyGeneratedState() *drivers.State {
	handlerMock := &test_helpers.TemplateHandlerMock{}
	template := test_helpers.GetDefaultTemplateModel()
	handlerMock.On("GenerateTemplate", mock.Anything).Return(dtos.TemplateLoadDto{Template: &template}, nil)
	state := drivers.NewUIState(
		handlerMock,
		test_helpers.NewFileSystemHandler(),
		test_helpers.NewRegenerationHandler(),

		false)

	state.Generate()
	state.UpdateState(func(dto *editor_state_model.EditorState) { dto.TemplateName = gofakeit.ProductName() })
	return state
}
