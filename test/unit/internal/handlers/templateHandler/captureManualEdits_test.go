package templateHandler_test

import (
	"errors"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWhenCaptureValidationRejectsTheState_ReturnsTheRejection(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	state := editor_state_model.NewDefaultEditorStateModel()
	rejection := errors.New(gofakeit.Sentence(3))
	fixture.stateHandler.On("ValidateEditorState", state, true).
		Return(editor_state_dto.EditorStateValidationDto{State: state, Rejection: rejection})

	// Act
	_, err := fixture.handler.CaptureManualEdits(dtos.ManualEditCaptureRequestDto{State: toDto(state)})

	// Assert
	assert.ErrorIs(t, err, rejection)
}

func TestWhenCaptureValidationRejectsTheState_NeverCaptures(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	state := editor_state_model.NewDefaultEditorStateModel()
	fixture.stateHandler.On("ValidateEditorState", state, true).
		Return(editor_state_dto.EditorStateValidationDto{State: state, Rejection: errors.New(gofakeit.Sentence(3))})

	// Act
	_, _ = fixture.handler.CaptureManualEdits(dtos.ManualEditCaptureRequestDto{State: toDto(state)})

	// Assert
	fixture.manualEditDelta.AssertNotCalled(t, "Capture", mock.Anything)
}

func TestWhenCaptureTemplateNameIsEmpty_ReturnsNoTemplateNameError(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	state := editor_state_model.NewDefaultEditorStateModel()
	configuration := config.NewGeneratorConfig()
	configuration.TemplateName = ""
	fixture.stateHandler.On("ValidateEditorState", state, true).
		Return(editor_state_dto.EditorStateValidationDto{State: state})
	fixture.mapper.On("FromEditorState", state).Return(configuration)

	// Act
	_, err := fixture.handler.CaptureManualEdits(dtos.ManualEditCaptureRequestDto{State: toDto(state)})

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrNoTemplateName)
}

func TestWhenCaptureGenerationFails_ReturnsTheGenerationError(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	state := editor_state_model.NewDefaultEditorStateModel()
	generationError := errors.New(gofakeit.Sentence(3))
	fixture.stateHandler.On("ValidateEditorState", state, true).
		Return(editor_state_dto.EditorStateValidationDto{State: state})
	fixture.mapper.On("FromEditorState", state).Return(namedConfiguration())
	fixture.templateGenerator.On("SetConfiguration", mock.Anything).Return()
	fixture.templateGenerator.On("Generate").Return(nil, nil, generationError)

	// Act
	_, err := fixture.handler.CaptureManualEdits(dtos.ManualEditCaptureRequestDto{State: toDto(state)})

	// Assert
	assert.ErrorIs(t, err, generationError)
}

func TestWhenCaptureGenerationYieldsNoTemplate_ReturnsGeneratedTemplateInvalidError(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	state := editor_state_model.NewDefaultEditorStateModel()
	fixture.stateHandler.On("ValidateEditorState", state, true).
		Return(editor_state_dto.EditorStateValidationDto{State: state})
	fixture.mapper.On("FromEditorState", state).Return(namedConfiguration())
	fixture.templateGenerator.On("SetConfiguration", mock.Anything).Return()
	fixture.templateGenerator.On("Generate").Return(nil, nil)

	// Act
	_, err := fixture.handler.CaptureManualEdits(dtos.ManualEditCaptureRequestDto{State: toDto(state)})

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrGeneratedTemplateInvalid)
}

func TestWhenCaptureFails_NeverCaptures(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	state := editor_state_model.NewDefaultEditorStateModel()
	fixture.stateHandler.On("ValidateEditorState", state, true).
		Return(editor_state_dto.EditorStateValidationDto{State: state})
	fixture.mapper.On("FromEditorState", state).Return(namedConfiguration())
	fixture.templateGenerator.On("SetConfiguration", mock.Anything).Return()
	fixture.templateGenerator.On("Generate").Return(nil, nil)

	// Act
	_, _ = fixture.handler.CaptureManualEdits(dtos.ManualEditCaptureRequestDto{State: toDto(state)})

	// Assert
	fixture.manualEditDelta.AssertNotCalled(t, "Capture", mock.Anything)
}

func TestWhenCaptureSucceeds_PassesTheFreshAndEditedLayoutToTheDeltaService(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	state := editor_state_model.NewDefaultEditorStateModel()
	configuration := namedConfiguration()
	configuration.Topology = config.TopologySquare
	freshZones := []template_model.Zone{{Name: "Fresh"}}
	freshConnections := []template_model.Connection{{Name: "Fresh-Link"}}
	editedZones := []template_model.Zone{{Name: "Edited", Size: gofakeit.Float64Range(1, 3)}}
	editedConnections := []template_model.Connection{{Name: "Edited-Link"}}
	fixture.arrangeGeneration(state, configuration, newTemplateWithLayout(freshZones, freshConnections))
	fixture.manualEditDelta.On("Capture", mock.Anything).Return(editor_state_model.ManualEditSettings{})

	// Act
	_, _ = fixture.handler.CaptureManualEdits(dtos.ManualEditCaptureRequestDto{
		State:       toDto(state),
		Zones:       editedZones,
		Connections: editedConnections,
	})

	// Assert
	fixture.manualEditDelta.AssertCalled(t, "Capture", models.ManualEditCaptureRequest{
		FreshZones:       freshZones,
		FreshConnections: freshConnections,
		Zones:            editedZones,
		Connections:      editedConnections,
		Configuration:    configuration,
		NonDeterministic: false,
	})
}

func TestWhenCaptureConfigurationHasARandomTopology_FlagsTheLayoutAsNonDeterministic(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	state := editor_state_model.NewDefaultEditorStateModel()
	configuration := namedConfiguration()
	configuration.Topology = config.TopologyRandom
	fixture.arrangeGeneration(state, configuration, newTemplateWithLayout(nil, nil))
	fixture.manualEditDelta.On("Capture", mock.Anything).Return(editor_state_model.ManualEditSettings{})

	// Act
	_, _ = fixture.handler.CaptureManualEdits(dtos.ManualEditCaptureRequestDto{State: toDto(state)})

	// Assert
	fixture.manualEditDelta.AssertCalled(t, "Capture", models.ManualEditCaptureRequest{
		Configuration:    configuration,
		NonDeterministic: true,
	})
}

func TestWhenCaptureConfigurationUsesRandomPortals_FlagsTheLayoutAsNonDeterministic(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	state := editor_state_model.NewDefaultEditorStateModel()
	configuration := namedConfiguration()
	configuration.Topology = config.TopologySquare
	configuration.RandomPortals = true
	fixture.arrangeGeneration(state, configuration, newTemplateWithLayout(nil, nil))
	fixture.manualEditDelta.On("Capture", mock.Anything).Return(editor_state_model.ManualEditSettings{})

	// Act
	_, _ = fixture.handler.CaptureManualEdits(dtos.ManualEditCaptureRequestDto{State: toDto(state)})

	// Assert
	fixture.manualEditDelta.AssertCalled(t, "Capture", models.ManualEditCaptureRequest{
		Configuration:    configuration,
		NonDeterministic: true,
	})
}

func TestWhenTheGeneratedTemplateHasNoVariant_CapturePassesAnEmptyFreshLayout(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	state := editor_state_model.NewDefaultEditorStateModel()
	configuration := namedConfiguration()
	fixture.arrangeGeneration(state, configuration, &template_model.Template{})
	fixture.manualEditDelta.On("Capture", mock.Anything).Return(editor_state_model.ManualEditSettings{})

	// Act
	_, _ = fixture.handler.CaptureManualEdits(dtos.ManualEditCaptureRequestDto{State: toDto(state)})

	// Assert
	fixture.manualEditDelta.AssertCalled(t, "Capture", models.ManualEditCaptureRequest{
		Configuration:    configuration,
		NonDeterministic: configuration.HasNonDeterministicLayout(),
	})
}

func TestWhenCaptureSucceeds_ReturnsTheDeltaServiceResult(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	state := editor_state_model.NewDefaultEditorStateModel()
	captured := editor_state_model.ManualEditSettingsFromSnapshot([]template_model.Zone{{Name: gofakeit.Word()}}, nil)
	fixture.arrangeGeneration(state, namedConfiguration(), newTemplateWithLayout(nil, nil))
	fixture.manualEditDelta.On("Capture", mock.Anything).Return(captured)

	// Act
	result, _ := fixture.handler.CaptureManualEdits(dtos.ManualEditCaptureRequestDto{State: toDto(state)})

	// Assert
	assert.Equal(t, dtos.ManualEditCaptureDto{ManualEdits: captured}, result)
}

func TestWhenCaptureSucceeds_ReturnsNoError(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	state := editor_state_model.NewDefaultEditorStateModel()
	fixture.arrangeGeneration(state, namedConfiguration(), newTemplateWithLayout(nil, nil))
	fixture.manualEditDelta.On("Capture", mock.Anything).Return(editor_state_model.ManualEditSettings{})

	// Act
	_, err := fixture.handler.CaptureManualEdits(dtos.ManualEditCaptureRequestDto{State: toDto(state)})

	// Assert
	assert.NoError(t, err)
}
