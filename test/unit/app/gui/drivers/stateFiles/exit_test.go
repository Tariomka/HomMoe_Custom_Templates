package stateFiles_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWhenUnsavedAndExitPressedOnce_ApplicationDoesNotExit(t *testing.T) {
	t.Parallel()
	// Arrange
	state, exited := newUnsavedState()

	// Act
	state.Exit()

	// Assert
	assert.False(t, *exited)
}

func TestWhenUnsavedAndExitPressedOnce_WarningStatusIsSet(t *testing.T) {
	t.Parallel()
	// Arrange
	state, _ := newUnsavedState()

	// Act
	state.Exit()

	// Assert
	_, isError := state.GetStatus()
	assert.True(t, isError)
}

func TestWhenUnsavedAndExitPressedTwice_ApplicationExits(t *testing.T) {
	t.Parallel()
	// Arrange
	state, exited := newUnsavedState()
	state.Exit()

	// Act
	state.Exit()

	// Assert
	assert.True(t, *exited)
}

func TestWhenSavedAndExitPressed_ApplicationExits(t *testing.T) {
	t.Parallel()
	// Arrange
	state := drivers.NewUIState(
		&test_helpers.TemplateHandlerMock{},
		test_helpers.NewFileSystemHandler(),
		test_helpers.NewRegenerationHandler(),
		false)

	exited := false
	state.SetOnExit(func() { exited = true })

	// Act
	state.Exit()

	// Assert
	assert.True(t, exited)
}

func TestWhenEditsFollowExitConfirmation_ExitIsBlockedAgain(t *testing.T) {
	t.Parallel()
	// Arrange
	state, exited := newUnsavedState()
	state.Exit()
	state.UpdateState(func(dto *editor_state_model.EditorState) { dto.TemplateName = gofakeit.ProductName() })

	// Act
	state.Exit()

	// Assert
	assert.False(t, *exited)
}

// A manual layout is persisted state, so committing one after the warning has
// to demand a fresh confirmation just like a scalar edit does.
func TestWhenManualEditsFollowExitConfirmation_ExitIsBlockedAgain(t *testing.T) {
	t.Parallel()
	// Arrange
	state, exited, zones, connections := newUnsavedGeneratedState()
	state.Exit()

	// Act
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{Zones: zones, Connections: connections})
	state.Exit()

	// Assert
	assert.False(t, *exited)
}

// An untouched revert drops the committed layout, which is as much of a
// change as making one.
func TestWhenAnUntouchedRevertFollowsExitConfirmation_ExitIsBlockedAgain(t *testing.T) {
	t.Parallel()
	// Arrange
	state, exited, zones, connections := newUnsavedGeneratedState()
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{Zones: zones, Connections: connections})
	base, _ := state.PreviewBaseZones()
	state.Exit()

	// Act
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{
		Zones:        base.Zones,
		Connections:  base.Connections,
		RevertToBase: true,
	})
	state.Exit()

	// Assert
	assert.False(t, *exited)
}

// Re-applying the same layout commits nothing, so it must not disarm a
// confirmation the user already answered.
func TestWhenAnIdenticalApplyFollowsExitConfirmation_TheConfirmationStillHolds(t *testing.T) {
	t.Parallel()
	// Arrange
	state, exited, zones, connections := newUnsavedGeneratedState()
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{Zones: zones, Connections: connections})
	state.Exit()

	// Act
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{Zones: zones, Connections: connections})
	state.Exit()

	// Assert
	assert.True(t, *exited)
}

// A validation warning travels with an accepted template, so re-applying the
// same layout still commits nothing and must leave the answered confirmation
// standing.
func TestWhenAnIdenticalApplyWithAValidationWarningFollowsExitConfirmation_TheConfirmationStillHolds(t *testing.T) {
	t.Parallel()
	// Arrange
	handlerMock := &test_helpers.TemplateHandlerMock{}
	template := test_helpers.GetDefaultTemplateModel()
	handlerMock.On("GenerateTemplate", mock.Anything).Return(dtos.TemplateLoadDto{Template: &template}, nil)
	updatedTemplate := test_helpers.GetDefaultTemplateModel()
	handlerMock.On("UpdateTemplate", mock.Anything).
		Return(dtos.TemplateLoadDto{Template: &updatedTemplate}, gofakeit.ErrorValidation())
	state := drivers.NewUIState(
		handlerMock,
		test_helpers.NewFileSystemHandler(),
		test_helpers.NewRegenerationHandler(),
		false)

	state.Generate()
	exited := false
	state.SetOnExit(func() { exited = true })
	variant := template.Variants[0]
	edits := dtos.ZoneEditorZonesDto{Zones: variant.Zones, Connections: variant.Connections}
	state.ApplyEditedZones(edits)
	state.Exit()

	// Act
	state.ApplyEditedZones(edits)
	state.Exit()

	// Assert
	assert.True(t, exited)
}

// newUnsavedState returns a State with a generated template and an unsaved
// change, plus a flag pointer reporting whether Exit closed the application.
func newUnsavedState() (state *drivers.State, exited *bool) {
	state, exited, _, _ = newUnsavedGeneratedState()
	return state, exited
}

// newUnsavedGeneratedState additionally accepts manual updates and hands back
// the generated layout to edit.
func newUnsavedGeneratedState() (
	state *drivers.State,
	exited *bool,
	zones []template_model.Zone,
	connections []template_model.Connection) {
	handlerMock := &test_helpers.TemplateHandlerMock{}
	template := test_helpers.GetDefaultTemplateModel()
	handlerMock.On("GenerateTemplate", mock.Anything).Return(dtos.TemplateLoadDto{Template: &template}, nil)
	updatedTemplate := test_helpers.GetDefaultTemplateModel()
	handlerMock.On("UpdateTemplate", mock.Anything).
		Return(dtos.TemplateLoadDto{Template: &updatedTemplate}, nil)
	state = drivers.NewUIState(
		handlerMock,
		test_helpers.NewFileSystemHandler(),
		test_helpers.NewRegenerationHandler(),
		false)

	state.Generate()
	state.UpdateState(func(dto *editor_state_model.EditorState) { dto.TemplateName = gofakeit.ProductName() })

	exited = new(bool)
	state.SetOnExit(func() { *exited = true })
	variant := template.Variants[0]
	return state, exited, variant.Zones, variant.Connections
}
