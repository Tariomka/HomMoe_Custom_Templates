package stateManualEdits_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWhenNoTemplateWasGenerated_EditsAreIgnored(t *testing.T) {
	t.Parallel()
	// Arrange
	handlerMock := &test_helpers.TemplateHandlerMock{}
	state := drivers.NewUIState(
		handlerMock,
		test_helpers.NewFileSystemHandler(),
		test_helpers.NewRegenerationHandler(),

		false)

	// Act
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{})

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
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{Zones: zones, Connections: connections})

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
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{Zones: zones, Connections: connections})

	// Assert
	stateData := state.GetStateData()
	assert.True(t, stateData.HasManualEdits())
}

func TestWhenTemplateExists_CurrentEditorStateIsSentForUpdate(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock, zones, connections := newGeneratedState()
	expectedState := state.GetStateData()
	updatedTemplate := test_helpers.GetDefaultTemplate()
	var updateRequest dtos.TemplateUpdateDto
	handlerMock.On("UpdateTemplate", mock.Anything).
		Run(func(arguments mock.Arguments) {
			updateRequest = arguments.Get(0).(dtos.TemplateUpdateDto)
		}).
		Return(dtos.TemplateLoadDto{Template: &updatedTemplate}, nil)

	// Act
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{Zones: zones, Connections: connections})

	// Assert
	assert.Equal(t, &editor_state_dto.EditorStateDto{EditorState: expectedState}, updateRequest.EditorState)
}

func TestWhenTemplateExists_StatusReportsAppliedCounts(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock, zones, connections := newGeneratedState()
	updatedTemplate := test_helpers.GetDefaultTemplate()
	handlerMock.On("UpdateTemplate", mock.Anything).
		Return(dtos.TemplateLoadDto{Template: &updatedTemplate}, nil)

	// Act
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{Zones: zones, Connections: connections})

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
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{Zones: zones, Connections: connections})

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
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{Zones: zones, Connections: connections})

	// Assert
	message, isError := state.GetStatus()
	assert.True(t, isError, "status: %s", message)
}

// Storing an untouched base as a manual snapshot would pin it and reapply it
// over every later regeneration, undoing the revert the user asked for.
func TestWhenApplyingAnUntouchedRevertToBase_NoManualSnapshotIsStored(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock, zones, connections := newGeneratedState()
	updatedTemplate := test_helpers.GetDefaultTemplate()
	handlerMock.On("UpdateTemplate", mock.Anything).
		Return(dtos.TemplateLoadDto{Template: &updatedTemplate}, nil)
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{Zones: zones, Connections: connections})
	base, _ := state.PreviewBaseZones()

	// Act
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{
		Zones:        base.Zones,
		Connections:  base.Connections,
		RevertToBase: true,
	})

	// Assert
	stateData := state.GetStateData()
	assert.False(t, stateData.HasManualEdits())
}

// Edits made on top of the fresh base are ordinary manual edits and must
// survive later regenerations.
func TestWhenApplyingAnEditedRevertToBase_TheEditsAreStored(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock, _, _ := newGeneratedState()
	updatedTemplate := test_helpers.GetDefaultTemplate()
	handlerMock.On("UpdateTemplate", mock.Anything).
		Return(dtos.TemplateLoadDto{Template: &updatedTemplate}, nil)
	base, _ := state.PreviewBaseZones()
	editedZones := append([]entities.Zone(nil), base.Zones...)
	editedZones[0].ManualPosition = &[2]float64{0.1, 0.2}

	// Act
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{
		Zones:        editedZones,
		Connections:  base.Connections,
		RevertToBase: true,
	})

	// Assert
	stateData := state.GetStateData()
	assert.True(t, stateData.HasManualEdits())
}

// The flag is only trustworthy for the editor session that produced the base;
// a later apply must not pick up a stale preview.
func TestWhenApplyingWithoutARevert_TheManualSnapshotIsStoredAnyway(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock, _, _ := newGeneratedState()
	updatedTemplate := test_helpers.GetDefaultTemplate()
	handlerMock.On("UpdateTemplate", mock.Anything).
		Return(dtos.TemplateLoadDto{Template: &updatedTemplate}, nil)
	base, _ := state.PreviewBaseZones()

	// Act
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{Zones: base.Zones, Connections: base.Connections})

	// Assert
	stateData := state.GetStateData()
	assert.True(t, stateData.HasManualEdits())
}

// newGeneratedState returns a State holding the default template, plus its
// mock for further expectations, and the template's zones and connections to
// edit.
func newGeneratedState() (
	*drivers.State, *test_helpers.TemplateHandlerMock, []entities.Zone, []entities.Connection) {
	handlerMock := &test_helpers.TemplateHandlerMock{}
	template := test_helpers.GetDefaultTemplate()
	handlerMock.On("GenerateTemplate", mock.Anything).Return(dtos.TemplateLoadDto{Template: &template}, nil)
	state := drivers.NewUIState(
		handlerMock,
		test_helpers.NewFileSystemHandler(),
		test_helpers.NewRegenerationHandler(),

		false)

	state.Generate()
	return state, handlerMock, template.Variants[0].Zones, template.Variants[0].Connections
}
