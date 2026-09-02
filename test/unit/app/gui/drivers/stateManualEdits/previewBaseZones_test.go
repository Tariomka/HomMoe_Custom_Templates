package stateManualEdits_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWhenPreviewingTheBase_TheGeneratedZonesAreReturned(t *testing.T) {
	t.Parallel()
	// Arrange
	state, _, _, _ := newGeneratedState()
	expected := test_helpers.GetDefaultTemplateModel()

	// Act
	base, _ := state.PreviewBaseZones()

	// Assert
	assert.Equal(t, template_model.ToZoneEntities(expected.Variants[0].Zones), base.Zones)
}

func TestWhenPreviewingTheBase_TheGeneratedConnectionsAreReturned(t *testing.T) {
	t.Parallel()
	// Arrange
	state, _, _, _ := newGeneratedState()
	expected := test_helpers.GetDefaultTemplateModel()

	// Act
	base, _ := state.PreviewBaseZones()

	// Assert
	assert.Equal(t, template_model.ToConnectionEntities(expected.Variants[0].Connections), base.Connections)
}

func TestWhenPreviewingTheBase_SuccessIsReported(t *testing.T) {
	t.Parallel()
	// Arrange
	state, _, _, _ := newGeneratedState()

	// Act
	_, ok := state.PreviewBaseZones()

	// Assert
	assert.True(t, ok)
}

// Nothing may change on screen until the user applies, otherwise cancelling
// the editor cannot bring the edited template back.
func TestWhenPreviewingTheBase_TheLiveTemplateIsUntouched(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock, zones, connections := newGeneratedState()
	editedTemplate := test_helpers.GetDefaultTemplateModel()
	editedTemplate.Name = gofakeit.ProductName()
	handlerMock.On("UpdateTemplate", mock.Anything).
		Return(dtos.TemplateLoadDto{Template: &editedTemplate}, nil)
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{Zones: zones, Connections: connections})

	// Act
	state.PreviewBaseZones()

	// Assert
	assert.Equal(t, &editedTemplate, state.GetLastTemplate())
}

func TestWhenPreviewingTheBase_TheStoredManualEditsAreUntouched(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock, zones, connections := newGeneratedState()
	editedTemplate := test_helpers.GetDefaultTemplateModel()
	handlerMock.On("UpdateTemplate", mock.Anything).
		Return(dtos.TemplateLoadDto{Template: &editedTemplate}, nil)
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{Zones: zones, Connections: connections})

	// Act
	state.PreviewBaseZones()

	// Assert
	stateData := state.GetStateData()
	assert.True(t, stateData.HasManualEdits())
}

func TestWhenGenerationFails_NoBaseIsReported(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newFailingState()

	// Act
	_, ok := state.PreviewBaseZones()

	// Assert
	assert.False(t, ok)
}

func TestWhenGenerationFails_NoBaseZonesAreReturned(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newFailingState()

	// Act
	base, _ := state.PreviewBaseZones()

	// Assert
	assert.Empty(t, base.Zones)
}

// newFailingState returns a State whose generator always errors.
func newFailingState() *drivers.State {
	handlerMock := &test_helpers.TemplateHandlerMock{}
	handlerMock.On("GenerateTemplate", mock.Anything).Return(dtos.TemplateLoadDto{}, gofakeit.ErrorValidation())

	return drivers.NewUIState(
		handlerMock,
		test_helpers.NewFileSystemHandler(),
		test_helpers.NewRegenerationHandler(),
		false)
}
