package stateHandler_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/validators"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWhenStatePathIsEmpty_ReturnsNoOutputPathError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := handlers.NewStateHandler(
		&test_helpers.FileServiceMock{},
		newPassingValidator(),
	)

	// Act
	_, err := handler.LoadState("", true)

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrNoOutputPath)
}

func TestWhenStatePathIsWhitespaceOnly_ReturnsNoOutputPathError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := handlers.NewStateHandler(
		&test_helpers.FileServiceMock{},
		newPassingValidator(),
	)

	// Act
	_, err := handler.LoadState("  \t ", true)

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrNoOutputPath)
}

func TestWhenStatePathIsPadded_LoadsTheTrimmedPath(t *testing.T) {
	t.Parallel()
	// Arrange
	path := gofakeit.Word() + ".gen.json"
	state := editor_state_model.NewDefaultEditorStateModel()
	fileService := &test_helpers.FileServiceMock{}
	fileService.On("LoadSettingsFile", path).Return(&state, nil)
	handler := handlers.NewStateHandler(fileService, newPassingValidator())

	// Act
	_, _ = handler.LoadState("  "+path+"  ", true)

	// Assert
	fileService.AssertCalled(t, "LoadSettingsFile", path)
}

func TestWhenSettingsFileCannotBeLoaded_PropagatesTheError(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedError := errors.New(gofakeit.Sentence(3))
	fileService := &test_helpers.FileServiceMock{}
	fileService.On("LoadSettingsFile", mock.Anything).Return(nil, expectedError)
	handler := handlers.NewStateHandler(fileService, newPassingValidator())

	// Act
	_, err := handler.LoadState(gofakeit.Word(), true)

	// Assert
	assert.ErrorIs(t, err, expectedError)
}

func TestWhenSettingsFileCannotBeLoaded_ReturnsNoState(t *testing.T) {
	t.Parallel()
	// Arrange
	fileService := &test_helpers.FileServiceMock{}
	fileService.On("LoadSettingsFile", mock.Anything).Return(nil, errors.New(gofakeit.Sentence(3)))
	handler := handlers.NewStateHandler(fileService, newPassingValidator())

	// Act
	validation, _ := handler.LoadState(gofakeit.Word(), true)

	// Assert
	assert.Nil(t, validation)
}

func TestWhenSettingsFileIsLoaded_ReturnsTheValidatedState(t *testing.T) {
	t.Parallel()
	// Arrange
	loaded := editor_state_model.NewDefaultEditorStateModel()
	loaded.TemplateName = gofakeit.Word()
	fileService := &test_helpers.FileServiceMock{}
	fileService.On("LoadSettingsFile", mock.Anything).Return(&loaded, nil)
	handler := handlers.NewStateHandler(fileService, newPassingValidator())

	// Act
	validation, err := handler.LoadState(gofakeit.Word(), true)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, loaded, validation.State)
}

func TestWhenValidationReportsIssues_ReturnsThemAsWarnings(t *testing.T) {
	t.Parallel()
	// Arrange
	firstMessage := gofakeit.Sentence(3)
	secondMessage := gofakeit.Sentence(3)
	loaded := editor_state_model.NewDefaultEditorStateModel()
	fileService := &test_helpers.FileServiceMock{}
	fileService.On("LoadSettingsFile", mock.Anything).Return(&loaded, nil)
	handler := handlers.NewStateHandler(
		fileService,
		newValidatorReporting(firstMessage, secondMessage),
	)

	// Act
	validation, _ := handler.LoadState(gofakeit.Word(), false)

	// Assert
	assert.Equal(t, []string{firstMessage, secondMessage}, validation.Warnings)
}

// The load the editor performs first asks for no fixes, because the file's own
// tournament player count is what tells it a correction happened at all.
func TestWhenIssuesAreNotFixed_ReturnsTheFilesTournamentPlayerCount(t *testing.T) {
	t.Parallel()
	// Arrange
	loaded := editor_state_model.NewDefaultEditorStateModel()
	loaded.Tournament = true
	loaded.PlayerCount = 6
	fileService := &test_helpers.FileServiceMock{}
	fileService.On("LoadSettingsFile", mock.Anything).Return(&loaded, nil)
	handler := handlers.NewStateHandler(fileService, validators.NewEditorStateValidator())

	// Act
	validation, err := handler.LoadState(gofakeit.Word(), false)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 6, validation.State.PlayerCount)
}

func TestWhenIssuesAreFixed_ReturnsTheCorrectedTournamentPlayerCount(t *testing.T) {
	t.Parallel()
	// Arrange
	loaded := editor_state_model.NewDefaultEditorStateModel()
	loaded.Tournament = true
	loaded.PlayerCount = 6
	fileService := &test_helpers.FileServiceMock{}
	fileService.On("LoadSettingsFile", mock.Anything).Return(&loaded, nil)
	handler := handlers.NewStateHandler(fileService, validators.NewEditorStateValidator())

	// Act
	validation, err := handler.LoadState(gofakeit.Word(), true)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 2, validation.State.PlayerCount)
}

func TestWhenLoadedTopologyIsBlocked_ReturnsUnsupportedTopologyError(t *testing.T) {
	t.Parallel()
	for _, fixIssues := range []bool{false, true} {
		t.Run(fmt.Sprintf("FixIssues%v_ReturnsUnsupportedTopologyError", fixIssues), func(t *testing.T) {
			t.Parallel()
			// Arrange
			handler := newHandlerLoading("NotARealTopology")

			// Act
			_, err := handler.LoadState(gofakeit.Word(), fixIssues)

			// Assert
			assert.ErrorIs(t, err, common_errors.ErrUnsupportedTopology)
		})
	}
}

func TestWhenLoadedTopologyIsRetired_ErrorTextNamesTheRetiredTopology(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newHandlerLoading("Default")

	// Act
	_, err := handler.LoadState(gofakeit.Word(), false)

	// Assert
	assert.EqualError(t, err,
		`topology "Ring" (saved as "Default") has been retired and is no longer supported; `+
			"re-create the template with a supported topology")
}

func TestWhenLoadedTopologyIsBlocked_ReturnsNoState(t *testing.T) {
	t.Parallel()
	for _, fixIssues := range []bool{false, true} {
		t.Run(fmt.Sprintf("FixIssues%v_ReturnsNoState", fixIssues), func(t *testing.T) {
			t.Parallel()
			// Arrange
			handler := newHandlerLoading("Chain")

			// Act
			validation, _ := handler.LoadState(gofakeit.Word(), fixIssues)

			// Assert
			assert.Nil(t, validation)
		})
	}
}

func TestWhenLoadedTopologyIsEmpty_LoadsWithTheEmptyTopologyWarning(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newHandlerLoading("")

	// Act
	validation, err := handler.LoadState(gofakeit.Word(), true)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{"topology is empty; using Random"}, validation.Warnings)
}

// newHandlerLoading returns a handler with the real validator whose file
// service loads a default state carrying the given topology.
func newHandlerLoading(topology config.MapTopology) handler_interfaces.IStateHandler {
	loaded := editor_state_model.NewDefaultEditorStateModel()
	loaded.Topology = topology
	fileService := &test_helpers.FileServiceMock{}
	fileService.On("LoadSettingsFile", mock.Anything).Return(&loaded, nil)
	return handlers.NewStateHandler(fileService, validators.NewEditorStateValidator())
}
