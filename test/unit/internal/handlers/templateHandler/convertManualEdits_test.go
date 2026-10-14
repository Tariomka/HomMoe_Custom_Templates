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

func TestWhenConvertValidationRejectsTheBaseline_ReturnsTheRejection(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	baseline := editor_state_model.NewDefaultEditorStateModel()
	rejection := errors.New(gofakeit.Sentence(3))
	fixture.stateHandler.On("ValidateEditorState", baseline, true).
		Return(editor_state_dto.EditorStateValidationDto{State: baseline, Rejection: rejection})

	// Act
	_, err := fixture.handler.ConvertManualEdits(dtos.ManualEditConvertRequestDto{Baseline: toDto(baseline)})

	// Assert
	assert.ErrorIs(t, err, rejection)
}

func TestWhenConvertTemplateNameIsEmpty_ReturnsNoTemplateNameError(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	baseline := editor_state_model.NewDefaultEditorStateModel()
	configuration := config.NewGeneratorConfig()
	configuration.TemplateName = ""
	fixture.stateHandler.On("ValidateEditorState", baseline, true).
		Return(editor_state_dto.EditorStateValidationDto{State: baseline})
	fixture.mapper.On("FromEditorState", baseline).Return(configuration)

	// Act
	_, err := fixture.handler.ConvertManualEdits(dtos.ManualEditConvertRequestDto{Baseline: toDto(baseline)})

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrNoTemplateName)
}

func TestWhenConvertGenerationFails_ReturnsTheGenerationError(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	baseline := editor_state_model.NewDefaultEditorStateModel()
	generationError := errors.New(gofakeit.Sentence(3))
	fixture.stateHandler.On("ValidateEditorState", baseline, true).
		Return(editor_state_dto.EditorStateValidationDto{State: baseline})
	fixture.mapper.On("FromEditorState", baseline).Return(namedConfiguration())
	fixture.templateGenerator.On("SetConfiguration", mock.Anything).Return()
	fixture.templateGenerator.On("Generate").Return(nil, nil, generationError)

	// Act
	_, err := fixture.handler.ConvertManualEdits(dtos.ManualEditConvertRequestDto{Baseline: toDto(baseline)})

	// Assert
	assert.ErrorIs(t, err, generationError)
}

func TestWhenConvertGenerationYieldsNoTemplate_ReturnsGeneratedTemplateInvalidError(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	baseline := editor_state_model.NewDefaultEditorStateModel()
	fixture.stateHandler.On("ValidateEditorState", baseline, true).
		Return(editor_state_dto.EditorStateValidationDto{State: baseline})
	fixture.mapper.On("FromEditorState", baseline).Return(namedConfiguration())
	fixture.templateGenerator.On("SetConfiguration", mock.Anything).Return()
	fixture.templateGenerator.On("Generate").Return(nil, nil)

	// Act
	_, err := fixture.handler.ConvertManualEdits(dtos.ManualEditConvertRequestDto{Baseline: toDto(baseline)})

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrGeneratedTemplateInvalid)
}

func TestWhenConvertSucceeds_DecodesTheSnapshotAndPassesItAsTheEditedLayout(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	baseline := editor_state_model.NewDefaultEditorStateModel()
	configuration := namedConfiguration()
	configuration.Topology = config.TopologySquare
	editedZones := []template_model.Zone{{Name: "Edited", Size: gofakeit.Float64Range(1, 3)}}
	editedConnections := []template_model.Connection{{Name: "Edited-Link", From: "Edited", To: "Other", Length: 2}}
	freshZones := []template_model.Zone{{Name: "Fresh"}}
	freshConnections := []template_model.Connection{{Name: "Fresh-Link"}}
	fixture.arrangeGeneration(baseline, configuration, newTemplateWithLayout(freshZones, freshConnections))
	fixture.manualEditDelta.On("Capture", mock.Anything).Return(editor_state_model.ManualEditSettings{})

	// Act
	_, _ = fixture.handler.ConvertManualEdits(dtos.ManualEditConvertRequestDto{
		Baseline:    toDto(baseline),
		ManualEdits: editor_state_model.ManualEditSettingsFromSnapshot(editedZones, editedConnections),
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

func TestWhenConvertSucceeds_GeneratesFromTheBaselineSettings(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	baseline := editor_state_model.NewDefaultEditorStateModel()
	baseline.PlayerCount = gofakeit.IntRange(3, 6)
	fixture.arrangeGeneration(baseline, namedConfiguration(), newTemplateWithLayout(nil, nil))
	fixture.manualEditDelta.On("Capture", mock.Anything).Return(editor_state_model.ManualEditSettings{})

	// Act
	_, _ = fixture.handler.ConvertManualEdits(dtos.ManualEditConvertRequestDto{Baseline: toDto(baseline)})

	// Assert
	fixture.stateHandler.AssertCalled(t, "ValidateEditorState", baseline, true)
}

func TestWhenConvertSucceeds_MapsTheValidatedBaseline(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	baseline := editor_state_model.NewDefaultEditorStateModel()
	fixture.arrangeGeneration(baseline, namedConfiguration(), newTemplateWithLayout(nil, nil))
	fixture.manualEditDelta.On("Capture", mock.Anything).Return(editor_state_model.ManualEditSettings{})

	// Act
	_, _ = fixture.handler.ConvertManualEdits(dtos.ManualEditConvertRequestDto{Baseline: toDto(baseline)})

	// Assert
	fixture.mapper.AssertCalled(t, "FromEditorState", baseline)
}

func TestWhenConvertConfigurationHasARandomTopology_FlagsTheLayoutAsNonDeterministic(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	baseline := editor_state_model.NewDefaultEditorStateModel()
	configuration := namedConfiguration()
	configuration.Topology = config.TopologyRandom
	fixture.arrangeGeneration(baseline, configuration, newTemplateWithLayout(nil, nil))
	fixture.manualEditDelta.On("Capture", mock.Anything).Return(editor_state_model.ManualEditSettings{})

	// Act
	_, _ = fixture.handler.ConvertManualEdits(dtos.ManualEditConvertRequestDto{Baseline: toDto(baseline)})

	// Assert
	fixture.manualEditDelta.AssertCalled(t, "Capture", models.ManualEditCaptureRequest{
		Configuration:    configuration,
		NonDeterministic: true,
	})
}

func TestWhenConvertSucceeds_ReturnsTheDeltaServiceResult(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	baseline := editor_state_model.NewDefaultEditorStateModel()
	captured := editor_state_model.ManualEditSettingsFromSnapshot([]template_model.Zone{{Name: gofakeit.Word()}}, nil)
	fixture.arrangeGeneration(baseline, namedConfiguration(), newTemplateWithLayout(nil, nil))
	fixture.manualEditDelta.On("Capture", mock.Anything).Return(captured)

	// Act
	result, _ := fixture.handler.ConvertManualEdits(dtos.ManualEditConvertRequestDto{Baseline: toDto(baseline)})

	// Assert
	assert.Equal(t, dtos.ManualEditCaptureDto{ManualEdits: captured}, result)
}

func TestWhenConvertFails_NeverCaptures(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	baseline := editor_state_model.NewDefaultEditorStateModel()
	fixture.stateHandler.On("ValidateEditorState", baseline, true).
		Return(editor_state_dto.EditorStateValidationDto{State: baseline, Rejection: errors.New(gofakeit.Sentence(3))})

	// Act
	_, _ = fixture.handler.ConvertManualEdits(dtos.ManualEditConvertRequestDto{Baseline: toDto(baseline)})

	// Assert
	fixture.manualEditDelta.AssertNotCalled(t, "Capture", mock.Anything)
}
