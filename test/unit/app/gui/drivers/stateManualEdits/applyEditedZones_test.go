package stateManualEdits_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// newGeneratedState returns a State holding the default template, plus its
// mock for further expectations, and the template's zones and connections to
// edit.
func newGeneratedState() (
	*drivers.State, *test_helpers.TemplateHandlerMock, []entities.Zone, []entities.Connection) {
	handlerMock := &test_helpers.TemplateHandlerMock{}
	template := test_helpers.GetDefaultTemplate()
	handlerMock.On("GenerateTemplate", mock.Anything).Return(dtos.TemplateLoadDto{Template: &template}, nil)
	state := drivers.NewUIStateWithHandler(handlerMock)
	state.Generate()
	return state, handlerMock, template.Variants[0].Zones, template.Variants[0].Connections
}

func TestWhenNoTemplateWasGenerated_EditsAreIgnored(t *testing.T) {
	t.Parallel()
	// Arrange
	handlerMock := &test_helpers.TemplateHandlerMock{}
	state := drivers.NewUIStateWithHandler(handlerMock)

	// Act
	state.ApplyEditedZones(nil, nil)

	// Assert
	handlerMock.AssertNotCalled(t, "UpdateTemplate")
}

func TestWhenTemplateExists_UpdatedTemplateIsStored(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock, zones, connections := newGeneratedState()
	updatedTemplate := test_helpers.GetDefaultTemplate()
	updatedTemplate.Name = gofakeit.ProductName()
	handlerMock.On("UpdateTemplate", mock.Anything).
		Return(dtos.TemplateLoadDto{Template: &updatedTemplate}, nil)

	// Act
	state.ApplyEditedZones(zones, connections)

	// Assert
	assert.Equal(t, &updatedTemplate, state.GetLastTemplate())
}

func TestWhenTemplateExists_ManualEditsAreStoredInState(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock, zones, connections := newGeneratedState()
	updatedTemplate := test_helpers.GetDefaultTemplate()
	handlerMock.On("UpdateTemplate", mock.Anything).
		Return(dtos.TemplateLoadDto{Template: &updatedTemplate}, nil)

	// Act
	state.ApplyEditedZones(zones, connections)

	// Assert
	stateData := state.GetStateData()
	assert.True(t, stateData.HasManualEdits())
}

func TestWhenTemplateExists_StatusReportsAppliedCounts(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock, zones, connections := newGeneratedState()
	updatedTemplate := test_helpers.GetDefaultTemplate()
	handlerMock.On("UpdateTemplate", mock.Anything).
		Return(dtos.TemplateLoadDto{Template: &updatedTemplate}, nil)

	// Act
	state.ApplyEditedZones(zones, connections)

	// Assert
	message, _ := state.GetStatus()
	assert.Contains(t, message, "from the editor")
}

func TestWhenUpdateRejectsTemplate_LastTemplateIsKept(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock, zones, connections := newGeneratedState()
	previousTemplate := state.GetLastTemplate()
	handlerMock.On("UpdateTemplate", mock.Anything).
		Return(dtos.TemplateLoadDto{}, common_errors.ErrProvidedTemplateInvalid)

	// Act
	state.ApplyEditedZones(zones, connections)

	// Assert
	assert.Equal(t, previousTemplate, state.GetLastTemplate())
}

func TestWhenUpdateFailsWithOtherError_ErrorStatusAsksToFixBeforeExport(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock, zones, connections := newGeneratedState()
	updatedTemplate := test_helpers.GetDefaultTemplate()
	handlerMock.On("UpdateTemplate", mock.Anything).
		Return(dtos.TemplateLoadDto{Template: &updatedTemplate}, gofakeit.ErrorValidation())

	// Act
	state.ApplyEditedZones(zones, connections)

	// Assert
	message, isError := state.GetStatus()
	assert.True(t, isError, "status: %s", message)
}
