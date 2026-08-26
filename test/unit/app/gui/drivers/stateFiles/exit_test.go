package stateFiles_test

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

// newUnsavedState returns a State with a generated template and an unsaved
// change, plus a flag pointer reporting whether Exit closed the application.
func newUnsavedState() (state *drivers.State, exited *bool) {
	handlerMock := &test_helpers.TemplateHandlerMock{}
	template := test_helpers.GetDefaultTemplate()
	handlerMock.On("GenerateTemplate", mock.Anything).Return(dtos.TemplateLoadDto{Template: &template}, nil)
	state = drivers.NewUIState(
		handlerMock,
		test_helpers.NewFileSystemHandler(),
		test_helpers.NewRegenerationHandler(),

		false)

	state.Generate()
	state.UpdateState(func(dto *editor_state_model.EditorState) { dto.TemplateName = gofakeit.ProductName() })

	exited = new(bool)
	state.SetOnExit(func() { *exited = true })
	return state, exited
}
