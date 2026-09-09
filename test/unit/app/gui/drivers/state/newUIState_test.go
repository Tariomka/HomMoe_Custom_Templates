package state_test

import (
	"errors"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

// pickerHint is the recovery instruction every failed lookup must end with: the
// game only reads templates from its own folder (AGENTS.md 2.7), so a user whose
// folder was not found has to point the editor at it for this session.
const pickerHint = "Choose the game templates folder using the output folder picker before exporting."

func TestWhenStateIsCreated_StateDataIsDefault(t *testing.T) {
	t.Parallel()
	// Arrange
	handlerMock := &test_helpers.TemplateHandlerMock{}

	// Act
	state := drivers.NewUIState(
		handlerMock,
		test_helpers.NewFileSystemHandler(),
		test_helpers.NewRegenerationHandler(),

		false)

	// Assert
	assert.Equal(t, editor_state_model.NewDefaultEditorStateModel(), state.GetStateData())
}

func TestWhenStateIsCreated_NoDialogIsOpen(t *testing.T) {
	t.Parallel()
	// Arrange
	handlerMock := &test_helpers.TemplateHandlerMock{}

	// Act
	state := drivers.NewUIState(
		handlerMock,
		test_helpers.NewFileSystemHandler(),
		test_helpers.NewRegenerationHandler(),

		false)

	// Assert
	assert.False(t, state.GetDialogHost().IsOpen())
}

func TestWhenTemplateDirLookupIsSkipped_OutputPathIsEmpty(t *testing.T) {
	t.Parallel()
	// Arrange
	handlerMock := &test_helpers.TemplateHandlerMock{}

	// Act
	state := drivers.NewUIState(
		handlerMock,
		test_helpers.NewFileSystemHandler(),
		test_helpers.NewRegenerationHandler(),

		false)

	// Assert
	assert.Empty(t, state.GetOutputPath())
}

func TestWhenTemplateDirLookupIsSkipped_TheHandlerIsNeverAsked(t *testing.T) {
	t.Parallel()
	// Arrange
	fileSystemMock := &test_helpers.FileSystemHandlerMock{}

	// Act
	drivers.NewUIState(
		&test_helpers.TemplateHandlerMock{},
		fileSystemMock,
		test_helpers.NewRegenerationHandler(),
		false)

	// Assert
	fileSystemMock.AssertNotCalled(t, "FindGameTemplateDirectory")
}

func TestWhenTemplateDirLookupIsRequested_TheHandlerIsAsked(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedDirectory := gofakeit.Word()

	// Act
	_, fileSystemMock := newStateWithLookup(expectedDirectory, nil)

	// Assert
	fileSystemMock.AssertCalled(t, "FindGameTemplateDirectory")
}

func TestWhenTemplateDirIsFound_OutputPathIsTheDetectedDirectory(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedDirectory := gofakeit.Word()

	// Act
	state, _ := newStateWithLookup(expectedDirectory, nil)

	// Assert
	assert.Equal(t, expectedDirectory, state.GetOutputPath())
}

func TestWhenTemplateDirIsFound_NoErrorIsReported(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedDirectory := gofakeit.Word()

	// Act
	state, _ := newStateWithLookup(expectedDirectory, nil)

	// Assert
	_, isError := state.GetStatus()
	assert.False(t, isError)
}

func TestWhenTemplateDirIsNotFound_OutputPathStaysEmpty(t *testing.T) {
	t.Parallel()
	// Arrange
	wrapped := errors.Join(common_errors.ErrTemplatesDirNotFound, errors.New(gofakeit.Sentence(3)))

	// Act
	state, _ := newStateWithLookup("", wrapped)

	// Assert
	assert.Empty(t, state.GetOutputPath())
}

func TestWhenTemplateDirIsNotFound_TheStatusDirectsToThePicker(t *testing.T) {
	t.Parallel()
	// Arrange
	wrapped := errors.Join(common_errors.ErrTemplatesDirNotFound, errors.New(gofakeit.Sentence(3)))

	// Act
	state, _ := newStateWithLookup("", wrapped)

	// Assert
	message, _ := state.GetStatus()
	assert.Equal(t, "Game template directory not found. "+pickerHint, message)
}

func TestWhenTemplateDirIsNotFound_TheStatusIsAnError(t *testing.T) {
	t.Parallel()
	// Arrange
	wrapped := errors.Join(common_errors.ErrTemplatesDirNotFound, errors.New(gofakeit.Sentence(3)))

	// Act
	state, _ := newStateWithLookup("", wrapped)

	// Assert
	_, isError := state.GetStatus()
	assert.True(t, isError)
}

func TestWhenTemplateDirLookupFails_OutputPathStaysEmpty(t *testing.T) {
	t.Parallel()
	// Arrange
	lookupErr := errors.New(gofakeit.Sentence(3))

	// Act
	state, _ := newStateWithLookup("", lookupErr)

	// Assert
	assert.Empty(t, state.GetOutputPath())
}

func TestWhenTemplateDirLookupFails_TheStatusKeepsTheUnderlyingError(t *testing.T) {
	t.Parallel()
	// Arrange
	lookupErr := errors.New(gofakeit.Sentence(3))

	// Act
	state, _ := newStateWithLookup("", lookupErr)

	// Assert
	message, _ := state.GetStatus()
	assert.Equal(t, "Failed to find game template directory: "+lookupErr.Error()+". "+pickerHint, message)
}

func TestWhenTemplateDirLookupFails_TheStatusIsAnError(t *testing.T) {
	t.Parallel()
	// Arrange
	lookupErr := errors.New(gofakeit.Sentence(3))

	// Act
	state, _ := newStateWithLookup("", lookupErr)

	// Assert
	_, isError := state.GetStatus()
	assert.True(t, isError)
}

func TestWhenTemplateDirLookupReturnsABlankPath_OutputPathStaysEmpty(t *testing.T) {
	t.Parallel()
	// Arrange
	blankDirectory := "   "

	// Act
	state, _ := newStateWithLookup(blankDirectory, nil)

	// Assert
	assert.Empty(t, state.GetOutputPath())
}

func TestWhenTemplateDirLookupReturnsABlankPath_TheStatusDirectsToThePicker(t *testing.T) {
	t.Parallel()
	// Arrange
	blankDirectory := "   "

	// Act
	state, _ := newStateWithLookup(blankDirectory, nil)

	// Assert
	message, _ := state.GetStatus()
	assert.Equal(t, "Game template directory not found. "+pickerHint, message)
}

// A path handed back alongside an error is not a destination the editor may
// export into: only a clean lookup authorizes a target.
func TestWhenTemplateDirLookupReportsAPathAndAnError_OutputPathStaysEmpty(t *testing.T) {
	t.Parallel()
	// Arrange
	lookupErr := errors.New(gofakeit.Sentence(3))

	// Act
	state, _ := newStateWithLookup(gofakeit.Word(), lookupErr)

	// Assert
	assert.Empty(t, state.GetOutputPath())
}

func TestWhenTemplateDirLookupReportsAPathAndAnError_TheStatusIsAnError(t *testing.T) {
	t.Parallel()
	// Arrange
	lookupErr := errors.New(gofakeit.Sentence(3))

	// Act
	state, _ := newStateWithLookup(gofakeit.Word(), lookupErr)

	// Assert
	_, isError := state.GetStatus()
	assert.True(t, isError)
}

// newStateWithLookup supplies a deterministic lookup without changing the host environment.
func newStateWithLookup(directory string, lookupErr error) (*drivers.State, *test_helpers.FileSystemHandlerMock) {
	fileSystemMock := &test_helpers.FileSystemHandlerMock{}
	fileSystemMock.On("FindGameTemplateDirectory").Return(directory, lookupErr)

	return drivers.NewUIState(
		&test_helpers.TemplateHandlerMock{},
		fileSystemMock,
		test_helpers.NewRegenerationHandler(),
		true), fileSystemMock
}
