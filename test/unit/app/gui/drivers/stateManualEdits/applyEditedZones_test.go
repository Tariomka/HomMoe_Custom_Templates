package stateManualEdits_test

import (
	"fmt"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
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
	updatedTemplate := test_helpers.GetDefaultTemplateModel()
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
	updatedTemplate := test_helpers.GetDefaultTemplateModel()
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
	updatedTemplate := test_helpers.GetDefaultTemplateModel()
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
	updatedTemplate := test_helpers.GetDefaultTemplateModel()
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
	updatedTemplate := test_helpers.GetDefaultTemplateModel()
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
	updatedTemplate := test_helpers.GetDefaultTemplateModel()
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
	updatedTemplate := test_helpers.GetDefaultTemplateModel()
	handlerMock.On("UpdateTemplate", mock.Anything).
		Return(dtos.TemplateLoadDto{Template: &updatedTemplate}, nil)
	base, _ := state.PreviewBaseZones()
	editedZones := append([]template_model.Zone(nil), base.Zones...)
	editedZones[0].ManualPosition = new(data.NewVec2(0.1, 0.2))

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
	updatedTemplate := test_helpers.GetDefaultTemplateModel()
	handlerMock.On("UpdateTemplate", mock.Anything).
		Return(dtos.TemplateLoadDto{Template: &updatedTemplate}, nil)
	base, _ := state.PreviewBaseZones()

	// Act
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{Zones: base.Zones, Connections: base.Connections})

	// Assert
	stateData := state.GetStateData()
	assert.True(t, stateData.HasManualEdits())
}

// The applied layout is saved with the rest of the state, so committing it is
// an unsaved change even though the scalar options never moved.
func TestWhenTheFirstManualSnapshotIsCommitted_TheDocumentBecomesUnsaved(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock, zones, connections := newGeneratedState()
	expectAcceptedUpdate(handlerMock)

	// Act
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{Zones: zones, Connections: connections})

	// Assert
	assert.True(t, state.IsUnsaved())
}

func TestWhenNoTemplateWasGenerated_TheDocumentStaysSaved(t *testing.T) {
	t.Parallel()
	// Arrange
	state := drivers.NewUIState(
		&test_helpers.TemplateHandlerMock{},
		test_helpers.NewFileSystemHandler(),
		test_helpers.NewRegenerationHandler(),
		false)

	// Act
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{Zones: []template_model.Zone{{Name: "Zone A"}}})

	// Assert
	assert.False(t, state.IsUnsaved())
}

func TestWhenUpdateRejectsTemplate_NoManualSnapshotIsCommitted(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock, zones, connections := newGeneratedState()
	handlerMock.On("UpdateTemplate", mock.Anything).
		Return(dtos.TemplateLoadDto{}, common_errors.ErrProvidedTemplateInvalid)

	// Act
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{Zones: zones, Connections: connections})

	// Assert
	stateData := state.GetStateData()
	assert.False(t, stateData.HasManualEdits())
}

func TestWhenUpdateRejectsTemplate_TheDocumentStaysSaved(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock, zones, connections := newGeneratedState()
	handlerMock.On("UpdateTemplate", mock.Anything).
		Return(dtos.TemplateLoadDto{}, common_errors.ErrProvidedTemplateInvalid)

	// Act
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{Zones: zones, Connections: connections})

	// Assert
	assert.False(t, state.IsUnsaved())
}

// A template that comes back with validation warnings is still the live one,
// so its edits are committed and the document is dirty.
func TestWhenUpdateFailsWithOtherError_TheDocumentBecomesUnsaved(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock, zones, connections := newGeneratedState()
	updatedTemplate := test_helpers.GetDefaultTemplateModel()
	handlerMock.On("UpdateTemplate", mock.Anything).
		Return(dtos.TemplateLoadDto{Template: &updatedTemplate}, gofakeit.ErrorValidation())

	// Act
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{Zones: zones, Connections: connections})

	// Assert
	assert.True(t, state.IsUnsaved())
}

// The handler rebuilds the request's roads in place, so a revert compared
// after the update no longer recognises the base it was handed.
func TestWhenTheUpdateRewritesTheRequestRoads_TheUntouchedRevertStillClearsTheSnapshot(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock, zones, connections := newGeneratedState()
	expectRoadRewritingUpdate(handlerMock)
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{Zones: zones, Connections: connections})
	base, _ := state.PreviewBaseZones()

	// Act
	state.ApplyEditedZones(copyZoneSet(base))

	// Assert
	stateData := state.GetStateData()
	assert.False(t, stateData.HasManualEdits())
}

// A rejected apply consumes the previewed base like any other: the editor is
// expected to preview again rather than retry against a stale layout.
func TestWhenARejectedRevertIsRetried_TheConsumedBaseIsNoLongerRecognised(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock, _, _ := newGeneratedState()
	handlerMock.On("UpdateTemplate", mock.Anything).
		Return(dtos.TemplateLoadDto{}, common_errors.ErrProvidedTemplateInvalid).Once()
	expectAcceptedUpdate(handlerMock)
	base, _ := state.PreviewBaseZones()
	state.ApplyEditedZones(copyZoneSet(base))

	// Act
	state.ApplyEditedZones(copyZoneSet(base))

	// Assert
	stateData := state.GetStateData()
	assert.True(t, stateData.HasManualEdits())
}

func TestWhenAnUntouchedBaseIsAppliedWithoutManualEdits_NoSnapshotIsStored(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock, _, _ := newGeneratedState()
	expectAcceptedUpdate(handlerMock)
	base, _ := state.PreviewBaseZones()

	// Act
	state.ApplyEditedZones(copyZoneSet(base))

	// Assert
	stateData := state.GetStateData()
	assert.False(t, stateData.HasManualEdits())
}

func TestWhenAnUntouchedBaseIsAppliedWithoutManualEdits_TheDocumentStaysSaved(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock, _, _ := newGeneratedState()
	expectAcceptedUpdate(handlerMock)
	base, _ := state.PreviewBaseZones()

	// Act
	state.ApplyEditedZones(copyZoneSet(base))

	// Assert
	assert.False(t, state.IsUnsaved())
}

// newGeneratedState returns a State holding the default template, plus its
// mock for further expectations, and the template's zones and connections to
// edit.
func newGeneratedState() (
	*drivers.State, *test_helpers.TemplateHandlerMock, []template_model.Zone, []template_model.Connection) {
	handlerMock := &test_helpers.TemplateHandlerMock{}
	template := test_helpers.GetDefaultTemplateModel()
	handlerMock.On("GenerateTemplate", mock.Anything).Return(dtos.TemplateLoadDto{Template: &template}, nil)
	state := drivers.NewUIState(
		handlerMock,
		test_helpers.NewFileSystemHandler(),
		test_helpers.NewRegenerationHandler(),
		false)

	state.Generate()
	variant := template.Variants[0]
	return state,
		handlerMock,
		variant.Zones,
		variant.Connections
}

// expectAcceptedUpdate makes the handler accept every update, returning a
// fresh template as the live one.
func expectAcceptedUpdate(handlerMock *test_helpers.TemplateHandlerMock) {
	updatedTemplate := test_helpers.GetDefaultTemplateModel()
	handlerMock.On("UpdateTemplate", mock.Anything).
		Return(dtos.TemplateLoadDto{Template: &updatedTemplate}, nil)
}

// expectRoadRewritingUpdate accepts every update, first stamping a different
// road on each requested zone the way the real handler rebuilds them in place.
func expectRoadRewritingUpdate(handlerMock *test_helpers.TemplateHandlerMock) {
	updatedTemplate := test_helpers.GetDefaultTemplateModel()
	rebuildCount := 0
	handlerMock.On("UpdateTemplate", mock.Anything).
		Run(func(arguments mock.Arguments) {
			rebuildCount++
			request := arguments.Get(0).(dtos.TemplateUpdateDto)
			for index := range request.Zones {
				request.Zones[index].Roads = []template_model.Road{
					{Type: fmt.Sprintf("rebuilt-%d", rebuildCount)}}
			}
		}).
		Return(dtos.TemplateLoadDto{Template: &updatedTemplate}, nil)
}

// copyZoneSet mimics the zone editor dialog, which hands the driver its own
// top-level slices holding the same zone values.
func copyZoneSet(source dtos.ZoneEditorZonesDto) dtos.ZoneEditorZonesDto {
	return dtos.ZoneEditorZonesDto{
		Zones:        append([]template_model.Zone(nil), source.Zones...),
		Connections:  append([]template_model.Connection(nil), source.Connections...),
		RevertToBase: true,
	}
}
