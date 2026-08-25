package templateHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWhenTemplateNameIsEmpty_ReturnsNoTemplateNameError(t *testing.T) {
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
	_, err := fixture.handler.GenerateTemplate(fixture.editorStateMapper.ToDto(state))

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrNoTemplateName)
}

func TestWhenGenerationYieldsNoTemplate_ReturnsGeneratedTemplateInvalidError(t *testing.T) {
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
	_, err := fixture.handler.GenerateTemplate(fixture.editorStateMapper.ToDto(state))

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrGeneratedTemplateInvalid)
}

func TestWhenTemplateIsGenerated_ReturnsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	state := editor_state_model.NewDefaultEditorStateModel()
	expected := &entities.RmgTemplate{Name: gofakeit.Word()}
	fixture.stateHandler.On("ValidateEditorState", state, true).
		Return(editor_state_dto.EditorStateValidationDto{State: state})
	fixture.mapper.On("FromEditorState", state).Return(namedConfiguration())
	fixture.templateGenerator.On("SetConfiguration", mock.Anything).Return()
	fixture.templateGenerator.On("Generate").Return(expected, nil)

	// Act
	loadDto, err := fixture.handler.GenerateTemplate(fixture.editorStateMapper.ToDto(state))

	// Assert
	require.NoError(t, err)
	assert.Same(t, expected, loadDto.Template)
}

func TestWhenValidationAndGenerationBothWarn_ConcatenatesTheWarnings(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	state := editor_state_model.NewDefaultEditorStateModel()
	validationWarning := gofakeit.Sentence(3)
	generationWarning := gofakeit.Sentence(3)
	fixture.stateHandler.On("ValidateEditorState", state, true).
		Return(editor_state_dto.EditorStateValidationDto{State: state, Warnings: []string{validationWarning}})
	fixture.mapper.On("FromEditorState", state).Return(namedConfiguration())
	fixture.templateGenerator.On("SetConfiguration", mock.Anything).Return()
	fixture.templateGenerator.On("Generate").
		Return(&entities.RmgTemplate{}, []string{generationWarning})

	// Act
	loadDto, _ := fixture.handler.GenerateTemplate(fixture.editorStateMapper.ToDto(state))

	// Assert
	assert.Equal(t, []string{validationWarning, generationWarning}, loadDto.Warnings)
}

func TestWhenValidationFixesTheState_MapsTheFixedState(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	state := editor_state_model.NewDefaultEditorStateModel()
	fixedState := editor_state_model.NewDefaultEditorStateModel()
	fixedState.PlayerCount = state.PlayerCount + 1
	fixture.stateHandler.On("ValidateEditorState", state, true).
		Return(editor_state_dto.EditorStateValidationDto{State: fixedState})
	fixture.mapper.On("FromEditorState", fixedState).Return(namedConfiguration())
	fixture.templateGenerator.On("SetConfiguration", mock.Anything).Return()
	fixture.templateGenerator.On("Generate").Return(&entities.RmgTemplate{}, nil)

	// Act
	_, _ = fixture.handler.GenerateTemplate(fixture.editorStateMapper.ToDto(state))

	// Assert
	fixture.mapper.AssertCalled(t, "FromEditorState", fixedState)
}

func TestWhenStateIsMapped_ConfiguresTheGeneratorWithTheMappedConfiguration(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	state := editor_state_model.NewDefaultEditorStateModel()
	configuration := namedConfiguration()
	fixture.stateHandler.On("ValidateEditorState", state, true).
		Return(editor_state_dto.EditorStateValidationDto{State: state})
	fixture.mapper.On("FromEditorState", state).Return(configuration)
	fixture.templateGenerator.On("SetConfiguration", configuration).Return()
	fixture.templateGenerator.On("Generate").Return(&entities.RmgTemplate{}, nil)

	// Act
	_, _ = fixture.handler.GenerateTemplate(fixture.editorStateMapper.ToDto(state))

	// Assert
	fixture.templateGenerator.AssertCalled(t, "SetConfiguration", configuration)
}

// namedConfiguration returns a configuration that passes the empty-name guard.
func namedConfiguration() *config.GeneratorConfig {
	configuration := config.NewGeneratorConfig()
	configuration.TemplateName = gofakeit.Word()
	return configuration
}
